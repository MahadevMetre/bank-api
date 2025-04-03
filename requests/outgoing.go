package requests

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"regexp"
	"strings"
	"time"

	"bitbucket.org/paydoh/paydoh-commons/customvalidation"

	"bankapi/constants"
	"bankapi/responses"
	"bankapi/security"
)

type OutgoingSimVerificationRequest struct {
	ServiceBy     string `json:"ServiceBy" validate:"required"`
	TxnIdentifier string `json:"TxnIdentifier" validate:"required"`
	ApplicantId   string `json:"ApplicantId" validate:"required"`
	MobileNo      string `json:"MobileNo" validate:"required,numeric"`
	ServiceStatus string `json:"ServiceStatus" validate:"required"`
}

type OutgoingSmsVerificationRequest struct {
	ServiceBy     string `json:"ServiceBy" validate:"required"`
	TxnIdentifier string `json:"TxnIdentifier" validate:"required"`
	ApplicantId   string `json:"ApplicantId" validate:"required"`
	MobileNo      string `json:"MobileNo" validate:"required,numeric"`
	ServiceStatus string `json:"ServiceStatus" validate:"required"`
}

type OutGoingVcipInvokeRequest struct {
	ApplicantId   string `json:"ApplicantId" validate:"required,max=20"`
	MobileNumber  string `json:"MobileNo" validate:"required,numeric,max=10"`
	TxnIdentifier string `json:"TxnIdentifier" validate:"required,max=25"`
}

type OutgoingDemographicRequest struct {
	ApplicantId   string `json:"ApplicantId"`
	MobileNo      string `json:"MobileNo"`
	TxnIdentifier string `json:"TxnIdentifier"`
}

type OutgoingCreateBankAccountRequest struct {
	ApplicationID string         `json:"ApplicationId"`
	UserID        string         `json:"UserId"`
	CustomerDtls  []CustomerDtls `json:"Customer_Dtls"`
}
type CustomerDtls struct {
	ApplicantId        string `json:"ApplicantId"`
	Nationality        string `json:"Nationality"`
	NamePrefix         string `json:"NamePrefix"`
	MotherMaidenName   string `json:"MotherMaidenName"`
	MiddleName         string `json:"MiddleName"`
	MaritalStatus      string `json:"MaritalStatus"`
	MailState          string `json:"MailState" validate:"required"`
	MailPin            string `json:"MailPin"`
	AnnualTurnOver     string `json:"AnualTurnOver"`
	MailEmail          string `json:"MailEmail"`
	MailCtry           string `json:"MailCtry"`
	MailCity           string `json:"MailCity"`
	MailAddr1          string `json:"MailAddr1"`
	LastName           string `json:"LastName"`
	FatherName         string `json:"FatherName"`
	CustomerEducation  string `json:"CustomerEducation"`
	CountryResidence   string `json:"CountryResidence"`
	FirstName          string `json:"FirstName"`
	Gender             string `json:"Gender"`
	Dob                string `json:"DOB"`
	PermCustAddress1   string `json:"PermCustAddress1"`
	PermCustCity       string `json:"PermCustCity"`
	PermCustState      string `json:"PermCustState" validate:"required"`
	PermCustCountry    string `json:"PermCustCountry"`
	PermZipcode        string `json:"PermZipcode"`
	AgeCategory        string `json:"AgeCategory"`
	ParentMobNo        string `json:"ParentMobNo"`
	ParentRelationship string `json:"ParentRelationship"`
	PrimaryApplicant   string `json:"PrimaryApplicant"`
	ProfessionCode     string `json:"ProfessionCode"`
}

type OutgoingAddNomineeRequest struct {
	ApplicantId             string `json:"ApplicantId" validate:"required"`
	AccountNo               string `json:"AccountNo" validate:"required"`
	TxnIdentifier           string `json:"TxnIdentifier" validate:"required"`
	NomReqType              string `json:"NomReqType" validate:"required"`
	NomApplId               string `json:"NomApplId" validate:"required"`
	NomName                 string `json:"NomName" validate:"required"`
	NomRelation             string `json:"NomRelation" validate:"required"`
	NomDOB                  string `json:"NomDOB" validate:"required"`
	NomAddressL1            string `json:"NomAddressL1" validate:"required"`
	NomAddressL2            string `json:"NomAddressL2" validate:"required"`
	NomAddressL3            string `json:"NomAddressL3" validate:"required"`
	NomCity                 string `json:"NomCity" validate:"required"`
	NomState                string `json:"NomState" validate:"required"`
	NomCountry              string `json:"NomCountry" validate:"required"`
	NomZipcode              string `json:"NomZipcode" validate:"required"`
	GuardianName            string `json:"GuardianName,omitempty"`
	GuardianNomineeRelation string `json:"GuardianNomineeRelation,omitempty"`
	GuardianAddressL1       string `json:"GuardianAddressL1,omitempty"`
	GuardianAddressL2       string `json:"GuardianAddressL2,omitempty"`
	GuardianAddressL3       string `json:"GuardianAddressL3,omitempty"`
	GuardianCity            string `json:"GuardianCity,omitempty"`
	GuardianState           string `json:"GuardianState,omitempty"`
	GuardianCountry         string `json:"GuardianCountry,omitempty"`
	GuardianZipcode         string `json:"GuardianZipcode,omitempty"`
	ResendOtp               string `json:"ResendOtp"` // Mandatory. Y/N.
	RetryFlag               string `json:"RetryFlag"` // Mandatory. Y/N.
	OTP                     string `json:"Otp"`
}

type OutgoingNomineeOtpGenerationRequest struct {
	ApplicantId string `json:"ApplicantId"`
}

type OutgoingFetchNomineeRequest struct {
	ApplicantId   string `json:"ApplicantId" validate:"required"`
	AccountNo     string `json:"AccountNo" validate:"required"`
	TxnIdentifier string `json:"TxnIdentifier" validate:"required"`
}

type OutgoingSyncRequest struct {
	TxnIdentifier string `json:"TxnIdentifier" validate:"required"`
	FromDate      string `json:"FromDate" validate:"required"`
	ToDate        string `json:"ToDate" validate:"required"`
}

type OutgoingBeneficiarySearchRequest struct {
	ApplicantId   string `json:"ApplicantId" validate:"required"`
	TxnIdentifier string `json:"TxnIdentifier" validate:"required"`
	AccountNo     string `json:"AccountNo" validate:"required"`
}

type OutgoingAddBeneficiaryRequest struct {
	ApplicantId   string `json:"ApplicantId" validate:"required"`   // Unique identifier for the applicant. Mandatory.
	TxnIdentifier string `json:"TxnIdentifier" validate:"required"` // Unique identification of the request. Mandatory.
	AccountNo     string `json:"AccountNo" validate:"required"`     // Mandatory.
	BenfID        string `json:"BenfId" validate:"required"`        // Mandatory. Account No of the Applicant.
	BenfName      string `json:"BenfName" validate:"required"`      // Mandatory. Name of beneficiary.
	BenfIFSC      string `json:"BenfIfsc" validate:"required"`      // Mandatory. IFSC code of beneficiary.
	BenfAcctNo    string `json:"BenfAcctNo" validate:"required"`    // Mandatory. Beneficiary Account No.
	BenfAcctNo1   string `json:"BenfAcctNo1" validate:"required"`   // To be reconfirmed.
	BenfAcctType  string `json:"BenfAcctType" validate:"required"`  // Mandatory. Account type of beneficiary like Savings/Current etc.
	BenfMobNo     string `json:"BenfMobNo" validate:"required"`     // Mandatory. Mobile No of beneficiary.
	PaymentMode   string `json:"PaymentMode" validate:"required"`   // Mandatory. Payment mode like NEFT, IMPS, BOTH(If both IMPS and NEFT required).
	ResendOtp     string `json:"ResendOtp" validate:"required"`     // Mandatory. Y/N.
	RetryFlag     string `json:"RetryFlag" validate:"required"`     // Mandatory. Y/N.
	OTP           string `json:"Otp"`                               // Tag need to be passed mandatorily, even if empty. To be passed when applicant submits OTP.
}

type OutgoingBeneficiaryOtpRequest struct {
	ApplicantId   string `json:"ApplicantId"`
	TxnIdentifier string `json:"TxnIdentifier"`
	AccountNo     string `json:"AccountNo"`
	Otp           string `json:"Otp"`
}

type OutgoingPaymentRequest struct {
	ApplicantId   string `json:"ApplicantId" validate:"required,max=20"`
	TxnIdentifier string `json:"TxnIdentifier" validate:"required,max=25"`
	PaymentMode   string `json:"PaymentMode" validate:"required,oneof=NEFT IMPS IFT"`
	AccountNo     string `json:"AccountNo" validate:"required,len=16"`
	BenfId        string `json:"BenfId" validate:"required,max=20"`
	BenfName      string `json:"BenfName" validate:"required,max=100"`
	BenfIfsc      string `json:"BenfIfsc" validate:"required,len=11"`
	BenfAcctNo    string `json:"BenfAcctNo" validate:"required,len=25"`
	BenfAcctType  string `json:"BenfAcctType" validate:"required,oneof=SA CA"`
	BenfMobNo     string `json:"BenfMobNo" validate:"omitempty,len=10"`
	Amount        string `json:"Amount" validate:"required"`
	TxnRemarks    string `json:"TxnRemarks" validate:"required,max=20"`
	ResendOtp     string `json:"ResendOtp" validate:"required,oneof=Y N"`
	RetryFlag     string `json:"RetryFlag" validate:"omitempty,oneof=Y N"`
	Otp           string `json:"Otp"`
	QuickTransfer string `json:"QuickTransfer" validate:"required,oneof=Y N"`
}

type OutgoingPaymentRequestOTP struct {
	ApplicantId   string `json:"ApplicantId" validate:"required,max=20"`
	TxnIdentifier string `json:"TxnIdentifier" validate:"required,max=25"`
	AccountNo     string `json:"AccountNo" validate:"required,len=16"`
	Otp           string `json:"Otp"`
}

type OutgoingConsentRequest struct {
	ApplicantId     string `json:"ApplicantId" validate:"required,max=20"`
	TxnIdentifier   string `json:"TxnIdentifier" validate:"required,max=25"`
	ConsentType     string `json:"ConsentType" validate:"required,max=25"`
	ConsentProvided string `json:"ConsentProvided" validate:"required,max=3"`
	ConsentTime     string `json:"ConsentTime" validate:"required"`
}

type OutgoingVerifyNomineeOTP struct {
	ApplicantId   string `json:"ApplicantId" validate:"required"`
	TxnIdentifier string `json:"TxnIdentifier" validate:"required"`
	AccountNo     string `json:"AccountNo" validate:"required"`
	NomReqType    string `json:"NomReqType" validate:"required"`
	NomAppId      string `json:"NomApplId" validate:"required"`
	Otp           string `json:"Otp" validate:"required"`
}

