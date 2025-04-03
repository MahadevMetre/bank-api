package transaction_module

import (
	"bankapi/requests"
	"bankapi/stores"

	"bitbucket.org/paydoh/paydoh-commons/customerror"
	"bitbucket.org/paydoh/paydoh-commons/responses"
	"github.com/gin-gonic/gin"
)

// @Summary Api to get the transactions.
// @Description Api to get the transactions.
// @Tags transaction api
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "With the bearer started"
// @Param X-Device-Ip header string true "With the device ip"
// @Param X-OS header string true "With the os"
// @Param X-OS-Version header string true "With the os version"
// @Param X-Lat-Long header string true "With the lat long"
// @Param encryptedRequest body requests.EncryptedRequest true "Encrypted Request"
// @Success 200 {object} responses.MobileTeamSuccessResponse "success response"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Router /api/transaction/history [post]
func Transactions(c *gin.Context) {
	authValues, err := stores.GetAuthValue(c)
	if err != nil {
		responses.StatusForbidden(
			c,
			customerror.NewError(err),
		)
		return
	}

	// store
	store, err := stores.GetStores(c)
	if err != nil {
		responses.StatusInternalServerError(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	payload, err := stores.GetRequestPayload(c)
	if err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	request := requests.NewTransactionRequest()
	if err := request.ValidateEncrypted(payload.Payload); err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	request.UserId = authValues.UserId
	responseData, err := store.TransactionHistory.GetTransactionData(c.Request.Context(), *request)
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
		responseData,
		"successfully fetched transaction history.",
		"",
	)
}

// @Summary Api to get the transactions.
// @Description Api to get the transactions.
// @Tags transaction api
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "With the bearer started"
// @Param X-Device-Ip header string true "With the device ip"
// @Param X-OS header string true "With the os"
// @Param X-OS-Version header string true "With the os version"
// @Param X-Lat-Long header string true "With the lat long"
// @Param encryptedRequest body requests.EncryptedRequest true "Encrypted Request"
// @Success 200 {object} responses.MobileTeamSuccessResponse "success response"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Router /api/transaction/history/internal [Post]
func GetInternalTransactions(c *gin.Context) {
	store, err := stores.GetStores(c)
	if err != nil {
		responses.StatusInternalServerError(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	request := requests.NewTransactionRequest()
	if err := request.Validate(c); err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	responseData, err := store.TransactionHistory.GetTransactionData(c.Request.Context(), *request)
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
		responseData,
		"successfully fetched transaction history.",
		"",
	)
}

// @Summary Api to get the transactions.
// @Description Api to get the transactions.
// @Tags transaction api
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "With the bearer started"
// @Param X-Device-Ip header string true "With the device ip"
// @Param X-OS header string true "With the os"
// @Param X-OS-Version header string true "With the os version"
// @Param X-Lat-Long header string true "With the lat long"
// @Param encryptedRequest body requests.EncryptedRequest true "Encrypted Request"
// @Success 200 {object} responses.MobileTeamSuccessResponse "success response"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Router /api/transaction/history/get-details [post]
func GetTxnDetails(c *gin.Context) {
	authValues, err := stores.GetAuthValue(c)
	if err != nil {
		responses.StatusForbidden(
			c,
			customerror.NewError(err),
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

	payload, err := stores.GetRequestPayload(c)
	if err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	req := new(requests.TransactionDetailsRequest)
	if err := req.ValidateEncrypted(payload.Payload); err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	responseData, err := store.TransactionHistory.GetTransactionDetails(req, authValues.UserId)
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
		responseData,
		"successfully fetched transaction details.",
		"",
	)
}
