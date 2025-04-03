package requests

import (
	"bankapi/constants"
	"bankapi/security"
	"encoding/json"
	"errors"
	"strings"

	"bitbucket.org/paydoh/paydoh-commons/customvalidation"
)

type DebitCardGenerationRequest struct {
	DebitCardGenerationType string `json:"debitcard_generation_type"  validate:"required"` //virtual -physical-both
}

func NewDebitCardGenerationType() *DebitCardGenerationRequest {
	return &DebitCardGenerationRequest{}
}

func (r *DebitCardGenerationRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *DebitCardGenerationRequest) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *DebitCardGenerationRequest) Validate(payload string) error {
	if err := r.Unmarshal([]byte(payload)); err != nil {
		return err
	}
	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}
	return nil
}

type GenerateVirtualDebitcardOutGoingReq struct {
	ApplicantId   string       `json:"ApplicantId"`
	TxnIdentifier string       `json:"TxnIdentifier"`
	ServiceData   ServiceData1 `json:"serviceData"`
}

type ServiceData1 struct {
	AccountNumber string        `json:"accountNumber"`
	IssueCardForm IssueCardForm `json:"issueCardForm"`
}

type IssueCardForm struct {
	NameOnCard string `json:"nameOnCard"`
}

func NewGenerateVirtualDebitCardOutGoingReq() *GenerateVirtualDebitcardOutGoingReq {
	return &GenerateVirtualDebitcardOutGoingReq{}
}

func (r *GenerateVirtualDebitcardOutGoingReq) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *GenerateVirtualDebitcardOutGoingReq) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *GenerateVirtualDebitcardOutGoingReq) Bind(name, applicantId, accountNumber, txnIdetifier string) error {
	transactionId := ""
	if txnIdetifier == "" {
		transaction_id, err := security.GenerateRandomUUID(20)
		if err != nil {
			return err
		}
		transactionId = strings.ReplaceAll(transaction_id, "-", "")
	} else {
		transactionId = txnIdetifier
	}

	r.ApplicantId = applicantId
	r.TxnIdentifier = transactionId // strings.ReplaceAll(transactionId, "-", "")
	r.ServiceData.AccountNumber = accountNumber
	r.ServiceData.IssueCardForm.NameOnCard = strings.ToUpper(name)

	return nil
}

type GeneratePhysicalDebitCardOutGoingReq struct {
	ApplicantId   string `json:"ApplicantId"`
	TxnIdentifier string `json:"TxnIdentifier"`
	ProxyNumber   string `json:"ProxyNumber"`
	AccountNo     string `json:"AccountNo"`
}

func NewGeneratePhysicalDebitCardOutGoingReq() *GeneratePhysicalDebitCardOutGoingReq {
	return &GeneratePhysicalDebitCardOutGoingReq{}
}

func (r *GeneratePhysicalDebitCardOutGoingReq) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *GeneratePhysicalDebitCardOutGoingReq) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *GeneratePhysicalDebitCardOutGoingReq) Bind(applicantId, proxyNumber, accountNo, transactionID string) error {
	transactionId := ""
	if transactionID == "" {
		transactionsId, err := security.GenerateRandomUUID(25)
		if err != nil {
			return err
		}
		transactionId = strings.ReplaceAll(transactionsId, "-", "")
	} else {
		transactionId = transactionID
	}
	r.AccountNo = accountNo
	r.ApplicantId = applicantId
	r.ProxyNumber = proxyNumber
	r.TxnIdentifier = transactionId

	return nil

}

type GetTransactionLimitReq struct {
	TransactionType string `json:"transaction_type" validate:"required"` //Intenational  and Domestic
}

func NewGetTransactionLimit() *GetTransactionLimitReq {
	return &GetTransactionLimitReq{}
}

func (debitReq *GetTransactionLimitReq) Validate(payload string) error {
	if err := json.Unmarshal([]byte(payload), debitReq); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(debitReq); err != nil {
		return err
	}

	return nil
}

