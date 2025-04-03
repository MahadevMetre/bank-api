package responses

import (
	"bitbucket.org/paydoh/paydoh-commons/pkg/json"
)

type FetchAPIKeyResponse struct {
	ResId     string `json:resid`
	ShashKey  string `json:shashkey`
	Rc        string `json:rc`
	Desc      string `json:desc`
	Publickey string `json:publickey`
}

func NewFetchAPIKeyResponse() *FetchAPIKeyResponse {
	return &FetchAPIKeyResponse{}
}

func (r *FetchAPIKeyResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *FetchAPIKeyResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type UserLoginAndRegistrationResponse struct {
	CustomerID   string `json:cid`
	ErrorMessage string `json:desc`
	ErrorCode    string `json:rc`
}

func NewUserLoginRegistrationResponse() *UserLoginAndRegistrationResponse {
	return &UserLoginAndRegistrationResponse{}
}

func (r *UserLoginAndRegistrationResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *UserLoginAndRegistrationResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type GenerateDebitcardResponse struct {
	ErrorCode     string `json:"ErrorCode"`
	ErrorMessage  string `json:"ErrorMessage"`
	ApplicantId   string `json:"ApplicantId"`
	TxnIdentifier string `json:"TxnIdentifier"`
	ProxyNumber   string `json:"proxyNumber"`
	Cvv           string `json:"cvv"`
}

func NewGenerateDebitcardResponse() *GenerateDebitcardResponse {
	return &GenerateDebitcardResponse{}
}

func (r *GenerateDebitcardResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *GenerateDebitcardResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type DebitCardDetailRes struct {
	ApplicantId        string `json:"ApplicantId"`
	AccountNo          string `json:"AccountNo"`
	ProxyNumber        string `json:"proxyNumber"`
	EncryptedPAN       string `json:"encryptedPAN"`
	ExpiryDate         string `json:"expiryDate"`
	CardholderName     string `json:"cardholderName"`
	CardIssuanceStatus int    `json:"cardIssuanceStatus"`
	CvvValue           string `json:"cvvValue"`

	HasPhysicalDebitCard bool `json:"has_physical_debit_card"`
	HasVirtualDebitCard  bool `json:"has_virtual_debit_card"`
	IsPermanentlyBlocked bool `json:"is_permanently_blocked"`
}

func NewDebitCardDetailRes() *DebitCardDetailRes {
	return &DebitCardDetailRes{}
}

func (s *DebitCardDetailRes) Bind(res DebitcardDetailResponse, isPhysicalGenerated, isVirtualGenerated, isPermanentlyBlocked bool) error {

	s.ApplicantId = res.ApplicantId
	s.AccountNo = res.AccountNo
	s.ProxyNumber = res.ServiceData.CardData[0].ProxyNumber
	s.EncryptedPAN = res.ServiceData.CardData[0].EncryptedPAN
	s.ExpiryDate = res.ServiceData.CardData[0].ExpiryDate
	s.CardholderName = res.ServiceData.CardData[0].CardholderName
	s.CardIssuanceStatus = res.ServiceData.CardData[0].CardIssuanceStatus
	s.CvvValue = res.ServiceData.CardData[0].CvvValue
	s.IsPermanentlyBlocked = isPermanentlyBlocked
	s.HasPhysicalDebitCard = isPhysicalGenerated
	s.HasVirtualDebitCard = isVirtualGenerated
	return nil
}

type GeneratePhysicalDebitCardRes struct {
	ErrorCode     string `json:"ErrorCode"`
	ErrorMessage  string `json:"ErrorMessage"`
	ApplicantId   string `json:"ApplicantId"`
	TxnIdentifier string `json:"TxnIdentifier"`
	ProxyNumber   string `json:"proxyNumber"`
	AccountNo     string `json:"AccountNo"`
	PinIdentifier string `json:"PinIdentifier"`
}

func NewGeneratePhysicalDebitCardRes() *GeneratePhysicalDebitCardRes {
	return &GeneratePhysicalDebitCardRes{}
}

func (res *GeneratePhysicalDebitCardRes) Marshal() ([]byte, error) {
	return json.Marshal(res)
}

func (res *GeneratePhysicalDebitCardRes) UnMarshal(data []byte) error {
	return json.Unmarshal(data, res)
}

type GetTransactionLimitRes struct {
	// TransactionId string            `json:"transactionID"`
	ReqData []TxnResponseData `json:"data"`
}

type TxnResponseData struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	MaxLimit     string `json:"maxLimit"`
	SetValue     string `json:"setlimit"`
	TranTypes    string `json:"tranTypes"`
	StatusDc     string `json:"statusDc"`     //Delivery Channel Status
	DChBlkStatus string `json:"dchBlkStatus"` //Delivery Channel Block status 0 Unblock 1 block
	TransMStatus string `json:"transtatus"`   //Delivery Channel Block status 0 Unblock 1 block
}

func NewGetTransactionLimit() GetTransactionLimitRes {
	return GetTransactionLimitRes{}
}

func (res *GetTransactionLimitRes) StrucherRes(resp FetchTransactionResponse, txnType, transactionId string) *GetTransactionLimitRes {
	list := []TxnResponseData{}
	for _, v := range resp.TranList {
		c := TxnResponseData{}
		c.Name = v.DeliveryChannel // + " " + v.TranM[0].TranLabel
		c.Type = txnType
		c.MaxLimit = v.TranM[0].Max
		c.SetValue = v.TranM[0].Value
		c.TranTypes = v.TranTypes
		c.StatusDc = v.StatusDc
		c.DChBlkStatus = v.DchBlkStatus
		c.TransMStatus = v.TranM[0].Status
		list = append(list, c)
	}
	res.ReqData = list
	// res.TransactionId = transactionId
	return res
}

type EncryptedResponse struct {
	Response string `json:"encryptRes"`
	BankError
}

type ErrorResponse struct {
	ErrorCode    string `json:"ErrorCode"`
	ErrorMessage string `json:"ErrorMessage,omitempty"`
	HttpStatus   string `json:"HttpStatus,omitempty"`
}

func NewDebitCardEncryptedResponse() *EncryptedResponse {
	return &EncryptedResponse{}
}

func (r *EncryptedResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *EncryptedResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type ExternalKeyFetch struct {
	StatusCode string `json:"status code"`
	PublicKey  string `json:"publickey"`
	ResID      string `json:"resid"`
	ShashKey   string `json:"shashkey"`
	RC         string `json:"rc"`
	Desc       string `json:"desc"`
}

type LoginResponse struct {
	CustomerID string `json:"cid"`
	RC         string `json:"rc"`
	Desc       string `json:"desc"`
}

type AddCardResponse struct {
	RC                string `json:"rc"`
	ENID              string `json:"enid"`
	BrandType         string `json:"bt"`
	PaymentInstrument string `json:"pi"`
	OTPID             string `json:"otpid"`
	HashData          string `json:"hashdata"`
	Desc              string `json:"desc"`
}

type ListCardResponse struct {
	RC             string `json:"rc"`
	ENID           string `json:"enid"`
	MaskedCardNo   string `json:"mcn"`
	CardHolderName string `json:"cnm"`
	Status         string `json:"stat"`
	ExpiryDuration string `json:"exp"`
	BrandType      string `json:"bt"`
	CardType       string `json:"ct"`
	Desc           string `json:"desc"`
}

type ListCardControlResponses []ListCardControlResponse

type ListCardControlResponse struct {
	RC                       string `json:"rc"`
	CNID                     string `json:"cnid"`
	ENID                     string `json:"enid"`
	CID                      string `json:"cid"`
	CardHolderName           string `json:"cnm"`
	Status                   string `json:"stat"`
	NCFlag                   string `json:"ncflag"`
	BlockStatus              string `json:"bs"`
	InternationalBlockStatus string `json:"ibs"`
	Desc                     string `json:"desc"`
	Lang                     string `json:"lang"`
}

type FetchTransactionResponse struct {
	RC       string       `json:"rc"`
	Desc     string       `json:"desc"`
	TranList []TranDetail `json:"tranList"`
}
type TranDetail struct {
	DeliveryChannel string         `json:"deliveryChannel"`
	DeliveryIndex   string         `json:"deliveryIndex"`
	TranTypes       string         `json:"tranTypes"`
	StatusDc        string         `json:"statusDc"`
	DchBlkStatus    string         `json:"dchBlkStatus"`
	TranM           []DTransaction `json:"tranM"`
	Minordigit      string         `json:"minordigit"`
	Symbol          string         `json:"symbol"`
}

type DTransaction struct {
	Max       string `json:"max"`
	Status    string `json:"status"`
	Value     string `json:"value"`
	TranValue string `json:"tranValue"`
	TranLabel string `json:"tranLabel"`
}

type EditTransactionResponse struct {
	RC   string `json:"rc"`
	Desc string `json:"desc"`
}

type CardBlockResponse struct {
	RC   string `json:"rc"`
	Desc string `json:"desc"`
}

func (r *ExternalKeyFetch) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *ExternalKeyFetch) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *LoginResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}
func (r *AddCardResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *ListCardResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *FetchTransactionResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}
func (r *EditTransactionResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *CardBlockResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}
func NewKeyFetchResponse() *ExternalKeyFetch {
	return &ExternalKeyFetch{}
}
func NewLoginResponse() *LoginResponse {
	return &LoginResponse{}
}

func NewAddCardResponse() *AddCardResponse {
	return &AddCardResponse{}
}
func NewListCardResponse() []ListCardResponse {
	return []ListCardResponse{}
}
func NewListCardControllResponse() []ListCardControlResponse {
	return []ListCardControlResponse{}
}
func NewFetchTransactionResponse() *FetchTransactionResponse {
	return &FetchTransactionResponse{}
}
func NewEditTransactionResponse() *EditTransactionResponse {
	return &EditTransactionResponse{}
}
func NewCardBlockResponse() *CardBlockResponse {
	return &CardBlockResponse{}
}

type DebitCardBlockStatusResponse struct {
	BlockStatus              string `json:"domestic_block_status"`
	InternationalBlockStatus string `json:"international_block_status"`
	IsPermanentlyBlocked     string `json:"is_permanently_blocked"`
}

func NewDebitCardBlockStatusResponse() *DebitCardBlockStatusResponse {
	return &DebitCardBlockStatusResponse{}
}

func (r *DebitCardBlockStatusResponse) Bind(res []ListCardControlResponse, parmanentBlockStatus bool) error {
	r.InternationalBlockStatus = res[0].InternationalBlockStatus
	r.BlockStatus = res[0].BlockStatus

	if !parmanentBlockStatus {
		r.IsPermanentlyBlocked = "0"
	} else {
		r.IsPermanentlyBlocked = "1"
	}

	return nil
}

func (r *DebitCardBlockStatusResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *DebitCardBlockStatusResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}
