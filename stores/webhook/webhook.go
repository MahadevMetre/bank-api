package webhook

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"bitbucket.org/paydoh/paydoh-commons/database"
	commonSrv "bitbucket.org/paydoh/paydoh-commons/services"
	"github.com/redis/go-redis/v9"

	"bankapi/constants"
	"bankapi/models"
	"bankapi/requests"
	"bankapi/security"
	"bankapi/services"
	"bankapi/utils"
)

type Webhook interface {
	GetWebhookData(request *requests.WebhookRequest) (interface{}, error)
}

type WebhookStore struct {
	db                  *sql.DB
	memory              *database.InMemory
	notificationService *services.NotificationService
	bankService         *services.BankApiService
	LoggerService       *commonSrv.LoggerService
}

func NewWebhookStore(log *commonSrv.LoggerService, db *sql.DB, memory *database.InMemory) *WebhookStore {
	notificationService := services.NewNotificationService()
	bankService := services.NewBankApiService(log, memory)
	return &WebhookStore{
		db:                  db,
		memory:              memory,
		notificationService: notificationService,
		bankService:         bankService,
		LoggerService:       log,
	}
}

func (store *WebhookStore) UpdateKycUpdateData(ctx context.Context, request *requests.KycDataUpdateRequest) error {
	logData := &commonSrv.LogEntry{
		Action:     constants.WEBHOOK,
		RequestURI: "/api/webhook/kvb/kyc",
		Message:    "UpdateKycUpdateData log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	existing, err := models.FindOneKycUpdateData(store.db, request.UserId)

	if err != nil {
		logData.Message = "UpdateKycUpdateData: Error finding KYC update data"
		store.LoggerService.LogError(logData)

		if errors.Is(err, constants.ErrNoDataFound) {
			logData.Message = "UpdateKycUpdateData: KYC update data not found, inserting new data"
			store.LoggerService.LogError(logData)

			insertingKycUpdateData := models.NewKycUpdateData()
			insertingKycUpdateData.Bind(request)

			if err := models.InsertKycUpdateData(store.db, insertingKycUpdateData); err != nil {
				logData.Message = "UpdateKycUpdateData: Error inserting KYC update data"
				store.LoggerService.LogError(logData)
				return err
			}

			return nil
		}

		return err
	}

	if _, err := models.FindOneAndUpdateKycUpdateData(store.db, existing.UserId, request.Status, request.Acom, request.AStat); err != nil {
		logData.Message = "UpdateKycUpdateData: Error updating KYC update data"
		store.LoggerService.LogError(logData)
		return err
	}

	logData.Message = "UpdateKycUpdateData: KYC update data updated successfully"
	logData.EndTime = time.Now()
	store.LoggerService.LogInfo(logData)

	return nil
}

func (store *WebhookStore) GetWebhookData(ctx context.Context, request *requests.WebhookRequest) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.WEBHOOK,
		RequestURI: "/api/webhook/route-mobile",
		Message:    "GetWebhookData log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	referencingUser := models.NewUserDeviceData()
	existingRouteMobileData := models.NewRouteMobileData()
	err := store.db.QueryRow(
		`
			SELECT id, mobile_number, signing_key, user_id, device_id, sim_vendor_id, device_token, os, package_id, created_at, updated_at
			FROM user_device_data
			WHERE mobile_number = $1
		`,
		request.Sender,
	).Scan(
		&referencingUser.Id,
		&referencingUser.MobileNumber,
		&referencingUser.SigningKey,
		&referencingUser.UserId,
		&referencingUser.DeviceId,
		&referencingUser.SimVendorId,
		&referencingUser.DeviceToken,
		&referencingUser.OS,
		&referencingUser.PackageId,
		&referencingUser.CreatedAt,
		&referencingUser.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logData.Message = "GetWebhookData: No device registered with the given number"
			store.LoggerService.LogError(logData)
			return nil, err
		}
		logData.Message = "GetWebhookData: Error querying device data"
		store.LoggerService.LogError(logData)
		return nil, err
	}

	users := make([]requests.NotificationUser, 0)

	notificationUser := requests.NewNotificationUser()
	notificationUser.UserId = referencingUser.UserId
	notificationUser.DeviceToken = referencingUser.DeviceToken
	notificationUser.PackageId = referencingUser.PackageId
	notificationUser.OS = referencingUser.OS

	users = append(users, *notificationUser)

	decryptedSigningKey, err := security.Decrypt(referencingUser.SigningKey, []byte(constants.AesPassPhrase))

	if err != nil {
		logData.Message = "GetWebhookData: Error decrypting signing key"
		store.LoggerService.LogError(logData)
		return nil, err
	}

	savedMessage, err := store.memory.Get(referencingUser.UserId)

	if err != nil {
		if errors.Is(err, redis.Nil) {
			logData.Message = "GetWebhookData: No message with the given sender exists"
			store.LoggerService.LogError(logData)
			return nil, err
		}
		logData.Message = "GetWebhookData: Error getting message from memory"
		store.LoggerService.LogError(logData)
		return nil, err
	}

	decryptedMessage, err := security.Decrypt(savedMessage, []byte(decryptedSigningKey))

	if err != nil {
		logData.Message = "GetWebhookData: Error decrypting message"
		store.LoggerService.LogError(logData)
		return nil, err
	}

	decryptedRequestMessage, err := security.Decrypt(request.Message, []byte(decryptedSigningKey))

	if err != nil {
		logData.Message = "GetWebhookData: Error decrypting request message"
		store.LoggerService.LogError(logData)
		return nil, err
	}

	if decryptedRequestMessage != decryptedMessage {
		logData.Message = "GetWebhookData: Saved message and received message do not match"
		store.LoggerService.LogError(logData)
		return nil, err
	}

	err = store.db.QueryRow(
		`
			SELECT id, user_id, mobile_number, operator, circle, created_at, updated_at
			FROM route_mobile_data
			WHERE user_id = $1	
		`,
		referencingUser.UserId,
	).Scan(
		&existingRouteMobileData.Id,
		&existingRouteMobileData.UserId,
		&existingRouteMobileData.MobileNumber,
		&existingRouteMobileData.Operator,
		&existingRouteMobileData.Circle,
		&existingRouteMobileData.CreatedAt,
		&existingRouteMobileData.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			insertedRouteMobileData := models.NewRouteMobileData()

			err = store.db.QueryRow(
				`
					INSERT into route_mobile_data (user_id, mobile_number, operator, circle)
					VALUES($1, $2, $3, $4) RETURNING id, user_id, mobile_number, operator, circle, created_at, updated_at
				`,
				referencingUser.UserId,
				referencingUser.MobileNumber,
				request.Operator,
				request.Circle,
			).Scan(
				&insertedRouteMobileData.Id,
				&insertedRouteMobileData.UserId,
				&insertedRouteMobileData.MobileNumber,
				&insertedRouteMobileData.Operator,
				&insertedRouteMobileData.Circle,
				&insertedRouteMobileData.CreatedAt,
				&insertedRouteMobileData.UpdatedAt,
			)

			if err != nil {
				logData.Message = "GetWebhookData: Error inserting route mobile data"
				store.LoggerService.LogError(logData)
				return nil, err
			}

			if err := store.memory.Delete(referencingUser.UserId); err != nil {
				logData.Message = "GetWebhookData: Error deleting route mobile data from memory"
				store.LoggerService.LogError(logData)
				return nil, err
			}

			notification := requests.NewNotificationRequest()

			if err := notification.CreateNotificationPayload(users, "Verification successful", fmt.Sprintf("Successfully verified user %s", referencingUser.UserId), "background", "success"); err != nil {
				logData.Message = "GetWebhookData: Error creating notification payload"
				store.LoggerService.LogError(logData)
				return nil, err
			}

			_, err := json.Marshal(notification)

			if err != nil {
				logData.Message = "GetWebhookData: Error marshaling notification data"
				store.LoggerService.LogError(logData)
				return nil, err
			}

			if _, err := store.notificationService.SendNotification(notification); err != nil {
				logData.Message = "GetWebhookData: Error sending notification"
				store.LoggerService.LogError(logData)
				return nil, err
			}

			verificationRequest := requests.NewOutgoingSimVerificationRequest()

			verificationRequest.Bind(referencingUser.UserId, referencingUser.MobileNumber, "Verified")

			if err := verificationRequest.Validate(); err != nil {
				logData.Message = "GetWebhookData: Error validating verification request"
				store.LoggerService.LogError(logData)
				return nil, err
			}

			if _, err := store.bankService.VerifySim(ctx, verificationRequest); err != nil {
				logData.Message = "GetWebhookData: Error verifying SIM"
				store.LoggerService.LogError(logData)
				return nil, err
			}

			return map[string]interface{}{
				"message": "Successfully verified user",
			}, nil
		}

		logData.Message = "GetWebhookData: Error fetching route mobile data"
		store.LoggerService.LogError(logData)
		return nil, err

	}

	if request.Operator != "" && existingRouteMobileData.Operator != "" && existingRouteMobileData.Operator != request.Operator {
		logData.Message = "GetWebhookData: Notify user of SIM change"
		store.LoggerService.LogError(logData)

		if err := store.memory.Delete(referencingUser.UserId); err != nil {
			logData.Message = "GetWebhookData: Error deleting route mobile data from memory"
			store.LoggerService.LogError(logData)
			return nil, err
		}

		notification := requests.NewNotificationRequest()

		if err := notification.CreateNotificationPayload(users, "Verification failure", fmt.Sprintf("failed to verify user %s", referencingUser.UserId), "background", "failure"); err != nil {
			logData.Message = "GetWebhookData: Error creating notification payload"
			store.LoggerService.LogError(logData)
			return nil, err
		}

		_, err := json.Marshal(notification)

		if err != nil {
			logData.Message = "GetWebhookData: Error marshaling notification data"
			store.LoggerService.LogError(logData)
			return nil, err
		}

		if _, err := store.notificationService.SendNotification(notification); err != nil {
			logData.Message = "GetWebhookData: Error sending notification"
			store.LoggerService.LogError(logData)
			return nil, err
		}

		verificationRequest := requests.NewOutgoingSimVerificationRequest()

		verificationRequest.Bind(referencingUser.UserId, referencingUser.MobileNumber, "Verified")

		if err := verificationRequest.Validate(); err != nil {
			logData.Message = "GetWebhookData: Error validating verification request"
			store.LoggerService.LogError(logData)
			return nil, err
		}

		if _, err := store.bankService.VerifySim(ctx, verificationRequest); err != nil {
			logData.Message = "GetWebhookData: Error verifying SIM"
			store.LoggerService.LogError(logData)
			return nil, err
		}

		return map[string]interface{}{
			"message": "Successfully verified user",
		}, nil
	}

	if err := store.memory.Delete(referencingUser.UserId); err != nil {
		logData.Message = "GetWebhookData: Error deleting route mobile data from memory"
		store.LoggerService.LogError(logData)
		return nil, err
	}

	notification := requests.NewNotificationRequest()

	if err := notification.CreateNotificationPayload(users, "Verification successful", fmt.Sprintf("Successfully verified user %s", referencingUser.UserId), "background", "success"); err != nil {
		logData.Message = "GetWebhookData: Error creating notification payload"
		store.LoggerService.LogError(logData)
		return nil, err
	}

	_, err = json.Marshal(notification)

	if err != nil {
		logData.Message = "GetWebhookData: Error marshaling notification data"
		store.LoggerService.LogError(logData)
		return nil, err
	}

	if _, err := store.notificationService.SendNotification(notification); err != nil {
		logData.Message = "GetWebhookData: Error sending notification"
		store.LoggerService.LogError(logData)
		return nil, err
	}

	verificationRequest := requests.NewOutgoingSimVerificationRequest()

	verificationRequest.Bind(referencingUser.UserId, referencingUser.MobileNumber, "Verified")

	if err := verificationRequest.Validate(); err != nil {
		logData.Message = "GetWebhookData: Error validating verification request"
		store.LoggerService.LogError(logData)
		return nil, err
	}

	if _, err := store.bankService.VerifySim(ctx, verificationRequest); err != nil {
		logData.Message = "GetWebhookData: Error verifying SIM"
		store.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "GetWebhookData: Verification successful"
	logData.EndTime = time.Now()
	store.LoggerService.LogInfo(logData)

	return map[string]interface{}{
		"message": "Successfully verified user",
	}, nil
}

