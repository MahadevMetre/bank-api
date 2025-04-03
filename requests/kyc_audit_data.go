package requests

import (
	"bitbucket.org/paydoh/paydoh-commons/customvalidation"
	"github.com/gin-gonic/gin"
)

type KycAuditRequestData struct {
	UserId            string `json:"UserId"`
	MobileNo          string `json:"MobileNo"`
	CallBackName      string `json:"CallBackName"`
	ApplicantId       string `json:"ApplicantId"`
	SourcedBy         string `json:"SourcedBy"`
	VKYCAuditStatus   string `json:"VKYCAuditStatus"`
	ProductType       string `json:"ProductType"`
	AuditRejectReason string `json:"AuditRejectReason"`
}

func NewwKycAuditRequestData() *KycAuditRequestData {
	return &KycAuditRequestData{}

}

func (request *KycAuditRequestData) Validate(c *gin.Context) error {
	if err := customvalidation.ValidatePayload(c, request); err != nil {
		return err
	}

	return nil
}
