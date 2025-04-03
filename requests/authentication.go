package requests

import (
	"bankapi/security"
	"encoding/json"

	"bitbucket.org/paydoh/paydoh-commons/customvalidation"
	"github.com/gin-gonic/gin"
)

type AuthenticationRequest struct {
	Data     string `json:"data" validate:"required"`
	CrcValue string `json:"crc_value"`
}

type DataRequest struct {
	PackageId   string `json:"package_id" validate:"required"`
	DeviceToken string `json:"device_token" validate:"required"`
	DeviceId    string `json:"device_id" validate:"required"`
	SimVendorId string `json:"sim_vendor_id"`
}

func NewAuthenticationRequest() *AuthenticationRequest {
	return &AuthenticationRequest{}
}

func NewDataRequest() *DataRequest {
	return &DataRequest{}
}

func (request *AuthenticationRequest) Validate(c *gin.Context) error {
	if err := customvalidation.ValidatePayload(c, request); err != nil {
		return err
	}

	return nil
}

func (request *DataRequest) validate() error {
	if err := customvalidation.ValidateStruct(request); err != nil {
		return err
	}

	return nil
}

func (request *AuthenticationRequest) DecrypToData(paraphrase string) (*DataRequest, error) {

	dataRequest := &DataRequest{}

	decryptedData, err := security.Decrypt(request.Data, []byte(paraphrase))

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(decryptedData), dataRequest); err != nil {
		return nil, err
	}

	if err := dataRequest.validate(); err != nil {
		return nil, err
	}

	return dataRequest, nil
}