type OutBeneficiaryTemplateRequest struct {
	ApplicantId   string `json:"ApplicantId" validate:"required"`
	TxnIdentifier string `json:"TxnIdentifier" validate:"required"`
	AccountNo     string `json:"AccountNo" validate:"required"`
	BenfId        string `json:"BenfId" validate:"required"`
}

type OutgoingStatementRequest struct {
	ApplicantId   string `json:"ApplicantId"`
	TxnIdentifier string `json:"TxnIdentifier"`
	AccountNo     string `json:"KVBAccountNumber"`
	FromDate      string `json:"FromDate"`
	ToDate        string `json:"ToDate"`
}

type OutgoingRewardTransactionRequest struct {
	ProductType      string                `json:"productType"`
	ApplicationId    string                `json:"ApplicationId" validate:"required"`
	FundTransferDtls []FundTransferDetails `json:"Fund_Transfer_Dtls"`
}

type FundTransferDetails struct {
	Amount           string `json:"Amount"`
	Remarks          string `json:"Remarks"`
	TranId           string `json:"TranId"`
	ApplicationDate  string `json:"ApplicationDate"`
	ReqsTime         string `json:"Reqs_Time"`
	RetryFlag        string `json:"RetryFlag"`
	PaymentMode      string `json:"Payment_Mode"`
	MobileNo         string `json:"Mobile_No"`
	AccountNo        string `json:"Account No"`
	BenAcctNo        string `json:"Ben_Acct_No"`
	BenName          string `json:"Ben_Name"`
	TxnDate          string `json:"Txn_Date"`
	RemitterMobileNo string `json:"Remitter_Mobile_No"`
}

type ConsentRequestV2 struct {
	ApplicantId    string           `json:"ApplicantId"`
	TxnIdentifier  string           `json:"TxnIdentifier"`
	ConsentDetails []ConsentDetails `json:"ConsentDetails"`
}

type ConsentDetails struct {
	ConsentType     string `json:"ConsentType"`
	ConsentProvided string `json:"ConsentProvided"`
	ConsentTime     string `json:"ConsentTime"`
}

func NewOutGoingDemographicRequest() *OutgoingDemographicRequest {
	return &OutgoingDemographicRequest{}
}

func NewOutgoingCreateAccountRequest() *OutgoingCreateBankAccountRequest {
	return &OutgoingCreateBankAccountRequest{}
}

func NewOutgoingAddNomineeRequest() *OutgoingAddNomineeRequest {
	return &OutgoingAddNomineeRequest{}
}

func NewOutgoingIFSCSyncRequest() *OutgoingSyncRequest {
	return &OutgoingSyncRequest{}
}

func NewOutgoingBeneficiarySearchRequest() *OutgoingBeneficiarySearchRequest {
	return &OutgoingBeneficiarySearchRequest{}
}

func NewOutgoingAddBeneficiaryRequest() *OutgoingAddBeneficiaryRequest {
	return &OutgoingAddBeneficiaryRequest{}
}

func NewOutgoingBeneficiaryOTPRequest() *OutgoingBeneficiaryOtpRequest {
	return &OutgoingBeneficiaryOtpRequest{}
}

func NewOutgoingPaymentRequest() *OutgoingPaymentRequest {
	return &OutgoingPaymentRequest{}
}

func NewOutgoingPaymentRequestOTP() *OutgoingPaymentRequestOTP {
	return &OutgoingPaymentRequestOTP{}
}

func NewOutBeneficiaryTemplateRequest() *OutBeneficiaryTemplateRequest {
	return &OutBeneficiaryTemplateRequest{}
}

func NewOutgoingStatementRequest() *OutgoingStatementRequest {
	return &OutgoingStatementRequest{}
}

func NewOutgoingVcipRequest() *OutGoingVcipInvokeRequest {
	return &OutGoingVcipInvokeRequest{}
}

func NewOutgoingSimVerificationRequest() *OutgoingSimVerificationRequest {
	return &OutgoingSimVerificationRequest{}
}

func NewOutgoingRewardsTransactionRequest() *OutgoingRewardTransactionRequest {
	return &OutgoingRewardTransactionRequest{}
}

func NewOutgoingFetchNomineeRequest() *OutgoingFetchNomineeRequest {
	return &OutgoingFetchNomineeRequest{}
}

func (r *OutgoingSimVerificationRequest) Bind(userId, mobileNumber, status string) error {
	transactionId, err := security.GenerateRandomUUID(20)
	if err != nil {
		return err
	}
	r.ServiceBy = "PAYDOH"
	r.TxnIdentifier = transactionId
	r.MobileNo = mobileNumber
	r.ServiceStatus = status
	r.ApplicantId = userId

	return nil
}

func (r *OutGoingVcipInvokeRequest) Bind(applicantId, mobileNumber string) error {
	transactionId, err := security.GenerateRandomUUID(25)

	if err != nil {
		return err
	}

	r.ApplicantId = applicantId
	r.MobileNumber = mobileNumber
	r.TxnIdentifier = transactionId

	return nil
}

func (r *OutgoingDemographicRequest) Bind(applicantId, mobileNumber string) error {
	transactionId, err := security.GenerateRandomUUID(25)

	if err != nil {
		return err
	}

	r.ApplicantId = applicantId
	r.MobileNo = mobileNumber
	r.TxnIdentifier = transactionId

	return nil
}

func MapAddress(components ...string) string {
	var nonEmptyComponents []string

	for _, component := range components {
		trimmed := strings.TrimSpace(component)
		if trimmed != "" {
			nonEmptyComponents = append(nonEmptyComponents, trimmed)
		}
	}

	return strings.Join(nonEmptyComponents, ", ")
}

func (r *OutgoingCreateBankAccountRequest) Bind(
	userId,
	ApplicantId,
	email,
	middleName,
	mobileNumber string,
	request *CreateBankAccountRequest,
	response *responses.DemographicResponse,
	lastName string,
	dob string,
	stateCode string,
	IsAddrSameAsAdhaar bool,
) error {
	applicationId, err := security.GenerateRandomUUID(20)
	c := &CustomerDtls{}
	if err != nil {
		return err
	}

	uuidStr := strings.ReplaceAll(applicationId, "-", "")

	namePrefix := constants.GetNamePrefix(request.MaritalStatus, response.Root.UIDData.Poi.Gender)

	permanentAdr := MapAddress(
		response.Root.UIDData.Poa.Co,
		response.Root.UIDData.Poa.House,
		response.Root.Locality,
		response.Root.Street,
		response.Root.Landmark)

	r.ApplicationID = uuidStr
	c.ApplicantId = ApplicantId
	c.NamePrefix = namePrefix
	// c.ParentMobNo = mobileNumber //only for minor
	c.FirstName = response.FirstName
	c.MiddleName = response.MiddleName
	c.LastName = lastName
	c.AnnualTurnOver = request.AnnualTurnOver
	c.CountryResidence = request.CountryResidence
	c.CustomerEducation = request.CustomerEducation
	c.Gender = response.Root.UIDData.Poi.Gender
	c.Dob = dob
	c.AgeCategory = "A2"
	c.PrimaryApplicant = "Y"
	c.ProfessionCode = request.ProfessionCode

	c.MailEmail = email
	c.MotherMaidenName = request.MotherMaidenName
	c.MaritalStatus = request.MaritalStatus
	c.FatherName = ""
	c.Nationality = request.Nationality
	c.PermCustAddress1 = permanentAdr
	c.PermCustCountry = request.CountryResidence
	c.PermCustCity = response.Root.Vtc
	c.PermCustState = stateCode
	c.PermZipcode = response.Root.Pincode

	// communication address
	c.MailCtry = request.CountryResidence

	if !IsAddrSameAsAdhaar {
		c.MailAddr1 = MapAddress(
			request.CommunicationAddress.HouseNo,
			request.CommunicationAddress.StreetName,
			request.CommunicationAddress.Locality,
			request.CommunicationAddress.Landmark,
		)
		c.MailCity = request.CommunicationAddress.City
		c.MailState = stateCode
		c.MailPin = request.CommunicationAddress.PinCode
	} else {
		c.MailAddr1 = permanentAdr
		c.MailCity = strings.Split(response.Root.Vtc, " ")[0]
		c.MailState = stateCode
		c.MailPin = response.Root.Pincode
	}

	r.CustomerDtls = append(r.CustomerDtls, *c)
	r.UserID = "APIPAYDOH"

	return nil
}

func (r *OutgoingAddNomineeRequest) Bind(applicantId, nomApplicantId string, accountNumber string, nomineeCode string, request *AddNomineeRequest, txnId string) error {
	txnIdentifier := ""
	if txnId == "" {
		transactionId, err := security.GenerateRandomUUID(20)
		if err != nil {
			return err
		}
		txnIdentifier = strings.ReplaceAll(transactionId, "-", "")
	} else {
		txnIdentifier = txnId
	}

	if nomApplicantId == "" {
		nomApplId, err := security.GenerateRandomUUID(15)
		if err != nil {
			return err
		}

		nomIDStr := strings.ReplaceAll(nomApplId, "-", "")
		r.NomApplId = "nom" + nomIDStr
	} else {
		r.NomApplId = nomApplicantId
	}

	r.ApplicantId = applicantId
	r.TxnIdentifier = txnIdentifier //transactionId
	r.AccountNo = accountNumber
	r.NomAddressL1 = request.NomAddressL1
	r.NomAddressL2 = request.NomAddressL2
	r.NomAddressL3 = request.NomAddressL3

	r.NomCity = request.NomCity
	r.NomState = request.NomState
	r.NomCountry = request.NomCountry
	r.NomZipcode = request.NomZipcode
	r.NomName = request.NomName
	r.NomRelation = nomineeCode
	r.NomReqType = request.NomReqType
	r.NomDOB = request.NomDOB

	if request.GuardianName != "" {
		r.GuardianName = request.GuardianName
	}

	if request.GuardianNomineeRelation != "" {
		r.GuardianNomineeRelation = request.GuardianNomineeRelation
	}

	if request.GuardianAddressL1 != "" {
		r.GuardianAddressL1 = request.GuardianAddressL1
	}

	if request.GuardianAddressL2 != "" {
		r.GuardianAddressL2 = request.GuardianAddressL2
	}

	if request.GuardianAddressL3 != "" {
		r.GuardianAddressL3 = request.GuardianAddressL3
	}

	if request.GuardianCity != "" {
		r.GuardianCity = request.GuardianCity
	}

	if request.GuardianState != "" {
		r.GuardianState = request.GuardianState
	}

	if request.GuardianCountry != "" {
		r.GuardianCountry = request.GuardianCountry
	}

	if request.GuardianZipcode != "" {
		r.GuardianZipcode = request.GuardianZipcode
	}

	r.ResendOtp = request.ResendOtp
	r.RetryFlag = request.RetryFlag
	if request.Otp != "" {
		r.OTP = request.Otp
	} else {
		r.OTP = ""
	}
	return nil
}

