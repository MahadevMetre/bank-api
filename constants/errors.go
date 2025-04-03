package constants

import "strings"

const (
	RetryErrorMessage            = "We cannot process your request this time. Please try again later."
	InputErrorMessage            = "Please enter correct input."
	InputErrorIfscCodeMessage    = "Please enter correct ifsc code."
	UpiNoRecordFoundErrorMessage = "No Record Found."
	UpiInvalidMpinErrorMessage   = "Please enter correct Upi Mpin."
	AddressUpdateInProgressError = "Your address modification request is in progress, bank will notify once your address will be updated"
	AadhaarNumberMismatchError   = "Please enter the correct first six digits of your Aadhaar number."
)

const (
	RetryErrorMW9999 = "MW9999"
	RetryErrorMW9997 = "MW9997"
	RetryErrorMW9998 = "MW9998"
	RetryErrorMW999  = "MW999"
	RetryUpiError1   = "1"
	DocumentUpload99 = "99"
)

var RetryBankErrors = map[string]string{
	RetryErrorMW9999: RetryErrorMessage,
	RetryErrorMW9997: RetryErrorMessage,
	RetryErrorMW9998: RetryErrorMessage,
	RetryErrorMW999:  RetryErrorMessage,
	RetryUpiError1:   RetryErrorMessage,
	DocumentUpload99: RetryErrorMessage,
}

var RetryUpiErrors = map[string]string{
	RetryUpiError1: RetryErrorMessage,
	RetryErrorMW9997: RetryErrorMessage,
	RetryErrorMW9998: RetryErrorMessage,
}

const (
	ConsentErrorCodeMW0011         = "MW0011"
	ConsentErrorCodeMW0012         = "MW0012"
	ConsentErrorCodeMW0013         = "MW0013"
	ConsentErrorCodeMW0014         = "MW0014"
	ConsentErrorCodeMW0015         = "MW0015"
	ConsentErrorCodeMW0016         = "MW0016"
	DebitCardControlErrorCode91    = "91"
	DebitCardControlErrorCode02    = "02"
	DebitCardControlErrorCode17    = "17"
	DebitCardControlErrorCode10    = "10"
	DebitCardControlErrorCode01    = "01"
	DebitCardControlErrorCode11    = "11"
	DebitCardControlErrorCodeMW500 = "MW500"
)

var ConsentErrorMessages = map[string]string{
	ConsentErrorCodeMW0011:         InputErrorMessage,
	ConsentErrorCodeMW0012:         InputErrorMessage,
	ConsentErrorCodeMW0013:         InputErrorMessage,
	ConsentErrorCodeMW0014:         InputErrorMessage,
	ConsentErrorCodeMW0015:         InputErrorMessage,
	ConsentErrorCodeMW0016:         InputErrorMessage,
	DebitCardControlErrorCode91:    InputErrorMessage,
	DebitCardControlErrorCode02:    InputErrorMessage,
	DebitCardControlErrorCode17:    InputErrorMessage,
	DebitCardControlErrorCode10:    InputErrorMessage,
	DebitCardControlErrorCode01:    InputErrorMessage,
	DebitCardControlErrorCode11:    InputErrorMessage,
	DebitCardControlErrorCodeMW500: InputErrorMessage,
}

func GetConsentErrorMessage(errorCode string) (string, bool) {
	message, exists := ConsentErrorMessages[errorCode]
	return message, exists
}

const (
	AccountDetailFetchErrorCodeMW0001 = "MW0001"
	AccountDetailFetchErrorCodeMW0002 = "MW0002"
	AccountDetailFetchErrorCodeMW0004 = "MW0004"
)

var AccountDetailFetchErrorMessages = map[string]string{
	AccountDetailFetchErrorCodeMW0001: InputErrorMessage,
	AccountDetailFetchErrorCodeMW0002: InputErrorMessage,
	AccountDetailFetchErrorCodeMW0004: InputErrorMessage,
}

