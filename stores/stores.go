package stores

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"bitbucket.org/paydoh/paydoh-commons/amazon"
	"bitbucket.org/paydoh/paydoh-commons/database"
	"bitbucket.org/paydoh/paydoh-commons/pkg/storage"
	"bitbucket.org/paydoh/paydoh-commons/pkg/task"
	"bitbucket.org/paydoh/paydoh-commons/services"
	"github.com/gin-gonic/gin"

	"bankapi/constants"
	"bankapi/models"
	"bankapi/rpc"
	bankServices "bankapi/services"
	"bankapi/stores/address"
	"bankapi/stores/authentication"
	"bankapi/stores/authorization"
	"bankapi/stores/beneficiary"
	"bankapi/stores/consent"
	debitcard "bankapi/stores/debit_card"
	"bankapi/stores/demographic"
	"bankapi/stores/faq"
	"bankapi/stores/kyc"
	"bankapi/stores/kyc_audit_data"
	"bankapi/stores/mail"
	"bankapi/stores/nominee"
	"bankapi/stores/onboarding"
	"bankapi/stores/open"
	"bankapi/stores/payment_beneficiary"
	"bankapi/stores/statement"
	staticParameters "bankapi/stores/static_parameters"
	"bankapi/stores/transaction"
	"bankapi/stores/upi"
	"bankapi/stores/user_details"
	"bankapi/stores/webhook"
)

type Stores struct {
	Authorization      *authorization.AuthorizationStore
	Authentication     *authentication.AuthenticationStore
	Webhook            *webhook.WebhookStore
	Open               *open.OpenStore
	Onboarding         *onboarding.Store
	Kyc                *kyc.Store
	Demographic        *demographic.Store
	Nominee            *nominee.Store
	Beneficiary        *beneficiary.Store
	Consent            *consent.Store
	Upi                *upi.Store
	Statement          *statement.Store
	KycAuditStore      *kyc_audit_data.KycAuditStore
	Payment            *payment_beneficiary.PaymentCallbackStore
	UserDetail         *user_details.Store
	DebitCard          *debitcard.Store
	TransactionHistory *transaction.TransactionStore
	StaticParams       *staticParameters.Store
	AddressStore       *address.AddressStore
	EmailStore         *mail.Store
	AuditLogService    bankServices.AuditLogService
	FaqStore           faq.FAQStore
}

func NewStores(
	logSrv *services.LoggerService,
	db *sql.DB,
	memory *database.InMemory,
	mongo *database.Document,
	ctx context.Context,
	client rpc.PaymentServiceClient,
) *Stores {

	awsKeyId := constants.AWSAccessKeyID
	awsSecretKey := constants.AWSSecretAccessKey
	awsRegion := constants.AWSRegion

	s3Client, err := storage.NewS3Client(awsKeyId, awsSecretKey, awsRegion)
	if err != nil {
		fmt.Println(err)
	}

	newTaskEnqueuer := task.NewAsynqTaskEnqueuer(constants.RedisURL, constants.RedisUserName, constants.RedisPassword, constants.RedisDB)
	auditLogSrv := bankServices.NewAuditLogService(logSrv, newTaskEnqueuer)
	authorizationStore := authorization.NewAuthorizationStore(logSrv, db, ctx, mongo, memory, newTaskEnqueuer, auditLogSrv)
	authenticationStore := authentication.NewAuthenticationStore(logSrv, db, mongo, memory, auditLogSrv)
	webhookStore := webhook.NewWebhookStore(logSrv, db, memory)
	openStore := open.NewOpenStore(logSrv, db, mongo, memory, ctx, client, auditLogSrv)
	o := onboarding.NewStore(logSrv, db, mongo, memory, authorizationStore)
	k := kyc.NewStore(logSrv, db, mongo, memory, auditLogSrv)
	d := demographic.NewStore(logSrv, db, mongo, memory)
	n := nominee.NewStore(logSrv, db, mongo, memory, auditLogSrv)
	bn := beneficiary.NewStore(logSrv, db, mongo, memory)
	cn := consent.NewStore(logSrv, db, mongo, memory)
	u := upi.NewStore(logSrv, db, mongo, memory, memory, auditLogSrv)
	kas := kyc_audit_data.NewKycAuditStore(logSrv, db, mongo, memory)
	paymentCallback := payment_beneficiary.NewPaymentCallbackStore(logSrv, db, mongo, memory)
	userDetails := user_details.NewStore(logSrv, db, mongo, memory, memory)
	debitcard := debitcard.NewStore(logSrv, db, mongo, memory, s3Client, auditLogSrv, newTaskEnqueuer)
	txnHistory := transaction.NewTransactionStore(logSrv, db, mongo, memory)
	st := statement.NewStore(logSrv, db, mongo, memory, txnHistory)
	staticParameters := staticParameters.NewStore(logSrv)
	updateAddress := address.NewStore(logSrv, db, memory, auditLogSrv)
	mail := mail.NewStore(logSrv, db, mongo, memory)
	faqStore := faq.NewFAQStore(logSrv)
	return &Stores{
		Authorization:      authorizationStore,
		Authentication:     authenticationStore,
		Webhook:            webhookStore,
		Open:               openStore,
		Onboarding:         o,
		Kyc:                k,
		Demographic:        d,
		Nominee:            n,
		Beneficiary:        bn,
		Consent:            cn,
		Upi:                u,
		Statement:          st,
		KycAuditStore:      kas,
		Payment:            paymentCallback,
		UserDetail:         userDetails,
		DebitCard:          debitcard,
		TransactionHistory: txnHistory,
		StaticParams:       staticParameters,
		AddressStore:       updateAddress,
		EmailStore:         mail,
		AuditLogService:    auditLogSrv,
		FaqStore:           faqStore,
	}
}

