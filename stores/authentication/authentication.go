package authentication

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"bitbucket.org/paydoh/paydoh-commons/database"
	commonSrv "bitbucket.org/paydoh/paydoh-commons/services"

	"bankapi/constants"
	"bankapi/models"
	"bankapi/requests"
	"bankapi/security"
	"bankapi/services"
	"bankapi/utils"
)

type Authentication interface {
	InitiateSimVerification(mobileNumber, key string, request *requests.AuthenticationRequest) (interface{}, error)
}

// AuthenticationStore @impl Authentication
type AuthenticationStore struct {
	db            *sql.DB
	m             *database.Document
	memory        *database.InMemory
	LoggerService *commonSrv.LoggerService
	auditLogSrv   services.AuditLogService
}

func NewAuthenticationStore(log *commonSrv.LoggerService, db *sql.DB, m *database.Document, memory *database.InMemory, auditSrv services.AuditLogService) *AuthenticationStore {
	return &AuthenticationStore{
		db:            db,
		m:             m,
		memory:        memory,
		LoggerService: log,
		auditLogSrv:   auditSrv,
	}
}

func (store *AuthenticationStore) InitiateSimVerification(ctx context.Context, userId, deviceIp, os, osVersion, key string, request *requests.AuthenticationRequest) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.AUTHENTICATION,
		RequestURI: "/api/authentication/initiate-sim-verification",
		Message:    "InitiateSimVerification log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	dataRequest, err := request.DecrypToData(key)
	if err != nil {
		logData.Message = "InitiateSimVerification: Error while decrypting data"
		store.LoggerService.LogError(logData)
		return nil, err
	}

	existingDevice, err := models.FindOneDeviceByUserID(store.db, userId)
	if err != nil {
		if err == constants.ErrDeviceNotFound {
			insertingDevice := models.NewDeviceData()

			insertingDevice.UserId = userId

			encryptedDeviceid, err := security.Encrypt([]byte(dataRequest.DeviceId), []byte(key))
			if err != nil {
				logData.Message = "InitiateSimVerification: Error while encrypting device id"
				store.LoggerService.LogError(logData)
				return nil, err
			}

			insertingDevice.DeviceId = encryptedDeviceid
			insertingDevice.DeviceIp = sql.NullString{
				String: deviceIp,
				Valid:  deviceIp != "",
			}

			if dataRequest.SimVendorId != "" {
				encryptedSimVendorId, err := security.Encrypt([]byte(dataRequest.SimVendorId), []byte(key))

				if err != nil {
					logData.Message = "InitiateSimVerification: Error while encrypting sim vendor id"
					store.LoggerService.LogError(logData)
					return nil, err
				}

				insertingDevice.SimVendorId = sql.NullString{
					String: encryptedSimVendorId,
					Valid:  encryptedSimVendorId != "",
				}
			}

			insertingDevice.PackageId = dataRequest.PackageId
			insertingDevice.OS = sql.NullString{
				String: os,
				Valid:  os != "",
			}
			insertingDevice.OSVersion = sql.NullString{
				String: osVersion,
				Valid:  osVersion != "",
			}
			insertingDevice.DeviceToken = sql.NullString{
				String: dataRequest.DeviceToken,
				Valid:  dataRequest.DeviceToken != "",
			}

			if err := models.InsertDevice(store.db, insertingDevice); err != nil {
				logData.Message = "InitiateSimVerification: Error while inserting device data"
				store.LoggerService.LogError(logData)
				return nil, err
			}

			message := security.GenerateRandomCode(36)

			encryptedMessage, err := security.Encrypt([]byte(message), []byte(key))

			if err != nil {
				logData.Message = "InitiateSimVerification: Error while encrypting message"
				store.LoggerService.LogError(logData)
				return nil, err
			}

			err = store.memory.Set(insertingDevice.UserId, encryptedMessage, time.Minute*10)

			if err != nil {
				logData.Message = "InitiateSimVerification: Error while setting message in memory"
				store.LoggerService.LogError(logData)
				return nil, err
			}

			return map[string]interface{}{
				"message":       encryptedMessage,
				"mobile_number": constants.ROUTE_SMS_MOBILE_NUMBER,
			}, nil
		}

		logData.Message = "InitiateSimVerification: Error while querying device data"
		store.LoggerService.LogError(logData)
		return nil, err
	}

	// It's not needed as we discuss, while updating it should not update it's device data in db

	// updateDevice := models.NewDeviceData()

	// encryptedDeviceid, err := security.Encrypt([]byte(dataRequest.DeviceId), []byte(key))
	// if err != nil {
	// 	logData.Message = "InitiateSimVerification: Error while encrypting device id"
	// 	store.LoggerService.LogError(logData)
	// 	return nil, err
	// }

	// updateDevice.DeviceId = encryptedDeviceid
	// updateDevice.DeviceToken.String = dataRequest.DeviceToken

	// if err := models.UpdateDevice(store.db, updateDevice, existingDevice.UserId); err != nil {
	// 	logData.Message = "InitiateSimVerification: Error while updating device data"
	// 	store.LoggerService.LogError(logData)
	// 	return nil, err
	// }

	if err := models.UpdateUserOnboardingStatusV2(existingDevice.UserId, false); err != nil {
		logData.Message = "SmsVerification: Error updating user onboarding status:- " + err.Error()
		store.LoggerService.LogError(logData)
		return nil, err
	}

	if err := models.UpdateIsActiveAndSimVerifiedStatus(true, false, existingDevice.UserId); err != nil {
		logData.Message = "InitiateSimVerification: Error while updating active status"
		store.LoggerService.LogError(logData)
		return nil, err
	}

	message := security.GenerateRandomCode(36)

	encryptedMessage, err := security.Encrypt([]byte(message), []byte(key))
	if err != nil {
		logData.Message = "InitiateSimVerification: Error while encrypting message"
		store.LoggerService.LogError(logData)
		return nil, err
	}

	err = store.memory.Set(existingDevice.UserId, encryptedMessage, time.Minute*10)
	if err != nil {
		logData.Message = "InitiateSimVerification: Error while setting message in memory"
		store.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "InitiateSimVerification: Verification process completed successfully"
	logData.ResponseSize = len(encryptedMessage)
	logData.ResponseBody = string(encryptedMessage)
	logData.EndTime = time.Now()
	store.LoggerService.LogInfo(logData)

	return map[string]interface{}{
		"message":       encryptedMessage,
		"mobile_number": constants.ROUTE_SMS_MOBILE_NUMBER,
	}, nil
}

