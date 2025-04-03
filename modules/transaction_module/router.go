package transaction_module

import (
	"bankapi/middleware"

	responseMiddleware "bitbucket.org/paydoh/paydoh-commons/middleware"
	"github.com/gin-gonic/gin"
)

func Routes(app *gin.RouterGroup) {

	route := app.Group("/transaction")

	route.Use(middleware.AuthMiddleware())
	route.Use(middleware.DecryptMiddleware())
	route.Use(responseMiddleware.ResponseEncryptionMiddleware())

	route.POST("/history", Transactions)
	route.POST("/history/get-details", GetTxnDetails)
	route.GET("/recent-users", GetRecentTransactionUsers)
}
