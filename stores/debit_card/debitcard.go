package debitcard

import (
	"bankapi/constants"
	"bankapi/ftp_server"
	"bankapi/models"
	"bankapi/requests"
	"bankapi/responses"
	"bankapi/security"
	"bankapi/services"
	"bankapi/utils"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/paydoh/paydoh-commons/database"
	"bitbucket.org/paydoh/paydoh-commons/pkg/storage"
	"bitbucket.org/paydoh/paydoh-commons/pkg/task"
	commonSrv "bitbucket.org/paydoh/paydoh-commons/services"
	"bitbucket.org/paydoh/paydoh-commons/types"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type DebitCard interface {
	GetDebitCardData(authValues *models.AuthValues) (interface{}, error)
	NameOnDebitCard(fullName ...string) string
}

type Store struct {
	db                      *sql.DB
	m                       *database.Document
	memory                  *database.InMemory
	debitcardservice        *services.DebitcardApiService
	debitCardControlService *services.DebitcardControlApiService
	bankservice             *services.BankApiService
	LoggerService           *commonSrv.LoggerService
	s3Client                storage.S3Client
	cachedTrackingFile      *excelize.File
	lastFileDownloadTime    time.Time
	auditLogSrv             services.AuditLogService
	taskEnqueuer            task.TaskEnqueuer
}

type CardStatus struct {
	Date            string `xlsx:"Date" json:"date"`
	CardNumber      string `xlsx:"Card Number" json:"card_number"`
	CardHolderName  string `xlsx:"Card Holder Name" json:"card_holder_name"`
	ReferenceNumber string `xlsx:"Reference Number" json:"reference_number"`
	AWB             string `xlsx:"AWB" json:"awb"`
	CardType        string `xlsx:"Card Type" json:"card_type"`
	DispatchStatus  string `xlsx:"Dispatch Status" json:"dispatch_status"`
	DispatchMode    string `xlsx:"Dispatch Mode" json:"dispatch_mode"`
	DispatchDate    string `xlsx:"Dispatch Date" json:"dispatch_date"`
}

type FTPConfig struct {
	Host     string
	User     string
	Password string
	Port     int
	FilePath string
}

func NewStore(log *commonSrv.LoggerService, db *sql.DB, m *database.Document, memory *database.InMemory, s3Client storage.S3Client, auditLogSrv services.AuditLogService, taskEnqueuer task.TaskEnqueuer) *Store {
	debitcardService := services.NewDebitcardApiService(log, memory)
	debitcardcontrolService := services.NewDebitcardControlApiService(log, memory)
	bankService := services.NewBankApiService(log, memory)

	return &Store{
		db:                      db,
		m:                       m,
		memory:                  memory,
		debitcardservice:        debitcardService,
		LoggerService:           log,
		bankservice:             bankService,
		s3Client:                s3Client,
		debitCardControlService: debitcardcontrolService,
		auditLogSrv:             auditLogSrv,
		taskEnqueuer:            taskEnqueuer,
	}
}

func (s *Store) DebitCardGeneration(ctx context.Context, authValue *models.AuthValues, req *requests.DebitCardGenerationRequest) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.DEBITCARD,
		RequestURI: "/api/debitcard/generate",
		Message:    "DebitCardGeneration log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
		AppVersion: utils.GetAppVersionFromContext(ctx),
	}

	var transactionID string
	debitCardData, err := models.GetDebitCardData(s.db, authValue.UserId)
	if err != nil {
		if err != constants.ErrNoDataFound {
			s.ErrorLogData(logData, "DebitCardGeneration: Error while checking debit card data is available or not")
			return nil, err
		}
		transactionID = ""
		s.ErrorLogData(logData, "DebitCardGeneration: No DebitCard data is Available for this user")
	}

	if debitCardData != nil {
		if debitCardData.IsVirtualCardGenerated && debitCardData.IsPhysicalCardGenerated {
			s.ErrorLogData(logData, "DebitCardGeneration: Debit card already exists")
			return nil, errors.New("debit card already exists")
		}
	}

	onboardingStages, err := models.GetUserOnboardingStatus(authValue.UserId)
	if err != nil {
		s.ErrorLogData(logData, "DebitCardGeneration: Error while checking onboarding status")
		return nil, err
	}

	if !onboardingStages.IsDebitCardConsentComplete || !onboardingStages.IsDebitCardPaymentComplete {
		s.ErrorLogData(logData, "DebitCardGeneration: Debit card consent and payment not completed")
		return nil, errors.New("debit card consent and payment not completed")
	}

	accountDetail, err := models.GetUserAndAccountDetailByUserID(s.db, authValue.UserId)
	if err != nil {
		s.ErrorLogData(logData, "DebitCardGeneration: Error getting account details")
		return nil, err
	}

	kycConsent, err := models.FindKycConsentByUserId(s.db, authValue.UserId)
	if err != nil {
		s.ErrorLogData(logData, "DebitCardGeneration: Error getting kyc consent data")
		return nil, err
	}

	if kycConsent == nil || !kycConsent.VirtualDebitCardConsent {
		s.ErrorLogData(logData, "DebitCardGeneration: Debit card consent not provided")
		return nil, errors.New("debit card consent not provided")
	}

	if kycConsent.VirtualDebitCardConsent {
		if debitCardData != nil {
			if debitCardData.IsVirtualCardGenerated && req.DebitCardGenerationType == "virtual" {
				s.ErrorLogData(logData, "DebitCardGeneration: Virtual debit card already exists")
				return nil, errors.New("virtual debit card already exists")
			}
			transactionID = debitCardData.TxnIdentifyer.String
		}

		if debitCardData == nil || !debitCardData.IsVirtualCardGenerated {
			_, err := s.GenerateVirtualDebitCard(ctx, debitCardData, accountDetail, authValue.UserId, transactionID, logData)
			if err != nil {
				s.ErrorLogData(logData, "DebitCardGeneration: Error while generating virtual debit card")
				return nil, err
			}

		}
	}

	if req.DebitCardGenerationType == "physical" || req.DebitCardGenerationType == "both" {
		if kycConsent.PhysicalDebitCardConsent && kycConsent.VirtualDebitCardConsent {
			_, err = s.GeneratePhysicalDebitCard(ctx, accountDetail, authValue.UserId, req.DebitCardGenerationType, logData)
			if err != nil {
				s.ErrorLogData(logData, "DebitCardGeneration: Error while generating physical debit card")
				return nil, err
			}
		} else {
			s.ErrorLogData(logData, "DebitCardGeneration: Physical debit card consent not given")
			return nil, errors.New("physical debit card consent not provided")
		}
	}

	logData.Message = "DebitCardGeneration: Response received successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return nil, nil
}

