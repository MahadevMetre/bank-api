package consentmodule

import (
	"bankapi/middleware"

	"github.com/gin-gonic/gin"
)

func Routes(app *gin.RouterGroup) {
	api := app.Group("/consent")
	api.Use(middleware.AuthMiddleware())
	api.Use(middleware.DecryptMiddleware())

	{
		api.POST("/update", UpdateConsent)
	}
}
