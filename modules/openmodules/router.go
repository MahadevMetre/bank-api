package openmodules

import (
	"github.com/gin-gonic/gin"

	"bankapi/middleware"

	responseMiddleware "bitbucket.org/paydoh/paydoh-commons/middleware"
)

func Routes(api *gin.RouterGroup) {
	open := api.Group("/open")
	open.Use(middleware.AuthMiddleware())
	{

		open.GET("/relations", GetRelations)
		open.GET("/state", GetStates)
		open.GET("/cities/:state", GetCities)

		open.POST("/shipping-address", UpdateShippingAddress)
		open.GET("/shipping-address", GetShippingAddress)
		open.POST("/update-shipping-address", ShippingAddressUpdate)

		open.POST("/get-receipt-id", middleware.DecryptMiddleware(), responseMiddleware.ResponseEncryptionMiddleware(), GetReceiptId)
		open.POST("/payment-status", middleware.DecryptMiddleware(), responseMiddleware.ResponseEncryptionMiddleware(), UpdatePaymentStatus)
		open.GET("/payment-status", responseMiddleware.ResponseEncryptionMiddleware(), GetPaymentStatus)

		open.POST("/set-mpin", middleware.DecryptMiddleware(), responseMiddleware.ResponseEncryptionMiddleware(), SetMpin)
		open.POST("/verify-mpin", middleware.DecryptMiddleware(), responseMiddleware.ResponseEncryptionMiddleware(), VerifyMpin)
		open.POST("/reset-mpin", middleware.DecryptMiddleware(), responseMiddleware.ResponseEncryptionMiddleware(), ReSetMpin)
		open.POST("/verify-forgot-mpin", middleware.DecryptMiddleware(), responseMiddleware.ResponseEncryptionMiddleware(), verifyForgetMpinReset)
		open.POST("/update-mpin", middleware.DecryptMiddleware(), responseMiddleware.ResponseEncryptionMiddleware(), UpdateForgetMpinReset)

		open.GET("/ifsc-data/:bank", GetIfscData)
		open.GET("/ifsc-data/banks", GetBanks)
		open.GET("/ifsc-data/banks/:ifsc", GetIfscDataByIfscCode)

		open.POST("/update-fcm-token", UpdateFcmToken)
	}
}
