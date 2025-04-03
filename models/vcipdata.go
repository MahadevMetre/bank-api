package models

import (
	"bankapi/constants"
	"bankapi/requests"
	"database/sql"
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type VcipData struct {
	Id                    int64     `json:"id"`
	UserId                string    `json:"user_id"`
	FirstName             string    `json:"first_name"`
	MiddleName            string    `json:"middle_name"`
	LastName              string    `json:"last_name"`
	PanNumber             string    `json:"pan_number"`
	AadharReferenceNumber string    `json:"aadhar_reference_number"`
	VKYCCompletion        string    `json:"vkyc_completion"`
	VKYCAuditStatus       string    `json:"vkyc_audit_status"`
	AuditorRejectRemarks  string    `json:"auditor_reject_remarks"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

func NewVcipData() *VcipData {
	return &VcipData{}
}

func (v *VcipData) Bind(request *requests.IncomingVcipData, userId string) {
	v.UserId = userId
	v.FirstName = request.Firstname
	v.MiddleName = request.MiddleName
	v.LastName = request.LastName
	v.PanNumber = request.PanNumber
	v.AadharReferenceNumber = request.AadharReferenceNumber
	v.VKYCCompletion = request.VKYCCompletion
	v.VKYCAuditStatus = request.VKYCAuditStatus
	v.AuditorRejectRemarks = request.AuditorRejectRemarks
}

func (v *VcipData) Marshal() ([]byte, error) {
	return json.Marshal(v)
}

func (v *VcipData) Unmarshal(data []byte) error {
	return json.Unmarshal(data, v)
}

func InsertVcipData(db *sql.DB, vcipData *VcipData) error {

	if _, err := db.Exec(""+
		"INSERT INTO kyc_vcip_data (user_id, first_name, middle_name, last_name, pan_number, aadhar_reference_number, vkyc_completion, vkyc_audit_status, auditor_reject_remarks) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		vcipData.UserId,
		vcipData.FirstName,
		vcipData.MiddleName,
		vcipData.LastName,
		vcipData.PanNumber,
		vcipData.AadharReferenceNumber,
		vcipData.VKYCCompletion,
		vcipData.VKYCAuditStatus,
		vcipData.AuditorRejectRemarks,
	); err != nil {
		return err
	}

	return nil
}

func GetVcipDataByUserId(db *sql.DB, userId string) (*VcipData, error) {
	vcipdata := NewVcipData()

	row := db.QueryRow(
		"SELECT id, user_id, first_name, middle_name, last_name, pan_number, aadhar_reference_number, vkyc_completion, vkyc_audit_status, auditor_reject_remarks, created_at, updated_at FROM kyc_vcip_data WHERE user_id = $1",
		userId,
	)

	if err := row.Scan(
		&vcipdata.Id,
		&vcipdata.UserId,
		&vcipdata.FirstName,
		&vcipdata.MiddleName,
		&vcipdata.LastName,
		&vcipdata.PanNumber,
		&vcipdata.AadharReferenceNumber,
		&vcipdata.VKYCCompletion,
		&vcipdata.VKYCAuditStatus,
		&vcipdata.AadharReferenceNumber,
		&vcipdata.CreatedAt,
		&vcipdata.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrNoDataFound
		}

		return nil, err
	}

	return vcipdata, nil
}

func UpdateVcipData(db *sql.DB, updateModel *VcipData, userId string) error {
	var clause []string
	var params []interface{}

	if updateModel.FirstName != "" {
		clause = append(clause, "first_name = $"+strconv.Itoa(len(clause)+1))
		params = append(params, updateModel.FirstName)
	}

	if updateModel.MiddleName != "" {
		clause = append(clause, "middle_name = $"+strconv.Itoa(len(clause)+1))
		params = append(params, updateModel.MiddleName)
	}

	if updateModel.LastName != "" {
		clause = append(clause, "last_name = $"+strconv.Itoa(len(clause)+1))
		params = append(params, updateModel.LastName)
	}

	if updateModel.PanNumber != "" {
		clause = append(clause, "pan_number = $"+strconv.Itoa(len(clause)+1))
		params = append(params, updateModel.PanNumber)
	}

	if updateModel.AadharReferenceNumber != "" {
		clause = append(clause, "aadhar_reference_number = $"+strconv.Itoa(len(clause)+1))
		params = append(params, updateModel.AadharReferenceNumber)
	}

	if updateModel.VKYCCompletion != "" {
		clause = append(clause, "vkyc_completion = $"+strconv.Itoa(len(clause)+1))
		params = append(params, updateModel.VKYCCompletion)
	}

	if updateModel.VKYCAuditStatus != "" {
		clause = append(clause, "vkyc_audit_status = $"+strconv.Itoa(len(clause)+1))
		params = append(params, updateModel.VKYCAuditStatus)
	}

	if updateModel.AuditorRejectRemarks != "" {
		clause = append(clause, "auditor_reject_remarks = $"+strconv.Itoa(len(clause)+1))
		params = append(params, updateModel.AuditorRejectRemarks)
	}

	if len(clause) > 0 {
		clause = append(clause, "updated_at = NOW()")
		if _, err := db.Exec("UPDATE kyc_vcip_data SET "+strings.Join(clause, ", ")+" WHERE user_id ="+userId, params...); err != nil {
			return err
		}

		return nil
	}

	return nil
}
