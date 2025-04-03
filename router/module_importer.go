package router

import (
	"bankapi/middleware"
	"bankapi/modules/account_create_callback"
	addressupdate "bankapi/modules/address_update"
	"bankapi/modules/faq"
	"bankapi/modules/mail"
	"bankapi/modules/sms_callback"
	staticParameters "bankapi/modules/static_parameters"

	authenticationmodule "bankapi/modules/authentication_module"
	authorizationmodule "bankapi/modules/authorization_module"
	beneficiarymodule "bankapi/modules/beneficiary_module"
	consentmodule "bankapi/modules/consent_module"
	encryptRouter "bankapi/modules/data_encryption"
	debitCard "bankapi/modules/debitcard_module"
	"bankapi/modules/demographic_module"
	get_user_details "bankapi/modules/get_user_details"
	"bankapi/modules/kyc_audit_data"
	"bankapi/modules/kyc_module"
	"bankapi/modules/kyc_update_data_callback"
	"bankapi/modules/nominee_module"
	onboardingmodule "bankapi/modules/onboarding_module"
	"bankapi/modules/openmodules"
	"bankapi/modules/payment_callback"
	"bankapi/modules/statementmodule"
	transactionsRouter "bankapi/modules/transaction_module"
	"bankapi/modules/upi_module"
	webhookmodule "bankapi/modules/webhook_module"

	"github.com/gin-gonic/gin"
)

func ProjectModules(app *gin.Engine) {

	api := app.Group("/api")
	{
		authorizationmodule.Routes(api)
		authenticationmodule.Routes(api)
		webhookmodule.Routes(api)
		openmodules.Routes(api)
		onboardingmodule.Routes(api)
		kyc_module.Routes(api)
		demographic_module.Routes(api)
		nominee_module.Routes(api)
		beneficiarymodule.Routes(api)
		consentmodule.Routes(api)
		upi_module.Routes(api)
		statementmodule.Routes(api)
		encryptRouter.Routes(api)
		transactionsRouter.Routes(api)
		get_user_details.Routes(api)
		debitCard.Routes(api)
		addressupdate.Router(api)
		staticParameters.Routes(api)
		mail.Routes(api)
		faq.Routes(api)
	}

	// callback APIs
	callbackAPI := app.Group("/callback")
	sms_callback.Routes(callbackAPI)
	mail.CallbackRoutes(callbackAPI)
	kyc_update_data_callback.Routes(callbackAPI)

	payment_callback.Routes(callbackAPI)
	kyc_audit_data.Routes(callbackAPI)
	account_create_callback.Routes(callbackAPI)

	// internal service call apis
	{
		app.POST("api/transaction/history/internal", middleware.CallbackMiddleware(), transactionsRouter.GetInternalTransactions)
	}

}
