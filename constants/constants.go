package constants

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"bitbucket.org/paydoh/paydoh-commons/settings"

	"bankapi/responses"
)

func getGolangPort() int {
	portString := settings.Config("GOLANG_PORT")
	port, _ := strconv.Atoi(portString)

	return port
}

func GetNamePrefix(maritalStatus, gender string) string {
	if gender == "M" {
		return "Mr"
	} else if gender == "F" {
		switch maritalStatus {
		case "Single":
			return "Ms"
		case "Married":
			return "Mrs"
		case "Divorced":
			return "Ms"
		default:
			return "Ms"
		}
	}
	return "Mr"
}

func getRedisNetworkType() string {
	return settings.Config("REDIS_NETWORK_TYPE")
}

func getRedisURL() string {
	return settings.Config("REDIS_URL")
}

func getRedisUserName() string {
	return settings.Config("REDIS_USERNAME")
}

func getRedisPassword() string {
	return settings.Config("REDIS_PASSWORD")
}

func getRedisDB() int {
	redisDB := settings.Config("REDIS_DB")

	db, err := strconv.Atoi(redisDB)
	if err != nil {
		return 0
	}
	return db
}

func getPostgresHost() string {
	return settings.Config("POSTGRES_HOST")
}

func getPostgresPort() int {
	postgresPORT := settings.Config("POSTGRES_PORT")

	port, _ := strconv.Atoi(postgresPORT)

	return port
}

func getPostgresUserName() string {
	return settings.Config("POSTGRES_USER_NAME")
}

func getPostgresPassword() string {
	return settings.Config("POSTGRES_PASSWORD")
}

func getPostgresDatabase() string {
	return settings.Config("POSTGRES_DATABASE")
}

func getAesPassPhrase() string {
	return settings.Config("AES_PASSPHRASE")
}

func getJwtKey() string {
	return settings.Config("JWT_KEY")
}

func getCardKey() string {
	return settings.Config("CARD_KEY")
}

func getCardCvvKey() string {
	return settings.Config("CARD_CVV_KEY")
}

func getKVBUrl() string {
	return settings.Config("KVB_UAT_URL")
}

func getNotificationURL() string {
	return settings.Config("NOTIFICATION_SERVICE")
}

func getMongodbURI() string {
	return settings.Config("MONGODB_URI")
}

func getBankEncryptionKey() string {
	return settings.Config("BANK_ENCRYPTION_KEY")
}

func getCardControlEncryptionKey() string {
	return settings.Config("CARDCONTROL_ENCRYPTION_KEY")
}

func getAWSRegion() string {
	awsRegion := settings.Config("AWS_REGION")
	return awsRegion
}

func getAWSAccessKeyID() string {
	awsAccessKeyID := settings.Config("AWS_ACCESS_KEY_ID")
	return awsAccessKeyID
}

func getAWSSecretAccessKey() string {
	awsSecretAccessKey := settings.Config("AWS_SECRET_ACCESS_KEY")
	return awsSecretAccessKey
}

func getAWSBucketName() string {
	awsBucketName := settings.Config("AWS_BUCKET_NAME")
	return awsBucketName
}

func getAWSCloudFrontURL() string {
	awsCloudFrontURL := settings.Config("BANK_API_AWS_CLOUDFRONT_URL")
	return awsCloudFrontURL
}

func getKvbUserName() string {
	return settings.Config("KVB_USERNAME")
}

func getKvbPassword() string {
	return settings.Config("KVB_PASSWORD")
}

func getPaymentServiceURL() string {
	return settings.Config("PAYMENT_URL")
}

func getDebitCardPaymentAmount() string {
	return settings.Config("DEBIT_CARD_PAYMENT_AMT")
}

func getTollFreeNumber() string {
	return settings.Config("TOLL_FREE_NUMBER")
}

func getSupportMailId() string {
	return settings.Config("SUPPORT_MAIL_ID")
}

func getBranchName() string {
	return settings.Config("BRANCH_NAME")
}

func getIfscCode() string {
	return settings.Config("IFSC_CODE")
}

func getEmailVerificationExpireTime() string {
	return settings.Config("EMAIL_VERIFICATION_EXPIRE_TIME")
}

func getJwtExpTime() time.Duration {
	defaultJwtExpTime := 15 * time.Minute
	jwtExpTime := settings.Config("JWT_EXPIRE_TIME")

	expTime, err := strconv.Atoi(jwtExpTime)
	if err != nil || expTime <= 0 {
		return defaultJwtExpTime
	}

	return time.Duration(expTime) * time.Minute
}

