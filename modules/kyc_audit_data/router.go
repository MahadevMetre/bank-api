package kyc_audit_data

import (
	"bankapi/middleware"

	"github.com/gin-gonic/gin"
)

func Routes(app *gin.RouterGroup) {
	app.Use(middleware.CallbackMiddleware())
	app.POST("/audit-completion", KycAuditComplition)
}
