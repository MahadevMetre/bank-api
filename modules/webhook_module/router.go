package webhookmodule

import (
	"bankapi/middleware"

	"github.com/gin-gonic/gin"
)

func Routes(c *gin.RouterGroup) {
	webhookmodule := c.Group("/webhook")
	{
		webhookmodule.GET("/route-mobile", GetRouteMobileData)
		webhookmodule.POST("/kvb/vcip", UpdateVcipData)
		webhookmodule.POST("/kvb/kyc", KycUpdateState)
	}

	{
		webhookmodule.POST("/rewards-point", middleware.CallbackMiddleware(), ProvideRewardPoint)

	}
}
