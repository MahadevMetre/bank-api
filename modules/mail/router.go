package mail

import (
	"bankapi/middleware"

	"github.com/gin-gonic/gin"
)

func Routes(app *gin.RouterGroup) {
	api := app.Group("/email")
	api.Use(middleware.AuthMiddleware())
	api.Use(middleware.DecryptMiddleware())
	{
		api.POST("/sendverification", SendVerificationEmail)
		api.GET("/verification-status", CheckEmailVerificationStatus)
	}
}

func CallbackRoutes(app *gin.RouterGroup) {
	app.GET("/update-verify/:id", VerifyEmail)
}
