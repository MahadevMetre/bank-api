package requests

import (
	"encoding/json"
	"encoding/xml"
	"errors"

	"bitbucket.org/paydoh/paydoh-commons/customvalidation"
)

type CreateUPIRequest struct {
	DeviceCapability string `json:"device_capability" validate:"required"`
	Challenge        string `json:"challenge" validate:"omitempty"`
	Location         string `json:"location" validate:"required"`
}

type SetUpiPinRequest struct {
	UpiPin       string `json:"upi_pin" validate:"required"`
	Otp          string `json:"otp" validate:"required"`
	AtmPin       string `json:"atm_pin,omitempty"`
	TransId      string `json:"trans_id" validate:"required"`
	Cred_AADHAAR string `json:"Cred_AADHAAR"`
}

type ValidateVpaRequest struct {
	RPayeraddr string `json:"rpayer_addr" validate:"required"`
	RPayername string `json:"rpayer_name,omitempty"`
	PayerCode  string `json:"payer_code,omitempty"`
}

type PayMoneyWithVpaRequest struct {
	Payeeaddr       string `json:"payee_addr" validate:"required"`
	PayeeName       string `json:"payee_name" validate:"required"`
	PayerAmount     string `json:"payer_amount" validate:"required"`
	UpiPin          string `json:"upi_pin" validate:"required"`
	Remark          string `json:"remark" validate:"required"`
	TransId         string `json:"trans_id" validate:"required"`
	MccCode         string `json:"mcc_code" validate:"required"`
	TransactionType string `json:"transaction_type" validate:"required"`
}

type ReqBalEnqRequest struct {
	UpiPin  string `json:"upi_pin" validate:"required,numeric"`
	TransId string `json:"trans_id" validate:"required"`
}

type GetAllBankAccount struct {
	MobileNumber string `json:"mobile_number" validate:"required,numeric"`
}

type AadharReqlistaccount struct {
	AadharNumber string `json:"aadhar_number" validate:"required,numeric"`
}

type UpiTokenRequest struct {
	Challenge     string `json:"challenge" validate:"required"`
	ChallengeType string `json:"challenge_type" validate:"required"`
}

type UpiTokenXMLRequest struct {
	Challenge string `json:"challenge" validate:"omitempty"`
}

type SimBindingRequest struct {
	Type string `json:"type,omitempty"`
}

func NewSimBindingRequest() *SimBindingRequest {
	return &SimBindingRequest{}
}

func (r *SimBindingRequest) Validate(payload string) error {
	if err := json.Unmarshal([]byte(payload), r); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}
	return nil
}

func NewUpiTokenXMLRequest() *UpiTokenXMLRequest {
	return &UpiTokenXMLRequest{}
}

func (r *UpiTokenXMLRequest) Validate(payload string) error {
	if err := json.Unmarshal([]byte(payload), r); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}
	return nil
}

type UpiCollectMoneyRequest struct {
	UpiPin string `json:"upi_pin" validate:"required,numeric"`
}

type IncomingRemappingRequest struct {
	Challenge   string `json:"challenge" validate:"required"`
	DeviceToken string `json:"device_token,omitempty"`
}

func NewCreateUPIRequest() *CreateUPIRequest {
	return &CreateUPIRequest{}
}

func NewSetUpiPinRequest() *SetUpiPinRequest {
	return &SetUpiPinRequest{}
}

func NewValidateVpaRequest() *ValidateVpaRequest {
	return &ValidateVpaRequest{}
}

func NewPayMoneyWithVpaRequest() *PayMoneyWithVpaRequest {
	return &PayMoneyWithVpaRequest{}
}

func NewReqBalEnqRequest() *ReqBalEnqRequest {
	return &ReqBalEnqRequest{}
}

func NewGetAllBankAccount() *GetAllBankAccount {
	return &GetAllBankAccount{}
}

func NewAadharReqlistaccount() *AadharReqlistaccount {
	return &AadharReqlistaccount{}
}

