package services

import (
	"bankapi/constants"
	"bankapi/requests"
	"bankapi/responses"
	"bankapi/utils"
	"context"
	"time"

	"errors"
	"io"
	"net/http"

	"bitbucket.org/paydoh/paydoh-commons/database"
	"bitbucket.org/paydoh/paydoh-commons/httpservice"
	commonSrv "bitbucket.org/paydoh/paydoh-commons/services"
)

type DebitcardApiService struct {
	service       *httpservice.HttpService
	memory        *database.InMemory
	LoggerService *commonSrv.LoggerService
	BankService   *BankApiService
}

func NewDebitcardApiService(log *commonSrv.LoggerService, memory *database.InMemory) *DebitcardApiService {
	bankservice := NewBankApiService(log, memory)

	interceptor := &Interceptor{
		MaxRetries: 3,
		log:        log,
		bankSrv:    bankservice,
	}

	return &DebitcardApiService{
		service:       httpservice.NewHttpService(constants.KvbUatURL, interceptor),
		memory:        memory,
		LoggerService: log,
		BankService:   bankservice,
	}
}

func (s *DebitcardApiService) DebitCardVirtualGeneration(ctx context.Context, request *requests.GenerateVirtualDebitcardOutGoingReq) (*responses.GenerateDebitcardResponse, error) {
	startTime := time.Now()
	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		Latency:       time.Since(startTime).Seconds(),
		RequestMethod: "POST",
		RequestURI:    "/fintech/v2/card/debit/generate",
		Message:       "GetDebitCardGeneration log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.BankService.GenerateToken(ctx)

	if err != nil {
		logData.Message = "GetDebitCardGeneration: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := request.Marshal()

	if err != nil {
		logData.Message = "GetDebitCardGeneration: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/v2/card/debit/generate", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
	})

	if err != nil {
		logData.Message = "GetDebitCardGeneration: Error in POST request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	defer response.Body.Close()
	body, err = io.ReadAll(response.Body)

	logData.Message = "GetDebitCardGeneration: Error reading response body"
	if err != nil {
		logData.Message = "GetDebitCardGeneration: body read error"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.ResponseBody = string(body)
	if response.StatusCode != http.StatusOK {
		logData.Message = "GetDebitCardGeneration: Received non-OK response status"
		logData.ResponseSize = int(response.ContentLength)
		logData.EndTime = time.Now()
		logData.ResponseBody = string(body)
		s.LoggerService.LogError(logData)
		return nil, errors.New("unkown error")
	}

	responseData := responses.NewGenerateDebitcardResponse()
	if err := responseData.UnMarshal(body); err != nil {
		logData.Message = "GetDebitCardGeneration: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	if responseData.ErrorCode != "00" {
		logData.Message = "GetDebitCardDetails: Received error code from response\t\t" + responseData.ErrorCode
		s.LoggerService.LogError(logData)
		return nil, errors.New(responseData.ErrorMessage)
	}

	logData.Message = "GetDebitCardGeneration API call completed successfully"
	logData.ResponseSize = len(body)
	logData.EndTime = time.Now()
	logData.ResponseBody = string(body)
	s.LoggerService.LogInfo(logData)

	return responseData, nil
}

func (s *DebitcardApiService) DebitCardPhysicalGeneration(ctx context.Context, request *requests.GeneratePhysicalDebitCardOutGoingReq) (*responses.GeneratePhysicalDebitCardRes, error) {
	startTime := time.Now()
	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		Latency:       time.Since(startTime).Seconds(),
		RequestMethod: "POST",
		RequestURI:    "fintech/v2/card/debit/physical/generate",
		Message:       "GeneratePhysicalDebitCardGeneration log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.BankService.GenerateToken(ctx)
	if err != nil {
		logData.Message = "GeneratePhysicalDebitCardGeneration: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := request.Marshal()
	if err != nil {
		logData.Message = "GeneratePhysicalDebitCardGeneration: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/v2/card/debit/physical/generate", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
	})
	if err != nil {
		logData.Message = "GeneratePhysicalDebitCardGeneration: Error in POST request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	defer response.Body.Close()

	body, err = io.ReadAll(response.Body)
	if err != nil {
		logData.Message = "GeneratePhysicalDebitCardGeneration: Error reading response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.ResponseBody = string(body)
	if response.StatusCode != http.StatusOK {
		logData.Message = "GeneratePhysicalDebitCardGeneration: Received non-OK response status"
		logData.ResponseSize = int(response.ContentLength)
		logData.EndTime = time.Now()
		logData.ResponseBody = string(body)
		s.LoggerService.LogError(logData)
		return nil, errors.New("failed")
	}

	phyDebitCardRes := responses.NewGeneratePhysicalDebitCardRes()

	if err := phyDebitCardRes.UnMarshal(body); err != nil {
		logData.Message = "GeneratePhysicalDebitCardGeneration: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if phyDebitCardRes.ErrorCode != "0" && phyDebitCardRes.ErrorCode != "00" {
		logData.Message = "GeneratePhysicalDebitCardGeneration: Received error code from response"
		s.LoggerService.LogError(logData)
		return nil, errors.New(phyDebitCardRes.ErrorMessage)
	}

	logData.Message = "GeneratePhysicalDebitCardGeneration API call completed successfully"
	logData.ResponseSize = len(body)
	logData.EndTime = time.Now()
	logData.ResponseBody = string(body)
	s.LoggerService.LogInfo(logData)

	return phyDebitCardRes, nil
}

func (s *DebitcardApiService) GetDebitCardDetails(ctx context.Context, request *requests.GetDebitcardDetailRequest) (*responses.DebitcardDetailResponse, error) {
	startTime := time.Now()
	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		Latency:       time.Since(startTime).Seconds(),
		RequestMethod: "POST",
		RequestURI:    "/fintech/v2/card/virtual/fetch",
		Message:       "GetDebitCardDetails log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.BankService.GenerateToken(ctx)
	if err != nil {
		logData.Message = "GetDebitCardDetails: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := request.Marshal()
	if err != nil {
		logData.Message = "GetDebitCardDetails: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/v2/card/virtual/fetch", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
	})

	if err != nil {
		logData.Message = "GetDebitCardDetails: Error in POST request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	defer response.Body.Close()

	body, err = io.ReadAll(response.Body)

	if err != nil {
		logData.Message = "GetDebitCardDetails: Error reading response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.ResponseBody = string(body)
	s.LoggerService.LogInfo(logData)
	if response.StatusCode != http.StatusOK {
		logData.Message = "GetDebitCardDetails: Received non-OK response status"
		logData.ResponseSize = int(response.ContentLength)
		logData.EndTime = time.Now()
		logData.ResponseBody = string(body)
		s.LoggerService.LogError(logData)
		return nil, errors.New("failed")
	}

	debitcardDetailResponse := responses.NewDebitCardDetailResponseResponse()

	if err := debitcardDetailResponse.UnMarshal(body); err != nil {
		logData.Message = "GetDebitCardDetails: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if debitcardDetailResponse.ErrorCode != "0" && debitcardDetailResponse.ErrorCode != "00" {
		logData.Message = "GetDebitCardDetails: Received error code from response"
		s.LoggerService.LogError(logData)
		return nil, errors.New(debitcardDetailResponse.ErrorMessage)
	}

	logData.Message = "GetDebitCardDetails API call completed successfully"
	logData.ResponseSize = len(body)
	logData.EndTime = time.Now()
	logData.ResponseBody = string(body)
	s.LoggerService.LogInfo(logData)

	return debitcardDetailResponse, nil
}

func (s *DebitcardApiService) SetDebitCardPin(ctx context.Context, request *requests.SetDebitCardPin) (*responses.SetDebitCardPinResponse, error) {
	startTime := time.Now()
	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		Latency:       time.Since(startTime).Seconds(),
		RequestMethod: "POST",
		RequestURI:    "/fintech/v2/card/pin/set",
		Message:       "SetDetbitCardPin log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.BankService.GenerateToken(ctx)
	if err != nil {
		logData.Message = "SetDetbitCardPin: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	body, err := request.Marshal()
	if err != nil {
		logData.Message = "SetDetbitCardPin: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/v2/card/pin/set", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
	})
	if err != nil {
		logData.Message = "SetDetbitCardPin: Error in POST request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	defer response.Body.Close()

	body, err = io.ReadAll(response.Body)

	if err != nil {
		logData.Message = "SetDetbitCardPin: Error reading response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.ResponseBody = string(body)

	if response.StatusCode != http.StatusOK {
		logData.Message = "SetDetbitCardPin: Received non-OK response status"
		logData.ResponseSize = int(response.ContentLength)
		logData.EndTime = time.Now()
		logData.ResponseBody = string(body)
		s.LoggerService.LogError(logData)
		return nil, errors.New("failed")
	}

	debitcardSetPinResponse := responses.NewSetDebitCardPinResponse()

	if err := debitcardSetPinResponse.UnMarshal(body); err != nil {
		logData.Message = "SetDetbitCardPin: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if debitcardSetPinResponse.ErrorCode != "0" && debitcardSetPinResponse.ErrorCode != "00" {
		logData.Message = "SetDetbitCardPin: Received error code from response"
		s.LoggerService.LogError(logData)
		return nil, errors.New(debitcardSetPinResponse.ErrorMessage)
	}

	logData.Message = "SetDetbitCardPin API call completed successfully"
	logData.ResponseSize = len(body)
	logData.EndTime = time.Now()
	logData.ResponseBody = string(body)
	s.LoggerService.LogInfo(logData)
	return debitcardSetPinResponse, nil
}

func (s *DebitcardApiService) SendOTPForDebitCard(ctx context.Context, request *requests.SetDebitCardOTPReq) (*responses.SetDebitCardOTPResponse, error) {
	startTime := time.Now()
	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		Latency:       time.Since(startTime).Seconds(),
		RequestMethod: "POST",
		RequestURI:    "/fintech/v2/otp",
		Message:       "SendOTPForDebitCard log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.BankService.GenerateToken(ctx)
	if err != nil {
		logData.Message = "SendOTPForDebitCard: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	body, err := request.Marshal()
	if err != nil {
		logData.Message = "SendOTPForDebitCard: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/v2/otp", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
	})
	if err != nil {
		logData.Message = "SendOTPForDebitCard: Error in POST request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	defer response.Body.Close()

	body, err = io.ReadAll(response.Body)

	if err != nil {
		logData.Message = "SendOTPForDebitCard: Error reading response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.ResponseBody = string(body)

	if response.StatusCode != http.StatusOK {
		logData.Message = "SendOTPForDebitCard: Received non-OK response status"
		logData.ResponseSize = int(response.ContentLength)
		logData.EndTime = time.Now()
		logData.ResponseBody = string(body)
		s.LoggerService.LogError(logData)
		return nil, errors.New("failed")
	}

	debitCardOtpResponse := responses.NewSetDebitCardOTPResponse()

	if err := debitCardOtpResponse.UnMarshal(body); err != nil {
		logData.Message = "SendOTPForDebitCard: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if debitCardOtpResponse.ErrorCode != "0" && debitCardOtpResponse.ErrorCode != "00" {
		logData.Message = "SendOTPForDebitCard: Received error code from response"
		s.LoggerService.LogError(logData)
		return nil, errors.New(debitCardOtpResponse.ErrorMessage)
	}

	logData.Message = "SendOTPForDebitCard API call completed successfully"
	logData.ResponseSize = len(body)
	logData.EndTime = time.Now()
	logData.ResponseBody = string(body)
	s.LoggerService.LogInfo(logData)
	return debitCardOtpResponse, nil
}

func (s *DebitcardApiService) VerifyOTPForDebitCard(ctx context.Context, request *requests.SetDebitCardOTPReq) (*responses.VerifyDebitCardOTPResponse, error) {
	startTime := time.Now()
	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		Latency:       time.Since(startTime).Seconds(),
		RequestMethod: "POST",
		RequestURI:    "/fintech/v2/otp",
		Message:       "VerifyOTPForDebitCard log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.BankService.GenerateToken(ctx)
	if err != nil {
		logData.Message = "VerifyOTPForDebitCard: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	body, err := request.Marshal()
	if err != nil {
		logData.Message = "VerifyOTPForDebitCard: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/v2/otp", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
	})
	if err != nil {
		logData.Message = "VerifyOTPForDebitCard: Error in POST request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	defer response.Body.Close()

	body, err = io.ReadAll(response.Body)

	if err != nil {
		logData.Message = "VerifyOTPForDebitCard: Error reading response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.ResponseBody = string(body)

	if response.StatusCode != http.StatusOK {
		logData.Message = "VerifyOTPForDebitCard: Received non-OK response status"
		logData.ResponseSize = int(response.ContentLength)
		logData.EndTime = time.Now()
		logData.ResponseBody = string(body)
		s.LoggerService.LogError(logData)
		return nil, errors.New("failed")
	}

	verifyDebitCardOtpResponse := responses.NewVerifyDebitCardOTPResponse()

	if err := verifyDebitCardOtpResponse.UnMarshal(body); err != nil {
		logData.Message = "VerifyOTPForDebitCard: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if verifyDebitCardOtpResponse.ErrorCode != "0" && verifyDebitCardOtpResponse.ErrorCode != "00" {
		logData.Message = "VerifyOTPForDebitCard: Received error code from response"
		s.LoggerService.LogError(logData)
		return nil, errors.New(verifyDebitCardOtpResponse.ErrorMessage)
	}

	logData.Message = "VerifyOTPForDebitCard API call completed successfully"
	logData.ResponseSize = len(body)
	logData.EndTime = time.Now()
	logData.ResponseBody = string(body)
	s.LoggerService.LogInfo(logData)
	return verifyDebitCardOtpResponse, nil
}