func (r *OutgoingSyncRequest) Bind(startdate, enddate string) error {
	transactionId, err := security.GenerateRandomUUID(25)

	if err != nil {
		return err
	}

	r.TxnIdentifier = transactionId
	r.FromDate = startdate
	r.ToDate = enddate

	return nil
}

func (r *OutgoingBeneficiarySearchRequest) Bind(applicantId, accountNumber string) error {
	transactionId, err := security.GenerateRandomUUID(15)

	if err != nil {
		return err
	}

	r.TxnIdentifier = transactionId
	r.ApplicantId = applicantId
	r.AccountNo = accountNumber

	return nil
}

func (r *OutgoingRewardTransactionRequest) BindAndValidate(userId, amount, remarks, accountNumber, mobileNumber, beneficiaryName, transactionId string) error {
	r.ApplicationId = userId
	r.ProductType = "FINTECH"
	fundTransferDetail := make([]FundTransferDetails, 0)

	fundTransferDetail = append(fundTransferDetail, FundTransferDetails{
		TranId:           transactionId,
		Amount:           amount,
		ReqsTime:         time.Now().Format("02/01/2006 15:04:05"),
		Remarks:          remarks,
		MobileNo:         mobileNumber,
		ApplicationDate:  time.Now().Format("02/01/2006"),
		AccountNo:        os.Getenv("PAYDOH_ACCOUNT_NUMBER"),
		PaymentMode:      "IFT",
		BenAcctNo:        accountNumber,
		TxnDate:          time.Now().Format("02/01/2006"),
		RemitterMobileNo: mobileNumber,
		BenName:          beneficiaryName,
		RetryFlag:        "N",
	})

	r.FundTransferDtls = fundTransferDetail

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}

	return nil
}

func (r *OutgoingSimVerificationRequest) Validate() error {
	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}

	return nil
}

func (r *OutGoingVcipInvokeRequest) Validate() error {
	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}

	return nil
}

func NewOutgoingVerifyNomineeOTP() *OutgoingVerifyNomineeOTP {
	return &OutgoingVerifyNomineeOTP{}
}

func NewOutgoingConsentRequest() *OutgoingConsentRequest {
	return &OutgoingConsentRequest{}
}

func (r *OutgoingDemographicRequest) Validate() error {
	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}

	return nil
}

func (r *OutgoingAddNomineeRequest) Validate() error {
	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}

	return nil
}

func (r *OutgoingVerifyNomineeOTP) Validate() error {
	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}

	return nil
}

func (r *OutgoingVerifyNomineeOTP) Bind(nomineeApplId, applicantId string, accountNumber string, request *VerifyOtpRequest) error {
	r.ApplicantId = applicantId
	r.AccountNo = accountNumber
	r.NomReqType = request.ReqType
	r.NomAppId = nomineeApplId
	r.Otp = request.Otp

	return nil
}

func (r *OutgoingAddBeneficiaryRequest) Bind(applicantId, accountNumber string, request *AddNewBeneficiary) error {

	if request.TxnIdentifier == "" {
		transactionId, err := security.GenerateRandomUUID(15)
		if err != nil {
			return err
		}
		r.TxnIdentifier = transactionId
	} else {
		r.TxnIdentifier = request.TxnIdentifier
	}

	r.ApplicantId = applicantId
	r.AccountNo = accountNumber
	r.BenfID = request.BenfNickName
	r.BenfName = request.BenfName
	r.BenfAcctNo = request.BenfAcctNo
	r.BenfIFSC = request.BenfIFSC
	r.BenfAcctNo1 = request.BenfAcctNo1
	r.BenfAcctType = request.BenfAcctType
	r.BenfMobNo = request.BenfMobNo
	r.PaymentMode = request.PaymentMode
	r.ResendOtp = request.ResendOtp
	r.RetryFlag = request.RetryFlag
	if r.OTP != "" {
		r.OTP = request.OTP
	}

	r.OTP = request.OTP
	return nil
}

func (r *OutgoingBeneficiaryOtpRequest) Bind(applicantId string, accountNumber string, otp string) error {
	r.ApplicantId = applicantId
	r.AccountNo = accountNumber
	r.Otp = otp
	return nil
}

func (r *OutgoingPaymentRequest) Bind(applicantId, accountNumber string, request *PaymentRequest) error {
	transactionId, err := security.GenerateRandomUUID(15)
	if err != nil {
		return err
	}

	r.ApplicantId = applicantId
	r.AccountNo = accountNumber
	r.TxnRemarks = request.Remarks
	r.Amount = request.Amount
	r.BenfAcctNo = request.BenfAcctNo
	r.BenfAcctType = request.BenfAcctType
	r.BenfName = request.BenfName
	r.TxnIdentifier = transactionId
	r.BenfMobNo = request.BenfMobNo
	r.BenfIfsc = request.BenfIfsc
	r.BenfId = request.BenfId
	r.ResendOtp = request.ResendOtp
	r.RetryFlag = request.RetryFlag
	r.PaymentMode = request.PaymentMode
	r.QuickTransfer = request.QuickTransfer
	if r.Otp != "" {
		r.Otp = request.Otp
	}

	r.Otp = request.Otp

	return nil
}

func (r *OutgoingPaymentRequestOTP) Bind(applicantId, accountNumber, otp string) error {

	r.ApplicantId = applicantId
	r.AccountNo = accountNumber
	r.Otp = otp
	return nil
}

func (r *OutgoingConsentRequest) BindAndValidate(applicantId string) error {
	transactionId, err := security.GenerateRandomUUID(25)

	if err != nil {
		return err
	}

	r.ApplicantId = applicantId
	r.ConsentType = "PEP"
	r.ConsentProvided = "Yes"
	r.TxnIdentifier = transactionId

	currentTime := time.Now()
	formattedTime := currentTime.Format("02/Jan/2006 15:04:05")
	r.ConsentTime = formattedTime

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}

	return nil
}

func (r *OutgoingStatementRequest) BindAndValidate(userId, accountNumber, fromDate, toDate string) error {
	transactionId, err := security.GenerateRandomUUID(25)

	if err != nil {
		return err
	}

	r.TxnIdentifier = transactionId
	r.ApplicantId = userId
	r.AccountNo = accountNumber
	r.FromDate = fromDate
	r.ToDate = toDate

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}

	return nil
}

func (r *OutBeneficiaryTemplateRequest) BindAndValidate(applicantId, accountNumber, beneficiaryId string) error {
	r.AccountNo = accountNumber
	r.ApplicantId = applicantId
	r.BenfId = beneficiaryId

	return nil
}

func (r *OutgoingSimVerificationRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutGoingVcipInvokeRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutgoingDemographicRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutgoingSimVerificationRequest) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *OutGoingVcipInvokeRequest) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *OutgoingDemographicRequest) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *OutgoingCreateBankAccountRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutgoingCreateBankAccountRequest) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *OutgoingAddNomineeRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutgoingAddNomineeRequest) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *OutgoingVerifyNomineeOTP) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutgoingVerifyNomineeOTP) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *OutgoingSyncRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutgoingSyncRequest) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *OutgoingBeneficiarySearchRequest) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *OutgoingBeneficiarySearchRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutgoingAddBeneficiaryRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutgoingAddBeneficiaryRequest) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *OutgoingBeneficiaryOtpRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutgoingBeneficiaryOtpRequest) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *OutgoingPaymentRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutgoingPaymentRequest) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *OutgoingPaymentRequestOTP) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutgoingPaymentRequestOTP) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *OutgoingConsentRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *ConsentRequestV2) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutgoingConsentRequest) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *OutgoingStatementRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutgoingStatementRequest) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *OutgoingRewardTransactionRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutgoingRewardTransactionRequest) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *OutBeneficiaryTemplateRequest) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *OutBeneficiaryTemplateRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type OutgoingMobileMappingType0ApiRequest struct {
	MobileMapping MobileMappingType0 `json:"MobileMapping"`
}
type MobileMappingType0 struct {
	Type      string `json:"Type"`
	DeviceID  string `json:"DeviceID"`
	DeviceIP  string `json:"DeviceIP"`
	ChannelId string `json:"CHANNELID"`
}

func NewOutgoingMobileMappingType0ApiRequest() *OutgoingMobileMappingType0ApiRequest {
	return &OutgoingMobileMappingType0ApiRequest{}
}

func (r *OutgoingMobileMappingType0ApiRequest) Bind(deviceId, deviceIp string) error {

	r.MobileMapping.Type = "0"
	r.MobileMapping.DeviceID = deviceId
	r.MobileMapping.DeviceIP = deviceIp
	r.MobileMapping.ChannelId = "1"

	return nil
}

func (r *OutgoingMobileMappingType0ApiRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutgoingMobileMappingType0ApiRequest) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type OutgoingVerifyUserApiRequest struct {
	VerifyUser VerifyUserApi `json:"VerifyUser"`
}

type VerifyUserApi struct {
	Encryptinfo string `json:"Encryptinfo"`
	DeviceID    string `json:"DeviceID"`
	DeviceIP    string `json:"DeviceIP"`
	TransID     string `json:"TransID"`
	UniqueID    string `json:"UniqueID"`
	OSType      string `json:"OSType"`
	ChannelId   string `json:"CHANNELID"`
}

func NewOutgoingVerifyUserApiRequest() *OutgoingVerifyUserApiRequest {
	return &OutgoingVerifyUserApiRequest{}
}

func (r *OutgoingVerifyUserApiRequest) Bind(deviceIp, deviceId, os, transactionId, clientId string) error {
	// uniqueId, err := security.GenerateRandomUUID(16)

	// if err != nil {
	// 	return err
	// }

	encryptInfo := fmt.Sprintf("%s|%s|%s|%s|%s", deviceId, deviceIp, transactionId, clientId, os)

	r.VerifyUser.UniqueID = clientId
	r.VerifyUser.DeviceID = deviceId
	r.VerifyUser.DeviceIP = deviceIp
	r.VerifyUser.TransID = transactionId
	r.VerifyUser.OSType = strings.ToUpper(os)
	r.VerifyUser.Encryptinfo = encryptInfo
	r.VerifyUser.ChannelId = "1"

	return nil
}

func (r *OutgoingVerifyUserApiRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutgoingVerifyUserApiRequest) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type OutgoingMobileMappingType1ApiRequest struct {
	MobileMapping MobileMappingType1 `json:"MobileMapping"`
}

