package onboarding

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"bitbucket.org/paydoh/paydoh-commons/database"
	commonSrv "bitbucket.org/paydoh/paydoh-commons/services"
	"bitbucket.org/paydoh/paydoh-commons/types"

	"bankapi/constants"
	"bankapi/models"
	"bankapi/requests"
	"bankapi/responses"
	"bankapi/security"
	"bankapi/services"
	"bankapi/stores/authorization"
	"bankapi/stores/mail"
	"bankapi/utils"
)

type Onboarding interface {
	AddPersonalInformation(userId, key, deviceIp, os, osVersion string, request *requests.PersonalInformationRequest)
	GetPersonalInformation(userId, key, deviceIp, os, osVersion string)
}

type Store struct {
	db                  *sql.DB
	m                   *database.Document
	memory              *database.InMemory
	bankService         *services.BankApiService
	notificationService *services.NotificationService
	LoggerService       *commonSrv.LoggerService
	AuthStore           *authorization.AuthorizationStore
}

func NewStore(
	log *commonSrv.LoggerService,
	db *sql.DB,
	m *database.Document,
	memory *database.InMemory,
	authStore *authorization.AuthorizationStore,
) *Store {
	bankService := services.NewBankApiService(log, memory)
	notificationService := services.NewNotificationService()
	return &Store{
		db:                  db,
		memory:              memory,
		m:                   m,
		bankService:         bankService,
		notificationService: notificationService,
		LoggerService:       log,
		AuthStore:           authStore,
	}
}

