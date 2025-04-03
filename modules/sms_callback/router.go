package sms_callback

import (
	"github.com/gin-gonic/gin"
)

func Routes(app *gin.RouterGroup) {
	app.GET("/sms-received", SmsReceivedCallback)
}