type MobileMappingType1 struct {
	Type      string `json:"Type"`
	DeviceID  string `json:"DeviceID"`
	TransID   string `json:"TransID"`
	DeviceIP  string `json:"DeviceIP"`
	ClientID  string `json:"ClientID"`
	OSType    string `json:"OSType"`
	ChannelId string `json:"CHANNELID"`
}

func NewOutgoingMobileMappingType1ApiRequest() *OutgoingMobileMappingType1ApiRequest {
	return &OutgoingMobileMappingType1ApiRequest{}
}

func (r *OutgoingMobileMappingType1ApiRequest) Bind(deviceId, deviceIp, os, transactionId, clientId string) error {
	// clientId, err := security.GenerateRandomUUID(12)

	// if err != nil {
	// 	return err
	// }

	r.MobileMapping.Type = "1"
	r.MobileMapping.DeviceID = deviceId
	r.MobileMapping.TransID = transactionId
	r.MobileMapping.DeviceIP = deviceIp
	r.MobileMapping.ClientID = clientId
	r.MobileMapping.OSType = strings.ToUpper(os)
	r.MobileMapping.ChannelId = "1"

	return nil
}

func (r *OutgoingMobileMappingType1ApiRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutgoingMobileMappingType1ApiRequest) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type OutgoingLCValidatorApiRequest struct {
	Request string `json:"Request"`
}

func NewOutgoingLCValidatorApiRequest() *OutgoingLCValidatorApiRequest {
	return &OutgoingLCValidatorApiRequest{}
}

func (r *OutgoingLCValidatorApiRequest) Bind(mobileNumber, transactionId, deviceId, deviceIp, userName string) {
	r.Request = fmt.Sprintf("%s&%s&%s&%s&%s&%s", mobileNumber, transactionId, deviceId, deviceIp, userName, "123456")
}

func (r *OutgoingLCValidatorApiRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutgoingLCValidatorApiRequest) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type OutgoingProfileCreationApiRequest struct {
	ProfileCreation ProfileCreationApi `json:"ProfileCreation"`
}

type ProfileCreationApi struct {
	NAME               string `json:"NAME"`
	MOBILE             string `json:"MOBILE"`
	GENDER             string `json:"GENDER"`
	MAILID             string `json:"MAILID"`
	ADDRESS            string `json:"ADDRESS"`
	CITY               string `json:"CITY"`
	STATE              string `json:"STATE"`
	COUNTRY            string `json:"COUNTRY"`
	AADHAR             string `json:"AADHAR"`
	CryptoInfo         string `json:"CryptoInfo"`
	PROFILEPHOTO       string `json:"PROFILEPHOTO"`
	FileName           string `json:"FileName"`
	DeviceThumbInfo    string `json:"DeviceThumbInfo"`
	Notification_RegID string `json:"Notification_RegID"`
	ChannelId          string `json:"CHANNELID"`
}

func NewOutgoingProfileCreationApiRequest() *OutgoingProfileCreationApiRequest {
	return &OutgoingProfileCreationApiRequest{}
}

type AddressRequest struct {
	Address string `json:"Address"`
	City    string `json:"City"`
	State   string `json:"State"`
	Country string `json:"Country"`
}

func NewAddressRequest() *AddressRequest {
	return &AddressRequest{}
}

func (r *AddressRequest) Bind(address, city, state string) error {
	r.Address = address
	r.City = city
	r.State = state
	r.Country = "INDIA"

	return nil
}

func (r *OutgoingProfileCreationApiRequest) Bind(
	mobileNumber, deviceId, name, email, Gender, cryptoInfo string,
	response *responses.DemographicResponse,
	p *AddressRequest,
	packageId string,
) error {

	thumbInfo := fmt.Sprint(packageId, "|", "91"+mobileNumber, "|", deviceId, "|", "NA")

	r.ProfileCreation.NAME = name
	r.ProfileCreation.MOBILE = "91" + mobileNumber
	r.ProfileCreation.GENDER = Gender
	r.ProfileCreation.MAILID = email
	r.ProfileCreation.ADDRESS = p.Address
	r.ProfileCreation.CITY = strings.Replace(p.City, ".", "", -1)
	r.ProfileCreation.STATE = p.State
	r.ProfileCreation.COUNTRY = p.Country
	r.ProfileCreation.AADHAR = "NA"
	r.ProfileCreation.CryptoInfo = cryptoInfo
	r.ProfileCreation.PROFILEPHOTO = "NA"
	r.ProfileCreation.FileName = "NA"
	r.ProfileCreation.DeviceThumbInfo = thumbInfo
	r.ProfileCreation.Notification_RegID = deviceId
	r.ProfileCreation.ChannelId = "1"

	return nil
}

func (r *OutgoingProfileCreationApiRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutgoingProfileCreationApiRequest) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type OutgoingReqListKeysApiRequest struct {
	ReqListKeys ReqListKeysApi `json:"ReqListKeys"`
}

type ReqListKeysApi struct {
	Type        string `json:"Type"`
	CredType    string `json:"Cred_type"`
	CredSubType string `json:"Cred_subType"`
	CredCode    string `json:"Cred_code"`
	CredKi      string `json:"Cred_ki"`
	CredData    string `json:"Cred_Data"`
	DeviceIP    string `json:"DeviceIP"`
	ChannelId   string `json:"CHANNELID"`
}

func NewOutgoingReqListKeysApiRequest() *OutgoingReqListKeysApiRequest {
	return &OutgoingReqListKeysApiRequest{}
}

func (r *OutgoingReqListKeysApiRequest) Bind(deviceIp, deviceId, challengeType, packageId, mobileNumber, challenge string) error {
	input := fmt.Sprintf("%s|%s|%s|%s", deviceId, packageId, "91"+mobileNumber, challenge)

	reqlistkeys := &ReqListKeysApi{
		Type:        "1",
		CredType:    "Challenge",
		CredSubType: challengeType,
		CredCode:    "NPCI",
		CredKi:      "20150822",
		CredData:    input,
		DeviceIP:    deviceIp,
		ChannelId:   "1",
	}

	r.ReqListKeys = *reqlistkeys

	return nil
}

func (r *OutgoingReqListKeysApiRequest) NewBind() error {
	reqlistkeys := &ReqListKeysApi{
		Type: "0",
	}

	r.ReqListKeys = *reqlistkeys

	return nil
}

func (r *OutgoingReqListKeysApiRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutgoingReqListKeysApiRequest) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type OutgoingRemappingApiRequest struct {
	ReMapping ReMappingApi `json:"ReMapping"`
}
type ReMappingApi struct {
	MobileNo           string `json:"MobileNo"`
	DeviceID           string `json:"DeviceId"`
	Dguuid             string `json:"DGUUID"`
	OSVersion          string `json:"OSVersion"`
	OSType             string `json:"OSType"`
	Psppin             string `json:"PSPPIN"`
	DeviceIP           string `json:"DeviceIP"`
	Notification_RegID string `json:"Notification_RegID"`
	ChannelId          string `json:"CHANNELID"`
}

func NewOutgoingRemappingApiRequest() *OutgoingRemappingApiRequest {
	return &OutgoingRemappingApiRequest{}
}

func (r *OutgoingRemappingApiRequest) Bind(
	mobileNumber, deviceId, deviceIp, osVersion, os, emailid, dguuid string,
) error {

	// dguuid, err := security.GenerateRandomUUID(16)

	// if err != nil {
	// 	return err
	// }

	r.ReMapping.MobileNo = "91" + mobileNumber
	r.ReMapping.DeviceID = deviceId
	r.ReMapping.Dguuid = dguuid
	r.ReMapping.OSVersion = strings.ToUpper(os) + " " + strings.ToUpper(osVersion)
	r.ReMapping.OSType = strings.ToUpper(os)
	r.ReMapping.Psppin = "0"
	r.ReMapping.DeviceIP = deviceIp
	r.ReMapping.Notification_RegID = deviceId
	r.ReMapping.ChannelId = "1"

	return nil
}

func (r *OutgoingRemappingApiRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutgoingRemappingApiRequest) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type OutgoingExistingReqlistkeysApiRequest struct {
	ReqListKeys ExistingReqListKeysApi `json:"ReqListKeys"`
}
type ExistingReqListKeysApi struct {
	Type        string `json:"Type"`
	CredType    string `json:"Cred_type"`
	CredSubType string `json:"Cred_subType"`
	CredCode    string `json:"Cred_code"`
	CredKi      string `json:"Cred_ki"`
	CredData    string `json:"Cred_Data"`
	DeviceIP    string `json:"DeviceIP"`
	ChannelId   string `json:"CHANNELID"`
}

func NewOutgoingExistingReqlistkeysApiRequest() *OutgoingExistingReqlistkeysApiRequest {
	return &OutgoingExistingReqlistkeysApiRequest{}
}

func (r *OutgoingExistingReqlistkeysApiRequest) Bind(deviceIp, deviceId, packageId, mobileNumber, challenge string) error {
	input := fmt.Sprintf("%s|%s|%s|%s", deviceId, packageId, "91"+mobileNumber, challenge)

	r.ReqListKeys.Type = "1"
	r.ReqListKeys.CredType = "Challenge"
	r.ReqListKeys.CredSubType = "initial"
	r.ReqListKeys.CredCode = "NPCI"
	r.ReqListKeys.CredKi = "20150822"
	r.ReqListKeys.CredData = input
	r.ReqListKeys.DeviceIP = deviceIp
	r.ReqListKeys.ChannelId = "1"

	return nil
}

func (r *OutgoingExistingReqlistkeysApiRequest) AlternativeBind(deviceIp, deviceId, packageId, mobileNumber, challenge string) error {
	input := fmt.Sprintf("%s|%s|%s|%s", deviceId, packageId, "91"+mobileNumber, challenge)

	r.ReqListKeys.Type = "1"
	r.ReqListKeys.CredType = "CHALLENGE"
	r.ReqListKeys.CredSubType = "ROTATE"
	r.ReqListKeys.CredCode = "NPCI"
	r.ReqListKeys.CredKi = "20150822"
	r.ReqListKeys.CredData = input
	r.ReqListKeys.DeviceIP = deviceIp
	r.ReqListKeys.ChannelId = "1"

	return nil
}

func (r *OutgoingExistingReqlistkeysApiRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutgoingExistingReqlistkeysApiRequest) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type OutgoingCreateupiidRequestListAccountApiRequest struct {
	ReqListAccount CreateUpiidRequestListAccountApi `json:"ReqListAccount"`
}