func getRewardDatabaseName() string {
	return settings.Config("REWARD_DATABASE_NAME")
}

var (
	AWSRegion              = getAWSRegion()
	AWSAccessKeyID         = getAWSAccessKeyID()
	AWSSecretAccessKey     = getAWSSecretAccessKey()
	AWSBucketName          = getAWSBucketName()
	AWSCloudFrontURL       = getAWSCloudFrontURL()
	GolangPort             = getGolangPort()
	RedisNetworkType       = getRedisNetworkType()
	RedisURL               = getRedisURL()
	RedisUserName          = getRedisUserName()
	RedisPassword          = getRedisPassword()
	RedisDB                = getRedisDB()
	PostgresHost           = getPostgresHost()
	PostgresPort           = getPostgresPort()
	PostgresUsername       = getPostgresUserName()
	PostgresPassword       = getPostgresPassword()
	PostgresDatabase       = getPostgresDatabase()
	AesPassPhrase          = getAesPassPhrase()
	JwtKey                 = getJwtKey()
	KvbUatURL              = getKVBUrl()
	NotificationURL        = getNotificationURL()
	MongodbURI             = getMongodbURI()
	KvbUserName            = getKvbUserName()
	KvbPassword            = getKvbPassword()
	PaymentServiceURL      = getPaymentServiceURL()
	BankEncryptionKey      = getBankEncryptionKey()
	DebitCardPaymentAmount = getDebitCardPaymentAmount()
	CardKey                = getCardKey()
	CardCvvKey             = getCardCvvKey()

	IfscCode                    = getIfscCode()
	EmailVerificationExpireTime = getEmailVerificationExpireTime()
	JwtExpTime                  = getJwtExpTime()
	CardControlEncryptionKey    = getCardControlEncryptionKey()
	TollFreeNumber              = getTollFreeNumber()
	SupportMailID               = getSupportMailId()
	RazorPayID                  = getRazorPayID()
	MERCHANT_ID                 = getMerchantID()
	LONGCODE_ENCRYPT_TEXT       = getLongCodeEncryptText()
	KEystorePath                = getKeystorePath()
	KeyAlias                    = getKeyAlias()
	KeyPassword                 = getKeyPassword()
	KeySharedKey                = getKeySharedKey()
	BranchName                  = "Chennai Nandanam D.B.U"
	BranchAddress               = "AAVIN Illam, No. 3A, Pasumpon Muthuramalingam Salai, Nandanam, Chennai - 600035"
	BankName                    = "KARUR VYSYA BANK - KVB"
	RewardDatabaseName          = getRewardDatabaseName()
	PAYDOH_REWARDS              = "rewards"
)

var (
	ErrUserNotFound            = errors.New("user not found")
	ErrNoUsersFound            = errors.New("no users found")
	ErrDeviceNotFound          = errors.New("device not found")
	ErrNoDevicesFound          = errors.New("no devices found")
	ErrDeviceIpMismatch        = errors.New("device ip mismatch")
	ErrDeviceIdNotFound        = errors.New("device id not found")
	ErrNoRelationsFound        = errors.New("no relations found")
	ErrRelationNotFound        = errors.New("relation not found")
	ErrNoDataFound             = errors.New("no data found")
	ErrOtpIsRequired           = errors.New("otp is required")
	ErrBeneficiaryIdIsRequired = errors.New("beneficiary id is required")
	ErrKycConsentNotProvided   = errors.New("kyc consent for given number is not provided")
)

const (
	UPI                 string = "UPI"
	BANK                string = "BANK"
	CONSENT             string = "CONSENT"
	AUTHENTICATION      string = "SIM AUTHENTICATION"
	AUTHORIZATION       string = "AUTHORIZATION"
	BENEFICIARY         string = "BENEFICIARY"
	DEMOGRAPHIC         string = "DEMOGRAPHIC"
	KYC                 string = "KYC"
	KYC_AUDIT_DATA      string = "KYC_AUDIT_DATA"
	NOMINEE             string = "NOMINEE"
	ONBOARDING          string = "ONBOARDING"
	OPEN                string = "OPEN"
	PAYMENT_BENEFICIARY string = "PAYMENT_BENEFICIARY"
	STATEMENT           string = "STATEMENT"
	TRANSACTION_HISTORY string = "TRANSACTION_HISTORY"
	USER_DETAILS        string = "USER_DETAILS"
	WEBHOOK             string = "WEBHOOK"
	DEBITCARD           string = "DEBITCARD"
)