func (s *Store) GeneratePhysicalDebitCard(ctx context.Context, accountDetail *models.UserPersonalInformationAndAccountData, userID, requestType string, logData *commonSrv.LogEntry) (interface{}, error) {

	var transactionIDForPhysical string
	debitCardDataForPhysical, err := models.GetDebitCardData(s.db, userID)
	if err != nil {
		if err != constants.ErrNoDataFound {
			s.ErrorLogData(logData, "DebitCardGeneration: Error while checking debit card data is available or not")
			return nil, err
		}
		s.ErrorLogData(logData, "DebitCardGeneration: No DebitCard data is Available for this user")
		return nil, errors.New("virtual debit card is also required for a physical debit card")
	}

	if debitCardDataForPhysical != nil {
		if debitCardDataForPhysical.IsPhysicalCardGenerated {
			s.ErrorLogData(logData, "DebitCardGeneration: Physical debit card already exists.")
			return nil, errors.New("physical debit card already exists")
		}
		transactionIDForPhysical = debitCardDataForPhysical.PhysicalDebitCardTxnId.String
	}

	request := requests.NewGeneratePhysicalDebitCardOutGoingReq()

	if err := request.Bind(accountDetail.Applicant_id, debitCardDataForPhysical.Proxy_Number.String, accountDetail.AccountNumber, transactionIDForPhysical); err != nil {
		logData.Message = "DebitCardGeneration: Error binding physical debit card request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if debitCardDataForPhysical.PhysicalDebitCardTxnId.String == "" {
		if err := models.UpdateDebitCardData(s.db, &models.DebitCardData{
			UserID:                 userID,
			PhysicalDebitCardTxnId: types.FromString(request.TxnIdentifier),
		}); err != nil {
			logData.Message = "DebitCardGeneration: Error Updating phycical debit card txn Id"
			s.LoggerService.LogError(logData)
			return nil, err
		}
	}

	_, err = s.debitcardservice.DebitCardPhysicalGeneration(ctx, request)
	if err != nil {
		bankErr := s.bankservice.HandleBankSpecificError(err, func(errorCode string) (string, bool) {
			return constants.GetPhysicalDebitCardErrorMessage(errorCode)
		})

		if bankErr != nil {
			logData.Message = fmt.Sprintf("DebitCardGeneration: Bank error encountered (ErrorCode: %s)", bankErr.ErrorCode)
			s.LoggerService.LogError(logData)
			return nil, errors.New(bankErr.ErrorMessage)
		}

		return nil, err
	}

	if err := models.UpdateDebitCardData(s.db, &models.DebitCardData{
		UserID:                  userID,
		IsPhysicalCardGenerated: true,
	}); err != nil {
		logData.Message = "DebitCardGeneration: Error Updating phycical debit card generation status"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	return nil, nil
}

func (s *Store) GenerateVirtualDebitCard(ctx context.Context, debitCardData *models.DebitCardData, accountDetail *models.UserPersonalInformationAndAccountData, userID, transactionID string, logData *commonSrv.LogEntry) (interface{}, error) {
	virtualCardRequest := requests.NewGenerateVirtualDebitCardOutGoingReq()
	debitCardName := NameOnDebitCard(accountDetail.FirstName, accountDetail.MiddleName, accountDetail.LastName)

	if err := virtualCardRequest.Bind(debitCardName, accountDetail.Applicant_id, accountDetail.AccountNumber, transactionID); err != nil {
		s.ErrorLogData(logData, "DebitCardGeneration: Error binding virtual debit card request")
		return nil, err
	}

	if debitCardData == nil {
		if err := models.InsertDebitCardData(s.db, userID, virtualCardRequest.TxnIdentifier); err != nil {
			s.ErrorLogData(logData, "DebitCardGeneration: Error inserting debit card data")
			return nil, err
		}
	}

	var virtualCardRes *responses.GenerateDebitcardResponse
	var opErr error

	virtualCardRes, opErr = s.debitcardservice.DebitCardVirtualGeneration(ctx, virtualCardRequest)
	if opErr != nil {
		bankErr := s.bankservice.HandleBankSpecificError(opErr, func(errorCode string) (string, bool) {
			s.ErrorLogData(logData, fmt.Sprintf("Bank side error encountered (ErrorCode: %s)", errorCode))
			return constants.GetVirtualDebitCardErrorMessage(errorCode)
		})

		if bankErr != nil {
			if msg, retryable := constants.GetVirtualDebitCardErrorMessageRetry(bankErr.ErrorCode); retryable {
				err := utils.RetryFunc(func() error {
					s.ErrorLogData(logData, "DebitCardGeneration: Retrying virtual debit card generation")
					virtualCardRes, opErr = s.debitcardservice.DebitCardVirtualGeneration(ctx, virtualCardRequest)
					return opErr
				}, 2)

				if err != nil {
					s.ErrorLogData(logData, "DebitCardGeneration: Callback failed after retries")
					return nil, errors.New(msg)
				}
			} else {
				s.ErrorLogData(logData, fmt.Sprintf("DebitCardGeneration: Bank error encountered (ErrorCode: %s)", bankErr.ErrorCode))
				return nil, errors.New(bankErr.ErrorMessage)
			}
		}
		return nil, opErr
	}

	if err := models.UpdateDebitCardData(s.db, &models.DebitCardData{
		UserID:                 userID,
		Proxy_Number:           types.FromString(virtualCardRes.ProxyNumber),
		IsVirtualCardGenerated: true,
	}); err != nil {
		s.ErrorLogData(logData, "DebitCardGeneration: Error updating debit card data")
		return nil, err
	}

	// update onboarding status
	if err := models.UpdateUserOnboardingStatus(constants.DEBIT_CARD_GENERATION_STAGE, userID); err != nil {
		s.LoggerService.Logger.Error(err)
	}

	if err := s.ReferralRewardTransfer(userID); err != nil {
		s.LoggerService.Logger.Error(err)
	}

	return virtualCardRes, nil
}

func (s *Store) ReferralRewardTransfer(userID string) error {
	userData, err := s.m.FindOne(constants.MONGO_USER_DB,
		constants.MONGO_USER_COLLECTION, bson.M{
			"user_id": userID,
		}, bson.M{})
	if err != nil {
		return err
	}

	var user2 models.User
	if err := userData.Decode(&user2); err != nil {
		return err
	}

	if user2.ReferredBy != "" {
		_, err := s.m.FindOne(constants.RewardDatabaseName,
			constants.PAYDOH_REWARDS, bson.M{
				"userId": userID, "event": "ReferralFrom",
			}, bson.M{})
		if err != nil {
			if err.Error() != "mongo: no documents in result" {
				return err
			}
			if _, _, err := s.taskEnqueuer.EnqueueNow("referral:rewards", user2, "default"); err != nil {
				s.LoggerService.LogError(&commonSrv.LogEntry{
					Message: "processReward: Error enqueuing reward coin task for order: " + err.Error(),
				})
				return nil
			}
		}
	}

	return nil
}

func (s *Store) ErrorLogData(logData *commonSrv.LogEntry, message string) {
	logData.Message = message
	s.LoggerService.LogError(logData)
	return
}

func (s *Store) DebitCardDetail(ctx context.Context, authValue *models.AuthValues) (interface{}, error) {
	// (*responses.DebitCardDetailRes, error) {

	logData := &commonSrv.LogEntry{
		Action:     constants.DEBITCARD,
		RequestURI: "/api/debitcard/detail",
		Message:    "DebitCardDetail log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
		AppVersion: utils.GetAppVersionFromContext(ctx),
	}

	userData, err := models.GetUserDataByUserId(s.db, authValue.UserId)
	if err != nil {
		logData.Message = "DebitCardDetail: Error getting user data by user id"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	account, err := models.GetAccountDataByUserId(s.db, authValue.UserId)

	if err != nil {
		logData.Message = "DebitCardDetail: Error getting account data by user id"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	cardData, err := models.GetDebitCardData(s.db, authValue.UserId)

	if err != nil {
		logData.Message = "DebitCardDetail: Error getting debit card data by user id"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if cardData.Proxy_Number.String == "" || !cardData.Proxy_Number.Valid {
		logData.Message = "DebitCardDetail: debitcard not generated yet"
		s.LoggerService.LogError(logData)
		return nil, errors.New("DebitCard not generated yet")
	}

	req := requests.NewGetDebitcardDetailRequest()

	if err := req.GetDebitCardDetailReq_Bind(userData.ApplicantId, account.AccountNumber, cardData.Proxy_Number.String, cardData.TxnIdentifyer.String); err != nil {
		logData.Message = "DebitCardDetail: Error binding debit card detail request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	var result *responses.DebitcardDetailResponse
	var opErr error

	result, opErr = s.debitcardservice.GetDebitCardDetails(ctx, req)
	if opErr != nil {
		bankErr := s.bankservice.HandleBankSpecificError(opErr, func(errorCode string) (string, bool) {
			return constants.GetDebitCardFetchErrorMessage(errorCode)
		})

		if bankErr != nil {
			logData.Message = fmt.Sprintf("Bank error encountered (ErrorCode: %s)", bankErr.ErrorCode)
			s.LoggerService.LogError(logData)

			if msg, retryable := constants.GetDebitCardFetchRetryErrorMessage(bankErr.ErrorCode); retryable {
				err := utils.RetryFunc(func() error {
					result, opErr = s.debitcardservice.GetDebitCardDetails(ctx, req)
					return opErr
				}, 2)

				if err != nil {
					logData.Message = "DebitCardGeneration: Callback failed after retries"
					s.LoggerService.LogError(logData)
					return nil, errors.New(msg)
				}
			} else {
				return nil, errors.New(bankErr.ErrorMessage)
			}
		}

		return nil, opErr
	}

	result.ServiceData.CardData[0].EncryptedPAN, err = utils.NewAESEncryptionUtil().CardDecryption(result.ServiceData.CardData[0].EncryptedPAN, constants.CardKey)
	if err != nil {
		logData.Message = "GetDebitCardDetails: Error Encrypted Pan Unable to decrypt" + err.Error()
		s.LoggerService.LogError(logData)
		return nil, err
	}
	result.ServiceData.CardData[0].CvvValue, err = utils.NewAESEncryptionUtil().CardDecryption(result.ServiceData.CardData[0].CvvValue, constants.CardCvvKey)
	if err != nil {
		logData.Message = "GetDebitCardDetails: Error Cvv Unable to decrypt" + err.Error()
		s.LoggerService.LogError(logData)
		return nil, err
	}

	userRes := responses.NewDebitCardDetailRes()

	if err = userRes.Bind(*result, cardData.IsPhysicalCardGenerated, cardData.IsVirtualCardGenerated, cardData.IsPermanentlyBlocked); err != nil {
		logData.Message = "DebitCardDetail: Error binding response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	value, _ := json.Marshal(userRes)
	ency, err := security.Encrypt(value, []byte(authValue.Key))
	if err != nil {
		logData.Message = "DebitCardDetail: Error encrypting response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "DebitCardDetail: Response received successfully"
	logData.ResponseSize = len(ency)
	logData.EndTime = time.Now()
	logData.ResponseBody = string(ency)
	s.LoggerService.LogInfo(logData)

	return ency, nil
}

func (s *Store) GetTransactionLimit(ctx context.Context, auth *models.AuthValues, request *requests.GetTransactionLimitReq) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.DEBITCARD,
		RequestURI: "/api/debitcard/get-limit-list",
		Message:    "DebitCardDetail log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
		AppVersion: utils.GetAppVersionFromContext(ctx),
	}
	txnid, err := security.GenerateRandomUUID(20)
	if err != nil {
		logData.Message = "GetTransactionLimit: Error while Generating Transaction Id"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	transactionID := strings.ReplaceAll(txnid, "-", "")
	dIndex := ""

	publicKey, enid, cid, err := s.PublicKeyAndLogin(ctx, transactionID, auth)
	if err != nil {
		logData.Message = "GetTransactionLimit: Error while Getting Public Key"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	if request.TransactionType == "Intenational" {
		dIndex = constants.IntenationalIndex
		// dIndex = "11,12,13,14"
	} else if request.TransactionType == "Domestic" {
		dIndex = constants.DomesticIndex
	}

	ftchTransaction, err := s.FetchTransaction(ctx, auth.UserId, enid, publicKey, dIndex, transactionID, cid)
	if err != nil {
		logData.Message = "GetTransactionLimit: Error while fetching transactions Limits"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	res := responses.NewGetTransactionLimit()
	result := res.StrucherRes(*ftchTransaction, request.TransactionType, transactionID)
	value, _ := json.Marshal(result)
	encrypted, err := security.Encrypt(value, []byte(auth.Key))
	if err != nil {
		return nil, err
	}
	logData.Message = "GetTransactionLimit: Transaction limit Get successfull"
	logData.EndTime = time.Now()
	logData.ResponseBody = fmt.Sprintf("", result)
	s.LoggerService.LogInfo(logData)
	return encrypted, nil
}

func (s *Store) SetTransactionLimit(ctx context.Context, auth *models.AuthValues, request requests.RequestEditTransaction) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.DEBITCARD,
		RequestURI: "/api/debitcard/settxn-limit",
		Message:    "SetTransactionLimit log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
		AppVersion: utils.GetAppVersionFromContext(ctx),
	}

	txnid, err := security.GenerateRandomUUID(20)
	if err != nil {
		logData.Message = "GetTransactionLimit: Error while Generating Transaction Id"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	transactionID := strings.ReplaceAll(txnid, "-", "")

	publicKey, enid, cid, err := s.PublicKeyAndLogin(ctx, transactionID, auth)
	if err != nil {
		logData.Message = "SetTransactionLimit: Error while fetching public key"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	listCard, err := s.ListCardControl(ctx, cid, enid, publicKey, transactionID)
	if err != nil {
		logData.Message = "SetTransactionLimit: Error while fetching ListCardControl"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	accountData, err := models.GetUserAndAccountDetailByUserID(s.db, auth.UserId)
	if err != nil {
		logData.Message = "SetTransactionLimit: Error while fetching ListCardControl"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	debitCardData, err := models.GetDebitCardData(s.db, auth.UserId)
	if err != nil {
		logData.Message = "SetTransactionLimit: Error in EditTransactionList"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	var otpType string
	if request.ReqData != nil {
		if request.ReqData[0].Type == "Domestic" {
			otpType = "SetDomesticCardLimit"
		} else if request.ReqData[0].Type == "Intenational" {
			otpType = "SetInternationalCardLimit"
		} else {
			return nil, errors.New("please specify type")
		}
	}

	optReq := requests.NewSetDebitCardPinOTP()

	if err := optReq.Bind(accountData.Applicant_id, accountData.AccountNumber, transactionID, otpType, ""); err != nil {
		logData.Message = "SetTransactionLimit: Error while binding data from debitcard Pin set"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	var otpRes *responses.SetDebitCardOTPResponse
	var opErr error

	otpRes, opErr = s.debitcardservice.SendOTPForDebitCard(ctx, optReq)
	if opErr != nil {
		bankErr := s.bankservice.HandleBankSpecificError(opErr, func(errorCode string) (string, bool) {
			return constants.GetOTPErrorMessage(errorCode)
		})

		if bankErr != nil {
			logData.Message = fmt.Sprintf("Bank error encountered (ErrorCode: %s)", bankErr.ErrorCode)
			s.LoggerService.LogError(logData)

			if msg, retryable := constants.GetOTPRetryErrorMessage(bankErr.ErrorCode); retryable {
				err := utils.RetryFunc(func() error {
					txnID, err := security.GenerateRandomUUID(20)
					if err != nil {
						logData.Message = "SetTransactionLimit: Error while Generating Transaction Id"
						s.LoggerService.LogError(logData)
						return err
					}
					optReq.TxnIdentifier = txnID
					otpRes, opErr = s.debitcardservice.SendOTPForDebitCard(ctx, optReq)
					return opErr
				}, 2)

				if err != nil {
					logData.Message = "SetTransactionLimit: Callback failed after retries"
					s.LoggerService.LogError(logData)
					return nil, errors.New(msg)
				}
			} else {
				return nil, errors.New(bankErr.ErrorMessage)
			}
		}

		return nil, opErr
	}

	req := requests.NewEditTransactionRequest()
	if err := req.Bind(request, listCard[0].CNID, cid, debitCardData.Enrollment_id.String, publicKey); err != nil {
		logData.Message = "SetTransactionLimit: Error while Setting DebitCard pin"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	// Data Saved to Cache
	if err := s.memory.Set(fmt.Sprintf("user:debitcard:transaction:%s", auth.UserId), transactionID, time.Minute*5); err != nil {
		return nil, err
	}

	requestData, err := json.Marshal(req)
	if err != nil {
		logData.Message = "SetTransactionLimit: Error while fetching ListCardControl"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if err := s.memory.Set(fmt.Sprintf("user:debitcard_transaction_limit:request:%s", auth.UserId), string(requestData), time.Minute*5); err != nil {
		return nil, err
	}

	logData.Message = "SetTransaction: OTP sent Successfully"
	logData.ResponseBody = fmt.Sprintf("", otpRes)
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)
	return nil, nil
}

func (s *Store) GetCardStatus(ctx context.Context, auth *models.AuthValues) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.DEBITCARD,
		RequestURI: "/api/debitcard/get-card-status",
		Message:    "GetCardStatus log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
		AppVersion: utils.GetAppVersionFromContext(ctx),
	}

	txnid, err := security.GenerateRandomUUID(20)
	if err != nil {
		logData.Message = "GetTransactionLimit: Error while Generating Transaction Id"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	transactionID := strings.ReplaceAll(txnid, "-", "")
	publicKey, enid, cid, err := s.PublicKeyAndLogin(ctx, transactionID, auth)
	if err != nil {
		logData.Message = "GetCardStatus: Error while getting PublicKey"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	result, err := s.ListCardControl(ctx, cid, enid, publicKey, transactionID)
	if err != nil {
		logData.Message = "GetCardStatus: Error while getting ListCardControl"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	debitCardData, err := models.GetDebitCardData(s.db, auth.UserId)
	if err != nil {
		logData.Message = "GetCardStatus: Error while getting debitcard data from DB"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	res := responses.NewDebitCardBlockStatusResponse()

	if err := res.Bind(result, debitCardData.IsPermanentlyBlocked); err != nil {
		logData.Message = "GetCardStatus: Error while binding data for response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	value, _ := json.Marshal(res)
	encrypted, err := security.Encrypt(value, []byte(auth.Key))
	if err != nil {
		return nil, err
	}
	logData.Message = "GetCardStatus: Status retrieved successfully."
	logData.ResponseBody = fmt.Sprintf("", res)
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)
	return encrypted, nil
}

func (s *Store) SetCardStatus(ctx context.Context, auth *models.AuthValues, request *requests.SetCardStatus) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.DEBITCARD,
		Message:    "SetCardStatus log",
		RequestURI: "/api/debitcard/set-card-status",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
		AppVersion: utils.GetAppVersionFromContext(ctx),
	}

	debitCardData, err := models.GetDebitCardData(s.db, auth.UserId)
	if err != nil {
		logData.Message = "SetCardStatus: Error while getting debitcard data from DB " + err.Error()
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if debitCardData.IsPermanentlyBlocked {
		logData.Message = "SetCardStatus: DebitCard is Permanently Blocked"
		s.LoggerService.LogError(logData)
		return nil, errors.New("your debit card is permanently blocked")
	}

	txnid, err := security.GenerateRandomUUID(20)
	if err != nil {
		logData.Message = "SetCardStatus: Error while Generating Transaction Id"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	transactionID := strings.ReplaceAll(txnid, "-", "")
	publicKey, enid, cid, err := s.PublicKeyAndLogin(ctx, transactionID, auth)
	if err != nil {
		logData.Message = "SetCardStatus: Error while getting PublicKey"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	cardControl, err := s.ListCardControl(ctx, cid, enid, publicKey, transactionID)
	if err != nil {
		logData.Message = "SetCardStatus: Error while getting ListCardControl"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	accountData, err := models.GetUserAndAccountDetailByUserID(s.db, auth.UserId)
	if err != nil {
		logData.Message = "SetCardStatus: Error while fetching ListCardControl"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	otpType := ""
	if request.DomesticStatus == "1" || request.InternationalStatus == "1" {
		otpType = "SetCardBlock"
	} else if request.DomesticStatus == "0" || request.InternationalStatus == "0" {
		otpType = "SetCardUnblock"
	}

	optReq := requests.NewSetDebitCardPinOTP()

	if err := optReq.Bind(accountData.Applicant_id, accountData.AccountNumber, transactionID, otpType, ""); err != nil {
		logData.Message = "SetCardStatus: Error while binding data from debitcard Pin set"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	var otpRes *responses.SetDebitCardOTPResponse
	var opErr error

	otpRes, opErr = s.debitcardservice.SendOTPForDebitCard(ctx, optReq)
	if opErr != nil {
		bankErr := s.bankservice.HandleBankSpecificError(opErr, func(errorCode string) (string, bool) {
			return constants.GetOTPErrorMessage(errorCode)
		})

		if bankErr != nil {
			logData.Message = fmt.Sprintf("Bank error encountered (ErrorCode: %s)", bankErr.ErrorCode)
			s.LoggerService.LogError(logData)

			if msg, retryable := constants.GetOTPRetryErrorMessage(bankErr.ErrorCode); retryable {
				err := utils.RetryFunc(func() error {
					txnID, err := security.GenerateRandomUUID(20)
					if err != nil {
						logData.Message = "SetTransactionLimit: Error while Generating Transaction Id"
						s.LoggerService.LogError(logData)
						return err
					}
					optReq.TxnIdentifier = txnID
					otpRes, opErr = s.debitcardservice.SendOTPForDebitCard(ctx, optReq)
					return opErr
				}, 2)

				if err != nil {
					logData.Message = "SetTransactionLimit: Callback failed after retries"
					s.LoggerService.LogError(logData)
					return nil, errors.New(msg)
				}
			} else {
				return nil, errors.New(bankErr.ErrorMessage)
			}
		}

		return nil, opErr
	}

	req := requests.CardBlockRequest{}
	req.Bind(cid, enid, cardControl[0].CNID, request.InternationalStatus, request.DomesticStatus, publicKey)

	requestData, err := json.Marshal(req)
	if err != nil {
		logData.Message = "SetCardStatus: Error while Marshal request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if err := s.memory.Set(fmt.Sprintf("user:debitcard_block:request:%s", auth.UserId), string(requestData), time.Minute*5); err != nil {
		return nil, err
	}

	if err := s.memory.Set(fmt.Sprintf("user:debitcard:is_permanently_blocked:%s", auth.UserId), request.IsPermanentlyBlocked, time.Minute*5); err != nil {
		return nil, err
	}

	if err := s.memory.Set(fmt.Sprintf("user:debitcard:transaction:%s", auth.UserId), transactionID, time.Minute*5); err != nil {
		return nil, err
	}

	logData.Message = "SetTransaction: OTP sent Successfully"
	logData.ResponseBody = fmt.Sprintf("", otpRes)
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)
	return nil, nil
}

func (s *Store) EditTransaction(ctx context.Context, authValue *models.AuthValues, transactionId string) (*responses.EditTransactionResponse, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.DEBITCARD,
		Message:    "EditTransactionLimit log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
		AppVersion: utils.GetAppVersionFromContext(ctx),
	}

	requestData, err := s.memory.Get(fmt.Sprintf("user:debitcard_transaction_limit:request:%s", authValue.UserId))
	if err != nil {
		return nil, err
	}

	req := requests.NewEditTransactionRequest()
	if err := json.Unmarshal([]byte(requestData), req); err != nil {
		logData.Message = "EditTransactionLimit: Error in Unmarshalling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	res, err := s.debitCardControlService.EditTransaction(ctx, req, transactionId)
	if err != nil {
		logData.Message = "EditTransactionLimit: Error while Updating Transaction Limit"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if err := s.memory.Delete(fmt.Sprintf("user:debitcard_transaction_limit:request:%s", authValue.UserId)); err != nil {
		return nil, err
	}

	if res == nil {
		logData.Message = "EditTransactionLimit: Empty response from bank"
		s.LoggerService.LogError(logData)
		return nil, errors.New("empty response from bank")
	}

	// save audit log
	request, err := json.Marshal(req)
	if err != nil {
		logData.Message = "EditTransactionLimit:error while marshalling"
		logData.EndTime = time.Now()
		s.LoggerService.LogError(logData)
		return res, nil
	}

	encryptReq, err := security.Encrypt(request, []byte(authValue.Key))
	if err != nil {
		logData.Message = "UpdateAdressLog:error while encrypting request"
		logData.EndTime = time.Now()
		s.LoggerService.LogError(logData)
		return res, nil
	}

	if err := s.auditLogSrv.Save(ctx, &services.AuditLog{
		TransactionID:  transactionId,
		UserID:         authValue.UserId,
		RequestURL:     "/fintech/card/transaction/edit",
		HTTPMethod:     "POST",
		ResponseStatus: 200,
		Action:         constants.DEBIT_CARD_TRANSACTION_LIMIT,
		RequestBody:    encryptReq,
	}); err != nil {
		logData.Message = "EditTransactionLimit:error while saving audit log" + err.Error()
		logData.EndTime = time.Now()
		s.LoggerService.LogInfo(logData)
	}

	return res, nil
}

func (s *Store) CardBlock(ctx context.Context, authValue *models.AuthValues, transactionId, otpType string) (*responses.CardBlockResponse, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.DEBITCARD,
		Message:    "CardBlockUnblock log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
		AppVersion: utils.GetAppVersionFromContext(ctx),
	}

	requestData, err := s.memory.Get(fmt.Sprintf("user:debitcard_block:request:%s", authValue.UserId))
	if err != nil {
		return nil, err
	}

	request := &requests.CardBlockRequest{}
	if err := json.Unmarshal([]byte(requestData), request); err != nil {
		logData.Message = "EditTransactionLimit: Error in Unmarshalling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	res, err := s.debitCardControlService.CardBlock(ctx, request, transactionId)
	if err != nil {
		logData.Message = "CardBlockUnblock: Error while Updating Card Block Unblock Status"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if res == nil {
		logData.Message = "CardBlockUnblock: response getting empty from bank"
		s.LoggerService.LogError(logData)
		return nil, errors.New("response getting empty from bank ")
	}

	if otpType == "SetCardBlockPermanently" {
		debitCard := &models.DebitCardData{
			UserID:               authValue.UserId,
			IsPermanentlyBlocked: true,
		}
		if err := models.UpdateDebitCardData(s.db, debitCard); err != nil {
			logData.Message = "CardBlockUnblock: Error while Updating Permanently Card Block Status"
			s.LoggerService.LogError(logData)
			return nil, err
		}
	}

	if err := s.memory.Delete(fmt.Sprintf("user:debitcard_block:request:%s", authValue.UserId)); err != nil {
		return nil, err
	}

	// save audit log
	req, err := json.Marshal(request)
	if err != nil {
		logData.Message = "CardBlockUnblock:error while marshalling"
		logData.EndTime = time.Now()
		s.LoggerService.LogError(logData)
		return res, nil
	}

	encryptReq, err := security.Encrypt(req, []byte(authValue.Key))
	if err != nil {
		logData.Message = "UpdateAdressLog:error while encrypting request"
		logData.EndTime = time.Now()
		s.LoggerService.LogError(logData)
		return res, nil
	}

	if err := s.auditLogSrv.Save(ctx, &services.AuditLog{
		TransactionID:  transactionId,
		UserID:         authValue.UserId,
		RequestURL:     "/fintech/card/control/edit",
		HTTPMethod:     "POST",
		ResponseStatus: 200,
		Action:         constants.DEBIT_CARD_BLOCK_UNBLOCK,
		RequestBody:    encryptReq,
	}); err != nil {
		logData.Message = "CardBlockUnblock:error while saving audit log" + err.Error()
		logData.EndTime = time.Now()
		s.LoggerService.LogInfo(logData)
	}

	return res, nil
}

func (s *Store) PublicKeyAndLogin(ctx context.Context, transactionId string, auth *models.AuthValues) (string, string, string, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.DEBITCARD,
		Message:    "GetPublicKeyAndLogin log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
		AppVersion: utils.GetAppVersionFromContext(ctx),
	}

	publicKey := ""
	dcdata, err := models.GetDebitCardData(s.db, auth.UserId)
	if err != nil {
		logData.Message = "GetPublicKeyAndLogin: Error while getting DebitCard data by userid"
		s.LoggerService.LogError(logData)
		return "", "", "", errors.New("debit card not found")
	}

	if !dcdata.PublicKey.Valid {
		publicKey, err = s.KeyFetch(ctx, transactionId)
		if err != nil {
			logData.Message = "GetPublicKeyAndLogin: Error while fetching Public Key"
			s.LoggerService.LogError(logData)
			return "", "", "", err
		}

		if err := models.UpdateDebitCardData(s.db, &models.DebitCardData{
			UserID:    auth.UserId,
			PublicKey: types.FromString(publicKey),
		}); err != nil {
			logData.Message = "TrackDebitCardStatus: Error updating debit card data"
			s.LoggerService.LogError(logData)
			return "", "", "", err
		}
	} else {
		publicKey = dcdata.PublicKey.String
	}

	// First call handle login and registration with a single API Task #911
	login, err := s.LoginAndRegister(ctx, auth.UserId, publicKey, transactionId)
	if err != nil {
		logData.Message = "GetPublicKeyAndLogin: error while login and register"
		s.LoggerService.LogError(logData)
		return "", "", "", err
	}

	if login.Desc == "Registration success" {
		// This was a first-time registration
		logData.Message = "First time registration successful"
		s.LoggerService.LogInfo(logData)
	} else if login.Desc == "Login Successful" {
		// This was a login for an already registered user
		logData.Message = "This was a login for an already registered user"
		s.LoggerService.LogInfo(logData)
	}

	if login.CustomerID == "" { // If in first time customerID is not getting then it will go inside this
		login, err = s.LoginAndRegister(ctx, auth.UserId, publicKey, transactionId)
		if err != nil {
			logData.Message = "GetPublicKeyAndLogin: error while login and register"
			s.LoggerService.LogError(logData)
			return "", "", "", err
		}
	}

	enid := ""

	if !dcdata.Enrollment_id.Valid {
		debitCardData, err := s.getDebitCardDetails(ctx, auth)
		if err != nil {
			logData.Message = "GetPublicKeyAndLogin: Error while getting DebitCard data"
			s.LoggerService.LogError(logData)
			return "", "", "", err
		}
		if debitCardData.EncryptedPAN == "" {
			logData.Message = "GetPublicKeyAndLogin: DebitCard Number Getting Null"
			s.LoggerService.LogError(logData)
			return "", "", "", errors.New("Error while getting debitcard")
		}
		add_card, err := s.AddCard(ctx, auth.UserId, login.CustomerID, publicKey, transactionId, debitCardData)
		if err != nil {
			logData.Message = "GetPublicKeyAndLogin: Error while Adding card"
			s.LoggerService.LogError(logData)
			return "", "", "", err
		}
		enid = add_card.ENID
	} else {
		enid = dcdata.Enrollment_id.String
	}

	return publicKey, enid, login.CustomerID, nil
}

func (s *Store) ListCardControl(ctx context.Context, cid, enid, publicKey, transactionID string) ([]responses.ListCardControlResponse, error) {
	logData := &commonSrv.LogEntry{
		Action:    constants.DEBITCARD,
		Message:   "ListCardControl log",
		UserID:    utils.GetUserIDFromContext(ctx),
		RequestID: utils.GetRequestIDFromContext(ctx),
	}
	request := requests.ListCardControlRequest{}
	request.Bind(cid, enid, publicKey)
	response, err := s.debitCardControlService.ListCardControl(ctx, &request, transactionID)
	if err != nil {
		logData.Message = "ListCardControl:Error while getting ListCardControl Data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if response == nil {
		logData.Message = "ListCardControl: response getting empty from bank"
		s.LoggerService.LogError(logData)
		return nil, errors.New("response getting empty from bank ")
	}

	return response, nil
}

func (s *Store) KeyFetch(ctx context.Context, transactionID string) (string, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.DEBITCARD,
		Message:    "KeyFetchApi log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
		AppVersion: utils.GetAppVersionFromContext(ctx),
	}
	req := requests.KeyFetchRequest{
		ID: "66",
	}

	response, err := s.debitCardControlService.KeyFetch(ctx, req, transactionID)
	if err != nil {
		logData.Message = "KeyFetchApi: Error while getting PublicKey from bank" + err.Error()
		s.LoggerService.LogError(logData)
		return "", err
	}

	if response == nil {
		logData.Message = "GetPublicKeyAndLogin: PublicKey is getting null"
		s.LoggerService.LogError(logData)
		return "", errors.New("public key fetching error")
	}

	return response.PublicKey, nil
}

func (s *Store) LoginAndRegister(ctx context.Context, userId, publicKey, transactionID string) (*responses.LoginResponse, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.DEBITCARD,
		Message:    "LoginAndRegister log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
		AppVersion: utils.GetAppVersionFromContext(ctx),
	}
	UserDetail, err := models.GetUserAndAccountDetailByUserID(s.db, userId)
	if err != nil {
		logData.Message = "LoginAndRegister: Error while getting UserAnd Account Data by User ID"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	req := requests.LoginRequest{}
	name := NameOnDebitCard(UserDetail.FirstName, UserDetail.MiddleName, UserDetail.LastName)
	if name == "" {
		logData.Message = "LoginAndRegister: Error while getting Name on DebitCard"
		s.LoggerService.LogError(logData)
		return nil, errors.New("name should not be greater then 50")
	}

	req.Bind(UserDetail.MobileNumber, UserDetail.Email, UserDetail.AccountNumber, name, UserDetail.CustomerId, UserDetail.DateOfBirth)

	req.PublicKey = publicKey

	response, err := s.debitCardControlService.Login(ctx, &req, transactionID)
	if err != nil {
		logData.Message = "LoginAndRegister: Error while Login for DebitCard"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if response == nil {
		logData.Message = "LoginAndRegister: response getting empty from bank"
		s.LoggerService.LogError(logData)
		return nil, errors.New("response getting empty from bank ")
	}

	if response.CustomerID != "" {
		if err := models.UpdateDebitCardData(s.db, &models.DebitCardData{ //  if CID not empty then Update in DB
			UserID: userId,
			CID:    types.FromString(response.CustomerID),
		}); err != nil {
			logData.Message = "TrackDebitCardStatus: Error updating debit card data"
			s.LoggerService.LogError(logData)
			return nil, err
		}
	}

	return response, nil
}

func (s *Store) AddCard(ctx context.Context, userId, cid, publicKey, transactionID string, card *responses.DebitCardDetailRes) (*responses.AddCardResponse, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.DEBITCARD,
		Message:    "AddCard log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
		AppVersion: utils.GetAppVersionFromContext(ctx),
	}
	request := requests.AddCardRequest{}
	debitcard := card
	expYearIn4 := strings.Split(debitcard.ExpiryDate, "/")[1]
	month := strings.Split(debitcard.ExpiryDate, "/")[0]
	expYearIn2 := expYearIn4 //[len(expYearIn4)-2:]

	cardNo := debitcard.EncryptedPAN
	name := debitcard.CardholderName

	request.Bind(name, expYearIn2, month, cardNo, cid, publicKey)

	response, err := s.debitCardControlService.AddCard(ctx, &request, transactionID)
	if err != nil {
		logData.Message = "AddCard: Error while AddCard"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if response == nil {
		logData.Message = "AddCard: response getting empty from bank"
		s.LoggerService.LogError(logData)
		return nil, errors.New("response getting empty from bank ")
	}

	err = models.UpdateEnrollmentID(s.db, userId, response.ENID)
	if err != nil {
		logData.Message = "AddCard: Error While Updating EnrollmentId"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	return response, err
}

func (s *Store) ListCard(ctx context.Context, cid, publickey, transactionId string) ([]responses.ListCardResponse, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.DEBITCARD,
		Message:    "ListCard log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
		AppVersion: utils.GetAppVersionFromContext(ctx),
	}
	request := requests.ListCardRequest{}
	request.Bind(cid, publickey)
	response, err := s.debitCardControlService.ListCard(ctx, &request, transactionId)
	if err != nil {
		logData.Message = "ListCard: Error while List Card API"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if response == nil {
		logData.Message = "ListCard: response getting empty from bank"
		s.LoggerService.LogError(logData)
		return nil, errors.New("response getting empty from bank ")
	}

	return response, nil
}

func (s *Store) FetchTransaction(ctx context.Context, userId, enid, publicKey, deliveryChannel, transactionId, CustomerId string) (*responses.FetchTransactionResponse, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.DEBITCARD,
		Message:    "FetchTransaction log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
		AppVersion: utils.GetAppVersionFromContext(ctx),
	}

	request := requests.FetchTransactionRequest{}
	request.Bind(CustomerId, enid, publicKey, deliveryChannel)

	response, err := s.debitCardControlService.FetchTransaction(ctx, &request, transactionId)
	if err != nil {
		logData.Message = "FetchTransaction: Error while FetchingTransaction Limit Detail"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if response == nil {
		logData.Message = "FetchTransaction: response getting empty from bank"
		s.LoggerService.LogError(logData)
		return nil, errors.New("response getting empty from bank ")
	}

	return response, nil
}

func (s *Store) SetDebitCardPin(ctx context.Context, authValue *models.AuthValues, request *requests.SetDebitCardPinReq) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.DEBITCARD,
		RequestURI: "/api/debitcard/set-debitcard-pin",
		Message:    "SetDebitCardPin log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
		AppVersion: utils.GetAppVersionFromContext(ctx),
	}

	account, err := models.GetUserAndAccountDetailByUserID(s.db, authValue.UserId)

	if err != nil {
		logData.Message = "SetDebitCardPin: Error getting Account data by user id"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	_, err = models.GetDebitCardData(s.db, authValue.UserId)

	if err != nil {
		logData.Message = "SetDebitCardPin: Error getting DebitCard data by user id"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	optReq := requests.NewSetDebitCardPinOTP()

	if err := optReq.Bind(account.Applicant_id, account.AccountNumber, "", request.PinSetType, ""); err != nil {
		logData.Message = "SetDebitCardPin: Error while binding data from debitcard Pin set"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	var otpRes *responses.SetDebitCardOTPResponse
	var opErr error

	otpRes, opErr = s.debitcardservice.SendOTPForDebitCard(ctx, optReq)
	if opErr != nil {
		bankErr := s.bankservice.HandleBankSpecificError(opErr, func(errorCode string) (string, bool) {
			return constants.GetDebitCardFetchErrorMessage(errorCode)
		})

		if bankErr != nil {
			logData.Message = fmt.Sprintf("Bank error encountered (ErrorCode: %s)", bankErr.ErrorCode)
			s.LoggerService.LogError(logData)

			if msg, retryable := constants.GetDebitCardFetchRetryErrorMessage(bankErr.ErrorCode); retryable {
				err := utils.RetryFunc(func() error {
					txnID, err := security.GenerateRandomUUID(20)
					if err != nil {
						logData.Message = "SetDebitCardPin: Error while Generating Transaction Id"
						s.LoggerService.LogError(logData)
						return err
					}
					optReq.TxnIdentifier = txnID
					otpRes, opErr = s.debitcardservice.SendOTPForDebitCard(ctx, optReq)
					return opErr
				}, 2)

				if err != nil {
					logData.Message = "SetDebitCardPin: Callback failed after retries"
					s.LoggerService.LogError(logData)
					return nil, errors.New(msg)
				}
			} else {
				return nil, errors.New(bankErr.ErrorMessage)
			}
		}

		return nil, opErr
	}

	if err := s.memory.Set(fmt.Sprintf("user:debitcard:transaction:%s", authValue.UserId), otpRes.TxnIdentifier, time.Minute*5); err != nil {
		return nil, err
	}

	if err := s.memory.Set(fmt.Sprintf("user:debitcard:pin:%s", authValue.UserId), request.Pin, time.Minute*5); err != nil {
		return nil, err
	}

	logData.Message = "SetDebitCardPin: Response received successfully"
	logData.EndTime = time.Now()
	logData.ResponseBody = fmt.Sprintf("", otpRes)
	s.LoggerService.LogInfo(logData)

	return nil, nil
}

func (s *Store) VerifyDebitCardOTP(ctx context.Context, authValue *models.AuthValues, request *requests.DebitCardVerifyOtpReq) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.DEBITCARD,
		RequestURI: "/api/debitcard/verify-debitcard-pin",
		Message:    "VerifyDebitCardOTP log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
		AppVersion: utils.GetAppVersionFromContext(ctx),
	}
	userData, err := models.GetUserAndAccountDetailByUserID(s.db, authValue.UserId)
	if err != nil {
		logData.Message = "VerifyDebitCardOTP: Error getting User data by user id"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	txnIdentifier, err := s.memory.Get(fmt.Sprintf("user:debitcard:transaction:%s", authValue.UserId))
	if err != nil {
		return nil, err
	}

	optReq := requests.NewSetDebitCardPinOTP()

	if err := optReq.Bind(userData.Applicant_id, userData.AccountNumber, txnIdentifier, request.OtpType, request.Otp); err != nil {
		logData.Message = "VerifyDebitCardOTP: Error while binding data from debitcard Pin set"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	var optRes *responses.VerifyDebitCardOTPResponse
	var opErr error

	optRes, opErr = s.debitcardservice.VerifyOTPForDebitCard(ctx, optReq)
	if opErr != nil {
		bankErr := s.bankservice.HandleBankSpecificError(opErr, func(errorCode string) (string, bool) {
			return constants.GetOTPErrorMessage(errorCode)
		})

		if bankErr != nil {
			logData.Message = fmt.Sprintf("Bank error encountered (ErrorCode: %s)", bankErr.ErrorCode)
			s.LoggerService.LogError(logData)

			if msg, retryable := constants.GetOTPRetryErrorMessage(bankErr.ErrorCode); retryable {
				err := utils.RetryFunc(func() error {
					txnID, err := security.GenerateRandomUUID(20)
					if err != nil {
						logData.Message = "verifyDebitCardOTP: Error while Generating Transaction Id"
						s.LoggerService.LogError(logData)
						return err
					}
					optReq.TxnIdentifier = txnID
					optRes, opErr = s.debitcardservice.VerifyOTPForDebitCard(ctx, optReq)
					return opErr
				}, 2)

				if err != nil {
					logData.Message = "verifyDebitCardOTP: Callback failed after retries"
					s.LoggerService.LogError(logData)
					return nil, errors.New(msg)
				}
			} else {
				return nil, errors.New(bankErr.ErrorMessage)
			}
		}

		return nil, opErr
	}

	cardData, err := models.GetDebitCardData(s.db, authValue.UserId)

	if err != nil {
		logData.Message = "SetDebitCardPin: Error getting DebitCard data by user id"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	req := requests.NewGetDebitcardDetailRequest()

	if err := req.GetDebitCardDetailReq_Bind(userData.Applicant_id, userData.AccountNumber, cardData.Proxy_Number.String, cardData.TxnIdentifyer.String); err != nil {
		logData.Message = "SetDebitCardPin: Error while binding DebitCard data from debitcard detail"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if request.OtpType == "New" || request.OtpType == "Reset" {

		var result *responses.DebitcardDetailResponse
		var opErr error

		result, opErr = s.debitcardservice.GetDebitCardDetails(ctx, req)
		if opErr != nil {
			bankErr := s.bankservice.HandleBankSpecificError(opErr, func(errorCode string) (string, bool) {
				return constants.GetDebitCardFetchErrorMessage(errorCode)
			})

			if bankErr != nil {
				logData.Message = fmt.Sprintf("Bank error encountered (ErrorCode: %s)", bankErr.ErrorCode)
				s.LoggerService.LogError(logData)

				if msg, retryable := constants.GetDebitCardFetchRetryErrorMessage(bankErr.ErrorCode); retryable {
					err := utils.RetryFunc(func() error {
						result, opErr = s.debitcardservice.GetDebitCardDetails(ctx, req)
						return opErr
					}, 2)

					if err != nil {
						logData.Message = "verifyDebitCardOTP: Callback failed after retries"
						s.LoggerService.LogError(logData)
						return nil, errors.New(msg)
					}
				} else {
					return nil, errors.New(bankErr.ErrorMessage)
				}
			}

			return nil, opErr
		}

		pin, err := s.memory.Get(fmt.Sprintf("user:debitcard:pin:%s", authValue.UserId))
		if err != nil {
			return nil, err
		}
		req2 := requests.NewSetDebitCardPin()

		if err := req2.Bind(result, pin, txnIdentifier, request.OtpType); err != nil {
			logData.Message = "SetDebitCardPin: Error while binding data from debitcard Pin set"
			s.LoggerService.LogError(logData)
			return nil, err
		}

		_, err = s.debitcardservice.SetDebitCardPin(ctx, req2)

		if err != nil {
			logData.Message = "DebitCardDetail: Error while Setting DebitCard pin"
			s.LoggerService.LogError(logData)

			bankErr := s.bankservice.HandleBankSpecificError(err, func(errorCode string) (string, bool) {
				return constants.GetDebitCardSetPinErrorMessage(errorCode)
			})

			if bankErr != nil {
				logData.Message = fmt.Sprintf("Bank error encountered (ErrorCode: %s)", bankErr.ErrorCode)
				s.LoggerService.LogError(logData)

				if msg, retryable := constants.GetDebitCardSetPinRetryErrorMessage(bankErr.ErrorCode); retryable {
					retryErr := utils.RetryFunc(func() error {
						txnID, err := security.GenerateRandomUUID(20)
						if err != nil {
							logData.Message = "verifyDebitCardOTP: Error while Generating Transaction Id"
							s.LoggerService.LogError(logData)
							return err
						}
						req2.TxnIdentifier = txnID
						_, opErr := s.debitcardservice.SetDebitCardPin(ctx, req2)
						return opErr
					}, 2)

					if retryErr != nil {
						logData.Message = "SetDebitCardPin: Callback failed after retries"
						s.LoggerService.LogError(logData)
						return nil, errors.New(msg)
					}
				} else {
					return nil, errors.New(bankErr.ErrorMessage)
				}
			}

			return nil, err
		}

		// save audit log
		req2.PinNo = ""
		req2.EncryptedPAN = ""
		requestData, err := json.Marshal(req2)
		if err != nil {
			logData.Message = "UpdateAdressLog:error while marshalling"
			logData.EndTime = time.Now()
			s.LoggerService.LogError(logData)
			return nil, nil
		}
		encryptReq, err := security.Encrypt(requestData, []byte(authValue.Key))
		if err != nil {
			logData.Message = "UpdateAdressLog:error while encrypting request"
			logData.EndTime = time.Now()
			s.LoggerService.LogError(logData)
			return nil, nil
		}

		if err := s.auditLogSrv.Save(ctx, &services.AuditLog{
			TransactionID:  optReq.TxnIdentifier,
			UserID:         authValue.UserId,
			ApplicantID:    userData.Applicant_id,
			RequestURL:     "/api/debitcard/set-debitcard-pin",
			HTTPMethod:     "POST",
			ResponseStatus: 200,
			RequestBody:    encryptReq,
			Action:         constants.DEBIT_CARD_PIN_RESET,
		}); err != nil {
			logData.Message = "SetDebitCardPin:error while saving audit log" + err.Error()
			logData.EndTime = time.Now()
		}

	} else if request.OtpType == "SetDomesticCardLimit" || request.OtpType == "SetInternationalCardLimit" {
		_, err := s.EditTransaction(ctx, authValue, txnIdentifier)
		if err != nil {
			logData.Message = "SetTransactionLimit: Error in EditTransactionList"
			s.LoggerService.LogError(logData)
			return nil, err
		}
	} else if request.OtpType == "SetCardBlock" || request.OtpType == "SetCardUnblock" || request.OtpType == "SetCardBlockPermanently" {
		_, err = s.CardBlock(ctx, authValue, txnIdentifier, request.OtpType)
		if err != nil {
			logData.Message = "SetCardStatus: Error while getting CardBlock Status"
			s.LoggerService.LogError(logData)
			return nil, err
		}
	}

	//Delete Data from Cache
	if err := s.memory.Delete(fmt.Sprintf("user:user:debitcard:transaction:%s", authValue.UserId)); err != nil {
		return nil, err
	}

	logData.Message = "VerifyDebitCardOTP: Response received successfully"
	logData.EndTime = time.Now()
	logData.ResponseBody = fmt.Sprintf("", optRes)
	s.LoggerService.LogInfo(logData)

	return nil, nil
}

func (s *Store) TrackDebitCardStatus(ctx context.Context, authValue *models.AuthValues) (interface{}, error) {
	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		StartTime:  startTime,
		Action:     constants.DEBITCARD,
		RequestURI: "/api/debitcard/track-status",
		Message:    "TrackDebitCardStatus log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
		AppVersion: utils.GetAppVersionFromContext(ctx),
	}

	debitCardResponse, err := s.getDebitCardDetails(ctx, authValue)
	if err != nil {
		logData.Message = "TrackDebitCardStatus: Error getting debit card details"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	debitCardTrackingFile, err := s.getDebitCardTrackingFile()
	if err != nil {
		s.LoggerService.LogError(logData)
		return nil, fmt.Errorf("failed to get tracking file: %w", err)
	}
	defer debitCardTrackingFile.Close()

	cardStatus, err := s.findCardStatus(debitCardTrackingFile, debitCardResponse)
	if err != nil {
		logData.Message = "TrackDebitCardStatus: Error finding card status"
		s.LoggerService.LogError(logData)
		return nil, fmt.Errorf("failed to find card status: %w", err)
	}

	cardStatusByteData, err := json.Marshal(cardStatus)
	if err != nil {
		logData.Message = "TrackDebitCardStatus: Error marshaling card status"
		s.LoggerService.LogInfo(logData)
		return nil, fmt.Errorf("failed to marshal card status: %w", err)
	}

	encryptedData, err := security.Encrypt(cardStatusByteData, []byte(authValue.Key))
	if err != nil {
		logData.Message = "TrackDebitCardStatus: Error encrypting card status"
		s.LoggerService.LogInfo(logData)
		return nil, fmt.Errorf("failed to encrypt card status: %w", err)
	}

	if err := models.UpdateDebitCardData(s.db, &models.DebitCardData{
		UserID:         authValue.UserId,
		DeliveryStatus: types.FromString(cardStatus.DispatchStatus),
	}); err != nil {
		logData.Message = "TrackDebitCardStatus: Error updating debit card data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "TrackDebitCardStatus: Response received successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return encryptedData, nil
}

func (s *Store) getDebitCardDetails(ctx context.Context, authValue *models.AuthValues) (*responses.DebitCardDetailRes, error) {
	debitCardDetails, err := s.DebitCardDetail(ctx, authValue)
	if err != nil {
		return nil, err
	}

	decryptedData, err := security.Decrypt(debitCardDetails.(string), []byte(authValue.Key))
	if err != nil {
		return nil, err
	}

	var debitCardResponse responses.DebitCardDetailRes
	if err := json.Unmarshal([]byte(decryptedData), &debitCardResponse); err != nil {
		return nil, err
	}

	return &debitCardResponse, nil
	// return nil, nil
}

func (s *Store) getDebitCardTrackingFile() (*excelize.File, error) {
	// Check if we have a cached file and it's less than 1 hour old
	if s.cachedTrackingFile != nil && time.Since(s.lastFileDownloadTime) < time.Hour {
		return s.cachedTrackingFile, nil
	}

	var file *excelize.File
	var err error

	if trackingFilePath := os.Getenv("DEBIT_CARD_TRACK_STATUS_FILE_NAME"); trackingFilePath != "" {
		file, err = s.downloadFromS3(trackingFilePath)
	} else if ftpConfig := s.getFTPConfig(); ftpConfig != nil {
		file, err = s.downloadFromFTP(ftpConfig)
	} else if localPath := os.Getenv("LOCAL_PATH"); localPath != "" {
		file, err = excelize.OpenFile(localPath)
	} else {
		return nil, errors.New("no valid source for tracking file")
	}

	if err != nil {
		return nil, err
	}

	// Update cache
	s.cachedTrackingFile = file
	s.lastFileDownloadTime = time.Now()

	return file, nil
}

func (s *Store) getFTPConfig() *FTPConfig {
	host := os.Getenv("FTP_HOST")
	port := os.Getenv("FTP_PORT")
	username := os.Getenv("FTP_USER")
	password := os.Getenv("FTP_PASSWORD")
	filePath := os.Getenv("FTP_FILE_PATH")

	if host == "" || port == "" || username == "" || password == "" || filePath == "" {
		return nil
	}

	portInt, _ := strconv.Atoi(port)
	return &FTPConfig{
		Host:     host,
		Port:     portInt,
		User:     username,
		Password: password,
		FilePath: filePath,
	}
}

func (s *Store) downloadFromS3(fileName string) (*excelize.File, error) {
	bucketName := os.Getenv("DEBIT_CARD_TRACK_S3_BUCKET_NAME")
	data, err := s.s3Client.Download(context.Background(), bucketName, fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to download file from S3: %w", err)
	}

	return excelize.OpenReader(bytes.NewReader(data))
}

func (s *Store) downloadFromFTP(config *FTPConfig) (*excelize.File, error) {
	client := ftp_server.NewFTPClient(config.Host, config.Port, config.User, config.Password)
	if err := client.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to FTP: %w", err)
	}
	defer client.Disconnect()

	if err := client.DownloadFile("./tracking.xlsx", config.FilePath); err != nil {
		return nil, fmt.Errorf("failed to download file from FTP: %w", err)
	}

	return excelize.OpenFile("./tracking.xlsx")
}

func (s *Store) findCardStatus(file *excelize.File, cardDetails *responses.DebitCardDetailRes) (*CardStatus, error) {
	sheetName := file.GetSheetName(0)
	rows, err := file.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to read sheet: %w", err)
	}

	for _, row := range rows {
		cardStatus := newCardStatus(row)
		if cardStatus.matches(cardDetails) {
			return cardStatus, nil
		}
	}

	return nil, errors.New("debit card track status not found")
}

func (cs *CardStatus) matches(cardDetails *responses.DebitCardDetailRes) bool {
	return strings.EqualFold(strings.TrimSpace(cs.CardHolderName), strings.TrimSpace(cardDetails.CardholderName)) &&
		strings.TrimSpace(cs.CardNumber[len(cs.CardNumber)-4:]) == strings.TrimSpace(cardDetails.EncryptedPAN[len(cardDetails.EncryptedPAN)-4:])
}

func newCardStatus(row []string) *CardStatus {
	return &CardStatus{
		Date:            row[1],
		CardNumber:      row[3],
		CardHolderName:  row[4],
		ReferenceNumber: row[5],
		AWB:             row[6],
		CardType:        row[17],
		DispatchStatus:  row[18],
		DispatchMode:    row[19],
		DispatchDate:    row[20],
	}
}

func NameOnDebitCard(fullName ...string) string {
	var nameParts []string
	for _, n := range fullName {
		trimmed := strings.TrimSpace(n)
		if trimmed != "" {
			nameParts = append(nameParts, trimmed)
		}
	}

	name := strings.Join(nameParts, " ")

	if len(name) <= 50 { //full name
		return strings.TrimSpace(name)
	}

	if len(nameParts[0]+nameParts[2]) <= 50 { //first name and last name
		return strings.TrimSpace(nameParts[0] + " " + nameParts[2])
	}

	if len(nameParts[0]+nameParts[1]) <= 50 { //first name and middle name
		return strings.TrimSpace(nameParts[0] + " " + nameParts[1])
	}

	if len(nameParts[0]) <= 50 { //first name
		return strings.TrimSpace(nameParts[0])
	}

	return ""
}
