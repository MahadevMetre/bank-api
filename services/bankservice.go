package services

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"bitbucket.org/paydoh/paydoh-commons/database"
	"bitbucket.org/paydoh/paydoh-commons/httpservice"
	"bitbucket.org/paydoh/paydoh-commons/pkg/json"
	commonSrv "bitbucket.org/paydoh/paydoh-commons/services"

	"bankapi/constants"
	"bankapi/requests"
	"bankapi/responses"
	"bankapi/utils"
)

type BankApiService struct {
	service       *httpservice.HttpService
	Memory        *database.InMemory
	LoggerService *commonSrv.LoggerService
}

func NewBankApiService(log *commonSrv.LoggerService, memory *database.InMemory) *BankApiService {
	bankApiService := &BankApiService{
		Memory:        memory,
		LoggerService: log,
	}
	interceptor := &Interceptor{
		MaxRetries: 3,
		log:        log,
		bankSrv:    bankApiService,
	}

	bankApiService.service = httpservice.NewHttpService(constants.KvbUatURL, interceptor)
	return bankApiService
}

func (b *BankApiService) ExtractBankError(err error) *responses.BankErrorResponse {
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		var customErr *BankErrorWithRetry
		if errors.As(urlErr.Err, &customErr) {
			return customErr.BankError
		}
	}

	// Scenario 3- Beneficiary details sent to KVB, but not received the response. (Back Ofice TimeOut)
	// Senerio 12- OTP Recived and OTP verification request  initiated and receive fail  response Technical Error
	// Senerio 12- OTP Recived and OTP verification request  initiated and receive fail  response back Office Time Out
	// Check if the error is directly a BankErrorResponse
	var bankErr *responses.BankErrorResponse
	if errors.As(err, &bankErr) {
		return bankErr
	}

	return nil
}

type ErrorMessageMapper func(errorCode string) (string, bool)

func (b *BankApiService) HandleBankSpecificError(err error, mapper ErrorMessageMapper) *responses.BankErrorResponse {
	if bankErr := b.ExtractBankError(err); bankErr != nil {
		if msg, exists := mapper(bankErr.ErrorCode); exists {
			bankErr.ErrorMessage = msg
			return bankErr
		}
		return bankErr
	}
	return nil
}

func (s *BankApiService) GenerateToken(ctx context.Context) (*responses.TokenResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/nbfc/v1/oauth/cc/accesstoken",
		Message:       "GenerateToken log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	tokenResponseString, _ := s.Memory.Get("kvb_auth_token")

	logData.Message = "GenerateToken: Token retrieved from cache"
	logData.ResponseSize = len(tokenResponseString)
	logData.EndTime = time.Now()

	if tokenResponseString == "" || tokenResponseString == "null" {
		username := constants.KvbUserName
		password := constants.KvbPassword

		logData.Message = "GenerateToken: Retrieving new token"

		formFields := map[string]interface{}{
			"grant_type": "client_credentials",
		}

		authorizationDetails := httpservice.NewAuthDetails(httpservice.BasicAuth, username, password, "")

		response, err := s.service.PostFormNoFile("/nbfc/v1/oauth/cc/accesstoken", formFields, authorizationDetails)
		if err != nil {
			logData.Message = "GenerateToken: Error making POST request"
			s.LoggerService.LogError(logData)
			return nil, err
		}

		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			logData.Message = "GenerateToken: Received non-OK response status"
			logData.ResponseSize = int(response.ContentLength)
			logData.EndTime = time.Now()
			s.LoggerService.LogError(logData)

			if response.StatusCode == http.StatusUnauthorized {
				return nil, errors.New("unauthorized")
			}

			if response.StatusCode == http.StatusForbidden {
				return nil, errors.New("forbidden")
			}

			if response.StatusCode == http.StatusBadRequest {
				errorTokenResponse := responses.NewTokenErrorResponse()
				body, err := io.ReadAll(response.Body)
				if err != nil {
					logData.Message = "GenerateToken: Error reading error token response body"
					s.LoggerService.LogError(logData)
					return nil, err
				}

				if err := errorTokenResponse.Unmarshal(body); err != nil {
					logData.Message = "GenerateToken: Error unmarshaling error token response"
					s.LoggerService.LogError(logData)
					return nil, err
				}

				logData.Message = "GenerateToken: Received error from token endpoint"
				s.LoggerService.LogError(logData)

				return nil, errors.New(errorTokenResponse.ErrorDescription)
			}

			logData.Message = "GenerateToken: Unknown error occurred"
			s.LoggerService.LogError(logData)

			return nil, errors.New("unknown error")
		}

		body, err := io.ReadAll(response.Body)
		if err != nil {
			logData.Message = "GenerateToken: Error reading response body"
			s.LoggerService.LogError(logData)
			return nil, err
		}

		token := responses.NewTokenResponse()
		if err := token.Unmarshal(body); err != nil {
			logData.Message = "GenerateToken: Error unmarshaling token response"
			s.LoggerService.LogError(logData)
			return nil, err
		}

		tokenData, err := token.Marshal()
		if err != nil {
			logData.Message = "GenerateToken: Error marshaling token data"
			s.LoggerService.LogError(logData)
			return nil, err
		}

		expiry, err := strconv.Atoi(token.ExpiresIn)
		if err != nil {
			logData.Message = "GenerateToken: Error converting token expiry"
			s.LoggerService.LogError(logData)
			return nil, err
		}

		if err := s.Memory.Set("kvb_auth_token", string(tokenData), time.Duration(expiry)*time.Second); err != nil {
			logData.Message = "GenerateToken: Error caching token"
			s.LoggerService.LogError(logData)
			return nil, err
		}

		logData.Message = "GenerateToken: Token generated and cached successfully"
		logData.ResponseSize = len(body)
		logData.ResponseBody = string(body)
		logData.EndTime = time.Now()
		logData.Latency = time.Since(startTime).Seconds()
		s.LoggerService.LogInfo(logData)

		return token, nil
	}

	logData.Message = "GenerateToken: Token retrieved from cache and unmarshaled successfully"
	logData.ResponseSize = len(tokenResponseString)
	logData.EndTime = time.Now()
	logData.Latency = time.Since(startTime).Seconds()
	s.LoggerService.LogInfo(logData)

	tokenResponse := responses.NewTokenResponse()
	if err := tokenResponse.Unmarshal([]byte(tokenResponseString)); err != nil {
		logData.Message = "GenerateToken: Error unmarshaling cached token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	return tokenResponse, nil
}

