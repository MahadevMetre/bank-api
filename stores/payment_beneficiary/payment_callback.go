package payment_beneficiary

import (
	"bankapi/constants"
	"bankapi/models"
	"bankapi/requests"
	"bankapi/services"
	"bankapi/utils"
	"context"
	"database/sql"
	"time"

	"bitbucket.org/paydoh/paydoh-commons/database"
	commonSrv "bitbucket.org/paydoh/paydoh-commons/services"
	"bitbucket.org/paydoh/paydoh-commons/types"
)

type PaymentCallbackStore struct {
	db            *sql.DB
	m             *database.Document
	memory        *database.InMemory
	bankService   *services.BankApiService
	LoggerService *commonSrv.LoggerService
}

func NewPaymentCallbackStore(log *commonSrv.LoggerService, db *sql.DB, m *database.Document, memory *database.InMemory) *PaymentCallbackStore {
	bankService := services.NewBankApiService(log, memory)
	return &PaymentCallbackStore{
		db:            db,
		m:             m,
		memory:        memory,
		bankService:   bankService,
		LoggerService: log,
	}
}

func (p *PaymentCallbackStore) Update(ctx context.Context, reqData *requests.PaymentCallbackRequestData) error {
	logData := &commonSrv.LogEntry{
		Action:     constants.BANK,
		RequestURI: "/callback/normal-payment",
		Message:    "Update payment callback log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	// update transaction status in db
	if err := models.UpdateTransactionByTransID(p.db, &models.Transaction{
		TransactionID: reqData.CbsStatus[0].TransactionId,
		CBSStatus:     types.FromString(reqData.CbsStatus[0].Status),
		UTRRefNumber:  types.FromString(reqData.CbsStatus[0].UTR_Ref_Number),
	}); err != nil {
		logData.Message = "PaymentCallback Update: Error updating transaction status in db"
		p.LoggerService.LogError(logData)
		return err
	}

	logData.Message = "Update: Payment callback updated successfully"
	logData.EndTime = time.Now()
	p.LoggerService.LogInfo(logData)
	return nil
}
