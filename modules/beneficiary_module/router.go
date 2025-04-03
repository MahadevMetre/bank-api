package beneficiarymodule

import (
	"bankapi/middleware"

	"github.com/gin-gonic/gin"
)

func Routes(api *gin.RouterGroup) {
	beneficiary := api.Group("/beneficiary")
	beneficiary.Use(middleware.AuthMiddleware())
	beneficiary.Use(middleware.DecryptMiddleware())
	{
		// GET apis
		beneficiary.GET("/search", SearchBeneficiary)

		// POST apis
		beneficiary.POST("/add-beneficiary", AddNewBeneficiary)
		beneficiary.POST("/beneficiary-otp", BeneficiaryOTP)
		beneficiary.POST("/payment", BeneficiaryPayment)
		beneficiary.POST("/payment-otp", BeneficiaryPaymentOTP)
		beneficiary.POST("/payment-status", BeneficiaryPaymentStatus)
		beneficiary.POST("/quick-transfer-template", QuickTransferTemplate)
	}
}
