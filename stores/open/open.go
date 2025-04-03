package open

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/paydoh/paydoh-commons/database"
	"bitbucket.org/paydoh/paydoh-commons/httpservice"
	commonSrv "bitbucket.org/paydoh/paydoh-commons/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"bankapi/constants"
	"bankapi/models"
	"bankapi/requests"
	"bankapi/responses"
	"bankapi/rpc"
	"bankapi/security"
	"bankapi/services"
	"bankapi/utils"
)

type Open interface {
	GetRelations(userId, os, osVersion, key string) (interface{}, error)
	GetStates(userId, os, osVersion, key string) (interface{}, error)

	GetCities(userId, os, osVersion, key, state string) (interface{}, error)

	UpdateShippingAddress(authValues *models.AuthValues, request *requests.AddShippingAddress) (interface{}, error)

	GetShippingAddress(authValues *models.AuthValues) (interface{}, error)

	UpdatePaymentStatus(authValues *models.AuthValues, request *requests.AddPaymentStatusRequest) (interface{}, error)
	GetPaymentStatus(authValues *models.AuthValues) (interface{}, error)

	ShippingAddressUpdate(authValues *models.AuthValues, request *requests.UpdateShippingAddress) error
}

type OpenStore struct {
	db            *sql.DB
	bankService   *services.BankApiService
	m             *database.Document
	memory        *database.InMemory
	ctx           context.Context
	client        rpc.PaymentServiceClient
	LoggerService *commonSrv.LoggerService
	service       *httpservice.HttpService
	auditLogSrv   services.AuditLogService
}

func NewOpenStore(
	log *commonSrv.LoggerService,
	db *sql.DB,
	mongo *database.Document,
	memory *database.InMemory,
	ctx context.Context,
	client rpc.PaymentServiceClient,
	auditLogSrv services.AuditLogService,
) *OpenStore {
	bankService := services.NewBankApiService(log, memory)
	service := httpservice.NewHttpService(constants.PaymentServiceURL)
	return &OpenStore{
		db:            db,
		m:             mongo,
		memory:        memory,
		client:        client,
		bankService:   bankService,
		ctx:           ctx,
		LoggerService: log,
		service:       service,
		auditLogSrv:   auditLogSrv,
	}
}

func (o *OpenStore) GetRelations(ctx context.Context, authValues *models.AuthValues) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.OPEN,
		RequestURI: "/api/open/relations",
		Message:    "GetRelations log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	relations, err := models.FindAllActiveRelations(o.db)

	if err != nil {
		logData.Message = "GetRelations: Error finding active relations"
		o.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "GetRelations: Active relations found"
	logData.EndTime = time.Now()
	o.LoggerService.LogInfo(logData)

	return relations, nil
}

func (o *OpenStore) GetStates(ctx context.Context, authValues *models.AuthValues) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.OPEN,
		RequestURI: "/api/open/state",
		Message:    "GetStates log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	states, err := models.FindAllStatesDistinct(o.db)

	if err != nil {
		logData.Message = "GetStates: Error finding all states"
		o.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "GetStates: All states found"
	logData.EndTime = time.Now()
	o.LoggerService.LogInfo(logData)

	return states, nil
}

func (o *OpenStore) GetCities(ctx context.Context, authValues *models.AuthValues, state string) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.OPEN,
		RequestURI: "/api/open//cities/:state",
		Message:    "GetCities log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	// call all five api
	// request and response

	cities, err := models.FindCityByState(o.db, state)

	if err != nil {
		logData.Message = "GetCities: Error finding cities by state"
		o.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "GetCities: Cities found by state"
	logData.EndTime = time.Now()
	o.LoggerService.LogInfo(logData)

	return cities, nil
}

