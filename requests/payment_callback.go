package requests

import (
	"bitbucket.org/paydoh/paydoh-commons/customvalidation"
	"github.com/gin-gonic/gin"
)

type CbsStatus struct {
	ErrorCode         string `json:"ErrorCode"`
	ErrorMessage      string `json:"ErrorMessage"`
	ApplicationNumber string `json:"ApplicationNumber"`
	TransactionId     string `json:"TransactionId"`
	Status            string `json:"Status"`
	RetryFlag         string `json:"RetryFlag"`
	UTR_Ref_Number    string `json:"UTR_Ref_Number"`
}

type PaymentCallbackRequestData struct {
	ServiceName   string      `json:"serviceName"`
	ApplicationId string      `json:"ApplicationId"`
	CbsStatus     []CbsStatus `json:"CbsStatus"`
	ProductType   string      `json:"ProductType"`
	SourcedBy     string      `json:"SourcedBy"`
	CallBackName  string      `json:"CallBackName"`
}

func NewPaymentCallbackRequestData() *PaymentCallbackRequestData {
	return &PaymentCallbackRequestData{}

}

func (request *PaymentCallbackRequestData) Validate(c *gin.Context) error {
	if err := customvalidation.ValidatePayload(c, request); err != nil {
		return err
	}

	return nil
}
