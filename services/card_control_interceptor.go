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
	"net"
	"net/http"
	"net/url"
	"strings"

	commonSrv "bitbucket.org/paydoh/paydoh-commons/services"
)

type CardControlInterceptor struct {
	Transport  http.RoundTripper
	MaxRetries int
	log        *commonSrv.LoggerService
	bankSrv    *BankApiService
}

type BankCardControlErrorWithRetry struct {
	BankError   *responses.BankErrorResponse
	ShouldRetry bool
}

func (e *BankCardControlErrorWithRetry) Error() string {
	return e.BankError.ErrorMessage
}

func (i *CardControlInterceptor) shouldSkip(url string) bool {
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
func (i *CardControlInterceptor) RoundTrip(req *http.Request) (*http.Response, error) {
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
	transactionID := req.Header.Get("X-Transaction-ID")
	skipEncryption := i.shouldSkip(req.URL.Path)
	clonedReq := req.Clone(req.Context())
	bodyBytes, err := i.prepareRequestBody(clonedReq, skipEncryption, logEntry, transactionID)
	if err != nil {
		i.log.LogError(logEntry)
		return nil, err
	}
	req.Header.Del("X-Transaction-ID")
	clonedReq.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	clonedReq.ContentLength = int64(len(bodyBytes))
	clonedReq.Header = req.Header
	return i.processRequest(clonedReq, bodyBytes, skipEncryption, logEntry)
}

func (i *CardControlInterceptor) processRequest(clonedReq *http.Request, bodyBytes []byte, skipEncryption bool, logEntry *commonSrv.LogEntry) (*http.Response, error) {
	var lastErr error

	for attempt := 1; attempt <= i.MaxRetries; attempt++ {
		resp, err := i.Transport.RoundTrip(clonedReq)
		if err != nil {
			if urlErr, ok := err.(*url.Error); ok {
				// Check for DNS resolution errors
				if dnsErr, ok := urlErr.Err.(*net.DNSError); ok {
					fmt.Printf("DNS error: no such host: %v\n", dnsErr.Name)
					lastErr = dnsErr
					logEntry.Message = fmt.Sprintf("Retrying request (attempt %d/%d) due to DNS error", attempt, i.MaxRetries)
					i.log.LogInfo(logEntry)
					clonedReq.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
					continue
				}
				// Check for timeout errors
				if urlErr.Timeout() {
					fmt.Println("Request timed out")
					lastErr = urlErr
					logEntry.Message = fmt.Sprintf("Retrying request (attempt %d/%d) due to timeout", attempt, i.MaxRetries)
					i.log.LogInfo(logEntry)
					clonedReq.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
					continue
				}
				// Handle other URL errors
				lastErr = urlErr
				logEntry.Message = fmt.Sprintf("Retrying request (attempt %d/%d) due to URL error", attempt, i.MaxRetries)
				i.log.LogInfo(logEntry)
				clonedReq.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				continue
			}

			// Log non-URL errors
			// lastErr = fmt.Errorf("request failed: %w", err)
			logEntry.Message = fmt.Sprintf("Request failed due to an unexpected error %d/%d", attempt, i.MaxRetries)
			// i.log.LogError(logEntry)
			// continue
			// logEntry.ResponseStatusCode = resp.StatusCode
			lastErr = fmt.Errorf("request failed: %w", err)
			i.log.LogError(logEntry)
			continue
		}

		logEntry.ResponseStatusCode = resp.StatusCode

		if resp.StatusCode == http.StatusOK {
			response, err := i.handleResponse(resp, bodyBytes, skipEncryption, logEntry)
			if err != nil {
				if retryErr, ok := err.(*BankCardControlErrorWithRetry); ok && retryErr.ShouldRetry {
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
		} else if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusBadRequest || resp.StatusCode == http.StatusRequestTimeout {
			resp.Body.Close()
			logEntry.Message = fmt.Sprintf("Retrying request due to status %d (attempt %d/%d)", resp.StatusCode, attempt+1, i.MaxRetries)
			logEntry.ResponseStatusCode = resp.StatusCode
			i.log.LogInfo(logEntry)
			clonedReq.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			lastErr = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
			continue
		} else {
			lastErr = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
			i.log.LogError(logEntry)
			continue
		}
	}

	i.log.LogError(logEntry)

	fmt.Printf("max retries exceeded after %d attempts, last error: %v", i.MaxRetries, lastErr)
	return nil, lastErr
}

func (i *CardControlInterceptor) handleResponse(resp *http.Response, bodyBytes []byte, skipEncryption bool, logEntry *commonSrv.LogEntry) (*http.Response, error) {
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
			if err := i.processDebitCardControlError(encryptedBody, logEntry); err != nil {
				return nil, err
			}
		}
		i.log.LogError(logEntry)
		return nil, err
	}

	if err := i.processDebitCardControlError(decryptedBody, logEntry); err != nil {
		return nil, err
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

func (i *CardControlInterceptor) prepareRequestBody(req *http.Request, skipEncryption bool, logEntry *commonSrv.LogEntry, transactionID string) ([]byte, error) {
	var bodyBytes []byte
	var err error
	if req.Body != nil {
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			logEntry.Message = fmt.Sprintf("Failed to read request body: %v", err)
			return nil, fmt.Errorf("failed to read request body: %w", err)
		}
		req.Body.Close()
	}
	if !skipEncryption && len(bodyBytes) > 0 {
		encReq, err := utils.DebitCardGenerateEncryptedReq(bodyBytes, transactionID, constants.CardControlEncryptionKey)
		if err != nil {
			logEntry.Message = fmt.Sprintf("Failed to encrypt request: %v", err)
			return nil, err
		}
		req, err := json.Marshal(encReq)
		if err != nil {
			logEntry.Message = fmt.Sprintf("Failed to encrypt request: %v", err)
			return nil, err
		}
		logEntry.RequestBody = string(req)
		return req, nil
	}
	return bodyBytes, nil
}

func (i *CardControlInterceptor) handleUnauthorized(attempt int, logEntry *commonSrv.LogEntry) (*responses.TokenResponse, error) {
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

func (i *CardControlInterceptor) handleEncryptedResponse(encryptedBody []byte, logEntry *commonSrv.LogEntry) ([]byte, error) {
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
	decryptRespData, err := utils.DebitCardDecryptResponse(decryptResp.EncryptRes, constants.CardControlEncryptionKey)
	if err != nil {
		logEntry.Message = fmt.Sprintf("Failed to decrypt response: %v", err)
		return nil, fmt.Errorf("failed to decrypt response: %w", err)
	}
	return decryptRespData, nil
}

func (i *CardControlInterceptor) processDebitCardControlError(decryptedBody []byte, logEntry *commonSrv.LogEntry) error {

	var debitCardError responses.BankErrorResponse
	var array []responses.BankErrorResponse

	if err := json.Unmarshal(decryptedBody, &debitCardError); err == nil {
		return i.ProcessBankErrorResponse([]responses.BankErrorResponse{debitCardError}, logEntry)
	}

	if err := json.Unmarshal(decryptedBody, &array); err == nil {
		return i.ProcessBankErrorResponse(array, logEntry)
	}

	return fmt.Errorf("failed to unmarshal response into known formats")
}

func (i *CardControlInterceptor) ProcessBankErrorResponse(response []responses.BankErrorResponse, logEntry *commonSrv.LogEntry) error {
	if len(response) > 0 {
		debitCardError := response[0]
		if debitCardError.ErrorCode1 == "0" || debitCardError.ErrorCode1 == "00" ||
			debitCardError.ErrorCode == "0" || debitCardError.ErrorCode == "00" {
			return nil
		}

		if msg, exists := constants.RetryBankErrors[debitCardError.ErrorCode1]; exists {
			logEntry.Message = fmt.Sprintf("Retriable bank error: %s - %s", debitCardError.ErrorCode1, debitCardError.ErrorMessage1)
			i.log.LogError(logEntry)
			debitCardError.ErrorMessage = msg
			return &BankCardControlErrorWithRetry{BankError: &debitCardError, ShouldRetry: true}
		}

		if msg, exists := constants.RetryBankErrors[debitCardError.ErrorCode]; exists {
			logEntry.Message = fmt.Sprintf("Retriable bank error: %s - %s", debitCardError.ErrorCode1, debitCardError.ErrorMessage1)
			i.log.LogError(logEntry)
			debitCardError.ErrorMessage = msg
			return &BankCardControlErrorWithRetry{BankError: &debitCardError, ShouldRetry: true}
		}

		if msg, exists := constants.GetCardControlService(debitCardError.ErrorCode1); exists {
			logEntry.Message = fmt.Sprintf("Non-retriable bank error: %s - %s", debitCardError.ErrorCode1, debitCardError.ErrorMessage1)
			i.log.LogError(logEntry)
			debitCardError.ErrorMessage = msg
			return &BankCardControlErrorWithRetry{BankError: &debitCardError, ShouldRetry: false}
		}

		logEntry.Message = fmt.Sprintf("Non-retriable bank error: %s - %s", debitCardError.ErrorCode1, debitCardError.ErrorMessage1)
		i.log.LogError(logEntry)
		return &BankCardControlErrorWithRetry{BankError: &debitCardError, ShouldRetry: false}
	}

	logEntry.Message = "no response found"
	i.log.LogError(logEntry)
	return errors.New("no response found")
}
