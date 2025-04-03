package responses

import "bitbucket.org/paydoh/paydoh-commons/pkg/json"

type Transaction struct {
	ChequeNumber           string `json:"ChequeNumber"`
	CodeDRCR               string `json:"CodeDRCR"`
	PostingDate            string `json:"PostingDate"`
	RunningTotal           string `json:"RunningTotal"`
	TransactionAmount      string `json:"TransactionAmount"`
	TransactionDate        string `json:"TransactionDate"`
	TransactionDescription string `json:"TransactionDescription"`
	ValueDate              string `json:"ValueDate"`
	BranchCode             string `json:"BranchCode"`
	BranchName             string `json:"BranchName"`
	MccCode                string `json:"MccCode"`
	MnemonicCode           string `json:"MnemonicCode"`
	MnemonicDesc           string `json:"MnemonicDesc"`
}

type TxnBankError struct {
	ErrorCode    string
	ErrorMessage string
}

type TransactionResponse struct {
	Data []Transaction `json:"CasaTransactionDetails"`
	TxnBankError
}

func NewTransactionResponse() *TransactionResponse {
	return &TransactionResponse{}
}

func (r *TransactionResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *TransactionResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type EncryptedRes struct {
	EncryptRes string `json:"encryptRes"`
}

func NewEncryptedResponse() *EncryptedRes {
	return &EncryptedRes{}
}

func (r *EncryptedRes) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type CustomUpiAndBankError struct {
	ErrorCode        string       `json:"ErrorCode"`
	ErrorMessage     string       `json:"ErrorMessage"`
	HttpStatus       string       `json:"HttpStatus"`
	ErrorCodeRC      string       `json:"rc"`
	ErrorMessageDesc string       `json:"desc"`
	UpiError         UpiBaseError `json:"response"`
	DebitCardError   []map[string]interface{}
}

type UpiBaseError struct {
	ResponseCode    string `json:"ResponseCode"`
	ResponseMessage string `json:"ResponseMessage"`
}

func NewCustomUpiAndBankError() *CustomUpiAndBankError {
	return &CustomUpiAndBankError{}
}

func (r *CustomUpiAndBankError) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}
