package requests

import (
	"bankapi/security"
	"encoding/json"
	"errors"

	"bitbucket.org/paydoh/paydoh-commons/customvalidation"
	"github.com/gin-gonic/gin"
)

type PincodeDetails struct {
	Pincode string `json:"pincode" validate:"required,numeric,min=3,max=6"`
}

type EncryptedRequest struct {
	Data     string `json:"data"`
	CrcValue string `json:"crc_value"`
}

type PersonalInformationRequest struct {
	FirstName   string `json:"first_name" validate:"required"`
	MiddleName  string `json:"middle_name,omitempty"`
	LastName    string `json:"last_name" validate:"required"`
	Gender      string `json:"gender" validate:"required,oneof=male female other"`
	Email       string `json:"email" validate:"required,email"`
	DateOfBirth string `json:"date_of_birth" validate:"required"`
}

type CommunicationAddress struct {
	HouseNo    string `json:"house_no"`
	StreetName string `json:"street_name"`
	Locality   string `json:"locality"`
	Landmark   string `json:"landmark"`
	PinCode    string `json:"pin_code"`
	State      string `json:"state"`
	City       string `json:"city"`
}

type CreateBankAccountRequest struct {
	AnnualTurnOver       string                `json:"annual_turn_over" validate:"required"`
	MaritalStatus        string                `json:"marital_status" validate:"required"`
	CountryResidence     string                `json:"country_residence" validate:"required"`
	MotherMaidenName     string                `json:"mother_maiden_name" validate:"required"`
	CustomerEducation    string                `json:"customer_education" validate:"required"`
	Nationality          string                `json:"nationality" validate:"required"`
	ProfessionCode       string                `json:"ProfessionCode" validate:"required"`
	CommunicationAddress *CommunicationAddress `json:"communication_address,omitempty"`
	IsAddrSameAsAdhaar   bool                  `json:"is_addr_same_as_adhaar,omitempty"`
}

func NewEncryptedRequest() *EncryptedRequest {
	return &EncryptedRequest{}
}

func NewPersonalInformationRequest() *PersonalInformationRequest {
	return &PersonalInformationRequest{}
}

func NewCreateBankAccountRequest() *CreateBankAccountRequest {
	return &CreateBankAccountRequest{}
}

func (r *PersonalInformationRequest) Validate(requestPayload string) error {

	if err := json.Unmarshal([]byte(requestPayload), r); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}

	return nil
}

func (r *CreateBankAccountRequest) Validate(requestPayload string) error {
	if err := json.Unmarshal([]byte(requestPayload), r); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}

	if !r.IsAddrSameAsAdhaar && r.CommunicationAddress == nil {
		return errors.New("communication_address cannot be nil when is_addr_same_as_adhaar is false")
	}

	return nil
}

func NewPincodeDetails() *PincodeDetails {
	return &PincodeDetails{}
}

func (p *PincodeDetails) Validate(requestPayload string) error {
	if err := json.Unmarshal([]byte(requestPayload), p); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(p); err != nil {
		return err
	}

	return nil
}

func (r *EncryptedRequest) Validate(c *gin.Context) error {
	if err := customvalidation.ValidatePayload(c, r); err != nil {
		return err
	}

	return nil
}

func (r *EncryptedRequest) GetDecryptedInfo(singingKey string) (string, error) {

	decrypted, err := security.Decrypt(r.Data, []byte(singingKey))
	if err != nil {
		return "", err
	}
	return decrypted, nil
}

func (r *CommunicationAddress) BindForAddressUpdate(req *AddressUpdateRequest) error {

	r.HouseNo = req.Address1
	r.StreetName = req.Address2
	r.Locality = req.Address3
	r.Landmark = ""
	r.PinCode = req.PinCode
	r.State = req.State
	r.City = req.City

	return nil
}
