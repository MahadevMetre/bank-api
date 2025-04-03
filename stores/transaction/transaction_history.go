package transaction

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"bitbucket.org/paydoh/paydoh-commons/database"
	commonSrv "bitbucket.org/paydoh/paydoh-commons/services"
	"github.com/google/uuid"

	"bankapi/constants"
	"bankapi/models"
	"bankapi/requests"
	"bankapi/responses"
	"bankapi/services"
	"bankapi/utils"
)

type Transaction interface {
	InitiateSimVerification(mobileNumber, key string, request *requests.AuthenticationRequest) (interface{}, error)
}

type TransactionStore struct {
	db            *sql.DB
	m             *database.Document
	memory        *database.InMemory
	service       *services.BankApiService
	LoggerService *commonSrv.LoggerService
}

func NewTransactionStore(log *commonSrv.LoggerService, db *sql.DB, m *database.Document, memory *database.InMemory) *TransactionStore {
	return &TransactionStore{
		db:            db,
		m:             m,
		memory:        memory,
		LoggerService: log,
		service:       services.NewBankApiService(log, memory),
	}
}

func (ts *TransactionStore) GetTransactionData(ctx context.Context, req requests.TransactionRequest) ([]responses.Transaction, error) {
	logData := &commonSrv.LogEntry{
		Action:    constants.TRANSACTION_HISTORY,
		Message:   "GetTransactionData log",
		UserID:    utils.GetUserIDFromContext(ctx),
		RequestID: utils.GetRequestIDFromContext(ctx),
	}

	accountData, err := models.GetAccountDataByUserId(ts.db, req.UserId)
	if err != nil {
		logData.Message = "GetTransactionData: Error fetching transaction history"
		ts.LoggerService.LogError(logData)
		return nil, err
	}

	// today := time.Now()
	// oneMonthBefore := today.AddDate(0, -1, 0)
	// formattedDate := oneMonthBefore.Format("02-01-2006")
	// KVBAccountNumber: "1219155000160816",
	// FromDate: formattedDate,
	// ToDate:   today.Format("02-01-2006"),

	reqData := requests.KVBTransactionRequest{
		TxnBranch:        "1479",
		TxnIdentifier:    generateTxnIdentifier(),
		KVBAccountNumber: accountData.AccountNumber,
		ToDate:           req.ToDate,
		FromDate:         req.FromDate,
	}

    var response *responses.TransactionResponse
    var opErr error

    response, opErr = ts.service.FetchTransactionHistory(ctx, reqData)
    if opErr != nil {
        bankErr := ts.service.HandleBankSpecificError(opErr, func(errorCode string) (string, bool) {
            return constants.GetCasaTxnErrorMessage(errorCode)
        })

        if bankErr != nil {
            logData.Message = fmt.Sprintf("Bank error encountered (ErrorCode: %s)", bankErr.ErrorCode)
            ts.LoggerService.LogError(logData)

            if msg, retryable := constants.GetCasaTxnErrorRetryMessage(bankErr.ErrorCode); retryable {
                err := utils.RetryFunc(func() error {
                    reqData.TxnIdentifier = generateTxnIdentifier()
                    response, opErr = ts.service.FetchTransactionHistory(ctx, reqData)
                    return opErr
                }, 2)

                if err != nil {
                    logData.Message = "GetTransactionData failed after retries"
                    ts.LoggerService.LogError(logData)
                    return nil, errors.New(msg)
                }
            } else {
                return nil, errors.New(bankErr.ErrorMessage)
            }
        } else {
            return nil, opErr
        }
    }

	logData.Message = "GetTransactionData: Transaction data retrieved successfully"
	logData.EndTime = time.Now()
	ts.LoggerService.LogInfo(logData)

	return response.Data, nil
}

func generateTxnIdentifier() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func (ts *TransactionStore) GetTransactionDetails(req *requests.TransactionDetailsRequest, userId string) (interface{}, error) {
	urtRefNumber := utils.GetUpiUtrRefNumber(req.TransactionDescription)

	filterParam := &models.FilterParams{
		UtrRefNumber: urtRefNumber,
		TxnDate:      req.TransactionDate,
		Amount:       req.TransactionAmount,
	}

	txnData, err := models.FetchOneTransaction(ts.db, filterParam)
	if err != nil {
		return nil, err
	}

	return txnData, nil
}

func (s *TransactionStore) GetRecentTransactionUsers(ctx context.Context, authValues *models.AuthValues) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     "RECENT_TRANSACTION_USERS",
		RequestURI: "/api/transaction/recent-users",
		Message:    "GetRecentTransactionUsers log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
		StartTime:  time.Now(),
	}

	users, err := models.GetRecentTransactionUsers(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "GetRecentTransactionUsers: Error fetching recent transaction users"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "GetRecentTransactionUsers: Recent transaction users retrieved successfully"
	logData.EndTime = time.Now()
	logData.Latency = time.Since(logData.StartTime).Seconds()
	s.LoggerService.LogInfo(logData)

	return users, nil
}