func GetAccountDetailFetchErrorMessage(errorCode string) (string, bool) {
	message, exists := AccountDetailFetchErrorMessages[errorCode]
	return message, exists
}
func GetCardControlService(errorCode string) (string, bool) {
	message, exists := ConsentErrorMessages[errorCode]
	return message, exists
}

const (
	CasaTxnErrorCode3611   = "3611"
	CasaTxnErrorCode6083   = "6083"
	CasaTxnErrorCode2778   = "2778"
	CasaTxnErrorCodeMW0010 = "MW0010"
	CasaTxnErrorCodeMW0011 = "MW0011"
)

var CasaTxnErrorMessages = map[string]string{
	CasaTxnErrorCode2778:   InputErrorMessage,
	CasaTxnErrorCodeMW0010: InputErrorMessage,
	CasaTxnErrorCodeMW0011: InputErrorMessage,
}

var CasaTxnErrorForRetry = map[string]string{
	CasaTxnErrorCode3611: RetryErrorMessage,
	CasaTxnErrorCode6083: RetryErrorMessage,
}

func GetCasaTxnErrorMessage(errorCode string) (string, bool) {
	message, exists := CasaTxnErrorMessages[errorCode]
	return message, exists
}

func GetCasaTxnErrorRetryMessage(errorCode string) (string, bool) {
	message, exists := CasaTxnErrorForRetry[errorCode]
	return message, exists
}

const (
	DemographicErrorCodeMW0001 = "MW0001"
	DemographicErrorCodeMW0002 = "MW0002"
	DemographicErrorCodeMW0003 = "MW0003"
)

var DemographicErrorMessages = map[string]string{
	DemographicErrorCodeMW0001: InputErrorMessage,
	DemographicErrorCodeMW0002: InputErrorMessage,
	DemographicErrorCodeMW0003: InputErrorMessage,
}

func GetDemographicErrorMessage(errorCode string) (string, bool) {
	message, exists := DemographicErrorMessages[errorCode]
	return message, exists
}

const (
	VCIPInvokeErrorCode205        = "205"
	VCIPInvokeErrorCode200        = "200"
	VCIPInvokeErrorCode500        = "500"
	VCIPInvokeErrorCode200Success = "success"
	VCIPInvokeErrorCode200Failure = "failure"
)

var VcipInvokeErrorMessages = map[string]string{
	VCIPInvokeErrorCode205: "The given Mobile Number already registered. Please try again with another number",
	VCIPInvokeErrorCode200: "Close the session and Reinitiate the request since agent unavailable to start VKYC.",
}

var VcipInvokeErrorForRetry = map[string]string{
	VCIPInvokeErrorCode500: RetryErrorMessage,
}

func GetVcipInvokeErrorMessage(errorCode, status string) (string, bool) {
	message, exists := VcipInvokeErrorMessages[errorCode]
	if !exists {
		return "", false
	}

	if errorCode == VCIPInvokeErrorCode200 {
		if strings.ToLower(status) == VCIPInvokeErrorCode200Failure {
			return message, true
		}

		if strings.ToLower(status) == VCIPInvokeErrorCode200Success {
			return "", false
		}

		return message, true
	}

	return message, true
}

func GetVcipInvokeErrorRetryMessage(errorCode string) (string, bool) {
	message, exists := VcipInvokeErrorForRetry[errorCode]
	return message, exists
}

const (
	AccountCreationErrorCodeCI95  = "CI95"
	AccountCreationErrorCode80004 = "80004"
	AccountCreationErrorCode80026 = "80026"
	AccountCreationErrorCode93134 = "93134"
	AccountCreationErrorCode93155 = "93155"
	AccountCreationErrorCode93143 = "93143"
)

var AccountCreationErrorMessages = map[string]string{
	AccountCreationErrorCode80004: InputErrorMessage,
	AccountCreationErrorCode80026: InputErrorMessage,
	AccountCreationErrorCode93155: "We cannot process your request this time, Annual Turnover code Not Found.",
}

var AccountCreationRetryableErrorCodes = map[string]string{
	AccountCreationErrorCodeCI95:  RetryErrorMessage,
	AccountCreationErrorCode93134: RetryErrorMessage,
	AccountCreationErrorCode93143: RetryErrorMessage,
}

