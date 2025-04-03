package authorizationmodule

import (
	"bitbucket.org/paydoh/paydoh-commons/customerror"
	"bitbucket.org/paydoh/paydoh-commons/responses"
	"github.com/gin-gonic/gin"

	"bankapi/requests"
	"bankapi/stores"
)

// GetCurrentUserByUserId Get current user by user id
// @Tags authorization apis
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
// @Router /api/authenticated/current-user [POST]
func GetCurrentUserByUserId(c *gin.Context) {

	store, err := stores.GetStores(c)

	if err != nil {
		responses.StatusInternalServerError(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	request := requests.NewUserRequest()

	if err := request.Validate(c); err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	user, err := store.Authentication.GetCurrentUser(c.Request.Context(), request.UserId)

	if err != nil {
		responses.StatusUnauthorized(
			c,
			customerror.NewError(err),
		)
		return
	}

	responses.StatusOk(
		c,
		user,
		"successfully fetched user details",
		"",
	)
}

// GetCurrentUserMobileNumber Get current user by mobile number
// @Tags authorization apis
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
// @Router /api/authenticated/current-user-mobile [POST]
func GetCurrentUserMobileNumber(c *gin.Context) {

	store, err := stores.GetStores(c)

	if err != nil {
		responses.StatusInternalServerError(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	request := requests.NewUserMobileNumber()

	if err := request.Validate(c); err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	user, err := store.Authentication.GetCurrentUserByMobileNumber(c.Request.Context(), request.MobileNumber)

	if err != nil {
		responses.StatusUnauthorized(
			c,
			customerror.NewError(err),
		)
		return
	}

	responses.StatusOk(
		c,
		user,
		"successfully fetched user details",
		"",
	)
}
