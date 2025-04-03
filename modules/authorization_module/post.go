package authorizationmodule

import (
	"bitbucket.org/paydoh/paydoh-commons/customerror"
	"bitbucket.org/paydoh/paydoh-commons/responses"
	"github.com/gin-gonic/gin"

	"bankapi/requests"
	"bankapi/stores"
)

// @Summary Api get authorization Token for all the api
// @Description Using the mobile_number get the authorization token for the other apis.
// @Tags authorization apis
// @Accept  json
// @Produce  json
// @Param authorizationRequest body requests.AuthorizationRequest true "Authorization Request"
// @Success 200 {object} responses.MobileTeamSuccessResponse "success response"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Router /api/authorization/ [post]
func Authorization(c *gin.Context) {
	stores, err := stores.GetStores(c)

	if err != nil {
		responses.StatusInternalServerError(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	request := requests.NewAuthorizationRequest()

	if err := request.Validate(c); err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	if request.MobileNumber != "" {
		result, err := stores.Authorization.GetAuthorizationToken(c.Request.Context(), request)

		if err != nil {
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
			"successfully sent token to user",
			"",
		)
	}

	if request.UserId != "" {
		result, err := stores.Authorization.GetAuthorizationTokenByUserId(c.Request.Context(), request)

		if err != nil {
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
			"successfully sent token to user",
			"",
		)
	}
}
