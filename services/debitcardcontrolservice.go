package services

import (
	"bankapi/constants"
	"bankapi/requests"
	"bankapi/responses"
	"bankapi/utils"
	"context"
	"encoding/json"
	"time"

	"errors"
	"net/http"

	"bitbucket.org/paydoh/paydoh-commons/database"
	"bitbucket.org/paydoh/paydoh-commons/httpservice"
	commonSrv "bitbucket.org/paydoh/paydoh-commons/services"
)

type DebitcardControlApiService struct {
	service       *httpservice.HttpService
	memory        *database.InMemory
	LoggerService *commonSrv.LoggerService
	BankService   *BankApiService
}

func NewDebitcardControlApiService(log *commonSrv.LoggerService, memory *database.InMemory) *DebitcardControlApiService {
	bankservice := NewBankApiService(log, memory)

	interceptor := &CardControlInterceptor{
		MaxRetries: 3,
		Transport:  http.DefaultTransport,
		log:        log,
		bankSrv:    bankservice,
	}

	return &DebitcardControlApiService{
		service:       httpservice.NewHttpService(constants.KvbUatURL, interceptor),
		memory:        memory,
		LoggerService: log,
		BankService:   bankservice,
	}
}

func (s *DebitcardControlApiService) KeyFetch(ctx context.Context, req interface{}, transactionID string) (*responses.ExternalKeyFetch, error) {
	startTime := time.Now()
	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		Latency:       time.Since(startTime).Seconds(),
		RequestMethod: "POST",
		RequestURI:    "/fintech/card/key/fetch",
		Message:       "GetPublicKey log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.BankService.GenerateToken(ctx)

	if err != nil {
		logData.Message = "GetPublicKey: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := json.Marshal(req)
	if err != nil {
		logData.Message = "GetPublicKey: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.RequestBody = string(body)

	response, err := s.service.Post("/fintech/card/key/fetch", body, map[string]string{
		"Content-Type":     "application/json",
		"Authorization":    "Bearer " + token.AccessToken,
		"X-Transaction-ID": transactionID,
	})

	if err != nil {
		logData.Message = "GetPublicKey: Error in POST request"
		logData.ResponseBody = string(err.Error())
		s.LoggerService.LogError(logData)
		return nil, err
	}

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "GetPublicKey: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		logData.Message = "GetPublicKey: Received non-OK response status"
		logData.ResponseSize = int(response.ContentLength)
		logData.EndTime = time.Now()
		logData.ResponseBody = string(body)
		s.LoggerService.LogError(logData)
		return nil, errors.New("unknown error")
	}

	logData.ResponseBody = string(respData)
	KeyFetchResponse := responses.NewKeyFetchResponse()
	if err := KeyFetchResponse.UnMarshal(respData); err != nil {
		logData.Message = "GetPublicKey: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if KeyFetchResponse.RC != "00" {
		logData.ResponseSize = len(respData)
		logData.ResponseBody = string(respData)
		logData.EndTime = time.Now()
		logData.Latency = time.Since(startTime).Seconds()
		s.LoggerService.LogError(logData)
		return nil, errors.New(KeyFetchResponse.Desc)
	}

	logData.Message = "GetPublicKey: PublicKey generated and cached successfully"
	logData.ResponseSize = len(respData)
	logData.ResponseBody = string(respData)
	logData.EndTime = time.Now()
	logData.Latency = time.Since(startTime).Seconds()
	s.LoggerService.LogInfo(logData)

	return KeyFetchResponse, nil

}

func (s *DebitcardControlApiService) Login(ctx context.Context, req *requests.LoginRequest, transactionID string) (*responses.LoginResponse, error) {
	startTime := time.Now()
	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		Latency:       time.Since(startTime).Seconds(),
		RequestMethod: "POST",
		RequestURI:    "/fintech/card/login",
		Message:       "CardLogin log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.BankService.GenerateToken(ctx)
	if err != nil {
		logData.Message = "CardLogin: Error Gnerating Token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := json.Marshal(req)
	if err != nil {
		logData.Message = "CardLogin: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/card/login", body, map[string]string{
		"Content-Type":     "application/json",
		"Authorization":    "Bearer " + token.AccessToken,
		"X-Transaction-ID": transactionID,
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "CardLogin: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseBody = string(body)
	s.LoggerService.LogInfo(logData)
	if response.StatusCode != http.StatusOK {
		logData.Message = "CardLogin: Received non-OK response status"
		logData.ResponseSize = int(response.ContentLength)
		logData.EndTime = time.Now()
		logData.ResponseBody = string(body)
		s.LoggerService.LogError(logData)
		return nil, errors.New("unknown error")
	}

	logData.ResponseBody = string(respData)
	loginResponse := responses.NewLoginResponse()
	err = json.Unmarshal(respData, &loginResponse)
	if err != nil {
		logData.Message = "CardLogin:Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if loginResponse.RC != "00" {
		logData.Message = "CardLogin: Received error code from response " + loginResponse.RC

		logData.ResponseSize = len(respData)
		logData.ResponseBody = string(respData)
		logData.EndTime = time.Now()
		logData.Latency = time.Since(startTime).Seconds()
		s.LoggerService.LogError(logData)
		return nil, errors.New(loginResponse.Desc)
	}

	logData.Message = "CardLogin:  API call completed successfully"
	logData.ResponseSize = len(respData)
	logData.EndTime = time.Now()
	logData.ResponseBody = string(respData)
	s.LoggerService.LogInfo(logData)

	return loginResponse, nil

}

func (s *DebitcardControlApiService) AddCard(ctx context.Context, req *requests.AddCardRequest, transactionID string) (*responses.AddCardResponse, error) {
	startTime := time.Now()
	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		Latency:       time.Since(startTime).Seconds(),
		RequestMethod: "POST",
		RequestURI:    "/fintech/card/add",
		Message:       "AddCardAPI log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.BankService.GenerateToken(ctx)

	if err != nil {
		logData.Message = "AddCardAPI:Error in creating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := json.Marshal(req)
	if err != nil {
		logData.Message = "AddCardAPI: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.RequestBody = string(body)

	response, err := s.service.Post("/fintech/card/add", body, map[string]string{
		"Content-Type":     "application/json",
		"Authorization":    "Bearer " + token.AccessToken,
		"X-Transaction-ID": transactionID,
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "AddCardAPI: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseBody = string(body)
	s.LoggerService.LogInfo(logData)
	if response.StatusCode != http.StatusOK {
		logData.Message = "AddCardAPI: Received non-OK response status"
		logData.ResponseSize = int(response.ContentLength)
		logData.EndTime = time.Now()
		logData.ResponseBody = string(body)
		s.LoggerService.LogError(logData)
		return nil, errors.New("unknown error")
	}

	AddCardResponse := responses.NewAddCardResponse()

	if err := AddCardResponse.UnMarshal(respData); err != nil {
		logData.Message = "AddCardAPI:Error in Unmarshal"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if AddCardResponse.RC != "00" {
		logData.Message = "AddCardAPI: Received error code from response " + AddCardResponse.RC
		logData.ResponseSize = len(respData)
		logData.ResponseBody = string(respData)
		logData.EndTime = time.Now()
		logData.Latency = time.Since(startTime).Seconds()
		s.LoggerService.LogError(logData)
		return nil, errors.New(AddCardResponse.Desc)
	}

	logData.Message = "AddCardAPI: API call completed successfully"
	logData.ResponseSize = len(respData)
	logData.EndTime = time.Now()
	logData.ResponseBody = string(respData)
	s.LoggerService.LogInfo(logData)

	return AddCardResponse, nil
}

func (s *DebitcardControlApiService) ListCard(ctx context.Context, req *requests.ListCardRequest, transactionID string) ([]responses.ListCardResponse, error) {
	startTime := time.Now()
	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		Latency:       time.Since(startTime).Seconds(),
		RequestMethod: "POST",
		RequestURI:    "/fintech/card/list",
		Message:       "ListCardAPI log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.BankService.GenerateToken(ctx)

	if err != nil {
		logData.Message = "ListCardAPI :Error in creating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := json.Marshal(req)
	if err != nil {
		logData.Message = "ListCardAPI: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/card/list", body, map[string]string{
		"Content-Type":     "application/json",
		"Authorization":    "Bearer " + token.AccessToken,
		"X-Transaction-ID": transactionID,
	})

	if err != nil {
		logData.Message = "ListCardAPI :Error in POSt request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "ListCardAPI: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseBody = string(body)
	s.LoggerService.LogInfo(logData)
	if response.StatusCode != http.StatusOK {
		logData.Message = "ListCardAPI: Received non-OK response status"
		logData.ResponseSize = int(response.ContentLength)
		logData.EndTime = time.Now()
		logData.ResponseBody = string(body)
		s.LoggerService.LogError(logData)
		return nil, errors.New("unknown error")
	}

	ListCardResponse := responses.NewListCardResponse()

	err = json.Unmarshal(respData, &ListCardResponse)
	if err != nil {
		logData.Message = "ListCardAPI:Error in Unmarshal"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if ListCardResponse[0].RC != "00" {
		logData.Message = "ListCardAPI: Received error code from response\t\t" + ListCardResponse[0].RC
		logData.ResponseSize = len(respData)
		logData.ResponseBody = string(respData)
		logData.EndTime = time.Now()
		logData.Latency = time.Since(startTime).Seconds()
		s.LoggerService.LogError(logData)
		return nil, errors.New(ListCardResponse[0].Desc)
	}

	logData.Message = "ListCardAPI: API call completed successfully"
	logData.ResponseSize = len(respData)
	logData.EndTime = time.Now()
	logData.ResponseBody = string(respData)
	s.LoggerService.LogInfo(logData)

	return ListCardResponse, nil
}

func (s *DebitcardControlApiService) ListCardControl(ctx context.Context, req *requests.ListCardControlRequest, transactionID string) ([]responses.ListCardControlResponse, error) {
	startTime := time.Now()
	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		Latency:       time.Since(startTime).Seconds(),
		RequestMethod: "POST",
		RequestURI:    "/fintech/card/view",
		Message:       "ListCardView log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.BankService.GenerateToken(ctx)

	if err != nil {
		logData.Message = "ListCardView:Error in creating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := json.Marshal(req)
	if err != nil {
		logData.Message = "ListCardView: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/card/view", body, map[string]string{
		"Content-Type":     "application/json",
		"Authorization":    "Bearer " + token.AccessToken,
		"X-Transaction-ID": transactionID,
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "ListCardView: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		logData.Message = "ListCardView: Received non-OK response status"
		logData.ResponseSize = int(response.ContentLength)
		logData.EndTime = time.Now()
		logData.ResponseBody = string(body)
		s.LoggerService.LogError(logData)
		return nil, errors.New("unknown error")
	}

	ListCardControllResponse := responses.NewListCardControllResponse()

	err = json.Unmarshal(respData, &ListCardControllResponse)
	if err != nil {
		logData.Message = "ListCardView:Error in Unmarshal"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if ListCardControllResponse[0].RC != "00" {
		logData.Message = "ListCardView: Received error code from response\t\t" + ListCardControllResponse[0].RC
		logData.ResponseSize = len(respData)
		logData.ResponseBody = string(respData)
		logData.EndTime = time.Now()
		logData.Latency = time.Since(startTime).Seconds()
		s.LoggerService.LogError(logData)
		return nil, errors.New(ListCardControllResponse[0].Desc)
	}

	logData.Message = "ListCardView: API call completed successfully"
	logData.ResponseSize = len(respData)
	logData.EndTime = time.Now()
	logData.ResponseBody = string(respData)
	s.LoggerService.LogInfo(logData)

	return ListCardControllResponse, nil
}

func (s *DebitcardControlApiService) FetchTransaction(ctx context.Context, req *requests.FetchTransactionRequest, transactionID string) (*responses.FetchTransactionResponse, error) {
	startTime := time.Now()
	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		Latency:       time.Since(startTime).Seconds(),
		RequestMethod: "POST",
		RequestURI:    "/fintech/card/transaction/fetch",
		Message:       "CardTransactionFetch log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.BankService.GenerateToken(ctx)

	if err != nil {
		logData.Message = "CardTransactionFetch :Error in ceating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := json.Marshal(req)
	if err != nil {
		logData.Message = "CardTransactionFetch: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/card/transaction/fetch", body, map[string]string{
		"Content-Type":     "application/json",
		"Authorization":    "Bearer " + token.AccessToken,
		"X-Transaction-ID": transactionID,
	})

	if err != nil {
		logData.Message = "CardTransactionFetch: Error in POST request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	respData, err := utils.HandleResponse(response, err)
	logData.ResponseBody = string(respData)
	s.LoggerService.LogError(logData)
	if err != nil {
		logData.Message = "CardTransactionFetch: Error in POST request Handle"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		logData.Message = "CardTransactionFetch: Received non-OK response status"
		logData.ResponseSize = int(response.ContentLength)
		logData.EndTime = time.Now()
		logData.ResponseBody = string(body)
		s.LoggerService.LogError(logData)
		return nil, errors.New("unknown error")
	}

	FetchTransactionResponse := responses.NewFetchTransactionResponse()

	if err := FetchTransactionResponse.UnMarshal(respData); err != nil {
		logData.Message = "CardTransactionFetch :Error in Unmarshal"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if FetchTransactionResponse.RC != "00" {
		logData.Message = "CardTransactionFetch: Received error code from response\t\t" + FetchTransactionResponse.RC
		logData.ResponseSize = len(respData)
		logData.ResponseBody = string(respData)
		logData.EndTime = time.Now()
		logData.Latency = time.Since(startTime).Seconds()
		s.LoggerService.LogError(logData)
		return nil, errors.New(FetchTransactionResponse.Desc)
	}

	logData.Message = "CardTransactionFetch: API call completed successfully"
	logData.ResponseSize = len(respData)
	logData.EndTime = time.Now()
	logData.ResponseBody = string(respData)
	s.LoggerService.LogInfo(logData)

	return FetchTransactionResponse, nil
}

func (s *DebitcardControlApiService) EditTransaction(ctx context.Context, req *requests.EditTransactionRequest, transactionID string) (*responses.EditTransactionResponse, error) {
	startTime := time.Now()
	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		Latency:       time.Since(startTime).Seconds(),
		RequestMethod: "POST",
		RequestURI:    "/fintech/card/transaction/edit",
		Message:       "EditCardTransaction log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.BankService.GenerateToken(ctx)

	if err != nil {
		logData.Message = "EditCardTransaction:Error in creating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := json.Marshal(req)
	if err != nil {
		logData.Message = "EditCardTransaction: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/card/transaction/edit", body, map[string]string{
		"Content-Type":     "application/json",
		"Authorization":    "Bearer " + token.AccessToken,
		"X-Transaction-ID": transactionID,
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "EditCardTransaction: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseBody = string(body)
	s.LoggerService.LogInfo(logData)
	if response.StatusCode != http.StatusOK {
		logData.Message = "EditCardTransaction: Received non-OK response status"
		logData.ResponseSize = int(response.ContentLength)
		logData.EndTime = time.Now()
		logData.ResponseBody = string(body)
		s.LoggerService.LogError(logData)
		return nil, errors.New("unknown error")
	}

	EditTransactionResponse := responses.NewEditTransactionResponse()

	if err := EditTransactionResponse.UnMarshal(respData); err != nil {
		logData.Message = "EditCardTransaction:Error in Unmarshal"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	if EditTransactionResponse.RC != "00" {
		logData.Message = "EditCardTransaction: Received error code from response\t\t" + EditTransactionResponse.RC
		logData.ResponseSize = len(respData)
		logData.ResponseBody = string(respData)
		logData.EndTime = time.Now()
		logData.Latency = time.Since(startTime).Seconds()
		s.LoggerService.LogError(logData)
		return nil, errors.New(EditTransactionResponse.Desc)
	}

	logData.Message = "EditCardTransaction: API call completed successfully"
	logData.ResponseSize = len(respData)
	logData.EndTime = time.Now()
	logData.ResponseBody = string(respData)
	s.LoggerService.LogInfo(logData)

	return EditTransactionResponse, nil
}

func (s *DebitcardControlApiService) CardBlock(ctx context.Context, req *requests.CardBlockRequest, transactionID string) (*responses.CardBlockResponse, error) {
	startTime := time.Now()
	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		Latency:       time.Since(startTime).Seconds(),
		RequestMethod: "POST",
		RequestURI:    "/fintech/card/control/edit",
		Message:       "EditCardControl log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.BankService.GenerateToken(ctx)

	if err != nil {
		logData.Message = "EditCardControl:Error in creating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := json.Marshal(req)
	if err != nil {
		logData.Message = "EditCardControl: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/card/control/edit", body, map[string]string{
		"Content-Type":     "application/json",
		"Authorization":    "Bearer " + token.AccessToken,
		"X-Transaction-ID": transactionID,
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "EditCardControl: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseBody = string(body)
	s.LoggerService.LogInfo(logData)
	if response.StatusCode != http.StatusOK {
		logData.Message = "EditCardControl: Received non-OK response status"
		logData.ResponseSize = int(response.ContentLength)
		logData.EndTime = time.Now()
		logData.ResponseBody = string(body)
		s.LoggerService.LogError(logData)
		return nil, errors.New("unknown error")
	}

	CardBlockResponse := responses.NewCardBlockResponse()

	if err := CardBlockResponse.UnMarshal(respData); err != nil {
		logData.Message = "EditCardControl :Error in Unmarshal"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	if CardBlockResponse.RC != "00" {
		logData.Message = "EditCardControl: Received error code from response\t\t" + CardBlockResponse.RC
		logData.ResponseSize = len(respData)
		logData.ResponseBody = string(respData)
		logData.EndTime = time.Now()
		logData.Latency = time.Since(startTime).Seconds()
		s.LoggerService.LogError(logData)
		return nil, errors.New(CardBlockResponse.Desc)
	}

	logData.Message = "EditCardControl: API call completed successfully"
	logData.ResponseSize = len(respData)
	logData.EndTime = time.Now()
	logData.ResponseBody = string(respData)
	s.LoggerService.LogInfo(logData)

	return CardBlockResponse, nil
}
