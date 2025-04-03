package authorizationmodule

import (
	"github.com/gin-gonic/gin"
)

func Routes(app *gin.RouterGroup) {

	authorization := app.Group("/authorization")
	{
		authorization.POST("", Authorization)
	}

	userAuthenticated := authorization.Group("/authenticated")
	{
		userAuthenticated.POST("/current-user", GetCurrentUserByUserId)
		userAuthenticated.POST("/current-user-mobile", GetCurrentUserMobileNumber)
	}
}
