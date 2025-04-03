package get_user_details_module

import (
	"bankapi/middleware"

	"github.com/gin-gonic/gin"
)

func Routes(app *gin.RouterGroup) {
	user := app.Group("/user")
	user.Use(middleware.AuthMiddleware())

	user.GET("/get-details", UserDetails)

}