func (s *Stores) BindStore(awsIstance *amazon.Aws) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("authorization_store", s.Authorization)
		ctx.Set("authentication_store", s.Authentication)
		ctx.Set("webhook_store", s.Webhook)
		ctx.Set("open_store", s.Open)
		ctx.Set("onboarding_store", s.Onboarding)
		ctx.Set("kyc_store", s.Kyc)
		ctx.Set("demographic_store", s.Demographic)
		ctx.Set("nominee_store", s.Nominee)
		ctx.Set("beneficiary_store", s.Beneficiary)
		ctx.Set("consent_store", s.Consent)
		ctx.Set("upi_store", s.Upi)
		ctx.Set("statement_store", s.Statement)
		ctx.Set("aws_instance", awsIstance)
		ctx.Set("kyc_audit_store", s.KycAuditStore)
		ctx.Set("payment_store", s.Payment)
		ctx.Set("userdetail_store", s.UserDetail)
		ctx.Set("debitcard_store", s.DebitCard)
		ctx.Set("transaction_store", s.TransactionHistory)
		ctx.Set("static_parameters_store", s.StaticParams)
		ctx.Set("update_address_store", s.AddressStore)
		ctx.Set("email_store", s.EmailStore)
		ctx.Set("faq_store", s.FaqStore)
		ctx.Next()
	}
}