type CreateUpiidRequestListAccountApi struct {
	Linkvalue      string `json:"Linkvalue"`
	Payeraddr      string `json:"Payeraddr"`
	Payername      string `json:"Payername"`
	AccountIfsc    string `json:"ACCOUNT_IFSC"`
	CryptoInfo     string `json:"CryptoInfo"`
	DeviceIP       string `json:"DeviceIP"`
	ChannelId      string `json:"CHANNELID"`
	AadhaarConsent string `json:"AADHAARConsent"`
	DeviceApp      string `json:"DeviceApp"`
	DeviceGeocode  string `json:"DeviceGeocode"`
}

func NewOutgoingCreateupiidRequestListAccountApiRequest() *OutgoingCreateupiidRequestListAccountApiRequest {
	return &OutgoingCreateupiidRequestListAccountApiRequest{}
}

func (r *OutgoingCreateupiidRequestListAccountApiRequest) Bind(
	mobileNumber,
	deviceIp,
	userName,
	cryptoInfo,
	packageId,
	latLong string,
) error {
	var paddr string

	if len(userName) >= 3 {
		paddr = strings.ToLower(userName[:1]) + mobileNumber[:6]
	} else {
		paddr = strings.ToLower(userName) + mobileNumber[:3]
	}

	payeraddr := fmt.Sprintf("%s.paydoh@kvb", paddr)

	r.ReqListAccount.Linkvalue = "91" + mobileNumber
	r.ReqListAccount.Payeraddr = payeraddr
	r.ReqListAccount.Payername = "91" + mobileNumber
	r.ReqListAccount.CryptoInfo = cryptoInfo
	r.ReqListAccount.DeviceIP = deviceIp
	r.ReqListAccount.ChannelId = "1"
	r.ReqListAccount.AccountIfsc = "KVBL0001861"
	r.ReqListAccount.AadhaarConsent = "Y"
	r.ReqListAccount.DeviceApp = packageId
	r.ReqListAccount.DeviceGeocode = latLong

	return nil
}

func (r *OutgoingCreateupiidRequestListAccountApiRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutgoingCreateupiidRequestListAccountApiRequest) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type OutgoingPspAvailabilityApiRequest struct {
	PspAvailability PspAvailabilityApi `json:"PspAvailability"`
}
type PspAvailabilityApi struct {
	PayerAddr  string `json:"PayerAddr"`
	CryptoInfo string `json:"CryptoInfo"`
	ChannelId  string `json:"CHANNELID"`
}

func NewOutgoingPspAvailabilityApiRequest() *OutgoingPspAvailabilityApiRequest {
	return &OutgoingPspAvailabilityApiRequest{}
}

func (r *OutgoingPspAvailabilityApiRequest) Bind(
	deviceId,
	osversion,
	os,
	cryptoInfo string,
	req *OutgoingCreateupiidRequestListAccountApiRequest,
) error {

	r.PspAvailability.PayerAddr = req.ReqListAccount.Payeraddr
	r.PspAvailability.CryptoInfo = cryptoInfo
	r.PspAvailability.ChannelId = "1"

	return nil
}

func (r *OutgoingPspAvailabilityApiRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutgoingPspAvailabilityApiRequest) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type OutgoingAddBankApiRequest struct {
	AddBank AddBankApi `json:"AddBank"`
}
type AddBankApi struct {
	Mbeba            string `json:"mbeba"`
	PayerType        string `json:"PayerType"`
	OTPDType         string `json:"OTPDType"`
	MaskedAccNo      string `json:"MaskedAccNo"`
	PayerAddr        string `json:"PayerAddr"`
	PayerCode        string `json:"PayerCode"`
	DeviceMOBILE     string `json:"DeviceMOBILE"`
	OTPDLength       string `json:"OTPDLength"`
	DeviceLOCATION   string `json:"DeviceLOCATION"`
	Aeba             string `json:"aeba"`
	DeviceTYPE       string `json:"DeviceTYPE"`
	PayerSeqNum      string `json:"PayerSeqNum"`
	RegDetailsMOBILE string `json:"RegDetails_MOBILE"`
	DeviceCAPABILITY string `json:"DeviceCAPABILITY"`
	CryptoInfo       string `json:"CryptoInfo"`
	AccountAcnum     string `json:"ACCOUNT_ACNUM"`
	DeviceGEOCODE    string `json:"DeviceGEOCODE"`
	DeviceID         string `json:"DeviceID"`
	PspName          string `json:"PspName"`
	ATMDType         string `json:"ATMDType"`
	MPinDLength      string `json:"MPinDLength"`
	DeviceIP         string `json:"DeviceIP"`
	PayerName        string `json:"PayerName"`
	MPinDType        string `json:"MPinDType"`
	AccountType      string `json:"AccountType"`
	DeviceOS         string `json:"DeviceOS"`
	AccountActype    string `json:"ACCOUNT_ACTYPE"`
	AccountIfsc      string `json:"ACCOUNT_IFSC"`
	DeviceAPP        string `json:"DeviceAPP"`
	ATMDLength       string `json:"ATMDLength"`
	ChannelId        string `json:"CHANNELID"`
}

func NewOutgoingAddBankApiRequest() *OutgoingAddBankApiRequest {
	return &OutgoingAddBankApiRequest{}
}

func (r *OutgoingAddBankApiRequest) Bind(
	deviceIp,
	deviceId,
	Mobilenumber,
	osversion,
	os,
	name,
	latlog,
	cryptoInfo,
	location string,
	response *responses.CreateUpiIdRequestListAccountApiResponse,
	request *OutgoingCreateupiidRequestListAccountApiRequest,
	deviceCapability string,
	appId string,
	lastName string,
	accountNumber string,
) error {
	r.AddBank.PayerType = "PERSON"
	r.AddBank.PayerAddr = request.ReqListAccount.Payeraddr
	r.AddBank.PayerCode = "0000"
	r.AddBank.DeviceMOBILE = "91" + Mobilenumber
	r.AddBank.DeviceLOCATION = location
	r.AddBank.DeviceTYPE = "MOB"
	r.AddBank.PayerSeqNum = "1"
	r.AddBank.RegDetailsMOBILE = "91" + Mobilenumber
	r.AddBank.DeviceCAPABILITY = deviceCapability
	r.AddBank.CryptoInfo = cryptoInfo
	r.AddBank.DeviceGEOCODE = latlog
	r.AddBank.DeviceID = deviceId
	r.AddBank.PspName = "607100"
	r.AddBank.ATMDType = ""
	r.AddBank.DeviceIP = deviceIp
	r.AddBank.AccountType = "ACCOUNT"
	r.AddBank.DeviceOS = strings.ToUpper(os) + " " + strings.ToUpper(osversion)
	r.AddBank.DeviceAPP = appId
	r.AddBank.ATMDLength = ""
	r.AddBank.ChannelId = "1"

	for i := range response.Response.Response.Ns2RespListAccount.AccountList.Account {
		if response.Response.Response.Ns2RespListAccount.AccountList.Account[i].AccRefNumber == accountNumber {
			r.AddBank.Mbeba = response.Response.Response.Ns2RespListAccount.AccountList.Account[i].Mbeba
			r.AddBank.OTPDType = response.Response.Response.Ns2RespListAccount.AccountList.Account[i].CredsAllowed[1].SubType
			r.AddBank.MaskedAccNo = response.Response.Response.Ns2RespListAccount.AccountList.Account[i].MaskedAccnumber
			r.AddBank.Aeba = response.Response.Response.Ns2RespListAccount.AccountList.Account[i].Aeba
			r.AddBank.OTPDLength = response.Response.Response.Ns2RespListAccount.AccountList.Account[i].CredsAllowed[1].DLength
			r.AddBank.AccountAcnum = response.Response.Response.Ns2RespListAccount.AccountList.Account[i].AccRefNumber
			r.AddBank.MPinDLength = response.Response.Response.Ns2RespListAccount.AccountList.Account[i].CredsAllowed[0].DLength
			r.AddBank.PayerName = response.Response.Response.Ns2RespListAccount.AccountList.Account[i].Name
			r.AddBank.MPinDType = response.Response.Response.Ns2RespListAccount.AccountList.Account[i].CredsAllowed[0].Type
			r.AddBank.AccountActype = response.Response.Response.Ns2RespListAccount.AccountList.Account[i].AccType
			r.AddBank.AccountIfsc = response.Response.Response.Ns2RespListAccount.AccountList.Account[i].Ifsc

			break
		}
	}

	if len(r.AddBank.AccountAcnum) == 0 && r.AddBank.AccountIfsc == "" {
		r.AddBank.PayerName = fmt.Sprintf("%s %s", name, lastName)
		r.AddBank.AccountAcnum = accountNumber
		r.AddBank.AccountIfsc = "KVBL0001861"
		r.AddBank.AccountActype = "SAVINGS"
		r.AddBank.Aeba = "Y"
		r.AddBank.Mbeba = "N"
		r.AddBank.MaskedAccNo = maskAccountNumber(accountNumber)

		r.AddBank.OTPDType = response.Response.Response.Ns2RespListAccount.AccountList.Account[0].CredsAllowed[1].SubType
		r.AddBank.OTPDLength = response.Response.Response.Ns2RespListAccount.AccountList.Account[0].CredsAllowed[1].DLength
		r.AddBank.MPinDLength = response.Response.Response.Ns2RespListAccount.AccountList.Account[0].CredsAllowed[0].DLength
		r.AddBank.MPinDType = response.Response.Response.Ns2RespListAccount.AccountList.Account[0].CredsAllowed[0].Type
	}

	return nil
}

func maskAccountNumber(input string) string {
	re := regexp.MustCompile(`\b\d{10,28}\b`)

	return re.ReplaceAllStringFunc(input, func(match string) string {
		n := len(match)
		if n <= 4 {
			return match
		}
		return strings.Repeat("X", n-4) + match[n-4:]
	})
}

func (r *OutgoingAddBankApiRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutgoingAddBankApiRequest) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type OutgoingReqBalEnqApiRequest struct {
	ReqBalEnq ReqBalEnqApi `json:"ReqBalEnq"`
}
type ReqBalEnqApi struct {
	TxnID       string `json:"TxnID"`
	Payeraddr   string `json:"Payeraddr"`
	CredData    string `json:"Cred_Data"`
	CryptoInfo  string `json:"CryptoInfo"`
	GeoLocation string `json:"GeoLocation"`
	DeviceIP    string `json:"DeviceIP"`
	ChannelId   string `json:"CHANNELID"`
}

func NewOutgoingReqBalEnqApiRequest() *OutgoingReqBalEnqApiRequest {
	return &OutgoingReqBalEnqApiRequest{}
}

