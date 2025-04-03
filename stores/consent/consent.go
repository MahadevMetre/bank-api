package consent

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"bitbucket.org/paydoh/paydoh-commons/database"
	commonSrv "bitbucket.org/paydoh/paydoh-commons/services"

	"bankapi/constants"
	"bankapi/models"
	"bankapi/requests"
	"bankapi/services"
	"bankapi/utils"
)

type Consent interface {
	UpdateConsent(request *requests.ConsentRequest, authValues *models.AuthValues) (interface{}, error)
}

type Store struct {
	db            *sql.DB
	bankService   *services.BankApiService
	memory        *database.InMemory
	m             *database.Document
	LoggerService *commonSrv.LoggerService
}

func NewStore(log *commonSrv.LoggerService, db *sql.DB, m *database.Document, memory *database.InMemory) *Store {
	return &Store{
		db:            db,
		m:             m,
		memory:        memory,
		LoggerService: log,
		bankService:   services.NewBankApiService(log, memory),
	}
}

func (s *Store) UpdateConsent(ctx context.Context, request *requests.ConsentRequest, authValues *models.AuthValues) error {
	logData := &commonSrv.LogEntry{
		Action:     constants.CONSENT,
		RequestURI: "/api/consent/update",
		Message:    "Update Consent log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	existingDevice, err := models.GetUserDataByUserId(s.db, authValues.UserId)

	if err != nil {
		logData.Message = "UpdateConsent: Error getting user data"
		s.LoggerService.LogError(logData)
		return err
	}

	outGoingRequest := requests.NewOutgoingConsentRequest()
	if err := outGoingRequest.BindAndValidate(existingDevice.ApplicantId); err != nil {
		logData.Message = "UpdateConsent: Error binding and validating consent request"
		s.LoggerService.LogError(logData)
		return err
	}

	consentResponse, err := s.bankService.PostUserConsent(ctx, outGoingRequest, nil)
	if err != nil {
		if bankErr := s.bankService.HandleBankSpecificError(err, func(errorCode string) (string, bool) {
			return constants.GetConsentErrorMessage(errorCode)
		}); bankErr != nil {
			return errors.New(bankErr.ErrorMessage)
		}

		logData.Message = "UpdateConsent: Error posting user consent to bank service"
		s.LoggerService.LogError(logData)
		return err
	}

	// if consentResponse.ResponseHeader.ErrorCode != "0" {
	// 	logData.Message = "UpdateConsent: Error processing user consent"
	// 	s.LoggerService.LogError(logData)
	// 	return err
	// }

	responseBytes, err := json.Marshal(consentResponse)
	if err != nil {
		logData.Message = "UpdateConsent: Error marshaling response to JSON"
		s.LoggerService.LogError(logData)
		return err
	}

	logData.Message = "UpdateConsent: User consent processed successfully"
	logData.ResponseSize = len(responseBytes)
	logData.ResponseBody = string(responseBytes)
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return nil
}