func (store *WebhookStore) UpdateVcipData(ctx context.Context, request *requests.IncomingVcipData) error {
	logData := &commonSrv.LogEntry{
		Action:     constants.WEBHOOK,
		RequestURI: "/api/webhook/kvb/vcip",
		Message:    "UpdateVcipData log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	existing, err := models.GetVcipDataByUserId(store.db, request.ApplicantId)

	if err != nil {
		if err == constants.ErrNoDataFound {

			logData.Message = "UpdateVcipData: No existing VCIP data found, inserting new data"
			store.LoggerService.LogError(logData)

			insertingVcipData := models.NewVcipData()
			insertingVcipData.Bind(request, request.ApplicantId)

			if err := models.InsertVcipData(store.db, insertingVcipData); err != nil {
				logData.Message = "UpdateVcipData: Error inserting new VCIP data"
				store.LoggerService.LogError(logData)
				return err
			}

			return nil
		}

		logData.Message = "UpdateVcipData: Error fetching existing VCIP data"
		store.LoggerService.LogError(logData)
		return err
	}

	updateModel := models.NewVcipData()

	if request.PanNumber != existing.PanNumber {
		logData.Message = "UpdateVcipData: PAN number updated"
		store.LoggerService.LogError(logData)
		updateModel.PanNumber = request.PanNumber
	}

	if request.AadharReferenceNumber != existing.AadharReferenceNumber {
		logData.Message = "UpdateVcipData: Aadhar number updated"
		store.LoggerService.LogError(logData)
		updateModel.AadharReferenceNumber = request.AadharReferenceNumber
	}

	if request.VKYCCompletion != existing.VKYCCompletion {
		logData.Message = "UpdateVcipData: VKYC status updated"
		store.LoggerService.LogError(logData)
		updateModel.VKYCCompletion = request.VKYCCompletion
	}

	if request.VKYCAuditStatus != existing.VKYCAuditStatus {
		logData.Message = "UpdateVcipData: VKYC audit status updated"
		store.LoggerService.LogError(logData)
		updateModel.VKYCAuditStatus = request.VKYCAuditStatus
	}

	if request.AuditorRejectRemarks != existing.AuditorRejectRemarks {
		logData.Message = "UpdateVcipData: Auditor reject remarks updated"
		store.LoggerService.LogError(logData)
		updateModel.AuditorRejectRemarks = request.AuditorRejectRemarks
	}

	if err := models.UpdateVcipData(store.db, updateModel, existing.UserId); err != nil {
		logData.Message = "UpdateVcipData: Error updating VCIP data"
		store.LoggerService.LogError(logData)
		return err
	}

	logData.Message = "UpdateVcipData: Successfully updated VCIP data"
	logData.EndTime = time.Now()
	store.LoggerService.LogInfo(logData)
	return nil
}

func (s *WebhookStore) ProvideRewardsPoint(ctx context.Context, request *requests.RewardsRequest) error {
	logData := &commonSrv.LogEntry{
		Action:     constants.WEBHOOK,
		RequestURI: "/api/webhook/rewards-point",
		Message:    "ProvideRewardsPoint log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	userData, err := models.GetUserAndAccountDetailByUserID(s.db, request.BenfUserId)
	if err != nil {
		logData.Message = "Error fetching user account information"
		s.LoggerService.LogError(logData)
		return err
	}

	outgoingRewardsRequest := requests.NewOutgoingRewardsTransactionRequest()

	beneficiaryName := fmt.Sprintf("%s %s %s", userData.FirstName, userData.MiddleName, userData.LastName)

	amount := fmt.Sprintf("%f", request.Amount)

	if err := outgoingRewardsRequest.BindAndValidate(
		userData.Applicant_id,
		amount,
		request.Remarks,
		userData.AccountNumber,
		userData.MobileNumber,
		beneficiaryName,
		request.TranId,
	); err != nil {
		logData.Message = "Error binding and validating outgoing rewards request"
		s.LoggerService.LogError(logData)
		return err
	}

	response, err := s.bankService.RewardsTransfer(ctx, outgoingRewardsRequest)
	if err != nil {
		logData.Message = "Error during rewards transfer"
		s.LoggerService.LogError(logData)
		return err
	}

	if response.ErrorCode != "00" && response.ErrorCode != "0" {
		logData.Message = fmt.Sprintf("Rewards transfer failed: %s", response.ErrorMessage)
		s.LoggerService.LogError(logData)
		return errors.New(response.ErrorMessage)
	}

	byteudd, err := json.Marshal(response)
	if err != nil {
		logData.Message = "ProvideRewardsPoint: Error marshaling to JSON"
		s.LoggerService.LogError(logData)
		return err
	}

	logData.Message = "Rewards transfer successful"
	logData.ResponseSize = len(byteudd)
	logData.ResponseBody = string(byteudd)
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return nil
}
