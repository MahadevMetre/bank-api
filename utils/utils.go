package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"bitbucket.org/paydoh/paydoh-commons/database"
	"github.com/brianvoe/gofakeit"
)

// HandleResponse processes the HTTP response and returns the body or an error if any.
func HandleResponse(resp *http.Response, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return body, fmt.Errorf("bank error: status_code:- %d: and response body:- %s", resp.StatusCode, string(body))
	}

	return body, nil
}

func GenerateReferralCode(name string) (string, error) {

	src := rand.NewSource(time.Now().UnixNano())

	r := rand.New(src)

	firstLetter := strings.ToUpper(name[:1])

	letters := make([]byte, 2)
	for i := range letters {
		letters[i] = byte(r.Intn(26) + 'A')
	}

	digits := make([]byte, 3)
	for i := range digits {
		digits[i] = byte(r.Intn(10) + '0')
	}

	referralCode := fmt.Sprintf("%s%s%s", firstLetter, string(letters), string(digits))

	return referralCode, nil
}

func StructToMap(obj interface{}) (map[string]interface{}, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func ParseName(fullName string) map[string]string {
	nameParts := strings.Fields(fullName)
	result := make(map[string]string)

	switch len(nameParts) {
	case 0:
		return result
	case 1:
		result["first_name"] = nameParts[0]
	case 2:
		result["first_name"] = nameParts[0]
		result["last_name"] = nameParts[1]
	default:
		result["first_name"] = nameParts[0]
		result["last_name"] = nameParts[len(nameParts)-1]
		result["middle_name"] = strings.Join(nameParts[1:len(nameParts)-1], " ")
	}

	return result
}

func SetRequestHeaders(req *http.Request, bearerToken string) {
	req.Header.Add("Authorization", "Bearer "+bearerToken)
	req.Header.Add("X-Device-Ip", gofakeit.IPv4Address())
	req.Header.Add("X-OS", "android")
	req.Header.Add("X-OS-Version", "14")
	latLong := fmt.Sprintf("%f,%f", gofakeit.Latitude(), gofakeit.Longitude())
	req.Header.Add("X-Lat-Long", latLong)
}

func RemoveCountryCode(phoneNumber string) string {
	if strings.HasPrefix(phoneNumber, "91") {
		return phoneNumber[2:]
	}
	return phoneNumber
}

// GetUpiUtrRefNumber extracts the UTR number from a given string
func GetUpiUtrRefNumber(input string) string {
	prefixes := []string{"UPI-CR-", "UPI-DR-"}
	var start int

	for _, prefix := range prefixes {
		start = strings.Index(input, prefix)
		if start != -1 {
			start += len(prefix)
			break
		}
	}

	if start == -1 {
		return ""
	}

	end := strings.Index(input[start:], "-")
	if end == -1 {
		return ""
	}

	return input[start : start+end]
}

// GetUserIDFromContext retrieves the user_id from the given context.
func GetUserIDFromContext(ctx context.Context) string {
	value := ctx.Value("user_id")
	if value == nil {
		return ""
	}

	userID, ok := value.(string)
	if !ok {
		return ""
	}

	return userID
}

func GetRequestIDFromContext(ctx context.Context) string {
	value := ctx.Value("request_id")
	if value == nil {
		return ""
	}

	userID, ok := value.(string)
	if !ok {
		return ""
	}

	return userID
}

func GetSourceIPFromContext(ctx context.Context) string {
	value := ctx.Value("source_ip")
	if value == nil {
		return ""
	}

	sourceIP, ok := value.(string)
	if !ok {
		return ""
	}

	return sourceIP
}

func GetDeviceOSFromContext(ctx context.Context) string {
	value := ctx.Value("device_os")
	if value == nil {
		return ""
	}

	deviceOS, ok := value.(string)
	if !ok {
		return ""
	}

	return deviceOS
}

func GetDeviceIdFromContext(ctx context.Context) string {
	value := ctx.Value("device_id")
	if value == nil {
		return ""
	}

	deviceId, ok := value.(string)
	if !ok {
		return ""
	}

	return deviceId
}

func GetAppVersionFromContext(ctx context.Context) string {
	value := ctx.Value("app_version")
	if value == nil {
		return ""
	}

	appVersion, ok := value.(string)
	if !ok {
		return ""
	}

	return appVersion
}

func StructToJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func StructToJSONString(v interface{}) (string, error) {
	bytes, err := StructToJSON(v)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func GetPackageIdFromContext(ctx context.Context) string {
	value := ctx.Value("package_id")
	if value == nil {
		return ""
	}

	packageId, ok := value.(string)
	if !ok {
		return ""
	}

	return packageId
}

func GetUserAgentFromContext(ctx context.Context) string {
	value := ctx.Value("user_agent")
	if value == nil {
		return ""
	}

	userAgent, ok := value.(string)
	if !ok {
		return ""
	}

	return userAgent
}

func RetryFunc(operation func() error, maxRetries int) error {
	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		lastErr = operation()
		if lastErr == nil {
			return nil
		}
	}
	return lastErr
}

func CalculateTimeDifference(timeString string) string {
	timeString = convertMonthCase(timeString)

	layouts := []string{
		"2006-01-02 15:04:05.0",
		"02-Jan-2006 03.04:05.000000 PM",
		"2006-01-02T15:04:05.999999-07:00",
		"2006-01-02T15:04:05-07:00",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"02-Jan-2006 03.04:05.000000 PM",
		"02-Jan-06 03.04:05.000000 PM",
	}

	location, _ := time.LoadLocation("Asia/Kolkata")
	currentTime := time.Now().In(location)

	var parsedTime time.Time
	var err error
	for _, layout := range layouts {

		parsedTime, err = time.ParseInLocation(layout, timeString, location)
		if err == nil {
			break
		}
	}

	if err != nil {
		return timeString
	}

	duration := parsedTime.Sub(currentTime)

	if duration < 0 {
		return fmt.Sprintf("1 Minute")
	}

	if duration.Minutes() < 60 {
		return fmt.Sprintf("%.0f Minute", duration.Minutes())
	} else {
		return fmt.Sprintf("%.0f Hour", duration.Hours())
	}
}

func convertMonthCase(timeString string) string {
	months := map[string]string{
		"JAN": "Jan", "FEB": "Feb", "MAR": "Mar",
		"APR": "Apr", "MAY": "May", "JUN": "Jun",
		"JUL": "Jul", "AUG": "Aug", "SEP": "Sep",
		"OCT": "Oct", "NOV": "Nov", "DEC": "Dec",
	}

	result := timeString
	for upper, proper := range months {
		result = strings.Replace(result, upper, proper, 1)
	}
	return result
}

func SaveHashMapToRedis(ctx context.Context, redis *database.InMemory, userID, key string, data interface{}, deleteTime time.Duration) error {
	redisKey := fmt.Sprintf(key, userID)

	err := redis.GetClient().HSet(context.Background(), redisKey, data).Err()
	if err != nil {
		return fmt.Errorf("failed to save data to Redis: %w", err)
	}

	err = redis.GetClient().Expire(context.Background(), redisKey, deleteTime).Err()
	if err != nil {
		return fmt.Errorf("failed to set expiration time for Redis key: %w", err)
	}

	return nil
}

func GetHashDataFromRedis(ctx context.Context, redis *database.InMemory, userID, key string) (map[string]string, error) {
	data, err := redis.GetClient().HGetAll(context.Background(), fmt.Sprintf(key, userID)).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve data from Redis: %w", err)
	}

	if len(data) == 0 {
		return nil, errors.New("data not found in redis")
	}

	return data, nil
}

func DateFormat(dateString string) string {
	originalFormat := "02-Jan-2006 03.04:05.000000 PM"
	parsedTime, err := time.Parse(originalFormat, dateString)
	if err != nil {
		return ""
	}
	newFormat := "02-01-2006 03:04 PM"
	return parsedTime.Format(newFormat)
}
