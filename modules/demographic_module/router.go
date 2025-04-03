package demographic_module

import (
	"bankapi/middleware"
	"github.com/gin-gonic/gin"
)

func Routes(app *gin.RouterGroup) {
	demographic := app.Group("/demographic_module")
	demographic.Use(middleware.AuthMiddleware())
	demographic.Use(middleware.DecryptMiddleware())

	{
		demographic.GET("/fetch", GetDemographicData)
	}

}
