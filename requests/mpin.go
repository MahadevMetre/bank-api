package requests

import (
	"encoding/json"
	"errors"

	"bitbucket.org/paydoh/paydoh-commons/customvalidation"
)

type MpinRequest struct {
	Mpin string `json:"mpin" validate:"required"`
}

func NewMpinRequest() *MpinRequest {
	return &MpinRequest{}
}

func (r *MpinRequest) Validate(payload string) error {
	if err := json.Unmarshal([]byte(payload), r); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}
	if len(r.Mpin) != 4 {
		return errors.New("Mpin must be 4 digits")
	}

	return nil
}

type ResetMpinRequest struct {
	ExistingMpin string `json:"existing_mpin" validate:"required"`
	NewMpin      string `json:"new_mpin" validate:"required"`
}

func NewResetMpinRequest() *ResetMpinRequest {
	return &ResetMpinRequest{}
}

func (r *ResetMpinRequest) Validate(payload string) error {
	if err := json.Unmarshal([]byte(payload), r); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}
	if len(r.ExistingMpin) != 4 {
		return errors.New("ExistingMpin must be 4 digits")
	}
	if len(r.NewMpin) != 4 {
		return errors.New("NewMpin must be 4 digits")
	}

	return nil
}

type ForgotMpinResetRequest struct {
	AccountNumber    string `json:"account_number" validate:"required"`
	Email            string `json:"email" validate:"required,email"`
	MotherMaidenName string `json:"mother_maiden_name" validate:"required"`
	PinCode          string `json:"pin_code" validate:"required"`
}

func NewForgotMpinResetRequest() *ForgotMpinResetRequest {
	return &ForgotMpinResetRequest{}
}

func (r *ForgotMpinResetRequest) Validate(payload string) error {
	if err := json.Unmarshal([]byte(payload), r); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}
	if len(r.PinCode) != 6 {
		return errors.New("PinCode must be 6 digits")
	}

	return nil
}

type UpdateMpinRequest struct {
	Mpin string `json:"mpin" validate:"required"`
}

func NewUpdateMpinRequest() *UpdateMpinRequest {
	return &UpdateMpinRequest{}
}

func (r *UpdateMpinRequest) Validate(payload string) error {
	if err := json.Unmarshal([]byte(payload), r); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}
	if len(r.Mpin) != 4 {
		return errors.New("Mpin must be 4 digits")
	}

	return nil
}
