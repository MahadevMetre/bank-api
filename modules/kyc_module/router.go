package kyc_module

import (
	"bankapi/middleware"

	"github.com/gin-gonic/gin"
)

func Routes(app *gin.RouterGroup) {
	api := app.Group("/kyc")
	api.Use(middleware.AuthMiddleware())
	api.Use(middleware.DecryptMiddleware())
	{
		api.POST("/consent", UpdateKycConsent)
		api.GET("/consent", GetKycConsent)
		api.GET("/vcip-url", GetVcipUrl)
		api.GET("/get-update", GetKycUpdate)
	}
}
