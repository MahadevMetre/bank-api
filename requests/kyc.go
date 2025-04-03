package requests

import (
	"encoding/json"

	"bitbucket.org/paydoh/paydoh-commons/customvalidation"
	"github.com/gin-gonic/gin"
)

type KycConsentRequest struct {
	IndianResident           bool `json:"indian_resident"`
	NotPoliticallyExposed    bool `json:"not_politically_exposed"`
	AadharConsent            bool `json:"aadhar_consent" `
	VirtualDebitCardConsent  bool `json:"virtual_card_consent" `
	PhysicalDebitCardConsent bool `json:"physical_card_consent"`
	Aadhar2Consent           bool `json:"aadhar2_consent"`
	AddressChangeConsent     bool `json:"address_change_consent"`
	NominationConsent        bool `json:"nomination_consent"`
	LocationConsent          bool `json:"location_consent"`
	PrivacyPolicyConsent     bool `json:"privacy_policy_consent"`
	TermsAndCondition        bool `json:"terms_and_condition"`
}

type IncomingVcipData struct {
	ApplicantId           string `json:"ApplicantId" validate:"required"`
	MobileNumber          string `json:"MobileNo" validate:"required,numeric,min=10,max=10"`
	PanNumber             string `json:"PaNNumber" validate:"required"`
	Firstname             string `json:"FirstName" validate:"required"`
	MiddleName            string `json:"MiddleName" validate:"required"`
	LastName              string `json:"LastName" validate:"required"`
	AadharReferenceNumber string `json:"AadharReferenceNo" validate:"required"`
	VKYCCompletion        string `json:"VKYCCompletion" validate:"required"`
	VKYCAuditStatus       string `json:"VKYCAuditStatus" validate:"required"`
	AuditorRejectRemarks  string `json:"AuditorRejectRemarks" validate:"required"`
}

func NewKycConsentRequest() *KycConsentRequest {
	return &KycConsentRequest{}
}

func NewIncomingVcipData() *IncomingVcipData {
	return &IncomingVcipData{}
}

func (k *KycConsentRequest) Validate(payload string) error {
	if err := json.Unmarshal([]byte(payload), k); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(k); err != nil {
		return err
	}

	// if k.IndianResident && k.NotPoliticallyExposed && k.AadharConsent {
	// 	return errors.New("kyc consent is required")
	// }

	return nil
}

func (i *IncomingVcipData) Validate(c *gin.Context) error {
	if err := customvalidation.ValidatePayload(c, i); err != nil {
		return err
	}

	return nil
}