func (s *AuthenticationStore) GetSimVerificationStatus(ctx context.Context, userId string) (interface{}, error) {
	deviceData, err := models.FindOneDeviceByUserID(s.db, userId)
	if err != nil {
		s.LoggerService.LogInfo(&commonSrv.LogEntry{
			Message:   "GetSimVerificationStatus error:- " + err.Error(),
			UserID:    utils.GetUserIDFromContext(ctx),
			RequestID: utils.GetRequestIDFromContext(ctx),
		})
		return nil, err
	}

	resp := map[string]bool{
		"is_sim_verified": deviceData.IsSimVerified,
	}

	return resp, nil
}

func (s *AuthenticationStore) GetCurrentUser(ctx context.Context, userId string) (*models.UserDevicePersonal, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.AUTHENTICATION,
		RequestURI: "/api/authorization/authenticated/current-user",
		Message:    "Fetching current user data",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	data, err := models.FindOneUserPersonal(s.db, userId)

	if err != nil {
		logData.Message = "GetCurrentUser: Error while fetching user data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "GetCurrentUser: User data fetched successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return data, nil
}

func (s *AuthenticationStore) GetCurrentUserByMobileNumber(ctx context.Context, mobileNumber string) (*models.UserDevicePersonal, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.AUTHENTICATION,
		RequestURI: "/api/authorization/authenticated/current-user-mobile",
		Message:    "Fetching current user data by mobile number",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	data, err := models.FindOneUserPersonalMobileNumber(s.db, mobileNumber)

	if err != nil {
		logData.Message = "GetCurrentUserByMobileNumber: Error while fetching user data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "GetCurrentUserByMobileNumber: User data fetched successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return data, nil
}

func (s *AuthenticationStore) Logout(ctx context.Context, userId string) error {
	logoutEndpoint := "/api/authentication/logout"
	logData := &commonSrv.LogEntry{
		Action:     constants.AUTHENTICATION,
		RequestURI: logoutEndpoint,
		Message:    "Logging out user",
		UserID:     userId,
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	err := s.memory.Delete(fmt.Sprintf(constants.TokenKeyFormat, userId))
	if err != nil {
		logData.Message = "Logout: Error deleting user data from memory"
		s.LoggerService.LogError(logData)
		return err
	}

	user, err := models.GetUserDataByUserId(s.db, userId)
	if err != nil {
		logData.Message = "Logout: Error getting user data by user id"
		s.LoggerService.LogError(logData)
	}

	// saving audit log
	if err := s.auditLogSrv.Save(ctx, &services.AuditLog{
		UserID:         userId,
		ApplicantID:    user.ApplicantId,
		RequestURL:     logoutEndpoint,
		HTTPMethod:     "POST",
		ResponseStatus: 200,
		Action:         constants.LOGOUT,
	}); err != nil {
		logData.Message = "Logout: Error while saving audit log"
		s.LoggerService.LogError(logData)
	}

	logData.Message = "Logout: User logged out successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return nil
}
