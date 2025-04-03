package responses

import (
	"bitbucket.org/paydoh/paydoh-commons/pkg/json"
)

type GetDebitCardPaymentStatusResponse struct {
	Data    PaymentStatusData `json:"data"`
	Message string            `json:"message"`
	Status  uint64            `json:"status"`
}

type PaymentStatusData struct {
	Status  uint64                          `json:"status"`
	Message string                          `json:"message"`
	Data    []GetDebitCardPaymentStatusData `json:"data"`
}

type GetDebitCardPaymentStatusData struct {
	ReceiptID     string `json:"receipt_id"`
	OrderID       string `json:"order_id"`
	TransactionID string `json:"transaction_id"`
	TxnTimestamp  string `json:"txn_timestamp"`
	GatewayID     int64  `json:"gateway_id"`
	UserID        string `json:"user_id"`
	ApplicationID string `json:"application_id"`
	Amount        string `json:"amount"`
	Currency      string `json:"currency"`
	Remarks       string `json:"remarks"`
	PaymentType   string `json:"payment_type"`
	StatusID      int32  `json:"status_id"`
	TxnStatus     string `json:"txn_status"`
}

func NewGetDebitCardPaymentStatusResponse() *GetDebitCardPaymentStatusResponse {
	return &GetDebitCardPaymentStatusResponse{}
}

func (r *GetDebitCardPaymentStatusResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *GetDebitCardPaymentStatusResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type ReciptIDResponse struct {
	Data struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Data    struct {
			ReceiptID string `json:"receipt_id"`
			OrderID   string `json:"order_id"`
		} `json:"data"`
	} `json:"data"`
}

func NewGetReciptIDResponse() *ReciptIDResponse {
	return &ReciptIDResponse{}
}

func (r *ReciptIDResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *ReciptIDResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type GatewayServiceResponse struct {
	Data    *GatewayServiceData `json:"data"`
	Status  int                 `json:"status"`
	Message string              `json:"message"`
}

type GatewayServiceData struct {
	ReceiptId string `json:"receipt_id"`
	OrderId   string `json:"order_id"`
}

func NewGatewayServiceResponse() *GatewayServiceResponse {
	return &GatewayServiceResponse{
		Data: &GatewayServiceData{}, // Initialize Data here
	}
}

func (r *GatewayServiceResponse) Bind(res ReciptIDResponse) error {
	r.Data.OrderId = res.Data.Data.OrderID
	r.Data.ReceiptId = res.Data.Data.ReceiptID
	r.Status = res.Data.Status
	r.Message = res.Data.Message
	return nil
}

type PaymentStatusResponse struct {
	TxnStatus           string `json:"TxnStatus"`
	TxnIdentifier       string `json:"TxnIdentifier"`
	TxnRefNo            string `json:"TxnRefNo"`
	PaydohTransactionId string `json:"PaydohTransactionId"`
	TransactionAmount   string `json:"TxnAmount"`
}

func NewPaymentStatusResponse() *PaymentStatusResponse {
	return &PaymentStatusResponse{}
}

type EncryptRes struct {
	EncryptRes string `json:"encrypt_res"`
}

func NewEncryptResponse() *EncryptRes {
	return &EncryptRes{}
}

func (r *EncryptRes) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}
