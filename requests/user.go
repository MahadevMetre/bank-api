package requests

import (
	"bitbucket.org/paydoh/paydoh-commons/customvalidation"
	"github.com/gin-gonic/gin"
)

type UserRequest struct {
	UserId string `json:"user_id" validate:"required,max=20"`
}

type UserMobileNumber struct {
	MobileNumber string `json:"mobile_number" validate:"required,numeric,min=10,max=10"`
}

func NewUserRequest() *UserRequest {
	return &UserRequest{}
}

func NewUserMobileNumber() *UserMobileNumber {
	return &UserMobileNumber{}
}

func (request *UserRequest) Validate(c *gin.Context) error {
	if err := customvalidation.ValidatePayload(c, request); err != nil {
		return err
	}

	return nil
}

func (request *UserMobileNumber) Validate(c *gin.Context) error {
	if err := customvalidation.ValidatePayload(c, request); err != nil {
		return err
	}

	return nil
}

type UpdateFcmToken struct {
	FcmToken string `json:"fcm_token" validate:"required"`
}

func NewUpdateFcmToken() *UpdateFcmToken {
	return &UpdateFcmToken{}
}

func (i *UpdateFcmToken) Validate(c *gin.Context) error {
	if err := customvalidation.ValidatePayload(c, i); err != nil {
		return err
	}

	return nil
}
