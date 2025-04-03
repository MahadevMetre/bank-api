package nominee_module

import (
	"github.com/gin-gonic/gin"

	"bankapi/middleware"
)

func Routes(api *gin.RouterGroup) {
	nominee := api.Group("/nominee")
	nominee.Use(middleware.AuthMiddleware())
	nominee.Use(middleware.DecryptMiddleware())
	{
		nominee.POST("/add-nominee", AddNewNominee)
		nominee.POST("/verify-otp", VerifyNomineeOtp)
		nominee.GET("/fetch", FetchNominee)
	}
}
