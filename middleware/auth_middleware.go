package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"bitbucket.org/paydoh/paydoh-commons/customerror"
	"bitbucket.org/paydoh/paydoh-commons/jsonwebtoken"
	"bitbucket.org/paydoh/paydoh-commons/responses"
	"github.com/gin-gonic/gin"

	"bankapi/config"
	"bankapi/constants"
	"bankapi/models"
	"bankapi/security"
)

// AuthMiddleware validates the authorization header and sets the user_id, key, device_ip, os, os_version and lat_long in the gin context.
// It returns an error if the authorization header is not provided, the token is invalid, or the required headers are not provided.
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.Request.Header.Get("Authorization")
		if authorizationHeader == "" {
			responses.StatusUnauthorized(ctx, customerror.NewError(errors.New("authorization header is not provided")))
			ctx.Abort()
			return
		}

		if !strings.Contains(authorizationHeader, "Bearer") {
			responses.StatusUnauthorized(ctx, customerror.NewError(errors.New("invalid token !!!, token should be in Bearer <token> format")))
			ctx.Abort()
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) < 2 {
			responses.StatusUnauthorized(ctx, customerror.NewError(errors.New("invalid token !!!, token should be in Bearer <token> format")))
			ctx.Abort()
			return
		}
		token := headerParts[1]

		tokenData, err := VerifyJwtToken(token, ctx.Request.Header.Get("User-Agent"), ctx.Request.Header.Get("X-Device-Ip"))
		if err != nil {
			responses.StatusUnauthorized(
				ctx,
				customerror.NewError(err),
			)
			ctx.Abort()
			return
		}

		userID := tokenData.UserId
		signKey := tokenData.SigningKey

		deviceIp := ctx.Request.Header.Get("X-Device-Ip")
		if deviceIp == "" {
			responses.StatusBadRequest(ctx, customerror.NewError(errors.New("device ip is not provided")), "")
			ctx.Abort()
			return
		}

		os := ctx.Request.Header.Get("X-OS")
		if os == "" {
			responses.StatusBadRequest(ctx, customerror.NewError(errors.New("os is not provided")), "")
			ctx.Abort()
			return
		}

		osVersion := ctx.Request.Header.Get("X-OS-Version")
		if osVersion == "" {
			responses.StatusBadRequest(ctx, customerror.NewError(errors.New("os version is not provided")), "")
			ctx.Abort()
			return
		}

		latLong := ctx.Request.Header.Get("X-Lat-Long")
		if latLong == "" {
			responses.StatusBadRequest(ctx, customerror.NewError(errors.New("lat long is not provided")), "")
			ctx.Abort()
			return
		}

		skipURLs := []string{
			"/api/authentication/initiate-sim-verification",
			"/api/authentication/sim-verification-status",
			"/api/upi/remapping-upi-id",
			"/api/upi/simbinding/sms-verification",
			"/api/get-secrets",
		}

		deviceID := ctx.Request.Header.Get("X-Device-ID")
		if deviceID != "" {
			accountData, _ := models.GetAccountDataByUserIdV2(userID)

			if !containsURL(skipURLs, ctx.Request.RequestURI) && accountData != nil && accountData.UpiId.Valid {
				deviceData, err := models.FindOneDeviceByUserIDV2(userID)
				if err != nil {
					ctx.JSON(
						http.StatusOK,
						gin.H{
							"status":     http.StatusUnauthorized,
							"message":    "Internal Server Error",
							"ERROR_CODE": "ER001",
							"error":      customerror.NewError(err),
						},
					)
					ctx.Abort()
					return
				}

				decryptedDeviceId, err := security.Decrypt(deviceData.DeviceId, []byte(strings.TrimSpace(signKey)))
				if err != nil {
					responses.StatusInternalServerError(ctx, customerror.NewError(err), "")
					ctx.Abort()
					return
				}

				if deviceID != decryptedDeviceId && !deviceData.IsActive {
					ctx.JSON(
						http.StatusOK,
						gin.H{
							"status":     http.StatusUnauthorized,
							"message":    "Unauthorized, device id mismatch",
							"ERROR_CODE": "ER001",
							"error":      customerror.NewError(errors.New("unauthorized")),
						},
					)
					ctx.Abort()
					return
				}
			}
		}

		ctx.Set("user_id", userID)
		ctx.Set("key", signKey)
		ctx.Set("device_ip", deviceIp)
		ctx.Set("os", os)
		ctx.Set("os_version", osVersion)
		ctx.Set("lat_long", latLong)
		ctx.Next()
	}
}

func CallbackMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		apiKey := ctx.GetHeader("X-API-Key")

		if apiKey == "" {
			responses.StatusBadRequest(ctx, customerror.NewError(errors.New("api key is not provided")), "")
			ctx.Abort()
			return
		}

		key := os.Getenv("CALLBACK_API_KEY")

		if apiKey != key {
			responses.StatusBadRequest(ctx, customerror.NewError(errors.New("api key is not valid")), "")
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

// containsURL checks if the given URL exists in the slice of URLs
func containsURL(urls []string, url string) bool {
	for _, u := range urls {
		if u == url {
			return true
		}
	}
	return false
}

func VerifyJwtToken(token, userAgent, userIP string) (*jsonwebtoken.ResponseClaims, error) {
	tokenData, err := jsonwebtoken.VerifyWithClaims(token, constants.JwtKey)
	if err != nil {
		return nil, err
	}

	if tokenData.DeviceInfo != userAgent && tokenData.UserIp != userIP {
		return nil, errors.New("invalid token")
	}

	data, err := extractAndSwapUserIDAndSignKey(tokenData.UserId)
	if err != nil {
		return nil, err
	}

	// for logout(token blacklist)
	_, err = config.GetRedis().Get(fmt.Sprintf(constants.TokenKeyFormat, data["user_id"]))
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("generate new token to continue")
	}

	jwtData := &jsonwebtoken.ResponseClaims{}

	jwtData.UserId = data["user_id"]
	jwtData.SigningKey = data["sign_key"]

	return jwtData, nil
}

func extractAndSwapUserIDAndSignKey(userID string) (map[string]string, error) {
	tokenDataParts := strings.Split(userID, "|")
	if len(tokenDataParts) < 2 {
		return nil, errors.New("invalid token")
	}

	userId := tokenDataParts[0]
	signKey := tokenDataParts[1]

	swappedUserId := signKey[:3] + userId[3:]
	swappedSignKey := userId[:3] + signKey[3:]

	return map[string]string{"user_id": swappedUserId, "sign_key": swappedSignKey}, nil
}
