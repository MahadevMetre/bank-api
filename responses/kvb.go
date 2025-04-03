package responses

import (
	"fmt"
	"time"

	"bitbucket.org/paydoh/paydoh-commons/pkg/json"
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   string `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type SimVerificationResponse struct {
	ApplicantId  string `json:"ApplicantId"`
	ServiceType  string `json:"ServiceType"`
	ErrorCode    string `json:"ErrorCode"`
	ErrorMessage string `json:"ErrorMessage"`
}

type KycInvokeResponse struct {
	Status       string `json:"status"`
	Data         string `json:"data,omitempty"`
	Message      string `json:"message"`
	ResponseCode string `json:"responseCode"`
}

type DemographicResponse struct {
	ApplicantId   string `json:"ApplicantId"`
	ErrorCode     string `json:"ErrorCode"`
	ErrorMessage  string `json:"ErrorMessage"`
	FirstName     string `json:"FirstName"`
	LastName      string `json:"LastName"`
	MiddleName    string `json:"MiddleName"`
	PanNumber     string `json:"PanNumber"`
	Root          Root   `json:"Root"`
	TxnIdentifier string `json:"TxnIdentifier"`
}
type Poa struct {
	Co    string `json:"Co"`
	Dist  string `json:"Dist"`
	House string `json:"House"`
	State string `json:"State"`
}
type Poi struct {
	Dob    string `json:"Dob"`
	Email  string `json:"Email"`
	Gender string `json:"Gender"`
	Name   string `json:"Name"`
	Phone  string `json:"Phone"`
}
type UIDData struct {
	ErrorCode string `json:"ErrorCode"`
	Pht       string `json:"Pht"`
	Poa       Poa    `json:"Poa"`
	Poi       Poi    `json:"Poi"`
	UID       string `json:"UID"`
}
type Root struct {
	GeneratedKeyForKycResponse string  `json:"GeneratedKeyForKycResponse"`
	IaskRefID                  string  `json:"IaskRefID"`
	Landmark                   string  `json:"Landmark"`
	Locality                   string  `json:"Locality"`
	Pincode                    string  `json:"Pincode"`
	Postoffice                 string  `json:"Postoffice"`
	RespCode                   string  `json:"RespCode"`
	RespDesc                   string  `json:"RespDesc"`
	Ret                        string  `json:"Ret"`
	Rrn                        string  `json:"Rrn"`
	Street                     string  `json:"Street"`
	Subdistrict                string  `json:"Subdistrict"`
	Txn                        string  `json:"Txn"`
	UIDData                    UIDData `json:"UIDData"`
	UidaiAuthCode              string  `json:"UidaiAuthCode"`
	Vtc                        string  `json:"Vtc"`
}

type ImmediateCreateBankResponse struct {
	ErrorCode     string `json:"ErrorCode"`
	ErrorMessage  string `json:"ErrorMessage"`
	ApplicationID string `json:"ApplicationId"`
	ServiceName   string `json:"serviceName"`
}

type CallbackCreateBankResponse struct {
	SourcedBy     string      `json:"SourcedBy"`
	ProductType   string      `json:"ProductType"`
	ServiceName   string      `json:"serviceName"`
	ApplicationID string      `json:"ApplicationId"`
	CbsStatus     []CbsStatus `json:"CbsStatus"`
}
type CbsStatus struct {
	Status      string `json:"Status"`
	SuccErrCode string `json:"succErrCode"`
	AccountNo   int64  `json:"AccountNo"`
	ApplicantId string `json:"ApplicantId"`
	CustomerID  string `json:"CustomerId"`
	Message     string `json:"message"`
}

type OtpGenerationResponse struct {
	TxnIdentifier string `json:"TxnIdentifier"`
	AccountNo     string `json:"AccountNo"`
	ApplicantId   string `json:"ApplicantId"`
	ErrorCode     string `json:"ErrorCode"`
	ErrorMessage  string `json:"ErrorMessage"`
}

type OtpAuthenticationResponse struct {
	TxnIdentifier   string `json:"TxnIdentifier"`
	NomApplId       string `json:"NomApplId"`
	AccountNo       string `json:"AccountNo"`
	ApplicantId     string `json:"ApplicantId"`
	NomUpdateDtTime string `json:"NomUpdateDtTime"`
	NomCBSStatus    string `json:"NomCBSStatus"`
	ErrorCode       string `json:"ErrorCode"`
	ErrorMessage    string `json:"ErrorMessage"`
	TxnStatus       string `json:"TxnStatus"`
}

type TokenErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type IfscDataResponse struct {
	BanksIFSCDtls []BanksIFSCDtl `json:"Banks_IFSC_Dtls"`
	TxnIdentifier string         `json:"TxnIdentifier"`
	ErrorCode     string         `json:"ErrorCode"`
	ErrorMessage  string         `json:"ErrorMessage"`
}

type BanksIFSCDtl struct {
	BankName      string        `json:"BankName"`
	IFSCCode      string        `json:"IFSCCode"`
	BranchName    string        `json:"BranchName"`
	BranchCity    string        `json:"BranchCity"`
	BranchState   string        `json:"BranchState"`
	BranchCountry BranchCountry `json:"BranchCountry"`
	PaymentMode   PaymentMode   `json:"PaymentMode"`
}

type QuickTransferBeneficiaryAdditionResponse struct {
	ApplicantId     string    `json:"ApplicantId"`
	TxnIdentifier   string    `json:"TxnIdentifier"`
	ErrorCode       string    `json:"ErrorCode"`
	ErrorMessage    string    `json:"ErrorMessage"`
	AccountNo       string    `json:"AccountNo,omitempty"`
	BenfId          string    `json:"BenfId,omitempty"`
	TxnStatus       string    `json:"TxnStatus,omitempty"`
	ActivatedDtTime time.Time `json:"ActivatedDtTime,omitempty"`
}

type BranchCountry string

const (
	India BranchCountry = "INDIA"
)

type PaymentMode string

const (
	Imps PaymentMode = "IMPS"
	Neft PaymentMode = "NEFT"
	Rtgs PaymentMode = "RTGS"
	Ift  PaymentMode = "IFT"
)

type FetchBeneficiaryResponse struct {
	TxnIdentifier      string               `json:"TxnIdentifier"`
	BeneficiaryDetails []BeneficiaryDetails `json:"Beneficiary_Dtls,omitempty"`
	AccountNo          string               `json:"AccountNo"`
	ApplicantId        string               `json:"ApplicantId"`
	ErrorCode          string               `json:"ErrorCode"`
	ErrorMessage       string               `json:"ErrorMessage"`
}

type BeneficiaryDetails struct {
	BenfMob          string `json:"BenfMob"`
	BenfName         string `json:"BenfName"`
	BenfID           string `json:"BenfId"`
	BenfAcctNo       string `json:"BenfAcctNo"`
	BenfIFSC         string `json:"BenfIFSC"`
	BenfAcctType     string `json:"BenfAcctType"`
	PaymentMode      string `json:"PaymentMode"`
	BenfStatus       string `json:"BenfStatus"`
	BenfActivateTime string `json:"BenfActivateTime"`
}

type BeneficiarySubmissionResponse struct {
	ApplicantId     string `json:"ApplicantId"`
	TxnIdentifier   string `json:"TxnIdentifier"`
	AccountNo       string `json:"AccountNo"`
	TxnStatus       string `json:"Txn_Status"`
	ErrorCode       string `json:"ErrorCode"`
	ErrorMessage    string `json:"ErrorMessage"`
	ActivatedDtTime string `json:"ActivatedDtTime"`
}

type BeneficiaryOTPValidationResponse struct {
	ApplicantId      string `json:"ApplicantId"`
	TxnIdentifier    string `json:"TxnIdentifier"`
	BenfId           string `json:"BenfId"`
	AccountNo        string `json:"AccountNo"`
	TxnStatus        string `json:"Txn_Status"`
	ErrorCode        string `json:"ErrorCode"`
	ErrorMessage     string `json:"ErrorMessage"`
	ActivationDtTime string `json:"ActivatedDtTime"`
}

type PaymentSubmissionResponse struct {
	ApplicantId   string `json:"ApplicantId"`
	AccountNo     string `json:"AccountNo"`
	TxnStatus     string `json:"Txn_Status"`
	ErrorCode     string `json:"ErrorCode"`
	ErrorMessage  string `json:"ErrorMessage"`
	TxnIdentifier string `json:"TxnIdentifier"`
	TxnRefNo      string `json:"TxnRefNo"`
}

type PaymentSubmissionOtpResponse struct {
	ApplicantId   string    `json:"ApplicantId"`
	AccountNo     string    `json:"AccountNo"`
	TxnStatus     string    `json:"TxnStatus"`
	ErrorCode     string    `json:"ErrorCode"`
	ErrorMessage  string    `json:"ErrorMessage"`
	TxnIdentifier string    `json:"TxnIdentifier"`
	TxnRefNo      string    `json:"TxnRefNo"`
	TxnDateTime   time.Time `json:"TxnDateTime,omitempty"`
}

type ResponseHeader struct {
	ErrorCode        string   `json:"ErrorCode"`
	ErrorMessage     string   `json:"ErrorMessage"`
	HttpStatusCode   string   `json:"HttpStatusCode"`
	ValidationErrors []string `json:"ValidationErrors"`
}

type PlainResponse struct {
	TxnIdentifier string `json:"TxnIdentifier"`
	ApplicantId   string `json:"ApplicantId"`
}

type ConsentResponse struct {
	ResponseHeader ResponseHeader `json:"ResponseHeader"`
	PlainResponse  PlainResponse  `json:"PlainResponse"`
}

type ConsentResponseV2 struct {
	ErrorCode      string             `json:"ErrorCode"`
	ErrorMessage   string             `json:"ErrorMessage"`
	TxnIdentifier  string             `json:"TxnIdentifier"`
	ApplicantId    string             `json:"ApplicantId"`
	ConsentDetails []ConsentDetailsV2 `json:"ConsentDetails"`
}

type ConsentDetailsV2 struct {
	ConsentType     string `json:"ConsentType"`
	ConsentProvided string `json:"ConsentProvided"`
	ConsentTime     string `json:"ConsentTime"`
}

type StatementResponse struct {
	ErrorCode      string         `json:"ErrorCode"`
	ErrorMessage   string         `json:"ErrorMessage"`
	AccountDetails AccountDetails `json:"AccountDetails,omitempty"`
}

type TransactionDetail struct {
	POSTING_DATE       string `json:"POSTING_DATE"`
	REF_NO             string `json:"REF_NO"`
	RUN_BAL            string `json:"RUN_BAL"`
	VALUE_DATE         string `json:"VALUE_DATE"`
	AMOUNT_DEBIT       string `json:"AMOUNT_DEBIT"`
	TRANSACTION_BRANCH string `json:"TRANSACTION_BRANCH"`
	AMOUNT_CREDIT      string `json:"AMOUNT_CREDIT"`
	NARRATION          string `json:"NARRATION"`
}

type AccountDetails struct {
	TxnDetails    []TransactionDetail `json:"TxnDetails"`
	AccountNumber string              `json:"AccountNumber"`
	TxnIdentifier string              `json:"TxnIdentifier"`
	ApplicantId   string              `json:"ApplicantId"`
}

type FetchNomineeResponse struct {
	NomDOB              string `json:"NomDOB"`
	NomCity             string `json:"NomCity"`
	GuardianCountry     string `json:"GuardianCountry"`
	NomReqType          string `json:"NomReqType"`
	GuardianZipcode     string `json:"GuardianZipcode"`
	GuardianCity        string `json:"GuardianCity"`
	GuardianState       string `json:"GuardianState"`
	ApplicantID         string `json:"ApplicantId"`
	NomUpdateDtTime     string `json:"NomUpdateDtTime"`
	NomZipcode          string `json:"NomZipcode"`
	NomName             string `json:"NomName"`
	AccountNo           string `json:"AccountNo"`
	NomCBSStatus        string `json:"NomCBSStatus"`
	NomAppID            string `json:"NomAppId"`
	GuardianName        string `json:"GuardianName"`
	GuardianNomRelation string `json:"GuardianNomRelation"`
	GuardianAddressL3   string `json:"GuardianAddressL3"`
	GuardianAddressL1   string `json:"GuardianAddressL1"`
	GuardianAddressL2   string `json:"GuardianAddressL2"`
	NomAddressL1        string `json:"NomAddressL1"`
	NomAddressL3        string `json:"NomAddressL3"`
	NomAddressL2        string `json:"NomAddressL2"`
	NomRelation         string `json:"NomRelation"`
	TxnIdentifier       string `json:"TxnIdentifier"`
	NomState            string `json:"NomState"`
	NomCountry          string `json:"NomCountry"`
	ErrorCode           string `json:"ErrorCode"`
	ErrorMessage        string `json:"ErrorMessage"`
	TxnStatus           string `json:"TxnStatus"`
}

func NewTokenResponse() *TokenResponse {
	return &TokenResponse{}
}

func NewSimVerificationResponse() *SimVerificationResponse {
	return &SimVerificationResponse{}
}

func NewKycInvokeResponse() *KycInvokeResponse {
	return &KycInvokeResponse{}
}

func NewDemographicResponse() *DemographicResponse {
	return &DemographicResponse{}
}

func NewImmediateCreateBankResponse() *ImmediateCreateBankResponse {
	return &ImmediateCreateBankResponse{}
}

func NewCallbackCreatebankResponse() *CallbackCreateBankResponse {
	return &CallbackCreateBankResponse{}
}

func NewTokenErrorResponse() *TokenErrorResponse {
	return &TokenErrorResponse{}
}

func NewOtpGenartionResponse() *OtpGenerationResponse {
	return &OtpGenerationResponse{}
}

func NewOtpAuthenticationResponse() *OtpAuthenticationResponse {
	return &OtpAuthenticationResponse{}
}

func NewIfscDataResponse() *IfscDataResponse {
	return &IfscDataResponse{}
}

func NewBeneficiaryResponse() *FetchBeneficiaryResponse {
	return &FetchBeneficiaryResponse{}
}

func NewBeneficiarySubmissionResponse() *BeneficiarySubmissionResponse {
	return &BeneficiarySubmissionResponse{}
}

func NewBeneficiaryOTPValidationResponse() *BeneficiaryOTPValidationResponse {
	return &BeneficiaryOTPValidationResponse{}
}

func NewPaymentSubmissionResponse() *PaymentSubmissionResponse {
	return &PaymentSubmissionResponse{}
}

func NewPaymentSubmissionOTPResponse() *PaymentSubmissionOtpResponse {
	return &PaymentSubmissionOtpResponse{}
}

func NewConsentResponse() *ConsentResponse {
	return &ConsentResponse{}
}

func NewConsentResponseV2() *ConsentResponseV2 {
	return &ConsentResponseV2{}
}

func NewQuickTransferBeneficiaryResponse() *QuickTransferBeneficiaryAdditionResponse {
	return &QuickTransferBeneficiaryAdditionResponse{}
}

func NewStatementResponse() *StatementResponse {
	return &StatementResponse{}
}

func NewFetchNomineeResponse() *FetchNomineeResponse {
	return &FetchNomineeResponse{}
}

func (r *TokenResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *TokenResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *SimVerificationResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *SimVerificationResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *KycInvokeResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *KycInvokeResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *DemographicResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *DemographicResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *ImmediateCreateBankResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *ImmediateCreateBankResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *CallbackCreateBankResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *CallbackCreateBankResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *TokenErrorResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *TokenErrorResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *OtpGenerationResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OtpGenerationResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *OtpAuthenticationResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *OtpAuthenticationResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *IfscDataResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *IfscDataResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *FetchBeneficiaryResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *FetchBeneficiaryResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *BeneficiarySubmissionResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *BeneficiarySubmissionResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *BeneficiaryOTPValidationResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *BeneficiaryOTPValidationResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *PaymentSubmissionResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *PaymentSubmissionResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *PaymentSubmissionOtpResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *PaymentSubmissionOtpResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *ConsentResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *ConsentResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *ConsentResponseV2) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *PlainResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *PlainResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *ResponseHeader) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *ResponseHeader) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *QuickTransferBeneficiaryAdditionResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *QuickTransferBeneficiaryAdditionResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *StatementResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *StatementResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *FetchNomineeResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *FetchNomineeResponse) Unmarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *FetchNomineeResponse) Unset() {
	r.ErrorCode = ""
	r.ErrorMessage = ""
	r.TxnIdentifier = ""
	r.TxnStatus = ""
}

type RewardsFundsTransferResponse struct {
	ErrorCode    string `json:"ErrorCode"`
	ErrorMessage string `json:"ErrorMessage"`
	ApplicantId  string `json:"ApplicantId"`
	ServiceName  string `json:"serviceName"`
}

func NewFundTransferResponse() *RewardsFundsTransferResponse {
	return &RewardsFundsTransferResponse{}
}

func (r *RewardsFundsTransferResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *RewardsFundsTransferResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type MobileMappingType0ApiResponse struct {
	Response MobileMappingType0 `json:"Response"`
}

type MobileMappingType0 struct {
	ResponseCode     string `json:"ResponseCode"`
	ResponseMessage  string `json:"ResponseMessage"`
	TransID          string `json:"TransID,omitempty"`
	LongCodeMobileNo string `json:"LongCodeMobileNo,omitempty"`
	VMN              string `json:"VMN,omitempty"`
	ClientId         string `json:"Client_id,omitempty"`
}

func NewMobileMappingType0ApiResponse() *MobileMappingType0ApiResponse {
	return &MobileMappingType0ApiResponse{}
}

func (r *MobileMappingType0ApiResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *MobileMappingType0ApiResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type VerifyUserApiResponse struct {
	Response VerifyUserApi `json:"Response"`
}
type VerifyUserApi struct {
	ResponseCode    string `json:"ResponseCode"`
	ResponseMessage string `json:"ResponseMessage"`
}

func NewVerifyUserApiResponse() *VerifyUserApiResponse {
	return &VerifyUserApiResponse{}
}

func (r *VerifyUserApiResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *VerifyUserApiResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type MobileMappingType1ApiResponse struct {
	Response MobileMappingType1 `json:"Response"`
}

type MobileMappingType1 struct {
	ResponseCode    string `json:"ResponseCode"`
	ResponseMessage string `json:"ResponseMessage"`
	TransID         string `json:"TransID,omitempty"`
	MobileNo        string `json:"MobileNo,omitempty"`
}

func NewMobileMappingType1ApiResponse() *MobileMappingType1ApiResponse {
	return &MobileMappingType1ApiResponse{}
}

func (r *MobileMappingType1ApiResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *MobileMappingType1ApiResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type LcValidatorApiResponse struct {
	Response LcvResponse `json:"Response"`
}

type LcvResponse struct {
	ResponseCode    string `json:"ResponseCode"`
	ResponseMessage string `json:"ResponseMessage"`
	LoginRefID      string `json:"LoginRefID,omitempty"`
	ServerID        string `json:"ServerID,omitempty"`
	TransID         string `json:"TransID,omitempty"`
	MobileNo        string `json:"MobileNo,omitempty"`
}

func NewLcValidatorApiResponse() *LcValidatorApiResponse {
	return &LcValidatorApiResponse{}
}

func (r *LcValidatorApiResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *LcValidatorApiResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type ProfileCreationApiResponse struct {
	Response ProfileCreation `json:"Response"`
}

type ProfileCreation struct {
	ResponseCode    string `json:"ResponseCode"`
	ResponseMessage string `json:"ResponseMessage"`
}

func NewProfileCreationApiResponse() *ProfileCreationApiResponse {
	return &ProfileCreationApiResponse{}
}

func (r *ProfileCreationApiResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *ProfileCreationApiResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type NewRequestKeyListApiResponse struct {
	Response KeyListResponseData `json:"Response"`
}

type KeyListResponseData struct {
	ResponseCode    string               `json:"ResponseCode"`
	ResponseMessage string               `json:"ResponseMessage"`
	MessageID       string               `json:"MessageID,omitempty"`
	Response        ExistingUserResponse `json:"Response,omitempty"`
}

func NewRequestKeyListResponse() *NewRequestKeyListApiResponse {
	return &NewRequestKeyListApiResponse{}
}

func (r *NewRequestKeyListApiResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *NewRequestKeyListApiResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type RequestListKeysResponse struct {
	Response string `json:"Response"`
}

type ListKeyResponse struct {
	RequestListKeys []KeyResponse `json:"RequestListKeys"`
}

type KeyResponse struct {
	ResponseCode    string `json:"ResponseCode"`
	ResponseMessage string `json:"ResponseMessage"`
	MessageID       string `json:"MessageID,omitempty"`
	// Response        string `json:"Response,omitempty"`
}

func NewRequestListKeysResponse() *RequestListKeysResponse {
	return &RequestListKeysResponse{}
}

func (r *RequestListKeysResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *RequestListKeysResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type ReMappingApiResponse struct {
	Response ReMappingResponse `json:"Response"`
}

type ReMappingResponse struct {
	ResponseCode    string `json:"ResponseCode"`
	ResponseMessage string `json:"ResponseMessage"`
	LoginRefID      string `json:"LoginRefID,omitempty"`
	ServerID        string `json:"ServerID,omitempty"`
}

func NewReMappingApiResponse() *ReMappingApiResponse {
	return &ReMappingApiResponse{}
}

func (r *ReMappingApiResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *ReMappingApiResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type BankError struct {
	ErrorCode    string
	ErrorMessage string
	HttpStatus   string
}

type XmlRequestListKeyApiResponse struct {
	Response string `json:"Response"`
	BankError
}

type RequestListKey struct {
	ResponseCode    string `json:"ResponseCode"`
	ResponseMessage string `json:"ResponseMessage"`
	MessageID       string `json:"MessageID,omitempty"`
	Response        string `json:"Response,omitempty"`
}

func NewXmlRequestListKeyApiResponse() *XmlRequestListKeyApiResponse {
	return &XmlRequestListKeyApiResponse{}
}

func (r *XmlRequestListKeyApiResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *XmlRequestListKeyApiResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type ExistingUserReqListApiResponse struct {
	Response ExistreqlistResponse `json:"Response"`
}

type ExistreqlistResponse struct {
	ResponseCode    string               `json:"ResponseCode"`
	ResponseMessage string               `json:"ResponseMessage"`
	MessageID       string               `json:"MessageID,omitempty"`
	Response        ExistingUserResponse `json:"Response,omitempty"`
}

type ExistingUserResponse struct {
	XML             ExistingUserXML             `json:"?xml"`
	Ns2RespListKeys ExistingUserNs2RespListKeys `json:"ns2:RespListKeys"`
}

type ExistingUserNs2RespListKeys struct {
	XmlnsNs2 string           `json:"@xmlns:ns2"`
	XmlnsNs3 string           `json:"@xmlns:ns3"`
	Head     ExistingUserHead `json:"Head"`
	Resp     ExistingUserResp `json:"Resp"`
	Txn      ExistingUserTxn  `json:"Txn"`
	KeyList  KeyList          `json:"keyList"`
}

type KeyList struct {
	Key Key `json:"key"`
}

type Key struct {
	Code     string          `json:"@code"`
	Ki       string          `json:"@ki"`
	Owner    string          `json:"@owner"`
	Type     string          `json:"@type"`
	KeyValue KeyListKeyValue `json:"keyValue"`
}

type KeyListKeyValue struct {
	XmlnsXs  string `json:"@xmlns:xs"`
	XmlnsXsi string `json:"@xmlns:xsi"`
	XsiType  string `json:"@xsi:type"`
	Text     string `json:"#text"`
}

type Response struct {
	KeyList KeyList `json:"keyList"`
}

type ExistingUserXML struct {
	Version  string `json:"@version"`
	Encoding string `json:"@encoding"`
}

type ExistingUserHead struct {
	MsgID string `json:"@msgId"`
	OrgID string `json:"@orgId"`
	Ts    string `json:"@ts"`
	Ver   string `json:"@ver"`
}

type ExistingUserResp struct {
	ReqMsgID string `json:"@reqMsgId"`
	Result   string `json:"@result"`
}

type ExistingUserTxn struct {
	ID     string `json:"@id"`
	Note   string `json:"@note"`
	RefID  string `json:"@refId"`
	RefURL string `json:"@refUrl"`
	Ts     string `json:"@ts"`
	Type   string `json:"@type"`
}

func NewExistingUserReqListApiResponse() *ExistingUserReqListApiResponse {
	return &ExistingUserReqListApiResponse{}
}

func (r *ExistingUserReqListApiResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *ExistingUserReqListApiResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type CreateUpiIdRequestListAccountApiResponse struct {
	Response ReqListAccount `json:"Response"`
}

type ReqListAccount struct {
	ResponseCode    string         `json:"ResponseCode"`
	ResponseMessage string         `json:"ResponseMessage"`
	MessageID       string         `json:"MessageID,omitempty"`
	Response        NestedResponse `json:"Response,omitempty"`
}

type NestedResponse struct {
	XML                ReqListXML                `json:"?xml"`
	Ns2RespListAccount ReqListNs2RespListAccount `json:"ns2:RespListAccount"`
}

type ReqListXML struct {
	Version  string `json:"@version"`
	Encoding string `json:"@encoding"`
}

type ReqListNs2RespListAccount struct {
	XmlnsNs2    string             `json:"@xmlns:ns2"`
	XmlnsNs3    string             `json:"@xmlns:ns3"`
	Head        ReqListHead        `json:"Head"`
	Resp        ReqListResp        `json:"Resp"`
	Txn         ReqListTxn         `json:"Txn"`
	AccountList ReqListAccountList `json:"AccountList"`
	// Signature   ReqListSignature   `json:" Signature"`
}

type ReqListHead struct {
	MsgID string `json:"@msgId"`
	OrgID string `json:"@orgId"`
	Ts    string `json:"@ts"`
	Ver   string `json:"@ver"`
}

type ReqListResp struct {
	ReqMsgID string `json:"@reqMsgId"`
	Result   string `json:"@result"`
}

type ReqListTxn struct {
	ID     string `json:"@id"`
	Note   string `json:"@note"`
	RefID  string `json:"@refId"`
	RefURL string `json:"@refUrl"`
	Ts     string `json:"@ts"`
	Type   string `json:"@type"`
}

type AccountListResp []AccountList

// Custom unmarshalling for AccountList to handle both single Account and a list of Accounts
func (a *AccountListResp) UnmarshalJSON(data []byte) error {
	var accountSlice []AccountList
	if err := json.Unmarshal(data, &accountSlice); err == nil {
		*a = accountSlice
		return nil
	}

	var account AccountList
	if err := json.Unmarshal(data, &account); err == nil {
		*a = append(*a, account)
		return nil
	}

	return fmt.Errorf("failed to unmarshal Account, expected either a single Account or a slice of Accounts")
}

type ReqListAccountList struct {
	Account AccountListResp `json:"Account"`
}

type AccountList struct {
	AadhaarNumber   string `json:"@aadhaarNo"`
	AccRefNumber    string `json:"@accRefNumber"`
	AccType         string `json:"@accType"`
	Aeba            string `json:"@aeba"`
	Ifsc            string `json:"@ifsc"`
	MaskedAccnumber string `json:"@maskedAccnumber"`
	Mbeba           string `json:"@mbeba"`
	// Mmid            string              `json:"@mmid"`
	Name         string                `json:"@name"`
	CredsAllowed []ReqListCredsAllowed `json:"CredsAllowed"`
}

type ReqListCredsAllowed struct {
	DLength string `json:"@dLength"`
	Type    string `json:"@type"`
	DType   string `json:"@dType"`
	SubType string `json:"@subType"`
}

type ReqListSignature struct {
	Xmlns          string            `json:"@xmlns"`
	SignedInfo     ReqListSignedInfo `json:"SignedInfo"`
	SignatureValue string            `json:"SignatureValue"`
	KeyInfo        ReqListKeyInfo    `json:"KeyInfo"`
}

type ReqListCanonicalizationMethod struct {
	Algorithm string `json:"@Algorithm"`
}

type ReqListSignatureMethod struct {
	Algorithm string `json:"@Algorithm"`
}

type ReqListReference struct {
	URI          string              `json:"@URI"`
	Transforms   ReqListTransforms   `json:"Transforms"`
	DigestMethod ReqListDigestMethod `json:"DigestMethod"`
	DigestValue  string              `json:"DigestValue"`
}

type ReqListTransform struct {
	Algorithm string `json:"@Algorithm"`
}

type ReqListTransforms struct {
	Transform ReqListTransform `json:"Transform"`
}

type ReqListDigestMethod struct {
	Algorithm string `json:"@Algorithm"`
}

type ReqListSignedInfo struct {
	CanonicalizationMethod ReqListCanonicalizationMethod `json:"Canonicalization Method"`
	SignatureMethod        ReqListSignatureMethod        `json:"SignatureMethod"`
	Reference              ReqListReference              `json:"Reference"`
}

type ReqListKeyInfo struct {
	KeyValue ReqListKeyValue `json:"KeyValue"`
}

type ReqListKeyValue struct {
	RSAKeyValue ReqListRSAKeyValue `json:"RSAKeyValue"`
}

type ReqListRSAKeyValue struct {
	Modulus  string `json:"Modulus"`
	Exponent string `json:"Exponent"`
}

func NewCreateUpiIdRequestListAccountApiResponse() *CreateUpiIdRequestListAccountApiResponse {
	return &CreateUpiIdRequestListAccountApiResponse{}
}

func (r *CreateUpiIdRequestListAccountApiResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *CreateUpiIdRequestListAccountApiResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type PspAvailabilityApiResponse struct {
	Response PspAvailabilityResponse `json:"Response"`
}

type PspAvailabilityResponse struct {
	ResponseCode    string `json:"ResponseCode"`
	ResponseMessage string `json:"ResponseMessage"`
}

func NewPspAvailabilityApiResponse() *PspAvailabilityApiResponse {
	return &PspAvailabilityApiResponse{}
}

func (r *PspAvailabilityApiResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *PspAvailabilityApiResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type AddBankApiResponse struct {
	Response AddBank `json:"Response"`
}

type AddBank struct {
	ResponseCode    string `json:"ResponseCode"`
	ResponseMessage string `json:"ResponseMessage"`
	UpiId           string `json:"upi_id"`
}

func NewAddBankApiResponse() *AddBankApiResponse {
	return &AddBankApiResponse{}
}

func (r *AddBankApiResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *AddBankApiResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type ReqBalEnqApiResponse struct {
	Response ReqBalEnqResponse `json:"Response"`
}

type ReqBalEnqResponse struct {
	ResponseCode    string     `json:"ResponseCode"`
	ResponseMessage string     `json:"ResponseMessage"`
	MessageID       string     `json:"MessageID,omitempty"`
	Response        MyResponse `json:"Response,omitempty"`
}

type MyResponse struct {
	XML           MyXML           `json:"?xml"`
	Ns2RespBalEnq MyNs2RespBalEnq `json:"ns2:RespBalEnq"`
}

type MyNs2RespBalEnq struct {
	XmlnsNs2  string      `json:"@xmlns:ns2"`
	XmlnsNs3  string      `json:"@xmlns:ns3"`
	Head      MyHead      `json:"Head"`
	Resp      MyResp      `json:"Resp"`
	Txn       MyTxn       `json:"Txn"`
	Payer     MyPayer     `json:"Payer"`
	Signature MySignature `json:"Signature"`
}

type MyXML struct {
	Version  string `json:"@version"`
	Encoding string `json:"@encoding"`
}

type MyHead struct {
	MsgID string `json:"@msgId"`
	OrgID string `json:"@orgId"`
	Ts    string `json:"@ts"`
	Ver   string `json:"@ver"`
}

type MyResp struct {
	ReqMsgID string `json:"@reqMsgId"`
	Result   string `json:"@result"`
}

type MyScore struct {
	Provider string `json:"@provider"`
	Type     string `json:"@type"`
	Value    string `json:"@value"`
}

type MyRiskScores struct {
	Score MyScore `json:"Score"`
}

type MyTxn struct {
	ID         string       `json:"@id"`
	Note       string       `json:"@note"`
	RefID      string       `json:"@refId"`
	RefURL     string       `json:"@refUrl"`
	Ts         string       `json:"@ts"`
	Type       string       `json:"@type"`
	RiskScores MyRiskScores `json:"RiskScores"`
}

type MyData struct {
	Code string `json:"@code"`
	Ki   string `json:"@ki"`
	Text string `json:"#text"`
}

type MyBal struct {
	Data MyData `json:"Data"`
}

type MyPayer struct {
	Addr   string `json:"@addr"`
	Code   string `json:"@code"`
	Name   string `json:"@name"`
	SeqNum string `json:"@seqNum"`
	Type   string `json:"@type"`
	Bal    MyBal  `json:"Bal"`
}

type MyCanonicalizationMethod struct {
	Algorithm string `json:"@Algorithm"`
}

type MySignatureMethod struct {
	Algorithm string `json:"@Algorithm"`
}

type MyTransform struct {
	Algorithm string `json:"@Algorithm"`
}

type MyTransforms struct {
	Transform MyTransform `json:"Transform"`
}

type MyDigestMethod struct {
	Algorithm string `json:"@Algorithm"`
}

type MyReference struct {
	URI          string         `json:"@URI"`
	Transforms   MyTransforms   `json:"Transforms"`
	DigestMethod MyDigestMethod `json:"DigestMethod"`
	DigestValue  string         `json:"DigestValue"`
}

type MySignedInfo struct {
	CanonicalizationMethod MyCanonicalizationMethod `json:"CanonicalizationMethod"`
	SignatureMethod        MySignatureMethod        `json:"SignatureMethod"`
	Reference              MyReference              `json:"Reference"`
}

type MyRSAKeyValue struct {
	Modulus  string `json:"Modulus"`
	Exponent string `json:"Exponent"`
}

type MyKeyValue struct {
	RSAKeyValue MyRSAKeyValue `json:"RSAKeyValue"`
}

type MyKeyInfo struct {
	KeyValue MyKeyValue `json:"KeyValue"`
}

type MySignature struct {
	Xmlns          string       `json:"@xmlns"`
	SignedInfo     MySignedInfo `json:"SignedInfo"`
	SignatureValue string       `json:"SignatureValue"`
	KeyInfo        MyKeyInfo    `json:"KeyInfo"`
}

func NewReqBalEnqApiResponse() *ReqBalEnqApiResponse {
	return &ReqBalEnqApiResponse{}
}

func (r *ReqBalEnqApiResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *ReqBalEnqApiResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type AadharRequestListAccountApiResponse struct {
	Response RequestListAccountResponse `json:"Response"`
}

type RequestListAccountResponse struct {
	ResponseCode    string               `json:"ResponseCode"`
	ResponseMessage string               `json:"ResponseMessage"`
	MessageID       string               `json:"MessageID,omitempty"`
	Response        AadharReqListAccount `json:"Response,omitempty"`
}

type AadharReqListAccount struct {
	XML             XMLInfo         `json:"?xml"`
	RespListAccount RespListAccount `json:"ns2:RespListAccount"`
}

type XMLInfo struct {
	Version  string `json:"@version"`
	Encoding string `json:"@encoding"`
}

type RespListAccount struct {
	XMLNSNs2    string       `json:"@xmlns:ns2"`
	XMLNSNs3    string       `json:"@xmlns:ns3"`
	Head        Head         `json:"Head"`
	Resp        Resp         `json:"Resp"`
	Txn         Txn          `json:"Txn"`
	AccountList AccountLists `json:"AccountList"`
}

type Head struct {
	MsgID    string `json:"@msgId"`
	OrgID    string `json:"@orgId"`
	ProdType string `json:"@prodType"`
	Ts       string `json:"@ts"`
	Ver      string `json:"@ver"`
}

type Resp struct {
	Ac       string `json:"@ac"`
	Lk       string `json:"@lk"`
	ReqMsgID string `json:"@reqMsgId"`
	Result   string `json:"@result"`
	Sa       string `json:"@sa"`
}

type Txn struct {
	ID     string `json:"@id"`
	Note   string `json:"@note"`
	RefID  string `json:"@refId"`
	RefUrl string `json:"@refUrl"`
	Ts     string `json:"@ts"`
	Type   string `json:"@type"`
}

type AccountLists struct {
	Account AccountListResp `json:"Account"`
}

type Account struct {
	AadhaarNo       string         `json:"@aadhaarNo"`
	AccRefNumber    string         `json:"@accRefNumber"`
	AccType         string         `json:"@accType"`
	Aeba            string         `json:"@aeba"`
	Ifsc            string         `json:"@ifsc"`
	MaskedAccNumber string         `json:"@maskedAccnumber"`
	Mbeba           string         `json:"@mbeba"`
	Name            string         `json:"@name"`
	CredsAllowed    []CredsAllowed `json:"CredsAllowed"`
}

type CredsAllowed struct {
	DLength string `json:"@dLength"`
	DType   string `json:"@dType"`
	SubType string `json:"@subType"`
	Type    string `json:"@type"`
}

func NewAadharRequestListAccountApiResponse() *AadharRequestListAccountApiResponse {
	return &AadharRequestListAccountApiResponse{}
}

func (r *AadharRequestListAccountApiResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *AadharRequestListAccountApiResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type ReqOtpApiResponse struct {
	Response ReqOtp `json:"Response"`
}

type ReqOtp struct {
	ResponseCode    string    `json:"ResponseCode"`
	ResponseMessage string    `json:"ResponseMessage"`
	MessageID       string    `json:"MessageID,omitempty"`
	Response        GetReqOtp `json:"Response,omitempty"`
}

type GetReqOtp struct {
	XML  ReqOtpXML `json:"?xml"`
	Resp RespOtp   `json:"ns2:RespOtp"`
}

type ReqOtpXML struct {
	Version  string `json:"@version"`
	Encoding string `json:"@encoding"`
}

type ReqOtpHead struct {
	MsgID string `json:"@msgId"`
	OrgID string `json:"@orgId"`
	Ts    string `json:"@ts"`
	Ver   string `json:"@ver"`
}

type RespOtp struct {
	XmlnsNs2 string     `json:"@xmlns:ns2"`
	XmlnsNs3 string     `json:"@xmlns:ns3"`
	Head     ReqOtpHead `json:"Head"`
	Resp     ReqOtpResp `json:"Resp"`
	Txn      ReqOtpTxn  `json:"Txn"`
}

type ReqOtpResp struct {
	ReqMsgID string `json:"@reqMsgId"`
	Result   string `json:"@result"`
}

type ReqOtpTxn struct {
	ID      string `json:"@id"`
	Note    string `json:"@note"`
	RefID   string `json:"@refId"`
	RefURL  string `json:"@refUrl"`
	SubType string `json:"@subType"`
	Ts      string `json:"@ts"`
	Type    string `json:"@type"`
}

func NewReqOtpApiResponse() *ReqOtpApiResponse {
	return &ReqOtpApiResponse{}
}

func (r *ReqOtpApiResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *ReqOtpApiResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type ReqRegMobApiResponse struct {
	Response ReqRegMob `json:"Response"`
}

type ReqRegMob struct {
	ResponseCode    string `json:"ResponseCode"`
	ResponseMessage string `json:"ResponseMessage"`
	MessageID       string `json:"MessageID,omitempty"`
	// Response        ResponseData `json:"Response,omitempty"`
}

type ResponseData struct {
	XML           ReqRegXML           `json:"?xml"`
	Ns2RespRegMob ReqRegNs2RespRegMob `json:"ns2:RespRegMob"`
}

type ReqRegXML struct {
	Version  string `json:"@version"`
	Encoding string `json:"@encoding"`
}

type ReqRegHead struct {
	MsgID string `json:"@msgId"`
	OrgID string `json:"@orgId"`
	Ts    string `json:"@ts"`
	Ver   string `json:"@ver"`
}

type ReqRegResp struct {
	ReqMsgID string `json:"@reqMsgId"`
	Result   string `json:"@result"`
}

type ReqRegTxn struct {
	ID     string `json:"@id"`
	Note   string `json:"@note"`
	RefID  string `json:"@refId"`
	RefURL string `json:"@refUrl"`
	Ts     string `json:"@ts"`
	Type   string `json:"@type"`
}

type ReqRegNs2RespRegMob struct {
	XmlnsNs2 string     `json:"@xmlns:ns2"`
	XmlnsNs3 string     `json:"@xmlns:ns3"`
	Head     ReqRegHead `json:"Head"`
	Resp     ReqRegResp `json:"Resp"`
	Txn      ReqRegTxn  `json:"Txn"`
}

func NewReqRegMobApiResponse() *ReqRegMobApiResponse {
	return &ReqRegMobApiResponse{}
}

func (r *ReqRegMobApiResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *ReqRegMobApiResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type ReqValAddApiResponse struct {
	Response ReqValAdd `json:"Response"`
}

type ReqValAdd struct {
	ResponseCode    string         `json:"ResponseCode"`
	ResponseMessage string         `json:"ResponseMessage"`
	MessageID       string         `json:"MessageID,omitempty"`
	Response        ReqValResponse `json:"Response,omitempty"`
}

type ReqValResponse struct {
	XML           ReqValXML           `json:"?xml"`
	Ns2RespValAdd ReqValNs2RespValAdd `json:"ns2:RespValAdd"`
}

type ReqValNs2RespValAdd struct {
	XmlnsNs2 string     `json:"@xmlns:ns2"`
	XmlnsNs3 string     `json:"@xmlns:ns3"`
	Head     ReqValHead `json:"Head"`
	Resp     ReqValResp `json:"Resp"`
	Txn      ReqValTxn  `json:"Txn"`
}

type ReqValXML struct {
	Version  string `json:"@version"`
	Encoding string `json:"@encoding"`
}

type ReqValHead struct {
	MsgID    string `json:"@msgId"`
	OrgID    string `json:"@orgId"`
	ProdType string `json:"@prodType"`
	Ts       string `json:"@ts"`
	Ver      string `json:"@ver"`
}

type ReqValResp struct {
	IFSC                   string                 `json:"@IFSC"`
	AccType                string                 `json:"@accType"`
	Addr                   string                 `json:"@addr"`
	Code                   string                 `json:"@code"`
	ErrCode                string                 `json:"@errCode"`
	MaskName               string                 `json:"@maskName"`
	ReqMsgID               string                 `json:"@reqMsgId"`
	Result                 string                 `json:"@result"`
	Type                   string                 `json:"@type"`
	Merchant               ReqValMerchant         `json:"Merchant"`
	ReqValFeatureSupported ReqValFeatureSupported `json:"FeatureSupported"`
}

type ReqValTxn struct {
	ID     string `json:"@id"`
	Note   string `json:"@note"`
	RefID  string `json:"@refId"`
	RefURL string `json:"@refUrl"`
	Ts     string `json:"@ts"`
	Type   string `json:"@type"`
}

type ReqValMerchant struct {
	Identifier ReqValIdentifier `json:"Identifier"`
	Name       ReqValName       `json:"Name"`
	Ownership  ReqValOwnership  `json:"Ownership"`
}

type ReqValFeatureSupported struct {
	Value string `json:"@value"`
}

type ReqValIdentifier struct {
	MerchantGenre  string `json:"@merchantGenre"`
	MerchantType   string `json:"@merchantType"`
	Mid            string `json:"@mid"`
	OnBoardingType string `json:"@onBoardingType"`
	Sid            string `json:"@sid"`
	SubCode        string `json:"@subCode"`
	Tid            string `json:"@tid"`
}

type ReqValName struct {
	Brand string `json:"@brand"`
	Legal string `json:"@legal"`
}

type ReqValOwnership struct {
	Type string `json:"@type"`
}

func NewReqValAddApiResponse() *ReqValAddApiResponse {
	return &ReqValAddApiResponse{}
}

func (r *ReqValAddApiResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *ReqValAddApiResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type ReqPayApiResponse struct {
	Response ReqPayMoney `json:"Response"`
}

type ReqPayMoney struct {
	ResponseCode    string      `json:"ResponseCode"`
	ResponseMessage string      `json:"ResponseMessage"`
	MessageID       string      `json:"MessageID,omitempty"`
	Response        ReqResponse `json:"Response,omitempty"`
}

type ReqResponse struct {
	XML        ReqPayXML        `json:"?xml"`
	Ns2RespPay ReqPayNs2RespPay `json:"ns2:RespPay"`
}

type ReqPayNs2RespPay struct {
	XmlnsNs2 string     `json:"@xmlns:ns2"`
	XmlnsNs3 string     `json:"@xmlns:ns3"`
	Head     ReqPayHead `json:"Head"`
	Txn      ReqPayTxn  `json:"Txn"`
	Resp     ReqPayResp `json:"Resp"`
}

type ReqPayXML struct {
	Version  string `json:"@version"`
	Encoding string `json:"@encoding"`
}

type ReqPayHead struct {
	MsgID string `json:"@msgId"`
	OrgID string `json:"@orgId"`
	Ts    string `json:"@ts"`
	Ver   string `json:"@ver"`
}

type ReqPayScore struct {
	Provider string `json:"@provider"`
	Type     string `json:"@type"`
	Value    string `json:"@value"`
}

type ReqPayRiskScores struct {
	Score ReqPayScore `json:"Score"`
}

type ReqPayTxn struct {
	CustRef        string           `json:"@custRef"`
	ID             string           `json:"@id"`
	InitiationMode string           `json:"@initiationMode"`
	Note           string           `json:"@note"`
	Purpose        string           `json:"@purpose"`
	RefID          string           `json:"@refId"`
	RefURL         string           `json:"@refUrl"`
	Ts             string           `json:"@ts"`
	Type           string           `json:"@type"`
	RiskScores     ReqPayRiskScores `json:"RiskScores"`
}

type ReqPayRef struct {
	IFSC         string `json:"@IFSC"`
	AcNum        string `json:"@acNum"`
	AccType      string `json:"@accType"`
	Addr         string `json:"@addr"`
	ApprovalNum  string `json:"@approvalNum"`
	Code         string `json:"@code"`
	OrgAmount    string `json:"@orgAmount"`
	RegName      string `json:"@regName"`
	RespCode     string `json:"@respCode"`
	SeqNum       string `json:"@seqNum"`
	SettAmount   string `json:"@settAmount"`
	SettCurrency string `json:"@settCurrency"`
	Type         string `json:"@type"`
}

type ReqPayResp struct {
	ReqMsgID string      `json:"@reqMsgId"`
	Result   string      `json:"@result"`
	Ref      []ReqPayRef `json:"Ref"`
}

func NewReqPayApiResponse() *ReqPayApiResponse {
	return &ReqPayApiResponse{}
}

func (r *ReqPayApiResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *ReqPayApiResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type UpiPaymentResponse struct {
	UTR             string `json:"utr"`
	IFSC            string `json:"ifsc"`
	MobileNumber    string `json:"mobile_number"`
	AccountNumber   string `json:"account_number"`
	Remarks         string `json:"remarks"`
	Amount          string `json:"amount"`
	PayeeName       string `json:"payee_name"`
	TransactionTime string `json:"transaction_time"`
}

func NewUpiPaymentResponse() *UpiPaymentResponse {
	return &UpiPaymentResponse{}
}

func (r *UpiPaymentResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *UpiPaymentResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type AccountLinkApiResponse struct {
	Response AccountLink `json:"Response"`
}

type AccountLink struct {
	ResponseCode    string                `json:"ResponseCode"`
	ResponseMessage string                `json:"ResponseMessage"`
	Response        []LinkaccountResponse `json:"Response,omitempty"`
}

type LinkaccountResponse struct {
	Payername     string `json:"PAYERNAME"`
	Payeraddr     string `json:"PAYERADDR"`
	Accounttype   string `json:"ACCOUNTTYPE"`
	AccountIfsc   string `json:"ACCOUNT_IFSC"`
	AccountActype string `json:"ACCOUNT_ACTYPE"`
	AccountAcnum  string `json:"ACCOUNT_ACNUM"`
	MobileMobnum  string `json:"MOBILE_MOBNUM"`
	Regstatus     string `json:"REGSTATUS"`
	Dtype         string `json:"DTYPE"`
	Dlength       string `json:"DLENGTH"`
	Otpdtype      string `json:"OTPDTYPE"`
	Otpdlength    string `json:"OTPDLENGTH"`
	Mbeba         string `json:"mbeba"`
	Aeba          string `json:"aeba"`
	Maskedaccno   string `json:"maskedaccno"`
	Bankname      string `json:"BANKNAME"`
	Atmdtype      string `json:"ATMDTYPE"`
	Atmdlength    string `json:"ATMDLENGTH"`
}

func NewAccountLinkApiResponse() *AccountLinkApiResponse {
	return &AccountLinkApiResponse{}
}

func (r *AccountLinkApiResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *AccountLinkApiResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type UpiMoneyCollectDetailsResponse struct {
	Response CollectDetailsResponse `json:"Response"`
}

type CollectDetailsResponse struct {
	ResponseCode    string           `json:"ResponseCode"`
	ResponseMessage string           `json:"ResponseMessage"`
	Response        []CollectDetails `json:"Response"`
}

type CollectDetails struct {
	TransactionType  string  `json:"TRANSACTIONTYPE,omitempty"`
	TransactionID    string  `json:"TRANSACTIONID,omitempty"`
	TransRefID       string  `json:"TRANSREFID,omitempty"`
	PayeeName        string  `json:"PAYEENAME,omitempty"`
	PayeeAccountNo   string  `json:"PAYEEACCOUNTNO,omitempty"`
	PayeeAmount      string  `json:"PAYEEAMOUNT,omitempty"`
	PayerName        string  `json:"PAYERNAME,omitempty"`
	PayerAccountNo   string  `json:"PAYERACCOUNTNO,omitempty"`
	PayerAmount      string  `json:"PAYERAMOUNT,omitempty"`
	ExpiryDateTime   string  `json:"EXIPRYDATETIME,omitempty"`
	Remarks          string  `json:"REMARKS,omitempty"`
	PayeeAddr        string  `json:"PAYEEADDR,omitempty"`
	PayerAddr        string  `json:"PAYERADDR,omitempty"`
	DType            string  `json:"DTYPE,omitempty"`
	DLength          float64 `json:"DLENGTH,omitempty"`
	VerifiedMerchant string  `json:"VERIFIEDMERCHANT,omitempty"`
	RefURL           *string `json:"REFURL,omitempty"`
}

func NewUpiMoneyCollectDetailsResponse() *UpiMoneyCollectDetailsResponse {
	return &UpiMoneyCollectDetailsResponse{}
}

func (r *UpiMoneyCollectDetailsResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *UpiMoneyCollectDetailsResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type UpiMoneyCollectCountResponse struct {
	Response UpiCollectCount `json:"Response"`
}

type UpiCollectCount struct {
	ResponseCode    string         `json:"ResponseCode"`
	ResponseMessage string         `json:"ResponseMessage"`
	Response        []CollectCount `json:"Response"`
}

type CollectCount struct {
	Transactiontype string `json:"TRANSACTIONTYPE,omitempty"`
	Transid         string `json:"TRANSID,omitempty"`
	Payername       string `json:"PAYERNAME,omitempty"`
	Payeraddr       string `json:"PAYERADDR,omitempty"`
	Payeename       string `json:"PAYEENAME,omitempty"`
	Payeeaddr       string `json:"PAYEEADDR,omitempty"`
	Payeramount     string `json:"PAYERAMOUNT,omitempty"`
	Statusdesc      string `json:"STATUSDESC,omitempty"`
	Transactiondate string `json:"TRANSACTIONDATE,omitempty"`
	Remarks         string `json:"REMARKS,omitempty"`
	Errorcode       string `json:"ERRORCODE,omitempty"`
	Respcode        string `json:"RESPCODE,omitempty"`
	Payrefid        string `json:"PAYREFID,omitempty"`
	Refid           string `json:"refid,omitempty"`
	ErrorDesc       string `json:"ErrorDesc,omitempty"`
}

func NewUpiMoneyCollectCountResponse() *UpiMoneyCollectCountResponse {
	return &UpiMoneyCollectCountResponse{}
}

func (r *UpiMoneyCollectCountResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *UpiMoneyCollectCountResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type UpiMoneyCollectApprovalResponse struct {
	Response UpiMoneyApprovalResponse `json:"Response"`
}

type UpiMoneyApprovalResponse struct {
	ResponseCode    string `json:"ResponseCode"`
	ResponseMessage string `json:"ResponseMessage"`
	RefID           string `json:"RefID,omitempty"`
	PayerAddr       string `json:"PayerAddr,omitempty"`
	PayeeAddr       string `json:"PayeeAddr,omitempty"`
	Amount          string `json:"Amount,omitempty"`
}

func NewUpiMoneyCollectApprovalResponse() *UpiMoneyCollectApprovalResponse {
	return &UpiMoneyCollectApprovalResponse{}
}

func (r *UpiMoneyCollectApprovalResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *UpiMoneyCollectApprovalResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type AccountLinkedResponse struct {
	Response Responses `json:"Response"`
}

// Responses contains the response details
type Responses struct {
	RESPONSECODE    string      `json:"RESPONSECODE"`
	RESPONSEMESSAGE string      `json:"RESPONSESAGESAGE"`
	Response        []Response1 `json:"Response"`
}

// Response1 represents individual response details
type Response1 struct {
	PAYERNAME      string `json:"PAYERNAME"`
	PAYERADDR      string `json:"PAYERADDR"`
	ACCOUNTTYPE    string `json:"ACCOUNTYP"`
	ACCOUNT_IFSC   string `json:"ACCOUNT_IFSC"`
	ACCOUNT_ACTYPE string `json:"ACCOUNT_ACTYPE"`
	ACCOUNT_ACNUM  string `json:"ACCOUNT_ACNUM"`
	DEVICEMOBILE   string `json:"DEVICEMOBILE"`
	MOBILE_MOBNUM  string `json:"MOBILE_MOBNUM"`
	REGSTATUS      string `json:"REGSTATUS"`
	DTYPE          string `json:"DTYPE"`
	DLENGTH        string `json:"DLENGTH"`
	OTPDTYPE       string `json:"OTPDTYPE"`
	OTPDLENGTH     string `json:"OTPLENGTH"`
	Mbeba          string `json:"mbe"`
	Aeba           string `json:"ae"`
	Maskedaccno    string `json:"maskedaccnno"`
	BANKNAME       string `json:"BANKNAME"`
	ATMDTYPE       string `json:"ATM_TYPE"`
	ATMDLENGTH     string `json:"ATM_LINKT"`
}

func NewAccountLinkedResponseResponse() *AccountLinkedResponse {
	return &AccountLinkedResponse{}
}

func (r *AccountLinkedResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *AccountLinkedResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type DebitcardDetailResponse struct {
	ServiceData  ServiceData `json:"serviceData"`
	ApplicantId  string      `json:"ApplicantId"`
	AccountNo    string      `json:"AccountNo"`
	ErrorCode    string      `json:"ErrorCode"`
	ErrorMessage string      `json:"ErrorMessage"`
}

type ServiceData struct {
	CardData []CardData `json:"cardData"`
}

type CardData struct {
	ProxyNumber        string `json:"proxyNumber"`
	EncryptedPAN       string `json:"encryptedPAN"`
	ExpiryDate         string `json:"expiryDate"`
	CardholderName     string `json:"cardholderName"`
	CardIssuanceStatus int    `json:"cardIssuanceStatus"`
	CvvValue           string `json:"cvvValue"`
}

func NewDebitCardDetailResponseResponse() *DebitcardDetailResponse {
	return &DebitcardDetailResponse{}
}

func (r *DebitcardDetailResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *DebitcardDetailResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type AccountDetailResponse struct {
	MobileNo       string `json:"MobileNo"`
	AddressThree   string `json:"AddressThree"`
	OfficePhone    string `json:"OfficePhone"`
	Email          string `json:"Email"`
	ResidencePhone string `json:"ResidencePhone"`
	City           string `json:"City"`
	AccountNumber  string `json:"AccountNumber"`
	AccountStatus  string `json:"AccountStatus"`
	AccountTitle   string `json:"AccountTitle"`
	AddressOne     string `json:"AddressOne"`
	State          string `json:"State"`
	Dob            string `json:"Dob"`
	ITNum          string `json:"ITNum"`
	AddressTwo     string `json:"AddressTwo"`
	Country        string `json:"Country"`
	ErrorCode      string `json:"ErrorCode"`
	ErrorMessage   string `json:"ErrorMessage"`
	Pincode        string `json:"Pincode"`
	CustomerId     string `json:"CustomerId"`
	ApplicantId    string `json:"ApplicantId"`
	AccountBalance string `json:"BalBook"`
}

func NewAccountDetailResponse() *AccountDetailResponse {
	return &AccountDetailResponse{}
}

func (r *AccountDetailResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *AccountDetailResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type UploadAddressProofResponse struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	Status       string `json:"status"`
	Data         string `json:"data"`
}

func NewUploadAddressProofResponse() *UploadAddressProofResponse {
	return &UploadAddressProofResponse{}
}

func (r *UploadAddressProofResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *UploadAddressProofResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type UpdateAddressResponse struct {
	Req_Ref_No   string `json:"Req_Ref_No"`
	ErrorCode    string `json:"ErrorCode"`
	ErrorMessage string `json:"ErrorMessage"`
}

func NewUpdateAddressResponse() *UpdateAddressResponse {
	return &UpdateAddressResponse{}
}

func (r *UpdateAddressResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *UpdateAddressResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type SetDebitCardPinResponse struct {
	ApplicantId   string `json:"ApplicantId"`
	AccountNo     string `json:"AccountNo"`
	TxnIdentifier string `json:"TxnIdentifier"`
	PinReset      string `json:"PinReset"`
	ErrorCode     string `json:"ErrorCode"`
	ErrorMessage  string `json:"ErrorMessage"`
}

func NewSetDebitCardPinResponse() *SetDebitCardPinResponse {
	return &SetDebitCardPinResponse{}
}

func (r *SetDebitCardPinResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *SetDebitCardPinResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type SetDebitCardOTPResponse struct {
	ApplicantId   string `json:"ApplicantId"`
	AccountNo     string `json:"AccountNo"`
	TxnIdentifier string `json:"TxnIdentifier"`
	ServiceType   string `json:"ServiceType"`
	ErrorCode     string `json:"ErrorCode"`
	ErrorMessage  string `json:"ErrorMessage"`
}

func NewSetDebitCardOTPResponse() *SetDebitCardOTPResponse {
	return &SetDebitCardOTPResponse{}
}

func (r *SetDebitCardOTPResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *SetDebitCardOTPResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type VerifyDebitCardOTPResponse struct {
	ApplicantId   string `json:"ApplicantId"`
	AccountNo     string `json:"AccountNo"`
	TxnIdentifier string `json:"TxnIdentifier"`
	ServiceType   string `json:"ServiceType"`
	ErrorCode     string `json:"ErrorCode"`
	ErrorMessage  string `json:"ErrorMessage"`
}

func NewVerifyDebitCardOTPResponse() *VerifyDebitCardOTPResponse {
	return &VerifyDebitCardOTPResponse{}
}

func (r *VerifyDebitCardOTPResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *VerifyDebitCardOTPResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type UpiTransactionHistoryApiResponse struct {
	Response TransactionHistory `json:"Response"`
}

type TransactionHistory struct {
	ResponseCode    string                       `json:"ResponseCode"`
	ResponseMessage string                       `json:"ResponseMessage"`
	Response        []TransactionHistoryResponse `json:"Response,omitempty"`
}

type TransactionHistoryResponse struct {
	Transactiontype string `json:"TRANSACTIONTYPE"`
	Transid         string `json:"TRANSID"`
	Payername       string `json:"PAYERNAME"`
	Payeraddr       string `json:"PAYERADDR"`
	Payeename       string `json:"PAYEENAME"`
	Payeeaddr       string `json:"PAYEEADDR"`
	Payeramount     string `json:"PAYERAMOUNT"`
	Statusdesc      string `json:"STATUSDESC"`
	Transactiondate string `json:"TRANSACTIONDATE"`
	Remarks         string `json:"REMARKS"`
	Errorcode       string `json:"ERRORCODE"`
	Respcode        string `json:"RESPCODE"`
	Payrefid        string `json:"PAYREFID"`
	Refid           string `json:"REFID"`
	Mcccode         string `json:"MCCCODE"`
	ErrorDesc       string `json:"ErrorDesc,omitempty"`
}

func NewUpiTransactionHistoryApiResponse() *UpiTransactionHistoryApiResponse {
	return &UpiTransactionHistoryApiResponse{}
}

func (r *UpiTransactionHistoryApiResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *UpiTransactionHistoryApiResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type UpiChangePinApiResponse struct {
	Response ChangeUpiPinResponse `json:"Response"`
}

type ChangeUpiPinResponse struct {
	ResponseCode    string       `json:"ResponseCode"`
	ResponseMessage string       `json:"ResponseMessage"`
	MessageID       string       `json:"MessageID"`
	Response        ChangeUpiPin `json:"Response"`
}

type ChangeUpiPin struct {
	XML        ChangeUpiPinXMLInfo `json:"?xml"`
	RespSetCre RespSetCreInfo      `json:"ns2:RespSetCre"`
}

type ChangeUpiPinXMLInfo struct {
	Version  string `json:"@version"`
	Encoding string `json:"@encoding"`
}

type RespSetCreInfo struct {
	XMLNSNs2 string               `json:"@xmlns:ns2"`
	XMLNSNs3 string               `json:"@xmlns:ns3"`
	Head     ChangeUpiPinHeadInfo `json:"Head"`
	Resp     ChangeUpiPinRespInfo `json:"Resp"`
	Txn      ChangeUpiPinTxnInfo  `json:"Txn"`
}

type ChangeUpiPinHeadInfo struct {
	MsgID string `json:"@msgId"`
	OrgID string `json:"@orgId"`
	TS    string `json:"@ts"`
	Ver   string `json:"@ver"`
}

type ChangeUpiPinRespInfo struct {
	ReqMsgID string `json:"@reqMsgId"`
	Result   string `json:"@result"`
}

type ChangeUpiPinTxnInfo struct {
	ID     string `json:"@id"`
	Note   string `json:"@note"`
	RefID  string `json:"@refId"`
	RefURL string `json:"@refUrl"`
	TS     string `json:"@ts"`
	Type   string `json:"@type"`
}

func NewUpiChangePinApiResponse() *UpiChangePinApiResponse {
	return &UpiChangePinApiResponse{}
}

func (r *UpiChangePinApiResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *UpiChangePinApiResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}