type SetCardStatus struct {
	DomesticStatus       string `json:"domestic_block_status" validate:"required"`      // 0 – Unblock 1 - Block
	InternationalStatus  string `json:"international_block_status" validate:"required"` // 0 – Unblock 1 - Block
	IsPermanentlyBlocked string `json:"is_permanently_blocked" validate:"required"`     // 0 – Unblock 1 - Block
}

func NewSetCardStatus() *SetCardStatus {
	return &SetCardStatus{}
}

func (debitReq *SetCardStatus) Validate(payload string) error {
	if err := json.Unmarshal([]byte(payload), debitReq); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(debitReq); err != nil {
		return err
	}
	return nil
}

type EncryptedDebitCardReq struct {
	EncryptReq string `json:"encryptReq"`
}

func (req *EncryptedDebitCardReq) Bind(enc string) (EncryptedDebitCardReq, error) {
	req.EncryptReq = enc
	return *req, nil
}

func (r *EncryptedDebitCardReq) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

//Sachin

// API FOR FETCH KEY
type KeyFetchRequest struct {
	ID string `json:"id"`
}

type LoginRequest struct {
	StaticPassword string `json:"statpwd,omitempty" validate:"required"`
	Mobile         string `json:"mno,omitempty" validate:"required"`
	Email          string `json:"eid,omitempty" validate:"required"`
	AccountNumber  string `json:"accno,omitempty" validate:"required"`
	Name           string `json:"nm,omitempty" validate:"required"`
	CustomerID     string `json:"cid,omitempty" validate:"required"`
	DateOfBirth    string `json:"dob,omitempty" validate:"required"`
	PublicKey      string `json:"key,omitempty" validate:"required"`
}

//API FOR Add Card

type AddCardRequest struct {
	ExpiryMonth string `json:"eMonth" validate:"required"`
	MaxCard     string `json:"max_card" validate:"required"`
	CT          string `json:"ct" validate:"required"`
	BT          string `json:"bt" validate:"required"`
	Enrollid    string `json:"enrollid" validate:"required"`
	Name        string `json:"holdername" validate:"required"`
	ExpiryYear  string `json:"eYear" validate:"required"`
	DebitCardNo string `json:"cardNo" validate:"required"`
	CustomerId  string `json:"cid" validate:"required"`
	Publickey   string `json:"key" validate:"required"`
}

//API FOR List Card

type ListCardRequest struct {
	CustomerId string `json:"cid"  validate:"required"`
	PublicKey  string `json:"key"  validate:"required"`
}

// List Card Control API
type ListCardControlRequest struct {
	CustomerId   string `json:"cid"  validate:"required"`
	EnrollmentId string `json:"enid"  validate:"required"`
	PublicKey    string `json:"key"  validate:"required"`
}

// Fetch Transaction API

type FetchTransactionRequest struct {
	CustomerId       string `json:"cid"`
	EnrollmentId     string `json:"enid"`
	DeliveryChannels string `json:"deliveryIndex"`
	PublicKey        string `json:"key"`
}

type TranDetail struct {
	DeliveryChannel string        `json:"deliveryChannel"`
	DeliveryIndex   string        `json:"deliveryIndex"`
	TranTypes       string        `json:"tranTypes"`
	StatusDc        string        `json:"statusDc"`
	DchBlkStatus    string        `json:"dchBlkStatus"`
	TranM           []Transaction `json:"tranM"`
	Minordigit      string        `json:"minordigit"`
	Symbol          string        `json:"symbol"`
}

type Transaction struct {
	Max       string `json:"max"`
	Status    string `json:"status"`
	Value     string `json:"value"`
	TranValue string `json:"tranValue"`
	TranLabel string `json:"tranLabel"`
}

//Edit Transaction API

type EditTransactionRequest struct {
	CID      string       `json:"cid"  validate:"required"`
	ENID     string       `json:"enid"  validate:"required"`
	CNID     string       `json:"cnid"  validate:"required"`
	Stat     string       `json:"stat"  validate:"required"`
	NCFlag   string       `json:"ncFlag"  validate:"required"`
	TranList []TranDetail `json:"tranList"  validate:"required"`
	Key      string       `json:"key"  validate:"required"`
}

