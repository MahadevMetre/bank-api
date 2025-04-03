package requests

import (
	"bitbucket.org/paydoh/paydoh-commons/amazon"
	"bitbucket.org/paydoh/paydoh-commons/customvalidation"
	"errors"
	"github.com/gin-gonic/gin"
	"strings"
)

type AddShippingAddress struct {
	DocumentType string `json:"document_type" validate:"required"`
	Document     string `json:"document" validate:"required,url"`
	AddressLine1 string `json:"address_line_1" validate:"required"`
	StreetName   string `json:"street_name" validate:"required"`
	Locality     string `json:"locality" validate:"required"`
	Landmark     string `json:"landmark"`
	City         string `json:"city" validate:"required"`
	State        string `json:"state" validate:"required"`
	PinCode      string `json:"pin_code" validate:"required"`
	Country      string `json:"country" validate:"required"`
}

type UpdateShippingAddress struct {
	DocumentType string `json:"document_type,omitempty"`
	Document     string `json:"document,omitempty"`
	AddressLine1 string `json:"address_line_1,omitempty"`
	StreetName   string `json:"street_name,omitempty"`
	Locality     string `json:"locality,omitempty"`
	Landmark     string `json:"landmark,omitempty"`
	City         string `json:"city,omitempty"`
	State        string `json:"state,omitempty"`
	PinCode      string `json:"pin_code,omitempty"`
	Country      string `json:"country,omitempty"`
}

func NewAddShippingAddress() *AddShippingAddress {
	return &AddShippingAddress{}
}

func NewUpdateShippingAddress() *UpdateShippingAddress {
	return &UpdateShippingAddress{}
}

func (request *AddShippingAddress) Validate(c *gin.Context, aws *amazon.Aws) error {
	if !strings.Contains(c.Request.Header.Get("Content-Type"), "multipart/form-data") {
		return errors.New("invalid content type, use multipart/form-data")
	}

	documentLocation, err := amazon.ProcessAndSaveImageToAWS(
		c,
		"document",
		"document",
		"documents",
		1,
		aws,
	)

	if err != nil {
		return err
	}

	request.Document = documentLocation

	request.DocumentType = c.PostForm("document_type")
	request.AddressLine1 = c.PostForm("address_line_1")
	request.StreetName = c.PostForm("street_name")
	request.Locality = c.PostForm("locality")
	request.Landmark = c.PostForm("landmark")
	request.City = c.PostForm("city")
	request.State = c.PostForm("state")
	request.PinCode = c.PostForm("pin_code")
	request.Country = c.PostForm("country")

	if err := customvalidation.ValidateStruct(request); err != nil {
		return err
	}

	return nil
}

func (r *UpdateShippingAddress) Validate(c *gin.Context, aws *amazon.Aws) error {
	if !strings.Contains(c.Request.Header.Get("Content-Type"), "multipart/form-data") {
		return errors.New("invalid content type, use multipart/form-data")
	}

	documentLocation, _ := amazon.ProcessAndSaveImageToAWS(
		c,
		"document",
		"document",
		"documents",
		1,
		aws,
	)

	if documentLocation != "" {
		r.Document = documentLocation
	}

	if c.PostForm("document_type") != "" {
		r.DocumentType = c.PostForm("document_type")
	}

	if c.PostForm("address_line_1") != "" {
		r.AddressLine1 = c.PostForm("address_line_1")
	}

	if c.PostForm("street_name") != "" {
		r.StreetName = c.PostForm("street_name")
	}

	if c.PostForm("locality") != "" {
		r.Locality = c.PostForm("locality")
	}

	if c.PostForm("landmark") != "" {
		r.Landmark = c.PostForm("landmark")
	}

	if c.PostForm("city") != "" {
		r.City = c.PostForm("city")
	}

	if c.PostForm("state") != "" {
		r.State = c.PostForm("state")
	}

	if c.PostForm("pin_code") != "" {
		r.PinCode = c.PostForm("pin_code")
	}

	if c.PostForm("country") != "" {
		r.Country = c.PostForm("country")
	}

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}

	return nil
}
