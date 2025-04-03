package kyc

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"bitbucket.org/paydoh/paydoh-commons/database"
	commonSrv "bitbucket.org/paydoh/paydoh-commons/services"
	"github.com/google/uuid"

	"bankapi/config"
	"bankapi/constants"
	"bankapi/models"
	"bankapi/requests"
	"bankapi/responses"
	"bankapi/security"
	"bankapi/services"
	"bankapi/utils"
)

type Kyc interface {
	UpdateKycConsent(authValues *models.AuthValues, request *requests.KycConsentRequest) (interface{}, error)
	GetKycConsent(authValues *models.AuthValues) (interface{}, error)
}

type Store struct {
	db            *sql.DB
	m             *database.Document
	memory        *database.InMemory
	bankService   *services.BankApiService
	LoggerService *commonSrv.LoggerService
	auditLogSrv   services.AuditLogService
}

func NewStore(log *commonSrv.LoggerService, db *sql.DB, m *database.Document, memory *database.InMemory, auditLogSrv services.AuditLogService) *Store {
	bankService := services.NewBankApiService(log, memory)
	return &Store{
		db:            db,
		m:             m,
		memory:        memory,
		bankService:   bankService,
		LoggerService: log,
		auditLogSrv:   auditLogSrv,
	}
}

func (s *Store) UpdateKycConsent(ctx context.Context, authValues *models.AuthValues, request *requests.KycConsentRequest) error {
	logData := &commonSrv.LogEntry{
		Action:     constants.KYC,
		RequestURI: "/api/kyc/consent",
		Message:    "UpdateKycConsent log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	existingUserData, err := models.GetUserDataByUserId(s.db, authValues.UserId)

	if err != nil {
		logData.Message = "UpdateKycConsent: Error getting user data"
		s.LoggerService.LogError(logData)
		return err
	}

	kycConsentData, err := models.FindKycConsentByUserIdV2(s.db, existingUserData.UserId)
	if err != nil {
		if errors.Is(err, constants.ErrKycConsentNotProvided) {
			if err := models.InsertKycConsent(s.db, &models.KycConsent{
				UserId:                   existingUserData.UserId,
				IndianResident:           false,
				PoliticallyExposedPerson: false,
				AadharConsent:            false,
				VirtualDebitCardConsent:  false,
				PhysicalDebitCardConsent: false,
			}); err != nil {
				logData.Message = "UpdateKycConsent: Error while inserting kyc consent"
				s.LoggerService.LogError(logData)
				return err
			}

			smsData := &requests.OutgoingSmsVerificationRequest{
				ServiceBy:     "PAYDOH",
				ApplicantId:   existingUserData.ApplicantId,
				MobileNo:      existingUserData.MobileNumber,
				TxnIdentifier: uuid.New().String(),
				ServiceStatus: "Verified",
			}
			_, err = s.bankService.SmsVerification(ctx, smsData)
			if err != nil {
				logData.Message = "UpdateKycConsent: Error while sending sms verification"
				s.LoggerService.LogError(logData)
				return err
			}

			// updating only bank_sms_verification_status in db
			if err := models.UpdateKycConsent(s.db, &models.KycConsent{
				BankSmsVerificationStatus: true,
			}, existingUserData.UserId); err != nil {
				logData.Message = "UpdateKycConsent: Error while updating sms verification status in kyc consent"
				s.LoggerService.LogError(logData)
				return err
			}

			if err := s.ConsentAddition(ctx, request, existingUserData, authValues); err != nil {
				return err
			}

			return nil

		}

		return err
	}

	// sms verification api call as per flag
	if !kycConsentData.BankSmsVerificationStatus {
		smsData := &requests.OutgoingSmsVerificationRequest{
			ServiceBy:     "PAYDOH",
			ApplicantId:   existingUserData.ApplicantId,
			MobileNo:      existingUserData.MobileNumber,
			TxnIdentifier: uuid.New().String(),
			ServiceStatus: "Verified",
		}
		_, err = s.bankService.SmsVerification(ctx, smsData)
		if err != nil {
			logData.Message = "UpdateKycConsent: Error while sending sms verification"
			s.LoggerService.LogError(logData)
			return err
		}

		// updating only bank_sms_verification_status in db
		if err := models.UpdateKycConsent(s.db, &models.KycConsent{
			BankSmsVerificationStatus: true,
		}, existingUserData.UserId); err != nil {
			logData.Message = "UpdateKycConsent: Error while updating sms verification status in kyc consent"
			s.LoggerService.LogError(logData)
			return err
		}
	}

	if err := s.ConsentAddition(ctx, request, existingUserData, authValues); err != nil {
		return err
	}

	// return nil

	logData.Message = "UpdateKycConsent: Updated successfully"
	s.LoggerService.LogError(logData)
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)
	return nil
}

