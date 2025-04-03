package payment_callback

import (
	"bankapi/middleware"

	"github.com/gin-gonic/gin"
)

func Routes(app *gin.RouterGroup) {
	app.Use(middleware.CallbackMiddleware())
	app.POST("/normal-payment", PaymentCallBackAPI)
}
