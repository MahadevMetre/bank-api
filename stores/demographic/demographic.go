package demographic

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

	"bankapi/constants"
	"bankapi/models"
	"bankapi/requests"
	"bankapi/responses"
	"bankapi/security"
	"bankapi/services"
	"bankapi/utils"
)

type Demographic interface {
	GetDemographicData(authValues *models.AuthValues) (interface{}, error)
}

type Store struct {
	db                  *sql.DB
	m                   *database.Document
	memory              *database.InMemory
	bankService         *services.BankApiService
	notificationService *services.NotificationService
	LoggerService       *commonSrv.LoggerService
}

func NewStore(log *commonSrv.LoggerService, db *sql.DB, m *database.Document, memory *database.InMemory) *Store {
	bankService := services.NewBankApiService(log, memory)
	notificationService := services.NewNotificationService()
	return &Store{
		db:                  db,
		m:                   m,
		memory:              memory,
		bankService:         bankService,
		LoggerService:       log,
		notificationService: notificationService,
	}
}

func (s *Store) GetDemographicData(ctx context.Context, authValues *models.AuthValues) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.DEMOGRAPHIC,
		RequestURI: "/api/demographic_module/fetch",
		Message:    "GetDemographicData log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	existingDevice, err := models.GetUserDataByUserId(s.db, authValues.UserId)

	if err != nil {
		logData.Message = "GetDemographicData: Error getting user data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	request := requests.NewOutGoingDemographicRequest()

	if err := request.Bind(existingDevice.ApplicantId, existingDevice.MobileNumber); err != nil {
		logData.Message = "GetDemographicData: Error binding demographic request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	response, err := s.bankService.GetDemographicData(ctx, request)

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

	if response == nil {
		logData.Message = "GetDemographicData: Error getting demographic data from bank service"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if response.ErrorCode != "0" && response.ErrorCode != "00" {
		logData.Message = "GetDemographicData: Error processing demographic data"
		s.LoggerService.LogError(logData)
		return nil, errors.New(response.ErrorMessage)
	}

	// update onboarding status
	if err := models.UpdateUserOnboardingStatus(constants.DEMOGRAPHIC_FETCH_STAGE, authValues.UserId); err != nil {
		logData.Message = "Create audit error while updating onboarding status"
		s.LoggerService.LogError(logData)
	}

	userPersonalInformation := responses.NewUserInformationResponse()
	if err := userPersonalInformation.Bind(existingDevice.UserId, response); err != nil {
		logData.Message = "GetDemographicData: Error binding user personal information"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	piData, _ := models.GetPersonalInformation(s.db, authValues.UserId)

	dobData := piData.DateOfBirth
	data := strings.Split(dobData, "/")
	year := data[2]

	dob := s.FormatDOB(response.Root.UIDData.Poi.Dob)

	if year != dob {
		logData.Message = "GetDemographicData: Error processing demographic data"
		s.LoggerService.LogError(logData)
		return "", errors.New("data of birth does not match with demographic data")
	}

	byteudd, err := json.Marshal(userPersonalInformation)

	if err != nil {
		logData.Message = "GetDemographicData: Error marshaling user personal information to JSON"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	encrypted, err := security.Encrypt(byteudd, []byte(authValues.Key))

	if err != nil {
		logData.Message = "GetDemographicData: Error encrypting user personal information"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "GetDemographicData: User personal information encrypted successfully"
	logData.ResponseSize = len(byteudd)
	logData.ResponseBody = string(byteudd)
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return encrypted, nil
}

// format date of birth
func (s *Store) FormatDOB(dob string) string {
	if len(dob) > 4 {
		dobData := strings.Split(dob, "-")
		return dobData[len(dobData)-1]
	}
	return dob
}