func NewEditTransactionRequest() *EditTransactionRequest {
	return &EditTransactionRequest{}
}

func (req *EditTransactionRequest) Bind(request RequestEditTransaction, cnid, cid, enrollmentId, publicKey string) error {
	transadetail := []TranDetail{}
	for _, d := range request.ReqData {
		c := TranDetail{
			DeliveryChannel: d.Name, //strings.Split(d.Name, " ")[0],
			DeliveryIndex:   constants.TransactionTypes[ /*strings.Split(d.Name, " ")[0]*/ d.Name],
			TranTypes:       d.TranTypes,
			StatusDc:        d.StatusDc,
			DchBlkStatus:    d.DChBlkStatus,
			Minordigit:      "2",
			Symbol:          "INR",
		}

		label := ""
		switch d.TranTypes {
		case "00":
			label = "Purchase"
		case "01":
			label = "Withdrawal"
		}
		t := Transaction{
			Max:       d.MaxLimit,
			Status:    d.TransMStatus,
			Value:     d.SetValue,
			TranValue: d.TranTypes,
			TranLabel: label,
		}
		c.TranM = append(c.TranM, t)
		transadetail = append(transadetail, c)
	}

	req.Stat = "1"
	req.NCFlag = "0"
	req.CNID = cnid
	req.CID = cid
	req.ENID = enrollmentId
	req.Key = publicKey
	req.TranList = transadetail
	return nil
}

// Card Block API
type CardBlockRequest struct {
	CustomerId               string `json:"cid"  validate:"required"`
	EnrollmentId             string `json:"enid"  validate:"required"`
	CardControlId            string `json:"cnid"  validate:"required"`
	CardStatus               string `json:"stat"  validate:"required"`
	IsNewCard                string `json:"ncFlag"  validate:"required"`
	InternationalBlockStatus string `json:"ibs"  validate:"required"`
	BlockStatus              string `json:"bs"  validate:"required"`
	PublicKey                string `json:"key"  validate:"required"`
}

func (r *CardBlockRequest) Bind(cid, enid, cardcontrolId, international_status, domestic_status, publicKey string) error {
	r.CustomerId = cid
	r.EnrollmentId = enid
	r.CardControlId = cardcontrolId
	r.CardStatus = "1"
	r.IsNewCard = "1"
	r.InternationalBlockStatus = international_status
	r.BlockStatus = domestic_status
	r.PublicKey = publicKey
	return nil
}

func (k *KeyFetchRequest) Validate(payload string) error {
	if err := json.Unmarshal([]byte(payload), k); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(k); err != nil {
		return err
	}

	if k.ID == "" {
		return errors.New("id required")
	}

	return nil
}

func (k *LoginRequest) Validate(payload string) error {
	if err := json.Unmarshal([]byte(payload), k); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(k); err != nil {
		return err
	}

	return nil
}
func (k *AddCardRequest) Validate(payload string) error {
	if err := json.Unmarshal([]byte(payload), k); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(k); err != nil {
		return err
	}

	return nil
}

