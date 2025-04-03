package requests

import (
	"encoding/json"
	"errors"

	"bitbucket.org/paydoh/paydoh-commons/customvalidation"
	"github.com/gin-gonic/gin"
)

type AuthorizationRequest struct {
	MobileNumber string `json:"mobile_number,omitempty" validate:"required,numeric,min=10,max=10"`
	UserId       string `json:"user_id,omitempty"`
}

func NewAuthorizationRequest() *AuthorizationRequest {
	return &AuthorizationRequest{}
}

func (request *AuthorizationRequest) Validate(c *gin.Context) error {
	if err := customvalidation.ValidatePayload(c, request); err != nil {
		return err
	}

	if request.UserId == "" && request.MobileNumber == "" {
		return errors.New("user_id or mobile_number cannot be empty")
	}

	return nil
}

func (request *AuthorizationRequest) ToJSON() ([]byte, error) {
	return json.Marshal(request)
}
