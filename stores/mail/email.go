package mail

import (
	"bankapi/constants"
	"bankapi/models"
	"bankapi/requests"
	"bankapi/services"
	"bankapi/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"bitbucket.org/paydoh/paydoh-commons/database"
	commonSrv "bitbucket.org/paydoh/paydoh-commons/services"
)

type Store struct {
	db                 *sql.DB
	m                  *database.Document
	memory             *database.InMemory
	notificaionService *services.NotificationService
	LoggerService      *commonSrv.LoggerService
}

func NewStore(
	log *commonSrv.LoggerService,
	db *sql.DB,
	m *database.Document,
	memory *database.InMemory,
) *Store {
	notificationService := services.NewNotificationService()
	return &Store{
		db:                 db,
		m:                  m,
		memory:             memory,
		LoggerService:      log,
		notificaionService: notificationService,
	}
}

func (s *Store) GetEmailVerificationStatus(ctx context.Context, auth *models.AuthValues) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.NOMINEE,
		RequestURI: "/api/email/verification-status",
		Message:    "EmailverificationStatus log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	result, err := models.GetEmailVerificationStatus(s.db, auth.UserId)

	if err != nil {
		if errors.Is(err, constants.ErrNoDataFound) {
			logData.Message = "EmailverificationStatus: No data found of user mail"
			s.LoggerService.LogError(logData)
			return false, nil
		}
		logData.Message = "EmailverificationStatus: error in get email verification status " + err.Error()
		s.LoggerService.LogError(logData)
		return nil, err
	}
	if result.IsEmailVerified == true {
		logData.Message = "EmailverificationStatus: Email is Verified "
		logData.ResponseBody = fmt.Sprintf("", result)
		logData.EndTime = time.Now()
		s.LoggerService.LogInfo(logData)
		return true, nil
	}
	return false, nil
}

func (s *Store) SendVerificationEmail(ctx context.Context, req *requests.VerifyEmailReq, auth *models.AuthValues) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.NOMINEE,
		RequestURI: "/api/email/sendverification",
		Message:    "SendEmailVerficationLink log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	request := requests.NewSendVerificationMailReq()

	if err := request.Bind(req.EmailId, req.Name); err != nil {
		logData.Message = "SendEmailVerficationLink:Error while bind req for Send Mail"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	req.UserId = auth.UserId
	reqData, err := req.Marshal()
	if err != nil {
		logData.Message = "SendEmailVerficationLink:Error while marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	emailExpire, err := strconv.Atoi(constants.EmailVerificationExpireTime)
	if err != nil {
		emailExpire = 5
	}
	if err := s.memory.Set(fmt.Sprintf("user:emailverification:TxnId:%s", request.TxnId), reqData, time.Minute*time.Duration(emailExpire)); err != nil {
		logData.Message = "SendEmailVerficationLink:Error while saving txnId on Cache"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	result, err := s.notificaionService.SendVerificationEmail(request)
	if err != nil {
		logData.Message = "SendEmailVerficationLink:Error while Sending Verification Mail"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.Message = "SendEmailVerficationLink: Email Verification Mail Sent"
	logData.ResponseBody = fmt.Sprintf("", result)
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)
	return nil, nil
}

func (s *Store) VerifyEmail(ctx context.Context, txnId string) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.NOMINEE,
		RequestURI: "/callback/update-verify/:id",
		Message:    "MailVerificationCallback log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}
	txn, err := s.memory.Get(fmt.Sprintf("user:emailverification:TxnId%s", txnId))

	if err != nil {
		logData.Message = "MailVerificationCallback:Error while getting txn Id from Cache"
		s.LoggerService.LogError(logData)
		return nil, errors.New("transaction not found")
	}

	emailVerificationReq := requests.NewVerifyEmailReq()
	if err := emailVerificationReq.Unmarshal([]byte(txn)); err != nil {
		logData.Message = "MailVerificationCallback:Error while unmarshalling data"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	req := models.NewPersonalInformation()
	req.Email = emailVerificationReq.EmailId
	req.UserId = emailVerificationReq.UserId

	userDetail, err := models.GetPersonalInformation(s.db, emailVerificationReq.UserId)

	if err != nil {
		if errors.Is(err, constants.ErrNoDataFound) {
			if err := models.InsertEmailInPersonalInfo(s.db, req); err != nil {
				logData.Message = "MailVerificationCallback:Error while adding email in personal information " + err.Error()
				s.LoggerService.LogError(logData)
				return nil, err
			}
			logData.Message = "MailVerificationCallback:Email Verified Successfully"
			s.LoggerService.LogError(logData)
			return nil, nil
		}
		return nil, err
	}

	if userDetail.Email != req.Email {
		req.IsEmailVerified = true
		if err := models.UpdatePersonalInformation(s.db, req, req.UserId); err != nil {
			logData.Message = "MailVerificationCallback:Error while Updating email in personal information " + err.Error()
			s.LoggerService.LogError(logData)
			return nil, err
		}
		logData.Message = "EmailverificationStatus: Email is verified successfully"
		logData.ResponseBody = fmt.Sprintf("", req)
		logData.EndTime = time.Now()
		s.LoggerService.LogInfo(logData)
		return nil, nil
	}

	logData.Message = "EmailverificationStatus: Email is already Verified "
	logData.ResponseBody = fmt.Sprintf("", req)
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)
	return nil, errors.New("email is already verified")

}

func (s *Store) SendAccountInformation(ctx context.Context, userId string) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.NOMINEE,
		RequestURI: "/callback/update-verify/:id",
		Message:    "SendAccountInformationOnMail log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}
	accountData, err := models.GetUserAndAccountDetailByUserID(s.db, userId)
	if err != nil {
		logData.Message = "SendAccountInformationOnMail:Error while checking Email Verification Status " + err.Error()
		s.LoggerService.LogError(logData)
		return nil, err
	}

	req := requests.NewSendAccountInformationReq()
	if err := req.Bind(accountData.AccountNumber, accountData.FirstName+" "+accountData.LastName, accountData.Email, "Saving Account", constants.BranchName, constants.IfscCode); err != nil {
		logData.Message = "SendAccountInformationOnMail:Error while bindin data for Send Account Information " + err.Error()
		s.LoggerService.LogError(logData)
		return nil, err
	}

	result, err := s.notificaionService.SendAccountInformation(req)
	if err != nil {
		logData.Message = "SendAccountInformationOnMail:Error while sending Account Information " + err.Error()
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if err := models.AddAccountDetailSentOnMailStatus(s.db, userId); err != nil {
		logData.Message = "SendAccountInformationOnMail:Error while Adding Email Account Information Sent status" + err.Error()
		s.LoggerService.LogError(logData)
		return nil, err
	}

	return result, nil
}
