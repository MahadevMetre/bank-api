package requests

import (
	"bankapi/constants"
	"bankapi/security"
	"encoding/json"
	"strings"

	"bitbucket.org/paydoh/paydoh-commons/customvalidation"
)

type VerifyEmailReq struct {
	EmailId string `json:"email" validate:"required,email"`
	Name    string `json:"name" validate:"required"`
	UserId  string `json:"userId"`
}

func NewVerifyEmailReq() *VerifyEmailReq {
	return &VerifyEmailReq{}
}

func ValidateDomain(email string) bool {
	domain := ""
	if parts := strings.Split(email, "@"); len(parts) == 2 {
		domain = parts[1]
	}

	for _, allowedDomain := range constants.AllowedDomains {
		if domain == allowedDomain {
			return true
		}
	}

	return false
}

func (r *VerifyEmailReq) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *VerifyEmailReq) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *VerifyEmailReq) Validate(requestPayload string) error {

	if err := json.Unmarshal([]byte(requestPayload), r); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}

	return nil
}

type SendVerificationMailReq struct {
	Email string `json:"email" validate:"required"`
	TxnId string `json:"txnId" validate:"required"`
	Name  string `json:"name" validate:"required"`
}

func NewSendVerificationMailReq() *SendVerificationMailReq {
	return &SendVerificationMailReq{}
}

func (r *SendVerificationMailReq) Bind(email, name string) error {
	txnID, err := security.GenerateRandomUUID(40)
	if err != nil {
		return err
	}
	r.Email = email
	r.Name = name
	r.TxnId = strings.ReplaceAll(txnID, "-", "")

	return nil
}

type SendAccountInformationReq struct {
	Email             string `json:"email" validate:"required"`
	AccountHolderName string `json:"account_holder_name" validate:"required"`
	AccountNumber     string `json:"account_no" validate:"required"`
	AccountType       string `json:"account_type" validate:"required"`
	BranchName        string `json:"branch_name" validate:"required"`
	IfscCode          string `json:"ifsc_code" validate:"required"`
}

func NewSendAccountInformationReq() *SendAccountInformationReq {
	return &SendAccountInformationReq{}
}

func (request *SendAccountInformationReq) Bind(accountNumber, name, email, accountType, branchName, ifsc string) error {
	request.Email = email
	request.AccountHolderName = name
	request.AccountNumber = accountNumber
	request.AccountType = accountType
	request.BranchName = branchName
	request.IfscCode = ifsc
	return nil
}