func (s *BankApiService) VerifySim(ctx context.Context, request *requests.OutgoingSimVerificationRequest) (*responses.SimVerificationResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/sms-verification",
		Message:       "VerifySim log",
		RequestHost:   s.service.Host,
		AppVersion:    utils.GetAppVersionFromContext(ctx),
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "VerifySim: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := json.Marshal(request)
	if err != nil {
		logData.Message = "VerifySim: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/sms-verification", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "VerifySim: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	simVerificationResponse := responses.NewSimVerificationResponse()
	if err := simVerificationResponse.Unmarshal(respData); err != nil {
		logData.Message = "VerifySim: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "VerifySim API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return simVerificationResponse, nil
}

func (s *BankApiService) SmsVerification(ctx context.Context, request *requests.OutgoingSmsVerificationRequest) (*responses.SimVerificationResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/sms-verification",
		Message:       "SmsVerification log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "SmsVerification: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := json.Marshal(request)
	if err != nil {
		logData.Message = "SmsVerification: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/sms-verification", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "SmsVerification: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	simVerificationResponse := responses.NewSimVerificationResponse()
	if err := simVerificationResponse.Unmarshal(respData); err != nil {
		logData.Message = "SmsVerification: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "SmsVerification API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return simVerificationResponse, nil
}

func (s *BankApiService) VcipInvoke(ctx context.Context, request *requests.OutGoingVcipInvokeRequest) (*responses.KycInvokeResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/vcip/invoke",
		Message:       "VcipInvoke log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "VcipInvoke: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := json.Marshal(request)
	if err != nil {
		logData.Message = "VcipInvoke: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)

	response, err := s.service.Post("/fintech/vcip/invoke", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "VcipInvoke: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	kycInvokeResponse := responses.NewKycInvokeResponse()
	if err := kycInvokeResponse.Unmarshal(respData); err != nil {
		logData.Message = "VcipInvoke: Error unmarshaling response body"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "VcipInvoke API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return kycInvokeResponse, nil
}

func (s *BankApiService) GetDemographicData(ctx context.Context, request *requests.OutgoingDemographicRequest) (*responses.DemographicResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/demographic/fetch",
		Message:       "GetDemographicData log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "GetDemographicData: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := json.Marshal(request)
	if err != nil {
		logData.Message = "GetDemographicData: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/demographic/fetch", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "GetDemographicData: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	demographicResponse := responses.NewDemographicResponse()
	if err := demographicResponse.Unmarshal(respData); err != nil {
		logData.Message = "GetDemographicData: Error unmarshaling response body"
		logData.ResponseBody = string(body)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "GetDemographicData API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return demographicResponse, nil
}

func (s *BankApiService) CreatebankAccount(ctx context.Context, request *requests.OutgoingCreateBankAccountRequest) (*responses.ImmediateCreateBankResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/account-create",
		Message:       "CreatebankAccount log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "CreatebankAccount: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := json.Marshal(request)
	if err != nil {
		logData.Message = "CreatebankAccount: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/account-create", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "CreatebankAccount: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	immediateResponse := responses.NewImmediateCreateBankResponse()
	if err := immediateResponse.Unmarshal(respData); err != nil {
		logData.Message = "CreatebankAccount: Error unmarshaling response body"
		logData.ResponseBody = string(body)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if immediateResponse.ErrorCode != "0" {
		logData.Message = "CreatebankAccount: Received error from API"
		logData.ResponseSize = len(respData)
		logData.EndTime = time.Now()
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, errors.New(immediateResponse.ErrorMessage)
	}

	logData.Message = "CreatebankAccount API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return immediateResponse, nil
}

func (s *BankApiService) CreateAddNominee(ctx context.Context, request *requests.OutgoingAddNomineeRequest) (*responses.OtpGenerationResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/nominee-registration",
		Message:       "CreateAddNominee log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "CreateAddNominee: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := json.Marshal(request)
	if err != nil {
		logData.Message = "CreateAddNominee: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/nominee-registration", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "CreateAddNominee: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	addNomineeResponse := responses.NewOtpGenartionResponse()
	if err := addNomineeResponse.Unmarshal(respData); err != nil {
		logData.Message = "CreateAddNominee: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if addNomineeResponse.ErrorCode != "0" && addNomineeResponse.ErrorCode != "00" {
		logData.Message = "CreateAddNominee: Received error from API"
		logData.ResponseSize = len(respData)
		logData.EndTime = time.Now()
		s.LoggerService.LogError(logData)
		return nil, errors.New(addNomineeResponse.ErrorMessage)
	}

	logData.Message = "CreateAddNominee API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return addNomineeResponse, nil
}

func (s *BankApiService) VerifyNomineeOTP(ctx context.Context, request *requests.OutgoingVerifyNomineeOTP) (*responses.OtpAuthenticationResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/nominee-registration",
		Message:       "VerifyNomineeOTP log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "VerifyNomineeOTP: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := json.Marshal(request)
	if err != nil {
		logData.Message = "VerifyNomineeOTP: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/nominee-registration", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "VerifyNomineeOTP: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	otpAuthenticationResponse := responses.NewOtpAuthenticationResponse()
	if err := otpAuthenticationResponse.Unmarshal(respData); err != nil {
		logData.Message = "VerifyNomineeOTP: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		logData.ResponseBody = string(body)
		return nil, err
	}

	logData.Message = "VerifyNomineeOTP API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return otpAuthenticationResponse, nil
}

func (s *BankApiService) FetchNominee(ctx context.Context, request *requests.OutgoingFetchNomineeRequest) (*responses.FetchNomineeResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/nominee/fetch",
		Message:       "FetchNominee log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "FetchNominee: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := json.Marshal(request)
	if err != nil {
		logData.Message = "FetchNominee: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/nominee/fetch", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "FetchNominee: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	fetchNomineeResponse := responses.NewFetchNomineeResponse()
	if err := fetchNomineeResponse.Unmarshal(respData); err != nil {
		logData.Message = "FetchNominee: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if fetchNomineeResponse.ErrorCode != "0" && fetchNomineeResponse.ErrorCode != "00" {
		logData.Message = "FetchNominee: Received error code from response"
		s.LoggerService.LogError(logData)
		return nil, errors.New(fetchNomineeResponse.ErrorMessage)
	}

	logData.Message = "FetchNominee API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return fetchNomineeResponse, nil
}

func (s *BankApiService) GetIfscData(request *requests.OutgoingSyncRequest) (*responses.IfscDataResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/ifsc-sync",
		Message:       "GetIfscData log",
		RequestHost:   s.service.Host,
	}

	token, err := s.GenerateToken(context.Background())
	if err != nil {
		logData.Message = "GetIfscData: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := json.Marshal(request)
	if err != nil {
		logData.Message = "GetIfscData: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/ifsc-sync", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "GetIfscData: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	ifscDataResponse := responses.NewIfscDataResponse()
	if err := ifscDataResponse.Unmarshal(respData); err != nil {
		logData.Message = "GetIfscData: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if ifscDataResponse.ErrorCode != "0" {
		logData.Message = "GetIfscData: Received error code from response"
		s.LoggerService.LogError(logData)
		return nil, errors.New(ifscDataResponse.ErrorMessage)
	}

	logData.Message = "GetIfscData API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return ifscDataResponse, nil
}

func (s *BankApiService) GetBeneficiaries(ctx context.Context, request *requests.OutgoingBeneficiarySearchRequest) (*responses.FetchBeneficiaryResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/beneficiary-fetch",
		Message:       "GetBeneficiaries log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "GetBeneficiaries: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := request.Marshal()
	if err != nil {
		logData.Message = "GetBeneficiaries: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/beneficiary-fetch", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "GetBeneficiaries: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	fetchBeneficiaryResponse := responses.NewBeneficiaryResponse()
	if err := fetchBeneficiaryResponse.Unmarshal(respData); err != nil {
		logData.Message = "GetBeneficiaries: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	// if fetchBeneficiaryResponse.ErrorCode != "0" && fetchBeneficiaryResponse.ErrorCode != "00" {
	// 	logData.Message = "GetBeneficiaries: Received Non 0 or 00 error code from response"
	// 	s.LoggerService.LogError(logData)
	// 	return nil, errors.New(fetchBeneficiaryResponse.ErrorMessage)
	// }

	logData.Message = "GetBeneficiaries API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return fetchBeneficiaryResponse, nil
}

func (s *BankApiService) AddBeneficiary(ctx context.Context, request *requests.OutgoingAddBeneficiaryRequest) (*responses.BeneficiarySubmissionResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/beneficiary-registration",
		Message:       "AddBeneficiary log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "AddBeneficiary: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := request.Marshal()
	if err != nil {
		logData.Message = "AddBeneficiary: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/beneficiary-registration", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "AddBeneficiary: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	beneficiarySubmissionResponse := responses.NewBeneficiarySubmissionResponse()
	if err := beneficiarySubmissionResponse.Unmarshal(respData); err != nil {
		logData.Message = "AddBeneficiary: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	// Scenario 3- Beneficiary details sent to KVB, but not received the response. (Back Ofice TimeOut)
	if beneficiarySubmissionResponse.ErrorCode != "0" {
		logData.Message = "AddBeneficiary: Received error code from response"
		s.LoggerService.LogError(logData)

		bankErr := &responses.BankErrorResponse{
			ErrorCode: beneficiarySubmissionResponse.ErrorCode,
		}
		return nil, bankErr
	}

	logData.Message = "AddBeneficiary API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return beneficiarySubmissionResponse, nil
}

func (s *BankApiService) SubmitOtpBeneficiaryAddition(ctx context.Context, r *requests.OutgoingBeneficiaryOtpRequest) (*responses.BeneficiaryOTPValidationResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/beneficiary-registration",
		Message:       "SubmitOtpBeneficiaryAddition log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "SubmitOtpBeneficiaryAddition: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := r.Marshal()
	if err != nil {
		logData.Message = "SubmitOtpBeneficiaryAddition: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/beneficiary-registration", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "SubmitOtpBeneficiaryAddition: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	beneficiaryOtpValidationResponse := responses.NewBeneficiaryOTPValidationResponse()
	if err := beneficiaryOtpValidationResponse.Unmarshal(respData); err != nil {
		logData.Message = "SubmitOtpBeneficiaryAddition: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	// Senerio 12- OTP Recived and OTP verification request  initiated and receive fail  response Technical Error
	// Senerio 12- OTP Recived and OTP verification request  initiated and receive fail  response back Office Time Out
	if beneficiaryOtpValidationResponse.ErrorCode != "0" && beneficiaryOtpValidationResponse.ErrorCode != "00" {
		logData.Message = "SubmitOtpBeneficiaryAddition: Received error code from response"
		s.LoggerService.LogError(logData)

		bankErr := &responses.BankErrorResponse{
			ErrorCode: beneficiaryOtpValidationResponse.ErrorCode,
		}
		return nil, bankErr
	}

	logData.Message = "SubmitOtpBeneficiaryAddition API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return beneficiaryOtpValidationResponse, nil
}

func (s *BankApiService) PaymentSubmission(ctx context.Context, request *requests.OutgoingPaymentRequest) (*responses.PaymentSubmissionResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/payment",
		Message:       "PaymentSubmission log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "PaymentSubmission: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := request.Marshal()
	if err != nil {
		logData.Message = "PaymentSubmission: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/payment", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "PaymentSubmission: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	paymentSubmissionResponse := responses.NewPaymentSubmissionResponse()
	if err := paymentSubmissionResponse.Unmarshal(respData); err != nil {
		logData.Message = "PaymentSubmission: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if paymentSubmissionResponse.ErrorCode != "0" && paymentSubmissionResponse.ErrorCode != "00" {
		logData.Message = "PaymentSubmission: Received error code from response"
		s.LoggerService.LogError(logData)
		return nil, errors.New(paymentSubmissionResponse.ErrorMessage)
	}

	logData.Message = "PaymentSubmission API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return paymentSubmissionResponse, nil
}

func (s *BankApiService) PaymentSubmissionOTP(ctx context.Context, request *requests.OutgoingPaymentRequestOTP) (*responses.PaymentSubmissionOtpResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/payment",
		Message:       "PaymentSubmissionOTP log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "PaymentSubmissionOTP: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := request.Marshal()
	if err != nil {
		logData.Message = "PaymentSubmissionOTP: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/payment", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "PaymentSubmissionOTP: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	paymentSubmissionResponse := responses.NewPaymentSubmissionOTPResponse()
	if err := paymentSubmissionResponse.Unmarshal(respData); err != nil {
		logData.Message = "PaymentSubmissionOTP: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if paymentSubmissionResponse.ErrorCode != "0" && paymentSubmissionResponse.ErrorCode != "00" {
		logData.Message = "PaymentSubmissionOTP: Received error code from response"
		s.LoggerService.LogError(logData)
		return nil, errors.New(paymentSubmissionResponse.ErrorMessage)
	}

	logData.Message = "PaymentSubmissionOTP API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return paymentSubmissionResponse, nil
}

func (s *BankApiService) PostUserConsent(ctx context.Context, request *requests.OutgoingConsentRequest, consentDetails []requests.ConsentDetails) (*responses.ConsentResponseV2, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/v2/customer/consent",
		Message:       "PostUserConsent log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "PostUserConsent: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	requestConsent := requests.ConsentRequestV2{
		ApplicantId:    request.ApplicantId,
		TxnIdentifier:  request.TxnIdentifier,
		ConsentDetails: consentDetails,
	}

	requestData, err := requestConsent.Marshal()
	if err != nil {
		logData.Message = "PostUserConsent: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(requestData)

	response, err := s.service.Post("/fintech/v2/customer/consent", requestData, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "PostUserConsent: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	consentResponse := responses.NewConsentResponseV2()
	if err := consentResponse.Unmarshal(respData); err != nil {
		logData.Message = "PostUserConsent: Error unmarshaling response body"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	// if consentResponse.ResponseHeader.ErrorCode != "0" && consentResponse.ResponseHeader.ErrorCode != "00" {
	// 	logData.Message = "PostUserConsent: Received error code from response"
	// 	logData.ResponseBody = string(body)
	// 	s.LoggerService.LogError(logData)
	// 	return nil, errors.New(consentResponse.ResponseHeader.ErrorMessage)
	// }

	logData.Message = "PostUserConsent API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return consentResponse, nil
}

func (s *BankApiService) MobileMapping(ctx context.Context, request *requests.OutgoingMobileMappingType0ApiRequest) (*responses.MobileMappingType0ApiResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/upi/mobile-mapping",
		Message:       "MobileMapping Type 0 log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "MobileMapping Type 0: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := json.Marshal(request)
	if err != nil {
		logData.Message = "MobileMapping Type 0: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/upi/mobile-mapping", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "MobileMapping Type 0: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	mobileMapping0Response := responses.NewMobileMappingType0ApiResponse()
	if err := mobileMapping0Response.UnMarshal(respData); err != nil {
		logData.Message = "MobileMapping Type 0: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "MobileMapping Type 0 API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return mobileMapping0Response, nil
}

func (s *BankApiService) VerifyUpiService(ctx context.Context, request *requests.OutgoingVerifyUserApiRequest) (*responses.VerifyUserApiResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/upi/verify-user",
		Message:       "VerifyUpiService log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "VerifyUpiService: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := request.Marshal()
	if err != nil {
		logData.Message = "VerifyUpiService: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/upi/verify-user", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "VerifyUpiService: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	verifyUserResponse := responses.NewVerifyUserApiResponse()
	if err := verifyUserResponse.UnMarshal(respData); err != nil {
		logData.Message = "VerifyUpiService: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "VerifyUpiService API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return verifyUserResponse, nil
}

func (s *BankApiService) MobileMapping1(ctx context.Context, request *requests.OutgoingMobileMappingType1ApiRequest) (*responses.MobileMappingType1ApiResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/upi/mobile-mapping",
		Message:       "MobileMapping Type 1 log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "MobileMapping Type 1: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := json.Marshal(request)
	if err != nil {
		logData.Message = "MobileMapping Type 1: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/upi/mobile-mapping", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "MobileMapping Type 1: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	mobileMapping1Response := responses.NewMobileMappingType1ApiResponse()
	if err := mobileMapping1Response.UnMarshal(respData); err != nil {
		logData.Message = "MobileMapping Type 1: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "MobileMapping Type 1 API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return mobileMapping1Response, nil
}

func (s *BankApiService) LcValidator(ctx context.Context, request *requests.OutgoingLCValidatorApiRequest) (*responses.LcValidatorApiResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/upi/lc-validator",
		Message:       "LcValidator log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "LcValidator: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := request.Marshal()
	if err != nil {
		logData.Message = "LcValidator: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/upi/lc-validator", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "LcValidator: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	lcValidatorResponse := responses.NewLcValidatorApiResponse()
	if err := lcValidatorResponse.UnMarshal(respData); err != nil {
		logData.Message = "LcValidator: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "LcValidator API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return lcValidatorResponse, nil
}

func (s *BankApiService) ProfileCreation(ctx context.Context, request *requests.OutgoingProfileCreationApiRequest) (*responses.ProfileCreationApiResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/upi/profile-creation",
		Message:       "ProfileCreation log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "ProfileCreation: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := request.Marshal()
	if err != nil {
		logData.Message = "ProfileCreation: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/upi/profile-creation", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "ProfileCreation: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	profileCreationResponse := responses.NewProfileCreationApiResponse()
	if err := profileCreationResponse.UnMarshal(respData); err != nil {
		logData.Message = "ProfileCreation: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "ProfileCreation API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return profileCreationResponse, nil
}

func (s *BankApiService) ReMapping(ctx context.Context, request *requests.OutgoingRemappingApiRequest) (*responses.ReMappingApiResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/upi/remapping",
		Message:       "ReMapping log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "ReMapping: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := request.Marshal()
	if err != nil {
		logData.Message = "ReMapping: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/upi/remapping", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "ReMapping: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	reMappingResponse := responses.NewReMappingApiResponse()
	if err := reMappingResponse.UnMarshal(respData); err != nil {
		logData.Message = "ReMapping: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "ReMapping API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return reMappingResponse, nil
}

func (s *BankApiService) AlreadyUserRequestListKeys(ctx context.Context, request *requests.OutgoingReqListKeysApiRequest) (*responses.RequestListKeysResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/upi/req-list-keys",
		Message:       "AlreadyUserRequestListKeys log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "AlreadyUserRequestListKeys: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := request.Marshal()
	if err != nil {
		logData.Message = "AlreadyUserRequestListKeys: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/upi/req-list-keys", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "AlreadyUserRequestListKeys: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	var apiResponse responses.RequestListKeysResponse
	if err := json.Unmarshal(respData, &apiResponse); err != nil {
		logData.Message = "AlreadyUserRequestListKeys: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "AlreadyUserRequestListKeys API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return &apiResponse, nil
}

func (s *BankApiService) NewUserRequestListKeys(ctx context.Context, request *requests.OutgoingReqListKeysApiRequest) (*responses.NewRequestKeyListApiResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/upi/req-list-keys",
		Message:       "NewUserRequestListKeys log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "NewUserRequestListKeys: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := request.Marshal()
	if err != nil {
		logData.Message = "NewUserRequestListKeys: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/upi/req-list-keys", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "NewUserRequestListKeys: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	var apiResponse responses.NewRequestKeyListApiResponse
	if err := json.Unmarshal(respData, &apiResponse); err != nil {
		logData.Message = "NewUserRequestListKeys: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "NewUserRequestListKeys API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return &apiResponse, nil
}

func (s *BankApiService) ExistingUserRequestListKeys(ctx context.Context, request *requests.OutgoingExistingReqlistkeysApiRequest) (*responses.ExistingUserReqListApiResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/upi/req-list-keys",
		Message:       "ExistingUserRequestListKeys log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "ExistingUserRequestListKeys: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := request.Marshal()
	if err != nil {
		logData.Message = "ExistingUserRequestListKeys: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/upi/req-list-keys", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "ExistingUserRequestListKeys: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	var apiResponse responses.ExistingUserReqListApiResponse
	if err := json.Unmarshal(respData, &apiResponse); err != nil {
		logData.Message = "ExistingUserRequestListKeys: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "ExistingUserRequestListKeys API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return &apiResponse, nil
}

func (s *BankApiService) GetXmlRequestListKeys(ctx context.Context) (*responses.XmlRequestListKeyApiResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/upi/req-list-keys",
		Message:       "GetXmlRequestListKeys log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "GetXmlRequestListKeys: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	reqListKeys := map[string]map[string]string{
		"ReqListKeys": {
			"Type": "0",
		},
	}

	body, err := json.Marshal(reqListKeys)
	if err != nil {
		logData.Message = "GetXmlRequestListKeys: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/upi/req-list-keys", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "GetXmlRequestListKeys: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	apiResponse := responses.NewXmlRequestListKeyApiResponse()
	if err := apiResponse.UnMarshal(respData); err != nil {
		logData.Message = "GetXmlRequestListKeys: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if apiResponse.Response == "" && apiResponse.ErrorCode != "0" {
		logData.Message = "GetXmlRequestListKeys: error"
		s.LoggerService.LogError(logData)
		return nil, errors.New(apiResponse.ErrorMessage)
	}

	logData.Message = "GetXmlRequestListKeys API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return apiResponse, nil
}

func (s *BankApiService) CreateUpiIdRequestListAccounts(ctx context.Context, request *requests.OutgoingCreateupiidRequestListAccountApiRequest) (*responses.CreateUpiIdRequestListAccountApiResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/upi/req-list-account",
		Message:       "CreateUpiIdRequestListAccounts log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "CreateUpiIdRequestListAccounts: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := request.Marshal()
	if err != nil {
		logData.Message = "CreateUpiIdRequestListAccounts: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/upi/req-list-account", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "CreateUpiIdRequestListAccounts: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	UserRequestListAccount := responses.NewCreateUpiIdRequestListAccountApiResponse()
	if err := UserRequestListAccount.UnMarshal(respData); err != nil {
		logData.Message = "CreateUpiIdRequestListAccounts: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "CreateUpiIdRequestListAccounts API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return UserRequestListAccount, nil
}

func (s *BankApiService) RequestPspAvailability(ctx context.Context, request *requests.OutgoingPspAvailabilityApiRequest) (*responses.PspAvailabilityApiResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/upi/psp-availability",
		Message:       "RequestPspAvailability log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "RequestPspAvailability: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := request.Marshal()
	if err != nil {
		logData.Message = "RequestPspAvailability: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/upi/psp-availability", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "RequestPspAvailability: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	UserPspAvailabilitycheck := responses.NewPspAvailabilityApiResponse()
	if err := UserPspAvailabilitycheck.UnMarshal(respData); err != nil {
		logData.Message = "RequestPspAvailability: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "RequestPspAvailability API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return UserPspAvailabilitycheck, nil
}

func (s *BankApiService) RequestAddBankAccount(ctx context.Context, request *requests.OutgoingAddBankApiRequest) (*responses.AddBankApiResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/upi/add-bank",
		Message:       "RequestAddBankAccount log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "RequestAddBankAccount: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := request.Marshal()
	if err != nil {
		logData.Message = "RequestAddBankAccount: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/upi/add-bank", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "RequestAddBankAccount: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	UserAddBankAccount := responses.NewAddBankApiResponse()
	if err := UserAddBankAccount.UnMarshal(respData); err != nil {
		logData.Message = "RequestAddBankAccount: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "RequestAddBankAccount API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return UserAddBankAccount, nil
}

func (s *BankApiService) RequestCheckAccountBalance(ctx context.Context, request *requests.OutgoingReqBalEnqApiRequest) (*responses.ReqBalEnqApiResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/upi/req-bal-enq",
		Message:       "RequestCheckAccountBalance log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "RequestCheckAccountBalance: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := request.Marshal()
	if err != nil {
		logData.Message = "RequestCheckAccountBalance: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/upi/req-bal-enq", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "RequestCheckAccountBalance: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	UserCheckAccountBalance := responses.NewReqBalEnqApiResponse()
	if err := UserCheckAccountBalance.UnMarshal(respData); err != nil {
		logData.Message = "RequestCheckAccountBalance: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "RequestCheckAccountBalance API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return UserCheckAccountBalance, nil
}

func (s *BankApiService) AadharRequestListAccount(ctx context.Context, request *requests.OutgoingAadharRequestListAccountsApiRequest) (*responses.AadharRequestListAccountApiResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/upi/req-list-account",
		Message:       "AadharRequestListAccount log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "AadharRequestListAccount: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := request.Marshal()
	if err != nil {
		logData.Message = "AadharRequestListAccount: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/upi/req-list-account", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "AadharRequestListAccount: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	UserRequestListAccount := responses.NewAadharRequestListAccountApiResponse()
	if err := UserRequestListAccount.UnMarshal(respData); err != nil {
		logData.Message = "AadharRequestListAccount: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "AadharRequestListAccount API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return UserRequestListAccount, nil
}

func (s *BankApiService) SetUpiPinReqOtp(ctx context.Context, request *requests.OutgoingSetUpiPinReqOtpApiRequest) (*responses.ReqOtpApiResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/upi/req-otp",
		Message:       "SetUPIPinReqOTP log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "SetUPIPinReqOTP: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := request.Marshal()
	if err != nil {
		logData.Message = "SetUPIPinReqOTP: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/upi/req-otp", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "SetUPIPinReqOTP: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	Setupipinotp := responses.NewReqOtpApiResponse()
	if err := Setupipinotp.UnMarshal(respData); err != nil {
		logData.Message = "SetUPIPinReqOTP: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "SetUPIPinReqOTP API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return Setupipinotp, nil
}

func (s *BankApiService) SetUpiPinReqRegMobile(ctx context.Context, request *requests.OutgoingSetUpiPinReqRegMobApiRequest) (*responses.ReqRegMobApiResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/upi/req-reg-mob",
		Message:       "SetUPIPinReqRegMobile log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "SetUPIPinReqRegMobile: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := request.Marshal()
	if err != nil {
		logData.Message = "SetUPIPinReqRegMobile: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/upi/req-reg-mob", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "SetUPIPinReqRegMobile: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	Setupiregmobile := responses.NewReqRegMobApiResponse()
	if err := Setupiregmobile.UnMarshal(respData); err != nil {
		logData.Message = "SetUPIPinReqRegMobile: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "SetUPIPinReqRegMobile API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return Setupiregmobile, nil
}

func (s *BankApiService) ValidateVpaAddress(ctx context.Context, request *requests.OutgoingReqValAddApiRequest) (*responses.ReqValAddApiResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/upi/req-val-add",
		Message:       "ValidateVpaAddress log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "ValidateVpaAddress: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := request.Marshal()
	if err != nil {
		logData.Message = "ValidateVpaAddress: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/upi/req-val-add", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "ValidateVpaAddress: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	userValidateVpa := responses.NewReqValAddApiResponse()
	if err := userValidateVpa.UnMarshal(respData); err != nil {
		logData.Message = "ValidateVpaAddress: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "ValidateVpaAddress API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return userValidateVpa, nil
}

func (s *BankApiService) PayWithVpa(ctx context.Context, request *requests.OutgoingReqPayApiRequest) (*responses.ReqPayApiResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/upi/req-pay",
		Message:       "PayWithVpa log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "PayWithVpa: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := request.Marshal()
	if err != nil {
		logData.Message = "PayWithVpa: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/upi/req-pay", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "PayWithVpa: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	userPayWithVpa := responses.NewReqPayApiResponse()
	if err := userPayWithVpa.UnMarshal(respData); err != nil {
		logData.Message = "PayWithVpa: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "PayWithVpa API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return userPayWithVpa, nil
}

func (s *BankApiService) LinkBankAccount(ctx context.Context, request *requests.OutgoingAccountLinkApiRequest) (*responses.AccountLinkApiResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/upi/account-link",
		Message:       "LinkBankAccount log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "LinkBankAccount: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := request.Marshal()
	if err != nil {
		logData.Message = "LinkBankAccount: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/upi/account-link", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "LinkBankAccount: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	userLinkBankAccount := responses.NewAccountLinkApiResponse()
	if err := userLinkBankAccount.UnMarshal(respData); err != nil {
		logData.Message = "LinkBankAccount: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "LinkBankAccount API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return userLinkBankAccount, nil
}

func (s *BankApiService) QuickTransferTemplateAdd(ctx context.Context, request *requests.OutBeneficiaryTemplateRequest) (*responses.QuickTransferBeneficiaryAdditionResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/beneficiary-registration/template/quick-transfer",
		Message:       "QuickTransferTemplateAdd log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "QuickTransferTemplateAdd: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := request.Marshal()
	if err != nil {
		logData.Message = "QuickTransferTemplateAdd: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/beneficiary-registration/template/quick-transfer", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "QuickTransferTemplateAdd: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	quickTransferResponse := responses.NewQuickTransferBeneficiaryResponse()
	if err := quickTransferResponse.Unmarshal(respData); err != nil {
		logData.Message = "QuickTransferTemplateAdd: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if quickTransferResponse.ErrorCode != "0" && quickTransferResponse.ErrorCode != "00" {
		logData.Message = "QuickTransferTemplateAdd: Received error code from response"
		s.LoggerService.LogError(logData)
		return nil, errors.New(quickTransferResponse.ErrorMessage)
	}

	logData.Message = "QuickTransferTemplateAdd API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return quickTransferResponse, nil
}

func (s *BankApiService) GetBankStatement(ctx context.Context, request *requests.OutgoingStatementRequest) (*responses.StatementResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/account-statement",
		Message:       "GetBankStatement log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "GetBankStatement: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := request.Marshal()
	if err != nil {
		logData.Message = "GetBankStatement: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/account-statement", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "GetBankStatement: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	statementResponse := responses.NewStatementResponse()
	if err := statementResponse.Unmarshal(respData); err != nil {
		logData.Message = "GetBankStatement: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if statementResponse.ErrorCode != "0" && statementResponse.ErrorCode != "00" {
		logData.Message = "GetBankStatement: Received error code from response"
		s.LoggerService.LogError(logData)
		return nil, errors.New(statementResponse.ErrorMessage)
	}

	logData.Message = "GetBankStatement API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return statementResponse, nil
}

func (s *BankApiService) RewardsTransfer(ctx context.Context, request *requests.OutgoingRewardTransactionRequest) (*responses.RewardsFundsTransferResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/payment-direct",
		Message:       "RewardsTransfer log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "RewardsTransfer: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := request.Marshal()
	if err != nil {
		logData.Message = "RewardsTransfer: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/payment-direct", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "RewardsTransfer: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	fundTransferResponse := responses.NewFundTransferResponse()
	if err := fundTransferResponse.UnMarshal(respData); err != nil {
		logData.Message = "RewardsTransfer: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if fundTransferResponse.ErrorCode != "0" && fundTransferResponse.ErrorCode != "00" {
		logData.Message = "RewardsTransfer: Received error code from response"
		s.LoggerService.LogError(logData)
		return nil, errors.New(fundTransferResponse.ErrorMessage)
	}

	logData.Message = "RewardsTransfer API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return fundTransferResponse, nil
}

func (s *BankApiService) CollectDetails(ctx context.Context, request *requests.OutgoingUpiMoneyCollectDetailsApiRequest) (*responses.UpiMoneyCollectDetailsResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/upi/collect-details",
		Message:       "CollectDetails log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "CollectDetails: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := request.Marshal()
	if err != nil {
		logData.Message = "CollectDetails: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/upi/collect-details", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "CollectDetails: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	upiCollectDetails := responses.NewUpiMoneyCollectDetailsResponse()
	if err := upiCollectDetails.UnMarshal(respData); err != nil {
		logData.Message = "CollectDetails: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "CollectDetails API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return upiCollectDetails, nil
}

func (s *BankApiService) CollectCount(ctx context.Context, request *requests.OutgoingUpiMoneyCollectCountApiRequest) (*responses.UpiMoneyCollectCountResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/upi/collect-count",
		Message:       "CollectCount log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "CollectCount: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := request.Marshal()
	if err != nil {
		logData.Message = "CollectCount: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/upi/collect-count", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "CollectCount: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	upiCollectCount := responses.NewUpiMoneyCollectCountResponse()
	if err := upiCollectCount.UnMarshal(respData); err != nil {
		logData.Message = "CollectCount: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "CollectCount API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return upiCollectCount, nil
}

func (s *BankApiService) CollectApproval(ctx context.Context, request *requests.OutgoingUpiMoneyCollectApprovalApiRequest) (*responses.UpiMoneyCollectApprovalResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/upi/approval",
		Message:       "CollectApproval log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "CollectApproval: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := request.Marshal()
	if err != nil {
		logData.Message = "CollectApproval: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/upi/approval", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "CollectApproval: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	upiCollectApproval := responses.NewUpiMoneyCollectApprovalResponse()
	if err := upiCollectApproval.UnMarshal(respData); err != nil {
		logData.Message = "CollectApproval: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "CollectApproval API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return upiCollectApproval, nil
}

func (s *BankApiService) FetchTransactionHistory(ctx context.Context, requestData requests.KVBTransactionRequest) (*responses.TransactionResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/v2/casa/transaction/enquiry",
		Message:       "FetchTransactionHistory log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "FetchTransactionHistory: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := requestData.Marshal()
	if err != nil {
		logData.Message = "FetchTransactionHistory: Error marshaling request data"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/v2/casa/transaction/enquiry", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "FetchTransactionHistory: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	FetchTransactionHistory := responses.NewTransactionResponse()
	if err := FetchTransactionHistory.UnMarshal(respData); err != nil {
		logData.Message = "FetchTransactionHistory: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if len(FetchTransactionHistory.Data) == 0 && FetchTransactionHistory.ErrorCode != "0" {
		logData.Message = "FetchTransactionHistory: response data is not correct"
		s.LoggerService.LogError(logData)
		return nil, errors.New(FetchTransactionHistory.ErrorMessage)
	}

	logData.Message = "FetchTransactionHistory API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return FetchTransactionHistory, nil
}

func (s *BankApiService) GetUpiID(request *requests.GetUpiIDRequest) (*responses.AccountLinkedResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/upi/account-link",
		Message:       "GetUpiID log",
		RequestHost:   s.service.Host,
	}

	token, err := s.GenerateToken(context.Background())
	if err != nil {
		logData.Message = "GetUpiID: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := request.Marshal()
	if err != nil {
		logData.Message = "GetUpiID: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/upi/account-link", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "GetUpiID: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	accountlinkedResponse := responses.NewAccountLinkedResponseResponse()
	if err := accountlinkedResponse.UnMarshal(respData); err != nil {
		logData.Message = "GetUpiID: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if accountlinkedResponse.Response.RESPONSECODE != "0" && accountlinkedResponse.Response.RESPONSECODE != "00" {
		logData.Message = "GetUpiID: Received error code from response"
		s.LoggerService.LogError(logData)
		return nil, errors.New(accountlinkedResponse.Response.RESPONSEMESSAGE)
	}

	logData.Message = "GetUpiID API call completed successfully"
	logData.ResponseSize = len(respData)
	logData.EndTime = time.Now()
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)
	s.LoggerService.LogInfo(logData)

	return accountlinkedResponse, nil
}

func (s *BankApiService) GetAccountDetail(ctx context.Context, request *requests.GetAccountDetail) (*responses.AccountDetailResponse, error) {
	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/v2/account/details",
		Message:       "GetAccountDetail log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(context.Background())
	if err != nil {
		logData.Message = "GetAccountDetail: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := request.Marshal()
	if err != nil {
		logData.Message = "GetAccountDetail: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.RequestBody = string(body)

	response, err := s.service.Post("/fintech/v2/account/details", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "GetAccountDetail: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	accountDetailResponse := responses.NewAccountDetailResponse()

	if err := accountDetailResponse.UnMarshal(respData); err != nil {
		logData.Message = "GetAccountDetail: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if accountDetailResponse.ErrorCode != "0" && accountDetailResponse.ErrorCode != "00" {
		logData.Message = "GetAccountDetail: Received error code from response"
		s.LoggerService.LogError(logData)
		return nil, errors.New(accountDetailResponse.ErrorMessage)
	}

	logData.Message = "GetAccountDetail API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return accountDetailResponse, nil
}

func (s *BankApiService) UploadAddressProof(ctx context.Context, req *requests.UploadAddressProofReq) (*responses.UploadAddressProofResponse, error) {
	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/v2/dms/upload",
		Message:       "UploadAddress log",
		RequestHost:   s.service.Host,
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "UploadAddress: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := req.Marshal()
	if err != nil {
		logData.Message = "UploadAddress: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/v2/dms/upload", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "UploadAddress: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)
	uploadAddressProofResponse := responses.NewUploadAddressProofResponse()

	if err := uploadAddressProofResponse.UnMarshal(respData); err != nil {
		logData.Message = "UploadAddress: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if uploadAddressProofResponse.ErrorCode != "0" && uploadAddressProofResponse.ErrorCode != "00" {
		logData.Message = "UploadAddress: Received error code from response"
		s.LoggerService.LogError(logData)
		return nil, errors.New(uploadAddressProofResponse.ErrorMessage)
	}

	logData.Message = "UploadAddress API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)
	return uploadAddressProofResponse, nil
}

func (s *BankApiService) UpdateAddress(ctx context.Context, req *requests.UpdateAddressReq) (*responses.UpdateAddressResponse, error) {
	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.BANK,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/v2/address/update",
		Message:       "UpdateAdressBankLog log",
		RequestHost:   s.service.Host,
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "UpdateAdressBankLog: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := req.Marshal()
	if err != nil {
		logData.Message = "UpdateAdressBankLog: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/v2/address/update", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "UpdateAdressBankLog: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)
	updateAddressResponse := responses.NewUpdateAddressResponse()

	if err := updateAddressResponse.UnMarshal(respData); err != nil {
		logData.Message = "UpdateAdressBankLog: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	if updateAddressResponse.ErrorCode != "0" && updateAddressResponse.ErrorCode != "00" {
		logData.Message = "UpdateAdressBankLog: Received error code from response"
		s.LoggerService.LogError(logData)
		return nil, errors.New(updateAddressResponse.ErrorMessage)
	}

	logData.Message = "UpdateAdressBankLog: API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)
	return updateAddressResponse, nil
}

func (s *BankApiService) UpiTransactionHistory(ctx context.Context, request *requests.OutgoingUpiTransactionHistoryApiRequest) (*responses.UpiTransactionHistoryApiResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/upi/transaction-history",
		Message:       "UpiTransactionHistory log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "UpiTransactionHistory: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := request.Marshal()
	if err != nil {
		logData.Message = "UpiTransactionHistory: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/upi/transaction-history", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "UpiTransactionHistory: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	transactionHistoryResponse := responses.NewUpiTransactionHistoryApiResponse()
	if err := transactionHistoryResponse.UnMarshal(respData); err != nil {
		logData.Message = "UpiTransactionHistory: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "UpiTransactionHistory API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return transactionHistoryResponse, nil
}

func (s *BankApiService) UpiChangeUpiPin(ctx context.Context, request *requests.OutgoingUpiReqSetCreRequest) (*responses.UpiChangePinApiResponse, error) {

	startTime := time.Now()

	logData := &commonSrv.LogEntry{
		Action:        constants.UPI,
		StartTime:     startTime,
		RequestMethod: "POST",
		RequestURI:    "/fintech/upi/req-set-cre",
		Message:       "UpiChangeUpiPin log",
		RequestHost:   s.service.Host,
		UserID:        utils.GetUserIDFromContext(ctx),
		RequestID:     utils.GetRequestIDFromContext(ctx),
		AppVersion:    utils.GetAppVersionFromContext(ctx),
	}

	token, err := s.GenerateToken(ctx)
	if err != nil {
		logData.Message = "UpiChangeUpiPin: Error generating token"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	body, err := request.Marshal()
	if err != nil {
		logData.Message = "UpiChangeUpiPin: Error marshaling request"
		s.LoggerService.LogError(logData)
		return nil, err
	}
	logData.RequestBody = string(body)
	response, err := s.service.Post("/fintech/upi/req-set-cre", body, map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + token.AccessToken,
		"X-Request-ID":  utils.GetRequestIDFromContext(ctx),
		"X-User-ID":     utils.GetUserIDFromContext(ctx),
		"X-App-Version": utils.GetAppVersionFromContext(ctx),
	})

	respData, err := utils.HandleResponse(response, err)
	if err != nil {
		logData.Message = "UpiChangeUpiPin: Error in POST request"
		logData.ResponseBody = string(respData)
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.ResponseSize = len(respData)
	logData.Latency = time.Since(startTime).Seconds()
	logData.ResponseBody = string(respData)

	upiPinChangeResponse := responses.NewUpiChangePinApiResponse()
	if err := upiPinChangeResponse.UnMarshal(respData); err != nil {
		logData.Message = "UpiChangeUpiPin: Error unmarshaling response body"
		s.LoggerService.LogError(logData)
		return nil, err
	}

	logData.Message = "UpiChangeUpiPin API call completed successfully"
	logData.EndTime = time.Now()
	s.LoggerService.LogInfo(logData)

	return upiPinChangeResponse, nil
}
