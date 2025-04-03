package upi

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/paydoh/paydoh-commons/pkg/json"

	"bitbucket.org/paydoh/paydoh-commons/database"
	commonSrv "bitbucket.org/paydoh/paydoh-commons/services"
	"bitbucket.org/paydoh/paydoh-commons/types"

	"bankapi/constants"
	"bankapi/models"
	"bankapi/requests"
	"bankapi/responses"
	"bankapi/security"
	"bankapi/services"
	"bankapi/utils"
)

type Upi interface {
	SimBindingAndSmsVerification(authValues *models.AuthValues) (interface{}, error)
	CreateUpiID(authValues *models.AuthValues, request *requests.CreateUPIRequest) (interface{}, error)
	SetUPIPin(authValues *models.AuthValues, request *requests.SetUpiPinRequest) (interface{}, error)
	CheckAccountBalance(authValues *models.AuthValues, request *requests.ReqBalEnqRequest) (interface{}, error)
	ValidateVPA(authValues *models.AuthValues, request *requests.ValidateVpaRequest) (interface{}, error)
	ProcessPaymentWithVPA(authValues *models.AuthValues, request *requests.PayMoneyWithVpaRequest) (interface{}, error)
	LinkedAccountlist(authValues *models.AuthValues) (interface{}, error)
	PayeeNameGet(authValues *models.AuthValues, request *requests.GetAllBankAccount) (interface{}, error)
	AadharRequestListAccount(authValues *models.AuthValues, request *requests.AadharReqlistaccount) (interface{}, error)
	GetTransactionsHistoryList(authValues *models.AuthValues) (interface{}, error)
	GetUpiTokenXml(authValues *models.AuthValues, request *requests.UpiTokenRequest) (interface{}, error)
	GetUpiToken(authValues *models.AuthValues, challenge string) (interface{}, error)
	CollectUpiMoneyCount(authValues *models.AuthValues) (interface{}, error)
	CollectUpiMoney(authValues *models.AuthValues, request *requests.UpiCollectMoneyRequest) (interface{}, error)
}

type Store struct {
	db            *sql.DB
	m             *database.Document
	redis         *database.InMemory
	bankService   *services.BankApiService
	LoggerService *commonSrv.LoggerService
	auditLogSrv   services.AuditLogService
}

func NewStore(log *commonSrv.LoggerService, db *sql.DB, m *database.Document, memory *database.InMemory, redis *database.InMemory, auditLogSrv services.AuditLogService) *Store {
	bankService := services.NewBankApiService(log, memory)
	return &Store{
		db:            db,
		m:             m,
		redis:         redis,
		bankService:   bankService,
		LoggerService: log,
		auditLogSrv:   auditLogSrv,
	}
}

