package onboardingmodule

import (
	"bankapi/middleware"

	"github.com/gin-gonic/gin"
)

func Routes(api *gin.RouterGroup) {
	onboarding := api.Group("/onboarding")
	onboarding.Use(middleware.AuthMiddleware())
	onboarding.Use(middleware.DecryptMiddleware())
	{
		onboarding.POST("/personal-information", UpdatePersonalInformation)
		onboarding.GET("/personal-information", GetPersonalInformation)
		onboarding.POST("/create-account", CreateAccount)
		onboarding.GET("/get-account-details", GetAccountDetails)
		onboarding.GET("/user-status", GetUserOnboardingStatus)
		onboarding.POST("/pincode/details", GetPincodeDetails)
	}
}