func (s *Store) AddPersonalInformation(ctx context.Context, authValues *models.AuthValues, request *requests.PersonalInformationRequest) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.ONBOARDING,
		RequestURI: "/api/onboarding/personal-information",
		Message:    "AddPersonalInformation log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	existingPersonalInformation, err := models.GetPersonalInformation(s.db, authValues.UserId)

	if err != nil {
		logData.Message = "AddPersonalInformation: Error getting personal information"
		s.LoggerService.LogError(logData)
		if errors.Is(err, constants.ErrNoDataFound) {
			personalInformation := models.NewPersonalInformation()
			personalInformation.Bind(request, authValues.UserId)

			if err := models.InsertPersonalInformation(s.db, personalInformation); err != nil {
				logData.Message = "AddPersonalInformation: Error inserting personal information"
				s.LoggerService.LogError(logData)
				return nil, err
			}

			// update onboarding status
			if err := models.UpdateUserOnboardingStatus(constants.PERSONAL_DETAILS_STEP, authValues.UserId); err != nil {
				logData.Message = "AddPersonalInformation: error while updating onboarding status"
				s.LoggerService.LogError(logData)
			}

			updateData, err := utils.StructToMap(request)
			if err != nil {
				fmt.Println("Error converting struct to map:", err)
			}

			referralCode, err := utils.GenerateReferralCode(request.FirstName)
			if err != nil {
				fmt.Println("error while generating referral")
			}
			updateData["referral_code"] = referralCode

			if err = s.AuthStore.UpdateUserInMongoDB(authValues.UserId, updateData); err != nil {
				return nil, err
			}

			personalInfoData, err := personalInformation.Marshal()
			if err != nil {
				logData.Message = "AddPersonalInformation: Error marshalling personal information"
				s.LoggerService.LogError(logData)
				return nil, err
			}

			encrypted, err := security.Encrypt(personalInfoData, []byte(authValues.Key))
			if err != nil {
				logData.Message = "AddPersonalInformation: Error encrypting personal information"
				s.LoggerService.LogError(logData)
				return nil, err
			}

			return encrypted, nil
		}

		return nil, err
	}

	updateModel := models.NewPersonalInformation()

	if err := updateModel.BindUpdate(request, existingPersonalInformation); err != nil {
		logData.Message = "AddPersonalInformation: Error binding personal information"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if err := models.UpdatePersonalInformation(s.db, updateModel, authValues.UserId); err != nil {
		if err.Error() == "nothing to update" {
			// update onboarding status
			if err := models.UpdateUserOnboardingStatus(constants.PERSONAL_DETAILS_STEP, authValues.UserId); err != nil {
				logData.Message = "AddPersonalInformation: error while updating onboarding status"
				s.LoggerService.LogError(logData)
			}
			s.AddUserDetailInMonogDB(authValues.UserId, request)
			return nil, nil
		}

		logData.Message = "AddPersonalInformation: Error updating personal information"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	// update onboarding status
	if err := models.UpdateUserOnboardingStatus(constants.PERSONAL_DETAILS_STEP, authValues.UserId); err != nil {
		logData.Message = "AddPersonalInformation: error while updating onboarding status"
		s.LoggerService.LogError(logData)
	}
	s.AddUserDetailInMonogDB(authValues.UserId, request)
	personalInfoData, err := updateModel.Marshal()
	if err != nil {
		logData.Message = "AddPersonalInformation: Error marshalling personal information"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	encrypted, err := security.Encrypt(personalInfoData, []byte(authValues.Key))
	if err != nil {
		logData.Message = "AddPersonalInformation: Error encrypting personal information"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "AddPersonalInformation: Personal information encrypted successfully"
	logData.EndTime = time.Now()
	logData.ResponseSize = len(personalInfoData)
	logData.ResponseBody = string(personalInfoData)
	s.LoggerService.LogInfo(logData)

	return encrypted, nil
}

func (s *Store) AddUserDetailInMonogDB(userID string, request *requests.PersonalInformationRequest) error {
	updateData, err := utils.StructToMap(request)
	if err != nil {
		fmt.Println("Error converting struct to map:", err)
	}

	referralCode, err := utils.GenerateReferralCode(request.FirstName)
	if err != nil {
		fmt.Println("error while generating referral")
	}
	updateData["referral_code"] = referralCode

	if err = s.AuthStore.UpdateUserInMongoDB(userID, updateData); err != nil {
		return err
	}
	return nil
}

func (s *Store) GetPersonalInformation(ctx context.Context, authValues *models.AuthValues) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.ONBOARDING,
		RequestURI: "/api/onboarding/personal-information",
		Message:    "GetPersonalInformation log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	existingPersonalInformation, err := models.GetPersonalInformation(s.db, authValues.UserId)

	if err != nil {
		logData.Message = "GetPersonalInformation: Error getting personal information"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	personalInfoData, err := existingPersonalInformation.Marshal()

	if err != nil {
		logData.Message = "GetPersonalInformation: Error marshalling personal information"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	encrypted, err := security.Encrypt(personalInfoData, []byte(authValues.Key))

	if err != nil {
		logData.Message = "GetPersonalInformation: Error encrypting personal information"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "GetPersonalInformation: Personal information encrypted successfully"
	logData.EndTime = time.Now()
	logData.ResponseSize = len(personalInfoData)
	logData.ResponseBody = string(personalInfoData)
	s.LoggerService.LogInfo(logData)

	return encrypted, nil
}

func (s *Store) CreateBankAccount(ctx context.Context, authValues *models.AuthValues, request *requests.CreateBankAccountRequest) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.ONBOARDING,
		RequestURI: "/api/onboarding/create-account",
		Message:    "CreateBankAccount log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	existingUserData, err := models.GetUserDataByUserId(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "CreateBankAccount: Error getting user data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	// check if user has already account
	err = models.CheckAccountDataAvailability(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "CreateBankAccount: Error checking account data availability"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	demographicRequest := requests.NewOutGoingDemographicRequest()

	if err := demographicRequest.Bind(existingUserData.ApplicantId, existingUserData.MobileNumber); err != nil {
		logData.Message = "CreateBankAccount: Error binding demographic request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	demographicResponse, err := s.bankService.GetDemographicData(ctx, demographicRequest)
	if err != nil {
		bankErr := s.bankService.HandleBankSpecificError(err, func(errorCode string) (string, bool) {
			return constants.GetDemographicErrorMessage(errorCode)
		})

		if bankErr != nil {
			logData.Message = fmt.Sprintf("GetDemographicData: Bank error encountered (ErrorCode: %s)", bankErr.ErrorCode)
			s.LoggerService.LogError(logData)
			return nil, errors.New(bankErr.ErrorMessage)
		}

		return nil, err
	}

	personalInformation, err := models.GetPersonalInformation(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "CreateBankAccount: Error getting personal information"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	// check if previous steps has been completed!
	err = models.IsEligibleForAccountCreate(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "CreateBankAccount: Error checking eligibility for account creation"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	// update personal information with demographic data
	updatePI := models.NewPersonalInformation()

	if demographicResponse.MiddleName == "" || demographicResponse.LastName == "" {
		nameParts := utils.ParseName(demographicResponse.FirstName)

		if firstName, ok := nameParts["first_name"]; ok {
			updatePI.FirstName = firstName
		}

		if middleName, ok := nameParts["middle_name"]; ok {
			updatePI.MiddleName = middleName
		}

		if lastName, ok := nameParts["last_name"]; ok {
			updatePI.LastName = lastName
		}

	} else {
		updatePI.FirstName = demographicResponse.FirstName
		updatePI.MiddleName = demographicResponse.MiddleName
		updatePI.LastName = demographicResponse.LastName
	}

	updatePI.Email = demographicResponse.Root.UIDData.Poi.Email

	if demographicResponse.Root.Pincode == "" {
		updatePI.PinCode = types.NewNullableString(nil)
	} else {
		updatePI.PinCode = types.FromString(demographicResponse.Root.Pincode)
	}

	err = models.UpdatePersonalInformation(s.db, updatePI, authValues.UserId)
	if err != nil {
		logData.Message = "CreateBankAccount: Error updating personal information"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	// check if user are from given bank city data
	exists, err := models.CheckZipCodeExists(demographicResponse.Root.Pincode[:3])
	if err != nil {
		logData.Message = "CreateBankAccount: Error checking zip code exists"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if !exists {
		logData.Message = "CreateBankAccount: Account creation service is not available in this area"
		s.LoggerService.LogError(logData)
		return nil, errors.New("account creation service is not available in this area")
	}

	stateData, err := models.GetStateCodeByName(strings.ToLower(demographicResponse.Root.UIDData.Poa.State))
	if err != nil {
		logData.Message = "CreateBankAccount: Error getting state code"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	outGoingCreateBankAccountRequest := requests.NewOutgoingCreateAccountRequest()
	if err := outGoingCreateBankAccountRequest.Bind(
		existingUserData.UserId,
		existingUserData.ApplicantId,
		personalInformation.Email,
		personalInformation.MiddleName,
		existingUserData.MobileNumber,
		request,
		demographicResponse,
		personalInformation.LastName,
		personalInformation.DateOfBirth,
		stateData.StateCode,
		request.IsAddrSameAsAdhaar,
	); err != nil {
		logData.Message = "CreateBankAccount: Error binding outgoing create bank account request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	var response *responses.ImmediateCreateBankResponse
	var opErr error

	response, opErr = s.bankService.CreatebankAccount(ctx, outGoingCreateBankAccountRequest)
	if opErr != nil {
		bankErr := s.bankService.HandleBankSpecificError(opErr, func(errorCode string) (string, bool) {
			return constants.GetAccountCreationErrorMessage(errorCode)
		})

		if bankErr != nil {
			logData.Message = fmt.Sprintf("Bank error encountered (ErrorCode: %s)", bankErr.ErrorCode)
			s.LoggerService.LogError(logData)

			if errorMessage, exists := constants.GetAccountCreationErrorMessage(bankErr.ErrorCode); exists {
				return nil, errors.New(errorMessage)
			}

			if msg, retryable := constants.GetAccountCreationRetryErrorCode(bankErr.ErrorCode); retryable {
				err := utils.RetryFunc(func() error {
					response, opErr = s.bankService.CreatebankAccount(ctx, outGoingCreateBankAccountRequest)
					return opErr
				}, 2)

				if err != nil {
					logData.Message = "CreateBankAccount failed after retries"
					s.LoggerService.LogError(logData)
					return nil, errors.New(msg)
				}
			} else {
				return nil, errors.New(bankErr.ErrorMessage)
			}
		} else {
			return nil, opErr
		}
	}

	// // update onboarding status
	// if err := models.UpdateUserOnboardingStatus(constants.ACCOUNT_CREATION_STAGE, authValues.UserId); err != nil {
	// 	logData.Message = "CreateBankAccount: error while updating onboarding status"
	// 	s.LoggerService.LogError(logData)
	// }

	accountData, err := models.GetAccountDetailsV2(s.db, authValues.UserId)
	if err != nil {
		if errors.Is(err, constants.ErrNoDataFound) {
			var communicationAddr requests.CommunicationAddress
			if request.IsAddrSameAsAdhaar {
				communicationAddr.City = demographicResponse.Root.UIDData.Poa.Co
				communicationAddr.HouseNo = demographicResponse.Root.UIDData.Poa.House
				communicationAddr.Locality = demographicResponse.Root.Locality
				communicationAddr.StreetName = demographicResponse.Root.Street
				communicationAddr.Landmark = demographicResponse.Root.Landmark
				communicationAddr.PinCode = demographicResponse.Root.Pincode
				communicationAddr.State = demographicResponse.Root.UIDData.Poa.State
			} else {
				communicationAddr = *request.CommunicationAddress
			}
			if err := models.InsertNewAccount(
				s.db, authValues.UserId,
				outGoingCreateBankAccountRequest.ApplicationID,
				response.ServiceName,
				&communicationAddr,
				request.IsAddrSameAsAdhaar,
				request.MotherMaidenName,
				request.AnnualTurnOver,
				request.CustomerEducation,
				request.ProfessionCode,
				request.MaritalStatus,
			); err != nil {
				logData.Message = "CreateBankAccount: Error inserting newAccount"
				s.LoggerService.LogError(logData)
				return nil, err
			}
		}
	}

	if accountData != nil {
		accountUpdateData := &models.AccountDataUpdate{
			ApplicationId:        outGoingCreateBankAccountRequest.ApplicationID,
			CommunicationAddress: request.CommunicationAddress,
		}

		if accountData.IsAddrSameAsAdhaar != request.IsAddrSameAsAdhaar {
			accountUpdateData.IsAddrSameAsAdhaar = request.IsAddrSameAsAdhaar
		}

		if accountData.MotherMaidenName.String != request.MotherMaidenName {
			accountUpdateData.MotherMaidenName = request.MotherMaidenName
		}

		if accountData.AnnualTurnOver.String != request.AnnualTurnOver {
			accountUpdateData.AnnualTurnOver = request.AnnualTurnOver
		}

		if accountData.CustomerEducation.String != request.CustomerEducation {
			accountUpdateData.CustomerEducation = request.CustomerEducation
		}

		if accountData.ProfessionCode.String != request.ProfessionCode {
			accountUpdateData.ProfessionCode = request.ProfessionCode
		}

		if accountData.MaritalStatus.String != request.MaritalStatus {
			accountUpdateData.MaritalStatus = request.MaritalStatus
		}

		if err := models.UpdateAccountByUserId(accountUpdateData, accountData.UserId); err != nil {
			return nil, err
		}
	}

	shippingAddress := models.NewShippingAddress()
	shippingAddress.UserId = existingUserData.UserId
	shippingAddress.AddressLine1 = demographicResponse.Root.UIDData.Poa.Co + ", " + demographicResponse.Root.UIDData.Poa.House
	shippingAddress.Locality = demographicResponse.Root.Locality
	shippingAddress.StreetName = demographicResponse.Root.Street
	shippingAddress.Landmark = sql.NullString{
		String: demographicResponse.Root.Landmark,
		Valid:  demographicResponse.Root.Landmark != "",
	}

	shippingAddress.City = demographicResponse.Root.Vtc
	shippingAddress.State = demographicResponse.Root.UIDData.Poa.State
	shippingAddress.PinCode = demographicResponse.Root.Pincode
	shippingAddress.Country = request.CountryResidence
	shippingAddress.Document = ""
	shippingAddress.DocumentType = "Aadhard"

	if err := models.InsertShippingAddress(s.db, shippingAddress); err != nil {
		logData.Message = "CreateBankAccount: Error inserting shipping address"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	responseByte, err := response.Marshal()
	if err != nil {
		logData.Message = "CreateBankAccount: Error marshalling response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	encrypted, err := security.Encrypt(responseByte, []byte(authValues.Key))
	if err != nil {
		logData.Message = "CreateBankAccount: Error encrypting response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "CreateBankAccount: Create Bank Account Response encrypted"
	logData.EndTime = time.Now()
	logData.ResponseSize = len(responseByte)
	logData.ResponseBody = string(responseByte)
	s.LoggerService.LogInfo(logData)

	return encrypted, nil
}

func (s *Store) UpdateBankAccount(ctx context.Context, request *requests.AccountCreateCallbackRequest, metadataReq *requests.AccountCreateCallbackEncryptedRequest) error {
	logData := &commonSrv.LogEntry{
		Action:     constants.ONBOARDING,
		RequestURI: "/api/callback/account-create",
		Message:    "UpdateBankAccount log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	updateData := &models.AccountDataUpdate{
		SourcedBy:     metadataReq.SourcedBy,
		ProductType:   metadataReq.ProductType,
		ServiceName:   request.ServiceName,
		ApplicationId: request.ApplicationId,
		CustomerId:    request.CBSStatus[0].CustomerID,
		AccountNumber: request.CBSStatus[0].AccountNo,
		CallbackName:  metadataReq.CallbackName,
		Status:        request.CBSStatus[0].Status,
	}

	userID, err := s.getUserID(request)
	if err != nil {
		logData.Message = "UpdateBankAccount: Error getting user data"
		s.LoggerService.LogError(logData)
		return err
	}

	isUpdated := models.IsAccountAlreadyUpdated(s.db, request.ApplicationId)
	if isUpdated {
		logData.Message = "UpdateBankAccount: Account data already updated"
		s.LoggerService.LogError(logData)
		return errors.New("account data already updated")
	}

	err = models.UpdateAccount(s.db, updateData, request.ApplicationId)
	if err != nil {
		logData.Message = "UpdateBankAccount: Error updating account"
		s.LoggerService.LogError(logData)
		return err
	}

	if strings.ToLower(updateData.Status) != "success" {
		return errors.New("account creation callback failed")
	}

	// update onboarding status
	if err := models.UpdateUserOnboardingStatus(constants.ACCOUNT_CALLBACK_STEP, userID); err != nil {
		logData.Message = "CreateBankAccount: error while updating onboarding status"
		s.LoggerService.LogError(logData)
	}

	_, err = mail.NewStore(s.LoggerService, s.db, s.m, s.memory).SendAccountInformation(ctx, userID)
	if err != nil {
		logData.Message = "UpdateBankAccount: Error while sending email" + err.Error()
		s.LoggerService.LogError(logData)
	}

	logData.Message = "UpdateBankAccount: Account updated successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return nil
}

// getUserID retrieves the user ID either by applicant ID or application ID
func (s *Store) getUserID(request *requests.AccountCreateCallbackRequest) (string, error) {
	userData, err := models.GetUserDataByApplicantId(s.db, request.CBSStatus[0].ApplicantID)
	if err == nil {
		return userData.UserId, nil
	}

	if !errors.Is(err, constants.ErrUserNotFound) {
		return "", err
	}

	accountData, err := models.GetAccountDataByApplicationId(s.db, request.ApplicationId)
	if err != nil {
		return "", err
	}

	return accountData.UserId, nil
}

func (s *Store) GetAccountDetails(ctx context.Context, authValues *models.AuthValues) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.ONBOARDING,
		RequestURI: "/api/onboarding/get-account-details",
		Message:    "GetAccountDetails log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	resp, err := models.GetAccountDetailsV2(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "GetAccountDetails: Error getting account details"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	type ProfilePic struct {
		FileName string `json:"file_name"`
		FileSize string `json:"file_size"`
		TempName string `json:"temp_name"`
	}

	var profilePicUrl string
	if len(resp.ProfilePic) > 0 {
		var profilePic ProfilePic
		if err := json.Unmarshal(resp.ProfilePic, &profilePic); err != nil {
			logData.Message = "GetAccountDetails: Error unmarshaling profile pic"
			s.LoggerService.LogError(logData)
			return nil, err
		}

		if profilePic.TempName != "" {
			profilePicUrl = fmt.Sprintf("%s/profile/%s/%s",
				constants.AWSCloudFrontURL,
				authValues.UserId,
				profilePic.TempName)
		}
	}

	data := map[string]interface{}{
		"user_id":                      resp.UserId,
		"account_number":               resp.AccountNumber.String,
		"customer_id":                  resp.CustomerId.String,
		"upi_id":                       resp.UpiID.String,
		"status":                       resp.Status.String,
		"first_name":                   resp.FirstName,
		"middle_name":                  resp.MiddleName,
		"last_name":                    resp.LastName,
		"is_email_verified":            resp.IsEmailVerified,
		"is_account_detail_email_sent": resp.IsAccountDetailMailSent,
		"profile_pic_url":              profilePicUrl,
		"ifsc_code":                    constants.IfscCode,
		"bank_name":                    constants.BankName,
		"branch_name":                  constants.BranchName,
		"email":                        resp.Email,
		"gender":                       resp.Gender,
	}

	respData, err := json.Marshal(data)
	if err != nil {
		logData.Message = "GetAccountDetails: Error marshaling response data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	encrypted, err := security.Encrypt(respData, []byte(authValues.Key))
	if err != nil {
		logData.Message = "GetAccountDetails: Error encrypting response data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "GetAccountDetails: Response data encrypted"
	logData.EndTime = time.Now()
	logData.ResponseSize = len(respData)
	logData.ResponseBody = string(respData)
	s.LoggerService.LogInfo(logData)

	return encrypted, nil
}

func (s *Store) UserOnboardingStatus(ctx context.Context, authValues *models.AuthValues) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.ONBOARDING,
		RequestURI: "/api/onboarding/user-status",
		Message:    "UserOnboardingStatus log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	personalInformation, err := models.GetPersonalInformation(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "UserOnboardingStatus: Error getting personal information"
		s.LoggerService.LogError(logData)
	}

	resp, err := models.GetUserOnboardingStatus(authValues.UserId)
	if err != nil {
		logData.Message = "UserOnboardingStatus: Error getting onboarding details"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if personalInformation != nil && personalInformation.FirstName != "" && resp.CurrentStepName == "PERSONAL_DETAILS" {
		resp.CurrentStepName = ""
	}

	logData.Message = "UserOnboardingStatus: Response data fetched"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return resp, nil
}

func (s *Store) GetPincodeDetails(ctx context.Context, authValues *models.AuthValues, request *requests.PincodeDetails) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.ONBOARDING,
		RequestURI: "/api/onboarding/get-pincode-details",
		Message:    "GetPincodeDetails log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	data, err := models.GetCitiesByZipCode(request.Pincode[:3])
	if err != nil {
		logData.Message = "GetPincodeDetails: Error getting pincode details"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if len(data) == 0 {
		return nil, errors.New("kvb bank branch is not avaialble at this pincode")
	}

	logData.Message = "GetPincodeDetails: Response data fetched"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return data, nil

}
