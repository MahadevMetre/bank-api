package kyc_update_data_callback

import (
	"bankapi/stores"
	"html/template"
	"net/http"

	"bitbucket.org/paydoh/paydoh-commons/customerror"
	"bitbucket.org/paydoh/paydoh-commons/responses"
	"github.com/gin-gonic/gin"
)

// @Summary Kyc audit update callback API.
// @Tags CallBack API
// @Accept  json
// @Produce  json
// @Param Authorization header string true "API  key"
// @Router /bank/callback/kyc-update [get]
func KycUpdateData(c *gin.Context) {
	applicantId := c.Query("vkyccompletionApplicantId")
	status := c.Query("Status")
	astat := c.Query("AStat")
	acom := c.Query("ACom")

	s, err := stores.GetStores(c)

	if err != nil {
		responses.StatusInternalServerError(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	err = s.Kyc.AddKycUpdateData(c.Request.Context(), applicantId, status, acom, astat)
	if err != nil {
		responses.StatusBadRequest(
			c,
			customerror.NewError(err),
			"",
		)
		return
	}

	tmpl, err := template.ParseFiles("templates/loader.html")
	if err != nil {
		c.String(http.StatusInternalServerError, "Error loading template")
		return
	}

	c.Status(http.StatusOK)

	err = tmpl.Execute(c.Writer, gin.H{})
	if err != nil {
		c.String(http.StatusInternalServerError, "Error executing template")
		return
	}

}
