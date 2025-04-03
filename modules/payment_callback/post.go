package payment_callback

import (
	"bankapi/requests"
	"bankapi/stores"
	"fmt"

	"bitbucket.org/paydoh/paydoh-commons/customerror"
	"bitbucket.org/paydoh/paydoh-commons/responses"
	"github.com/gin-gonic/gin"
)

// @Summary Account create callback API.
// @Tags CallBack API
// @Accept  json
// @Produce  json
// @Param Authorization header string true "API  key"
//	@Param		user	body	requests.PaymentCallbackRequestData	true	"request data"
// @Success 200 {object} responses.MobileTeamSuccessResponseWithoutData "success response"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Router /bank/callback/normal-payment [post]
func PaymentCallBackAPI(c *gin.Context) {
	request := requests.NewPaymentCallbackRequestData()

	if err := request.Validate(c); err != nil {
		fmt.Println("ERRR REQ ", err)
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	// store
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

	err = store.Payment.Update(c.Request.Context(), request)
	if err != nil {
		fmt.Println("ERRR ", err)
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
