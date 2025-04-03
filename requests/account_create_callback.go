package requests

import (
	"bitbucket.org/paydoh/paydoh-commons/customvalidation"
	"github.com/gin-gonic/gin"
)

type AccountCreateCallbackEncryptedRequest struct {
	SourcedBy    string `json:"SourcedBy"`
	ProductType  string `json:"ProductType"`
	CallbackName string `json:"CallBackName"`
	EncryptedRes string `json:"encryptRes"`
}

type AccountCreateCallbackRequest struct {
	ServiceName   string      `json:"serviceName" validate:"required"`
	ApplicationId string      `json:"ApplicationId" validate:"required"`
	CBSStatus     []CBSStatus `json:"CbsStatus"`
}

type CBSStatus struct {
	Status      string `json:"Status"`
	SuccErrCode string `json:"succErrCode"`
	AccountNo   string `json:"AccountNo,omitempty"`
	ApplicantID string `json:"ApplicantId"`
	CustomerID  string `json:"CustomerId,omitempty"`
	Message     string `json:"message"`
}

func NewAccountCreateCallbackRequest() *AccountCreateCallbackRequest {
	return &AccountCreateCallbackRequest{}

}

func (request *AccountCreateCallbackRequest) Validate(c *gin.Context) error {
	if err := customvalidation.ValidatePayload(c, request); err != nil {
		return err
	}

	return nil
}