func (s *Store) ConsentAddition(ctx context.Context, request *requests.KycConsentRequest, existingUserData *models.UserData, authValues *models.AuthValues) error {
	logData := &commonSrv.LogEntry{
		Action:     constants.KYC,
		RequestURI: "/api/kyc/consent",
		Message:    "UpdateKycConsent log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	insertingData := models.NewKycConsent()
	insertingData.Bind(request, existingUserData.UserId)
	// customer consent bank api call
	outGoingRequest := requests.NewOutgoingConsentRequest()
	if err := outGoingRequest.BindAndValidate(existingUserData.ApplicantId); err != nil {
		logData.Message = "UpdateKycConsent: Error binding and validating consent request"
		s.LoggerService.LogError(logData)
		return err
	}

	consentDetailsData := []requests.ConsentDetails{}

	if request.IndianResident {
		s.addConsent("FATCA", &consentDetailsData, outGoingRequest)
	}

	if request.NotPoliticallyExposed {
		s.addConsent("PEP", &consentDetailsData, outGoingRequest)
	}

	if request.AadharConsent {
		s.addConsent("AADHAAR1", &consentDetailsData, outGoingRequest)
	}

	if request.PhysicalDebitCardConsent {
		s.addConsent("PHYSICAL_DEBIT_CARD", &consentDetailsData, outGoingRequest)
	}

	if request.VirtualDebitCardConsent {
		s.addConsent("VIRTUAL_DEBIT_CARD", &consentDetailsData, outGoingRequest)
	}

	if request.Aadhar2Consent {
		s.addConsent("AADHAAR2", &consentDetailsData, outGoingRequest)
	}

	if request.AddressChangeConsent {
		s.addConsent("ADDRESS_CHANGE", &consentDetailsData, outGoingRequest)
	}

	if request.NominationConsent {
		s.addConsent("NOMINATION", &consentDetailsData, outGoingRequest)
	}
	if request.LocationConsent {
		s.addConsent("LOCATION", &consentDetailsData, outGoingRequest)
	}

	if request.PrivacyPolicyConsent {
		s.addConsent("PRIVACY_POLICY", &consentDetailsData, outGoingRequest)
	}

	if request.TermsAndCondition {
		s.addConsent("TERMS_AND_CONDITIONS", &consentDetailsData, outGoingRequest)
	}

	if len(consentDetailsData) > 0 {
		_, err := s.bankService.PostUserConsent(ctx, outGoingRequest, consentDetailsData)
		if err != nil {
			if bankErr := s.bankService.HandleBankSpecificError(err, func(errorCode string) (string, bool) {
				return constants.GetConsentErrorMessage(errorCode)
			}); bankErr != nil {
				return errors.New(bankErr.ErrorMessage)
			}

			logData.Message = "UpdateKycConsent: Error posting user consent to bank service"
			s.LoggerService.LogError(logData)
			return err
		}

		if err := models.UpsertKycConsent(s.db, existingUserData.UserId, insertingData); err != nil {
			logData.Message = "UpdateKycConsent: Error insert/update kyc consent"
			s.LoggerService.LogError(logData)
			return err
		}

		if request.VirtualDebitCardConsent {
			if err := models.UpdateUserOnboardingStatus(constants.DEBIT_CARD_CONSENT_STAGE, existingUserData.UserId); err != nil {
				logData.Message = "UpdateKycConsent: error while updating onboarding status"
				s.LoggerService.LogError(logData)
			}
		}

		// save audit log
		requestConsent := requests.ConsentRequestV2{
			ApplicantId:    outGoingRequest.ApplicantId,
			TxnIdentifier:  outGoingRequest.TxnIdentifier,
			ConsentDetails: consentDetailsData,
		}

		requestData, err := json.Marshal(requestConsent)
		if err != nil {
			logData.Message = "UpdateKycConsent:error while marshalling consent"
			s.LoggerService.LogError(logData)
			return err
		}

		encryptReq, err := security.Encrypt(requestData, []byte(authValues.Key))
		if err != nil {
			logData.Message = "UpdateAdressLog:error while encrypting request"
			logData.EndTime = time.Now()
			s.LoggerService.LogError(logData)
			return nil
		}

		if err := s.auditLogSrv.Save(ctx, &services.AuditLog{
			TransactionID:  outGoingRequest.TxnIdentifier,
			UserID:         existingUserData.UserId,
			ApplicantID:    existingUserData.ApplicantId,
			RequestURL:     "/fintech/v2/customer/consent",
			HTTPMethod:     "POST",
			ResponseStatus: 200,
			Action:         constants.BANK_CONSENT,
			RequestBody:    encryptReq,
		}); err != nil {
			logData.Message = "UpdateKycConsent:error while saving audit log"
			s.LoggerService.LogError(logData)
		}

		return nil
	}

	return errors.New("consent not provided")
}

func (s *Store) addConsent(consentType string, consentDetailsData *[]requests.ConsentDetails, outGoingRequest *requests.OutgoingConsentRequest) {
	consent := requests.ConsentDetails{
		ConsentType:     consentType,
		ConsentProvided: "Yes",
		ConsentTime:     outGoingRequest.ConsentTime,
	}
	*consentDetailsData = append(*consentDetailsData, consent)
}

func (s *Store) GetKycConsent(ctx context.Context, authValues *models.AuthValues) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.KYC,
		RequestURI: "/api/kyc/consent",
		Message:    "GetKycConsent log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	kycConsent, err := models.FindKycConsentByUserId(s.db, authValues.UserId)

	if err != nil {
		logData.Message = "GetKycConsent: Error finding kyc consent by user id"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if kycConsent == nil {
		logData.Message = "GetKycConsent: Kyc consent not found"
		s.LoggerService.LogError(logData)
		return nil, errors.New("kyc consent not found")
	}

	kycConsentData, err := kycConsent.Marshal()

	if err != nil {
		logData.Message = "GetKycConsent: Error marshaling kyc consent data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	encrypted, err := security.Encrypt(kycConsentData, []byte(authValues.Key))

	if err != nil {
		logData.Message = "GetKycConsent: Error encrypting kyc consent data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "GetKycConsent: Kyc consent data encrypted successfully"
	logData.ResponseSize = len(kycConsentData)
	logData.ResponseBody = string(kycConsentData)
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return encrypted, nil
}

func (s *Store) GetVcipUrl(ctx context.Context, authValues *models.AuthValues) (string, error) {

	logData := &commonSrv.LogEntry{
		Action:     constants.KYC,
		RequestURI: "/api/kyc/vcip-url",
		Message:    "GetVcipUrl log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	existingUserData, err := models.GetUserDataByUserId(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "GetVcipUrl: Error getting user data"
		s.LoggerService.LogError(logData)
		return "", err
	}

	kycConsent, err := models.FindKycConsentByUserId(s.db, authValues.UserId)

	if err != nil {
		logData.Message = "GetVcipUrl: Error finding kyc consent by user id"
		s.LoggerService.LogError(logData)
		return "", err
	}

	if kycConsent == nil {
		logData.Message = "GetVcipUrl: Kyc consent not found"
		s.LoggerService.LogError(logData)
		return "", errors.New("kyc consent not found")
	}

	if !kycConsent.AadharConsent {
		logData.Message = "GetVcipUrl: Kyc consent for aadhar not found"
		s.LoggerService.LogError(logData)
		return "", errors.New("kyc consent for aadhar not found")
	}

	if !kycConsent.IndianResident {
		logData.Message = "GetVcipUrl: Kyc consent for indian resident not found"
		s.LoggerService.LogError(logData)
		return "", errors.New("kyc consent for indian resident not found")
	}

	if !kycConsent.PoliticallyExposedPerson {
		logData.Message = "GetVcipUrl: Kyc consent for political exposed person not found"
		s.LoggerService.LogError(logData)
		return "", errors.New("kyc consent for political exposed person not found")
	}

	outgoingRequest := requests.NewOutgoingVcipRequest()
	if err := outgoingRequest.Bind(existingUserData.ApplicantId, existingUserData.MobileNumber); err != nil {
		logData.Message = "GetVcipUrl: Error binding outgoing request"
		s.LoggerService.LogError(logData)
		return "", err
	}

	var response *responses.KycInvokeResponse
	var lastRetryMessage string

	response, opErr := s.bankService.VcipInvoke(ctx, outgoingRequest)
	if opErr != nil {
		logData.Message = "GetVcipUrl: Error invoking vcip"
		s.LoggerService.LogError(logData)
		return "", opErr
	}

	if errorMessage, exists := constants.GetVcipInvokeErrorMessage(response.ResponseCode, response.Status); exists {
		logData.Message = fmt.Sprintf("Non-retryable error detected. ResponseCode: %s", response.ResponseCode)
		s.LoggerService.LogError(logData)
		return "", errors.New(errorMessage)
	}

	retryErr := utils.RetryFunc(func() error {
		var opErr error

		if response == nil || response.Data == "" {
			lastRetryMessage = constants.RetryErrorMessage
			return errors.New(lastRetryMessage)
		}

		if retryMessage, retryable := constants.GetVcipInvokeErrorRetryMessage(response.ResponseCode); retryable {
			logData.Message = fmt.Sprintf("Retryable error detected. ResponseCode: %s, retrying...", response.ResponseCode)
			s.LoggerService.LogError(logData)

			response, opErr = s.bankService.VcipInvoke(ctx, outgoingRequest)
			if opErr != nil {
				return opErr
			}

			lastRetryMessage = retryMessage
			return errors.New(retryMessage)
		}

		return nil
	}, 2)

	if retryErr != nil {
		logData.Message = fmt.Sprintf("VCIP API invocation failed: %s", lastRetryMessage)
		s.LoggerService.LogError(logData)
		return "", errors.New(lastRetryMessage)
	}

	// update onboarding status
	if err := models.UpdateUserOnboardingStatus(constants.AGENT_URL_STEP, authValues.UserId); err != nil {
		logData.Message = "GetVcipUrl: error while updating onboarding status"
		s.LoggerService.LogError(logData)
	}

	logData.Message = "GetVcipUrl: Vcip invocation successful"
	logData.ResponseSize = len(response.Data)
	logData.ResponseBody = string(response.Data)
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return response.Data, nil
}

func (s *Store) AddKycUpdateData(ctx context.Context, applicantId, status, acom, astat string) error {
	logData := &commonSrv.LogEntry{
		Action:     constants.KYC,
		RequestURI: "/api/callback/kyc-update",
		Message:    "AddKycUpdateData log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	db := config.GetDB()

	userData, err := models.GetUserDataByApplicantId(db, applicantId)
	if err != nil {
		logData.Message = "AddKycUpdateData: Error getting user data by applicant id"
		s.LoggerService.LogError(logData)
		return err
	}

	// check if already kyc consent completed
	_, err = models.FindKycConsentByUserId(db, userData.UserId)
	if err != nil {
		logData.Message = "AddKycUpdateData: Error finding kyc consent by user id"
		s.LoggerService.LogError(logData)
		return errors.New("kyc consent data not found for this user")
	}

	acomInt, err := strconv.Atoi(acom)
	if err != nil {
		logData.Message = "AddKycUpdateData: Error converting acom to int"
		s.LoggerService.LogError(logData)
		fmt.Println("Error:", err)
	}

	astatInt, err := strconv.Atoi(astat)
	if err != nil {
		logData.Message = "AddKycUpdateData: Error converting astat to int"
		s.LoggerService.LogError(logData)
		fmt.Println("Error:", err)
	}

	kycUpdateData := &models.KycUpdateData{
		UserId: userData.UserId,
		Acom:   acomInt,
		Astat:  astatInt,
		Status: status,
	}

	if _, err := models.FindOneKycUpdateData(db, userData.UserId); err != nil {
		if errors.Is(err, constants.ErrNoDataFound) {
			err = models.InsertKycUpdateData(db, kycUpdateData)
			if err != nil {
				logData.Message = "AddKycUpdateData: Error inserting kyc update data"
				s.LoggerService.LogError(logData)
				return err
			}

			// update onboarding status
			if err := models.UpdateUserOnboardingStatus(constants.KYC_CALLBACK_STEP, userData.UserId); err != nil {
				logData.Message = "GetVcipUrl: error while updating onboarding status"
				s.LoggerService.LogError(logData)
			}

			return nil
		}

		logData.Message = "AddKycUpdateData: Error finding kyc update data"
		s.LoggerService.LogError(logData)
		return err
	}

	err = models.UpdateKycUpdateData(db, kycUpdateData)
	if err != nil {
		logData.Message = "AddKycUpdateData: Error updating kyc update data"
		s.LoggerService.LogError(logData)
		return err
	}

	// update onboarding status
	if err := models.UpdateUserOnboardingStatus(constants.KYC_CALLBACK_STEP, userData.UserId); err != nil {
		logData.Message = "GetVcipUrl: error while updating onboarding status"
		s.LoggerService.LogError(logData)
	}

	logData.Message = "AddKycUpdateData: Kyc update data inserted successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return nil
}

func (s *Store) GetKycUpdateData(ctx context.Context, authValues *models.AuthValues) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.KYC,
		RequestURI: "/api/kyc/get-update",
		Message:    "GetKycUpdateData log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	kycUpdateData, err := models.GetKycUpdateDataByUserId(s.db, authValues.UserId)

	if err != nil {
		logData.Message = "GetKycUpdateData: Error getting kyc update data by user id"
		s.LoggerService.LogError(logData)
		return nil, errors.New("kyc update data not found for this user")
	}

	responseData := map[string]interface{}{
		"status": kycUpdateData.Status,
	}

	jsonData, err := json.Marshal(responseData)

	if err != nil {
		logData.Message = "GetKycUpdateData: Error marshaling response data to json"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	encrypted, err := security.Encrypt(jsonData, []byte(authValues.Key))

	if err != nil {
		logData.Message = "GetKycUpdateData: Error encrypting response data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "GetKycUpdateData: Response data encrypted successfully"
	logData.ResponseSize = len(jsonData)
	logData.ResponseBody = string(jsonData)
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return encrypted, nil
}
