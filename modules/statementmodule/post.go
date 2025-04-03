package statementmodule

import (
	"net/http"

	"bitbucket.org/paydoh/paydoh-commons/customerror"
	"bitbucket.org/paydoh/paydoh-commons/responses"
	"github.com/gin-gonic/gin"

	"bankapi/requests"
	"bankapi/stores"
)

// @Summary Api to get bank statement
// @Tags Account statement Api
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param Authorization header string true "With the bearer started"
// @Param X-Device-Ip header string true "With the device ip"
// @Param X-OS header string true "With the os"
// @Param X-OS-Version header string true "With the os version"
// @Param X-Lat-Long header string true "With the lat long"
// @Param encryptedRequest body requests.EncryptedRequest true "Encrypted Request"
// @Success 200 "ok"
// @Router /api/statement/get-statement [post]
func GetStatement(c *gin.Context) {
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

	payload, err := stores.GetRequestPayload(c)
	if err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	request := requests.NewStatementRequest()
	if err := request.Validate(payload.Payload); err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	result, err := store.Statement.GetAccountStatement(c.Request.Context(), authValues, request)
	if err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	if request.Type == "pdf" {
		c.Header("Content-Type", "application/pdf")
		c.Header("Content-Disposition", "attachment; filename=bank_statement.pdf")
	} else if request.Type == "csv" {
		c.Header("Content-Type", "text/csv")
		c.Header("Content-Disposition", "attachment; filename=bank_statement.csv")
	}

	byteData, ok := result.([]byte)
	if !ok {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	c.Data(http.StatusOK, "application/pdf", byteData)

}