func (s *Store) SimBindingAndSmsVerification(ctx context.Context, authValues *models.AuthValues, request *requests.SimBindingRequest) (interface{}, error) {
	// upiData, err := models.GetAccountDataByUserId(s.db, authValues.UserId)
	// if err != nil {
	// 	return nil, err
	// }

	// if upiData.UpiId.String != "" {
	// 	return nil, errors.New("the UPI ID already exists for this user; therefore, this process has already been completed")
	// }

	requestUri := "/api/upi/simbinding/sms-verification"

	if request.Type == "" {
		request.Type = "n"
	}

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:     constants.UPI,
		RequestURI: requestUri,
		Message:    "SimBinding And SmsVerification log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	// data, err := models.FindOneTransIdByUserId(s.db, authValues.UserId)
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		logData.Message = "SimBindingAndSmsVerification: TransactionId not found for the user proceeding further"
	// 		s.LoggerService.LogError(logData)
	// 	} else {
	// 		logData.Message = "SimBindingAndSmsVerification: Error Transaction Id " + err.Error()
	// 		s.LoggerService.LogError(logData)
	// 		return nil, err
	// 	}
	// }

	// if data != nil {
	// 	if data.LoginRefId.String == "" {
	// 		_, err := models.DeleteClientIDByUserId(authValues.UserId)
	// 		if err != nil {
	// 			logData.Message = "SimBindingAndSmsVerification: Error deleting Client Data device data" + err.Error()
	// 			s.LoggerService.LogError(logData)
	// 			return nil, err
	// 		}
	// 	}
	// }

	existingUserData, err := models.GetSingleDeviceDataByUserIDV2(authValues.UserId)
	if err != nil {
		logData.Message = "SimBindingAndSmsVerification: Error getting user data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	var decryptedDeviceId string
	if request.Type == "y" {
		decryptedDeviceId = utils.GetDeviceIdFromContext(ctx)
		if decryptedDeviceId == "" {
			logData.Message = "SimBindingAndSmsVerification: Device ID not found in context"
			s.LoggerService.LogError(logData)
			return nil, errors.New("device ID not found in context")
		}
	} else {
		decryptedDeviceId, err = security.Decrypt(existingUserData.DeviceId, []byte(authValues.Key))
		if err != nil {
			logData.Message = "SimBindingAndSmsVerification: Error decrypting device ID"
			s.LoggerService.LogError(logData)
			return nil, err
		}
	}

	mobileMapping0Request := requests.NewOutgoingMobileMappingType0ApiRequest()
	if err := mobileMapping0Request.Bind(decryptedDeviceId, authValues.DeviceIp); err != nil {
		logData.Message = "SimBindingAndSmsVerification: Error binding in mobileMapping0Request "
		s.LoggerService.LogError(logData)
		return nil, err
	}

	mobileMapping0Response, err := s.bankService.MobileMapping(ctx, mobileMapping0Request)
	if err != nil {
		logData.Message = "SimBindingAndSmsVerification: Error calling MobileMapping type 0 bank service to Create UPI ID"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if mobileMapping0Response.Response.ResponseCode != "0" {
		logData.Message = "SimBindingAndSmsVerification: Error in mobile mapping type 0 response"
		s.LoggerService.LogError(logData)
		return nil, errors.New(mobileMapping0Response.Response.ResponseMessage)
	}

	UuidClientId, err := s.GenerateClientId(authValues.OS)
	if err != nil {
		logData.Message = "SimBindingAndSmsVerification: Error generating client_id"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if request.Type == "n" {
		_, err := models.SaveTransIdAndClientIdByUserId(s.db, authValues.UserId, mobileMapping0Response.Response.TransID, UuidClientId)
		if err != nil {
			logData.Message = "SimBindingAndSmsVerification: Error saving client ID and Trans ID"
			s.LoggerService.LogError(logData)
			return nil, err
		}
	} else if request.Type == "y" {
		_, err := s.SaveTransIdAndClientIdToRedis(ctx, authValues.UserId, mobileMapping0Response.Response.TransID, UuidClientId)
		if err != nil {
			logData.Message = "SimBindingAndSmsVerification: Error saving client ID and Trans ID"
			s.LoggerService.LogError(logData)
			return nil, err
		}
	}

	mobileMapping0Response.Response.ClientId = UuidClientId

	responseBytes, err := mobileMapping0Response.Marshal()

	if err != nil {
		return nil, err
	}

	verifyUserRequest := requests.NewOutgoingVerifyUserApiRequest()
	if err := verifyUserRequest.Bind(authValues.DeviceIp, decryptedDeviceId, authValues.OS, mobileMapping0Response.Response.TransID, UuidClientId); err != nil {
		logData.Message = "CreateUpiID: Error binding in verifyUserRequest"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	verifyUserResponse, err := s.bankService.VerifyUpiService(ctx, verifyUserRequest)
	if err != nil {
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if verifyUserResponse.Response.ResponseCode != "0" {
		logData.Message = "CreateUpiID: Error verifying UPI service"
		s.LoggerService.LogError(logData)
		return nil, errors.New(verifyUserResponse.Response.ResponseMessage)
	}

	encrypted, err := security.Encrypt(responseBytes, []byte(authValues.Key))
	if err != nil {
		return nil, err
	}
	body, err := json.Marshal(mobileMapping0Request)
	if err != nil {
		logData.Message = "MobileMapping Type 0: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	encryptedReq, err := security.Encrypt([]byte(string(body)), []byte(authValues.Key))

	if err != nil {
		logData.Message = "CreateUpiID: Error encrypting response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if err := s.auditLogSrv.Save(ctx, &services.AuditLog{
		UserID:         authValues.UserId,
		RequestURL:     requestUri,
		HTTPMethod:     "POST",
		ResponseStatus: 200,
		Action:         constants.SIM_BINDING_SMS_VERIFICATION,
		RequestBody:    encryptedReq,
	}); err != nil {
		logData.Message = "error while saving audit log" + err.Error()
		logData.EndTime = time.Now()
		s.LoggerService.LogInfo(logData)
	}

	logData.Message = "SimBindingAndSmsVerification: Response encrypted successfully"
	logData.ResponseSize = len(responseBytes)
	logData.EndTime = time.Now()
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(responseBytes)
	s.LoggerService.LogInfo(logData)

	return encrypted, nil
}

func (s *Store) SaveTransIdAndClientIdToRedis(ctx context.Context, userID string, transactionID, clientID string) (map[string]string, error) {

	redisKey := fmt.Sprintf(constants.UpiSimBindingKey, userID)

	data := map[string]string{
		"transaction_id": transactionID,
		"client_id":      clientID,
	}

	err := s.redis.GetClient().HSet(ctx, redisKey, data).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to save client data to Redis: %w", err)
	}

	err = s.redis.GetClient().Expire(ctx, redisKey, constants.UpiSimBindingTTL).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to set expiration time for Redis key: %w", err)
	}

	savedData, err := s.redis.GetClient().HGetAll(ctx, redisKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve client data from Redis for verification: %w", err)
	}

	return savedData, nil
}

func (s *Store) waitForLongSMSCode() error {
	data := os.Getenv("LONG_SMS_WAIT_TIME")
	if data == "" {
		return fmt.Errorf("LONG_SMS_WAIT_TIME environment variable not set or empty")
	}

	waitTime, err := strconv.Atoi(data)
	if err != nil {
		return fmt.Errorf("failed to parse LONG_SMS_WAIT_TIME: %v", err)
	}

	time.Sleep(time.Second * time.Duration(waitTime))
	return nil
}

func (s *Store) CreateUpiID(ctx context.Context, authValues *models.AuthValues, request *requests.CreateUPIRequest) (interface{}, error) {
	upiData, err := models.GetAccountDataByUserId(s.db, authValues.UserId)
	if err != nil {
		return nil, err
	}

	if upiData.UpiId.String != "" {
		return nil, errors.New("UPI ID already exists")
	}

	// wait for long-sms code verification
	if err := s.waitForLongSMSCode(); err != nil {
		return nil, err
	}

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		UserID:        authValues.UserId,
		DeviceIP:      authValues.DeviceIp,
		LatLong:       authValues.LatLong,
		StartTime:     startTime,
		RequestID:     utils.GetRequestIDFromContext(ctx),
		RequestMethod: "POST",
		RequestURI:    "/api/upi/create-upi-id",
		Message:       "Create Upi ID log",
	}

	existingUserData, err := models.UpdateDeviceData(s.db, authValues.UserId, authValues.DeviceIp, authValues.OS, authValues.OSVersion, authValues.LatLong)
	if err != nil {
		logData.Message = "CreateUpiID: Error updating device data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	personalInformation, err := models.GetPersonalInformation(s.db, existingUserData.UserId)
	if err != nil {
		logData.Message = "CreateUpiID: Error fetching personal information"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	decryptedDeviceId, err := security.Decrypt(existingUserData.DeviceId, []byte(authValues.Key))
	if err != nil {
		logData.Message = "CreateUpiID: Error decrypting device id"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	MobileMappingType0Transid, err := models.FindOneTransIdByUserId(s.db, existingUserData.UserId)
	if err != nil {
		logData.Message = "CreateUpiID: Error fetching mobilemapping type0 trans id"
		s.LoggerService.LogError(logData)
		return "", err
	}

	UserClientID, err := models.FindOneClientIDByUserId(s.db, existingUserData.UserId)
	if err != nil {
		logData.Message = "CreateUpiID: Error fetching user client id"
		s.LoggerService.LogError(logData)
		return "", err
	}

	mobileMapping1Request := requests.NewOutgoingMobileMappingType1ApiRequest()
	if err := mobileMapping1Request.Bind(decryptedDeviceId, existingUserData.DeviceIp, existingUserData.OS, MobileMappingType0Transid.TransId.String, UserClientID.ClientId.String); err != nil {
		logData.Message = "CreateUpiID: Error binding in mobileMapping1Request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	mobileMapping1Response, err := s.bankService.MobileMapping1(ctx, mobileMapping1Request)
	if err != nil {
		logData.Message = "CreateUpiID: Error calling MobileMapping type 1 bank service to Create UPI ID"
		s.LoggerService.LogError(logData)
		// delete client id data
		models.DeleteClientIDByUserId(authValues.UserId)
		return nil, err
	}

	if mobileMapping1Response.Response.ResponseCode == "1" {
		logData.Message = "CreateUpiID: Error in MobileMapping type 1"
		s.LoggerService.LogError(logData)
		// delete client id data
		models.DeleteClientIDByUserId(authValues.UserId)
		return nil, errors.New("sms not received")
	}

	if mobileMapping1Response.Response.ResponseCode == "4" || mobileMapping1Response.Response.ResponseCode == "5" {
		remappingRequest := requests.NewOutgoingRemappingApiRequest()
		if err := remappingRequest.Bind(
			existingUserData.MobileNumber,
			decryptedDeviceId,
			existingUserData.DeviceIp,
			existingUserData.OSVersion,
			existingUserData.OS,
			personalInformation.Email,
			UserClientID.ClientId.String,
		); err != nil {
			logData.Message = "CreateUpiID: Error binding in remappingRequest"
			s.LoggerService.LogError(logData)
			return nil, err
		}

		remappingResponse, err := s.bankService.ReMapping(ctx, remappingRequest)
		if err != nil {
			logData.Message = "CreateUpiID: Error calling ReMapping bank service to Create UPI ID"
			s.LoggerService.LogError(logData)
			// delete client id data
			models.DeleteClientIDByUserId(authValues.UserId)
			return nil, err
		}

		// responseCode not getting in remapping api response
		if remappingResponse.Response.ResponseCode != "0" {
			logData.Message = "CreateUpiID: Error in remapping"
			s.LoggerService.LogError(logData)
			// delete client id data
			models.DeleteClientIDByUserId(authValues.UserId)
			return nil, errors.New(remappingResponse.Response.ResponseMessage)
		}

		if _, err := models.SaveServerIdByUserId(s.db, existingUserData.UserId, remappingResponse.Response.ServerID, remappingResponse.Response.LoginRefID, MobileMappingType0Transid.TransId.String, UserClientID.ClientId.String); err != nil {
			logData.Message = "CreateUpiID: Error saving server ID"
			s.LoggerService.LogError(logData)
			return nil, err
		}

		existingUserReqListKeysRequest := requests.NewOutgoingExistingReqlistkeysApiRequest()
		if err := existingUserReqListKeysRequest.Bind(existingUserData.DeviceIp, decryptedDeviceId, existingUserData.PackageId, existingUserData.MobileNumber, request.Challenge); err != nil {
			logData.Message = "CreateUpiID: Error binding in existingUserReqListKeysRequest ist time"
			s.LoggerService.LogError(logData)
			return nil, err
		}

		existingUserReqListKeysResponse, err := s.bankService.ExistingUserRequestListKeys(ctx, existingUserReqListKeysRequest)
		if err != nil {
			logData.Message = "CreateUpiID: Error calling ExistingUserRequestListKeys 1st time bank service to Create UPI ID"
			// delete client id data
			models.DeleteClientIDByUserId(authValues.UserId)
			s.LoggerService.LogError(logData)
			return nil, err
		}

		if existingUserReqListKeysResponse.Response.ResponseCode != "0" {
			logData.Message = "CreateUpiID: Error in existing user request list keys"
			s.LoggerService.LogError(logData)
			// delete client id data
			models.DeleteClientIDByUserId(authValues.UserId)
			return nil, errors.New(existingUserReqListKeysResponse.Response.ResponseMessage)
		}

		// // check upi id is already created
		// accountData, err := models.GetAccountDataByUserId(s.db, authValues.UserId)
		// if err != nil {
		// 	logData.Message = "CreateUpiID: upi already created"
		// 	s.LoggerService.LogError(logData)
		// 	return nil, err
		// }

		// responseMap := map[string]map[string]interface{}{
		// 	"Response": {
		// 		"ResponseCode":    "0",
		// 		"ResponseMessage": "Success",
		// 		"upi_id":          accountData.UpiId,
		// 	},
		// }

		// responseMapBytes, err := json.Marshal(responseMap)
		// if err != nil {
		// 	return nil, err
		// }

		// encrypted, err := security.Encrypt(responseMapBytes, []byte(authValues.Key))
		// if err != nil {
		// 	logData.Message = "CreateUpiID: Error calling ProfileCreation bank service to Create UPI ID"
		// 	s.LoggerService.LogError(logData)
		// 	return nil, err
		// }

		// return encrypted, nil

	} else if mobileMapping1Response.Response.ResponseCode == "0" {

		UserClientID, err := models.FindOneClientIDByUserId(s.db, existingUserData.UserId)
		if err != nil {
			return nil, err
		}

		if UserClientID.LoginRefId.String != "" {
			loginRefID := UserClientID.LoginRefId.String
			logData.Message = "CreateUPIID: Reusing existing loginRefID"
			logData.ResponseBody = loginRefID
			s.LoggerService.LogInfo(logData)

		} else {

			lcValidatorRequest := requests.NewOutgoingLCValidatorApiRequest()
			lcValidatorRequest.Bind(
				"91"+existingUserData.MobileNumber,
				UserClientID.ClientId.String, //MobileMappingType0Transid.TransId.String,
				decryptedDeviceId,
				existingUserData.DeviceIp,
				personalInformation.FirstName,
			)

			lcValidatorResponse, err := s.bankService.LcValidator(ctx, lcValidatorRequest)
			if err != nil {
				logData.Message = "CreateUpiID: Error calling LC validator bank service to Create UPI ID"
				s.LoggerService.LogError(logData)
				// delete client id data
				models.DeleteClientIDByUserId(authValues.UserId)
				return nil, err
			}

			if lcValidatorResponse.Response.ResponseCode != "0" {
				logData.Message = "CreateUpiID: Error in LC validator"
				s.LoggerService.LogError(logData)
				// delete client id data
				models.DeleteClientIDByUserId(authValues.UserId)
				return nil, errors.New(lcValidatorResponse.Response.ResponseMessage)
			}

			if _, err := models.SaveServerIdByUserId(s.db, existingUserData.UserId, lcValidatorResponse.Response.ServerID, lcValidatorResponse.Response.LoginRefID, MobileMappingType0Transid.TransId.String, UserClientID.ClientId.String); err != nil {
				logData.Message = "CreateUpiID: Error saving server ID"
				s.LoggerService.LogError(logData)
				return nil, err
			}
		}

		// demographic api call
		demographicRequest := requests.NewOutGoingDemographicRequest()
		if err := demographicRequest.Bind(existingUserData.ApplicantId, existingUserData.MobileNumber); err != nil {
			logData.Message = "ProcessPaymentWithVPA: Error binding in demographic request"
			s.LoggerService.LogError(logData)
			return nil, err
		}

		demographicResponse, demographicErr := s.bankService.GetDemographicData(ctx, demographicRequest)
		if demographicErr != nil {
			// delete client id data
			models.DeleteClientIDByUserId(authValues.UserId)

			bankErr := s.bankService.HandleBankSpecificError(demographicErr, func(errorCode string) (string, bool) {
				return constants.GetDemographicErrorMessage(errorCode)
			})

			if bankErr != nil {
				logData.Message = fmt.Sprintf("GetDemographicData: Bank error encountered (ErrorCode: %s)", bankErr.ErrorCode)
				s.LoggerService.LogError(logData)
				return nil, errors.New(bankErr.ErrorMessage)
			}
			return nil, demographicErr
		}

		addressCo, addressVtc, addressState := "NA", "NA", "NA"
		if demographicResponse != nil {
			addressCo = demographicResponse.Root.UIDData.Poa.Co
			addressVtc = demographicResponse.Root.Vtc
			addressState = demographicResponse.Root.UIDData.Poa.State
		}

		addressReq := requests.NewAddressRequest()
		if err := addressReq.Bind(addressCo, addressVtc, addressState); err != nil {
			return nil, err
		}

		cryptoInfo, err := s.GenerateCryptoInfo(ctx, authValues)
		if err != nil {
			logData.Message = "CreateUpiID: Error generating CryptoInfo"
			s.LoggerService.LogError(logData)
			return nil, err
		}

		userProfileCreationRequest := requests.NewOutgoingProfileCreationApiRequest()
		if err := userProfileCreationRequest.Bind(
			existingUserData.MobileNumber,
			decryptedDeviceId,
			personalInformation.FirstName,
			personalInformation.Email,
			personalInformation.Gender,
			cryptoInfo,
			demographicResponse,
			addressReq,
			existingUserData.PackageId,
		); err != nil {
			logData.Message = "CreateUpiID: Error binding in userProfileCreationRequest"
			s.LoggerService.LogError(logData)
			return nil, err
		}

		profileCreationResponse, err := s.bankService.ProfileCreation(ctx, userProfileCreationRequest)
		if err != nil {
			logData.Message = "CreateUpiID: Error calling ProfileCreation bank service to Create UPI ID"
			s.LoggerService.LogError(logData)
			// delete client id data
			models.DeleteClientIDByUserId(authValues.UserId)
			return nil, err
		}

		if profileCreationResponse.Response.ResponseCode != "0" {
			logData.Message = "CreateUpiID: Error in profile creation"
			s.LoggerService.LogError(logData)
			// delete client id data
			models.DeleteClientIDByUserId(authValues.UserId)
			return nil, errors.New(profileCreationResponse.Response.ResponseMessage)
		}

		userReqListKeys := requests.NewOutgoingReqListKeysApiRequest()
		if err := userReqListKeys.NewBind(); err != nil {
			logData.Message = "CreateUpiID: Error binding in userReqListKeys"
			s.LoggerService.LogError(logData)
			return nil, err
		}

		_, err = s.bankService.AlreadyUserRequestListKeys(ctx, userReqListKeys)
		if err != nil {
			logData.Message = "CreateUpiID: Error calling AlreadyUserRequestListKeys bank service to Create UPI ID"
			s.LoggerService.LogError(logData)
			return nil, err
		}

		// if userReqListKeysResponse.Response.ResponseCode != "0" {
		// 	logData.Message = "CreateUpiID: Error in user request list keys"
		// 	s.LoggerService.LogError(logData)
		// 	return nil, errors.New(userReqListKeysResponse.Response.ResponseMessage)
		// }

	}

	cryptoInfo, err := s.GenerateCryptoInfo(ctx, authValues)
	if err != nil {
		logData.Message = "CreateUpiID: Error generating CryptoInfo"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	userListAccountRequest := requests.NewOutgoingCreateupiidRequestListAccountApiRequest()
	if err := userListAccountRequest.Bind(
		existingUserData.MobileNumber,
		existingUserData.DeviceIp,
		personalInformation.FirstName,
		cryptoInfo,
		existingUserData.PackageId,
		authValues.LatLong,
	); err != nil {
		logData.Message = "CreateUpiID: Error binding in userListAccountRequest"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	userListAccountResponse, err := s.bankService.CreateUpiIdRequestListAccounts(ctx, userListAccountRequest)
	if err != nil {
		logData.Message = "CreateUpiID: Error calling CreateUpiIdRequestListAccounts bank service to Create UPI ID"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if userListAccountResponse.Response.ResponseCode != "0" {
		logData.Message = "CreateUpiID: Error in creating UPI ID request list accounts"
		s.LoggerService.LogError(logData)
		return nil, errors.New(userListAccountResponse.Response.ResponseMessage)
	}

	userCheckPspRequest := requests.NewOutgoingPspAvailabilityApiRequest()
	if err := userCheckPspRequest.Bind(
		decryptedDeviceId,
		existingUserData.OSVersion,
		existingUserData.OS,
		cryptoInfo,
		userListAccountRequest,
	); err != nil {
		logData.Message = "CreateUpiID: Error binding in userCheckPspRequest"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	userCheckPspResponse, err := s.bankService.RequestPspAvailability(ctx, userCheckPspRequest)
	if err != nil {
		logData.Message = "CreateUpiID: Error calling RequestPspAvailability bank service to Create UPI ID"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if userCheckPspResponse.Response.ResponseCode != "0" {
		logData.Message = "CreateUpiID: Error checking PSP availability"
		s.LoggerService.LogError(logData)
		return nil, errors.New(userCheckPspResponse.Response.ResponseMessage)
	}

	userAddBankRequest := requests.NewOutgoingAddBankApiRequest()
	if err := userAddBankRequest.Bind(
		existingUserData.DeviceIp,
		decryptedDeviceId,
		existingUserData.MobileNumber,
		existingUserData.OSVersion,
		existingUserData.OS,
		personalInformation.FirstName,
		existingUserData.LatLong,
		cryptoInfo,
		request.Location,
		userListAccountResponse,
		userListAccountRequest,
		request.DeviceCapability,
		existingUserData.PackageId,
		personalInformation.LastName,
		upiData.AccountNumber,
	); err != nil {
		logData.Message = "CreateUpiID: Error binding in userAddBankRequest"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	userAddBankResponse, err := s.bankService.RequestAddBankAccount(ctx, userAddBankRequest)
	if err != nil {
		logData.Message = "CreateUpiID: Error calling Add Bank Account bank service to Create UPI ID"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if userAddBankResponse.Response.ResponseCode != "0" {
		logData.Message = "CreateUpiID: Error adding bank account"
		s.LoggerService.LogError(logData)
		return nil, errors.New(userAddBankResponse.Response.ResponseMessage)
	}

	userAddBankResponse.Response.UpiId = userListAccountRequest.ReqListAccount.Payeraddr

	// update upi id in db
	models.UpdateAccountByUserId(&models.AccountDataUpdate{UpiId: userAddBankResponse.Response.UpiId}, authValues.UserId)

	// existingUserReqListKeysRequest := requests.NewOutgoingExistingReqlistkeysApiRequest()
	// if err := existingUserReqListKeysRequest.Bind(existingUserData.DeviceIp, decryptedDeviceId, existingUserData.PackageId, existingUserData.MobileNumber, request.Challenge); err != nil {
	// 	logData.Message = "CreateUpiID: Error binding in existingUserReqListKeysRequest 2nd time"
	// 	s.LoggerService.LogError(logData)
	// 	return nil, err
	// }

	// existingUserReqListKeysResponse, err := s.bankService.ExistingUserRequestListKeys(existingUserReqListKeysRequest)
	// if err != nil {
	// 	logData.Message = "CreateUpiID: Error calling Existing User Request List Keys 2nd time bank service to Create UPI ID"
	// 	s.LoggerService.LogError(logData)
	// 	return nil, err
	// }

	// if existingUserReqListKeysResponse.Response.ResponseCode != "0" {
	// 	logData.Message = "CreateUpiID: Error in existing user request list keys after bank account addition"
	// 	s.LoggerService.LogError(logData)
	// 	return nil, errors.New(existingUserReqListKeysResponse.Response.ResponseMessage)
	// }

	// update onboarding status
	if err := models.UpdateUserOnboardingStatus(constants.UPI_GENERATION_STAGE, authValues.UserId); err != nil {
		logData.Message = "CreateUpiID: error while updating onboarding status"
		s.LoggerService.LogError(logData)
	}

	// // Store UPI token in Redis
	// xmlString := existingUserReqListKeysResponse.Response.Response.Ns2RespListKeys.KeyList.Key.KeyValue.Text
	// hexToken := s.convertUpiToken(xmlString)
	// s.SaveUPIToken(authValues.UserId, authValues.DeviceIp, hexToken)

	responseBytes, err := userAddBankResponse.Marshal()
	if err != nil {
		return nil, err
	}

	encrypted, err := security.Encrypt(responseBytes, []byte(authValues.Key))
	if err != nil {
		return nil, err
	}

	logData.Message = "CreateUpiID: Response encrypted successfully"
	logData.ResponseSize = len(responseBytes)
	logData.EndTime = time.Now()
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(responseBytes)
	s.LoggerService.LogInfo(logData)

	return encrypted, nil
}

func (s *Store) RemapExistingUpiID(ctx context.Context, authValues *models.AuthValues, request *requests.IncomingRemappingRequest) (interface{}, error) {
	// wait for long-sms code verification
	if err := s.waitForLongSMSCode(); err != nil {
		return nil, err
	}

	requestUri := "/api/upi/remapping-upi-id"
	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		UserID:        authValues.UserId,
		DeviceIP:      authValues.DeviceIp,
		LatLong:       authValues.LatLong,
		StartTime:     startTime,
		RequestID:     utils.GetRequestIDFromContext(ctx),
		RequestMethod: "POST",
		RequestURI:    requestUri,
		Message:       "RemapExistingUpiID log",
	}

	userData, err := models.GetUserDataByUserId(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "RemapExistingUpiID: Error fetching user"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	existingDeviceData, err := models.FindOneDeviceByUserID(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "RemapExistingUpiID: Error fetching user"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	personalInformation, err := models.GetPersonalInformation(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "RemapExistingUpiID: Error fetching personal information"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	// userClientId, err := models.FindOneClientIDByUserId(s.db, existingUserData.UserId)
	// if err != nil {
	// 	logData.Message = "RemapExistingUpiID: Error fetching user client ID"
	// 	s.LoggerService.LogError(logData)
	// 	return nil, err
	// }

	// mobileMappingtype0TransId, err := models.FindOneTransIdByUserId(s.db, existingUserData.UserId)
	// if err != nil {
	// 	logData.Message = "RemapExistingUpiID: Error fetching mobilemapping type0 trans id"
	// 	s.LoggerService.LogError(logData)
	// 	return "", err
	// }

	clientData, err := s.getClientDataFromRedis(ctx, authValues.UserId)
	if err != nil {
		logData.Message = "RemapExistingUpiID: Error getting client data from Redis"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	transationID := clientData["transaction_id"]
	clientID := clientData["client_id"]

	mobileMapping1Request := requests.NewOutgoingMobileMappingType1ApiRequest()
	if err := mobileMapping1Request.Bind(utils.GetDeviceIdFromContext(ctx), authValues.DeviceIp, authValues.OS, transationID, clientID); err != nil {
		logData.Message = "RemapExistingUpiID: Error binding in mobileMapping1Request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	mobileMapping1Response, err := s.bankService.MobileMapping1(ctx, mobileMapping1Request)
	if err != nil {
		logData.Message = "RemapExistingUpiID: Error calling MobileMapping type 1 bank service"
		s.LoggerService.LogError(logData)
		models.UpdateDevice(s.db, &models.DeviceData{IsActive: false}, authValues.UserId)
		return nil, err
	}

	if mobileMapping1Response.Response.ResponseCode == "1" {
		logData.Message = "RemapExistingUpiID: Error in MobileMapping type 1 - SMS not received"
		s.LoggerService.LogError(logData)
		return nil, errors.New("sms not received")
	}

	var remappingResponse *responses.ReMappingApiResponse
	if mobileMapping1Response.Response.ResponseCode == "4" || mobileMapping1Response.Response.ResponseCode == "5" {
		remappingRequest := requests.NewOutgoingRemappingApiRequest()
		if err := remappingRequest.Bind(
			userData.MobileNumber,
			utils.GetDeviceIdFromContext(ctx),
			authValues.DeviceIp,
			authValues.OSVersion,
			authValues.OS,
			personalInformation.Email,
			clientID,
		); err != nil {
			logData.Message = "RemapExistingUpiID: Error binding in remappingRequest"
			s.LoggerService.LogError(logData)
			models.UpdateDevice(s.db, &models.DeviceData{IsActive: false}, authValues.UserId)
			return nil, err
		}

		remappingResponse, err = s.bankService.ReMapping(ctx, remappingRequest)
		if err != nil {
			logData.Message = "RemapExistingUpiID: Error calling ReMapping bank service"
			s.LoggerService.LogError(logData)
			models.UpdateDevice(s.db, &models.DeviceData{IsActive: false}, authValues.UserId)
			return nil, err
		}

		// responseCode not getting in remapping api response
		if remappingResponse.Response.ResponseCode != "0" {
			logData.Message = "RemapExistingUpiID: Error in remapping api response"
			s.LoggerService.LogError(logData)
			return nil, errors.New(remappingResponse.Response.ResponseMessage)
		}
	}

	existingUserReqListKeysRequest := requests.NewOutgoingExistingReqlistkeysApiRequest()
	if err := existingUserReqListKeysRequest.Bind(
		authValues.DeviceIp,
		utils.GetDeviceIdFromContext(ctx),
		utils.GetPackageIdFromContext(ctx),
		userData.MobileNumber,
		request.Challenge,
	); err != nil {
		logData.Message = "RemapExistingUpiID: Error binding in ExistingUserReqListKeysRequest"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	existingUserReqListKeysResponse, err := s.bankService.ExistingUserRequestListKeys(ctx, existingUserReqListKeysRequest)
	if err != nil {
		logData.Message = "RemapExistingUpiID: Error calling ExistingUserRequestListKeys"
		s.LoggerService.LogError(logData)
		models.UpdateDevice(s.db, &models.DeviceData{IsActive: false}, authValues.UserId)
		return nil, err
	}

	if existingUserReqListKeysResponse.Response.ResponseCode != "0" {
		logData.Message = "RemapExistingUpiID: Error in ExistingUserReqListKeys response"
		s.LoggerService.LogError(logData)
		return nil, errors.New(existingUserReqListKeysResponse.Response.ResponseMessage)
	}

	// update device data
	updateDeviceData := &models.DeviceData{
		IsActive:  true,
		OS:        sql.NullString{String: authValues.OS, Valid: true},
		OSVersion: sql.NullString{String: authValues.OSVersion, Valid: true},
		DeviceIp:  sql.NullString{String: authValues.DeviceIp, Valid: true},
	}

	currentDeviceID := utils.GetDeviceIdFromContext(ctx)
	if currentDeviceID != "" {
		encryptedDeviceID, err := security.Encrypt([]byte(currentDeviceID), []byte(authValues.Key))
		if err != nil {
			logData.Message = "SetUPIPin: Error encrypting device ID"
			s.LoggerService.LogError(logData)
			return nil, err
		}

		updateDeviceData.DeviceId = string(encryptedDeviceID)
	}

	currentPackageID := utils.GetPackageIdFromContext(ctx)
	if currentPackageID != "" {
		updateDeviceData.PackageId = currentPackageID
	}

	if request.DeviceToken != "" {
		updateDeviceData.DeviceToken = sql.NullString{
			String: request.DeviceToken,
			Valid:  true,
		}
	}

	if err := models.UpdateDevice(s.db, updateDeviceData, authValues.UserId); err != nil {
		logData.Message = "RemapExistingUpiID: Error updating device"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if remappingResponse != nil {
		if _, err := models.SaveServerIdByUserId(s.db, existingDeviceData.UserId, remappingResponse.Response.ServerID, remappingResponse.Response.LoginRefID, transationID, clientID); err != nil {
			logData.Message = "RemapExistingUpiID: Error saving server ID"
			s.LoggerService.LogError(logData)
			return nil, err
		}
	}

	_, err = s.SaveServerIdAndLoginRefIdToRedis(ctx, existingDeviceData.UserId, remappingResponse.Response.ServerID, remappingResponse.Response.LoginRefID)
	if err != nil {
		logData.Message = "RemapExistingUpiID: Error saving data to Redis"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	responseBytes, err := existingUserReqListKeysResponse.Marshal()
	if err != nil {
		logData.Message = "RemapExistingUpiID: Error marshalling response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := json.Marshal(mobileMapping1Request)
	if err != nil {
		logData.Message = "MobileMapping Type 0: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	encryptedReq, err := security.Encrypt([]byte(string(body)), []byte(authValues.Key))

	if err != nil {
		logData.Message = "RemapExistingUpiID: Error encrypting response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if err := s.auditLogSrv.Save(ctx, &services.AuditLog{
		UserID:         authValues.UserId,
		RequestURL:     requestUri,
		HTTPMethod:     "POST",
		ResponseStatus: 200,
		Action:         constants.REMAPPING_UPI_ID,
		RequestBody:    encryptedReq,
	}); err != nil {
		logData.Message = "error while saving audit log" + err.Error()
		logData.EndTime = time.Now()
		s.LoggerService.LogInfo(logData)
	}

	logData.Message = "RemapExistingUpiID: Response encrypted successfully"
	logData.ResponseSize = len(responseBytes)
	logData.EndTime = time.Now()
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(responseBytes)
	s.LoggerService.LogInfo(logData)

	return nil, nil
}

func (s *Store) SaveServerIdAndLoginRefIdToRedis(ctx context.Context, userID string, serverID string, loginRefID string) (map[string]string, error) {
	redisKey := fmt.Sprintf(constants.UpiSimBindingKey, userID)

	data := map[string]string{
		"server_id":    serverID,
		"login_ref_id": loginRefID,
	}

	err := s.redis.GetClient().HSet(context.Background(), redisKey, data).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to save serverid and loginrefid data to Redis: %w", err)
	}

	err = s.redis.GetClient().Expire(context.Background(), redisKey, constants.UpiSimBindingTTL).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to set expiration time for Redis key: %w", err)
	}

	return data, nil
}

func (s *Store) getClientDataFromRedis(ctx context.Context, userID string) (map[string]string, error) {
	data, err := s.redis.GetClient().HGetAll(context.Background(), fmt.Sprintf(constants.UpiSimBindingKey, userID)).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve client data from Redis: %w", err)
	}

	if len(data) == 0 {
		return nil, errors.New("client data not found in redis")
	}

	return data, nil
}

func (s *Store) SetUPIPin(ctx context.Context, authValues *models.AuthValues, request *requests.SetUpiPinRequest) (interface{}, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		UserID:        authValues.UserId,
		DeviceIP:      authValues.DeviceIp,
		LatLong:       authValues.LatLong,
		StartTime:     startTime,
		RequestID:     utils.GetRequestIDFromContext(ctx),
		RequestMethod: "POST",
		RequestURI:    "/api/upi/set-upi-pin",
		Message:       "Set UPI Pin log",
	}

	existingUserData, err := models.FindOneDeviceByUserID(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "SetUPIPin: Error getting user data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	userData, err := models.GetUserDataByUserId(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "SetUPIPin: Error getting user data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	personalInformation, err := models.GetPersonalInformation(s.db, existingUserData.UserId)
	if err != nil {
		logData.Message = "SetUPIPin: Error fetching personal information"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	cryptoInfo, err := s.GenerateCryptoInfo(ctx, authValues)
	if err != nil {
		logData.Message = "SetUPIPin: Error generating CryptoInfo"
		s.LoggerService.LogError(logData)
		return nil, fmt.Errorf("error generating CryptoInfo: %v", err)
	}

	userListAccountRequest := requests.NewOutgoingCreateupiidRequestListAccountApiRequest()
	if err := userListAccountRequest.Bind(
		userData.MobileNumber,
		existingUserData.DeviceIp.String,
		personalInformation.FirstName,
		cryptoInfo,
		existingUserData.PackageId,
		authValues.LatLong,
	); err != nil {
		logData.Message = "SetUPIPin: Error binding create UPI ID request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	setUserUpiPinRegMobile := requests.NewOutgoingSetUpiPinReqRegMobApiRequest()
	if err := setUserUpiPinRegMobile.Bind(
		userData.MobileNumber,
		existingUserData.DeviceIp.String,
		authValues.LatLong,
		cryptoInfo,
		request.TransId,
		request.Otp,
		request.UpiPin,
		request.AtmPin,
		userListAccountRequest,
		request.Cred_AADHAAR,
	); err != nil {
		logData.Message = "SetUPIPin: Error binding set UPI PIN request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	setUserUpiPinRegMobileResponse, err := s.bankService.SetUpiPinReqRegMobile(ctx, setUserUpiPinRegMobile)
	if err != nil {
		logData.Message = "SetUPIPin: Error calling bank service to set UPI PIN"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if setUserUpiPinRegMobileResponse.Response.ResponseCode != "0" {
		logData.Message = "SetUPIPin: Received error code from bank service for set UPI PIN"
		s.LoggerService.LogError(logData)
		return nil, errors.New(setUserUpiPinRegMobileResponse.Response.ResponseMessage)
	}

	responseBytes, err := setUserUpiPinRegMobileResponse.Marshal()
	if err != nil {
		logData.Message = "SetUPIPin: Error marshaling response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	// update onboarding status
	if err := models.UpdateUserOnboardingStatus(constants.UPI_PIN_SETUP_STAGE, authValues.UserId); err != nil {
		logData.Message = "SetUPIPin: error while updating onboarding status"
		s.LoggerService.LogError(logData)
	}

	encrypted, err := security.Encrypt(responseBytes, []byte(authValues.Key))
	if err != nil {
		logData.Message = "SetUPIPin: Error encrypting response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "SetUPIPin: Response encrypted successfully"
	logData.ResponseSize = len(responseBytes)
	logData.EndTime = time.Now()
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(responseBytes)
	s.LoggerService.LogInfo(logData)

	return encrypted, nil
}

func (s *Store) CheckAccountBalance(ctx context.Context, authValues *models.AuthValues, request *requests.ReqBalEnqRequest) (interface{}, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		UserID:        authValues.UserId,
		DeviceIP:      authValues.DeviceIp,
		LatLong:       authValues.LatLong,
		StartTime:     startTime,
		RequestID:     utils.GetRequestIDFromContext(ctx),
		RequestMethod: "POST",
		RequestURI:    "/api/upi/account-balance",
		Message:       "Check account balance log",
	}

	existingUserData, err := models.FindOneDeviceByUserID(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "CheckAccountBalance: Error getting user data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	userData, err := models.GetUserDataByUserId(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "CheckAccountBalance: Error getting user data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	personalInformation, err := models.GetPersonalInformation(s.db, existingUserData.UserId)
	if err != nil {
		logData.Message = "CheckAccountBalance: Error fetching personal information"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "CheckAccountBalance: Personal information fetched successfully"
	s.LoggerService.LogInfo(logData)

	cryptoInfo, err := s.GenerateCryptoInfo(ctx, authValues)
	if err != nil {
		logData.Message = "CheckAccountBalance: Error generating CryptoInfo"
		s.LoggerService.LogError(logData)
		return nil, fmt.Errorf("error generating CryptoInfo: %v", err)
	}

	logData.Message = "CheckAccountBalance: CryptoInfo generated successfully"
	s.LoggerService.LogInfo(logData)

	userCreateUpiIdListAccountRequest := requests.NewOutgoingCreateupiidRequestListAccountApiRequest()
	if err := userCreateUpiIdListAccountRequest.Bind(
		userData.MobileNumber,
		existingUserData.DeviceIp.String,
		personalInformation.FirstName,
		cryptoInfo,
		existingUserData.PackageId,
		authValues.LatLong,
	); err != nil {
		logData.Message = "CheckAccountBalance: Error binding create UPI ID request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "CheckAccountBalance: Bind create UPI ID request successful"
	s.LoggerService.LogInfo(logData)

	userCheckBankBalanceRequest := requests.NewOutgoingReqBalEnqApiRequest()
	if err := userCheckBankBalanceRequest.Bind(
		existingUserData.DeviceIp.String,
		authValues.LatLong,
		cryptoInfo,
		request.TransId,
		request.UpiPin,
		userCreateUpiIdListAccountRequest,
	); err != nil {
		logData.Message = "CheckAccountBalance: Error binding check bank balance request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "CheckAccountBalance: Bind check bank balance request successful"
	s.LoggerService.LogInfo(logData)

	userCheckBankBalanceResponse, err := s.bankService.RequestCheckAccountBalance(ctx, userCheckBankBalanceRequest)
	if err != nil {
		logData.Message = "CheckAccountBalance: Error calling bank service to check account balance"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "CheckAccountBalance: Bank service called successfully to check account balance"
	s.LoggerService.LogInfo(logData)

	if userCheckBankBalanceResponse.Response.ResponseCode == "0" &&
		strings.EqualFold(userCheckBankBalanceResponse.Response.ResponseMessage, constants.UpiErrorMessageInvalidMpin) {
		logData.Message = "CheckAccountBalance: Invalid MPIN entered"
		s.LoggerService.LogError(logData)
		return nil, errors.New(constants.UpiInvalidMpinErrorMessage)
	}

	responseBytes, err := userCheckBankBalanceResponse.Marshal()
	if err != nil {
		logData.Message = "CheckAccountBalance: Error marshaling response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "CheckAccountBalance: Response marshaled successfully"
	s.LoggerService.LogInfo(logData)

	encrypted, err := security.Encrypt(responseBytes, []byte(authValues.Key))
	if err != nil {
		logData.Message = "CheckAccountBalance: Error encrypting response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "CheckAccountBalance: Response encrypted successfully"
	logData.ResponseSize = len(responseBytes)
	logData.EndTime = time.Now()
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(responseBytes)
	s.LoggerService.LogInfo(logData)

	return encrypted, nil
}

func (s *Store) ValidateVPA(ctx context.Context, authValues *models.AuthValues, request *requests.ValidateVpaRequest) (interface{}, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		UserID:        authValues.UserId,
		DeviceIP:      authValues.DeviceIp,
		LatLong:       authValues.LatLong,
		StartTime:     startTime,
		RequestID:     utils.GetRequestIDFromContext(ctx),
		RequestMethod: "POST",
		RequestURI:    "/api/pay-val-vpa",
		Message:       "Validating VPA log",
	}

	existingUserData, err := models.FindOneDeviceByUserID(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "ValidateVPA: Error getting user data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	userData, err := models.GetUserDataByUserId(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "ValidateVPA: Error getting user data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	personalInformation, err := models.GetPersonalInformation(s.db, existingUserData.UserId)
	if err != nil {
		logData.Message = "ValidateVPA: Error fetching personal information"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	cryptoInfo, err := s.GenerateCryptoInfo(ctx, authValues)
	if err != nil {
		logData.Message = "ValidateVPA: Error generating CryptoInfo"
		s.LoggerService.LogError(logData)
		return nil, fmt.Errorf("error generating CryptoInfo: %v", err)
	}

	userListAccountRequest := requests.NewOutgoingCreateupiidRequestListAccountApiRequest()
	if err := userListAccountRequest.Bind(
		userData.MobileNumber,
		existingUserData.DeviceIp.String,
		personalInformation.FirstName,
		cryptoInfo,
		existingUserData.PackageId,
		authValues.LatLong,
	); err != nil {
		logData.Message = "ValidateVPA: Error binding create UPI ID request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	upiData, err := models.GetAccountDataByUserId(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "ValidateVPA: error while getting account data"
		return nil, err
	}

	// get receiver upi handler
	if !strings.Contains(request.RPayeraddr, "@") {
		rUpiData, _ := models.GetUserAndAccountDetailByMobileNumber(s.db, request.RPayeraddr)
		if rUpiData != nil {
			request.RPayeraddr = rUpiData.UpiId.String
		}
	}

	userValidatevpa := requests.NewOutgoingReqValAddApiRequest()
	if err := userValidatevpa.Bind(
		existingUserData.DeviceIp.String,
		personalInformation.FirstName,
		cryptoInfo,
		userData.MobileNumber,
		userListAccountRequest,
		request,
		upiData.UpiId.String,
		authValues.LatLong,
	); err != nil {
		logData.Message = "ValidateVPA: Error binding validate VPA address request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	userValidatevpaResponse, err := s.bankService.ValidateVpaAddress(ctx, userValidatevpa)
	if err != nil {
		logData.Message = "ValidateVPA: Error calling bank service to validate VPA address"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "ValidateVPA: Bank service called successfully to validate VPA address"
	s.LoggerService.LogInfo(logData)

	if userValidatevpaResponse.Response.ResponseCode != "0" {
		logData.Message = "ValidateVPA: Received error code from bank service for validate VPA address"
		s.LoggerService.LogError(logData)
		return nil, errors.New(userValidatevpaResponse.Response.ResponseMessage)
	}

	responseBytes, err := userValidatevpaResponse.Marshal()
	if err != nil {
		logData.Message = "ValidateVPA: Error marshaling response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	encrypted, err := security.Encrypt(responseBytes, []byte(authValues.Key))
	if err != nil {
		logData.Message = "ValidateVPA: Error encrypting response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "ValidateVPA: Response encrypted successfully"
	logData.ResponseSize = len(responseBytes)
	logData.EndTime = time.Now()
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(responseBytes)
	s.LoggerService.LogInfo(logData)

	return encrypted, nil
}

func (s *Store) ProcessPaymentWithVPA(ctx context.Context, authValues *models.AuthValues, request *requests.PayMoneyWithVpaRequest) (interface{}, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		UserID:        authValues.UserId,
		DeviceIP:      authValues.DeviceIp,
		LatLong:       authValues.LatLong,
		StartTime:     startTime,
		RequestID:     utils.GetRequestIDFromContext(ctx),
		Latency:       time.Since(startTime).Seconds(),
		RequestMethod: "POST",
		RequestURI:    "/api/upi/pay-vpa",
		Message:       "payment with VPA log",
	}

	existingUserData, err := models.FindOneDeviceByUserID(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "ProcessPaymentWithVPA: Error getting user data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	userData, err := models.GetUserDataByUserId(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "ProcessPaymentWithVPA: Error getting user data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	personalInformation, err := models.GetPersonalInformation(s.db, existingUserData.UserId)
	if err != nil {
		logData.Message = "ProcessPaymentWithVPA: Error fetching personal information"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "ProcessPaymentWithVPA: Personal information fetched successfully"
	s.LoggerService.LogInfo(logData)

	cryptoInfo, err := s.GenerateCryptoInfo(ctx, authValues)
	if err != nil {
		logData.Message = "ProcessPaymentWithVPA: Error generating CryptoInfo"
		s.LoggerService.LogError(logData)
		return nil, fmt.Errorf("error generating CryptoInfo: %v", err)
	}

	upiData, err := models.GetAccountDataByUserId(s.db, authValues.UserId)
	if err != nil {
		return nil, err
	}

	// get receiver upi handler
	if !strings.Contains(request.Payeeaddr, "@") {
		rUpiData, _ := models.GetUserAndAccountDetailByMobileNumber(s.db, request.Payeeaddr)
		if rUpiData != nil {
			request.Payeeaddr = rUpiData.UpiId.String
		}
	}

	userPaymentwithVpaRequest := requests.NewOutgoingReqPayApiRequest()
	if err := userPaymentwithVpaRequest.Bind(
		userData.MobileNumber,
		existingUserData.DeviceIp.String,
		personalInformation.FirstName,
		authValues.LatLong,
		cryptoInfo,
		request,
		upiData.UpiId.String,
	); err != nil {
		logData.Message = "ProcessPaymentWithVPA: Error binding payment with VPA request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	// update transaction data in db
	if err := models.UpdateTransactionByTransID(s.db, &models.Transaction{
		TransactionID:   userPaymentwithVpaRequest.ReqPay.TxnID,
		Amount:          types.FromString(request.PayerAmount),
		PaymentMode:     models.PaymentModeUPI,
		TransactionDesc: types.FromString(request.Remark),
		UPIPayeeAddr:    types.FromString(request.Payeeaddr),
		UPIPayeeName:    types.FromString(request.PayeeName),
	}); err != nil {
		logData.Message = "ProcessPaymentWithVPA: Error inserting transaction details"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	userPaymentwithVpaResponse, err := s.bankService.PayWithVpa(ctx, userPaymentwithVpaRequest)
	if err != nil {
		// update transaction status in db
		if err := models.UpdateTransactionByTransID(s.db, &models.Transaction{
			TransactionID: userPaymentwithVpaRequest.ReqPay.TxnID,
			CBSStatus:     types.FromString("Failure"),
		}); err != nil {
			logData.Message = "PaymentCallback Update: Error updating transaction status in db"
			s.LoggerService.LogError(logData)
		}
		logData.Message = "ProcessPaymentWithVPA: Error calling bank service to process payment with VPA"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if userPaymentwithVpaResponse.Response.ResponseCode != "0" {
		// update transaction status in db
		if err := models.UpdateTransactionByTransID(s.db, &models.Transaction{
			TransactionID: userPaymentwithVpaRequest.ReqPay.TxnID,
			CBSStatus:     types.FromString("Failure"),
		}); err != nil {
			logData.Message = "PaymentCallback Update: Error updating transaction status in db"
			s.LoggerService.LogError(logData)
		}
		logData.Message = "ProcessPaymentWithVPA: Received error code from bank service for payment with VPA"
		s.LoggerService.LogError(logData)
		return nil, errors.New(userPaymentwithVpaResponse.Response.ResponseMessage)
	}

	// Extract necessary data from response
	resp := userPaymentwithVpaResponse.Response.Response.Ns2RespPay
	payeeInfo := resp.Resp.Ref[1]
	utr := resp.Txn.CustRef
	ifsc := payeeInfo.IFSC

	mobileNumber := request.Payeeaddr
	if strings.Contains(mobileNumber, "@") {
		parts := strings.Split(mobileNumber, ".")
		if len(parts) > 0 {
			mobileNumber = parts[0]
		}
	}

	accountNumber := payeeInfo.AcNum
	remarks := resp.Txn.Note
	amount := payeeInfo.OrgAmount
	payeeName := payeeInfo.RegName
	transactionTime := resp.Txn.Ts

	// Create the frontend response structure
	frontendResponse := &responses.UpiPaymentResponse{
		UTR:             utr,
		IFSC:            ifsc,
		MobileNumber:    mobileNumber,
		AccountNumber:   accountNumber,
		Remarks:         remarks,
		Amount:          amount,
		PayeeName:       payeeName,
		TransactionTime: transactionTime,
	}

	responseBytes, err := frontendResponse.Marshal()
	if err != nil {
		logData.Message = "ProcessPaymentWithVPA: Error marshaling response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	// update transaction status in db
	if err := models.UpdateTransactionByTransID(s.db, &models.Transaction{
		TransactionID: userPaymentwithVpaRequest.ReqPay.TxnID,
		UTRRefNumber:  types.FromString(utr),
		CBSStatus:     types.FromString("Success"),
	}); err != nil {
		logData.Message = "PaymentCallback Update: Error updating transaction status in db"
		s.LoggerService.LogError(logData)
	}

	logData.Message = "ProcessPaymentWithVPA: Response marshaled successfully"
	s.LoggerService.LogInfo(logData)

	encrypted, err := security.Encrypt(responseBytes, []byte(authValues.Key))
	if err != nil {
		logData.Message = "ProcessPaymentWithVPA: Error encrypting response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "ProcessPaymentWithVPA: Response encrypted successfully"
	logData.ResponseSize = len(responseBytes)
	logData.EndTime = time.Now()
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(responseBytes)
	s.LoggerService.LogInfo(logData)

	return encrypted, nil
}

func (s *Store) LinkedAccountlist(ctx context.Context, authValues *models.AuthValues) (interface{}, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		UserID:        authValues.UserId,
		DeviceIP:      authValues.DeviceIp,
		LatLong:       authValues.LatLong,
		StartTime:     startTime,
		RequestID:     utils.GetRequestIDFromContext(ctx),
		RequestMethod: "GET",
		RequestURI:    "/api/upi/fetch-account-list",
		Message:       "Fetch User Bank Accounts List log",
	}

	userData, err := models.GetUserDataByUserId(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "LinkedAccountlist: Error getting user data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "LinkedAccountlist: Device data updated successfully"
	s.LoggerService.LogInfo(logData)

	cryptoInfo, err := s.GenerateCryptoInfo(ctx, authValues)
	if err != nil {
		logData.Message = "LinkedAccountlist: Error generating CryptoInfo"
		s.LoggerService.LogError(logData)
		return nil, fmt.Errorf("error generating CryptoInfo: %v", err)
	}

	userFetchAllBankAccountsRequest := requests.NewOutgoingAccountLinkApiRequest()
	if err := userFetchAllBankAccountsRequest.Bind(
		userData.MobileNumber,
		cryptoInfo,
	); err != nil {
		logData.Message = "LinkedAccountlist: Error binding fetch all bank accounts request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "LinkedAccountlist: Bind fetch all bank accounts request successful"
	s.LoggerService.LogInfo(logData)

	userFetchAllBankAccountsResponse, err := s.bankService.LinkBankAccount(ctx, userFetchAllBankAccountsRequest)
	if err != nil {
		logData.Message = "LinkedAccountlist: Error calling bank service to fetch all bank accounts"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "LinkedAccountlist: Bank service called successfully to fetch all bank accounts"
	s.LoggerService.LogInfo(logData)

	if userFetchAllBankAccountsResponse.Response.ResponseCode != "0" {
		logData.Message = "LinkedAccountlist: Received error code from bank service"
		s.LoggerService.LogError(logData)
		return nil, errors.New(userFetchAllBankAccountsResponse.Response.ResponseMessage)
	}

	logData.Message = "LinkedAccountlist: API call completed successfully"
	responseBytes, err := userFetchAllBankAccountsResponse.Marshal()
	if err != nil {
		logData.Message = "LinkedAccountlist: Error marshaling response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	encrypted, err := security.Encrypt(responseBytes, []byte(authValues.Key))
	if err != nil {
		logData.Message = "LinkedAccountlist: Error encrypting response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "LinkedAccountlist: Response encrypted successfully"
	logData.ResponseSize = len(responseBytes)
	logData.EndTime = time.Now()
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(responseBytes)
	s.LoggerService.LogInfo(logData)

	return encrypted, nil
}

func (s *Store) PayeeNameGet(ctx context.Context, authValues *models.AuthValues, request *requests.GetAllBankAccount) (interface{}, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		UserID:        authValues.UserId,
		DeviceIP:      authValues.DeviceIp,
		LatLong:       authValues.LatLong,
		StartTime:     startTime,
		RequestID:     utils.GetRequestIDFromContext(ctx),
		RequestMethod: "POST",
		RequestURI:    "/api/upi/payeename",
		Message:       "Get Payee Name log",
	}

	// existingUserData, err := models.UpdateDeviceData(s.db, authValues.UserId, authValues.DeviceIp, authValues.OS, authValues.OSVersion, authValues.LatLong)
	// if err != nil {
	// 	logData.Message = "PayeeNameGet: Error updating device data"
	// 	s.LoggerService.LogError(logData)
	// 	return nil, err
	// }

	// logData.Message = "PayeeNameGet: Device data updated successfully"
	// s.LoggerService.LogInfo(logData)

	cryptoInfo, err := s.GenerateCryptoInfo(ctx, authValues)
	if err != nil {
		logData.Message = "PayeeNameGet: Error generating CryptoInfo"
		s.LoggerService.LogError(logData)
		return nil, fmt.Errorf("error generating CryptoInfo: %v", err)
	}

	logData.Message = "PayeeNameGet: CryptoInfo generated successfully"
	s.LoggerService.LogInfo(logData)

	userGetAllBankAccountsRequest := requests.NewOutgoingAccountLinkApiRequest()
	if err := userGetAllBankAccountsRequest.GetPayeeNameBind(
		request.MobileNumber,
		cryptoInfo,
	); err != nil {
		logData.Message = "PayeeNameGet: Error binding GetPayeeName request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "PayeeNameGet: Bind GetPayeeName request successful"
	s.LoggerService.LogInfo(logData)

	userGetAllBankAccountsResponse, err := s.bankService.LinkBankAccount(ctx, userGetAllBankAccountsRequest)
	if err != nil {
		logData.Message = "PayeeNameGet: Error calling bank service for GetPayeeName"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "PayeeNameGet: Bank service called successfully for GetPayeeName"
	s.LoggerService.LogInfo(logData)

	if userGetAllBankAccountsResponse.Response.ResponseCode != "0" {
		logData.Message = "PayeeNameGet: Received error code from bank service"
		s.LoggerService.LogError(logData)
		return nil, errors.New(userGetAllBankAccountsResponse.Response.ResponseMessage)
	}

	logData.Message = "PayeeNameGet: API call completed successfully"
	responseBytes, err := userGetAllBankAccountsResponse.Marshal()
	if err != nil {
		logData.Message = "PayeeNameGet: Error marshaling response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(responseBytes)
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	encrypted, err := security.Encrypt(responseBytes, []byte(authValues.Key))
	if err != nil {
		logData.Message = "PayeeNameGet: Error encrypting response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "PayeeNameGet: Response encrypted successfully"
	logData.ResponseSize = len(responseBytes)
	logData.EndTime = time.Now()
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(responseBytes)
	s.LoggerService.LogInfo(logData)

	return encrypted, nil
}

func (s *Store) AadharRequestListAccount(ctx context.Context, authValues *models.AuthValues, request *requests.AadharReqlistaccount) (interface{}, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		UserID:        authValues.UserId,
		DeviceIP:      authValues.DeviceIp,
		LatLong:       authValues.LatLong,
		StartTime:     startTime,
		RequestID:     utils.GetRequestIDFromContext(ctx),
		RequestMethod: "POST",
		RequestURI:    "/api/upi/aadhar-verification",
		Message:       "Verifying Aadhar log",
	}

	existingUserData, err := models.FindOneDeviceByUserID(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "AadharRequestListAccount: Error getting user data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	userData, err := models.GetUserDataByUserId(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "AadharRequestListAccount: Error getting user data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	personalInformation, err := models.GetPersonalInformation(s.db, existingUserData.UserId)
	if err != nil {
		logData.Message = "AadharRequestListAccount: Error fetching personal information"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	cryptoInfo, err := s.GenerateCryptoInfo(ctx, authValues)
	if err != nil {
		logData.Message = "AadharRequestListAccount: Error generating CryptoInfo"
		s.LoggerService.LogError(logData)
		return nil, fmt.Errorf("error generating CryptoInfo: %v", err)
	}

	upiData, err := models.GetAccountDataByUserId(s.db, authValues.UserId)
	if err != nil {
		return nil, err
	}

	userListAccountRequest := requests.NewOutgoingAadharRequestListAccountsApiRequest()
	if err := userListAccountRequest.Bind(
		userData.MobileNumber,
		existingUserData.DeviceIp.String,
		personalInformation.FirstName,
		cryptoInfo,
		upiData.UpiId.String,
		existingUserData.PackageId,
		authValues.LatLong,
	); err != nil {
		logData.Message = "AadharRequestListAccount: Error binding list accounts request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	userFetchAllBankAccountsResponse, err := s.bankService.AadharRequestListAccount(ctx, userListAccountRequest)
	if err != nil {
		logData.Message = "AadharRequestListAccount: Error calling bank service for list accounts"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if userFetchAllBankAccountsResponse.Response.ResponseCode != "0" {
		logData.Message = "AadharRequestListAccount: Received error code from bank service"
		s.LoggerService.LogError(logData)
		return nil, errors.New(userFetchAllBankAccountsResponse.Response.ResponseMessage)
	}

	logData.Message = "AadharRequestListAccount: API call completed successfully"
	responseBytes, err := userFetchAllBankAccountsResponse.Marshal()
	if err != nil {
		logData.Message = "AadharRequestListAccount: Error marshaling response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(responseBytes)
	logData.EndTime = time.Now()

	// Compare Aadhar first six digits
	userEnteredAadharFirstSix := request.AadharNumber[:6]
	apiResponseAadharFirstSix := userFetchAllBankAccountsResponse.Response.Response.RespListAccount.AccountList.Account[0].AadhaarNumber[:6]

	if userEnteredAadharFirstSix != apiResponseAadharFirstSix {
		logData.Message = "AadharRequestListAccount: user entered Aadhaar first six digits did not match the AadharRequestListAccount response value"
		s.LoggerService.LogError(logData)
		return nil, errors.New(constants.AadhaarNumberMismatchError)
	}

	SetUserUpiPinOtp := requests.NewOutgoingSetUpiPinReqOtpApiRequest()
	if err := SetUserUpiPinOtp.Bind(
		userData.MobileNumber,
		existingUserData.DeviceIp.String,
		authValues.LatLong,
		cryptoInfo,
		upiData.UpiId.String,
	); err != nil {
		logData.Message = "AadharRequestListAccount: Error binding set UPI PIN OTP request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	SetUserUpiPinOtpResponse, err := s.bankService.SetUpiPinReqOtp(ctx, SetUserUpiPinOtp)
	if err != nil {
		logData.Message = "AadharRequestListAccount: Error calling bank service for set UPI PIN OTP"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if SetUserUpiPinOtpResponse.Response.ResponseCode != "0" {
		logData.Message = "AadharRequestListAccount: Received error code from bank service for set UPI PIN OTP"
		s.LoggerService.LogError(logData)
		return nil, errors.New(SetUserUpiPinOtpResponse.Response.ResponseMessage)
	}

	logData.Message = "AadharRequestListAccount: API call completed successfully"
	responseBytes, err = userFetchAllBankAccountsResponse.Marshal()
	if err != nil {
		logData.Message = "AadharRequestListAccount: Error marshaling response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(responseBytes)
	logData.EndTime = time.Now()

	encrypted, err := security.Encrypt(responseBytes, []byte(authValues.Key))
	if err != nil {
		logData.Message = "AadharRequestListAccount: Error encrypting response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "AadharRequestListAccount: Response encrypted successfully"
	logData.ResponseSize = len(responseBytes)
	logData.EndTime = time.Now()
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(responseBytes)
	s.LoggerService.LogInfo(logData)

	return encrypted, nil
}

func (s *Store) GetUpiTokenXml(ctx context.Context, authValues *models.AuthValues, request *requests.UpiTokenXMLRequest) (interface{}, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		UserID:        authValues.UserId,
		DeviceIP:      authValues.DeviceIp,
		LatLong:       authValues.LatLong,
		StartTime:     startTime,
		RequestID:     utils.GetRequestIDFromContext(ctx),
		RequestMethod: "POST",
		RequestURI:    "/api/upi/get-token-xml",
		Message:       "Fetching UPI token XML log",
	}

	deviceData, err := models.FindOneDeviceByUserID(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "GetUpiTokenXml: Error getting user data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	newToken := ""

	if request.Challenge != "" {
		key := fmt.Sprintf("upi_token:%s:%s", authValues.UserId, deviceData.DeviceId)
		if _, err = s.bankService.Memory.Get(key); err != nil {
			logData.Message = "GetUpiTokenXml: Error retrieving UPI token from memory"
			s.LoggerService.LogError(logData)

			var newData interface{}
			newData, err = s.GetUpiToken(ctx, authValues, request.Challenge, "rotate")
			if err != nil {
				return nil, err
			}
			newToken = newData.(string)
		}
	}

	existingUserListKeysResponse, err := s.bankService.GetXmlRequestListKeys(ctx)
	if err != nil {
		logData.Message = "GetUpiTokenXml: Error fetching XML request list keys"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	var innerResp responses.RequestListKey

	// Unmarshal the inner JSON string inside the "Response" field
	err = json.Unmarshal([]byte(existingUserListKeysResponse.Response), &innerResp)
	if err != nil {
		log.Fatalf("Error unmarshalling inner JSON: %v", err)
	}

	// if existingUserListKeysResponse.Response.RequestListKeys[0].ResponseCode != "0" && existingUserListKeysResponse.Response.RequestListKeys[0].ResponseCode != "00" {
	// 	logData.Message = "GetUpiTokenXml: Received non-successful response code"
	// 	s.LoggerService.LogError(logData)
	// 	return nil, errors.New(existingUserListKeysResponse.Response.RequestListKeys[0].ResponseMessage)
	// }

	// xmlString := existingUserListKeysResponse.Response.RequestListKeys[0].Response

	responseData := map[string]interface{}{
		"token":      newToken,
		"xmlPayload": innerResp.Response,
	}
	if newToken == "" {
		responseData["token"] = nil
	}

	// jsonBytes, err := json.Marshal(responseData)
	// if err != nil {
	// 	logData.Message = "GetUpiTokenXml: Error marshaling response data to JSON"
	// 	s.LoggerService.LogError(logData)
	// 	return nil, err
	// }

	// encrypted, err := security.Encrypt(jsonBytes, []byte(authValues.Key))
	// if err != nil {
	// 	logData.Message = "GetUpiTokenXml: Error encrypting JSON response"
	// 	s.LoggerService.LogError(logData)
	// 	return nil, errors.New("failed to encrypt token")
	// }

	logData.Message = "GetUpiTokenXml: Response encrypted successfully"
	logData.EndTime = time.Now()
	logData.Latency = time.Since(startTime).Seconds()
	s.LoggerService.LogInfo(logData)

	return responseData, nil
}

func (s *Store) GetUpiToken(ctx context.Context, authValues *models.AuthValues, challenge, challengeType string) (interface{}, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		UserID:        authValues.UserId,
		DeviceIP:      authValues.DeviceIp,
		LatLong:       authValues.LatLong,
		StartTime:     startTime,
		RequestID:     utils.GetRequestIDFromContext(ctx),
		RequestMethod: "POST",
		RequestURI:    "/api/upi/get-token",
		Message:       "Get Upi Token log",
	}

	existingUserData, err := models.FindOneDeviceByUserID(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "GetUpiToken: Error getting user data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	userData, err := models.GetUserDataByUserId(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "SetUPIPin: Error getting user data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	key := fmt.Sprintf("upi_token:%s:%s", authValues.UserId, existingUserData.DeviceId)
	tokenData, _ := s.bankService.Memory.Get(key)
	if len(tokenData) > 0 || tokenData == "null" {
		// encrypted, err := security.Encrypt([]byte(tokenData), []byte(authValues.Key))
		if err != nil {
			logData.Message = "GetUpiToken: Error encrypting existing token data"
			s.LoggerService.LogError(logData)
			return nil, err
		}
		return tokenData, nil
	}

	decryptedDeviceId, err := security.Decrypt(existingUserData.DeviceId, []byte(authValues.Key))
	if err != nil {
		logData.Message = "CreateUpiID: Error decrypting device id"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	newUserRequestListkeys := requests.NewOutgoingReqListKeysApiRequest()
	if err := newUserRequestListkeys.Bind(authValues.DeviceIp, decryptedDeviceId, challengeType, existingUserData.PackageId, userData.MobileNumber, challenge); err != nil {
		logData.Message = "GetUpiToken: Error binding request for list keys"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "GetUpiToken: Bind request for list keys successful"

	userRequestListkeysResponse, err := s.bankService.NewUserRequestListKeys(ctx, newUserRequestListkeys)
	if err != nil {
		logData.Message = "GetUpiToken: Error calling bank service for list keys"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "GetUpiToken: Bank service called successfully for list keys"

	keyResponse := userRequestListkeysResponse.Response
	if keyResponse.ResponseCode != "0" {
		logData.Message = "GetUpiToken: Received error code from list keys response"
		s.LoggerService.LogError(logData)
		return nil, errors.New(keyResponse.ResponseMessage)
	}

	// Save UPI token
	tokenData = keyResponse.Response.Ns2RespListKeys.KeyList.Key.KeyValue.Text
	// hexToken := s.convertUpiToken(tokenData)

	if tokenData == "" {
		return nil, errors.New("getting empty token from api")
	}

	err = s.SaveUPIToken(ctx, authValues.UserId, existingUserData.DeviceId, tokenData)
	if err != nil {
		logData.Message = "GetUpiToken: Error saving UPI token"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	// encrypted, err := security.Encrypt([]byte(tokenData), []byte(authValues.Key))
	// if err != nil {
	// 	logData.Message = "GetUpiToken: Error encrypting UPI token"
	// 	s.LoggerService.LogError(logData)
	// 	return nil, errors.New("failed to encrypt token")
	// }

	logData.Message = "GetUpiToken: Response encrypted successfully"
	logData.ResponseSize = len(tokenData)
	logData.EndTime = time.Now()
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(tokenData)
	s.LoggerService.LogInfo(logData)

	return tokenData, nil
}

func (s *Store) SaveAndExtractToken(ctx context.Context, xmlData, userId, deviceIp string) (interface{}, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		UserID:        userId,
		DeviceIP:      deviceIp,
		StartTime:     startTime,
		RequestID:     utils.GetRequestIDFromContext(ctx),
		RequestMethod: "POST",
		RequestURI:    "/api/upi/save-extract-token",
		Message:       "Save and extract UPI token log",
	}

	var respListKeys requests.RespListKeys
	err := xml.Unmarshal([]byte(xmlData), &respListKeys)
	if err != nil {
		logData.Message = "SaveAndExtractToken: Error unmarshaling XML"
		s.LoggerService.LogError(logData)
		return nil, fmt.Errorf("error unmarshaling XML: %v", err)
	}

	logData.Message = "SaveAndExtractToken: XML unmarshaled successfully"
	s.LoggerService.LogInfo(logData)

	tokenData := respListKeys.KeyList.Key.KeyValue
	hexToken := s.convertUpiToken(ctx, tokenData)

	err = s.SaveUPIToken(ctx, userId, deviceIp, hexToken)
	if err != nil {
		logData.Message = "SaveAndExtractToken: Error saving UPI token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "SaveAndExtractToken: UPI token saved successfully"
	logData.ResponseSize = len(tokenData)
	logData.EndTime = time.Now()
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(tokenData)
	s.LoggerService.LogInfo(logData)

	return tokenData, nil
}

// convert upi token to hex string.
func (s *Store) convertUpiToken(ctx context.Context, tokenData string) string {

	//startTime := time.Now()

	// logData := &commonSrv.LogEntry{
	// 	Action:        constants.UPI,
	// 	StartTime:     startTime,
	// 	RequestID:     utils.GetRequestIDFromContext(ctx),
	// 	RequestMethod: "Internal Token conversion Function",
	// 	RequestURI:    "/store/convertUpiToken",
	// 	Message:       "Converting UPI token to hex format log",
	// }

	// byteArray := []byte(tokenData)

	// if len(byteArray) < 16 {
	// 	byteArray = append(byteArray, make([]byte, 16-len(byteArray))...)
	// } else if len(byteArray) > 32 {
	// 	byteArray = byteArray[:32]
	// }

	// hexToken := hex.EncodeToString(tokenData)

	// logData.Message = "UPI token converted successfully"
	// logData.EndTime = time.Now()
	// logData.Latency = time.Since(startTime).Seconds()
	// logData.ResponseBody = string(hexToken)
	// logData.ResponseSize = len(hexToken)
	// s.LoggerService.LogInfo(logData)

	return tokenData
}

// save upi token in redis.
func (s *Store) SaveUPIToken(ctx context.Context, userId, deviceIp, tokenData string) error {

	// startTime := time.Now()

	// logData := &commonSrv.LogEntry{
	// 	Action:        constants.UPI,
	// 	UserID:        userId,
	// 	DeviceIP:      deviceIp,
	// 	StartTime:     startTime,
	// 	RequestID:     utils.GetRequestIDFromContext(ctx),
	// 	RequestMethod: "Internal Token Saving Function",
	// 	RequestURI:    "/store/SaveUPIToken",
	// 	Message:       "Saving UPI token to memory log",
	// }

	// key := fmt.Sprintf("upi_token:%s:%s", userId, deviceIp)
	// tokenTTL := time.Hour * 24 * 30

	// err := s.bankService.Memory.Set(key, tokenData, tokenTTL)
	// if err != nil {
	// 	logData.Message = fmt.Sprintf("Error saving UPI token: %v", err)
	// 	s.LoggerService.LogError(logData)
	// 	return err
	// }

	// logData.Message = "UPI token saved successfully"
	// logData.ResponseSize = len(tokenData)
	// logData.EndTime = time.Now()
	// logData.Latency = time.Since(startTime).Seconds()
	// logData.ResponseBody = string(tokenData)
	// s.LoggerService.LogInfo(logData)

	return nil
}

func (s *Store) CollectUpiMoneyCount(ctx context.Context, authValues *models.AuthValues) (interface{}, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		UserID:        authValues.UserId,
		DeviceIP:      authValues.DeviceIp,
		LatLong:       authValues.LatLong,
		StartTime:     startTime,
		RequestID:     utils.GetRequestIDFromContext(ctx),
		RequestMethod: "GET",
		RequestURI:    "/api/upi/collect-count",
		Message:       "Fetching collect UPI money count log",
	}

	userData, err := models.GetUserDataByUserId(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "SetUPIPin: Error getting user data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "CollectUpiMoneyCount: Device data updated successfully"
	s.LoggerService.LogInfo(logData)

	cryptoInfo, err := s.GenerateCryptoInfo(ctx, authValues)
	if err != nil {
		logData.Message = "CollectUpiMoneyCount: Error generating CryptoInfo"
		s.LoggerService.LogError(logData)
		return nil, fmt.Errorf("error generating CryptoInfo: %v", err)
	}

	collectcount := requests.NewOutgoingUpiMoneyCollectCountApiRequest()
	if err := collectcount.Bind(userData.MobileNumber, cryptoInfo); err != nil {
		logData.Message = "CollectUpiMoneyCount: Error binding collect count request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "CollectUpiMoneyCount: Bind request successful"
	s.LoggerService.LogInfo(logData)

	collectcountResponse, err := s.bankService.CollectCount(ctx, collectcount)
	if err != nil {
		logData.Message = "CollectUpiMoneyCount: Error calling bank service"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "CollectUpiMoneyCount: Bank service called successfully"
	s.LoggerService.LogInfo(logData)

	if collectcountResponse.Response.ResponseCode != "0" {
		logData.Message = "CollectUpiMoneyCount: Received error code from response"
		s.LoggerService.LogError(logData)
		return nil, errors.New(collectcountResponse.Response.ResponseMessage)
	}

	logData.Message = "CollectUpiMoneyCount: API call completed successfully"
	responseBytes, err := collectcountResponse.Marshal()
	if err != nil {
		logData.Message = "CollectUpiMoneyCount: Error marshaling response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(responseBytes)
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	encrypted, err := security.Encrypt(responseBytes, []byte(authValues.Key))
	if err != nil {
		logData.Message = "CollectUpiMoneyCount: Error encrypting response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "CollectUpiMoneyCount: Response encrypted successfully"
	logData.ResponseSize = len(responseBytes)
	logData.EndTime = time.Now()
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(responseBytes)
	s.LoggerService.LogInfo(logData)

	return encrypted, nil
}

func (s *Store) CollectUpiMoney(ctx context.Context, authValues *models.AuthValues, request *requests.UpiCollectMoneyRequest) (interface{}, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		UserID:        authValues.UserId,
		DeviceIP:      authValues.DeviceIp,
		LatLong:       authValues.LatLong,
		StartTime:     startTime,
		RequestID:     utils.GetRequestIDFromContext(ctx),
		RequestMethod: "POST",
		RequestURI:    "/api/upi/collect-money",
		Message:       "Collect Upi Money log",
	}

	userData, err := models.GetUserDataByUserId(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "CollectUpiMoney: Error getting user data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "CollectUpiMoney: Device data updated successfully"
	s.LoggerService.LogInfo(logData)

	cryptoInfo, err := s.GenerateCryptoInfo(ctx, authValues)
	if err != nil {
		logData.Message = "CollectUpiMoney: Error generating CryptoInfo"
		s.LoggerService.LogError(logData)
		return nil, fmt.Errorf("error generating CryptoInfo: %v", err)
	}

	collectdetails := requests.NewOutgoingUpiMoneyCollectDetailsApiRequest()
	if err := collectdetails.Bind(userData.MobileNumber, cryptoInfo); err != nil {
		logData.Message = "CollectUpiMoney: Error binding collect details request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "CollectUpiMoney: Bind collect details request successful"
	s.LoggerService.LogInfo(logData)

	collectdetailsResponse, err := s.bankService.CollectDetails(ctx, collectdetails)
	if err != nil {
		logData.Message = "CollectUpiMoney: Error calling bank service for collect details"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "CollectUpiMoney: Bank service called successfully for collect details"
	s.LoggerService.LogInfo(logData)

	if collectdetailsResponse.Response.ResponseCode != "0" {
		logData.Message = "CollectUpiMoney: Received error code from collect details response"
		s.LoggerService.LogError(logData)
		return nil, errors.New(collectdetailsResponse.Response.ResponseMessage)
	}

	collectApproval := requests.NewOutgoingUpiMoneyCollectApprovalApiRequest()
	if err := collectApproval.Bind(userData.MobileNumber, cryptoInfo, request.UpiPin, collectdetailsResponse); err != nil {
		logData.Message = "CollectUpiMoney: Error binding collect approval request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "CollectUpiMoney: Bind collect approval request successful"
	s.LoggerService.LogInfo(logData)

	collectApprovalResponse, err := s.bankService.CollectApproval(ctx, collectApproval)
	if err != nil {
		logData.Message = "CollectUpiMoney: Error calling bank service for collect approval"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "CollectUpiMoney: Bank service called successfully for collect approval"
	s.LoggerService.LogInfo(logData)

	if collectApprovalResponse.Response.ResponseCode != "0" {
		logData.Message = "CollectUpiMoney: Received error code from collect approval response"
		s.LoggerService.LogError(logData)
		return nil, errors.New(collectApprovalResponse.Response.ResponseMessage)
	}

	responseBytes, err := collectApprovalResponse.Marshal()
	if err != nil {
		logData.Message = "CollectUpiMoney: Error marshaling response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	encrypted, err := security.Encrypt(responseBytes, []byte(authValues.Key))
	if err != nil {
		logData.Message = "CollectUpiMoney: Error encrypting response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "CollectUpiMoney: Response encrypted successfully"
	logData.ResponseSize = len(responseBytes)
	logData.EndTime = time.Now()
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(responseBytes)
	s.LoggerService.LogInfo(logData)

	return encrypted, nil
}

// upi transactions
func (s *Store) generateTransactionID(ctx context.Context, authValues *models.AuthValues) (string, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		UserID:        authValues.UserId,
		RequestID:     utils.GetRequestIDFromContext(ctx),
		StartTime:     startTime,
		RequestMethod: "GET",
		RequestURI:    "/generate/transaction-id",
		Message:       "Generate Transaction ID log",
	}

	randomBytes := make([]byte, 14)
	_, err := rand.Read(randomBytes)
	if err != nil {
		logData.Message = "generateTransactionID: Error generating random bytes"
		s.LoggerService.LogError(logData)
		return "", fmt.Errorf("error generating random bytes: %v", err)
	}

	randomHex := hex.EncodeToString(randomBytes)
	result := fmt.Sprintf("PAYDOH%s", randomHex)

	logData.Message = "generateTransactionID: Transaction ID generated successfully"
	logData.ResponseSize = len(result)
	logData.EndTime = time.Now()
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(result)
	s.LoggerService.LogInfo(logData)

	return result, nil
}

func (s *Store) CreateTransaction(ctx context.Context, authValues *models.AuthValues, credType string) (string, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		UserID:        authValues.UserId,
		DeviceIP:      authValues.DeviceIp,
		LatLong:       authValues.LatLong,
		StartTime:     startTime,
		RequestID:     utils.GetRequestIDFromContext(ctx),
		RequestMethod: "POST",
		RequestURI:    "/create/transaction",
		Message:       "Create Transaction log",
	}

	txnId, err := s.generateTransactionID(ctx, authValues)
	if err != nil {
		logData.Message = "CreateTransaction: Error generating transaction ID"
		s.LoggerService.LogError(logData)
		return "", err
	}

	if strings.ToLower(credType) == "pay" {
		err = models.InsertTransaction(s.db, &models.Transaction{
			UserID:        authValues.UserId,
			TransactionID: txnId,
			PaymentMode:   models.PaymentModeUPI,
		})
		if err != nil {
			logData.Message = "CreateTransaction: Error creating UPI transaction"
			s.LoggerService.LogError(logData)
			return "", err
		}
	} else {
		transaction := &models.UPITransaction{
			UserID:        authValues.UserId,
			TransactionID: txnId,
			CredType:      credType,
		}

		err = models.CreateUPITransaction(s.db, transaction)
		if err != nil {
			logData.Message = "CreateTransaction: Error creating UPI transaction"
			s.LoggerService.LogError(logData)
			return "", err

		}
	}

	responseData := map[string]interface{}{
		"txnId": [1]string{txnId},
	}
	responseDataByte, _ := json.Marshal(responseData)
	encrypted, err := security.Encrypt(responseDataByte, []byte(authValues.Key))

	if err != nil {
		logData.Message = "CreateTransaction: Error encrypting response data"
		s.LoggerService.LogError(logData)
		return "", err
	}

	logData.Message = "CreateTransaction: Transaction created successfully"
	logData.ResponseSize = len(responseDataByte)
	logData.EndTime = time.Now()
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(responseDataByte)
	s.LoggerService.LogInfo(logData)

	return encrypted, nil
}

func (s *Store) UpiTransactionsHistoryList(ctx context.Context, authValues *models.AuthValues, request *requests.IncomingUpiTransactionHistoryApiRequest) (interface{}, error) {
	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		StartTime:     startTime,
		UserID:        authValues.UserId,
		RequestID:     utils.GetRequestIDFromContext(ctx),
		RequestMethod: "POST",
		RequestURI:    "/api/upi/transaction-history",
		Message:       "UpiTransactionsHistoryList log",
	}

	existingUserData, err := models.GetUserDataByUserId(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "UpiTransactionsHistoryList: Error getting user data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	cryptoInfo, err := s.GenerateCryptoInfo(ctx, authValues)
	if err != nil {
		logData.Message = "UpiTransactionsHistoryList: Error generating CryptoInfo"
		s.LoggerService.LogError(logData)
		return nil, fmt.Errorf("error generating CryptoInfo: %v", err)
	}

	upiTransactionListRequest := requests.NewOutgoingUpiTransactionHistoryApiRequest()
	if err := upiTransactionListRequest.Bind(
		existingUserData.MobileNumber,
		cryptoInfo,
		request,
	); err != nil {
		logData.Message = "UpiTransactionsHistoryList: Error binding in checking upi transaction list request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	upiTransactionListResponse, err := s.bankService.UpiTransactionHistory(ctx, upiTransactionListRequest)
	if err != nil {
		logData.Message = "UpiTransactionsHistoryList: Error calling bank service for getting upi transaction list"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if upiTransactionListResponse.Response.ResponseCode == constants.RetryUpiError1 &&
		(strings.EqualFold(upiTransactionListResponse.Response.ResponseMessage, constants.UpiErrorMessageNoRecordsFound) ||
			strings.EqualFold(upiTransactionListResponse.Response.ResponseMessage, constants.UpiErrorMessageNoDataFound)) {

		logData.Message = "UpiTransactionsHistoryList: No data found in the response"
		s.LoggerService.LogInfo(logData)
		return nil, errors.New(constants.UpiNoRecordFoundErrorMessage)
	}

	responseBytes, err := upiTransactionListResponse.Marshal()
	if err != nil {
		logData.Message = "UpiTransactionsHistoryList: Error marshaling response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	encrypted, err := security.Encrypt(responseBytes, []byte(authValues.Key))
	if err != nil {
		logData.Message = "UpiTransactionsHistoryList: Error encrypting response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "UpiTransactionsHistoryList: Response encrypted successfully"
	logData.ResponseSize = len(responseBytes)
	logData.EndTime = time.Now()
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(responseBytes)
	s.LoggerService.LogInfo(logData)

	return encrypted, nil
}

func (s *Store) ChangeUpiPin(ctx context.Context, authValues *models.AuthValues, request *requests.IncomingUpiChangeUpiPinRequest) (interface{}, error) {
	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		StartTime:     startTime,
		UserID:        authValues.UserId,
		RequestID:     utils.GetRequestIDFromContext(ctx),
		RequestMethod: "POST",
		RequestURI:    "/api/upi/change-upi-pin",
		Message:       "ChangeUpiPin log",
	}

	existingUserData, err := models.GetUserDataByUserId(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "ChangeUpiPin: Error getting user data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	cryptoInfo, err := s.GenerateCryptoInfo(ctx, authValues)
	if err != nil {
		logData.Message = "ChangeUpiPin: Error generating CryptoInfo"
		s.LoggerService.LogError(logData)
		return nil, fmt.Errorf("error generating CryptoInfo: %v", err)
	}

	upiData, err := models.GetAccountDataByUserId(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "ChangeUpiPin: Error getting user data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	upiChangeupiPinRequest := requests.NewOutgoingUpiReqSetCreRequest()
	if err := upiChangeupiPinRequest.Bind(
		existingUserData.MobileNumber,
		upiData.UpiId.String,
		authValues.DeviceIp,
		cryptoInfo,
		authValues.LatLong,
		request,
	); err != nil {
		logData.Message = "ChangeUpiPin: Error binding in change upi pin request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	upiChangeupiPinResponse, err := s.bankService.UpiChangeUpiPin(ctx, upiChangeupiPinRequest)
	if err != nil {
		logData.Message = "ChangeUpiPin: Error calling bank service for change upi pin"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if upiChangeupiPinResponse.Response.ResponseCode != "0" {
		logData.Message = "ChangeUpiPin: Received error code from bank service for change upi pin"
		s.LoggerService.LogError(logData)
		return nil, errors.New(upiChangeupiPinResponse.Response.ResponseMessage)
	}

	responseBytes, err := upiChangeupiPinResponse.Marshal()
	if err != nil {
		logData.Message = "ChangeUpiPin: Error marshaling response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if err := s.auditLogSrv.Save(ctx, &services.AuditLog{
		UserID:         authValues.UserId,
		RequestURL:     "/api/upi/change-upi-pin",
		HTTPMethod:     "POST",
		ResponseStatus: 200,
		Action:         constants.CHANGE_UPI_PIN,
	}); err != nil {
		logData.Message = "error while saving audit log" + err.Error()
		logData.EndTime = time.Now()
		s.LoggerService.LogInfo(logData)
	}

	logData.Message = "ChangeUpiPin: Response encrypted successfully"
	logData.ResponseSize = len(responseBytes)
	logData.EndTime = time.Now()
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(responseBytes)
	s.LoggerService.LogInfo(logData)

	return nil, nil
}
