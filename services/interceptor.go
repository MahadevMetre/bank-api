package services

import (
	"bankapi/constants"
	"bankapi/responses"
	"bankapi/utils"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	commonSrv "bitbucket.org/paydoh/paydoh-commons/services"
)

type Interceptor struct {
	Transport  http.RoundTripper
	MaxRetries int
	log        *commonSrv.LoggerService
	bankSrv    *BankApiService
}

type BankErrorWithRetry struct {
	BankError   *responses.BankErrorResponse
	ShouldRetry bool
}

func (e *BankErrorWithRetry) Error() string {
	return e.BankError.ErrorMessage
}

func (i *Interceptor) shouldSkip(url string) bool {
	skipPatterns := []string{
		"/oauth/cc/accesstoken",
	}
	for _, pattern := range skipPatterns {
		if strings.Contains(url, pattern) {
			return true
		}
	}
	return false
}

func (i *Interceptor) RoundTrip(req *http.Request) (*http.Response, error) {
	if i.Transport == nil {
		i.Transport = http.DefaultTransport
	}

	logEntry := &commonSrv.LogEntry{
		Action:      "INTERCEPTOR",
		Message:     "Bank API Request",
		RequestHost: req.Host,
		RequestURI:  req.URL.Path,
		RequestID:   req.Header.Get("X-Request-ID"),
		UserID:      req.Header.Get("X-User-ID"),
		AppVersion:  req.Header.Get("X-App-Version"),
	}

	req.Header.Del("X-Request-ID")
	req.Header.Del("X-User-ID")
	req.Header.Del("X-App-Version")

	skipEncryption := i.shouldSkip(req.URL.Path)
	clonedReq := req.Clone(req.Context())

	bodyBytes, err := i.prepareRequestBody(clonedReq, skipEncryption, logEntry)
	if err != nil {
		i.log.LogError(logEntry)
		return nil, err
	}

	clonedReq.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	clonedReq.ContentLength = int64(len(bodyBytes))
	clonedReq.Header = req.Header

	return i.processRequest(clonedReq, bodyBytes, skipEncryption, logEntry)
}

func (i *Interceptor) processRequest(clonedReq *http.Request, bodyBytes []byte, skipEncryption bool, logEntry *commonSrv.LogEntry) (*http.Response, error) {
	var lastErr error
	for attempt := 1; attempt <= i.MaxRetries; attempt++ {
		resp, err := i.Transport.RoundTrip(clonedReq)
		if err != nil {
			lastErr = fmt.Errorf("request failed: %w", err)
			i.log.LogError(logEntry)
			continue
		}

		logEntry.ResponseStatusCode = resp.StatusCode

		if resp.StatusCode == http.StatusOK {
			response, err := i.handleResponse(resp, bodyBytes, skipEncryption, logEntry)
			if err != nil {
				if retryErr, ok := err.(*BankErrorWithRetry); ok && retryErr.ShouldRetry {
					lastErr = retryErr
					logEntry.Message = fmt.Sprintf("Retrying request (attempt %d/%d)", attempt+1, i.MaxRetries)
					i.log.LogInfo(logEntry)
					clonedReq.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
					continue
				}
				return nil, err
			}
			return response, nil
		} else if resp.StatusCode == http.StatusUnauthorized {
			resp.Body.Close()
			tokenData, err := i.handleUnauthorized(attempt, logEntry)
			if err != nil {
				lastErr = err
				i.log.LogError(logEntry)
				continue
			}

			clonedReq.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			clonedReq.Header.Set("Authorization", "Bearer "+tokenData.AccessToken)
			continue
		} else {
			return resp, nil
		}
	}
	fmt.Printf("max retries exceeded after %d attempts, last error: %v", i.MaxRetries, lastErr)
	return nil, lastErr
}

func (i *Interceptor) handleResponse(resp *http.Response, bodyBytes []byte, skipEncryption bool, logEntry *commonSrv.LogEntry) (*http.Response, error) {
	if skipEncryption {
		logEntry.Message = "Request completed (unencrypted)"
		i.log.LogInfo(logEntry)
		return resp, nil
	}

	encryptedBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logEntry.Message = fmt.Sprintf("Failed to read response body: %v", err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	resp.Body.Close()

	decryptedBody, err := i.handleEncryptedResponse(encryptedBody, logEntry)
	if err != nil {
		if err.Error() == "'encryptRes' key is not getting in response" {
			if err := i.processBankError(encryptedBody, logEntry); err != nil {
				return nil, err
			}
		}
		i.log.LogError(logEntry)
		return nil, err
	}

	logEntry.Message = fmt.Sprintf("Decrypted response body: %s", string(decryptedBody))
	i.log.LogInfo(logEntry)

	// Scenario 3- Beneficiary details sent to KVB, but not received the response. (Back Ofice TimeOut)
	// Senerio 12- OTP Recived and OTP verification request  initiated and receive fail  response Technical Error
	// Senerio 12- OTP Recived and OTP verification request  initiated and receive fail  response back Office Time Out
	if resp.Request.URL.Path == "/fintech/upi/req-list-keys" || resp.Request.URL.Path == "/fintech/upi/transaction-history" || resp.Request.URL.Path == "/fintech/beneficiary-registration" {
		logEntry.Message = "skipping upi api processing error check"
		i.log.LogInfo(logEntry)
	} else {
		// handle bank error
		if err := i.processBankError(decryptedBody, logEntry); err != nil {
			return nil, err
		}
	}

	// Create new response with decrypted body
	newResp := *resp
	newResp.Body = io.NopCloser(bytes.NewBuffer(decryptedBody))
	newResp.ContentLength = int64(len(decryptedBody))

	logEntry.Message = "Request completed (encrypted)"
	logEntry.RequestBody = string(bodyBytes)
	logEntry.ResponseBody = string(encryptedBody)
	i.log.LogInfo(logEntry)

	return &newResp, nil
}

func (i *Interceptor) prepareRequestBody(req *http.Request, skipEncryption bool, logEntry *commonSrv.LogEntry) ([]byte, error) {
	var bodyBytes []byte
	var err error

	if req.Body != nil {
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			logEntry.Message = fmt.Sprintf("Failed to read request body: %v", err)
			return nil, fmt.Errorf("failed to read request body: %w", err)
		}
		req.Body.Close()

		logEntry.Message = fmt.Sprintf("Original request body before encryption: %s", string(bodyBytes))
		i.log.LogInfo(logEntry)
	}

	if !skipEncryption && len(bodyBytes) > 0 {
		encReq, err := utils.GenerateEncryptedReqV2(bodyBytes, constants.BankEncryptionKey)
		if err != nil {
			logEntry.Message = fmt.Sprintf("Failed to encrypt request: %v", err)
			return nil, err
		}
		logEntry.RequestBody = string(encReq)
		return encReq, nil
	}

	return bodyBytes, nil
}

