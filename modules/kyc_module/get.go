package kyc_module

import (
	"bitbucket.org/paydoh/paydoh-commons/customerror"
	"bitbucket.org/paydoh/paydoh-commons/responses"
	"github.com/gin-gonic/gin"

	"bankapi/stores"
)

// @Summary Api to get kyc consent
// @Tags KYC Apis
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
// @Router /api/kyc/consent [get]
func GetKycConsent(c *gin.Context) {
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

	result, err := store.Kyc.GetKycConsent(c.Request.Context(), authValues)

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
		"successfully retrieved kyc consent",
		"",
	)
}

// @Summary Api to get vcip url
// @Tags KYC Apis
// @Accept  json
// @Produce html
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "With the bearer started"
// @Param X-Device-Ip header string true "With the device ip"
// @Param X-OS header string true "With the os"
// @Param X-OS-Version header string true "With the os version"
// @Param X-Lat-Long header string true "With the lat long"
// @Success 301 "Moved Permanently"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Success 200 {object} responses.MobileTeamSuccessResponse "success response"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Failure 401 "Unauthorized"
// @Router /api/kyc/vcip-url [get]
func GetVcipUrl(c *gin.Context) {
	s, err := stores.GetStores(c)

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

	result, err := s.Kyc.GetVcipUrl(c.Request.Context(), authValues)

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
		"successfully retrieved kyc data",
		"",
	)

	// responses.StatusMovedPermanently(
	// 	c,
	// 	"",
	// 	result,
	// )
}

// @Summary Api to get vcip url
// @Tags KYC Apis
// @Accept  json
// @Produce html
// @Produce json
// @Security ApiKeyAuth
// @Param Authorization header string true "With the bearer started"
// @Param X-Device-Ip header string true "With the device ip"
// @Param X-OS header string true "With the os"
// @Param X-OS-Version header string true "With the os version"
// @Param X-Lat-Long header string true "With the lat long"
// @Success 301 "Moved Permanently"
// @Success 200 {object} responses.MobileTeamSuccessResponse "success response"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Failure 401 "Unauthorized"
// @Router /api/kyc/get-update [get]
func GetKycUpdate(c *gin.Context) {
	s, err := stores.GetStores(c)
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

	result, err := s.Kyc.GetKycUpdateData(c.Request.Context(), authValues)
	if err != nil {
		responses.StatusBadRequestV3(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	responses.StatusOk(
		c,
		result,
		"successfully retrieved kyc update data",
		"",
	)
}