func (r *OutgoingReqBalEnqApiRequest) Bind(
	deviceIp,
	latlog,
	CryptoInfo,
	txnid,
	UpiPin string,
	req *OutgoingCreateupiidRequestListAccountApiRequest,
) error {

	r.ReqBalEnq.TxnID = txnid
	r.ReqBalEnq.Payeraddr = req.ReqListAccount.Payeraddr
	r.ReqBalEnq.CredData = UpiPin
	r.ReqBalEnq.CryptoInfo = CryptoInfo
	r.ReqBalEnq.GeoLocation = latlog
	r.ReqBalEnq.DeviceIP = deviceIp
	r.ReqBalEnq.ChannelId = "1"

	return nil
}

func (r *OutgoingReqBalEnqApiRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutgoingReqBalEnqApiRequest) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type OutgoingAadharRequestListAccountsApiRequest struct {
	ReqListAccount ReqListAccountApi `json:"ReqListAccount"`
}
type ReqListAccountApi struct {
	Linkvalue      string `json:"Linkvalue"`
	Payeraddr      string `json:"Payeraddr"`
	Payername      string `json:"Payername"`
	AccountIfsc    string `json:"ACCOUNT_IFSC"`
	CryptoInfo     string `json:"CryptoInfo"`
	DeviceIP       string `json:"DeviceIP"`
	ChannelId      string `json:"CHANNELID"`
	Aeba           string `json:"aeba"`
	AadhaarConsent string `json:"AADHAARConsent"`
	DeviceApp      string `json:"DeviceApp"`
	DeviceGeocode  string `json:"DeviceGeocode"`
}

func NewOutgoingAadharRequestListAccountsApiRequest() *OutgoingAadharRequestListAccountsApiRequest {
	return &OutgoingAadharRequestListAccountsApiRequest{}
}

func (r *OutgoingAadharRequestListAccountsApiRequest) Bind(
	mobileNumber,
	deviceIp,
	userName,
	CryptoInfo,
	payerAddr,
	packageId,
	latLong string,
) error {

	r.ReqListAccount.Linkvalue = "91" + mobileNumber
	r.ReqListAccount.Payeraddr = payerAddr
	r.ReqListAccount.Payername = "91" + mobileNumber
	r.ReqListAccount.CryptoInfo = CryptoInfo
	r.ReqListAccount.DeviceIP = deviceIp
	r.ReqListAccount.ChannelId = "1"
	r.ReqListAccount.AccountIfsc = "KVBL0001861"
	r.ReqListAccount.Aeba = "Y"
	r.ReqListAccount.AadhaarConsent = "Y"
	r.ReqListAccount.DeviceApp = packageId
	r.ReqListAccount.DeviceGeocode = latLong

	return nil
}

func (r *OutgoingAadharRequestListAccountsApiRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutgoingAadharRequestListAccountsApiRequest) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type OutgoingSetUpiPinReqOtpApiRequest struct {
	ReqOtp RequestOtpApi `json:"ReqOtp"`
}
type RequestOtpApi struct {
	MobileNo    string `json:"MobileNo"`
	Payeraddr   string `json:"Payeraddr"`
	CryptoInfo  string `json:"CryptoInfo"`
	DeviceIP    string `json:"DeviceIP"`
	GeoLocation string `json:"GeoLocation"`
	ChannelId   string `json:"CHANNELID"`
	FormatType  string `json:"FormatType"`
}

func NewOutgoingSetUpiPinReqOtpApiRequest() *OutgoingSetUpiPinReqOtpApiRequest {
	return &OutgoingSetUpiPinReqOtpApiRequest{}
}

func (r *OutgoingSetUpiPinReqOtpApiRequest) Bind(
	mobileNumber,
	deviceIp,
	latlog,
	CryptoInfo string,
	payerAddr string,
) error {

	r.ReqOtp.MobileNo = "91" + mobileNumber
	r.ReqOtp.Payeraddr = payerAddr
	r.ReqOtp.CryptoInfo = CryptoInfo
	r.ReqOtp.DeviceIP = deviceIp
	r.ReqOtp.GeoLocation = latlog
	r.ReqOtp.ChannelId = "1"
	r.ReqOtp.FormatType = "FORMAT3"

	return nil
}

func (r *OutgoingSetUpiPinReqOtpApiRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutgoingSetUpiPinReqOtpApiRequest) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type OutgoingSetUpiPinReqRegMobApiRequest struct {
	ReqRegMob RequestRegMobApi `json:"ReqRegMob"`
}
type RequestRegMobApi struct {
	MobileNo             string `json:"MobileNo"`
	TxnID                string `json:"TxnID"`
	Payeraddr            string `json:"Payeraddr"`
	CryptoInfo           string `json:"CryptoInfo"`
	RegDetailsEXPDATE    string `json:"RegDetails_EXPDATE"`
	RegDetailsCARDDIGITS string `json:"RegDetails_CARDDIGITS"`
	CredDataOTP          string `json:"Cred_DataOTP"`
	DeviceIP             string `json:"DeviceIP"`
	GeoLocation          string `json:"GeoLocation"`
	CredDataMPIN         string `json:"Cred_DataMPIN"`
	CredDataATMPIN       string `json:"Cred_DataATMPIN"`
	ChannelId            string `json:"CHANNELID"`
	Cred_AADHAAR         string `json:"Cred_AADHAAR"`
	FormatType           string `json:"FormatType"`
}

func NewOutgoingSetUpiPinReqRegMobApiRequest() *OutgoingSetUpiPinReqRegMobApiRequest {
	return &OutgoingSetUpiPinReqRegMobApiRequest{}
}

func (r *OutgoingSetUpiPinReqRegMobApiRequest) Bind(
	mobileNumber,
	deviceIp,
	latlog,
	CryptoInfo,
	txnid,
	Otp,
	UpiPin,
	AtmPin string,
	s *OutgoingCreateupiidRequestListAccountApiRequest,
	credAdhaar string,
) error {

	r.ReqRegMob.MobileNo = "91" + mobileNumber
	r.ReqRegMob.TxnID = txnid
	r.ReqRegMob.Payeraddr = s.ReqListAccount.Payeraddr
	r.ReqRegMob.CryptoInfo = CryptoInfo
	r.ReqRegMob.RegDetailsEXPDATE = ""
	r.ReqRegMob.RegDetailsCARDDIGITS = ""
	r.ReqRegMob.CredDataOTP = Otp
	r.ReqRegMob.DeviceIP = deviceIp
	r.ReqRegMob.GeoLocation = latlog
	r.ReqRegMob.CredDataMPIN = UpiPin
	r.ReqRegMob.CredDataATMPIN = AtmPin
	r.ReqRegMob.ChannelId = "1"
	r.ReqRegMob.Cred_AADHAAR = credAdhaar
	r.ReqRegMob.FormatType = "FORMAT3"

	return nil
}

func (r *OutgoingSetUpiPinReqRegMobApiRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutgoingSetUpiPinReqRegMobApiRequest) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type OutgoingReqValAddApiRequest struct {
	ReqValAdd ReqValAddApi `json:"ReqValAdd"`
}
type ReqValAddApi struct {
	Payeraddr           string `json:"Payeraddr"`
	Payername           string `json:"Payername"`
	PayerCode           string `json:"PayerCode"`
	PayerSeqNum         string `json:"PayerSeqNum"`
	Payertype           string `json:"Payertype"`
	Infoid              string `json:"Infoid"`
	Infotype            string `json:"Infotype"`
	InfoverifiedName    string `json:"InfoverifiedName"`
	InfoverifiedAddress string `json:"InfoverifiedAddress"`
	RPayeraddr          string `json:"RPayeraddr"`
	RPayername          string `json:"RPayername"`
	RPayerSeqNum        string `json:"RPayerSeqNum"`
	CryptoInfo          string `json:"CryptoInfo"`
	DeviceIP            string `json:"DeviceIP"`
	MobileNo            string `json:"MobileNo"`
	CHANNELID           string `json:"CHANNELID"`
	DeviceGeocode       string `json:"DeviceGeocode"`
}

func NewOutgoingReqValAddApiRequest() *OutgoingReqValAddApiRequest {
	return &OutgoingReqValAddApiRequest{}
}

func (r *OutgoingReqValAddApiRequest) Bind(
	deviceIp,
	name,
	CryptoInfo,
	mobilenumber string,
	req2 *OutgoingCreateupiidRequestListAccountApiRequest,
	request *ValidateVpaRequest,
	payerAddr,
	latLong string,
) error {

	r.ReqValAdd.Payeraddr = payerAddr
	r.ReqValAdd.Payername = name

	if request.PayerCode != "" {
		r.ReqValAdd.PayerCode = request.PayerCode
	} else {
		r.ReqValAdd.PayerCode = "0000"
	}

	r.ReqValAdd.PayerSeqNum = "1"
	r.ReqValAdd.Payertype = "PERSON"
	r.ReqValAdd.Infoid = "1"
	r.ReqValAdd.Infotype = "ACCOUNT"
	r.ReqValAdd.InfoverifiedName = name
	r.ReqValAdd.InfoverifiedAddress = "TRUE"
	r.ReqValAdd.CHANNELID = "1"

	mobileNumber := request.RPayeraddr
	if strings.Contains(mobileNumber, "@") {
		r.ReqValAdd.RPayeraddr = mobileNumber
	} else {
		formattedMobileNumber := fmt.Sprintf("%s@mapper.npci", mobileNumber)
		r.ReqValAdd.RPayeraddr = formattedMobileNumber
	}

	if request.RPayername != "" {
		r.ReqValAdd.RPayername = request.RPayername
	} else {
		r.ReqValAdd.RPayername = "Paydoh"
	}

	r.ReqValAdd.RPayerSeqNum = "1"
	r.ReqValAdd.CryptoInfo = CryptoInfo
	r.ReqValAdd.DeviceIP = deviceIp
	r.ReqValAdd.MobileNo = "91" + mobilenumber
	r.ReqValAdd.DeviceGeocode = latLong

	return nil
}

func (r *OutgoingReqValAddApiRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutgoingReqValAddApiRequest) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type OutgoingReqPayApiRequest struct {
	ReqPay ReqPayApi `json:"ReqPay"`
}

type ReqPayApi struct {
	MobileNo        string `json:"MobileNo"`
	TransactionType string `json:"TransactionType"`
	Payeeaddr       string `json:"Payeeaddr"`
	CryptoInfo      string `json:"CryptoInfo"`
	Purpose         string `json:"Purpose"`
	DeviceIP        string `json:"DeviceIP"`
	PayerAmount     string `json:"PayerAmount"`
	Refurl          string `json:"Refurl"`
	Payeename       string `json:"Payeename"`
	TxnID           string `json:"TxnID"`
	Payeraddr       string `json:"Payeraddr"`
	ExpiryDateTime  string `json:"ExpiryDateTime"`
	Remarks         string `json:"Remarks"`
	InitiationMode  string `json:"InitiationMode"`
	GeoLocation     string `json:"GeoLocation"`
	Expiry          string `json:"Expiry"`
	MCCCode         string `json:"MCCCode"`
	CredData        string `json:"Cred_Data"`
	ChannelId       string `json:"CHANNELID"`
}

