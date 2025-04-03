package staticparameters

import (
	"bitbucket.org/paydoh/paydoh-commons/middleware"
	"github.com/gin-gonic/gin"
)

func Routes(api *gin.RouterGroup) {
	api.GET("/static-parameters", middleware.ResponseEncryptionMiddleware(), GetStaticParamter)
	api.GET("/get-secrets", middleware.ResponseEncryptionMiddleware(), GetSecrets)
}
