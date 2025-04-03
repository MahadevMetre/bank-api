package address

import (
	"bankapi/constants"
	"bankapi/models"
	"bankapi/requests"
	"bankapi/security"
	"bankapi/services"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"bitbucket.org/paydoh/paydoh-commons/database"
	commonSrv "bitbucket.org/paydoh/paydoh-commons/services"
)

type AddressStore struct {
	db                  *sql.DB
	memory              *database.InMemory
	LoggerService       *commonSrv.LoggerService
	bankService         *services.BankApiService
	notificationService *services.NotificationService
	auditLogSrv         services.AuditLogService
}

func NewStore(log *commonSrv.LoggerService, db *sql.DB, memory *database.InMemory, auditLogSrv services.AuditLogService) *AddressStore {
	bankService := services.NewBankApiService(log, memory)
	notification := services.NewNotificationService()
	return &AddressStore{
		db:                  db,
		memory:              memory,
		LoggerService:       log,
		bankService:         bankService,
		notificationService: notification,
		auditLogSrv:         auditLogSrv,
	}
}

func (s *AddressStore) UpdateAddress(ctx context.Context, req *requests.AddressUpdateRequest, authValue *models.AuthValues) (interface{}, error) {
	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:     constants.UPI,
		StartTime:  startTime,
		RequestURI: "/api/address-update",
		Message:    "Update Address log",
	}

	extUpdateReq, err := models.GetAddressUpdateDataByUserId(s.db, authValue.UserId)
	if err != nil {
		if err != constants.ErrNoDataFound {
			logData.Message = "UpdateAdressLog: Error while getting Address Update Data " + err.Error()
			s.LoggerService.LogError(logData)
			return nil, err
		}
		logData.Message = "UpdateAdressLog: No Existing request found for Address Update"
		s.LoggerService.LogError(logData)
	}

	if len(extUpdateReq) > 0 {
		logData.Message = "UpdateAdressLog: Error while getting Address Update Data "
		s.LoggerService.LogError(logData)
		return nil, errors.New("Address Update Request already in process")
	}

	account, err := models.GetAccountDataByUserId(s.db, authValue.UserId)
	if err != nil {
		logData.Message = "UpdateAdressLog: Error while getting Account Data" + err.Error()
		s.LoggerService.LogError(logData)
		return nil, err
	}
	if account.AccountNumber == "" {
		logData.Message = "UpdateAdressLog: Error while getting Account Data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	docUploadReq := requests.NewUploadaddressProofReq()
	if req.DocumentFront == nil {
		logData.Message = "UpdateAdressLog: Front Document Copy is Required "
		s.LoggerService.LogError(logData)
		return nil, errors.New("Front Document Copy is Required")
	}
	if err := docUploadReq.Bind(req.DocumentFront, account.CustomerId); err != nil {
		logData.Message = "UpdateAdressLog: Error while Binding request " + err.Error()
		s.LoggerService.LogError(logData)
		return nil, err
	}

	docUploadRes1, err := s.bankService.UploadAddressProof(ctx, docUploadReq)
	if err != nil {
		logData.Message = "UpdateAdressLog: Error while Uploading Aadhar Front Side Image " + err.Error()
		s.LoggerService.LogError(logData)
		return nil, errors.New("Address Upload Failed")
	}
	if req.DocumentBack == nil {
		logData.Message = "UpdateAdressLog: Back Document Copy is Required"
		s.LoggerService.LogError(logData)
		return nil, errors.New("Back Document Copy is Required")
	}
	docUploadReq2 := requests.NewUploadaddressProofReq()
	if err := docUploadReq2.Bind(req.DocumentBack, account.CustomerId); err != nil {
		logData.Message = "UpdateAdressLog: Error while Binding request " + err.Error()
		s.LoggerService.LogError(logData)
		return nil, err
	}

	docUploadRes2, err := s.bankService.UploadAddressProof(ctx, docUploadReq2)
	if err != nil {
		logData.Message = "UpdateAdressLog: Error while Uploading Aadhar Back Side Image " + err.Error()
		s.LoggerService.LogError(logData)
		return nil, errors.New("Address Upload Failed")
	}

	exists, err := models.CheckZipCodeExists(req.PinCode[:3])
	if err != nil {
		logData.Message = "UpdateAdressLog: Error checking zip code exists"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if !exists {
		logData.Message = "UpdateAdressLog: the provided pincode is not within a deliverable area"
		s.LoggerService.LogError(logData)
		return nil, errors.New("the provided pincode is not within a deliverable area")
	}

	cityData, err := models.FindCityDataByCityName(s.db, strings.ToLower(req.City))
	if err != nil {
		logData.Message = "UpdateAdressLog:Failed to find selected city data" + err.Error()
		s.LoggerService.LogError(logData)
		return nil, errors.New("Failed to find selected city data" + req.City)
	}

	if len(cityData) <= 0 {
		logData.Message = "UpdateAdressLog: selected city data is not available"
		s.LoggerService.LogError(logData)
		return nil, errors.New("Adddress city data not found 1" + req.City)
	}

	stateData, err := models.GetStateCodeByName(strings.ToLower(req.State))
	if err != nil {
		logData.Message = "UpdateAdressLog: Failed to find selected state data"
		s.LoggerService.LogError(logData)
		return nil, errors.New("Failed to find selected state data" + err.Error())
	}

	updateAdddressReq := requests.NewUpdateAddressReq()

	if err := updateAdddressReq.Bind(req, account.CustomerId, req.CityCode, cityData[0].City, stateData.StateName, docUploadRes1.Data, docUploadRes2.Data, stateData.StateCode, docUploadReq, docUploadReq2); err != nil {
		logData.Message = "UpdateAdressLog: Failed to Bind data for Address Update"
		s.LoggerService.LogError(logData)
		return nil, errors.New("Address bind Data Failed")
	}

	addUpdateRes, err := s.bankService.UpdateAddress(ctx, updateAdddressReq)
	if err != nil {
		logData.Message = "UpdateAdressLog: Error while Update Address Call to bank " + err.Error()
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if addUpdateRes.ErrorCode == constants.AddressUpdateErrorCode91 &&
		strings.EqualFold(addUpdateRes.ErrorMessage, constants.AddressUpdateErrorMessage91) {
		logData.Message = "UpdateAdressLog: Address Update request went to bank we will get shortly"
		s.LoggerService.LogError(logData)
		return nil, errors.New(constants.AddressUpdateInProgressError)
	}

	communicationAddress := requests.CommunicationAddress{}
	if err := communicationAddress.BindForAddressUpdate(req); err != nil {
		logData.Message = "UpdateAdressLog: Error while Binding Communication Address " + err.Error()
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if err := models.InsertAddressUpdate(s.db, authValue.UserId, addUpdateRes.Req_Ref_No, &communicationAddress); err != nil {
		logData.Message = "UpdateAdressLog: Error while Saving Update Address Data " + err.Error()
		s.LoggerService.LogError(logData)
	}

	if err := models.GenerateNotification(authValue.UserId, "address_update_request", "", "address_update"); err != nil {
		logData.Message = "UpdateAdressLog: Error while Generating Notification " + err.Error()
		s.LoggerService.LogError(logData)
	}

	// save audit log
	requestData, err := json.Marshal(updateAdddressReq)
	if err != nil {
		logData.Message = "UpdateAdressLog:error while marshalling"
		logData.EndTime = time.Now()
		s.LoggerService.LogError(logData)
		return addUpdateRes, nil
	}

	encryptReq, err := security.Encrypt(requestData, []byte(authValue.Key))
	if err != nil {
		logData.Message = "UpdateAdressLog:error while encrypting request"
		logData.EndTime = time.Now()
		s.LoggerService.LogError(logData)
		return addUpdateRes, nil
	}

	if err := s.auditLogSrv.Save(ctx, &services.AuditLog{
		UserID:         authValue.UserId,
		RequestURL:     "/api/address/update",
		HTTPMethod:     "POST",
		ResponseStatus: 200,
		Action:         constants.ADDRESS_UPDATE,
		RequestBody:    encryptReq,
	}); err != nil {
		logData.Message = "error while saving audit log" + err.Error()
		logData.EndTime = time.Now()
		s.LoggerService.LogInfo(logData)
	}

	return addUpdateRes, nil
}

func (s *AddressStore) GetAddressUpdateStatus(authValues *models.AuthValues) (interface{}, error) {
	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:     constants.UPI,
		StartTime:  startTime,
		RequestURI: "/api/address-update",
		Message:    "Update Address log",
	}

	extUpdateReq, err := models.GetAddressUpdateRequestListByUserId(s.db, authValues.UserId)
	if err != nil {
		if err != constants.ErrNoDataFound {
			logData.Message = "UpdateAdressLog: Error while getting Address Update Data " + err.Error()
			s.LoggerService.LogError(logData)
			return nil, err
		}
		logData.Message = "UpdateAdressLog: No Existing request found for Address Update"
		s.LoggerService.LogError(logData)
	}

	logData.Message = "UpdateAdressLog API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)
	return extUpdateReq, nil
}
