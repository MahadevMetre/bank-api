package authenticationmodule

import (
	"github.com/gin-gonic/gin"

	"bankapi/middleware"
)

func Routes(app *gin.RouterGroup) {

	authentication := app.Group("/authentication")
	authentication.Use(middleware.AuthMiddleware())
	{
		authentication.POST("/initiate-sim-verification", InitiateSimVerification)
		authentication.GET("/sim-verification-status", GetSimVerificationStatus)
		authentication.POST("/logout", UserLogout)
	}
}