func GetAccountCreationRetryErrorCode(errorCode string) (string, bool) {
	message, exists := AccountCreationRetryableErrorCodes[errorCode]
	return message, exists
}

func GetAccountCreationErrorMessage(errorCode string) (string, bool) {
	message, exists := AccountCreationErrorMessages[errorCode]
	return message, exists
}

const (
	IfscSyncErrorCode01 = "01"
)

var IfscSyncErrorMessages = map[string]string{
	IfscSyncErrorCode01: InputErrorMessage,
}

func GetIfscSyncErrorMessage(errorCode string) (string, bool) {
	message, exists := IfscSyncErrorMessages[errorCode]
	return message, exists
}

const (
	NomineeErrorCode99     = "99"
	NomineeErrorCodeMW0002 = "MW0002"
	NomineeErrorCodeMW0004 = "MW0004"
	NomineeErrorCodeMW0005 = "MW0005"
	NomineeErrorCodeMW0008 = "MW0008"
	NomineeErrorCodeMW0010 = "MW0010"
	NomineeErrorCode90     = "90"
)

var NomineeErrorMessages = map[string]string{
	NomineeErrorCodeMW0002: InputErrorMessage,
	NomineeErrorCodeMW0004: InputErrorMessage,
	NomineeErrorCodeMW0005: InputErrorMessage,
	NomineeErrorCodeMW0008: InputErrorMessage,
	NomineeErrorCodeMW0010: InputErrorMessage,
}

var NomineeErrorRetryMessages = map[string]string{
	NomineeErrorCode90: RetryErrorMessage,
	NomineeErrorCode99: RetryErrorMessage,
}

func GetNomineeErrorMessage(errorCode string) (string, bool) {
	message, exists := NomineeErrorMessages[errorCode]
	return message, exists
}

func GetNomineeErrorRetryMessages(errorCode string) (string, bool) {
	message, exists := NomineeErrorRetryMessages[errorCode]
	return message, exists
}

const (
	BeneficiaryErrorCodeMW0001 = "MW0001"
	BeneficiaryErrorCodeMW0002 = "MW0002"
	BeneficiaryErrorCodeMW0003 = "MW0003"
)

var BeneficiaryErrorMessages = map[string]string{
	BeneficiaryErrorCodeMW0001: InputErrorMessage,
	BeneficiaryErrorCodeMW0002: InputErrorMessage,
	BeneficiaryErrorCodeMW0003: InputErrorMessage,
}

// Scenario 3- Beneficiary details sent to KVB, but not received the response. (Back Ofice TimeOut)
var BeneficiaryRetryErrorMessages = map[string]string{
	RetryErrorMW9999: RetryErrorMessage,
	RetryErrorMW9997: RetryErrorMessage,
	RetryErrorMW9998: RetryErrorMessage,
}

func GetBeneficiaryErrorMessage(errorCode string) (string, bool) {
	message, exists := BeneficiaryErrorMessages[errorCode]
	return message, exists
}

// Scenario 3- Beneficiary details sent to KVB, but not received the response. (Back Ofice TimeOut)
func GetBeneficiaryRetryErrorMessage(errorCode string) (string, bool) {
	message, exists := BeneficiaryRetryErrorMessages[errorCode]
	return message, exists
}

const (
	PaymentCallbackErrorCodeMW0011 = "MW0011"
	PaymentCallbackErrorCodeMW0014 = "MW0014"
	PaymentCallbackErrorCodeMW0016 = "MW0016"
	PaymentCallbackErrorCodeMW0017 = "MW0017"
	PaymentCallbackErrorCodeMW0018 = "MW0018"
	PaymentCallbackErrorCodeMW0021 = "MW0021"
	PaymentCallbackErrorCodeMW0023 = "MW0023"
	PaymentCallbackErrorCodeMW0024 = "MW0024"
	PaymentCallbackErrorCodeMW0027 = "MW0027"
	PaymentCallbackErrorCodeMW0028 = "MW0028"
)

