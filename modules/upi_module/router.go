package upi_module

import (
	"github.com/gin-gonic/gin"

	"bankapi/middleware"
)

func Routes(app *gin.RouterGroup) {
	upi := app.Group("/upi")
	upi.Use(middleware.AuthMiddleware())
	upi.Use(middleware.DecryptMiddleware())
	{
		//POST
		upi.POST("/create-upi-id", CreateUpiId)
		upi.POST("/remapping-upi-id", RemapUpiId)
		upi.POST("/aadhar-verification", UsersReqlistaccount)
		upi.POST("/set-upi-pin", SetUpiPin)
		upi.POST("/pay-val-vpa", ValidateVpaPayment)
		upi.POST("/pay-vpa", PayWithVpa)
		upi.POST("/account-balance", GetAccountBalance)
		upi.POST("/payeename", GetPayeeName)
		upi.POST("/get-upi-token", GetUpiToken)
		upi.POST("/get-upi-xml-payload", GetUpiXmlToken)
		upi.POST("/collect-money", UpiCollectMoney)
		upi.POST("/transaction-id", GetUpiTransactionId)
		upi.POST("/transaction-history", GetUpiTransactionDetails)
		upi.POST("/change-upi-pin", ChangeUpiPin)

		//GET
		upi.GET("/fetch-account-list", GetAccountLists)
		upi.GET("/collect-count", UpiCollectCount)
		upi.POST("/simbinding/sms-verification", UpiSimBindingAndSmsVerification)

	}
}
