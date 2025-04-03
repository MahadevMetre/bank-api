package kyc_update_data_callback

import (
	"github.com/gin-gonic/gin"
)

func Routes(app *gin.RouterGroup) {
	app.GET("/kyc-update", KycUpdateData)
}
