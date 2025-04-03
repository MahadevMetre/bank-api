package requests

import (
	"encoding/json"

	"bitbucket.org/paydoh/paydoh-commons/customvalidation"
	"github.com/gin-gonic/gin"
)

type TransactionRequest struct {
	UserId   string `json:"user_id,omitempty"`
	FromDate string `json:"from_date"`
	ToDate   string `json:"to_date"`
	Limit    int    `json:"limit"`
}

func NewTransactionRequest() *TransactionRequest {
	return &TransactionRequest{}
}

func (r *TransactionRequest) Validate(c *gin.Context) error {
	if err := customvalidation.ValidatePayload(c, r); err != nil {
		return err
	}

	return nil
}

func (r *TransactionRequest) ValidateEncrypted(payload string) error {
	if err := json.Unmarshal([]byte(payload), r); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}

	return nil
}

type TransactionDetailsRequest struct {
	CodeDRCR               string `json:"CodeDRCR"`
	TransactionAmount      string `json:"TransactionAmount"`
	TransactionDate        string `json:"TransactionDate"`
	TransactionDescription string `json:"TransactionDescription"`
}

func (r *TransactionDetailsRequest) Validate(c *gin.Context) error {
	if err := customvalidation.ValidatePayload(c, r); err != nil {
		return err
	}

	return nil
}

func (r *TransactionDetailsRequest) ValidateEncrypted(payload string) error {
	if err := json.Unmarshal([]byte(payload), r); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}

	return nil
}

// bank
type KVBTransactionRequest struct {
	TxnBranch        string `json:"TxnBranch" validate:"required"`
	TxnIdentifier    string `json:"TxnRefNo" validate:"required"`
	KVBAccountNumber string `json:"AccountNumber" validate:"required"`
	FromDate         string `json:"FromDate" validate:"required"`
	ToDate           string `json:"ToDate" validate:"required"`
}

func (r *KVBTransactionRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// UPI transaction
type UpiTransaction struct {
	CredType string `json:"cred_type" validate:"required"`
}

func NewUpiTransactionRequest() *UpiTransaction {
	return &UpiTransaction{}
}

func (r *UpiTransaction) Validate(payload string) error {
	if err := json.Unmarshal([]byte(payload), r); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}
	return nil
}

type EncryptedReq struct {
	EncryptReq string `json:"encryptReq"`
}

func (req *EncryptedReq) Bind(enc string) (EncryptedReq, error) {
	req.EncryptReq = enc
	return *req, nil
}

func (r *EncryptedReq) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
