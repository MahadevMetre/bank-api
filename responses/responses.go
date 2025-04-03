package responses

import (
	"errors"

	"bitbucket.org/paydoh/paydoh-commons/pkg/json"
)

type UserInformationResponse struct {
	UserId      string `json:"user_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	MiddleName  string `json:"middle_name"`
	VTC         string `json:"vtc"`
	CareOf      string `json:"care_of"`
	House       string `json:"house"`
	LandMark    string `json:"landmark"`
	Street      string `json:"street"`
	PostOffice  string `json:"post_office"`
	District    string `json:"district"`
	State       string `json:"state"`
	PinCode     string `json:"pincode"`
	Gender      string `json:"gender"`
	DateOfBirth string `json:"date_of_birth"`
}

type FetchNomineeResponseData struct {
	NomDOB              string `json:"nominee_dob"`
	NomCity             string `json:"nominee_city"`
	GuardianCountry     string `json:"guardian_country"`
	NomReqType          string `json:"nominee_request_type"`
	GuardianZipcode     string `json:"guardian_zip_code"`
	GuardianCity        string `json:"guardian_city"`
	GuardianState       string `json:"guardian_state"`
	ApplicantID         string `json:"applicant_id"`
	NomUpdateDtTime     string `json:"nominee_update_dt_time"`
	NomZipcode          string `json:"nominee_zip_code"`
	NomName             string `json:"nominee_name"`
	AccountNo           string `json:"account_number"`
	NomAppID            string `json:"nominee_applicant_id"`
	GuardianName        string `json:"guardian_name"`
	GuardianNomRelation string `json:"guardian_nominee_relation"`
	GuardianAddressL3   string `json:"guardian_address_line_3"`
	GuardianAddressL1   string `json:"guardian_address_line_1"`
	GuardianAddressL2   string `json:"guardian_address_line_2"`
	NomAddressL1        string `json:"nominee_address_line_1"`
	NomAddressL3        string `json:"nominee_address_line_3"`
	NomAddressL2        string `json:"nominee_address_line_2"`
	NomRelation         string `json:"nominee_relation"`
	NomState            string `json:"nominee_state"`
	NomCountry          string `json:"nominee_country"`
	NomineeUpdateTime   string `json:"nominee_update_time"`
	NomineeActive       bool   `json:"is_active"`
}

func NewUserInformationResponse() *UserInformationResponse {
	return &UserInformationResponse{}
}

func NewFetchNomineeData() *FetchNomineeResponseData {
	return &FetchNomineeResponseData{}
}

func (r *FetchNomineeResponseData) Bind(d *FetchNomineeResponse) error {
	r.AccountNo = d.AccountNo
	r.ApplicantID = d.ApplicantID
	r.NomAppID = d.NomAppID
	r.NomAddressL1 = d.NomAddressL1
	r.NomAddressL2 = d.NomAddressL2
	r.NomAddressL3 = d.NomAddressL3
	r.NomCity = d.NomCity
	r.NomCountry = d.NomCountry
	r.NomDOB = d.NomDOB
	r.NomName = d.NomName
	r.NomReqType = d.NomReqType
	r.NomState = d.NomState
	r.NomZipcode = d.NomZipcode
	r.NomRelation = d.NomRelation
	r.GuardianAddressL1 = d.GuardianAddressL1
	r.GuardianAddressL2 = d.GuardianAddressL2
	r.GuardianAddressL3 = d.GuardianAddressL3
	r.GuardianCity = d.GuardianCity
	r.GuardianCountry = d.GuardianCountry
	r.GuardianName = d.GuardianName
	r.GuardianNomRelation = d.GuardianNomRelation
	r.GuardianState = d.GuardianState
	r.GuardianZipcode = d.GuardianZipcode
	r.NomineeUpdateTime = d.NomUpdateDtTime
	r.NomineeActive = true

	return nil
}

func (r *FetchNomineeResponseData) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *FetchNomineeResponseData) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

func (r *UserInformationResponse) Bind(userId string, d *DemographicResponse) error {
	r.UserId = userId
	r.FirstName = d.FirstName
	r.MiddleName = d.MiddleName
	r.LastName = d.LastName
	r.VTC = d.Root.Vtc
	r.LandMark = d.Root.Landmark
	r.CareOf = d.Root.UIDData.Poa.Co
	r.House = d.Root.UIDData.Poa.House
	r.LandMark = d.Root.Landmark
	r.Street = d.Root.Street
	r.PostOffice = d.Root.Postoffice
	r.District = d.Root.UIDData.Poa.Dist
	r.State = d.Root.UIDData.Poa.State
	r.PinCode = d.Root.Pincode
	r.Gender = d.Root.UIDData.Poi.Gender
	r.DateOfBirth = d.Root.UIDData.Poi.Dob

	return nil
}

func (r *UserInformationResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r *UserInformationResponse) UnMarshal(data []byte) error {
	return json.Unmarshal(data, r)
}

type UserNotificationPref struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Data    Datas
}
type Datas struct {
	ID        string `json:"_id"`
	CreatedAt string `json:"createdat"`
	Updatedat string `json:"updatedat"`
	Userid    string `json:"userid"`
	Prefrence []Prefrence
}

type Prefrence struct {
	Enable           bool `json:"enabled"`
	ID               bool `json:"id"`
	NotificationType bool `json:"notificationtype"`
}

func NewUserNotificationPref() *UserNotificationPref {
	return &UserNotificationPref{}
}

func (res *UserNotificationPref) Decode(data []byte) error {

	if err := json.Unmarshal(data, res); err != nil {
		return errors.New(err.Error())
	}

	return nil
}
