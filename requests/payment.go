package requests

import (
	"encoding/json"

	"bitbucket.org/paydoh/paydoh-commons/customvalidation"
)

type GetReceiptIdRequest struct {
	GatewayId string `json:"gateway_id" validate:"required"`
	Amount    string `json:"amount" validate:"required"`
	Currency  string `json:"currency" validate:"required"`
	Remarks   string `json:"remarks" validate:"required"`
}

type AddPaymentStatusRequest struct {
	TransactionStatus string `json:"transaction_status" validate:"required"`
	TransactionId     string `json:"transaction_id" validate:"required"`
	ReceiptId         string `json:"receipt_id" validate:"required"`
}

func NewGatewayReceiptIdRequest() *GetReceiptIdRequest {
	return &GetReceiptIdRequest{}
}

func NewAddPaymentStatusRequest() *AddPaymentStatusRequest {
	return &AddPaymentStatusRequest{}
}

func (r *GetReceiptIdRequest) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *AddPaymentStatusRequest) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *GetReceiptIdRequest) Validate(payload string) error {
	if err := r.Unmarshal([]byte(payload)); err != nil {
		return err
	}
	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}
	return nil
}

func (r *AddPaymentStatusRequest) Validate(payload string) error {
	if err := r.Unmarshal([]byte(payload)); err != nil {
		return err
	}
	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}
	return nil
}

type OutgoingGetReceiptIdRequest struct {
	GatewayId     string `json:"gateway_id"`
	UserId        string `json:"user_id"`
	ApplicationId string `json:"application_id"`
	Amount        string `json:"amount"`
	Currency      string `json:"currency"`
	Remarks       string `json:"remarks"`
	PaymentType   string `json:"payment_type"`
}

type OutgoingPaymentStatusRequest struct {
	ReceiptId         string `json:"receipt_id"`
	TransactionId     string `json:"transaction_id"`
	StatusId          uint32 `json:"status_id"`
	TransactionStatus string `json:"transaction_status"`
	TxnTimestamp      string `json:"txn_timestamp"`
}

type OutgoingGetDebitCardPaymentStatusRequest struct {
	UserId string `json:"user_id"`
}
