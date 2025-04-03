package consentmodule

import (
	"bitbucket.org/paydoh/paydoh-commons/customerror"
	"bitbucket.org/paydoh/paydoh-commons/responses"
	"github.com/gin-gonic/gin"

	"bankapi/requests"
	"bankapi/stores"
)

// @Summary Api to Update consent
// @Tags Consent Api
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "With the bearer started"
// @Param X-Device-Ip header string true "With the device ip"
// @Param X-OS header string true "With the os"
// @Param X-OS-Version header string true "With the os version"
// @Param X-Lat-Long header string true "With the lat long"
// @Param encryptedRequest body requests.EncryptedRequest true "Encrypted Request"
// @Success 200 {object} responses.MobileTeamSuccessResponseWithoutData "success response"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Router /api/consent/update [post]
func UpdateConsent(c *gin.Context) {
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

	request := requests.NewConsentRequest()

	if err := request.DecryptAndValidatePayload(requestPayload.Payload); err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	if err := store.Consent.UpdateConsent(c.Request.Context(), request, authValues); err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	responses.StatusOk(
		c,
		nil,
		"Consent updated successfully",
		"",
	)
}
