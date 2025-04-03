package user_details

import (
	"bankapi/constants"
	"bankapi/models"
	"bankapi/requests"
	"bankapi/security"
	"bankapi/services"
	"bankapi/utils"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"bitbucket.org/paydoh/paydoh-commons/database"
	commonSrv "bitbucket.org/paydoh/paydoh-commons/services"
)

type UserDetails interface {
	GetUserDetailData(authValues *models.AuthValues) (interface{}, error)
}

type Store struct {
	db            *sql.DB
	m             *database.Document
	bankService   *services.BankApiService
	debitService  *services.DebitcardApiService
	LoggerService *commonSrv.LoggerService
}

func NewStore(log *commonSrv.LoggerService, db *sql.DB, m *database.Document, memory *database.InMemory, redis *database.InMemory) *Store {
	bankService := services.NewBankApiService(log, memory)
	debitcardService := services.NewDebitcardApiService(log, memory)

	return &Store{
		db:            db,
		m:             m,
		bankService:   bankService,
		debitService:  debitcardService,
		LoggerService: log,
	}
}
func (s *Store) GetUserDetailData(ctx context.Context, authValues *models.AuthValues) (interface{}, error) {
	startTime := time.Now()
	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		Latency:       time.Since(startTime).Seconds(),
		RequestMethod: "POST",
		RequestURI:    "/api/user/get-details",
		Message:       "GetUserDetail log",
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
	}
	personalInformation, err := models.GetUserAndAccountDetailByUserID(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "GetUserDetail: Error while getting Personal Information Detail by userId"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	request := requests.NewGetAccountDetailRequest()
	request.AccountNo = personalInformation.AccountNumber
	request.ApplicantId = personalInformation.Applicant_id

	accountDetail, err := s.bankService.GetAccountDetail(ctx, request)
	if err != nil {
		if bankErr := s.bankService.HandleBankSpecificError(err, func(errorCode string) (string, bool) {
			return constants.GetAccountDetailFetchErrorMessage(errorCode)
		}); bankErr != nil {
			return nil, errors.New(bankErr.ErrorMessage)
		}

		logData.Message = "GetUserDetail: Error while getting Account Detail from BankAPI"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	nominee, err := models.FindOneVerfiedNomineeByUserID(s.db, authValues.UserId)
	if err != nil && err.Error() != "sql: no rows in result set" {
		logData.Message = "GetUserDetail: Error while getting Nominee Detail by userId"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	nomineeName := ""
	if nominee != nil && nominee.NomName.Valid {
		nomineeName = nominee.NomName.String
	}

	userDetail := models.ResponseUserDetails{}
	if err := userDetail.Bind(personalInformation, accountDetail, nomineeName); err != nil {
		logData.Message = "GetUserDetail: Error while Binding User Detail"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	value, _ := json.Marshal(userDetail)
	encryptedvalue, err := security.Encrypt(value, []byte(authValues.Key))
	if err != nil {
		return nil, err
	}
	logData.Message = "GetUserDetail: Response encrypted successfully"
	logData.ResponseSize = len(value)
	logData.EndTime = time.Now()
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(value)
	s.LoggerService.LogInfo(logData)
	return encryptedvalue, err

}
