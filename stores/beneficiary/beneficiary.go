package beneficiary

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
	"bankapi/utils"
)

type Beneficiary interface {
	FetchBeneficiary(authValues *models.AuthValues) (interface{}, error)
}

type Store struct {
	db            *sql.DB
	memory        *database.InMemory
	m             *database.Document
	redis         *database.InMemory
	bankService   *services.BankApiService
	LoggerService *commonSrv.LoggerService
}

func NewStore(log *commonSrv.LoggerService, db *sql.DB, m *database.Document, memory *database.InMemory) *Store {
	return &Store{
		db:            db,
		memory:        memory,
		m:             m,
		redis:         memory,
		LoggerService: log,
		bankService:   services.NewBankApiService(log, memory),
	}
}

func (s *Store) FetchBeneficiary(ctx context.Context, authValues *models.AuthValues) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.BENEFICIARY,
		RequestURI: "api/beneficiary/search",
		Message:    "FetchBeneficiary log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	existingDevice, err := models.GetUserDataByUserId(s.db, authValues.UserId)

	if err != nil {
		logData.Message = "FetchBeneficiary: Error getting user data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	existingAccount, err := models.GetAccountDataByUserId(s.db, existingDevice.UserId)

	if err != nil {
		logData.Message = "FetchBeneficiary: Error getting account data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	request := requests.NewOutgoingBeneficiarySearchRequest()

	if err := request.Bind(existingDevice.ApplicantId, existingAccount.AccountNumber); err != nil {
		logData.Message = "FetchBeneficiary: Error binding beneficiary search request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	response, err := s.bankService.GetBeneficiaries(ctx, request)

	if err != nil {
		logData.Message = "FetchBeneficiary: Error calling bank service"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if response == nil {
		logData.Message = "FetchBeneficiary: Error getting beneficiary error"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	// Scenario 13 If user select inactive beneficiary , data should be pre populated in the beneficiary form
	data, err := models.FindBeneficiariesByUserId(s.db, existingDevice.UserId)
	if err != nil {
		logData.Message = "FetchBeneficiary: Error getting existing beneficiaries from database"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	// Scenario 13 If user select inactive beneficiary , data should be pre populated in the beneficiary form
	// and also sending to this bind function from bank response and database response both
	b := models.NewBeneficiaryDetail()
	beneficiaries, err := b.Bind(response, data)
	if err != nil {
		logData.Message = "FetchBeneficiary: Error binding beneficiary response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	// benfIds := make([]any, len(beneficiaries))

	// for i := 0; i < len(beneficiaries); i++ {
	// 	benfIds[i] = beneficiaries[i].BenfId
	// }

	// benfDtos, err := models.GetBeneficiariesIn(s.db, benfIds)

	// if err != nil {
	// 	return nil, err
	// }

	byteudd, err := json.Marshal(beneficiaries)

	if err != nil {
		logData.Message = "FetchBeneficiary: Error marshaling beneficiary response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	encrypted, err := security.Encrypt(byteudd, []byte(authValues.Key))

	if err != nil {
		logData.Message = "FetchBeneficiary: Error encrypting beneficiary response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "FetchBeneficiary: Beneficiary response encrypted successfully"
	logData.ResponseSize = len(byteudd)
	logData.ResponseBody = string(byteudd)
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return encrypted, nil
}

func (s *Store) AddBeneficiary(ctx context.Context, authValues *models.AuthValues, r *requests.AddNewBeneficiary) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.BENEFICIARY,
		RequestURI: "api/beneficiary/add-beneficiary",
		Message:    "AddBeneficiary log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	existingDevice, err := models.GetUserDataByUserId(s.db, authValues.UserId)

	if err != nil {
		logData.Message = "AddBeneficiary: Error getting user data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if strings.ToLower(r.PaymentMode) == "ift" && r.BenfIFSC[:4] != "KVBL" {
		return nil, errors.New("invalid ifsc code for ift payment mode")
	}

	account, err := models.GetAccountDetails(s.db, existingDevice.UserId)
	if err != nil {
		logData.Message = "AddBeneficiary: Error while getting benificary account details"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if r.BenfAcctNo == account.AccountNumber {
		return nil, errors.New("self transfer not allowed")
	}

	existingBeneficiaries, err := models.FindBeneficiaryByNameAndIfscCode(s.db,
		authValues.UserId,
		r.BenfName,
		r.BenfIFSC,
		r.BenfNickName,
		r.BenfMobNo,
		r.BenfAcctNo,
	)

	if err != nil && !errors.Is(err, constants.ErrNoDataFound) {
		logData.Message = "AddBeneficiary: FindBeneficiaryByNameAndIfscCode Error fetching beneficiaries from the database"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if existingBeneficiaries != nil && existingBeneficiaries.IsActive {
		logData.Message = "AddBeneficiary: Beneficiary nickname already exists, please enter a new beneficiary nickname"
		s.LoggerService.LogError(logData)
		return nil, errors.New("Beneficiary nickname already exists, please enter a new beneficiary nickname")
	}

	_, err = models.GetIFSCData(s.db, r.BenfIFSC)
	if err != nil {
		if errors.Is(err, constants.ErrNoDataFound) {
			logData.Message = "AddBeneficiary: Invalid IFSC code provided"
			s.LoggerService.LogError(logData)
			return nil, errors.New(constants.InputErrorMessage)
		}

		logData.Message = "AddBeneficiary: Error fetching IFSC data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	// FetchBankBeneficiary is internal function to get the response from bank api.
	// we need to check if data is not there in our db if not we need to insert to out db
	beneficiariesResponse, err := s.FetchBankBeneficiary(ctx, authValues)
	if err != nil {
		if !errors.Is(err, constants.ErrNoDataFound) {
			logData.Message = "AddBeneficiary: Error fetching beneficiary data"
			s.LoggerService.LogError(logData)
			return nil, err
		}
	} else {

		str, ok := beneficiariesResponse.(string)
		if !ok {
			logData.Message = "AddBeneficiary: Error converting beneficiariesResponse to string"
			s.LoggerService.LogError(logData)
			return nil, errors.New("invalid response type")
		}

		DecryptedResponse, err := security.Decrypt(str, []byte(authValues.Key))
		if err != nil {
			logData.Message = "AddBeneficiary: Error decrypting response"
			s.LoggerService.LogError(logData)
			return nil, err
		}

		var beneficiariesResponseDetails []models.BeneficiaryDetail
		err = json.Unmarshal([]byte(DecryptedResponse), &beneficiariesResponseDetails)
		if err != nil {
			logData.Message = "AddBeneficiary: Error unmarshaling beneficiary response"
			s.LoggerService.LogError(logData)
			return nil, err
		}

		for _, beneficiary := range beneficiariesResponseDetails {
			if beneficiary.BenfAcctNo == r.BenfAcctNo && beneficiary.BenfMob == r.BenfMobNo {
				existingBeneficiary, err := models.BeneficiaryByBenfId(s.db, beneficiary.BenfId)
				if err != nil && !errors.Is(err, constants.ErrNoDataFound) {
					logData.Message = "AddBeneficiary: Error checking existing beneficiary"
					s.LoggerService.LogError(logData)
					return nil, err
				}

				if existingBeneficiary != nil {
					if existingBeneficiary.ActivatedDtTime == "" && !existingBeneficiary.IsActive {
						if err := models.UpdateBeneficiary(s.db, &models.BeneficiaryDTO{
							UserId:          authValues.UserId,
							BenfId:          beneficiary.BenfId,
							IsActive:        true,
							ActivatedDtTime: beneficiary.BenfActivateTime,
							BenfNickName:    types.FromString(beneficiary.BenfId),
						}); err != nil { // after success response data update in DB
							logData.Message = "ValidateOTPBeneficiary: Error updating beneficiary details"
							s.LoggerService.LogError(logData)
							return nil, err
						}
					}
					logData.Message = "AddBeneficiary: Beneficiary already exists for this account number"
					s.LoggerService.LogInfo(logData)
					return nil, errors.New("Beneficiary already exists for this account number")
				}

				logData.Message = "AddBeneficiary: Beneficiary not found in the database, inserting new beneficiary"
				s.LoggerService.LogInfo(logData)

				newBeneficiary := models.NewBeneficiaryDTO()
				if err := newBeneficiary.BankBindData(
					beneficiary.BenfId,
					authValues.UserId,
					beneficiary.BenfName,
					beneficiary.BenfId,
					beneficiary.BenfMob,
					beneficiary.BenfAcctNo,
					beneficiary.BenfIFSC,
					beneficiary.BenfAcctType,
					beneficiary.PaymentMode,
					beneficiary.BenfActivateTime,
					true,
				); err != nil {
					logData.Message = "AddBeneficiary: Error binding beneficiary data"
					s.LoggerService.LogError(logData)
					return nil, err
				}

				newBeneficiary.BenfAccountNo = beneficiary.BenfAcctNo
				if _, err := models.InsertBeneficiaryDetails(s.db, newBeneficiary); err != nil {
					logData.Message = "AddBeneficiary: Error inserting beneficiary details"
					s.LoggerService.LogError(logData)
					return nil, err
				}

				return "Beneficiary updated to database successfully", nil
			}
		}
	}

	existingAccount, err := models.GetAccountDataByUserId(s.db, existingDevice.UserId)

	if err != nil {
		logData.Message = "AddBeneficiary: Error getting account data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	request := requests.NewOutgoingAddBeneficiaryRequest()

	// Scenario 4- Beneficiary details sent to KVB, but Error Received (Invalid Input)
	// Scenario 5- Beneficiary details sent to KVB, but Error Received (Invalid Customer Data)
	// Scenario 6- Beneficiary details sent to KVB, but Error Received (Beneficiary ID already available for the applicant )
	if err := request.Bind(existingDevice.ApplicantId, existingAccount.AccountNumber, r); err != nil {
		logData.Message = "AddBeneficiary: Error binding add beneficiary request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	// Scenario 2 - If Beneficiary data sent to KVB and success res received and left the half way without completion, retried with new request with change in value
	// to check the requests body and DB avail data is same or not we are retrving data
	existingBeneficiary, err := models.GetBeneficiaryByTxnId(s.db, r.TxnIdentifier, authValues.UserId)
	if err != nil && !errors.Is(err, constants.ErrNoDataFound) {
		logData.Message = "AddBeneficiary 1: FindBeneficiaryByNameAndIfscCode Error fetching beneficiaries from the database"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	maskedAccount := strings.Repeat("X", len(r.BenfAcctNo)-4) + r.BenfAcctNo[len(r.BenfAcctNo)-4:]

	if existingBeneficiary != nil {
		request.TxnIdentifier = existingBeneficiary.TxnIdentifier
		// Senerio 1 - If Beneficiary data sent to KVB and success response received and left the half way without completion
		// Senerio 10 - Recived and OTP verification request  initiated and receive fail response OTP Expired
		if checkRequestDataIsSame(r, existingBeneficiary) {
			request.ResendOtp = "Y"
			request.RetryFlag = "N"
		} else {
			// Senerio 2 - If Beneficiary data sent to KVB and success res received and left the half way without completion, retried with new request with change in value
			request.ResendOtp = "N"
			request.RetryFlag = "N"
		}
	}

	var response *responses.BeneficiarySubmissionResponse
	var opErr error

	response, opErr = s.bankService.AddBeneficiary(ctx, request)
	if opErr != nil {
		bankErr := s.bankService.HandleBankSpecificError(opErr, func(errorCode string) (string, bool) {
			return constants.GetBeneficiaryErrorMessage(errorCode)
		})

		if bankErr != nil {
			logData.Message = fmt.Sprintf("AddBeneficiary: Bank error encountered (ErrorCode: %s)", bankErr.ErrorCode)
			s.LoggerService.LogError(logData)

			if msg, retryable := constants.GetBeneficiaryRetryErrorMessage(bankErr.ErrorCode); retryable {
				// Scenario 7- Beneficiary details sent to KVB, but Error Received (Technical issue)
				if bankErr.ErrorCode == constants.RetryErrorMW9999 {
					request.ResendOtp = "N"
					request.RetryFlag = "N"
				}

				// Scenario 3- Beneficiary details sent to KVB, but not received the response. (Back Ofice TimeOut)
				if bankErr.ErrorCode == constants.RetryErrorMW9997 || bankErr.ErrorCode == constants.RetryErrorMW9998 {
					request.ResendOtp = "N"
					request.RetryFlag = "Y"
				}

				// If error code is retryable, retry the request with change the resend otp and retry flag
				retryCount := 0
				retryErr := utils.RetryFunc(func() error {
					retryCount++
					logData.Message = fmt.Sprintf("AddBeneficiary: Retry attempt %d with txnid: %s", retryCount, request.TxnIdentifier)
					s.LoggerService.LogInfo(logData)

					response, opErr = s.bankService.AddBeneficiary(ctx, request)
					if opErr != nil {
						logData.Message = "AddBeneficiary: Error calling bank service add beneficiary"
						s.LoggerService.LogError(logData)
						return opErr
					}

					if response == nil {
						logData.Message = fmt.Sprintf("AddBeneficiary: Unknown error during retry attempt %d", retryCount)
						s.LoggerService.LogError(logData)
						return errors.New("unknown error during retry")
					}

					if response.ErrorCode != "0" && response.ErrorCode != "00" {
						logData.Message = fmt.Sprintf("AddBeneficiary: Error during retry attempt %d (ErrorCode: %s)", retryCount, response.ErrorCode)
						s.LoggerService.LogError(logData)
						return errors.New(response.ErrorMessage)
					}

					return nil
				}, 2)

				if retryErr != nil {
					logData.Message = fmt.Sprintf("AddBeneficiary: Adding beneficiary failed after %d retries", retryCount)
					s.LoggerService.LogError(logData)
					return nil, errors.New(msg)
				}

				logData.Message = fmt.Sprintf("AddBeneficiary: Request succeeded after %d retries", retryCount)
				s.LoggerService.LogInfo(logData)
				return response, nil
			} else {
				return nil, errors.New(bankErr.ErrorMessage)
			}
		} else {
			return nil, opErr
		}
	}

	if response == nil {
		logData.Message = "AddBeneficiary: Unknown error"
		s.LoggerService.LogError(logData)
		return nil, errors.New(constants.RetryErrorMessage)
	}

	// Senerio 0 - Request Sucessfully submited and Response Received For Beneficiary addition to KVB for OTP verification we are saving this data in redis
	// For Better approch we are passing this Key from the constants file
	benfKey := fmt.Sprintf(constants.BeneficiaryKey, authValues.UserId)
	benfData := map[string]string{
		"benf_id":       request.BenfID,
		"transactionId": request.TxnIdentifier,
	}

	if existingBeneficiary != nil {
		// Senerio 10 - Recived and OTP verification request  initiated and receive fail response OTP Expired
		if !checkRequestDataIsSame(r, existingBeneficiary) {
			// Updating beneficiary details
			beneficiaryDTO := &models.BeneficiaryDTO{
				UserId:        authValues.UserId,
				BenfId:        r.BenfNickName,
				BenfName:      r.BenfName,
				BenfNickName:  types.FromString(r.BenfNickName),
				BenfMobNo:     r.BenfMobNo,
				BenfAccountNo: maskedAccount,
				BenfIfsc:      r.BenfIFSC,
				BenfAcctType:  r.BenfAcctType,
				PaymentMode:   r.PaymentMode,
			}
			if err := models.UpdateBeneficiary(s.db, beneficiaryDTO); err != nil {
				logData.Message = "AddBeneficiary: Error updating beneficiary details"
				s.LoggerService.LogError(logData)
				return nil, err
			}
		}
	} else {
		// Insert new record after success response from bank for new beneficiary
		if err := s.InsertBeneficiaryDetails(r, authValues.UserId, request.TxnIdentifier); err != nil {
			// Handle duplicate key error (adjust based on your DB's actual error message)
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				logData.Message = "AddBeneficiary: Beneficiary already exists with this nikname"
				s.LoggerService.LogError(logData)
				return nil, errors.New("Beneficiary already exists with this nikname")
			}
			logData.Message = "AddBeneficiary: Error inserting beneficiary details"
			s.LoggerService.LogError(logData)
			return nil, err
		}
	}

	// Convert the map to a properly formatted JSON string
	benfDataJSON, err := json.MarshalIndent(benfData, "", "  ")
	if err != nil {
		logData.Message = fmt.Sprintf("JSON Marshalling Error: %v", err)
		return nil, err
	}

	// Store in Redis and handle errors
	// Senerio 0 - Request Sucessfully submited and Response Received For Beneficiary addition to KVB for OTP verification we are saving this data in redis
	// For Better approch we are passing this Key from the constants file
	if err := s.redis.Set(benfKey, benfDataJSON, constants.BeneficiaryTTL); err != nil {
		logData.Message = fmt.Sprintf("Error saving beneficiary data in Redis: %v", err)
		return nil, err
	}

	byteudd, err := json.Marshal(response)

	if err != nil {
		logData.Message = "AddBeneficiary: Error marshaling response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	encrypted, err := security.Encrypt(byteudd, []byte(authValues.Key))

	if err != nil {
		logData.Message = "AddBeneficiary: Error encrypting response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "AddBeneficiary: Response encrypted successfully"
	logData.ResponseSize = len(byteudd)
	logData.ResponseBody = string(byteudd)
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	// scenario 0 success
	return encrypted, nil
}

// checking request is same as the data available in the database
// Scenario 2 - If Beneficiary data sent to KVB and success res received and left the half way without completion, retried with new request with change in value
// to check the requests body and DB avail data is same or not we are retrving data
func checkRequestDataIsSame(r *requests.AddNewBeneficiary, existingBeneficiary *models.BeneficiaryDTO) bool {
	maskedAccount := strings.Repeat("X", len(r.BenfAcctNo)-4) + r.BenfAcctNo[len(r.BenfAcctNo)-4:]
	if existingBeneficiary.BenfName != r.BenfName && existingBeneficiary.BenfMobNo != r.BenfMobNo && existingBeneficiary.BenfNickName.String != r.BenfNickName && existingBeneficiary.BenfId != r.BenfNickName && existingBeneficiary.BenfAccountNo != maskedAccount && existingBeneficiary.BenfIfsc != r.BenfIFSC {
		return false
	} else {
		return true
	}
}

// Helper function to insert the beneficiary details into the database
func (s *Store) InsertBeneficiaryDetails(r *requests.AddNewBeneficiary, userId, txnId string) error {
	benf := models.NewBeneficiaryDTO()
	if err := benf.BindData(r.BenfNickName, userId, r.BenfName, r.BenfNickName, r.BenfMobNo, r.BenfAcctNo, r.BenfIFSC, r.BenfAcctType, r.PaymentMode, "", txnId); err != nil {
		return err
	}

	if _, err := models.InsertBeneficiaryDetails(s.db, benf); err != nil {
		return err
	}

	return nil
}

func (s *Store) ValidateOTPBeneficiary(ctx context.Context, authValues *models.AuthValues, r *requests.AddBeneficiaryOtpRequest) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.BENEFICIARY,
		RequestURI: "/api/beneficiary/beneficiary-otp",
		Message:    "ValidateOTPBeneficiary log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	existingDevice, err := models.GetUserDataByUserId(s.db, authValues.UserId)

	if err != nil {
		logData.Message = "ValidateOTPBeneficiary: Error getting user data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	existingAccount, err := models.GetAccountDataByUserId(s.db, existingDevice.UserId)

	if err != nil {
		logData.Message = "ValidateOTPBeneficiary: Error getting account data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	request := requests.NewOutgoingBeneficiaryOTPRequest()

	if err := request.Bind(existingDevice.ApplicantId, existingAccount.AccountNumber, r.Otp); err != nil {
		logData.Message = "ValidateOTPBeneficiary: Error binding OTP request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	// Senerio 0 - Request Sucessfully submited and Response Received For Beneficiary addition to KVB for OTP verification we are saving this data in redis
	// For Better approch we are passing this Key from the constants file
	benfData, err := s.memory.Get(fmt.Sprintf(constants.BeneficiaryKey, authValues.UserId))
	if err != nil {
		return nil, err
	}

	if len(benfData) == 0 {
		return nil, errors.New("beneficiary data not found in redis")
	}

	var benfDataMap map[string]string
	err = json.Unmarshal([]byte(benfData), &benfDataMap)
	if err != nil {
		logData.Message = "ValidateOTPBeneficiary: Error unmarshalling beneficiary data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	transactionId := benfDataMap["transactionId"]
	request.TxnIdentifier = transactionId
	var response *responses.BeneficiaryOTPValidationResponse
	var opErr error

	response, opErr = s.bankService.SubmitOtpBeneficiaryAddition(ctx, request)
	if opErr != nil {
		bankErr := s.bankService.HandleBankSpecificError(opErr, func(errorCode string) (string, bool) {
			return constants.GetBeneficiaryErrorMessage(errorCode)
		})

		if bankErr != nil {
			logData.Message = fmt.Sprintf("ValidateOTPBeneficiary: Bank error encountered (ErrorCode: %s)", bankErr.ErrorCode)
			s.LoggerService.LogError(logData)

			beneficiaryId := benfDataMap["benf_id"]
			fetchResponse, fetchErr := s.VerifyAndConfirmBeneficiary(ctx, authValues, beneficiaryId)
			if fetchErr == nil && fetchResponse != nil {
				logData.Message = "ValidateOTPBeneficiary: Beneficiary data found, returning success"
				s.LoggerService.LogInfo(logData)
				return fetchResponse, nil
			}

			// Senerio 12-OTP Recived and OTP verification request  initiated and receive fail  response Technical Error
			//Senerio 12-OTP Recived and OTP verification request  initiated and receive fail  response back Office Time Out
			if bankErr.ErrorCode == constants.RetryErrorMW9999 || bankErr.ErrorCode == constants.RetryErrorMW9997 || bankErr.ErrorCode == constants.RetryErrorMW9998 {
				logData.Message = "ValidateOTPBeneficiary: Retrying OTP submission"
				s.LoggerService.LogInfo(logData)

				retryCount := 0
				retryErr := utils.RetryFunc(func() error {
					retryCount++
					logData.Message = fmt.Sprintf("ValidateOTPBeneficiary: Retry attempt %d for OTP submission", retryCount)
					s.LoggerService.LogInfo(logData)

					retryResponse, err := s.bankService.SubmitOtpBeneficiaryAddition(ctx, request)
					if err != nil {
						logData.Message = fmt.Sprintf("ValidateOTPBeneficiary: Error during retry %d (Error: %v)", retryCount, err)
						s.LoggerService.LogError(logData)
						return err
					}

					if retryResponse.ErrorCode != "0" && retryResponse.ErrorCode != "00" {
						logData.Message = fmt.Sprintf("ValidateOTPBeneficiary: ErrorCode %s during retry %d", retryResponse.ErrorCode, retryCount)
						s.LoggerService.LogError(logData)
						return errors.New(retryResponse.ErrorMessage)
					}

					response = retryResponse
					return nil
				}, 2)

				if retryErr != nil {
					finalCheckResponse, finalCheckErr := s.VerifyAndConfirmBeneficiary(ctx, authValues, beneficiaryId)
					if finalCheckErr == nil && finalCheckResponse != nil {
						logData.Message = "ValidateOTPBeneficiary: Beneficiary found after all retries failed"
						s.LoggerService.LogInfo(logData)
						return finalCheckResponse, nil
					}

					logData.Message = fmt.Sprintf("ValidateOTPBeneficiary: OTP submission failed after %d retries, and beneficiary not found", retryCount)
					s.LoggerService.LogError(logData)
					return nil, errors.New(constants.RetryErrorMessage)
				}

				logData.Message = fmt.Sprintf("ValidateOTPBeneficiary: OTP submission succeeded after %d retries", retryCount)
				s.LoggerService.LogInfo(logData)
				return response, nil
			}

			return nil, errors.New(constants.InputErrorMessage)
		}
		return nil, opErr
	}

	// Senerio 0 - Request Sucessfully submited and Response Received For Beneficiary addition to KVB for OTP verification
	// after OTP verification we are marking as Active Here
	if err := models.UpdateBeneficiary(s.db, &models.BeneficiaryDTO{
		UserId:          authValues.UserId,
		BenfId:          benfDataMap["benf_id"],
		IsActive:        true,
		ActivatedDtTime: response.ActivationDtTime,
	}); err != nil {
		logData.Message = "ValidateOTPBeneficiary: Error updating beneficiary details"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	// Senerio 0 - Request Sucessfully submited and Response Received For Beneficiary addition to KVB for OTP verification we are saving this data in redis
	// For Better approch we are passing this Key from the constants file
	benfKey := fmt.Sprintf(constants.BeneficiaryKey, authValues.UserId)

	// Directly delete from Redis after user accesses the data
	s.redis.Delete(benfKey)

	byteudd, err := json.Marshal(response)

	if err != nil {
		logData.Message = "ValidateOTPBeneficiary: Error marshaling response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	encrypted, err := security.Encrypt(byteudd, []byte(authValues.Key))

	if err != nil {
		logData.Message = "ValidateOTPBeneficiary: Error encrypting response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if err := models.GenerateNotification(authValues.UserId, "beneficiary_addition", utils.CalculateTimeDifference(response.ActivationDtTime), "beneficiary_addition"); err != nil {
		logData.Message = "UpdateAdressLog: Error while Generating Notification " + err.Error()
		s.LoggerService.LogError(logData)
	}

	logData.Message = "ValidateOTPBeneficiary: Response encrypted successfully"
	logData.ResponseSize = len(byteudd)
	logData.ResponseBody = string(byteudd)
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)
	// scenarion 9 success
	return encrypted, nil
}

func (s *Store) BeneficiaryPayment(ctx context.Context, authValues *models.AuthValues, r *requests.PaymentRequest) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.BENEFICIARY,
		RequestURI: "/api/beneficiary/payment",
		Message:    "BeneficiaryPayment log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	existingDevice, err := models.GetUserDataByUserId(s.db, authValues.UserId)

	if err != nil {
		logData.Message = "BeneficiaryPayment: Error getting user data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if strings.ToLower(r.PaymentMode) == "ift" && r.BenfIfsc[:4] != "KVBL" {
		return nil, errors.New("invalid ifsc code for ift payment mode")
	}

	existingAccount, err := models.GetAccountDataByUserId(s.db, existingDevice.UserId)

	if err != nil {
		logData.Message = "BeneficiaryPayment: Error getting account data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if r.BenfAcctNo == existingAccount.AccountNumber {
		return nil, errors.New("self transfer not allowed")
	}

	ifscData, err := models.GetIFSCData(s.db, r.BenfIfsc)
	if err != nil {
		if errors.Is(err, constants.ErrNoDataFound) {
			logData.Message = "BeneficiaryPayment: Invalid IFSC code provided"
			s.LoggerService.LogError(logData)
			return nil, errors.New(constants.InputErrorMessage)
		}

		logData.Message = "BeneficiaryPayment: Error fetching IFSC data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = fmt.Sprintf("BeneficiaryPayment: IFSC data found: %v", ifscData.IfscCode)
	s.LoggerService.LogInfo(logData)

	request := requests.NewOutgoingPaymentRequest()

	if err := request.Bind(existingDevice.ApplicantId, existingAccount.AccountNumber, r); err != nil {
		logData.Message = "BeneficiaryPayment: Error binding payment request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if request.ResendOtp == "Y" {
		txn, err := s.memory.Get(fmt.Sprintf("beneficiary:payment:transaction:%s", authValues.UserId))
		if err != nil {
			return nil, err
		}
		if txn == "null" {
			return nil, err
		}
		request.TxnIdentifier = txn
	}

	var response *responses.PaymentSubmissionResponse
	var opErr error

	response, opErr = s.bankService.PaymentSubmission(ctx, request)
	if opErr != nil {
		bankErr := s.bankService.HandleBankSpecificError(opErr, func(errorCode string) (string, bool) {
			return constants.GetPaymentCallbackErrorMessage(errorCode)
		})

		if bankErr != nil {
			logData.Message = fmt.Sprintf("Bank error encountered (ErrorCode: %s)", bankErr.ErrorCode)
			s.LoggerService.LogError(logData)

			if bankErr.ErrorCode == constants.PaymentCallbackErrorCodeMW0014 {
				response = &responses.PaymentSubmissionResponse{
					ApplicantId:   request.ApplicantId,
					AccountNo:     request.AccountNo,
					TxnStatus:     response.TxnStatus,
					ErrorCode:     bankErr.ErrorCode,
					ErrorMessage:  bankErr.ErrorMessage,
					TxnIdentifier: request.TxnIdentifier,
					TxnRefNo:      response.TxnRefNo,
				}
			} else if msg, retryable := constants.GetPaymentCallbackRetryErrorMessage(bankErr.ErrorCode); retryable {
				err := utils.RetryFunc(func() error {
					response, opErr = s.bankService.PaymentSubmission(ctx, request)
					return opErr
				}, 2)

				if err != nil {
					logData.Message = "BeneficiaryPayment: Callback failed after retries"
					s.LoggerService.LogError(logData)
					return nil, errors.New(msg)
				}
			} else {
				return nil, errors.New(bankErr.ErrorMessage)
			}
		}

		return nil, opErr
	}

	if response == nil {
		logData.Message = "BeneficiaryPayment: Unknown error"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	_, err = models.FindOneTransactionByUserAndTransactionId(s.db, authValues.UserId, request.TxnIdentifier)

	if err != nil {
		if !errors.Is(err, constants.ErrNoDataFound) {
			logData.Message = "BeneficiaryPayment: Error finding transaction details"
			s.LoggerService.LogError(logData)
			return nil, err
		}

		// add quick transfer beneficiary in benficiaries table db with inactive flag
		var benfID string
		var benfUUID string

		if strings.ToLower(r.QuickTransfer) == "y" {
			existingBeneficiary, err := models.FindBeneficiaryByNameAndIfscCodeV2(s.db,
				authValues.UserId,
				r.BenfName,
				request.BenfIfsc,
				r.BenfNickName,
				r.BenfMobNo,
			)
			if err != nil && !errors.Is(err, constants.ErrNoDataFound) {
				logData.Message = "BeneficiaryPayment: FindBeneficiaryByNameAndIfscCode Error fetching beneficiaries from the database"
				s.LoggerService.LogError(logData)
				return nil, err
			}

			if existingBeneficiary == nil {
				if r.BenfNickName != "" {
					benfID = r.BenfNickName
				} else {
					benfID = formatBenfID(request.BenfName, request.BenfAcctNo)
				}
				benfUUID, err = models.InsertBeneficiaryDetails(s.db, &models.BeneficiaryDTO{
					BenfId:        benfID,
					BenfNickName:  types.FromString(benfID),
					UserId:        authValues.UserId,
					BenfName:      request.BenfName,
					BenfMobNo:     request.BenfMobNo,
					BenfIfsc:      request.BenfIfsc,
					BenfAcctType:  request.BenfAcctType,
					PaymentMode:   request.PaymentMode,
					IsActive:      false,
					BenfAccountNo: request.BenfAcctNo,
				})
				if err != nil {
					logData.Message = "BeneficiaryPayment: Error inserting beneficiary details"
					s.LoggerService.LogError(logData)
					return nil, err
				}
			}
		} else {
			ben, err := models.GetBeneficiaryByID(s.db, r.BenfId, request.BenfAcctNo, request.BenfMobNo, request.BenfIfsc)
			if err != nil {
				logData.Message = "GetBeneficiaryByID: Error fetching beneficiary details"
				s.LoggerService.LogError(logData)
			}
			if ben == nil {
				benfUUID = r.BenfId
			} else {
				benfUUID = ben.Id.String()
			}
		}

		// save transaction details to db
		if err := models.InsertTransaction(s.db, &models.Transaction{
			UserID:          authValues.UserId,
			TransactionID:   request.TxnIdentifier,
			Amount:          types.FromString(request.Amount),
			PaymentMode:     models.PaymentMode(request.PaymentMode),
			TransactionDesc: types.FromString(request.TxnRemarks),
			BeneficiaryID:   types.FromString(benfUUID),
		}); err != nil {
			logData.Message = "BeneficiaryPayment: Error inserting transaction details"
			s.LoggerService.LogError(logData)
			return nil, err
		}
	}
	// save transaction id to cache
	err = s.memory.Set(fmt.Sprintf("beneficiary:payment:transaction:%s", authValues.UserId), request.TxnIdentifier, time.Hour)
	if err != nil {
		return nil, err
	}

	byteudd, err := json.Marshal(response)

	if err != nil {
		logData.Message = "BeneficiaryPayment: Error marshaling response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	encrypted, err := security.Encrypt(byteudd, []byte(authValues.Key))

	if err != nil {
		logData.Message = "BeneficiaryPayment: Error encrypting response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "BeneficiaryPayment: Response encrypted successfully"
	logData.ResponseSize = len(byteudd)
	logData.ResponseBody = string(byteudd)
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return encrypted, nil
}

func (s *Store) BeneficiaryPaymentOTP(ctx context.Context, authValues *models.AuthValues, otp string) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.BENEFICIARY,
		RequestURI: "/api/beneficiary/payment-otp",
		Message:    "Initiating beneficiary payment OTP",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	existingDevice, err := models.GetUserDataByUserId(s.db, authValues.UserId)

	if err != nil {
		logData.Message = "BeneficiaryPaymentOTP: Error getting user data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	existingAccount, err := models.GetAccountDataByUserId(s.db, existingDevice.UserId)

	if err != nil {
		logData.Message = "BeneficiaryPaymentOTP: Error getting account data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	request := requests.NewOutgoingPaymentRequestOTP()

	if err := request.Bind(existingDevice.ApplicantId, existingAccount.AccountNumber, otp); err != nil {
		logData.Message = "BeneficiaryPaymentOTP: Error binding payment OTP request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	// get transaction id from cache
	txn, err := s.memory.Get(fmt.Sprintf("beneficiary:payment:transaction:%s", authValues.UserId))
	if err != nil {
		return nil, err
	}
	request.TxnIdentifier = txn

	var response *responses.PaymentSubmissionOtpResponse
	var opErr error

	response, opErr = s.bankService.PaymentSubmissionOTP(ctx, request)
	if opErr != nil {
		bankErr := s.bankService.HandleBankSpecificError(opErr, func(errorCode string) (string, bool) {
			return constants.GetPaymentCallbackErrorMessage(errorCode)
		})

		if bankErr != nil {
			logData.Message = fmt.Sprintf("Bank error encountered (ErrorCode: %s)", bankErr.ErrorCode)
			s.LoggerService.LogError(logData)

			if msg, retryable := constants.GetPaymentCallbackRetryErrorMessage(bankErr.ErrorCode); retryable {
				err := utils.RetryFunc(func() error {
					response, opErr = s.bankService.PaymentSubmissionOTP(ctx, request)
					return opErr
				}, 2)

				if err != nil {
					logData.Message = "BeneficiaryPaymentOTP: Callback failed after retries"
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

	if response == nil {
		logData.Message = "BeneficiaryPaymentOTP: Unknown error"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	// update transaction status in db
	if err := models.UpdateTransactionByTransID(s.db, &models.Transaction{
		TransactionID: request.TxnIdentifier,
		OTPStatus:     types.FromString(response.TxnStatus),
	}); err != nil {
		logData.Message = "BeneficiaryPaymentOTP: Error updating transaction status in db"
		s.LoggerService.LogError(logData)
	}

	if response.ErrorCode != "0" && response.ErrorCode != "00" {
		logData.Message = "BeneficiaryPaymentOTP: Error processing payment OTP"
		s.LoggerService.LogError(logData)
		return nil, errors.New(response.ErrorMessage)
	}

	response.TxnDateTime = time.Now()

	byteudd, err := json.Marshal(response)

	if err != nil {
		logData.Message = "BeneficiaryPaymentOTP: Error marshaling response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	encrypted, err := security.Encrypt(byteudd, []byte(authValues.Key))

	if err != nil {
		logData.Message = "BeneficiaryPaymentOTP: Error encrypting response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "BeneficiaryPaymentOTP: Response encrypted successfully"
	logData.ResponseSize = len(byteudd)
	logData.ResponseBody = string(byteudd)
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return encrypted, nil
}

func (s *Store) CreateQuickTransferTemplate(ctx context.Context, authValues *models.AuthValues, r *requests.QuickTransferBeneficiaryRegistrationRequest) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.BENEFICIARY,
		RequestURI: "/api/beneficiary/quick-transfer-template",
		Message:    "CreateQuickTransferTemplate LOG",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	existingDevice, err := models.GetUserDataByUserId(s.db, authValues.UserId)

	if err != nil {
		logData.Message = "CreateQuickTransferTemplate: Error getting user data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	existingAccount, err := models.GetAccountDataByUserId(s.db, existingDevice.UserId)

	if err != nil {
		logData.Message = "CreateQuickTransferTemplate: Error getting account data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	outgoingTemplateRequest := requests.NewOutBeneficiaryTemplateRequest()

	if err := outgoingTemplateRequest.BindAndValidate(existingDevice.ApplicantId, existingAccount.AccountNumber, r.BenficiaryId); err != nil {
		logData.Message = "CreateQuickTransferTemplate: Error binding and validating template request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	// save transaction id to cache
	txn, err := s.memory.Get(fmt.Sprintf("beneficiary:payment:transaction:%s", authValues.UserId))
	if err != nil {
		return nil, err
	}
	outgoingTemplateRequest.TxnIdentifier = txn

	response, err := s.bankService.QuickTransferTemplateAdd(ctx, outgoingTemplateRequest)

	if err != nil {
		bankErr := s.bankService.HandleBankSpecificError(err, func(errorCode string) (string, bool) {
			return constants.GetQuickTransferBeneficiaryErrorMessage(errorCode)
		})

		if bankErr != nil {
			logData.Message = fmt.Sprintf("CreateQuickTransferTemplate: Bank error encountered (ErrorCode: %s)", bankErr.ErrorCode)
			s.LoggerService.LogError(logData)
			return nil, errors.New(bankErr.ErrorMessage)
		}

		return nil, err
	}

	if response == nil {
		logData.Message = "CreateQuickTransferTemplate: Unknown error"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if response.ErrorCode != "0" && response.ErrorCode != "00" {
		logData.Message = "CreateQuickTransferTemplate: Error processing quick transfer template"
		s.LoggerService.LogError(logData)
		return nil, errors.New(response.ErrorMessage)
	}

	byteudd, err := json.Marshal(response)

	if err != nil {
		logData.Message = "CreateQuickTransferTemplate: Error marshaling response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	encrypted, err := security.Encrypt(byteudd, []byte(authValues.Key))

	if err != nil {
		logData.Message = "CreateQuickTransferTemplate: Error encrypting response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "CreateQuickTransferTemplate: Response encrypted successfully"
	logData.ResponseSize = len(byteudd)
	logData.ResponseBody = string(byteudd)
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return encrypted, nil
}

func (s *Store) CheckBeneficiaryPaymentStatus(ctx context.Context, authValues *models.AuthValues, txnid string) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.BENEFICIARY,
		RequestURI: "/api/beneficiary/payment-status",
		Message:    "CheckPaymentStatus Log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
		StartTime:  time.Now(),
	}

	paymentData, err := models.FindOneTransactionByUserAndTransactionId(s.db, authValues.UserId, txnid)
	if err != nil {
		logData.Message = "CheckPaymentStatus Log: Getting error while finding transaction " + err.Error()
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if paymentData == nil {
		logData.Message = "CheckPaymentStatus Log: payment detail not found "
		logData.EndTime = time.Now()
		s.LoggerService.LogError(logData)
		return nil, errors.New("payment detail not found")
	}

	response := responses.NewPaymentStatusResponse()
	response.TxnStatus = paymentData.CBSStatus.String
	response.TxnIdentifier = paymentData.TransactionID
	response.TxnRefNo = paymentData.UTRRefNumber.String
	response.TransactionAmount = paymentData.Amount.String
	response.PaydohTransactionId = paymentData.ID.String()

	jsonData, err := json.Marshal(response)
	if err != nil {
		logData.Message = "CheckPaymentStatus Log: error while marshalling transaction data" + err.Error()
		s.LoggerService.LogError(logData)
		return nil, err
	}
	encrypted, err := security.Encrypt(jsonData, []byte(authValues.Key))
	if err != nil {
		logData.Message = "CheckPaymentStatus Log: error while encrypting transaction data" + err.Error()
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.Message = "CheckPaymentStatus Log: Response encrypted successfully"
	logData.ResponseSize = len(encrypted)
	logData.ResponseBody = string(encrypted)
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)
	return encrypted, nil
}

func formatBenfID(benfName, benfAcctNo string) string {
	benfName = strings.TrimSpace(benfName)
	if len(benfName) >= 4 {
		return benfName[:4] + benfAcctNo[len(benfAcctNo)-4:]
	}
	return benfName + benfAcctNo[len(benfAcctNo)-4:]
}

// Helper FetchBankBeneficiary calls the bank service to fetch the list of beneficiaries for the user
func (s *Store) FetchBankBeneficiary(ctx context.Context, authValues *models.AuthValues) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.BENEFICIARY,
		RequestURI: "Internal fetch beneficiary",
		Message:    "FetchBankBeneficiary log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	existingDevice, err := models.GetUserDataByUserId(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "FetchBankBeneficiary: Error getting user data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	existingAccount, err := models.GetAccountDataByUserId(s.db, existingDevice.UserId)
	if err != nil {
		logData.Message = "FetchBankBeneficiary: Error getting account data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	request := requests.NewOutgoingBeneficiarySearchRequest()
	if err := request.Bind(existingDevice.ApplicantId, existingAccount.AccountNumber); err != nil {
		logData.Message = "FetchBankBeneficiary: Error binding beneficiary search request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	response, err := s.bankService.GetBeneficiaries(ctx, request)
	if err != nil {
		logData.Message = "FetchBankBeneficiary: Error calling bank service"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if response == nil {
		logData.Message = "FetchBankBeneficiary: Error getting beneficiary error"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	b := models.NewBeneficiaryDetail()
	beneficiaries, err := b.FetchBind(response)
	if err != nil {
		logData.Message = "FetchBankBeneficiary: Error binding beneficiary response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	byteudd, err := json.Marshal(beneficiaries)
	if err != nil {
		logData.Message = "FetchBankBeneficiary: Error marshaling beneficiary response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	encrypted, err := security.Encrypt(byteudd, []byte(authValues.Key))
	if err != nil {
		logData.Message = "FetchBankBeneficiary: Error encrypting beneficiary response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "FetchBankBeneficiary: Beneficiary response encrypted successfully"
	logData.ResponseSize = len(byteudd)
	logData.ResponseBody = string(byteudd)
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return encrypted, nil
}

// Helper function VerifyAndConfirmBeneficiary checks if the beneficiary is present in the list of beneficiaries from bank response
// if data is present then it will return success or else it will retry the request
func (s *Store) VerifyAndConfirmBeneficiary(ctx context.Context, authValues *models.AuthValues, beneficiaryId string) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.BENEFICIARY,
		RequestURI: "Internal verify beneficiary",
		Message:    "VerifyAndConfirmBeneficiary log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	beneficiariesResponse, err := s.FetchBankBeneficiary(ctx, authValues)
	if err != nil {
		logData.Message = "VerifyAndConfirmBeneficiary: Error fetching beneficiaries"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	str, ok := beneficiariesResponse.(string)
	if !ok {
		logData.Message = "VerifyAndConfirmBeneficiary: Error converting beneficiariesResponse to string"
		s.LoggerService.LogError(logData)
		return nil, errors.New("invalid response type")
	}

	decryptedResponse, err := security.Decrypt(str, []byte(authValues.Key))
	if err != nil {
		logData.Message = "VerifyAndConfirmBeneficiary: Error decrypting response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	var beneficiaries []models.BeneficiaryDetail
	err = json.Unmarshal([]byte(decryptedResponse), &beneficiaries)
	if err != nil {
		logData.Message = "VerifyAndConfirmBeneficiary: Error unmarshaling beneficiary response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	for _, benef := range beneficiaries {
		if benef.BenfId == beneficiaryId {
			logData.Message = fmt.Sprintf("VerifyAndConfirmBeneficiary: Beneficiary with ID %s found", beneficiaryId)
			return nil, nil
		}
	}

	logData.Message = fmt.Sprintf("VerifyAndConfirmBeneficiary: Beneficiary with ID %s not found", beneficiaryId)
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return nil, errors.New("beneficiary not found")
}
