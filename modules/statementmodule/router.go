package statementmodule

import (
	"github.com/gin-gonic/gin"

	"bankapi/middleware"
)

func Routes(api *gin.RouterGroup) {
	rt := api.Group("/statement")
	rt.Use(middleware.AuthMiddleware())
	rt.Use(middleware.DecryptMiddleware())
	{
		rt.POST("/get-statement", GetStatement)
	}
}
