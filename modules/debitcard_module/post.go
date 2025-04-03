package debitcardmodule

import (
	"bankapi/requests"
	"bankapi/stores"

	"bitbucket.org/paydoh/paydoh-commons/customerror"
	"bitbucket.org/paydoh/paydoh-commons/responses"
	"github.com/gin-gonic/gin"
)

// @Summary API to Generate Debit Card
// @Tags DebitCard API
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer started"
// @Param X-Device-Ip header string true "With the device IP"
// @Param X-OS header string true "With the OS"
// @Param X-OS-Version header string true "With the OS version"
// @Param X-Lat-Long header string true "With the lat long"
// @Success 200 {object} responses.MobileTeamSuccessResponseWithoutData "success response"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Router /api/debitcard/generate [POST]
func GenerateDebitcard(c *gin.Context) {
	s, err := stores.GetStores(c)
	if err != nil {
		responses.StatusBadRequest(c, customerror.NewError(err), "")
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

	request := requests.NewDebitCardGenerationType()
	if err := request.Validate(requestPayload.Payload); err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	authValue, err := stores.GetAuthValue(c)
	if err != nil {
		responses.StatusBadRequest(c, customerror.NewError(err), "")
		return
	}

	result, err := s.DebitCard.DebitCardGeneration(c.Request.Context(), authValue, request)

	if err != nil {
		responses.StatusBadRequest(c, customerror.NewError(err), "")
		return
	}

	responses.StatusOk(c, result, "DebitCard Generated Successfully", "")
}

// @Summary API to Debit Card Detail
// @Tags DebitCard API
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer started"
// @Param X-Device-Ip header string true "With the device IP"
// @Param X-OS header string true "With the OS"
// @Param X-OS-Version header string true "With the OS version"
// @Param X-Lat-Long header string true "With the lat long"
// @Success 200 {object} responses.MobileTeamSuccessResponse "success response"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Router /api/debitcard/detail [get]
func GetDebitCardDetails(c *gin.Context) {
	s, err := stores.GetStores(c)

	if err != nil {
		responses.StatusInternalServerError(c, customerror.NewError(err), "")
	}

	authValues, err := stores.GetAuthValue(c)
	if err != nil {
		responses.StatusBadRequest(c, customerror.NewError(err), "")
		return
	}

	result, err := s.DebitCard.DebitCardDetail(c.Request.Context(), authValues)

	if err != nil {
		responses.StatusBadRequest(c, customerror.NewError(err), "")
		return
	}

	responses.StatusOk(c, result, "Successfull", "")
}

// @Summary API to Get Transaction Limit
// @Tags DebitCard API
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer started"
// @Param X-Device-Ip header string true "With the device IP"
// @Param X-OS header string true "With the OS"
// @Param X-OS-Version header string true "With the OS version"
// @Param X-Lat-Long header string true "With the lat long"
// @Success 200 {object} responses.MobileTeamSuccessResponse "success response"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Router /api/debitcard/get-limit-list [post]
func GetTransactionLimit(c *gin.Context) {
	s, err := stores.GetStores(c)
	if err != nil {
		responses.StatusBadRequest(c, err.Error(), "")
		return
	}
	requestPayload, err := stores.GetRequestPayload(c)
	if err != nil {
		responses.StatusBadRequest(c, err.Error(), "")
		return
	}

	request := requests.NewGetTransactionLimit()
	if err := request.Validate(requestPayload.Payload); err != nil {
		responses.StatusBadRequest(c, err.Error(), "")
		return
	}
	auth, err := stores.GetAuthValue(c)
	if err != nil {
		responses.StatusBadRequest(c, err.Error(), "")
		return
	}
	response, err := s.DebitCard.GetTransactionLimit(c.Request.Context(), auth, request)
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
		response,
		"Successfully Get Transaction Limit",
		"",
	)
}

// @Summary API to Set Transaction Limit
// @Tags DebitCard API
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer started"
// @Param X-Device-Ip header string true "With the device IP"
// @Param X-OS header string true "With the OS"
// @Param X-OS-Version header string true "With the OS version"
// @Param X-Lat-Long header string true "With the lat long"
// @Success 200 {object} responses.MobileTeamSuccessResponseWithoutData "success response"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Router /api/debitcard/set-txn-limit [post]
func SetTransactonLimit(c *gin.Context) {
	s, err := stores.GetStores(c)

	if err != nil {
		responses.StatusBadRequest(c, customerror.NewError(err), "")
		return
	}

	requestayload, err := stores.GetRequestPayload(c)
	if err != nil {
		responses.StatusBadRequest(c, customerror.NewError(err), "")
		return
	}

	request := requests.NewRequestEditTransaction()

	if err := request.Validation(requestayload.Payload); err != nil {
		responses.StatusBadRequest(c, customerror.NewError(err), "")
		return
	}

	auth, err := stores.GetAuthValue(c)
	if err != nil {
		responses.StatusBadRequest(c, customerror.NewError(err), "")
		return
	}

	result, err := s.DebitCard.SetTransactionLimit(c.Request.Context(), auth, request)
	if err != nil {
		responses.StatusBadRequest(c, customerror.NewError(err), "")
		return
	}

	responses.StatusOk(c, result, "OTP Sent Successfully", "")

}

