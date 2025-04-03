package kyc_audit_data

import (
	"bankapi/constants"
	"bankapi/models"
	"bankapi/requests"
	"bankapi/services"
	"bankapi/utils"
	"context"
	"database/sql"
	"errors"
	"time"

	"bitbucket.org/paydoh/paydoh-commons/database"
	commonSrv "bitbucket.org/paydoh/paydoh-commons/services"
)

type KycAuditStore struct {
	db            *sql.DB
	m             *database.Document
	memory        *database.InMemory
	bankService   *services.BankApiService
	LoggerService *commonSrv.LoggerService
}

func NewKycAuditStore(log *commonSrv.LoggerService, db *sql.DB, m *database.Document, memory *database.InMemory) *KycAuditStore {
	bankService := services.NewBankApiService(log, memory)
	return &KycAuditStore{
		db:            db,
		m:             m,
		memory:        memory,
		bankService:   bankService,
		LoggerService: log,
	}
}

func (s *KycAuditStore) Create(ctx context.Context, reqData *requests.KycAuditRequestData) error {
	logData := &commonSrv.LogEntry{
		Action:     constants.KYC_AUDIT_DATA,
		RequestURI: "/callback/audit-complition",
		Message:    "Create KYC audit creation",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	data, err := models.GetUserDataByApplicantId(s.db, reqData.ApplicantId)
	if err != nil {
		logData.Message = "Create KYC audit creation: Error getting user data by applicant id"
		s.LoggerService.LogError(logData)
		return err
	}

	existsData, _ := models.GetRecordByUserId(data.UserId)
	if existsData != nil {
		logData.Message = "Create KYC audit creation: Error kyc audit data already found for this user"
		s.LoggerService.LogError(logData)
		return errors.New("kyc audit data exists for this user")
	}

	// check if already kyc consent completed
	_, err = models.FindKycConsentByUserId(s.db, data.UserId)
	if err != nil {
		logData.Message = "Create KYC audit creation: Error finding kyc consent by user id"
		s.LoggerService.LogError(logData)
		return errors.New("kyc consent data not found for this user")
	}

	// kyc_update_data status should be 'success', should be present in DB
	kycUpdate, err := models.GetKycUpdateDataByUserId(s.db, data.UserId)
	if err != nil {
		logData.Message = "Create KYC audit creation: Error getting kyc update data by user id"
		s.LoggerService.LogError(logData)
		return errors.New("kyc update data not found for this user")
	}

	if kycUpdate.Status != "SUCCESS" {
		logData.Message = "Create KYC audit creation: Error kyc update has failed for this user"
		s.LoggerService.LogError(logData)
		return errors.New("kyc failed for this user")
	}

	reqData.UserId = data.UserId
	model := models.NewKycAuditData(reqData)

	err = model.CreateRecord(s.db)
	if err != nil {
		return err
	}

	// update is_active account status
	err = models.UpdateAccountByUserId(&models.AccountDataUpdate{
		IsActive: true,
	}, data.UserId)
	if err != nil {
		logData.Message = "Create KYC audit creation: Error updating user account status"
		s.LoggerService.LogError(logData)
		return err
	}

	// update onboarding status
	if err := models.UpdateUserOnboardingStatus(constants.AUDIT_CALLBACK_STEP, data.UserId); err != nil {
		logData.Message = "Create audit error while updating onboarding status"
		s.LoggerService.LogError(logData)
	}

	if kycUpdate.Status == "SUCCESS" {
		if err := models.GenerateNotification(data.UserId, "AccountActivation", "", ""); err != nil {
			logData.Message = "UpdateAdressLog: Error while Generating Notification " + err.Error()
			s.LoggerService.LogError(logData)
		}
	}

	logData.Message = "Create: KYC audit record created successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return nil

}