func StreamCallbackResponse(reader io.Reader, callbackChan chan<- *responses.CallbackCreateBankResponse) {
	decoder := json.NewDecoder(reader)
	for {
		callbackResponse := responses.NewCallbackCreatebankResponse()
		if err := decoder.Decode(&callbackResponse); err != nil {
			if err != io.EOF {
				fmt.Println("Error decoding callback response:", err)
			}
			break
		}

		// Send parsed response to the channel
		callbackChan <- callbackResponse
	}
	// Close the channel to signal that no more responses will be sent
	close(callbackChan)
}

// stages
const (
	SIM_VERIFICATION_STAGE      = "SIM_VERIFICATION"
	KYC_CONSENT_STAGE           = "KYC_CONSENT"
	DEMOGRAPHIC_FETCH_STAGE     = "DEMOGRAPHIC_FETCH"
	ACCOUNT_CREATION_STAGE      = "ACCOUNT_CREATION"
	DEBIT_CARD_CONSENT_STAGE    = "DEBIT_CARD_CONSENT"
	DEBIT_CARD_PAYMENT_STAGE    = "DEBIT_CARD_PAYMENT"
	DEBIT_CARD_GENERATION_STAGE = "DEBIT_CARD_GENERATION"
	UPI_GENERATION_STAGE        = "UPI_GENERATION"
	UPI_PIN_SETUP_STAGE         = "UPI_PIN_SETUP"
	M_PIN_SETUP_STAGE           = "M_PIN_SETUP"
)

// onboarding steps
const (
	//sim verification
	AUTHORIZATION_STEP    = "AUTHORIZATION"
	PERSONAL_DETAILS_STEP = "PERSONAL_DETAILS"

	//kyc consent
	AGENT_URL_STEP    = "AGENT_URL"
	KYC_CALLBACK_STEP = "KYC_CALLBACK"

	// account creation
	ACCOUNT_CALLBACK_STEP = "ACCOUNT_CALLBACK"
	AUDIT_CALLBACK_STEP   = "AUDIT_CALLBACK"

	EnterMobile       = "Enter_Mobile"
	EnterPersonalInfo = "Enter_Personal_Info"
	Consent           = "Consent"
	KYCInvoke         = "KYC_Invoke"
	KYCRedirection    = "KYC_Redirection"
	KYCAudit          = "KYC_Audit"
	Demographic       = "Demographic"
	AccountCreate     = "Account_Create"
	UPIIDCreate       = "UPI_ID_Create"
	SetUPIPin         = "Set_UPI_PIN"
	SetMPIN           = "Set_MPIN"
)

var (
	MONGO_USER_DB               = os.Getenv("MONGO_USER_DB")
	MONGO_USER_COLLECTION       = "users"
	NotificationPrefrenceMaster = "notification_prefrence_master"
	NotificationPref_Collection = "notification_prefrence"
)

var (
	UpiPaymentTxnDescPREFIX = []string{"UPI-CR-", "UPI-DR-"}
)

const (
	IntenationalIndex = "11,22,33,14"
	DomesticIndex     = "01,02,03,04"
)

const (
	TokenKeyFormat = "token:%s"

	ATM                      = "01"
	POS                      = "02"
	ECOMMERCE                = "03"
	CONTACTLESS              = "04"
	INTERNATIONALATM         = "11"
	INTERNATIONALPOS         = "22"
	INTERNATIONALECOMMERCE   = "33"
	INTERNATIONALCONTACTLESS = "14"

	// INTERNATIONALPOS         = "12"
	// INTERNATIONALECOMMERCE   = "13"

)

var TransactionTypes = map[string]string{
	"ATM":                       ATM,
	"POS":                       POS,
	"ECOMMERCE":                 ECOMMERCE,
	"CONTACTLESS":               CONTACTLESS,
	"INTERNATIONAL ATM":         INTERNATIONALATM,
	"INTERNATIONAL POS":         INTERNATIONALPOS,
	"INTERNATIONAL ECOMMERCE":   INTERNATIONALECOMMERCE,
	"INTERNATIONAL CONTACTLESS": INTERNATIONALCONTACTLESS,
}

var AllowedDomains = []string{
	"gmail.com",
	"outlook.com",
	"yahoo.com",
	"yahoo.co.in",
	"hotmail.com",
	"in.com",
	"zoho.com",
	"rediffmail.com",
	"icloud.com",
}

const (
	AuditLogType = "audit_logs"
)

const (
	BeneficiaryFetchResponseNoRecordFound = "01"
)

