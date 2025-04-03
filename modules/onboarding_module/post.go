package onboardingmodule

import (
	"bitbucket.org/paydoh/paydoh-commons/customerror"
	"bitbucket.org/paydoh/paydoh-commons/responses"
	"github.com/gin-gonic/gin"

	"bankapi/requests"
	"bankapi/stores"
)

// @Summary Api to update personal information
// @Tags Onboarding Apis
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "With the bearer started"
// @Param X-Device-Ip header string true "With the device ip"
// @Param X-OS header string true "With the os"
// @Param X-OS-Version header string true "With the os version"
// @Param X-Lat-Long header string true "With the lat long"
// @Param requestBody body requests.PersonalInformationRequest true "Request Body Data"
// @Param encryptedRequest body requests.EncryptedRequest true "Encrypted Request"
// @Success 200 {object} responses.MobileTeamSuccessResponse "success response"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Router /api/onboarding/personal-information [post]
func UpdatePersonalInformation(c *gin.Context) {
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

	request := requests.NewPersonalInformationRequest()

	if err := request.Validate(requestPayload.Payload); err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	result, err := store.Onboarding.AddPersonalInformation(
		c.Request.Context(),
		authValues,
		request,
	)

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
		"successfully added / updated personal information",
		"",
	)
}

// @Summary Api to Create bank account
// @Tags Onboarding Apis
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "With the bearer started"
// @Param X-Device-Ip header string true "With the device ip"
// @Param X-OS header string true "With the os"
// @Param X-OS-Version header string true "With the os version"
// @Param X-Lat-Long header string true "With the lat long"
// @Param encryptedRequest body requests.EncryptedRequest true "Encrypted Request"
// @Param createAccountRequest body requests.CreateBankAccountRequest false "Create Bank Account Request"
// @Success 200 {object} responses.MobileTeamSuccessResponse "success response"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Router /api/onboarding/create-account [post]
func CreateAccount(c *gin.Context) {
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
		responses.StatusBadRequest(c,
			customerror.NewError(err),
			"",
		)
		return
	}

	request := requests.NewCreateBankAccountRequest()

	if err := request.Validate(requestPayload.Payload); err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	result, err := s.Onboarding.CreateBankAccount(c.Request.Context(), authValues, request)

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
		"successfully created bank account",
		"",
	)
}

// @Summary Api to Create bank account
// @Tags Onboarding Apis
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "With the bearer started"
// @Param X-Device-Ip header string true "With the device ip"
// @Param X-OS header string true "With the os"
// @Param X-OS-Version header string true "With the os version"
// @Param X-Lat-Long header string true "With the lat long"
// @Param encryptedRequest body requests.EncryptedRequest true "Encrypted Request"
// @Param createAccountRequest body requests.CreateBankAccountRequest false "Create Bank Account Request"
// @Success 200 {object} responses.MobileTeamSuccessResponse "success response"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Router /api/onboarding/get-account-details [post]
func GetAccountDetails(c *gin.Context) {
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

	result, err := s.Onboarding.GetAccountDetails(c.Request.Context(), authValues)
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
		"successfully fetch account details",
		"",
	)
}

// @Summary Api to get pincode details
// @Tags Onboarding Apis
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
// @Router /api/onboarding/pincode-details [POST]
func GetPincodeDetails(c *gin.Context) {
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
		responses.StatusBadRequest(c,
			customerror.NewError(err),
			"",
		)
		return
	}

	request := requests.NewPincodeDetails()
	if err := request.Validate(requestPayload.Payload); err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	result, err := s.Onboarding.GetPincodeDetails(c.Request.Context(), authValues, request)
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
		"successfully fetch user pincode details",
		"",
	)
}