func (i *Interceptor) handleUnauthorized(attempt int, logEntry *commonSrv.LogEntry) (*responses.TokenResponse, error) {
	if err := i.bankSrv.Memory.Delete("kvb_auth_token"); err != nil {
		logEntry.Message = fmt.Sprintf("Error deleting auth token on attempt %d: %v", attempt+1, err)
		return nil, err
	}

	ctx := context.Background()
	tokenData, err := i.bankSrv.GenerateToken(ctx)
	if err != nil {
		logEntry.Message = fmt.Sprintf("Error generating new token on attempt %d: %v", attempt+1, err)
		return nil, err
	}

	logEntry.Message = fmt.Sprintf("Received 401 Unauthorized on attempt %d, retrying with new token", attempt+1)
	return tokenData, nil
}

func (i *Interceptor) handleEncryptedResponse(encryptedBody []byte, logEntry *commonSrv.LogEntry) ([]byte, error) {
	logEntry.ResponseBody = string(encryptedBody)

	if len(encryptedBody) == 0 {
		return []byte{}, nil
	}

	decryptResp := new(responses.EncryptedRes)
	if err := json.Unmarshal(encryptedBody, decryptResp); err != nil {
		logEntry.Message = fmt.Sprintf("Failed to unmarshal encrypted response: %v", err)
		return nil, err
	}

	if decryptResp.EncryptRes == "" {
		return nil, errors.New("'encryptRes' key is not getting in response")
	}

	decryptRespData, err := utils.DecryptResponse(decryptResp.EncryptRes, constants.BankEncryptionKey)
	if err != nil {
		logEntry.Message = fmt.Sprintf("Failed to decrypt response: %v", err)
		return nil, fmt.Errorf("failed to decrypt response: %w", err)
	}

	return decryptRespData, nil
}

func (i *Interceptor) processBankError(decryptedBody []byte, logEntry *commonSrv.LogEntry) error {
	var bankError responses.BankErrorResponse
	if err := json.Unmarshal(decryptedBody, &bankError); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if bankError == (responses.BankErrorResponse{}) || bankError.ErrorCode == "0" || bankError.ErrorCode == "00" || bankError.ErrorCode == "01" ||
		strings.EqualFold(bankError.ErrorMessage, "Success") || bankError.UpiError.ResponseCode == "0" ||
		bankError.UpiError.ResponseCode == "2" || bankError.UpiError.ResponseCode == "4" || bankError.UpiError.ResponseCode == "5" {
		return nil
	}

	// check retry bank error
	if msg, exists := constants.RetryBankErrors[bankError.ErrorCode]; exists {
		logEntry.Message = fmt.Sprintf("Retriable bank error: %s - %s", bankError.ErrorCode, bankError.ErrorMessage)
		logEntry.ResponseBody = string(decryptedBody)
		i.log.LogError(logEntry)
		bankError.ErrorMessage = msg
		return &BankErrorWithRetry{BankError: &bankError, ShouldRetry: true}
	}

	// check retry upi error
	if msg, exists := constants.RetryUpiErrors[bankError.UpiError.ResponseCode]; exists {
		logEntry.Message = fmt.Sprintf("Retriable UPI error: %s - %s", bankError.UpiError.ResponseCode, bankError.UpiError.ResponseMessage)
		logEntry.ResponseBody = string(decryptedBody)
		i.log.LogError(logEntry)
		bankError.ErrorMessage = msg
		return &BankErrorWithRetry{BankError: &bankError, ShouldRetry: true}
	}

	logEntry.Message = fmt.Sprintf("Non-retriable bank error: %s - %s", bankError.ErrorCode, bankError.ErrorMessage)
	i.log.LogError(logEntry)
	return &BankErrorWithRetry{BankError: &bankError, ShouldRetry: false}
}