const (
	BankApiTechnicalError                          = "MW9999"
	BankApiTechnicalError1                         = "MW9998"
	BankApiUnknownError                            = "MW999"
	BankApiOtherTechnicalError                     = "91"
	AccountCreationApiTurnoverError                = "93155"
	AccountDetailFetchInvalidAccountNumberError    = "MW0004"
	DocumentUploadInvalidInputError                = "MW0030"
	BeneficiaryInvalidInputError                   = "MW0001"
	BankApiInvalidInputError                       = "MW0002"
	BankApiRecordNotExistError                     = "01"
	CasaTransactionHistoryNoAccountFoundError      = "2778"
	ConsentAplicantDataNotAvailableError           = "MW0016"
	DemographicFetchNoDataFoundError               = "MW0003"
	NomineeFailureError                            = "MW004"
	NomineeInvalidInputError                       = "MW0005"
	NomineeFetchInvalidInputError                  = "MW0001"
	PaymentInvalidBeneficiaryDetailError           = "MW0011"
	BankApiBackOfficeTimeOutError                  = "MW9997"
	VcipInvokeNoAgentsAvailableError               = "200"
	VcipInvokeNumberAlreadyExistError              = "205"
	VcipInvokeInternalServerError                  = "500"
	DebitCardTechnicalError                        = "99"
	DebitCardVirtualGenerationInvalidInputError    = "MW0005"
	DebitCardPhysicalInvalidCustomerDataError      = "MW0031"
	DebitCardInvalidInputError                     = "MW0030"
	DebitCardFetchInvalidCustomerDetailsError      = "MW0010"
	IfscSyncNoRecordFoundError                     = "01"
	DebitPinsetOrResetInvalidApplicantDetailsError = "MW0050"
	AddDebitCardFailedToAddCardError               = "11"
	AddDebitCardInvalidInputError                  = "MW500"
	EditTransactionNotSavedDataError               = "01"
	DebitCardBlockCardNotEditedError               = "02"
	OtpGenerationAndOtpValidationError             = "90"
	AddressUpdateRequestAlreadyPendingError        = "91"
	NomineeBankTokenExpiredError                   = "90"
	CardPinSetOtpExpiredError                      = "MW0054"
	CardPinResetOtpUtilizedError                   = "MW0055"
	CardListViewCardNotAvailableError              = "01"
	CardControlEditNotAbleToEditError              = "2"
)

const (
	UpiTechnicalIssueError = "1"
	UpiCryptoInfoError     = "2"
	UpiInvalidMpinError    = "0"
	UpiSmsNotReceivedError = "SMS not received"
)

const (
	UpiSimBindingKey = "upi:simbinding:%s"
	NomineeKey       = "nominee:%s"
	UpiSimBindingTTL = 10 * time.Minute
 
	// Senerio 0 - Request Sucessfully submited and Response Received For Beneficiary addition to KVB for OTP verification we are saving this data in redis
    // For Better approch we are passing this Key from the constants file
	BeneficiaryKey = "user:beneficary:%s"
	BeneficiaryTTL = 10 * time.Minute
)

var (
	ROUTE_SMS_MOBILE_NUMBER = getRouteSmsMobileNumber()
	LOG_LEVEL               = getLogLevel()
)

func getRouteSmsMobileNumber() string {
	awsAccessKeyID := settings.Config("ROUTE_SMS_MOBILE_NUMBER")
	return awsAccessKeyID
}

func getLogLevel() string {
	logLevel := settings.Config("LOG_LEVEL")
	if logLevel == "" {
		return "debug"
	}
	return logLevel
}

func getRazorPayID() string {
	razorPayID := settings.Config("RAZOR_KEY_ID")
	return razorPayID
}

func getMerchantID() string {
	merchantID := settings.Config("MERCHANT_ID")
	return merchantID
}

func getLongCodeEncryptText() string {
	longCodeEncryptText := settings.Config("LONGCODE_ENCRYPT_TEXT")
	return longCodeEncryptText
}

func getKeyAlias() string {
	keyAlias := settings.Config("KEY_ALIAS")
	return keyAlias
}

func getKeyPassword() string {
	keyPassword := settings.Config("KEY_PASSWORD")
	return keyPassword
}

func getKeySharedKey() string {
	keySharedKey := settings.Config("KEY_SHARED_KEY")
	return keySharedKey
}

func getKeystorePath() string {
	keystorePath := settings.Config("KEYSTORE_PATH_PASSWORD")
	return keystorePath
}
