package account_create_callback

import (
	"bankapi/middleware"

	"github.com/gin-gonic/gin"
)

func Routes(app *gin.RouterGroup) {
	app.Use(middleware.CallbackMiddleware())
	app.POST("/account-create", AccountCreateCallback)
}
