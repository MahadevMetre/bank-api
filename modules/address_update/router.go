package addressupdate

import (
	"bankapi/middleware"

	"github.com/gin-gonic/gin"
)

func Router(app *gin.RouterGroup) {
	api := app.Group("/address")
	api.Use(middleware.AuthMiddleware())

	api.POST("/update", AddressUpdate)
	api.GET("/get-status", GetAddressUpdateStatus)
}