// @Summary API to Get Card Status
// @Tags DebitCard API
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer started"
// @Param X-Device-Ip header string true "With the device IP"
// @Param X-OS header string true "With the OS"
// @Param X-OS-Version header string true "With the OS version"
// @Param X-Lat-Long header string true "With the lat long"
// @Success 200 {object} responses.MobileTeamSuccessResponse "success response"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Router /api/debitcard/get-card-status [get]
func GetCardStatus(c *gin.Context) {
	s, err := stores.GetStores(c)
	if err != nil {
		responses.StatusBadRequest(c, customerror.NewError(err), "")
		return
	}

	auth, err := stores.GetAuthValue(c)
	if err != nil {
		responses.StatusBadRequest(c, customerror.NewError(err), "")
		return
	}
	response, err := s.DebitCard.GetCardStatus(c.Request.Context(), auth)
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
		response,
		"Status retrieved successfully",
		"",
	)
}

// @Summary API to Set Card Status
// @Tags DebitCard API
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer started"
// @Param X-Device-Ip header string true "With the device IP"
// @Param X-OS header string true "With the OS"
// @Param X-OS-Version header string true "With the OS version"
// @Param X-Lat-Long header string true "With the lat long"
// @Success 200 {object} responses.MobileTeamSuccessResponseWithoutData "success response"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Router /api/debitcard/set-card-status [post]
func SetCardStatus(c *gin.Context) {
	s, err := stores.GetStores(c)
	if err != nil {
		responses.StatusBadRequest(c, customerror.NewError(err), "")
		return
	}
	requestPayload, err := stores.GetRequestPayload(c)
	if err != nil {
		responses.StatusBadRequest(c, customerror.NewError(err), "")
		return
	}

	request := requests.NewSetCardStatus()
	if err := request.Validate(requestPayload.Payload); err != nil {
		responses.StatusBadRequest(c, customerror.NewError(err), "")
		return
	}
	auth, err := stores.GetAuthValue(c)
	if err != nil {
		responses.StatusBadRequest(c, customerror.NewError(err), "")
		return
	}

	response, err := s.DebitCard.SetCardStatus(c.Request.Context(), auth, request)
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
		response,
		"OTP Sent Successfully",
		"",
	)
}

// @Summary API to Set Debit Card Pin
// @Tags DebitCard API
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer started"
// @Param X-Device-Ip header string true "With the device IP"
// @Param X-OS header string true "With the OS"
// @Param X-OS-Version header string true "With the OS version"
// @Param X-Lat-Long header string true "With the lat long"
// @Success 200 {object} responses.MobileTeamSuccessResponseWithoutData "success response"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Router /api/debitcard/set-debitcard-pin [post]
func SetDebitCardPin(c *gin.Context) {
	s, err := stores.GetStores(c)
	if err != nil {
		responses.StatusBadRequest(c, customerror.NewError(err), "")
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
	authValue, err := stores.GetAuthValue(c)
	if err != nil {
		responses.StatusBadRequest(c, customerror.NewError(err), "")
		return
	}
	req := requests.NewSetDebitCardPinReq()
	if err := req.Validate(requestPayload.Payload); err != nil {
		responses.StatusBadRequest(c, customerror.NewError(err), "")
		return
	}
	result, err := s.DebitCard.SetDebitCardPin(c.Request.Context(), authValue, req)
	if err != nil {
		responses.StatusBadRequest(c, customerror.NewError(err), "")
		return
	}
	responses.StatusOk(c, result, "OTP sent Successfully", "")
}

// @Summary API to Verify Debit Card Pin
// @Tags DebitCard API
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer started"
// @Param X-Device-Ip header string true "With the device IP"
// @Param X-OS header string true "With the OS"
// @Param X-OS-Version header string true "With the OS version"
// @Param X-Lat-Long header string true "With the lat long"
// @Success 200 {object} responses.MobileTeamSuccessResponseWithoutData "success response"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Router /api/debitcard/verify-otp [post]
func VerifyOTP(c *gin.Context) {
	s, err := stores.GetStores(c)
	if err != nil {
		responses.StatusBadRequest(c, customerror.NewError(err), "")
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
	authValue, err := stores.GetAuthValue(c)
	if err != nil {
		responses.StatusBadRequest(c, customerror.NewError(err), "")
		return
	}
	req := requests.NewDebitCardVerifyOtpReq()
	if err := req.Validate(requestPayload.Payload); err != nil {
		responses.StatusBadRequest(c, customerror.NewError(err), "")
		return
	}

	result, err := s.DebitCard.VerifyDebitCardOTP(c.Request.Context(), authValue, req)
	if err != nil {
		responses.StatusBadRequest(c, customerror.NewError(err), "")
		return
	}
	responses.StatusOk(c, result, "OTP verified Successfully", "")
}

// @Summary API to Track Debit Card Status
// @Tags DebitCard API
// @Accept json
// @Produce json
// @Param Authorization header string true "With the bearer started"
// @Param X-Device-Ip header string true "With the device IP"
// @Param X-OS header string true "With the OS"
// @Param X-OS-Version header string true "With the OS version"
// @Param X-Lat-Long header string true "With the lat long"
// @Success 200 {object} responses.MobileTeamSuccessResponse "success response"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Router /api/debitcard/track-status [get]
func TrackDebitCardStatus(c *gin.Context) {
	s, err := stores.GetStores(c)
	if err != nil {
		responses.StatusInternalServerError(c, customerror.NewError(err), "")
		return
	}

	authValues, err := stores.GetAuthValue(c)
	if err != nil {
		responses.StatusBadRequest(c, customerror.NewError(err), "")
		return
	}

	result, err := s.DebitCard.TrackDebitCardStatus(c.Request.Context(), authValues)
	if err != nil {
		responses.StatusBadRequest(c, customerror.NewError(err), "")
		return
	}

	responses.StatusOk(c, result, "Debit Card Status Retrieved Successfully", "")
}