func NewOutgoingReqPayApiRequest() *OutgoingReqPayApiRequest {
	return &OutgoingReqPayApiRequest{}
}

func (r *OutgoingReqPayApiRequest) Bind(
	mobileNumber,
	deviceIp,
	FirstName,
	latlog,
	CryptoInfo string,
	request *PayMoneyWithVpaRequest,
	payerAddr string,
) error {

	r.ReqPay.MobileNo = "91" + mobileNumber
	r.ReqPay.TransactionType = request.TransactionType

	mobileNumber = request.Payeeaddr
	if strings.Contains(mobileNumber, "@") {
		r.ReqPay.Payeeaddr = mobileNumber
	} else {
		formattedMobileNumber := fmt.Sprintf("%s@mapper.npci", mobileNumber)
		r.ReqPay.Payeeaddr = formattedMobileNumber
	}

	r.ReqPay.CryptoInfo = CryptoInfo
	r.ReqPay.Purpose = "00"
	r.ReqPay.DeviceIP = deviceIp
	r.ReqPay.PayerAmount = request.PayerAmount
	r.ReqPay.Refurl = "http://www.karurvysyabank.in/"
	r.ReqPay.Payeename = request.PayeeName
	r.ReqPay.TxnID = request.TransId
	r.ReqPay.Payeraddr = payerAddr
	r.ReqPay.ExpiryDateTime = ""
	r.ReqPay.Remarks = request.Remark
	r.ReqPay.InitiationMode = "00"
	r.ReqPay.GeoLocation = latlog
	r.ReqPay.Expiry = ""
	r.ReqPay.MCCCode = request.MccCode
	r.ReqPay.CredData = request.UpiPin
	r.ReqPay.ChannelId = "1"

	return nil
}

func (r *OutgoingReqPayApiRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutgoingReqPayApiRequest) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type OutgoingAccountLinkApiRequest struct {
	AccountLink AccountLinkApi `json:"AccountLink"`
}
type AccountLinkApi struct {
	MobileNo   string `json:"MobileNo"`
	CryptoInfo string `json:"CryptoInfo"`
	ChannelId  string `json:"CHANNELID"`
}

func NewOutgoingAccountLinkApiRequest() *OutgoingAccountLinkApiRequest {
	return &OutgoingAccountLinkApiRequest{}
}

func (r *OutgoingAccountLinkApiRequest) Bind(
	mobilenumber,
	CryptoInfo string,

) error {

	r.AccountLink.MobileNo = "91" + mobilenumber
	r.AccountLink.CryptoInfo = CryptoInfo
	r.AccountLink.ChannelId = "1"

	return nil
}

func (r *OutgoingAccountLinkApiRequest) GetPayeeNameBind(
	MobileNumber,
	CryptoInfo string,
) error {

	r.AccountLink.MobileNo = "91" + MobileNumber
	r.AccountLink.CryptoInfo = CryptoInfo
	r.AccountLink.ChannelId = "1"

	return nil
}

func (r *OutgoingAccountLinkApiRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutgoingAccountLinkApiRequest) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *OutgoingFetchNomineeRequest) Bind(applicantId, account, txnIdentifier string) error {

	if txnIdentifier == "" {
		transactionId, err := security.GenerateRandomUUID(15)
		if err != nil {
			return err
		}
		r.TxnIdentifier = transactionId
	} else {
		r.TxnIdentifier = txnIdentifier
	}

	r.AccountNo = account
	r.ApplicantId = applicantId

	return nil
}

func (r *OutgoingFetchNomineeRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutgoingFetchNomineeRequest) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type OutgoingUpiMoneyCollectDetailsApiRequest struct {
	CollectDetails UpiCollectDetailsApi `json:"CollectDetails"`
}

type UpiCollectDetailsApi struct {
	MobileNo   string `json:"MobileNo"`
	CryptoInfo string `json:"CryptoInfo"`
	ChannelId  string `json:"CHANNELID"`
}

func NewOutgoingUpiMoneyCollectDetailsApiRequest() *OutgoingUpiMoneyCollectDetailsApiRequest {
	return &OutgoingUpiMoneyCollectDetailsApiRequest{}
}

func (r *OutgoingUpiMoneyCollectDetailsApiRequest) Bind(
	MobileNumber,
	CryptoInfo string,
) error {

	r.CollectDetails.MobileNo = "91" + MobileNumber
	r.CollectDetails.CryptoInfo = CryptoInfo
	r.CollectDetails.ChannelId = "1"

	return nil
}

func (r *OutgoingUpiMoneyCollectDetailsApiRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutgoingUpiMoneyCollectDetailsApiRequest) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type OutgoingUpiMoneyCollectCountApiRequest struct {
	CollectCount UpiCollectCountApi `json:"CollectCount"`
}

type UpiCollectCountApi struct {
	MobileNo   string `json:"MobileNo"`
	CryptoInfo string `json:"CryptoInfo"`
	ChannelId  string `json:"CHANNELID"`
}

func NewOutgoingUpiMoneyCollectCountApiRequest() *OutgoingUpiMoneyCollectCountApiRequest {
	return &OutgoingUpiMoneyCollectCountApiRequest{}
}

func (r *OutgoingUpiMoneyCollectCountApiRequest) Bind(
	MobileNumber,
	CryptoInfo string,
) error {

	r.CollectCount.MobileNo = "91" + MobileNumber
	r.CollectCount.CryptoInfo = CryptoInfo
	r.CollectCount.ChannelId = "1"

	return nil
}

func (r *OutgoingUpiMoneyCollectCountApiRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutgoingUpiMoneyCollectCountApiRequest) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type OutgoingUpiMoneyCollectApprovalApiRequest struct {
	Approval UpiMoneyApprovalApi `json:"Approval"`
}

type UpiMoneyApprovalApi struct {
	MobileNo   string `json:"MobileNo"`
	Payeraddr  string `json:"Payeraddr"`
	Type       string `json:"Type"`
	CryptoInfo string `json:"CryptoInfo"`
	OrgTransID string `json:"OrgTransId"`
	CredData   string `json:"CredData"`
	ChannelId  string `json:"CHANNELID"`
}

func NewOutgoingUpiMoneyCollectApprovalApiRequest() *OutgoingUpiMoneyCollectApprovalApiRequest {
	return &OutgoingUpiMoneyCollectApprovalApiRequest{}
}

func (r *OutgoingUpiMoneyCollectApprovalApiRequest) Bind(
	mobilenumber,
	Cryptoinfo,
	Creddata string,
	response *responses.UpiMoneyCollectDetailsResponse,
) error {

	r.Approval.MobileNo = "91" + mobilenumber
	r.Approval.Payeraddr = response.Response.Response[0].PayerAddr
	r.Approval.Type = "0"
	r.Approval.CryptoInfo = Cryptoinfo
	r.Approval.OrgTransID = response.Response.Response[0].TransactionID
	r.Approval.CredData = Creddata
	r.Approval.ChannelId = "1"

	return nil
}

func (r *OutgoingUpiMoneyCollectApprovalApiRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutgoingUpiMoneyCollectApprovalApiRequest) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type GetUpiIDRequest struct {
	AccountLink AccountLinkdata `json:"AccountLink"`
}

type AccountLinkdata struct {
	MobileNo   string `json:"MobileNo"`
	CryptoInfo string `json:"CryptoInfo"`
	CHANELID   string `json:"CHANELID"`
}

func NewGetUpiIDRequest() *GetUpiIDRequest {
	return &GetUpiIDRequest{}
}

type GetDebitcardDetailRequest struct {
	ApplicantId   string `json:"ApplicantId"`
	AccountNo     string `json:"AccountNo"`
	TxnIdentifier string `json:"TxnIdentifier"`
	ProxyNumber   string `json:"ProxyNumber"`
}

func NewGetDebitcardDetailRequest() *GetDebitcardDetailRequest {
	return &GetDebitcardDetailRequest{}
}

func (r *GetUpiIDRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *GetUpiIDRequest) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *GetDebitcardDetailRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *GetDebitcardDetailRequest) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

// func (r *GetUpiIDRequest) GetUPIIdReq_Bind(
// 	mobileNumber, deviceId, osversion, os, name, email string,
// ) error {
// 	clientId, err := security.GenerateRandomUUID(12)
// 	if err != nil {
// 		return err
// 	}
// 	serverID, err := security.GenerateRandomUUID(12)
// 	if err != nil {
// 		return err
// 	}
// 	cryptoInfo := fmt.Sprintf("%s|%s|%s|%s|%s|%s", deviceId, clientId, serverID, osversion, os, "108")
// 	r.AccountLink.CHANELID = "1"
// 	r.AccountLink.CryptoInfo = cryptoInfo
// 	r.AccountLink.MobileNo = "91" + mobileNumber

// 	return nil
// }

func (r *GetDebitcardDetailRequest) GetDebitCardDetailReq_Bind(
	ApplicantId, AccountNo, ProxyNumber, transactionId string,
) error {
	transactionID := security.GenerateRandomCode(18)
	r.ApplicantId = ApplicantId
	r.TxnIdentifier = transactionID
	r.AccountNo = AccountNo
	r.ProxyNumber = ProxyNumber

	return nil
}

type GetAccountDetail struct {
	ApplicantId string `json:"ApplicantId" validate:"required"`
	AccountNo   string `json:"AccountNumber" validate:"required,len=16,numeric"`
}

func NewGetAccountDetailRequest() *GetAccountDetail {
	return &GetAccountDetail{}
}

func (r *GetAccountDetail) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *GetAccountDetail) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type UploadAddressProofReq struct {
	CustomerID string  `json:"CustomerID" validate:"required"`
	Extention  string  `json:"Ext" validate:"required"`
	FileName   string  `json:"FileName" validate:"required"`
	FileString string  `json:"FileString" validate:"required"`
	FolderName string  `json:"FolderName" validate:"required"`
	MineType   *string `json:"MineType" validate:"required"`
}

func NewUploadaddressProofReq() *UploadAddressProofReq {
	return &UploadAddressProofReq{}
}

