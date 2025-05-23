package addressupdate

import (
	"bankapi/stores"

	"bitbucket.org/paydoh/paydoh-commons/customerror"
	"bitbucket.org/paydoh/paydoh-commons/responses"
	"github.com/gin-gonic/gin"
)

// GetAddressUpdateStatus Get address update status
// @Tags address apis
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
// @Router /api/address/get-status [GET]
func GetAddressUpdateStatus(c *gin.Context) {
	s, err := stores.GetStores(c)
	if err != nil {
		responses.StatusBadRequest(c, err.Error(), "")
	}

	authValues, err := stores.GetAuthValue(c)
	if err != nil {
		responses.StatusUnauthorized(
			c,
			customerror.NewError(err),
		)
		return
	}

	result, err := s.AddressStore.GetAddressUpdateStatus(authValues)
	if err != nil {
		responses.StatusUnauthorized(
			c,
			customerror.NewError(err),
		)
		return
	}

	responses.StatusOk(c, result, "Address update status get successfully", "")
}
