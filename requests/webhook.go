package requests

import (
	"encoding/json"
	"strings"

	"bitbucket.org/paydoh/paydoh-commons/customvalidation"
	"github.com/gin-gonic/gin"
)

type WebhookRequest struct {
	Sender    string `form:"sender" json:"sender" validate:"required"`
	Message   string `form:"message" json:"message" validate:"required"`
	Operator  string `form:"operator" json:"operator"`
	Circle    string `form:"circle" json:"circle"`
	Coding    string `form:"coding" json:"coding"`
	Timestamp string `form:"timestamp" json:"timestamp"`
}

type KycDataUpdateRequest struct {
	UserId string `form:"user_id" json:"user_id" validate:"required"`
	Status string `form:"status" json:"status" validate:"required"`
	AStat  string `form:"astat" json:"astat" validate:"required"`
	Acom   string `form:"acom" json:"acom" validate:"required"`
}

func NewWebhookRequest() *WebhookRequest {
	return &WebhookRequest{}
}

func NewKycUpdateRequest() *KycDataUpdateRequest {
	return &KycDataUpdateRequest{}
}

func (request *WebhookRequest) Validate(c *gin.Context) error {
	if err := customvalidation.ValidateQuery(c, request); err != nil {
		return err
	}

	cleanedNumber := request.Sender
	if strings.HasPrefix(request.Sender, "91") && len(request.Sender) == 12 {
		cleanedNumber = request.Sender[2:] // Remove the first two characters ("91")
	}

	request.Sender = cleanedNumber

	return nil
}

func (request *KycDataUpdateRequest) Validate(c *gin.Context) error {
	if err := customvalidation.ValidatePayload(c, request); err != nil {
		return err
	}

	return nil
}

func (request *KycDataUpdateRequest) Marshal() ([]byte, error) {
	return json.Marshal(request)
}

func (request *KycDataUpdateRequest) UnMarshal(data []byte) error {
	return json.Unmarshal(data, request)
}
