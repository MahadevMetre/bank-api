package sms_callback

import (
	"bankapi/stores"

	"bitbucket.org/paydoh/paydoh-commons/customerror"
	"bitbucket.org/paydoh/paydoh-commons/responses"
	"github.com/gin-gonic/gin"
)

// @Summary Api to verify sms
// @Tags Sms API
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "With the bearer started"
// @Param X-Device-Ip header string true "With the device ip"
// @Param X-OS header string true "With the os"
// @Param X-OS-Version header string true "With the os version"
// @Param X-Lat-Long header string true "With the lat long"
// @Param sender query string true "With the sender"
// @Param message query string true "With the message"
// @Param operator query string true "With the operator"
// @Success 200 {object} responses.MobileTeamSuccessResponseWithoutData "success response"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Router /api/sms-received [get]
func SmsReceivedCallback(c *gin.Context) {
	sender := c.Query("sender")
	message := c.Query("message")
	simOperator := c.Query("operator")

	store, err := stores.GetStores(c)
	if err != nil {
		responses.StatusInternalServerError(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	if err := store.Authorization.SmsVerification(sender, message, simOperator); err != nil {
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
