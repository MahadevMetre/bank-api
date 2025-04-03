package openmodules

import (
	"bitbucket.org/paydoh/paydoh-commons/customerror"
	"bitbucket.org/paydoh/paydoh-commons/responses"
	"github.com/gin-gonic/gin"

	"bankapi/requests"
	"bankapi/stores"
)

// @Summary Api to add shipping address
// @Tags Open Bank apis
// @Accept  x-www-form-urlencoded
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "With the bearer started"
// @Param X-Device-Ip header string true "With the device ip"
// @Param X-OS header string true "With the os"
// @Param X-OS-Version header string true "With the os version"
// @Param documet_type formData string true "Document type"
// @Param document formData file true "Document"
// @Param address_line_1 formData string true "Address line 1"
// @Param street_name formData string true "Street name"
// @Param locality formData string true "Locality"
// @Param landmark formData string false "Landmark"
// @Param city formData string true "City"
// @Param state formData string true "State"
// @Param pin_code formData string true "Pin code"
// @Param country formData string true "Country"
// @Success 200 {object} responses.MobileTeamSuccessResponseWithoutData "success response"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Router /api/open/shipping-address [post]
func UpdateShippingAddress(c *gin.Context) {
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

	awsInstance, err := stores.GetAwsInstance(c)

	if err != nil {
		responses.StatusInternalServerError(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	request := requests.NewAddShippingAddress()

	if err := request.Validate(c, awsInstance); err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	if err := s.Open.UpdateShippingAddress(c.Request.Context(), authValues, request); err != nil {
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
		"successfully updated shipping address",
		"",
	)
}

// @Summary Api to update shipping address
// @Tags Open Bank apis
// @Accept  x-www-form-urlencoded
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "With the bearer started"
// @Param X-Device-Ip header string true "With the device ip"
// @Param X-OS header string true "With the os"
// @Param X-OS-Version header string true "With the os version"
// @Param documet_type formData string false "Document type"
// @Param document formData file false "Document"
// @Param address_line_1 formData string false "Address line 1"
// @Param street_name formData string false "Street name"
// @Param locality formData string false "Locality"
// @Param landmark formData string false "Landmark"
// @Param city formData string false "City"
// @Param state formData string false "State"
// @Param pin_code formData string false "Pin code"
// @Param country formData string false "Country"
// @Success 200 {object} responses.MobileTeamSuccessResponseWithoutData "success response"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Router /api/open/update-shipping-address [post]
func ShippingAddressUpdate(c *gin.Context) {
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

	awsInstance, err := stores.GetAwsInstance(c)

	if err != nil {
		responses.StatusInternalServerError(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	request := requests.NewUpdateShippingAddress()

	if err := request.Validate(c, awsInstance); err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	if err := s.Open.ShippingAddressUpdate(c.Request.Context(), authValues, request); err != nil {
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
		"successfully updated shipping address",
		"",
	)
}

// @Summary Api to get receipt id
// @Tags Open Bank apis
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "With the bearer started"
// @Param X-Device-Ip header string true "With the device ip"
// @Param X-OS header string true "With the os"
// @Param X-OS-Version header string true "With the os version"
// // @Param getReceiptId body requests.GetReceiptID true "Get receipt id Request"
// @Success 200 {object} responses.MobileTeamSuccessResponse "success response"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Router /api/open/get-receipt-id [post]
func GetReceiptId(c *gin.Context) {
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

	requestPayload, err := stores.GetRequestPayload(c)

	if err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	request := requests.NewGatewayReceiptIdRequest()

	if err := request.Validate(requestPayload.Payload); err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	result, err := s.Open.GetReceiptID(c.Request.Context(), authValues, request)

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
		"successfully got receipt id",
		"",
	)

}

// @Summary Api to add payment status
// @Tags Open Bank apis
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "With the bearer started"
// @Param X-Device-Ip header string true "With the device ip"
// @Param X-OS header string true "With the os"
// @Param X-OS-Version header string true "With the os version"
// @Param addPaymentStatus body requests.AddPaymentStatusRequest true "Add payment status Request"
// @Success 200 {object} responses.MobileTeamSuccessResponse "success response"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Router /api/open/payment-status [post]
func UpdatePaymentStatus(c *gin.Context) {
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

	requestPayload, err := stores.GetRequestPayload(c)

	if err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	request := requests.NewAddPaymentStatusRequest()

	if err := request.Validate(requestPayload.Payload); err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	result, err := s.Open.UpdatePaymentStatus(c.Request.Context(), authValues, request)
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
		"successfully updated payment status",
		"",
	)
}

// @Summary Api to set mpin
// @Tags Open Bank apis
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "With the bearer started"
// @Param X-Device-Ip header string true "With the device ip"
// @Param X-OS header string true "With the os"
// @Param X-OS-Version header string true "With the os version"
// @Param setMpinRequest body requests.MpinRequest true "Add mpin request"
// @Success 200 {object} responses.MobileTeamSuccessResponseWithoutData "success response"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Router /api/open/set-mpin [post]
func SetMpin(c *gin.Context) {
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

	requestPayload, err := stores.GetRequestPayload(c)

	if err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			authValues.Key,
		)
		return
	}

	request := requests.NewMpinRequest()

	if err := request.Validate(requestPayload.Payload); err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			authValues.Key,
		)
		return
	}

	if err := s.Open.SetMpin(c.Request.Context(), request, authValues); err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			authValues.Key,
		)
		return
	}

	responses.StatusOk(
		c,
		nil,
		"successfully set mpin",
		authValues.Key,
	)
}

