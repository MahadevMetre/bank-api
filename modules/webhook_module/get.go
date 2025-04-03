package webhookmodule

import (
	"bitbucket.org/paydoh/paydoh-commons/customerror"
	"bitbucket.org/paydoh/paydoh-commons/responses"
	"github.com/gin-gonic/gin"

	"bankapi/requests"
	"bankapi/stores"
)

// @Summary Web hook for getting api senders mobile data from Route mobile.
// @Tags Route Mobile Web hook
// @Accept  json
// @Produce  json
// @Param   sender     query    string     true        "Sender"
// @Param message query string true "Message"
// @Param operator query string false "Operator"
// @Param circle query string false "Circle"
// @Success 200 {object} responses.MobileTeamSuccessResponse "success response"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Router /api/webhook/route-mobile [get]
func GetRouteMobileData(c *gin.Context) {

	store, err := stores.GetStores(c)

	if err != nil {
		responses.StatusInternalServerError(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	request := requests.NewWebhookRequest()

	if err := request.Validate(c); err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	result, err := store.Webhook.GetWebhookData(c.Request.Context(), request)

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
		"Successfully recived message",
		"",
	)
}

func GetVcipData(c *gin.Context) {

}
