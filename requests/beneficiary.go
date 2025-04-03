package requests

import (
	"encoding/json"
	"errors"

	"bitbucket.org/paydoh/paydoh-commons/customvalidation"
	"github.com/leebenson/conform"

	"bankapi/constants"
)

type AddNewBeneficiary struct { // Mandatory. Account No of the Applicant.
	BenfName      string `json:"beneficiary_name" validate:"required,max=100"`
	BenfNickName  string `json:"beneficiary_nickname" validate:"required,max=100"`                        // Mandatory. Name of beneficiary.
	BenfIFSC      string `json:"beneficiary_ifsc" validate:"required,ifsc_code,len=11"`                   // Mandatory. IFSC code of beneficiary.
	BenfAcctNo    string `json:"beneficiary_account_number" validate:"required,numeric,min=12,max=28"`    // Mandatory. Beneficiary Account No.
	BenfAcctNo1   string `json:"account_number_reconfirmation" validate:"required,numeric,min=12,max=28"` // To be reconfirmed.
	BenfAcctType  string `json:"beneficiary_account_type" validate:"required"`                            // Mandatory. Account type of beneficiary like Savings/Current etc.
	BenfMobNo     string `json:"beneficiary_mobile_number" validate:"required,mobile_number,len=10"`      // Mandatory. Mobile No of beneficiary.
	PaymentMode   string `json:"payment_mode" validate:"required,oneof=NEFT IMPS IFT"`                    // Mandatory. Payment mode like NEFT, IMPS, BOTH(If both IMPS and NEFT required).
	ResendOtp     string `json:"resend_otp" validate:"required"`                                          // Mandatory. Y/N.
	RetryFlag     string `json:"retry_flag" validate:"required"`                                          // Mandatory. Y/N.
	OTP           string `json:"otp" validate:"omitempty"`                                                // Tag need to be passed mandatorily, even if empty. To be passed when applicant submits OTP.
	TxnIdentifier string `json:"txn_identifier" validate:"omitempty"`
}

type AddBeneficiaryOtpRequest struct {
	Otp string `json:"otp" validate:"required"` // Mandatory. OTP received by the applicant to be passed.
}

type PaymentRequest struct {
	PaymentMode   string `json:"payment_mode" validate:"required,oneof=NEFT IMPS IFT"`
	BenfId        string `json:"beneficiary_id"`
	BenfNickName  string `json:"beneficiary_nickname,omitempty"`
	BenfName      string `json:"beneficiary_name" validate:"required"`
	BenfIfsc      string `json:"beneficiary_ifsc" validate:"required,ifsc_code,len=11"`
	BenfAcctNo    string `json:"account_number" validate:"required,numeric,min=12,max=28"`
	BenfAcctType  string `json:"beneficiary_account_type" validate:"required"`
	BenfMobNo     string `json:"beneficiary_mobile_number" validate:"required,mobile_number,len=10"`
	Amount        string `json:"amount" validate:"required,numeric,gt=0"`
	Remarks       string `json:"remarks" validate:"required,max=100"`
	ResendOtp     string `json:"resend_otp" validate:"required"`
	RetryFlag     string `json:"retry_flag" validate:"required"`
	Otp           string `json:"otp"`
	QuickTransfer string `json:"quick_transfer" validate:"required"`
}

type QuickTransferBeneficiaryRegistrationRequest struct {
	BenficiaryId string `json:"beneficiary_id" validate:"required,len=20"`
}

func NewAddNewBeneficiaryRequest() *AddNewBeneficiary {
	return &AddNewBeneficiary{}
}

func NewAddBeneficiaryOtpRequest() *AddBeneficiaryOtpRequest {
	return &AddBeneficiaryOtpRequest{}
}

func NewPaymentRequest() *PaymentRequest {
	return &PaymentRequest{}
}

func NewQuickTransferBeneficiaryRegistrationRequest() *QuickTransferBeneficiaryRegistrationRequest {
	return &QuickTransferBeneficiaryRegistrationRequest{}
}

func (r *AddNewBeneficiary) Validate(payload string) error {
	if err := r.Unmarshal([]byte(payload)); err != nil {
		return err
	}

	if err := conform.Strings(r); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}

	if r.ResendOtp == "" {
		return errors.New("resend OTP is required")
	}

	if r.RetryFlag == "" {
		return errors.New("retry flag is required")
	}

	if r.ResendOtp == "Y" && r.OTP == "" {
		return constants.ErrOtpIsRequired
	}

	return nil
}

func (r *AddBeneficiaryOtpRequest) Validate(payload string) error {
	if err := r.Unmarshal([]byte(payload)); err != nil {
		return err
	}

	if err := conform.Strings(r); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}

	return nil
}

func (r *PaymentRequest) Validate(payload string) error {
	if err := r.Unmarshal([]byte(payload)); err != nil {
		return err
	}

	if r.ResendOtp == "" {
		return errors.New("resend OTP is required")
	}

	if r.RetryFlag == "" {
		return errors.New("retry flag is required")
	}

	// if r.QuickTransfer != "Y" && r.BenfId == "" {
	// 	return constants.ErrBeneficiaryIdIsRequired
	// }

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}

	// if r.ResendOtp == "Y" && r.Otp == "" {
	// 	return constants.ErrOtpIsRequired
	// }

	return nil
}

func (r *QuickTransferBeneficiaryRegistrationRequest) Validate(payload string) error {
	if err := r.Unmarshal([]byte(payload)); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}

	return nil
}

func (r *AddNewBeneficiary) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *AddNewBeneficiary) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *AddBeneficiaryOtpRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *AddBeneficiaryOtpRequest) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *PaymentRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *PaymentRequest) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *QuickTransferBeneficiaryRegistrationRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *QuickTransferBeneficiaryRegistrationRequest) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type CheckPaymentStatus struct {
	TransactionId string `json:"txnId" validate:"required"`
}

func NewCheckPaymentStatusRequest() *CheckPaymentStatus {
	return &CheckPaymentStatus{}
}

func (r *CheckPaymentStatus) Validate(payload string) error {
	if err := r.Unmarshal([]byte(payload)); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}

	return nil
}

func (r *CheckPaymentStatus) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *CheckPaymentStatus) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}