func (o *OpenStore) UpdateShippingAddress(ctx context.Context, authValues *models.AuthValues, request *requests.AddShippingAddress) error {
	logData := &commonSrv.LogEntry{
		Action:     constants.OPEN,
		RequestURI: "/api/open/shipping-address",
		Message:    "UpdateShippingAddress log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	_, err := models.FindShippingAddressByUserId(o.db, authValues.UserId)

	if err != nil {
		logData.Message = "UpdateShippingAddress: Error finding shipping address by user id"
		o.LoggerService.LogError(logData)

		if err == constants.ErrNoDataFound {

			shippingaddress := models.NewShippingAddress()
			shippingaddress.Bind(request, authValues.UserId)

			if err := models.InsertShippingAddress(o.db, shippingaddress); err != nil {
				logData.Message = "UpdateShippingAddress: Error inserting shipping address"
				o.LoggerService.LogError(logData)
				return err
			}

			return nil
		}

		return err
	}

	logData.Message = "UpdateShippingAddress: Shipping address for the user already exists"
	logData.EndTime = time.Now()
	o.LoggerService.LogError(logData)

	return errors.New("shipping address for the user already exists")
}

func (s *OpenStore) GetShippingAddress(ctx context.Context, authValues *models.AuthValues) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.OPEN,
		RequestURI: "/api/open/shipping-address",
		Message:    "GetShippingAddress log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	shippingAddress, err := models.FindShippingAddressByUserId(s.db, authValues.UserId)

	if err != nil {
		logData.Message = "GetShippingAddress: Error finding shipping address by user id"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "GetShippingAddress: Shipping address found"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return shippingAddress, nil
}

func (s *OpenStore) GetReceiptID(ctx context.Context, authValues *models.AuthValues, request *requests.GetReceiptIdRequest) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.OPEN,
		RequestURI: "/api/get-receipt-id",
		Message:    "GetReceiptID log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	existingUser, err := models.FindOneDeviceByUserID(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "GetReceiptID: Error getting user data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	requestPayload := &requests.OutgoingGetReceiptIdRequest{
		GatewayId:     request.GatewayId,
		UserId:        existingUser.UserId,
		ApplicationId: existingUser.PackageId,
		Amount:        request.Amount,
		Currency:      request.Currency,
		PaymentType:   "debitcard",
		Remarks:       request.Remarks,
	}

	jsonData, err := json.Marshal(requestPayload)
	if err != nil {
		logData.Message = "GetReceiptID: Error marshalling request payload"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(jsonData)
	encryptReq, err := security.Encrypt(jsonData, []byte(constants.AesPassPhrase))
	if err != nil {
		logData.Message = "GetReceiptID: Error while encrypting  request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	requestbody := map[string]string{
		"data": encryptReq,
	}

	body, err := json.Marshal(requestbody)
	if err != nil {
		logData.Message = "GetReceiptID: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	response, err := s.service.Post("/api/payment/get-receipt-id", body, map[string]string{
		"Content-Type": "application/json",
	})
	if err != nil {
		logData.Message = "GetReceiptID: Error making HTTP request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	respData, err := utils.HandleResponse(response, nil)
	if err != nil {
		logData.Message = "GetReceiptID: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		logData.Message = "GetReceiptID: Error in receipt id response"
		s.LoggerService.LogError(logData)
		return nil, errors.New("failed to get receipt id")
	}

	logData.Message = "GetReceiptID: Receipt id successfully received"
	logData.ResponseBody = string(respData)
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	decryptResp := new(responses.EncryptRes)
	if err := json.Unmarshal(respData, decryptResp); err != nil {
		logData.Message = "GetReceiptID: Error while encrypting  request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	decryptRes, err := security.Decrypt(decryptResp.EncryptRes, []byte(constants.AesPassPhrase))
	if err != nil {
		logData.Message = "GetReceiptID: Error while encrypting  request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	reciptId := responses.NewGetReciptIDResponse()
	if err = reciptId.Unmarshal([]byte(decryptRes)); err != nil {
		return nil, err
	}

	respon := responses.NewGatewayServiceResponse()
	if err = respon.Bind(*reciptId); err != nil {
		return nil, err
	}

	return respon, nil
}

func (s *OpenStore) UpdatePaymentStatus(ctx context.Context, authValues *models.AuthValues, request *requests.AddPaymentStatusRequest) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.OPEN,
		RequestURI: "/api/open/payment-status",
		Message:    "UpdatePaymentStatus log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	statusId := 0
	if request.TransactionStatus == "success" {
		statusId = 1
	}

	requestPayload := &requests.OutgoingPaymentStatusRequest{
		ReceiptId:         request.ReceiptId,
		TransactionId:     request.TransactionId,
		StatusId:          uint32(statusId),
		TransactionStatus: request.TransactionStatus,
		TxnTimestamp:      time.Now().Format("02-01-2006 15:04:05"),
	}

	jsonData, err := json.Marshal(requestPayload)
	if err != nil {
		logData.Message = "UpdatePaymentStatus: Error marshalling request payload"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	encryptReq, err := security.Encrypt(jsonData, []byte(constants.AesPassPhrase))
	if err != nil {
		logData.Message = "GetReceiptID: Error while encrypting  request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	requestbody := map[string]string{
		"data": encryptReq,
	}

	body, err := json.Marshal(requestbody)
	if err != nil {
		logData.Message = "GetReceiptID: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.RequestBody = string(jsonData)

	response, err := s.service.Post("/api/payment/update-payment-status", body, map[string]string{
		"Content-Type": "application/json",
	})
	if err != nil {
		logData.Message = "UpdatePaymentStatus: Error making HTTP request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	respData, err := utils.HandleResponse(response, nil)
	if err != nil {
		logData.Message = "UpdatePaymentStatus: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		logData.Message = "UpdatePaymentStatus: Error in payment status response"
		s.LoggerService.LogError(logData)
		return nil, errors.New("failed to update payment status")
	}

	decryptResp := new(responses.EncryptRes)
	if err := json.Unmarshal(respData, decryptResp); err != nil {
		logData.Message = "UpdatePaymentStatus: Error while unmarshalling  request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	decryptRes, err := security.Decrypt(decryptResp.EncryptRes, []byte(constants.AesPassPhrase))
	if err != nil {
		logData.Message = "UpdatePaymentStatus: Error while encrypting  request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	s.LoggerService.LogInfo(logData)

	// update onboarding status
	if strings.ToLower(request.TransactionStatus) == "success" {
		if err := models.UpdateUserOnboardingStatus(constants.DEBIT_CARD_PAYMENT_STAGE, authValues.UserId); err != nil {
			logData.Message = "UpdateKycConsent: error while updating onboarding status"
			s.LoggerService.LogError(logData)
		}
	}
	logData.RequestBody = string(decryptRes)
	logData.Message = "UpdatePaymentStatus: Payment status successfully updated"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return nil, nil
}

func (s *OpenStore) GetPaymentStatus(ctx context.Context, authValues *models.AuthValues) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.OPEN,
		RequestURI: "/api/open/payment-status",
		Message:    "GetPaymentStatus log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	requestPayload := &requests.OutgoingGetDebitCardPaymentStatusRequest{
		UserId: authValues.UserId,
	}
	jsonData, err := json.Marshal(requestPayload)
	if err != nil {
		logData.Message = "UpdatePaymentStatus: Error marshalling request payload"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	encryptReq, err := security.Encrypt(jsonData, []byte(constants.AesPassPhrase))
	if err != nil {
		logData.Message = "GetReceiptID: Error while encrypting  request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	requestbody := map[string]string{
		"data": encryptReq,
	}

	jsonDataReq, err := json.Marshal(requestbody)
	if err != nil {
		logData.Message = "GetPaymentStatus: Error marshalling request payload"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	response, err := s.service.Post("/api/payment/get-payment-status", jsonDataReq, map[string]string{
		"Content-Type": "application/json",
	})
	if err != nil {
		logData.Message = "GetPaymentStatus: Error making HTTP request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	respData, err := utils.HandleResponse(response, nil)
	if err != nil {
		logData.Message = "GetPaymentStatus: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		logData.Message = "GetPaymentStatus: Error in payment status response"
		s.LoggerService.LogError(logData)
		return nil, errors.New("failed to get payment status")
	}

	decryptResp := new(responses.EncryptRes)
	if err := json.Unmarshal(respData, decryptResp); err != nil {
		logData.Message = "GetReceiptID: Error while encrypting  request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	decryptRes, err := security.Decrypt(decryptResp.EncryptRes, []byte(constants.AesPassPhrase))
	if err != nil {
		logData.Message = "GetReceiptID: Error while encrypting  request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	var res responses.GetDebitCardPaymentStatusResponse
	err = json.Unmarshal([]byte(decryptRes), &res)
	if err != nil {
		logData.Message = "GetPaymentStatus: Error unmarshalling response"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if len(res.Data.Data) == 0 {
		logData.Message = "GetPaymentStatus: No data found"
		s.LoggerService.LogError(logData)
		return nil, errors.New("no data found")
	}

	if len(res.Data.Data) == 1 {
		logData.Message = "GetPaymentStatus: Single data found"
		s.LoggerService.LogError(logData)
		return res.Data.Data[0], nil
	}

	logData.Message = "GetPaymentStatus: Multiple data found"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return res.Data.Data, nil
}

func (s *OpenStore) ShippingAddressUpdate(ctx context.Context, authValues *models.AuthValues, request *requests.UpdateShippingAddress) error {
	logData := &commonSrv.LogEntry{
		Action:     constants.OPEN,
		RequestURI: "/api/open/update-shipping-address",
		Message:    "ShippingAddressUpdate log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	shippingAddress, err := models.FindShippingAddressByUserId(s.db, authValues.UserId)

	if err != nil {
		logData.Message = "ShippingAddressUpdate: Error finding shipping address"
		s.LoggerService.LogError(logData)
		return err
	}

	if err := models.UpdateShippingAddressByUserId(s.db, request, shippingAddress.UserId); err != nil {
		logData.Message = "ShippingAddressUpdate: Error updating shipping address"
		s.LoggerService.LogError(logData)
		return err
	}

	logData.Message = "ShippingAddressUpdate: Shipping address updated"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return nil
}

func (s *OpenStore) SetMpin(ctx context.Context, request *requests.MpinRequest, authValues *models.AuthValues) error {
	logData := &commonSrv.LogEntry{
		Action:     constants.OPEN,
		RequestURI: "/api/open/set-mpin",
		Message:    "SetMpin log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	_, err := models.FindOneMpinByUserId(s.db, authValues.UserId)

	if err != nil {
		logData.Message = "SetMpin: Error finding mpin"
		s.LoggerService.LogError(logData)

		if errors.Is(err, constants.ErrNoDataFound) {
			logData.Message = "SetMpin: Mpin not found, inserting new mpin"
			s.LoggerService.LogError(logData)

			if err := models.InsertMpin(s.db, request, authValues.UserId); err != nil {
				logData.Message = "SetMpin: Error inserting mpin"
				s.LoggerService.LogError(logData)
				return err
			}

			// update onboarding status
			if err := models.UpdateUserOnboardingStatus(constants.M_PIN_SETUP_STAGE, authValues.UserId); err != nil {
				logData.Message = "SetMpin: error while updating onboarding status"
				s.LoggerService.LogError(logData)
			}

			return nil
		}

		return err
	}

	logData.Message = "SetMpin: Mpin already set"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return errors.New("mpin already set")
}

func (s *OpenStore) GetIfscData(ctx context.Context, bankName string) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.OPEN,
		RequestURI: "/api/open/ifsc-data/:bank",
		Message:    "GetIfscData log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	ifscData, err := models.GetIfscByBankName(s.db, bankName)

	if err != nil {
		logData.Message = "GetIfscData: Error getting ifsc data"
		s.LoggerService.LogError(logData)

		if errors.Is(err, constants.ErrNoDataFound) {
			logData.Message = "GetIfscData: Ifsc data not found"
			s.LoggerService.LogError(logData)
			return nil, constants.ErrNoDataFound
		}

		return nil, err
	}

	logData.Message = "GetIfscData: Ifsc data found"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return ifscData, nil
}

func (s *OpenStore) GetIfscDataByIfscCode(ctx context.Context, authValues *models.AuthValues, ifsccode string) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.OPEN,
		RequestURI: "/api/open/ifsc-data/banks/:ifsc",
		Message:    "GetIfscDataByIfscCode log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	ifscifo, err := models.GetIFSCData(s.db, ifsccode)
	if err != nil {
		logData.Message = "GetIfscDataByIfscCode: Error getting ifsc data"
		s.LoggerService.LogError(logData)

		if errors.Is(err, constants.ErrNoDataFound) {
			logData.Message = "GetIfscDataByIfscCode: Ifsc data not found"
			s.LoggerService.LogError(logData)
			return nil, errors.New(constants.InputErrorIfscCodeMessage)
		}

		return nil, err
	}

	logData.Message = "GetIfscDataByIfscCode: Ifsc data found"
	s.LoggerService.LogInfo(logData)

	return ifscifo, nil
}

func (s *OpenStore) GetIfscBanks(ctx context.Context) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.OPEN,
		RequestURI: "/api/open/ifsc-data/banks",
		Message:    "GetIfscBanks log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	ifscBanks, err := models.GetDistinctBankNames(s.db)

	if err != nil {
		logData.Message = "GetIfscBanks: Error getting ifsc banks"
		s.LoggerService.LogError(logData)

		if errors.Is(err, constants.ErrNoDataFound) {
			logData.Message = "GetIfscBanks: Ifsc banks not found"
			s.LoggerService.LogError(logData)
			return nil, constants.ErrNoDataFound
		}

		return nil, err
	}

	logData.Message = "GetIfscBanks: Ifsc banks found"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return ifscBanks, nil
}

func (s *OpenStore) SyncIfscApi(ctx context.Context) error {
	logData := &commonSrv.LogEntry{
		Action:    constants.OPEN,
		Message:   "SyncIfscApi log",
		UserID:    utils.GetUserIDFromContext(ctx),
		RequestID: utils.GetRequestIDFromContext(ctx),
	}

	singleResult, err := s.m.FindOne("sync_automation", "ifsc", bson.M{}, bson.M{})

	if err != nil {
		logData.Message = "SyncIfscApi: Error finding sync automation data"
		s.LoggerService.LogError(logData)

		if errors.Is(err, mongo.ErrNoDocuments) {
			logData.Message = "SyncIfscApi: No sync automation data found, creating new request"
			s.LoggerService.LogError(logData)

			request := requests.NewOutgoingIFSCSyncRequest()

			today := time.Now()
			yesterday := today.AddDate(0, 0, -1)

			if err := request.Bind(yesterday.Format("02-Jan-2006"), today.Format("02-Jan-2006")); err != nil {
				logData.Message = "SyncIfscApi: Error binding request"
				s.LoggerService.LogError(logData)
				return err
			}

			response, err := s.bankService.GetIfscData(request)

			if err != nil {
				bankErr := s.bankService.HandleBankSpecificError(err, func(errorCode string) (string, bool) {
					return constants.GetIfscSyncErrorMessage(errorCode)
				})

				if bankErr != nil {
					logData.Message = fmt.Sprintf("SyncIfscApi: Bank error encountered (ErrorCode: %s)", bankErr.ErrorCode)
					s.LoggerService.LogError(logData)
					return errors.New(bankErr.ErrorMessage)
				}

				return err
			}

			ifscData := make([]models.IFSCDto, 0)

			for _, ifsc := range response.BanksIFSCDtls {
				ifscDto := models.NewIFSCDto()
				if err := ifscDto.Bind(&ifsc); err != nil {
					logData.Message = "SyncIfscApi: Error binding ifsc data"
					s.LoggerService.LogError(logData)
					return err
				}
				ifscData = append(ifscData, *ifscDto)
			}

			if err := models.InsertManyIfscData(s.db, ifscData); err != nil {
				logData.Message = "SyncIfscApi: Error inserting ifsc data"
				s.LoggerService.LogError(logData)
				return err
			}

			lastSyncData := models.NewLastSyncIfscData()
			lastSyncData.Id = primitive.NewObjectID()
			lastSyncData.LastSynced = today.Format("02-Jan-2006")
			lastSyncData.CreatedAt = time.Now()
			lastSyncData.UpdatedAt = time.Now()

			if _, err := s.m.InsertOne("sync_automation", "ifsc", lastSyncData, false); err != nil {
				logData.Message = "SyncIfscApi: Error inserting last sync data"
				s.LoggerService.LogError(logData)
				return err
			}

			return nil
		}
		return nil
	}

	lastSyncData := models.NewLastSyncIfscData()

	if err := singleResult.Decode(lastSyncData); err != nil {
		logData.Message = "SyncIfscApi: Error decoding last sync data"
		s.LoggerService.LogError(logData)
		return err
	}

	request := requests.NewOutgoingIFSCSyncRequest()

	today := time.Now()

	if err := request.Bind(lastSyncData.LastSynced, today.Format("02-Jan-2006")); err != nil {
		logData.Message = "SyncIfscApi: Error binding request"
		s.LoggerService.LogError(logData)
		return err
	}

	response, err := s.bankService.GetIfscData(request)

	if err != nil {
		bankErr := s.bankService.HandleBankSpecificError(err, func(errorCode string) (string, bool) {
			return constants.GetIfscSyncErrorMessage(errorCode)
		})

		if bankErr != nil {
			logData.Message = fmt.Sprintf("SyncIfscApi: Bank error encountered (ErrorCode: %s)", bankErr.ErrorCode)
			s.LoggerService.LogError(logData)
			return errors.New(bankErr.ErrorMessage)
		}

		return err
	}

	ifscData := make([]models.IFSCDto, 0)

	for _, ifsc := range response.BanksIFSCDtls {
		ifscDto := models.NewIFSCDto()
		if err := ifscDto.Bind(&ifsc); err != nil {
			logData.Message = "SyncIfscApi: Error binding ifsc data"
			s.LoggerService.LogError(logData)
			return err
		}
		ifscData = append(ifscData, *ifscDto)
	}

	if err := models.InsertManyIfscData(s.db, ifscData); err != nil {
		logData.Message = "SyncIfscApi: Error inserting ifsc data"
		s.LoggerService.LogError(logData)
		return err
	}

	lastSyncData.LastSynced = today.Format("02-Jan-2006")
	lastSyncData.UpdatedAt = time.Now()

	if _, err := s.m.UpdateOne("sync_automation", "ifsc", bson.M{
		"_id": lastSyncData.Id,
	}, lastSyncData, false, false); err != nil {
		logData.Message = "SyncIfscApi: Error updating last sync data"
		s.LoggerService.LogError(logData)
		return err
	}

	logData.Message = "SyncIfscApi: Ifsc data synced successfully"
	s.LoggerService.LogInfo(logData)

	return nil
}

func (s *OpenStore) VerifyMpin(ctx context.Context, request *requests.MpinRequest, authValues *models.AuthValues) error {
	logData := &commonSrv.LogEntry{
		Action:     constants.OPEN,
		RequestURI: "/api/open/verify-mpin",
		Message:    "VerifyMpin log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	maxAttempts, err := strconv.Atoi(os.Getenv("MAX_MPIN_ATTEMPTS"))
	if err != nil {
		return fmt.Errorf("failed to parse MAX_MPIN_ATTEMPTS: %v", err)
	}

	// lockPeriodMinutes, err := strconv.Atoi(os.Getenv("MPIN_LOCK_PERIOD"))
	// if err != nil {
	// 	return fmt.Errorf("failed to parse MPIN_LOCK_PERIOD: %v", err)
	// }

	// lockPeriod := time.Duration(lockPeriodMinutes) * time.Minute

	mpinData, err := models.FindOneMpinByUserId(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "VerifyMpin: Error fetching MPIN"
		s.LoggerService.LogError(logData)
		if errors.Is(err, constants.ErrNoDataFound) {
			return errors.New("user not found")
		}
		return err
	}

	attemptData, err := models.GetMpinAttempts(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "VerifyMpin: Error fetching attempt data"
		s.LoggerService.LogError(logData)
		return err
	}

	if attemptData == nil {
		attemptData, err = models.CreateMpinAttempts(s.db, authValues.UserId)
		if err != nil {
			logData.Message = "VerifyMpin: Error creating attempt data"
			s.LoggerService.LogError(logData)
			return err
		}
	}

	if attemptData.Attempts >= maxAttempts {
		logData.Message = "VerifyMpin: Account locked due to too many incorrect attempts"
		s.LoggerService.LogError(logData)
		return errors.New("account locked. please reset your MPIN")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(mpinData.MPIN), []byte(request.Mpin)); err != nil {

		attemptData.Attempts++
		attemptData.LastAttempt = sql.NullTime{Time: time.Now(), Valid: true}

		if err := models.UpdateMpinAttempts(s.db, attemptData); err != nil {
			logData.Message = "VerifyMpin: Error updating attempt data"
			s.LoggerService.LogError(logData)
			return err
		}

		remainingAttempts := maxAttempts - attemptData.Attempts
		logData.Message = "VerifyMpin: Incorrect MPIN"
		s.LoggerService.LogError(logData)

		if remainingAttempts <= 0 {
			return errors.New("account locked. please reset your MPIN")
		}

		return fmt.Errorf("incorrect MPIN. You have %d attempts remaining", remainingAttempts)
	}

	if err := models.ResetMpinAttempts(s.db, authValues.UserId); err != nil {
		logData.Message = "VerifyMpin: Error resetting attempt data"
		s.LoggerService.LogError(logData)
		return err
	}

	logData.Message = "VerifyMpin: MPIN verified successfully"
	s.LoggerService.LogInfo(logData)

	return nil
}

func (s *OpenStore) ResetMpin(ctx context.Context, request *requests.ResetMpinRequest, authValues *models.AuthValues) error {
	logData := &commonSrv.LogEntry{
		Action:     constants.OPEN,
		RequestURI: "/api/open/reset-mpin",
		Message:    "ResetMpin log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	storedMpin, err := models.FindOneMpinByUserId(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "ResetMpin: Error finding existing MPIN"
		s.LoggerService.LogError(logData)
		return err
	}

	if request.NewMpin == storedMpin.MPIN {
		logData.Message = "ResetMpin: MPIN already reset"
		s.LoggerService.LogInfo(logData)
		return errors.New("MPIN already reset with the same value")
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedMpin.MPIN), []byte(request.ExistingMpin))
	if err != nil {
		logData.Message = "ResetMpin: Existing MPIN is incorrect"
		s.LoggerService.LogError(logData)
		return errors.New("incorrect existing MPIN")
	}

	if request.ExistingMpin == request.NewMpin {
		logData.Message = "ResetMpin: New MPIN cannot be the same as the existing MPIN"
		s.LoggerService.LogError(logData)
		return errors.New("new MPIN cannot be the same as the existing MPIN")
	}

	newMpinHash, err := bcrypt.GenerateFromPassword([]byte(request.NewMpin), bcrypt.DefaultCost)
	if err != nil {
		logData.Message = "ResetMpin: Error Encrypting new MPIN"
		s.LoggerService.LogError(logData)
		return err
	}

	if err := models.UpdateMpinAfterReset(s.db, string(newMpinHash), authValues.UserId); err != nil {
		logData.Message = "ResetMpin: Error updating MPIN in the database"
		s.LoggerService.LogError(logData)
		return err
	}

	// save audit log

	if err := s.auditLogSrv.Save(ctx, &services.AuditLog{
		UserID:         utils.GetUserIDFromContext(ctx),
		RequestURL:     "/api/open/reset-mpin",
		HTTPMethod:     "POST",
		ResponseStatus: 200,
		Action:         constants.RESET_MPIN,
	}); err != nil {
		logData.Message = "VerifyForgottenMpinResetRequest: Error saving audit log"
		s.LoggerService.LogError(logData)
	}

	logData.Message = "ResetMpin: MPIN successfully reset"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return nil
}

func (s *OpenStore) VerifyForgottenMpinResetRequest(ctx context.Context, request *requests.ForgotMpinResetRequest, authValues *models.AuthValues) error {
	logData := &commonSrv.LogEntry{
		Action:     constants.OPEN,
		RequestURI: "/api/open/verify-forgot-mpin",
		Message:    "VerifyForgottenMpinResetRequest log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	accountDetails, err := models.GetAccountDataByUserIdV2(authValues.UserId)
	if err != nil {
		logData.Message = "VerifyForgottenMpinResetRequest: Error fetching account details"
		s.LoggerService.LogError(logData)
		return err
	}

	personalInformation, err := models.GetPersonalInformation(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "VerifyForgottenMpinResetRequest: Error fetching personal information"
		s.LoggerService.LogError(logData)
		return err
	}

	if request.AccountNumber != accountDetails.AccountNumber {
		logData.Message = "VerifyForgottenMpinResetRequest: Account number does not match"
		s.LoggerService.LogError(logData)
		return errors.New("account number does not match")
	}

	if request.Email != personalInformation.Email {
		logData.Message = "VerifyForgottenMpinResetRequest: Email does not match"
		s.LoggerService.LogError(logData)
		return errors.New("email does not match")
	}

	if !strings.EqualFold(request.MotherMaidenName, accountDetails.MotherMaidenName.String) {
		logData.Message = "VerifyForgottenMpinResetRequest: Mother maiden name does not match"
		s.LoggerService.LogError(logData)
		return errors.New("mother maiden name does not match")
	}

	if accountDetails.IsAddrSameAsAadhaar {
		if request.PinCode != personalInformation.PinCode.String {
			logData.Message = "VerifyForgottenMpinResetRequest: pincode does not match personal information"
			s.LoggerService.LogError(logData)
			return errors.New("pincode does not match")
		}
	} else {

		var userCommunicationAddress requests.CommunicationAddress
		if err := json.Unmarshal([]byte(accountDetails.CommunicationAddress.String), &userCommunicationAddress); err != nil {
			logData.Message = "VerifyForgottenMpinResetRequest: Error unmarshaling communication address"
			s.LoggerService.LogError(logData)
			return errors.New("invalid communication address format")
		}

		if request.PinCode != userCommunicationAddress.PinCode {
			logData.Message = "VerifyForgottenMpinResetRequest: pincode does not match communication address"
			s.LoggerService.LogError(logData)
			return errors.New("pincode does not match")
		}
	}

	// save audit log

	if err := s.auditLogSrv.Save(ctx, &services.AuditLog{
		UserID:         utils.GetUserIDFromContext(ctx),
		RequestURL:     "/api/open/verify-forgot-mpin",
		HTTPMethod:     "POST",
		ResponseStatus: 200,
		Action:         constants.FORGOT_MPIN,
	}); err != nil {
		logData.Message = "VerifyForgottenMpinResetRequest: Error saving audit log"
		s.LoggerService.LogError(logData)
	}

	logData.Message = "VerifyForgottenMpinResetRequest: User details verified successfully for forgotten MPIN reset"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return nil
}

func (s *OpenStore) UpdateMpin(ctx context.Context, request *requests.UpdateMpinRequest, authValues *models.AuthValues) error {
	logData := &commonSrv.LogEntry{
		Action:     constants.OPEN,
		RequestURI: "/api/open/update-mpin",
		Message:    "UpdateMpin log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	storedMpin, err := models.FindOneMpinByUserId(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "ResetMpin: Error finding existing MPIN"
		s.LoggerService.LogError(logData)
		return err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedMpin.MPIN), []byte(request.Mpin)); err == nil {
		logData.Message = "UpdateMpin: MPIN already reset with the same value"
		s.LoggerService.LogInfo(logData)
		return errors.New("your new mpin matches your current one. please enter a different mpin")
	}

	newMpinHash, err := bcrypt.GenerateFromPassword([]byte(request.Mpin), bcrypt.DefaultCost)
	if err != nil {
		logData.Message = "UpdateMpin: Error Encrypting new MPIN"
		s.LoggerService.LogError(logData)
		return err
	}

	if err := models.UpdateMpinAfterForgottenResetByUserId(s.db, string(newMpinHash), authValues.UserId); err != nil {
		logData.Message = "UpdateMpin: Error updating MPIN"
		s.LoggerService.LogError(logData)
		return err
	}

	if err := models.ResetMpinAttempts(s.db, authValues.UserId); err != nil {
		logData.Message = "UpdateMpin: Error resetting MPIN attempts"
		s.LoggerService.LogError(logData)
		return err
	}

	logData.Message = "UpdateMpin: MPIN updated and attempts reset successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	// save audit log

	if err := s.auditLogSrv.Save(ctx, &services.AuditLog{
		UserID:         utils.GetUserIDFromContext(ctx),
		RequestURL:     "/api/open/update-mpin",
		HTTPMethod:     "POST",
		ResponseStatus: 200,
		Action:         constants.UPDATE_MPIN,
	}); err != nil {
		logData.Message = "UpdateMpin: Error saving audit log"
		s.LoggerService.LogError(logData)
	}

	return nil
}

func (s *OpenStore) UpdateFcmToken(ctx context.Context, request *requests.UpdateFcmToken, authValues *models.AuthValues) (interface{}, error) {
	logData := &commonSrv.LogEntry{
		Action:     constants.OPEN,
		RequestURI: "/update-fcm-token",
		Message:    "UpdateFCMToken log",
		UserID:     utils.GetUserIDFromContext(ctx),
		RequestID:  utils.GetRequestIDFromContext(ctx),
	}

	userDeviceData, err := models.FindOneDeviceByUserID(s.db, authValues.UserId)
	if err != nil {
		logData.Message = "UpdateFCMToken: Error finding user device data"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if userDeviceData.DeviceToken.Valid && userDeviceData.DeviceToken.String == request.FcmToken {
		logData.Message = "UpdateFCMToken: Device token already exists"
		s.LoggerService.LogError(logData)
		return nil, errors.New("device token already exists")
	}

	if err := models.UpdateDevice(s.db, &models.DeviceData{
		DeviceToken: sql.NullString{
			String: request.FcmToken,
			Valid:  true,
		},
	}, authValues.UserId); err != nil {
		return nil, err
	}
	return nil, nil
}
