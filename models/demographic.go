package models

import (
	"bankapi/constants"
	"bankapi/responses"
	"database/sql"
	"encoding/json"
)

type DemographicData struct {
	UserID                     string `json:"user_id"`
	TxnIdentifier              string `json:"txn_identifier"`
	FirstName                  string `json:"first_name"`
	MiddleName                 string `json:"middle_name"`
	LastName                   string `json:"last_name"`
	PanNumber                  string `json:"pan_number"`
	Ret                        string `json:"ret"`
	GeneratedKeyForKYCResponse string `json:"generated_key_for_kyc_response"`
	PinCode                    string `json:"pin_code"`
	RespDesc                   string `json:"resp_desc"`
	PostOffice                 string `json:"post_office"`
	VTC                        string `json:"vtc"`
	UIDAIAuthCode              string `json:"uidai_auth_code"`
	IASKRefID                  string `json:"iask_ref_id"`
	Locality                   string `json:"locality"`
	Txn                        string `json:"txn"`
	RRN                        string `json:"rrn"`
	SubDistrict                string `json:"sub_district"`
	Street                     string `json:"street"`
	Landmark                   string `json:"landmark"`
	RespCode                   string `json:"resp_code"`
	DemographicResponseID      int    `json:"demographic_response_id"`
	UID                        string `json:"uid"`
	PHT                        []byte `json:"pht"`
	Dist                       string `json:"dist"`
	State                      string `json:"state"`
	Co                         string `json:"co"`
	House                      string `json:"house"`
	Gender                     string `json:"gender"`
	Phone                      string `json:"phone"`
	DOB                        string `json:"dob"`
	PName                      string `json:"pname"`
	Email                      string `json:"email"`
}

type UserDemographicData struct {
	UserID          string          `json:"user_id"`
	DemographicData DemographicInfo `json:"demographic_data"`
}

type DemographicInfo struct {
	TxnIdentifier         string             `json:"txn_identifier"`
	FirstName             string             `json:"first_name"`
	MiddleName            string             `json:"middle_name"`
	LastName              string             `json:"last_name"`
	PanNumber             string             `json:"pan_number"`
	Ret                   string             `json:"ret"`
	GeneratedKey          string             `json:"generated_key_for_kyc_response"`
	PinCode               string             `json:"pin_code"`
	RespDesc              string             `json:"resp_desc"`
	PostOffice            string             `json:"post_office"`
	Vtc                   string             `json:"vtc"`
	UidaiAuthCode         string             `json:"uidai_auth_code"`
	IaskRefID             string             `json:"iask_ref_id"`
	Locality              string             `json:"locality"`
	Txn                   string             `json:"txn"`
	Rrn                   string             `json:"rrn"`
	SubDistrict           string             `json:"sub_district"`
	Street                string             `json:"street"`
	Landmark              string             `json:"landmark"`
	RespCode              string             `json:"resp_code"`
	DemographicResponseID int                `json:"demographic_response_id"`
	UidData               DemographicUIDData `json:"uidData"`
}

type DemographicUIDData struct {
	Uid       string         `json:"uid"`
	Poa       DemographicPOA `json:"poa"`
	Poi       DemographicPOI `json:"poi"`
	Pht       string         `json:"pht"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt string         `json:"updated_at"`
}

type DemographicPOA struct {
	Dist  string `json:"dist"`
	State string `json:"state"`
	Co    string `json:"co"`
	House string `json:"house"`
}

type DemographicPOI struct {
	Gender string `json:"gender"`
	Phone  string `json:"phone"`
	Dob    string `json:"dob"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

func NewDemographicData() *DemographicData {
	return &DemographicData{}
}

func NewUserDemographicData() *UserDemographicData {
	return &UserDemographicData{}
}

func (r *DemographicData) Bind(p *responses.DemographicResponse) error {
	r.UserID = p.ApplicantId
	r.TxnIdentifier = p.TxnIdentifier
	r.FirstName = p.FirstName
	r.LastName = p.LastName
	r.MiddleName = p.MiddleName
	r.PanNumber = p.PanNumber
	r.RespDesc = p.ErrorMessage // Assuming ErrorMessage maps to RespDesc
	r.PostOffice = p.Root.Postoffice
	r.VTC = p.Root.Vtc
	r.UIDAIAuthCode = p.Root.UidaiAuthCode
	r.IASKRefID = p.Root.IaskRefID
	r.Locality = p.Root.Locality
	r.Txn = p.Root.Txn
	r.RRN = p.Root.Rrn
	r.SubDistrict = p.Root.Subdistrict
	r.Street = p.Root.Street
	r.Landmark = p.Root.Landmark
	r.RespCode = p.Root.RespCode
	r.DemographicResponseID = 0 // Assuming this is initialized separately
	r.UID = p.Root.UIDData.UID
	r.PHT = []byte(p.Root.UIDData.Pht)
	r.Dist = p.Root.UIDData.Poa.Dist
	r.State = p.Root.UIDData.Poa.State
	r.Co = p.Root.UIDData.Poa.Co
	r.House = p.Root.UIDData.Poa.House
	r.Gender = p.Root.UIDData.Poi.Gender
	r.Phone = p.Root.UIDData.Poi.Phone
	r.DOB = p.Root.UIDData.Poi.Dob
	r.PName = p.Root.UIDData.Poi.Name
	r.Email = p.Root.UIDData.Poi.Email
	r.Ret = p.Root.Ret
	r.GeneratedKeyForKYCResponse = p.Root.GeneratedKeyForKycResponse
	r.PinCode = p.Root.Pincode
	r.RespDesc = p.Root.RespDesc
	r.PostOffice = p.Root.Postoffice
	r.VTC = p.Root.Vtc
	r.UIDAIAuthCode = p.Root.UidaiAuthCode
	r.IASKRefID = p.Root.IaskRefID
	r.Locality = p.Root.Locality
	r.Txn = p.Root.Txn
	r.RRN = p.Root.Rrn
	r.SubDistrict = p.Root.Subdistrict
	r.Street = p.Root.Street
	r.Landmark = p.Root.Landmark
	r.RespCode = p.Root.RespCode
	return nil
}

func InsertDemographicData(db *sql.DB, data *DemographicData) error {
	if _, err := db.Exec("CALL insert_demographic_data($1)", data); err != nil {
		return err
	}

	return nil
}

func GetUserDemographicData(db *sql.DB, userId string) (*UserDemographicData, error) {
	var userID string
	var demographicJSON string

	row := db.QueryRow("SELECT user_id, demographic_data FROM user_demographic_data WHERE user_id = $1", userId)

	if err := row.Scan(
		&userID,
		&demographicJSON,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrNoDataFound
		}
		return nil, err
	}

	userDemographicData := NewUserDemographicData()
	userDemographicData.UserID = userID

	if err := json.Unmarshal([]byte(demographicJSON), &userDemographicData.DemographicData); err != nil {
		return nil, err
	}

	return userDemographicData, nil
}