func (r *UploadAddressProofReq) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *UploadAddressProofReq) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *UploadAddressProofReq) Bind(file *multipart.FileHeader, customerId string) error {
	ext, extValid := security.IsValidFileExtension(file.Filename)
	if !extValid {
		return errors.New("Invalid file type. Allowed types are JPG, JPEG, PDF.")
	}
	tempDir := os.TempDir()
	aadharTempFile, err := os.CreateTemp(tempDir, "aadhar-*"+ext)
	if err != nil {
		return errors.New("Failed to save Aadhar document1")
	}

	defer os.Remove(aadharTempFile.Name())

	src, err := file.Open()
	if err != nil {
		return errors.New("Failed to open Aadhar document")
	}
	defer src.Close()

	fileInfo, err := aadharTempFile.Stat()
	if err != nil {
		return err
	}

	fileSize := fileInfo.Size()
	maxSize := int64(256 * 1024)

	if fileSize > maxSize {
		return errors.New("File size should not be less then 256KB.")
	}

	if _, err := io.Copy(aadharTempFile, src); err != nil {
		return errors.New("Failed to save Aadhar document")
	}

	filePath := aadharTempFile.Name()
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return errors.New("failed to read file")
	}

	// Perform virus scan on the documents
	isClean, description, err := security.ScanDocument(aadharTempFile.Name())
	if err != nil {
		return errors.New("Error scanning document: " + err.Error())
	}
	if !isClean {
		return errors.New("document contains a virus:" + description)
	}

	base64Encoded := base64.StdEncoding.EncodeToString(fileContent)
	// fileNameWithExt := filepath.Base(file.Filename)
	// fileName := strings.TrimSuffix(fileNameWithExt, filepath.Ext(fileNameWithExt))
	// r.ApplicationClass = fileName
	// r.ApplicationID = customerId
	extWithoutDot := strings.TrimPrefix(ext, ".")
	r.Extention = extWithoutDot
	r.FileName = fmt.Sprintf("%s-%s-%s", "PAYDOH", customerId, strings.ReplaceAll(file.Filename, " ", ""))
	r.CustomerID = customerId
	r.FileString = base64Encoded
	r.FolderName = "PAYDOH"
	// r.CategoryName = "PAYDOH"
	// r.ApplicationName = "PAYDOH"
	return nil
}

type UpdateAddressReq struct {
	ServiceName string     `json:"ServiceName"`
	BranchCode  string     `json:"BranchCode"`
	CustomerID  string     `json:"CustomerID"`
	Channel     string     `json:"Channel"`
	Parameter   Parameters `json:"Parameters"`
}
type Parameters struct {
	Address1     string `json:"Address1"`
	Address2     string `json:"Address2"`
	Address3     string `json:"Address3"`
	CityCode     string `json:"City"`
	StateCode    string `json:"State"`
	Country      string `json:"Country"`
	CityName     string `json:"CityName"`
	StateName    string `json:"StateName"`
	CountryName  string `json:"CountryName"`
	ZipCode      string `json:"ZipCode"`
	FileGenID    string `json:"File_Gen_ID"`
	FileGenDate  string `json:"File_Gen_Date"`
	FileGenName  string `json:"File_Gen_Name"`
	FileGenID2   string `json:"File_Gen_ID2"`
	FileGenName2 string `json:"File_Gen_Name2"`
}

func NewUpdateAddressReq() *UpdateAddressReq {
	return &UpdateAddressReq{}
}

func (r *UpdateAddressReq) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *UpdateAddressReq) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *UpdateAddressReq) Bind(req *AddressUpdateRequest, customerId, cityCode, cityName, stateName, fileGenID1, fileGenID2, stateCode string, fileGenID1Name, fileGenID2Name *UploadAddressProofReq) error {
	currentTime := time.Now()
	fileGenDate := currentTime.Format("02-Jan-2006")
	addressServiceType := "COMM_ADDR_UPDATE"
	if req.AddressType == "permanent" {
		addressServiceType = "PERM_ADDR_UPDATE"
	} else {
		addressServiceType = "COMM_ADDR_UPDATE"
	}

	r.ServiceName = addressServiceType
	r.BranchCode = "4908"
	r.CustomerID = customerId
	r.Channel = "P"
	r.Parameter.Address1 = req.Address1
	r.Parameter.Address2 = req.Address2
	r.Parameter.Address3 = req.Address3
	r.Parameter.CityCode = cityCode
	r.Parameter.StateCode = stateCode
	r.Parameter.Country = "IN"
	r.Parameter.CityName = cityName
	r.Parameter.StateName = stateName
	r.Parameter.CountryName = "India"
	r.Parameter.ZipCode = req.PinCode
	r.Parameter.FileGenID = fileGenID1
	r.Parameter.FileGenDate = fileGenDate
	r.Parameter.FileGenName = fileGenID1Name.FileName
	r.Parameter.FileGenID2 = fileGenID2
	r.Parameter.FileGenName2 = fileGenID2Name.FileName
	return nil
}

type SetDebitCardPin struct {
	ApplicantId    string `json:"ApplicantId"`
	AccountNo      string `json:"AccountNo"`
	TxnIdentifier  string `json:"TxnIdentifier"`
	OtpServiceType string `json:"OtpServiceType"`
	ProxyNumber    string `json:"ProxyNumber"`
	EncryptedPAN   string `json:"EncryptedPAN"`
	PinNo          string `json:"PinNo"`
	PinReset       string `json:"PinReset"`
}

func NewSetDebitCardPin() *SetDebitCardPin {
	return &SetDebitCardPin{}
}

func (r *SetDebitCardPin) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *SetDebitCardPin) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *SetDebitCardPin) Bind(
	response *responses.DebitcardDetailResponse, pin, txnIdentifier, pinType string,
) error {

	otpServiceType := ""
	PinReset := ""
	if pinType == "New" {
		// otpServiceType = "PinSet"
		otpServiceType = "Card PIN SET"
		PinReset = "N"
	} else if pinType == "Reset" {
		otpServiceType = "Card PIN RESET"
		PinReset = "Y"
	} else {
		return errors.New("Please Pass Appropriate Pin Type 1")
	}

	r.ApplicantId = response.ApplicantId

	r.AccountNo = response.AccountNo
	r.TxnIdentifier = txnIdentifier
	r.OtpServiceType = otpServiceType
	r.ProxyNumber = response.ServiceData.CardData[0].ProxyNumber
	r.EncryptedPAN = response.ServiceData.CardData[0].EncryptedPAN
	r.PinNo = pin
	r.PinReset = PinReset

	return nil
}

type SetDebitCardOTPReq struct {
	ApplicantId   string `json:"ApplicantId"`
	AccountNo     string `json:"AccountNo"`
	TxnIdentifier string `json:"TxnIdentifier"`
	ServiceType   string `json:"ServiceType"`
	Otp           string `json:"Otp"`
}

func NewSetDebitCardPinOTP() *SetDebitCardOTPReq {
	return &SetDebitCardOTPReq{}
}

func (r *SetDebitCardOTPReq) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *SetDebitCardOTPReq) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *SetDebitCardOTPReq) Bind(
	ApplicantId, account_no, txnid, otpType, otp string,
) error {
	transactionID := ""
	if txnid == "" {
		txnID, err := security.GenerateRandomUUID(20)
		if err != nil {
			return err
		}
		transactionID = txnID
	} else {
		transactionID = txnid
	}

	var serviceType string

	switch otpType {
	case "New":
		serviceType = "Card PIN SET"
	case "Reset":
		serviceType = "Card PIN RESET"
	case "SetDomesticCardLimit":
		serviceType = "Card DOM LIMIT CHANGE"
	case "SetInternationalCardLimit":
		serviceType = "Card INT LIMIT CHANGE"
	case "SetCardBlock":
		serviceType = "Card Block"
	case "SetCardUnblock":
		serviceType = "Card Unblock"
	case "SetCardBlockPermanently":
		serviceType = "Card Block"
	default:
		serviceType = ""
		return errors.New("Please Pass Appropriate Pin Type 2 " + otpType)
	}

	r.ApplicantId = ApplicantId

	r.AccountNo = account_no
	r.TxnIdentifier = strings.ReplaceAll(transactionID, "-", "")
	r.ServiceType = serviceType
	if otp != "" {
		r.Otp = otp
	}
	return nil
}

type OutgoingUpiTransactionHistoryApiRequest struct {
	TransactionHistory UpiTransactionhistoryApi `json:"TransactionHistory"`
}
type UpiTransactionhistoryApi struct {
	FromDate   string `json:"FromDate"`
	ToDate     string `json:"ToDate"`
	MobileNo   string `json:"MobileNo"`
	CryptoInfo string `json:"CryptoInfo"`
	ChannelId  string `json:"CHANNELID"`
}

func NewOutgoingUpiTransactionHistoryApiRequest() *OutgoingUpiTransactionHistoryApiRequest {
	return &OutgoingUpiTransactionHistoryApiRequest{}
}

func (r *OutgoingUpiTransactionHistoryApiRequest) Bind(
	mobilenumber,
	cryptoInfo string,
	request *IncomingUpiTransactionHistoryApiRequest,
) error {

	r.TransactionHistory.FromDate = request.FromDate
	r.TransactionHistory.ToDate = request.ToDate
	r.TransactionHistory.MobileNo = mobilenumber
	r.TransactionHistory.CryptoInfo = cryptoInfo
	r.TransactionHistory.ChannelId = "1"

	return nil
}

func (r *OutgoingUpiTransactionHistoryApiRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutgoingUpiTransactionHistoryApiRequest) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type OutgoingUpiReqSetCreRequest struct {
	ReqSetCre UpiReqSetCreRequest `json:"ReqSetCre"`
}

type UpiReqSetCreRequest struct {
	MobileNo    string `json:"MobileNo"`
	TxnID       string `json:"TxnID"`
	Payeraddr   string `json:"Payeraddr"`
	CryptoInfo  string `json:"CryptoInfo"`
	NewCredData string `json:"NewCred_Data"`
	DeviceIP    string `json:"DeviceIP"`
	GeoLocation string `json:"GeoLocation"`
	CredData    string `json:"Cred_Data"`
	ChannelID   string `json:"CHANNELID"`
}

func NewOutgoingUpiReqSetCreRequest() *OutgoingUpiReqSetCreRequest {
	return &OutgoingUpiReqSetCreRequest{}
}

func (r *OutgoingUpiReqSetCreRequest) Bind(
	mobileNumber,
	payerAddr,
	deviceIp,
	cryptoInfo,
	latLong string,
	request *IncomingUpiChangeUpiPinRequest) error {

	r.ReqSetCre.MobileNo = mobileNumber
	r.ReqSetCre.TxnID = request.TransId
	r.ReqSetCre.Payeraddr = payerAddr
	r.ReqSetCre.CryptoInfo = cryptoInfo
	r.ReqSetCre.NewCredData = request.NewUpiPin
	r.ReqSetCre.DeviceIP = deviceIp
	r.ReqSetCre.GeoLocation = latLong
	r.ReqSetCre.CredData = request.OldUpiPin
	r.ReqSetCre.ChannelID = "1"

	return nil
}

func (r *OutgoingUpiReqSetCreRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OutgoingUpiReqSetCreRequest) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}
