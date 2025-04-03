package requests

import (
	"encoding/json"

	"bitbucket.org/paydoh/paydoh-commons/customvalidation"
	"github.com/gin-gonic/gin"
)

type ConsentRequest struct {
	ConsentType     string `json:"consent_type" validate:"required,max=20,oneof=FATCA PEP AADHAR1 LOCATION_TRACKING PRIVACY_POLICY TERMS_AND_CONDITIONS AADHAR2 ADDRESS_CHANGE NOMINATION"`
	ConsentProvided string `json:"consent_provided" validate:"required,max=3"`
}

func NewConsentRequest() *ConsentRequest {
	return &ConsentRequest{}
}

func (r *ConsentRequest) Validate(c *gin.Context) error {
	if err := customvalidation.ValidatePayload(c, r); err != nil {
		return err
	}

	return nil
}

func (r *ConsentRequest) DecryptAndValidatePayload(payload string) error {
	if err := r.Unmarshal([]byte(payload)); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}

	return nil
}

func (r *ConsentRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *ConsentRequest) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}