func (k *ListCardRequest) Validate(payload string) error {
	if err := json.Unmarshal([]byte(payload), k); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(k); err != nil {
		return err
	}

	return nil
}
func (k *ListCardControlRequest) Validate(payload string) error {
	if err := json.Unmarshal([]byte(payload), k); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(k); err != nil {
		return err
	}

	return nil
}
func (k *FetchTransactionRequest) Validate(payload string) error {
	if err := json.Unmarshal([]byte(payload), k); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(k); err != nil {
		return err
	}

	return nil
}
func (k *EditTransactionRequest) Validate(payload string) error {
	if err := json.Unmarshal([]byte(payload), k); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(k); err != nil {
		return err
	}

	return nil
}
func (k *CardBlockRequest) Validate(payload string) error {
	if err := json.Unmarshal([]byte(payload), k); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(k); err != nil {
		return err
	}

	return nil
}
func (r *KeyFetchRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *LoginRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
func (r *AddCardRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *ListCardRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
func (r *ListCardControlRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *FetchTransactionRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
func (r *EditTransactionRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *CardBlockRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *KeyFetchRequest) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *LoginRequest) Bind(MobileNumber, Email, AccountNumber, Name, CustomerId, dob string) error {
	r.StaticPassword = security.GenerateRandomCode(10)
	r.Mobile = "+91" + MobileNumber
	r.Email = Email
	r.AccountNumber = AccountNumber
	r.Name = Name
	r.CustomerID = CustomerId
	r.DateOfBirth = dob

	return nil
}

func (r *AddCardRequest) Bind(Name, ExpireYear string, month string, cardNo string, cid string, key string) error {
	r.Name = Name

	r.ExpiryYear = ExpireYear
	r.ExpiryMonth = month
	r.DebitCardNo = cardNo
	r.CustomerId = cid
	r.Publickey = key
	return nil
}
func (r *ListCardRequest) Bind(CustomerId string, publicKey string) error {
	r.CustomerId = CustomerId
	r.PublicKey = publicKey
	return nil
}
func (r *ListCardControlRequest) Bind(CustomerId, enid, publicKey string) error {
	r.CustomerId = CustomerId
	r.EnrollmentId = enid
	r.PublicKey = publicKey
	return nil
}

func (r *FetchTransactionRequest) Bind(CustomerId, enid, publicKey, deliveryChannel string) error {
	r.CustomerId = CustomerId
	r.EnrollmentId = enid
	r.PublicKey = publicKey
	r.DeliveryChannels = deliveryChannel
	return nil
}

type RequestEditTransaction struct {
	ReqData []RequestData `json:"data"`
}

type RequestData struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	MaxLimit     string `json:"maxLimit"`
	SetValue     string `json:"setlimit"`
	TranTypes    string `json:"tranTypes"`
	StatusDc     string `json:"statusDc"`     //Delivery Channel Status
	DChBlkStatus string `json:"dchBlkStatus"` //Delivery Channel Block status 0 Unblock 1 block
	TransMStatus string `json:"transtatus"`
}

func NewRequestEditTransaction() RequestEditTransaction {
	return RequestEditTransaction{}
}

func (s *RequestEditTransaction) Validation(payload string) error {

	if err := json.Unmarshal([]byte(payload), s); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(s); err != nil {
		return err
	}
	return nil
}

type SetDebitCardPinReq struct {
	Pin        string `json:"pin" validate:"required,numeric,min=4,max=4"`
	PinSetType string `json:"pin_set_type" validate:"required"`
}

func NewSetDebitCardPinReq() *SetDebitCardPinReq {
	return &SetDebitCardPinReq{}
}

func (r *SetDebitCardPinReq) Validate(payload string) error {
	if err := json.Unmarshal([]byte(payload), r); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}

	if r.Pin != "" && len(r.Pin) != 4 {
		return errors.New("DebitCard Pin Must be 4 Digit")
	}

	if r.PinSetType == "" {
		return errors.New("Pin Setup Type is required")
	} else if r.PinSetType != "New" && r.PinSetType != "Reset" {
		return errors.New("Pin Setup Type is not valid")
	}

	return nil
}

type DebitCardVerifyOtpReq struct {
	Otp     string `json:"otp" validate:"required"`
	OtpType string `json:"otp_type" validate:"required"`
}

func NewDebitCardVerifyOtpReq() *DebitCardVerifyOtpReq {
	return &DebitCardVerifyOtpReq{}
}

func (r *DebitCardVerifyOtpReq) Validate(payload string) error {
	if err := json.Unmarshal([]byte(payload), r); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}

	if r.Otp != "" && len(r.Otp) != 8 {
		return errors.New("OTP must be 8 Digit")
	}

	validOtpTypes := map[string]bool{
		"New":                       true,
		"Reset":                     true,
		"SetDomesticCardLimit":      true,
		"SetInternationalCardLimit": true,
		"SetCardBlock":              true,
		"SetCardUnblock":            true,
		"SetCardBlockPermanently":   true,
	}

	if r.OtpType == "" {
		return errors.New("Pin Setup Type is required")
	} else if !validOtpTypes[r.OtpType] {
		return errors.New("Pin Setup Type is not valid")
	}

	return nil
}
