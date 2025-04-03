package requests

import (
	"encoding/json"
	"errors"

	"bitbucket.org/paydoh/paydoh-commons/customvalidation"
)

type AddNomineeRequest struct {
	NomReqType              string `json:"nominee_request_type" validate:"required,max=20"`
	NomName                 string `json:"nominee_name" validate:"required,max=100"`
	NomRelation             string `json:"nominee_relation" validate:"required,max=20"`
	NomDOB                  string `json:"nominee_dob" validate:"required"`
	NomAddressL1            string `json:"nominee_address_line_1" validate:"required,max=105"`
	NomAddressL2            string `json:"nominee_address_line_2" validate:"required,max=105"`
	NomAddressL3            string `json:"nominee_address_line_3" validate:"required,max=105"`
	NomCity                 string `json:"nominee_city" validate:"required,max=50"`
	NomState                string `json:"nominee_state" validate:"required,max=20"`
	NomCountry              string `json:"nominee_country" validate:"required,max=20"`
	NomZipcode              string `json:"nominee_zip_code" validate:"required,len=6"`
	GuardianName            string `json:"guardian_name,omitempty"`
	GuardianNomineeRelation string `json:"guardian_nominee_relation,omitempty"`
	GuardianAddressL1       string `json:"guardian_address_line_1,omitempty"`
	GuardianAddressL2       string `json:"guardian_address_line_2,omitempty"`
	GuardianAddressL3       string `json:"guardian_address_line_3,omitempty"`
	GuardianCity            string `json:"guardian_city,omitempty"`
	GuardianState           string `json:"guardian_state,omitempty"`
	GuardianCountry         string `json:"guardian_country,omitempty"`
	GuardianZipcode         string `json:"guardian_zip_code,omitempty"`
	ResendOtp               string `json:"resend_otp" validate:"required,oneof=Y N"`
	RetryFlag               string `json:"retry_flag" validate:"required,oneof=Y N"`
	Otp                     string `json:"otp"`
	NomineeMobileNumber     string `json:"nominee_mobile_no"`
}

type VerifyOtpRequest struct {
	Otp     string `json:"otp,omitempty"`
	ReqType string `json:"nominee_request_type" validate:"required"`
}

type VerifyBeneficiaryOTPRequest struct {
	Otp string `json:"otp" validate:"required"`
}

type VerifyPaymentOTPRequest struct {
	Otp string `json:"otp" validate:"required"`
}

func NewAddNomineeRequest() *AddNomineeRequest {
	return &AddNomineeRequest{}
}

func NewVerifyOtpRequest() *VerifyOtpRequest {
	return &VerifyOtpRequest{}
}

func NewVerifyBeneficiaryOTPRequest() *VerifyBeneficiaryOTPRequest {
	return &VerifyBeneficiaryOTPRequest{}
}

func NewVerifyPaymentOTPRequest() *VerifyPaymentOTPRequest {
	return &VerifyPaymentOTPRequest{}
}

func (r *AddNomineeRequest) Validate(payload string) error {
	if err := json.Unmarshal([]byte(payload), r); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}

	if r.GuardianName != "" && len(r.GuardianName) > 100 {
		return errors.New("guardian name must be less than 100 characters")
	}

	if r.GuardianNomineeRelation != "" && len(r.GuardianNomineeRelation) > 20 {
		return errors.New("guardian relation must be less than 20 characters")
	}

	if r.GuardianAddressL1 != "" && len(r.GuardianAddressL1) > 105 {
		return errors.New("guardian address line 1 must be less than 105 characters")
	}

	if r.GuardianAddressL2 != "" && len(r.GuardianAddressL2) > 105 {
		return errors.New("guardian address line 2 must be less than 105 characters")
	}

	if r.GuardianAddressL3 != "" && len(r.GuardianAddressL3) > 105 {
		return errors.New("guardian address line 3 must be less than 105 characters")
	}

	if r.GuardianCity != "" && len(r.GuardianCity) > 50 {
		return errors.New("guardian city must be less than 50 characters")
	}

	if r.GuardianState != "" && len(r.GuardianState) > 20 {
		return errors.New("guardian state must be less than 20 characters")
	}

	if r.GuardianCountry != "" && len(r.GuardianCountry) > 20 {
		return errors.New("guardian country must be less than 20 characters")
	}

	if r.GuardianZipcode != "" && len(r.GuardianZipcode) > 6 && len(r.GuardianZipcode) < 6 {
		return errors.New("guardian zipcode must be equal to 6 characters")
	}

	// if r.ResendOtp == "Y" && r.Otp == "" {
	// 	return constants.ErrOtpIsRequired
	// }

	if r.RetryFlag == "" {
		return errors.New("retry flag is required")
	}

	return nil
}

func (r *VerifyOtpRequest) Validate(payload string) error {
	if err := json.Unmarshal([]byte(payload), r); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}
	return nil
}

func (r *VerifyBeneficiaryOTPRequest) Validate(payload string) error {
	if err := json.Unmarshal([]byte(payload), r); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}
	return nil
}

func (r *VerifyPaymentOTPRequest) Validate(payload string) error {
	if err := json.Unmarshal([]byte(payload), r); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}
	return nil
}