var PaymentCallbackErrorMessages = map[string]string{
	PaymentCallbackErrorCodeMW0016: InputErrorMessage,
	PaymentCallbackErrorCodeMW0017: InputErrorMessage,
}

var PaymentCallbackRetryErrorMessages = map[string]string{
	PaymentCallbackErrorCodeMW0011: RetryErrorMessage,
	PaymentCallbackErrorCodeMW0018: RetryErrorMessage,
	PaymentCallbackErrorCodeMW0021: RetryErrorMessage,
	PaymentCallbackErrorCodeMW0023: RetryErrorMessage,
	PaymentCallbackErrorCodeMW0024: RetryErrorMessage,
	PaymentCallbackErrorCodeMW0027: RetryErrorMessage,
	PaymentCallbackErrorCodeMW0028: RetryErrorMessage,
}

func GetPaymentCallbackErrorMessage(errorCode string) (string, bool) {
	message, exists := PaymentCallbackErrorMessages[errorCode]
	return message, exists
}

func GetPaymentCallbackRetryErrorMessage(errorCode string) (string, bool) {
	message, exists := PaymentCallbackRetryErrorMessages[errorCode]
	return message, exists
}

const (
	QuickTransferBeneficiaryErrorCodeMW0002 = "MW0002"
)

var QuickTransferBeneficiaryErrorMessages = map[string]string{
	QuickTransferBeneficiaryErrorCodeMW0002: InputErrorMessage,
}

func GetQuickTransferBeneficiaryErrorMessage(errorCode string) (string, bool) {
	message, exists := QuickTransferBeneficiaryErrorMessages[errorCode]
	return message, exists
}

const (
	VirtualDebitCardErrorCodeAN3101 = "AN3101"
	VirtualDebitCardErrorCodeMW0013 = "MW0013"
	VirtualDebitCardErrorCodeMW0014 = "MW0014"
	VirtualDebitCardErrorCodeMW0001 = "MW0001"
	VirtualDebitCardErrorCodeCE4000 = "CE4000"
	VirtualDebitCardErrorCodePX1103 = "PX1103"
	VirtualDebitCardErrorCodeLN1100 = "LN1100"
)

var VirtualDebitCardErrorMessages = map[string]string{
	VirtualDebitCardErrorCodeAN3101: InputErrorMessage,
	VirtualDebitCardErrorCodeMW0013: InputErrorMessage,
	VirtualDebitCardErrorCodeMW0014: InputErrorMessage,
	VirtualDebitCardErrorCodeMW0001: InputErrorMessage,
}

var VirtualDebitCardErrorMessagesRetry = map[string]string{
	VirtualDebitCardErrorCodeCE4000: RetryErrorMessage,
	VirtualDebitCardErrorCodePX1103: RetryErrorMessage,
	VirtualDebitCardErrorCodeLN1100: RetryErrorMessage,
}

func GetVirtualDebitCardErrorMessage(errorCode string) (string, bool) {
	message, exists := VirtualDebitCardErrorMessages[errorCode]
	return message, exists
}

func GetVirtualDebitCardErrorMessageRetry(errorCode string) (string, bool) {
	message, exists := VirtualDebitCardErrorMessagesRetry[errorCode]
	return message, exists
}

const (
	DebitCardFetchErrorCodeMW0010 = "MW0010"
	DebitCardFetchErrorCodeMW0030 = "MW0030"
	DebitCardFetchErrorCodePX1103 = "PX1103"
)

var DebitCardFetchErrorMessages = map[string]string{
	DebitCardFetchErrorCodeMW0010: InputErrorMessage,
	DebitCardFetchErrorCodeMW0030: InputErrorMessage,
}

var DebitCardFetchErrorCodeRetry = map[string]string{
	DebitCardFetchErrorCodePX1103: RetryErrorMessage,
}

func GetDebitCardFetchErrorMessage(errorCode string) (string, bool) {
	message, exists := DebitCardFetchErrorMessages[errorCode]
	return message, exists
}

