package account_create_callback

import (
	"bankapi/constants"
	"bankapi/requests"
	"bankapi/stores"
	"bankapi/utils"
	"encoding/json"
	"net/http"

	"bitbucket.org/paydoh/paydoh-commons/customerror"
	"bitbucket.org/paydoh/paydoh-commons/responses"
	"github.com/gin-gonic/gin"
)

// @Summary Account create callback API.
// @Tags CallBack API
// @Accept  json
// @Produce  json
// @Param Authorization header string true "API  key"
//	@Param		user	body	requests.AccountCreateCallbackRequest	true	"request data"
// @Success 200 {object} responses.MobileTeamSuccessResponseWithoutData "success response"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Router /bank/callback/account-create [post]
func AccountCreateCallback(c *gin.Context) {
	encryptedReq := new(requests.AccountCreateCallbackEncryptedRequest)

	if err := c.ShouldBindJSON(encryptedReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	respData, err := utils.DecryptResponse(encryptedReq.EncryptedRes, constants.BankEncryptionKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	request := requests.NewAccountCreateCallbackRequest()
	if err := json.Unmarshal(respData, request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	if err := request.Validate(c); err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
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

	err = store.Onboarding.UpdateBankAccount(c.Request.Context(), request, encryptedReq)
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
		"success",
		"",
	)
}
