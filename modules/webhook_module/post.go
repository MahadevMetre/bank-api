package webhookmodule

import (
	"bitbucket.org/paydoh/paydoh-commons/customerror"
	"bitbucket.org/paydoh/paydoh-commons/responses"
	"github.com/gin-gonic/gin"

	"bankapi/requests"
	"bankapi/stores"
)

// @Summary Web hook for updating vcip data from onboarding
// @Tags Webhook Apis
// @Accept  json
// @Produce  json
// @Param vcipWebhookResquest body requests.IncomingVcipData true "Vcip data Request"
// @Success 200 "ok"
// @Router /api/webhook/kvb/vcip [post]
func UpdateVcipData(c *gin.Context) {

	s, err := stores.GetStores(c)

	if err != nil {
		responses.StatusInternalServerError(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	request := requests.NewIncomingVcipData()

	if err := request.Validate(c); err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	if err := s.Webhook.UpdateVcipData(c.Request.Context(), request); err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	responses.StatusCreated(
		c,
		nil,
		"Successfully updated vcip data",
		"",
	)
}

// @Summary Web hook for updating vcip data from onboarding
// @Tags Webhook Apis
// @Accept  json
// @Produce  json
// @Success 200 "ok"
// @Router /api/webhook/kvb/kvb/kyc [post]
func KycUpdateState(c *gin.Context) {
	s, err := stores.GetStores(c)

	if err != nil {
		responses.StatusInternalServerError(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	request := requests.NewKycUpdateRequest()

	if err := request.Validate(c); err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	if err := s.Webhook.UpdateKycUpdateData(c.Request.Context(), request); err != nil {
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
		"Successfully updated vcip data",
		"",
	)
}

// @Summary Web hook for updating vcip data from onboarding
// @Tags Webhook Apis
// @Accept  json
// @Produce  json
// @Param vcipWebhookResquest body requests.RewardsRequest true "Vcip data Request"
// @Success 200 {object} responses.MobileTeamSuccessResponseWithoutData "success response"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Router /api/webhook/rewards-point [post]
func ProvideRewardPoint(c *gin.Context) {

	s, err := stores.GetStores(c)

	if err != nil {
		responses.StatusInternalServerError(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	request := requests.NewRewardsRequest()

	if err := request.Validate(c); err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	if err := s.Webhook.ProvideRewardsPoint(c.Request.Context(), request); err != nil {
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
		"Success",
		"",
	)
}
