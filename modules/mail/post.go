package mail

import (
	"bankapi/requests"
	"bankapi/stores"
	"errors"
	"net/http"
	"strings"

	"bitbucket.org/paydoh/paydoh-commons/customerror"
	"bitbucket.org/paydoh/paydoh-commons/responses"
	"github.com/gin-gonic/gin"
)

// @Summary Api to send verification email
// @Tags Mail API
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "With the bearer started"
// @Param X-Device-Ip header string true "With the device ip"
// @Param X-OS header string true "With the os"
// @Param X-OS-Version header string true "With the os version"
// @Param X-Lat-Long header string true "With the lat long"
// @Param request body requests.VerifyEmailReq true "Send verification email"
// @Success 200 {object} responses.MobileTeamSuccessResponseWithoutData "success response"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Router /api/email/sendverification [post]
func SendVerificationEmail(c *gin.Context) {
	store, err := stores.GetStores(c)

	if err != nil {
		responses.StatusInternalServerError(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	authValues, err := stores.GetAuthValue(c)

	if err != nil {
		responses.StatusUnauthorized(
			c,
			customerror.NewError(err),
		)
		return
	}

	requestPayload, err := stores.GetRequestPayload(c)

	if err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	request := requests.NewVerifyEmailReq()

	if err := request.Validate(requestPayload.Payload); err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	if !requests.ValidateDomain(request.EmailId) {
		responses.StatusBadRequest(
			c,
			customerror.NewError(errors.New("Invalid Email ID.")),
			"",
		)
		return
	}

	result, err := store.EmailStore.SendVerificationEmail(c.Request.Context(), request, authValues)

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
		result,
		"verification mail sent successfully",
		"",
	)
}

// @Summary Api to verify email
// @Tags Mail API
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "With the bearer started"
// @Param X-Device-Ip header string true "With the device ip"
// @Param X-OS header string true "With the os"
// @Param X-OS-Version header string true "With the os version"
// @Param X-Lat-Long header string true "With the lat long"
// @Success 200 {object} responses.MobileTeamSuccessResponseWithoutData "success response"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Router /api/email/update-verify/:id [get]
func VerifyEmail(c *gin.Context) {
	store, err := stores.GetStores(c)

	if err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	path := strings.Split(c.Request.URL.Path, "/")
	id := path[len(path)-1]

	_, err = store.EmailStore.VerifyEmail(c.Request.Context(), id)

	if err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	c.HTML(http.StatusOK, "emailverified.html", nil)
}