func GetStores(ctx *gin.Context) (*Stores, error) {
	ao, aOk := ctx.MustGet("authorization_store").(*authorization.AuthorizationStore)

	if !aOk {
		return nil, fmt.Errorf("authorization store not bound")
	}

	au, auOk := ctx.MustGet("authentication_store").(*authentication.AuthenticationStore)

	if !auOk {
		return nil, fmt.Errorf("authentication store not bound")
	}

	wb, wOk := ctx.MustGet("webhook_store").(*webhook.WebhookStore)

	if !wOk {
		return nil, fmt.Errorf("webhook store not bound")
	}

	os, osOk := ctx.MustGet("open_store").(*open.OpenStore)

	if !osOk {
		return nil, fmt.Errorf("open store not bound")
	}

	onboardingStore, onOk := ctx.MustGet("onboarding_store").(*onboarding.Store)

	if !onOk {
		return nil, fmt.Errorf("onboarding store not bound")
	}

	kycStore, kOk := ctx.MustGet("kyc_store").(*kyc.Store)

	if !kOk {
		return nil, fmt.Errorf("kyc store not bound")
	}

	demographicStore, dOk := ctx.MustGet("demographic_store").(*demographic.Store)

	if !dOk {
		return nil, fmt.Errorf("demographic_module store not bound")
	}

	nomineeStore, nOk := ctx.MustGet("nominee_store").(*nominee.Store)

	if !nOk {
		return nil, errors.New("nominee store not found")
	}

	beneficiaryStore, bnOk := ctx.MustGet("beneficiary_store").(*beneficiary.Store)

	if !bnOk {
		return nil, errors.New("beneficiary store not found")
	}

	consentStore, cnOk := ctx.MustGet("consent_store").(*consent.Store)

	if !cnOk {
		return nil, errors.New("consent store not found")
	}

	upiStore, uOk := ctx.MustGet("upi_store").(*upi.Store)

	if !uOk {
		return nil, errors.New("upi store not found")
	}

	statementStore, stOk := ctx.MustGet("statement_store").(*statement.Store)

	if !stOk {
		return nil, errors.New("statement store not found")
	}

	kycAuditStore, stOk := ctx.MustGet("kyc_audit_store").(*kyc_audit_data.KycAuditStore)

	if !stOk {
		return nil, errors.New("Kyc Audit store not found")
	}

	paymentStore, stOk := ctx.MustGet("payment_store").(*payment_beneficiary.PaymentCallbackStore)

	if !stOk {
		return nil, errors.New("payment store not found")
	}

	userDetails, stOk := ctx.MustGet("userdetail_store").(*user_details.Store)
	if !stOk {
		return nil, errors.New("userdetail store not found")
	}

	debitcard, stOk := ctx.MustGet("debitcard_store").(*debitcard.Store)
	if !stOk {
		return nil, errors.New("debitcard store not found")
	}
	txnHistory, stOk := ctx.MustGet("transaction_store").(*transaction.TransactionStore)
	if !stOk {
		return nil, errors.New("transaction store not found")
	}

	staticParameter, staticOk := ctx.MustGet("static_parameters_store").(*staticParameters.Store)
	if !staticOk {
		return nil, errors.New("Static parameters store not found")
	}
	updateAddress, stOk := ctx.MustGet("update_address_store").(*address.AddressStore)
	if !stOk {
		return nil, errors.New("update address store not found")
	}

	email, mailuOk := ctx.MustGet("email_store").(*mail.Store)

	if !mailuOk {
		return nil, fmt.Errorf("Email store not bound")
	}

	faqStore, ok := ctx.MustGet("faq_store").(faq.FAQStore)
	if !ok {
		return nil, fmt.Errorf("FAQ store not bound")
	}

	return &Stores{
		Authorization:      ao,
		Authentication:     au,
		Webhook:            wb,
		Open:               os,
		Onboarding:         onboardingStore,
		Kyc:                kycStore,
		Demographic:        demographicStore,
		Nominee:            nomineeStore,
		Beneficiary:        beneficiaryStore,
		Consent:            consentStore,
		Upi:                upiStore,
		Statement:          statementStore,
		KycAuditStore:      kycAuditStore,
		Payment:            paymentStore,
		UserDetail:         userDetails,
		DebitCard:          debitcard,
		TransactionHistory: txnHistory,
		StaticParams:       staticParameter,
		AddressStore:       updateAddress,
		EmailStore:         email,
		FaqStore:           faqStore,
	}, nil
}

func GetRequestPayload(ctx *gin.Context) (*models.RequestPayload, error) {
	decrypted, dOk := ctx.MustGet("decrypted").(string)

	if !dOk {
		return nil, fmt.Errorf("request payload not found")
	}

	return &models.RequestPayload{
		Payload: decrypted,
	}, nil
}

func GetAwsInstance(ctx *gin.Context) (*amazon.Aws, error) {
	awsInstance, aOk := ctx.MustGet("aws_instance").(*amazon.Aws)

	if !aOk {
		return nil, fmt.Errorf("aws instance not found")
	}

	return awsInstance, nil
}

func GetAuthValue(ctx *gin.Context) (*models.AuthValues, error) {
	userId, uOk := ctx.MustGet("user_id").(string)

	if !uOk {
		return nil, fmt.Errorf("auth values not found")
	}

	key, kOk := ctx.MustGet("key").(string)

	if !kOk {
		return nil, fmt.Errorf("auth values not found")
	}

	deviceIp, dOk := ctx.MustGet("device_ip").(string)

	if !dOk {
		return nil, fmt.Errorf("device ip not found")
	}

	os, osOk := ctx.MustGet("os").(string)

	if !osOk {
		return nil, fmt.Errorf("os not found")
	}

	osVersion, osVersionOk := ctx.MustGet("os_version").(string)

	if !osVersionOk {
		return nil, fmt.Errorf("os version not found")
	}

	latLong, latLongOk := ctx.MustGet("lat_long").(string)

	if !latLongOk {
		return nil, fmt.Errorf("lat long not found")
	}

	return &models.AuthValues{
		UserId:    userId,
		Key:       key,
		DeviceIp:  deviceIp,
		OS:        os,
		OSVersion: osVersion,
		LatLong:   latLong,
	}, nil
}

func (s *Stores) StartPeriodicTask() {
	go func(store *Stores) {
		now := time.Now()
		nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		durationUntilMidnight := nextMidnight.Sub(now)

		<-time.After(durationUntilMidnight)

		ticker := time.NewTicker(24 * time.Hour)

		defer ticker.Stop()

		for range ticker.C {
			store.Open.SyncIfscApi(context.Background())
		}
	}(s)
}
