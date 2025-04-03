package data_encryption

import (
	"bankapi/middleware"

	"github.com/gin-gonic/gin"
)

func Routes(app *gin.RouterGroup) {
	app.Use(middleware.AuthMiddleware())
	app.POST("/data-encrypt", DataEncryption)
	app.POST("/data-decrypt", DataDecryption)
	// app.POST("/bank-encrypt", BankEncryption)
	// app.POST("/bank-decrypt", BankDecryption)
}
