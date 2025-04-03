package kyc_audit_data

import (
	"bankapi/requests"
	"bankapi/stores"
	"fmt"

	"bitbucket.org/paydoh/paydoh-commons/customerror"
	"bitbucket.org/paydoh/paydoh-commons/responses"
	"github.com/gin-gonic/gin"
)

// @Summary Kyc audit complition callback API.
// @Tags CallBack API
// @Accept  json
// @Produce  json
// @Param Authorization header string true "API  key"
//	@Param		user	body	requests.KycAuditRequestData	true	"request data"
// @Success 200 {object} responses.MobileTeamSuccessResponseWithoutData "success response"
// @Failure 400 {object} responses.MobileTeamErrorResponse "Error response for Bad Request"
// @Failure 500 {object} responses.MobileTeamErrorResponse "Error response for Internal Server Error"
// @Router /bank/callback/audit-complition [post]
func KycAuditComplition(c *gin.Context) {
	request := requests.NewwKycAuditRequestData()

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

	err = store.KycAuditStore.Create(c.Request.Context(), request)
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
		"success",
		"",
	)
}