// @Summary Api to verify mpin
// @Tags Open Bank apis
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "With the bearer started"
// @Param X-Device-Ip header string true "With the device ip"
// @Param X-OS header string true "With the os"
// @Param X-OS-Version header string true "With the os version"
// @Param setMpinRequest body requests.MpinRequest true "Add mpin request"
// @Success 200 {object} responses.MobileTeamSuccessResponseWithoutData "success response"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Router /api/open/verify-mpin [post]
func VerifyMpin(c *gin.Context) {
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

	requestPayload, err := stores.GetRequestPayload(c)

	if err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			authValues.Key,
		)
		return
	}

	request := requests.NewMpinRequest()

	if err := request.Validate(requestPayload.Payload); err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			authValues.Key,
		)
		return
	}

	if err := s.Open.VerifyMpin(c.Request.Context(), request, authValues); err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			authValues.Key,
		)
		return
	}

	responses.StatusOk(
		c,
		nil,
		"successfully verified mpin",
		authValues.Key,
	)
}

// @Summary Api to reset mpin
// @Tags Open Bank apis
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "With the bearer started"
// @Param X-Device-Ip header string true "With the device ip"
// @Param X-OS header string true "With the os"
// @Param X-OS-Version header string true "With the os version"
// @Param resetMpinRequest body requests.ResetMpinRequest true "Add reset-mpin request"
// @Success 200 {object} responses.MobileTeamSuccessResponseWithoutData "success response"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Router /api/open/reset-mpin [post]
func ReSetMpin(c *gin.Context) {
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

	requestPayload, err := stores.GetRequestPayload(c)
	if err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			authValues.Key,
		)
		return
	}

	request := requests.NewResetMpinRequest()
	if err := request.Validate(requestPayload.Payload); err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			authValues.Key,
		)
		return
	}

	if err := s.Open.ResetMpin(c.Request.Context(), request, authValues); err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			authValues.Key,
		)
		return
	}

	responses.StatusOk(
		c,
		nil,
		"successfully reset mpin",
		authValues.Key,
	)
}

// @Summary API to verify reset forgotten M-PIN
// @Tags Open Bank APIs
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "With the bearer started"
// @Param X-Device-Ip header string true "With the device ip"
// @Param X-OS header string true "With the os"
// @Param X-OS-Version header string true "With the os version"
// @Param X-Device-Ip header string true "With the device IP"
// @Param X-OS header string true "With the OS"
// @Param X-OS-Version header string true "With the OS version"
// @Param forgotMpinRequest body requests.ForgotMpinResetRequest true "Add forgotten M-PIN reset request"
// @Success 200 {object} responses.MobileTeamSuccessResponseWithoutData "success response"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Router /api/open/verify-forgot-mpin [post]
func verifyForgetMpinReset(c *gin.Context) {
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

	requestPayload, err := stores.GetRequestPayload(c)
	if err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			authValues.Key,
		)
		return
	}

	request := requests.NewForgotMpinResetRequest()
	if err := request.Validate(requestPayload.Payload); err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			authValues.Key,
		)
		return
	}

	if err := s.Open.VerifyForgottenMpinResetRequest(c.Request.Context(), request, authValues); err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			authValues.Key,
		)
		return
	}

	responses.StatusOk(
		c,
		nil,
		"User details have been successfully verified for resetting the forgotten M-PIN",
		authValues.Key,
	)
}

// @Summary API to reset forgotten M-PIN
// @Tags Open Bank APIs
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "With the bearer started"
// @Param X-Device-Ip header string true "With the device IP"
// @Param X-OS header string true "With the OS"
// @Param X-OS-Version header string true "With the OS version"
// @Param UpdateMpinRequest body requests.UpdateMpinRequest true "Update M-PIN after reset request"
// @Success 200 {object} responses.MobileTeamSuccessResponseWithoutData "success response"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Router /api/open/update-mpin [post]
func UpdateForgetMpinReset(c *gin.Context) {
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

	requestPayload, err := stores.GetRequestPayload(c)
	if err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			authValues.Key,
		)
		return
	}

	request := requests.NewUpdateMpinRequest()
	if err := request.Validate(requestPayload.Payload); err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			authValues.Key,
		)
		return
	}

	if err := s.Open.UpdateMpin(c.Request.Context(), request, authValues); err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			authValues.Key,
		)
		return
	}

	responses.StatusOk(
		c,
		nil,
		"Successfully updated Mpin for user after forgotten M-PIN reset",
		authValues.Key,
	)
}

func UpdateFcmToken(c *gin.Context) {
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

	request := requests.NewUpdateFcmToken()
	if err := request.Validate(c); err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	_, err = s.Open.UpdateFcmToken(c.Request.Context(), request, authValues)
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
		nil,
		"Successfully updated FCM token",
		"",
	)
}