func NewUpiTokenRequest() *UpiTokenRequest {
	return &UpiTokenRequest{}
}

func NewUpiCollectMoneyRequest() *UpiCollectMoneyRequest {
	return &UpiCollectMoneyRequest{}
}

func NewIncomingRemappingRequest() *IncomingRemappingRequest {
	return &IncomingRemappingRequest{}
}

func (r *CreateUPIRequest) Validate(payloads string) error {
	if err := json.Unmarshal([]byte(payloads), r); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}
	return nil
}

func (r *SetUpiPinRequest) Validate(payload string) error {
	if err := json.Unmarshal([]byte(payload), r); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}

	return nil
}

func (r *ValidateVpaRequest) Validate(payload string) error {
	if err := json.Unmarshal([]byte(payload), r); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}
	return nil
}

func (r *PayMoneyWithVpaRequest) Validate(payload string) error {
	if err := json.Unmarshal([]byte(payload), r); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}
	// if len(r.UpiPin) != 6 {
	// 	return errors.New("upi pin must be 6 digits")
	// }

	return nil
}

func (r *ReqBalEnqRequest) Validate(payload string) error {
	if err := json.Unmarshal([]byte(payload), r); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}
	if len(r.UpiPin) != 6 {
		return errors.New("upi pin must be 6 digits")
	}

	return nil
}

func (r *GetAllBankAccount) Validate(payload string) error {
	if err := json.Unmarshal([]byte(payload), r); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}
	if len(r.MobileNumber) != 10 {
		return errors.New("mobile number must be 10 digits")
	}

	return nil
}

func (r *AadharReqlistaccount) Validate(payload string) error {
	if err := json.Unmarshal([]byte(payload), r); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}
	return nil
}

func (r *UpiTokenRequest) Validate(payload string) error {
	if err := json.Unmarshal([]byte(payload), r); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}
	return nil
}

func (r *UpiCollectMoneyRequest) Validate(payload string) error {
	if err := json.Unmarshal([]byte(payload), r); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}
	if len(r.UpiPin) != 6 {
		return errors.New("upi pin must be 6 digits")
	}

	return nil
}

func (r *IncomingRemappingRequest) Validate(payload string) error {
	if err := json.Unmarshal([]byte(payload), r); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}
	return nil
}

type IncomingUpiTransactionHistoryApiRequest struct {
	FromDate string `json:"from_date" validate:"required"`
	ToDate   string `json:"to_date" validate:"required"`
}

func NewIncomingUpiTransactionHistoryApiRequest() *IncomingUpiTransactionHistoryApiRequest {
	return &IncomingUpiTransactionHistoryApiRequest{}
}

func (r *IncomingUpiTransactionHistoryApiRequest) Validate(payload string) error {
	if err := json.Unmarshal([]byte(payload), r); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}
	return nil
}

type IncomingUpiChangeUpiPinRequest struct {
	OldUpiPin string `json:"old_upi_pin" validate:"required"`
	NewUpiPin string `json:"new_upi_pin" validate:"required"`
	TransId   string `json:"trans_id" validate:"required"`
}

func NewIncomingUpiChangeUpiPinRequest() *IncomingUpiChangeUpiPinRequest {
	return &IncomingUpiChangeUpiPinRequest{}
}

func (r *IncomingUpiChangeUpiPinRequest) Validate(payload string) error {
	if err := json.Unmarshal([]byte(payload), r); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}
	return nil
}

type RespListKeys struct {
	XMLName  xml.Name `xml:"RespListKeys"`
	Xmlns    string   `xml:"xmlns,attr"`
	XmlnsNs3 string   `xml:"xmlns:ns3,attr"`
	KeyList  KeyList  `xml:"keyList"`
}

type KeyList struct {
	Key Key `xml:"key"`
}

type Key struct {
	Code     string `xml:"code,attr"`
	Ki       string `xml:"ki,attr"`
	Owner    string `xml:"owner,attr"`
	Type     string `xml:"type,attr"`
	KeyValue string `xml:"keyValue"`
}
