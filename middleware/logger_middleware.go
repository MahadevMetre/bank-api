package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	commonConstants "bitbucket.org/paydoh/paydoh-commons/constants"
	"bitbucket.org/paydoh/paydoh-commons/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func LoggerMiddleware(loggerService *services.LoggerService) gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/swagger") || c.Request.URL.Path == "/health" {
			c.Next()
			return
		}

		// Check if request ID is already present in the headers
		requestId := c.GetHeader("X-Request-ID")

		if requestId == "" {
			requestId = uuid.New().String()
		}

		startTime := time.Now()

		userId := getUserIDFromToken(c)

		deviceIp, _ := GetHeaderData(c, "X-Device-Ip")

		deviceOS, _ := GetHeaderData(c, "X-OS")

		latLongData, _ := GetHeaderData(c, "X-Lat-Long")

		appVersion, _ := GetHeaderData(c, "X-App-Version")

		// copy of the response writer
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		// Create a map to store header information
		headerMap := make(map[string]string)

		// Iterate through all headers
		for name, values := range c.Request.Header {
			// Join multiple values with comma if present
			headerMap[name] = values[0]
			if len(values) > 1 {
				headerMap[name] = fmt.Sprintf("%s", values)
			}
		}

		// Convert the map to a JSON string
		jsonHeaders, err := json.MarshalIndent(headerMap, "", "  ")
		if err != nil {
			fmt.Println("Error marshalling headers:", err)
		}

		// Log before processing the request
		logData := &services.LogEntry{
			Action:        commonConstants.MIDDLEWARE_START,
			UserID:        userId,
			DeviceIP:      deviceIp,
			LatLong:       latLongData,
			StartTime:     startTime,
			RequestMethod: c.Request.Method,
			RequestURI:    c.Request.RequestURI,
			RequestHeader: string(jsonHeaders),
			UserAgent:     c.Request.UserAgent(),
			RequestID:     requestId,
			RequestHost:   c.Request.Host,
			Message:       "middleware start log",
			AppVersion:    appVersion,
		}

		if loggerService.EnableRequestLog {
			byteData, err := ReadRequestBody(c.Request)
			if err != nil {
				return
			}

			logData.RequestBody = string(byteData)
		}
		loggerService.LogInfo(logData)

		ctx := c.Request.Context()

		ctx = context.WithValue(ctx, "request_id", requestId)
		ctx = context.WithValue(ctx, "user_id", userId)
		ctx = context.WithValue(ctx, "source_ip", deviceIp)
		ctx = context.WithValue(ctx, "app_version", appVersion)
		ctx = context.WithValue(ctx, "device_os", deviceOS)
		ctx = context.WithValue(ctx, "device_id", c.Request.Header.Get("X-Device-ID"))
		ctx = context.WithValue(ctx, "package_id", c.Request.Header.Get("X-Package-ID"))
		ctx = context.WithValue(ctx, "user_agent", c.Request.Header.Get("User-Agent"))

		c.Request = c.Request.WithContext(ctx)

		c.Next()

		endTime := time.Now()

		// Log after processing the request
		resLogData := &services.LogEntry{
			Action:        commonConstants.MIDDLEWARE_END,
			UserID:        userId,
			DeviceIP:      deviceIp,
			LatLong:       latLongData,
			StartTime:     time.Now(),
			EndTime:       endTime,
			RequestMethod: c.Request.Method,
			RequestURI:    c.Request.RequestURI,
			UserAgent:     c.Request.UserAgent(),
			RequestID:     requestId,
			ResponseSize:  blw.body.Len(),
			RequestHost:   c.Request.Host,
			Latency:       endTime.Sub(startTime).Seconds(),
			Message:       "middleware end log",
		}

		if loggerService.EnableResponseLog {
			resLogData.ResponseBody = blw.body.String()
		}

		loggerService.LogInfo(resLogData)

	}
}

// reads the request body
func ReadRequestBody(r *http.Request) ([]byte, error) {
	if r.Body == nil {
		return nil, nil
	}
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	return bodyBytes, nil
}

// custom response writer that captures the response body
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func ExtractBearerToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return ""
	}

	return parts[1]
}

func getUserIDFromToken(ctx *gin.Context) string {
	token := ExtractBearerToken(ctx)
	if token == "" {
		return ""
	}

	tokenData, err := VerifyJwtToken(token, ctx.Request.Header.Get("User-Agent"), ctx.Request.Header.Get("X-Device-Ip"))
	if err != nil {
		return ""
	}

	return tokenData.UserId
}

func GetHeaderData(ctx *gin.Context, key string) (string, error) {
	keyData := ctx.Request.Header.Get(key)
	if keyData == "" {
		return "", nil
	}

	return keyData, nil
}
