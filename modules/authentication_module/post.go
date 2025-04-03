package authenticationmodule

import (
	"fmt"

	"bitbucket.org/paydoh/paydoh-commons/customerror"
	"bitbucket.org/paydoh/paydoh-commons/responses"
	"github.com/gin-gonic/gin"

	"bankapi/requests"
	"bankapi/stores"
)

// @Summary Api get message token by sending encrypted data and crc value
// @Tags Authentication api
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "With the bearer started"
// @Param X-Device-Ip header string true "With the device ip"
// @Param X-OS header string true "With the os"
// @Param X-OS-Version header string true "With the os version"
// @Param X-Lat-Long header string true "With the lat long"
// @Param authenticationRequest body requests.AuthenticationRequest true "Authentication Request"
// @Success 200 {object} responses.MobileTeamSuccessResponse "success response"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Router /api/authentication/initiate-sim-verification/ [post]
func InitiateSimVerification(c *gin.Context) {

	authValues, err := stores.GetAuthValue(c)

	if err != nil {
		fmt.Println("ATUH ERR ", err)
		responses.StatusForbidden(
			c,
			customerror.NewError(err),
		)
		return
	}

	store, err := stores.GetStores(c)

	if err != nil {
		fmt.Println("STORRE ERRR ", err)
		responses.StatusInternalServerError(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	request := requests.NewAuthenticationRequest()

	if err := request.Validate(c); err != nil {
		fmt.Println("ERRR REQ ", err)
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	result, err := store.Authentication.InitiateSimVerification(c.Request.Context(), authValues.UserId, authValues.DeviceIp, authValues.OS, authValues.OSVersion, authValues.Key, request)

	if err != nil {
		fmt.Println("ERR STORE ", err)
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	responses.StatusOk(
		c,
		result,
		"successfully created device / updated device",
		"",
	)
}

// GetSimVerificationStatus Get sim verification status
// @Tags Authentication api
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "With the bearer started"
// @Param X-Device-Ip header string true "With the device ip"
// @Param X-OS header string true "With the os"
// @Param X-OS-Version header string true "With the os version"
// @Param X-Lat-Long header string true "With the lat long"
// @Success 200 {object} responses.MobileTeamSuccessResponse "success response"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Router /api/authentication/sim-verification-status [GET]
func GetSimVerificationStatus(ctx *gin.Context) {
	authValues, err := stores.GetAuthValue(ctx)
	if err != nil {
		responses.StatusForbidden(
			ctx,
			customerror.NewError(err),
		)
		return
	}

	store, err := stores.GetStores(ctx)
	if err != nil {
		responses.StatusInternalServerError(
			ctx,
			customerror.NewError(err),
			"",
		)
		return
	}

	respData, err := store.Authentication.GetSimVerificationStatus(ctx.Request.Context(), authValues.UserId)
	if err != nil {
		responses.StatusInternalServerError(
			ctx,
			customerror.NewError(err),
			"",
		)
		return
	}

	responses.StatusOk(
		ctx,
		respData,
		"success",
		"",
	)
}

// @Summary User logout
// @Tags Authentication api
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "With the bearer started"
// @Success 200 {object} responses.MobileTeamSuccessResponseWithoutData "success response"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Router /api/authentication/logout [post]
func UserLogout(c *gin.Context) {
	authValues, err := stores.GetAuthValue(c)
	if err != nil {
		responses.StatusForbidden(
			c,
			customerror.NewError(err),
		)
		return
	}

	store, err := stores.GetStores(c)
	if err != nil {
		responses.StatusInternalServerError(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	err = store.Authentication.Logout(c.Request.Context(), authValues.UserId)
	if err != nil {
		responses.StatusInternalServerError(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	responses.StatusOk(
		c,
		nil,
		"User logged out successfully",
		"",
	)
}
