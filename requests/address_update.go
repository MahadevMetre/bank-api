package requests

import (
	"mime/multipart"
)

type AddressUpdateRequest struct {
	AddressProof  string                `form:"address_proof" validate:"required"`
	Address1      string                `form:"address1" validate:"required"`
	Address2      string                `form:"address2" validate:"required"`
	Address3      string                `form:"address3" validate:"required"`
	City          string                `form:"city" validate:"required"`
	PinCode       string                `form:"pincode" validate:"required"`
	CityCode      string                `form:"citycode" validate:"required"`
	State         string                `form:"state" validate:"required"`
	Country       string                `form:"country" validate:"required"`
	AddressType   string                `form:"addressType" validate:"required"`
	DocumentFront *multipart.FileHeader `form:"document_front" validate:"required"`
	DocumentBack  *multipart.FileHeader `form:"document_back" validate:"required"`
}

func NewAddressUpdateRequest() *AddressUpdateRequest {
	return &AddressUpdateRequest{}
}
