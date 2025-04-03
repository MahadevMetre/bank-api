package nominee

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"bankapi/constants"
	"bankapi/models"
	"bankapi/requests"
	"bankapi/responses"
	"bankapi/security"
	"bankapi/services"
	"bankapi/utils"

	"bitbucket.org/paydoh/paydoh-commons/database"
	commonSrv "bitbucket.org/paydoh/paydoh-commons/services"
	"bitbucket.org/paydoh/paydoh-commons/types"
)

type Nominee interface {
	AddNominee(authValues *models.AuthValues, request *requests.AddNomineeRequest) (interface{}, error)
}

type Store struct {
	db            *sql.DB
	m             *database.Document
	memory        *database.InMemory
	bankService   *services.BankApiService
	LoggerService *commonSrv.LoggerService
	auditLogSrv   services.AuditLogService
}

func NewStore(
	log *commonSrv.LoggerService,
	db *sql.DB,
	m *database.Document,
	memory *database.InMemory,
	auditLogSrv services.AuditLogService,
) *Store {
	bankService := services.NewBankApiService(log, memory)
	return &Store{
		db:            db,
		m:             m,
		memory:        memory,
		LoggerService: log,
		bankService:   bankService,
		auditLogSrv:   auditLogSrv,
	}
}

func (s *Store) AddNominee(ctx context.Context, authValues *models.AuthValues, request *requests.AddNomineeRequest) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.NOMINEE,
		RequestURI: "/api/nominee/add-nominee",
		Message:    "AddNominee log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	userData, err := models.GetUserAndAccountDetailByUserID(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "AddNominee: Error getting user data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	existingNominee, err := models.FindOneNomineeByUserID(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "AddNominee: Error fetching existing nominee from the database"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	txnId := ""
	nomApplId := ""
	if existingNominee != nil {
		if existingNominee.IsActive && strings.ToLower(request.NomReqType) == "new" {
			logData.Message = "AddNominee: Nominee already exists. Please try updating."
			s.LoggerService.LogError(logData)
			return nil, errors.New("nominee already exists please try updating")
		}

		if !existingNominee.TxnIdentifier.Valid {
			txnId = ""
		} else {
			txnId = existingNominee.TxnIdentifier.String
		}

		if !existingNominee.IsOtpSent && !existingNominee.IsActive {
			nomApplId = existingNominee.NomApplicantID.String
		}

		if existingNominee.IsOtpSent && !existingNominee.IsActive {
			nomApplId = existingNominee.NomApplicantID.String
		}
	}

	nomineesResponse, err := s.FetchNominee(ctx, authValues)
	if err != nil {
		logData.Message = fmt.Sprintf("FetchNominee: Error calling bank service for fetching nominee %v", err)
		s.LoggerService.LogError(logData)
	}

	if nomineesResponse != nil {
		str, ok := nomineesResponse.(string)
		if !ok {
			logData.Message = "AddNominee: Error converting nomineesResponse to string"
			s.LoggerService.LogError(logData)
			return nil, errors.New("invalid response type")
		}

		DecryptedResponse, err := security.Decrypt(str, []byte(authValues.Key))
		if err != nil {
			logData.Message = "AddNominee: Error decrypting response"
			s.LoggerService.LogError(logData)
			return nil, err
		}

		var nomineesResponseDetails responses.FetchNomineeResponseData
		err = json.Unmarshal([]byte(DecryptedResponse), &nomineesResponseDetails)
		if err != nil {
			logData.Message = "AddNominee: Error unmarshaling nominee response"
			s.LoggerService.LogError(logData)
			return nil, err
		}

		if nomineesResponseDetails.NomName == request.NomName &&
			nomineesResponseDetails.NomRelation == request.NomRelation &&
			nomineesResponseDetails.NomDOB == request.NomDOB {
			logData.Message = "AddNominee: Nominee with the same name, relation, and DOB already exists"
			s.LoggerService.LogError(logData)
			return nil, errors.New("nominee with the same name, relation, and date of birth already exists")
		}
	}

	rel, err := models.FindNomineeByNomineeRelation(s.db, request.NomRelation)
	if err != nil {
		logData.Message = "AddNominee: Error finding nominee relation"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	outgoingRequest := requests.NewOutgoingAddNomineeRequest()

	if err := outgoingRequest.Bind(userData.Applicant_id, nomApplId, userData.AccountNumber, rel.NomineeCode, request, txnId); err != nil {
		logData.Message = "AddNominee: Error binding outgoing request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	nominee := models.Nominee{
		AccountDataID:       userData.Id,
		UserId:              authValues.UserId,
		NomApplicantID:      types.FromString(outgoingRequest.NomApplId),
		NomReqType:          types.FromString(outgoingRequest.NomReqType),
		TxnIdentifier:       types.FromString(outgoingRequest.TxnIdentifier),
		NomName:             types.FromString(outgoingRequest.NomName),
		DateOfBirth:         types.FromString(outgoingRequest.NomDOB),
		Relation:            types.FromString(outgoingRequest.NomRelation),
		Address1:            types.FromString(outgoingRequest.NomAddressL1),
		Address2:            types.FromString(outgoingRequest.NomAddressL2),
		Address3:            types.FromString(outgoingRequest.NomAddressL3),
		City:                types.FromString(outgoingRequest.NomCity),
		Pincode:             types.FromString(outgoingRequest.NomZipcode),
		NomineeMobileNumber: types.FromString(request.NomineeMobileNumber),
		IsVerified:          false,
		IsActive:            false,
		IsOtpSent:           false,
	}

	if existingNominee == nil {
		_, err := models.InsertNominee(nominee)
		if err != nil {
			logData.Message = "AddNominee: Error while Updating Nomineedetail in DB " + err.Error()
			s.LoggerService.LogError(logData)
			return nil, err
		}
	} else {
		err = models.UpdateNomineeByUserId(s.db, &nominee)
		if err != nil {
			logData.Message = "AddNominee: Error while Updating Nomineedetail in DB " + err.Error()
			s.LoggerService.LogError(logData)
			return nil, err
		}
	}

	var response *responses.OtpGenerationResponse
	var opErr error

	response, opErr = s.bankService.CreateAddNominee(ctx, outgoingRequest)
	if opErr != nil {
		bankErr := s.bankService.HandleBankSpecificError(opErr, func(errorCode string) (string, bool) {
			return constants.GetNomineeErrorMessage(errorCode)
		})

		if bankErr != nil {
			logData.Message = fmt.Sprintf("Bank error encountered (ErrorCode: %s)", bankErr.ErrorCode)
			s.LoggerService.LogError(logData)

			if errorMessage, exists := constants.GetNomineeErrorMessage(bankErr.ErrorCode); exists {
				return nil, errors.New(errorMessage)
			}

			if msg, retryable := constants.GetNomineeErrorRetryMessages(bankErr.ErrorCode); retryable {
				err := utils.RetryFunc(func() error {
					response, opErr = s.bankService.CreateAddNominee(ctx, outgoingRequest)
					return opErr
				}, 2)

				if err != nil {
					logData.Message = "AddNominee: CreateAddNominee failed after retries"
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

	if response.ErrorCode != "0" && response.ErrorCode != "00" {
		logData.Message = "AddNominee: Error in response"
		s.LoggerService.LogError(logData)
		return nil, errors.New(response.ErrorMessage)
	}

	responseBytes, err := response.Marshal()

	if err != nil {
		logData.Message = "AddNominee: Error marshalling response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	data := map[string]string{
		"nomApplId":     outgoingRequest.NomApplId,
		"transaction":   response.TxnIdentifier,
		"applicantId":   userData.Applicant_id,
		"accountNumber": userData.AccountNumber,
	}

	if err := utils.SaveHashMapToRedis(ctx, s.memory, authValues.UserId, constants.NomineeKey, data, time.Minute*5); err != nil {
		logData.Message = "AddNominee: Error while Edding Transaction Id Cache"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	nomineeUpdate := models.Nominee{
		UserId:    authValues.UserId,
		IsOtpSent: true,
	}

	if err = models.UpdateNomineeByUserId(s.db, &nomineeUpdate); err != nil {
		logData.Message = "AddNominee: Error while updating Data for OTP sent"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	encrypted, err := security.Encrypt(responseBytes, []byte(authValues.Key))

	if err != nil {
		logData.Message = "AddNominee: Error encrypting response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "AddNominee: Response encrypted successfully"
	logData.ResponseSize = len(responseBytes)
	logData.ResponseBody = string(responseBytes)
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return encrypted, nil
}

func (s *Store) VerifyNomineeOtp(ctx context.Context, authValues *models.AuthValues, request *requests.VerifyOtpRequest) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.NOMINEE,
		RequestURI: "/api/nominee/verify-otp",
		Message:    "VerifyNomineeOtp log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	outgoingRequest := requests.NewOutgoingVerifyNomineeOTP()

	// get nomineApplId from cache
	clientData, err := utils.GetHashDataFromRedis(ctx, s.memory, authValues.UserId, constants.NomineeKey)
	if err != nil {
		logData.Message = "VerifyNomineeOtp: Error getting data from Redis"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	nomApplicantId := clientData["nomApplId"]
	txnData := clientData["transaction"]
	accountNumber := clientData["accountNumber"]
	applicantId := clientData["applicantId"]

	if err := outgoingRequest.Bind(nomApplicantId, applicantId, accountNumber, request); err != nil {
		logData.Message = "VerifyNomineeOtp: Error binding outgoing request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	outgoingRequest.TxnIdentifier = txnData

	var response *responses.OtpAuthenticationResponse
	var opErr error

	response, opErr = s.bankService.VerifyNomineeOTP(ctx, outgoingRequest)
	if opErr != nil {
		bankErr := s.bankService.HandleBankSpecificError(opErr, func(errorCode string) (string, bool) {
			return constants.GetNomineeErrorMessage(errorCode)
		})

		if bankErr != nil {
			logData.Message = fmt.Sprintf("Bank error encountered (ErrorCode: %s)", bankErr.ErrorCode)
			s.LoggerService.LogError(logData)

			if errorMessage, exists := constants.GetNomineeErrorMessage(bankErr.ErrorCode); exists {
				return nil, errors.New(errorMessage)
			}

			if msg, retryable := constants.GetNomineeErrorRetryMessages(bankErr.ErrorCode); retryable {
				err := utils.RetryFunc(func() error {
					response, opErr = s.bankService.VerifyNomineeOTP(ctx, outgoingRequest)
					return opErr
				}, 2)

				if err != nil {
					logData.Message = "VerifyNomineeOtp failed after retries"
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

	if response.ErrorCode != "0" && response.ErrorCode != "00" {
		if response.ErrorCode == "90" || strings.Contains(strings.ToLower(response.ErrorMessage), "otp") {
			nomineeUpdate := models.Nominee{
				UserId:         authValues.UserId,
				NomCBSStatus:   types.FromString(response.NomCBSStatus),
				NomUpdatedTime: types.FromString(response.NomUpdateDtTime),
				IsVerified:     false,
			}
			if err = models.UpdateNomineeByUserId(s.db, &nomineeUpdate); err != nil {
				logData.Message = "AddNominee: Error while updating Data for OTP sent"
				s.LoggerService.LogError(logData)
				return nil, err
			}
		}
		logData.Message = "VerifyNomineeOtp: Error in response"
		s.LoggerService.LogError(logData)
		return nil, errors.New(response.ErrorMessage)
	}

	// update nominee
	nomineeUpdate := models.Nominee{
		UserId:         authValues.UserId,
		NomCBSStatus:   types.FromString(response.NomCBSStatus),
		NomUpdatedTime: types.FromString(response.NomUpdateDtTime),
		IsVerified:     true,
		IsActive:       true,
		IsOtpSent:      true,
	}

	if err = models.UpdateNomineeByUserId(s.db, &nomineeUpdate); err != nil {
		logData.Message = "AddNominee: Error while updating Data for OTP sent"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	responseBytes, err := response.Marshal()
	if err != nil {
		logData.Message = "VerifyNomineeOtp: Error marshalling response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	encrypted, err := security.Encrypt(responseBytes, []byte(authValues.Key))

	if err != nil {
		logData.Message = "VerifyNomineeOtp: Error encrypting response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	requestBytes, err := outgoingRequest.Marshal()
	if err != nil {
		logData.Message = "VerifyNomineeOtp: Error marshalling response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	encryptedReq, err := security.Encrypt(requestBytes, []byte(authValues.Key))

	if err != nil {
		logData.Message = "VerifyNomineeOtp: Error encrypting response"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	var action string
	if strings.ToLower(request.ReqType) == "new" {
		action = constants.NEW_NOMINEE
	} else {
		action = constants.UPDATE_NOMINEE
	}

	// save audit log
	if err := s.auditLogSrv.Save(ctx, &services.AuditLog{
		TransactionID:  outgoingRequest.TxnIdentifier,
		UserID:         authValues.UserId,
		ApplicantID:    applicantId,
		RequestURL:     "/api/nominee/verify-otp",
		HTTPMethod:     "POST",
		ResponseStatus: 200,
		Action:         action,
		RequestBody:    encryptedReq,
	}); err != nil {
		logData.Message = "VerifyNomineeOtp: Error saving audit log"
		s.LoggerService.LogError(logData)
	}

	if err := models.GenerateNotification(authValues.UserId, "nominee_addition", utils.CalculateTimeDifference(response.NomUpdateDtTime), "nominee_addition"); err != nil {
		logData.Message = "VerifyNomineeOtp: Error while Generating Notification " + err.Error()
		s.LoggerService.LogError(logData)
	}

	if err := s.memory.Delete(fmt.Sprintf(constants.NomineeKey, authValues.UserId)); err != nil {
		logData.Message = "VerifyNomineeOtp:Failed to delete Transaction Id from cache"
		s.LoggerService.LogError(logData)
	}

	logData.Message = "VerifyNomineeOtp: Response encrypted successfully"
	logData.ResponseSize = len(responseBytes)
	logData.ResponseBody = string(responseBytes)
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return encrypted, nil
}

func (s *Store) FetchNominee(ctx context.Context, authValues *models.AuthValues) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.NOMINEE,
		RequestURI: "/api/nominee/fetch",
		Message:    "FetchNominee log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	userData, err := models.GetUserAndAccountDetailByUserID(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "FetchNominee: Error getting user data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	nomineeDetail, err := models.FindOneNomineeByUserID(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "FetchNominee: Error fetching nominee details from database"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	var transactionID string
	if nomineeDetail != nil {
		transactionID = nomineeDetail.TxnIdentifier.String
	}

	outgoingRequest := requests.NewOutgoingFetchNomineeRequest()

	if err := outgoingRequest.Bind(userData.Applicant_id, userData.AccountNumber, transactionID); err != nil {
		logData.Message = "FetchNominee: Error binding outgoing request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	response, err := s.bankService.FetchNominee(ctx, outgoingRequest)
	if err != nil {
		logData.Message = fmt.Sprintf("FetchNominee: Error calling bank service for fetching nominee %v", err.Error())
		s.LoggerService.LogError(logData)
	}

	if nomineeDetail == nil && response != nil {
		nominee := models.Nominee{
			AccountDataID:  userData.Id,
			UserId:         authValues.UserId,
			NomApplicantID: types.FromString(response.NomAppID),
			NomReqType:     types.FromString(response.NomReqType),
			DateOfBirth:    types.FromString(response.NomDOB),
			Relation:       types.FromString(response.NomRelation),
			Address1:       types.FromString(response.NomAddressL1),
			Address2:       types.FromString(response.NomAddressL2),
			Address3:       types.FromString(response.NomAddressL3),
			City:           types.FromString(response.NomCity),
			Pincode:        types.FromString(response.NomZipcode),
			NomCBSStatus:   types.FromString(response.NomCBSStatus),
			NomUpdatedTime: types.FromString(response.NomUpdateDtTime),
			IsVerified:     true,
			IsActive:       true,
			IsOtpSent:      true,
		}

		_, err = models.InsertNominee(nominee)
		if err != nil {
			logData.Message = "FetchNominee: Error inserting new nominee"
			s.LoggerService.LogError(logData)
			return nil, err
		}
	}

	respData := responses.NewFetchNomineeData()

	if nomineeDetail != nil {
		respData.AccountNo = userData.AccountNumber
		respData.ApplicantID = userData.Applicant_id
		respData.NomAppID = nomineeDetail.NomApplicantID.String
		respData.NomAddressL1 = nomineeDetail.Address1.String
		respData.NomAddressL2 = nomineeDetail.Address2.String
		respData.NomAddressL3 = nomineeDetail.Address3.String
		respData.NomCity = nomineeDetail.City.String
		respData.NomCountry = "India"
		respData.NomDOB = nomineeDetail.DateOfBirth.String
		respData.NomName = nomineeDetail.NomName.String
		respData.NomReqType = nomineeDetail.NomReqType.String
		respData.NomState = ""
		respData.NomZipcode = nomineeDetail.Pincode.String
		respData.NomRelation = nomineeDetail.Relation.String
		respData.GuardianAddressL1 = ""
		respData.GuardianAddressL2 = ""
		respData.GuardianAddressL3 = ""
		respData.GuardianCity = ""
		respData.GuardianCountry = ""
		respData.GuardianName = ""
		respData.GuardianNomRelation = ""
		respData.GuardianState = ""
		respData.GuardianZipcode = ""
		respData.NomineeUpdateTime = utils.DateFormat(nomineeDetail.NomUpdatedTime.String)
		respData.NomineeActive = nomineeDetail.IsActive
	} else {
		if response != nil {
			if err := respData.Bind(response); err != nil {
				logData.Message = "FetchNominee: Error binding response"
				s.LoggerService.LogError(logData)
				return nil, err
			}
		} else {
			logData.Message = "FetchNominee: Received nil response from bank service"
			s.LoggerService.LogError(logData)
			return nil, errors.New("No Data Found")
		}
	}

	byteudd, err := respData.Marshal()
	if err != nil {
		logData.Message = "FetchNominee: Error marshalling response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	encrypted, err := security.Encrypt(byteudd, []byte(authValues.Key))
	if err != nil {
		logData.Message = "FetchNominee: Error encrypting response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "FetchNominee: Response encrypted successfully"
	logData.ResponseBody = string(byteudd)
	logData.ResponseSize = len(byteudd)
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return encrypted, nil
}
