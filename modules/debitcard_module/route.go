package debitcardmodule

import (
	"bankapi/middleware"

	"github.com/gin-gonic/gin"
)

func Routes(app *gin.RouterGroup) {
	debitcard := app.Group("/debitcard")

	debitcard.Use(middleware.AuthMiddleware())
	debitcard.Use(middleware.DecryptMiddleware())

	{
		debitcard.POST("/generate", GenerateDebitcard)
		debitcard.GET("/detail", GetDebitCardDetails)
		debitcard.POST("/set-debitcard-pin", SetDebitCardPin)
		debitcard.POST("/verify-otp", VerifyOTP)
		debitcard.GET("/track-status", TrackDebitCardStatus)

		debitcard.POST("/get-limit-list", GetTransactionLimit)
		debitcard.POST("/set-txn-limit", SetTransactonLimit)
		debitcard.GET("/get-card-status", GetCardStatus)
		debitcard.POST("/set-card-status", SetCardStatus)
	}
}