func GetDebitCardFetchRetryErrorMessage(errorCode string) (string, bool) {
	message, exists := DebitCardFetchErrorCodeRetry[errorCode]
	return message, exists
}

const (
	PhysicalDebitCardErrorCodeMW0031 = "MW0031"
	PhysicalDebitCardErrorCodeMW0030 = "MW0030"
)

var PhysicalDebitCardErrorMessages = map[string]string{
	PhysicalDebitCardErrorCodeMW0031: InputErrorMessage,
	PhysicalDebitCardErrorCodeMW0030: InputErrorMessage,
}

func GetPhysicalDebitCardErrorMessage(errorCode string) (string, bool) {
	message, exists := PhysicalDebitCardErrorMessages[errorCode]
	return message, exists
}

const (
	OTPErrorCodeMW0029 = "MW0029"
	OTPErrorCodeMW0028 = "MW0028"
	OTPErrorCodeMW0030 = "MW0030"
	OTPErrorCode90     = "90"
)

var OTPErrorMessages = map[string]string{
	OTPErrorCode90: "OTP Expired. Please initiate Resent OTP.",
}

var OTPRetryErrorMessages = map[string]string{
	OTPErrorCodeMW0029: RetryErrorMessage,
	OTPErrorCodeMW0028: RetryErrorMessage,
	OTPErrorCodeMW0030: RetryErrorMessage,
}

func GetOTPErrorMessage(errorCode string) (string, bool) {
	message, exists := OTPErrorMessages[errorCode]
	return message, exists
}

func GetOTPRetryErrorMessage(errorCode string) (string, bool) {
	message, exists := OTPRetryErrorMessages[errorCode]
	return message, exists
}

const (
	DebitCardSetPinErrorCodeMW0030 = "MW0030"
	DebitCardSetPinErrorCodeMW0053 = "MW0053"
	DebitCardSetPinErrorCodeMW0054 = "MW0054"
	DebitCardSetPinErrorCodeMW0055 = "MW0055"
	DebitCardSetPinErrorCodeMW0056 = "MW0056"
	DebitCardSetPinErrorCodeMW0057 = "MW0057"
	DebitCardSetPinErrorCodeMW0052 = "MW0052"
	DebitCardSetPinErrorCodeMW0051 = "MW0051"
	DebitCardSetPinErrorCodeMW0050 = "MW0050"
	DebitCardSetPinErrorCodeMW0049 = "MW0049"
)

var DebitCardSetPinRetryErrorMessages = map[string]string{
	DebitCardSetPinErrorCodeMW0030: RetryErrorMessage,
	DebitCardSetPinErrorCodeMW0053: RetryErrorMessage,
	DebitCardSetPinErrorCodeMW0054: RetryErrorMessage,
	DebitCardSetPinErrorCodeMW0055: RetryErrorMessage,
	DebitCardSetPinErrorCodeMW0056: RetryErrorMessage,
	DebitCardSetPinErrorCodeMW0057: RetryErrorMessage,
	DebitCardSetPinErrorCodeMW0052: RetryErrorMessage,
	DebitCardSetPinErrorCodeMW0051: RetryErrorMessage,
	DebitCardSetPinErrorCodeMW0050: RetryErrorMessage,
}

var DebitCardSetPinErrorMessages = map[string]string{
	DebitCardSetPinErrorCodeMW0049: RetryErrorMessage,
}

func GetDebitCardSetPinRetryErrorMessage(errorCode string) (string, bool) {
	message, exists := DebitCardSetPinRetryErrorMessages[errorCode]
	return message, exists
}

func GetDebitCardSetPinErrorMessage(errorCode string) (string, bool) {
	message, exists := DebitCardSetPinErrorMessages[errorCode]
	return message, exists
}

const (
	UpiErrorMessageNoRecordsFound = "no records found"
	UpiErrorMessageNoDataFound    = "no data found"
	UpiErrorMessageInvalidMpin    = "invalid mpin"
)

const (
	AddressUpdateErrorCode91      = "91"
	AddressUpdateErrorMessage91   = "your request is already in pending status"
)
