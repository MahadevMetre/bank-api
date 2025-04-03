package models

import (
	"bankapi/config"
	"bankapi/constants"
	"bankapi/requests"
	"database/sql"
	"log"

	"github.com/google/uuid"
)

type KycAuditData struct {
	Id                uuid.UUID `json:"id"`
	UserId            string    `json:"UserId"`
	MobileNo          string    `json:"MobileNo"`
	CallBackName      string    `json:"CallBackName"`
	ApplicantId       string    `json:"ApplicantId"`
	SourcedBy         string    `json:"SourcedBy"`
	VKYCAuditStatus   string    `json:"VKYCAuditStatus"`
	ProductType       string    `json:"ProductType"`
	AuditRejectReason string    `json:"AuditRejectReason"`
}

func NewKycAuditData(
	data *requests.KycAuditRequestData) *KycAuditData {
	return &KycAuditData{
		Id:                uuid.New(),
		UserId:            data.UserId,
		MobileNo:          data.MobileNo,
		CallBackName:      data.CallBackName,
		ApplicantId:       data.ApplicantId,
		SourcedBy:         data.SourcedBy,
		VKYCAuditStatus:   data.VKYCAuditStatus,
		ProductType:       data.ProductType,
		AuditRejectReason: data.AuditRejectReason,
	}
}

// create kyc audit data record in DB
func (k *KycAuditData) CreateRecord(db *sql.DB) error {
	sqlStatement := `
		INSERT INTO kyc_audit_data (id, user_id, mobile_no, callback_name, applicant_id, sourced_by, vkyc_audit_status, product_type, audit_reject_reason)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id`

	var id uuid.UUID
	err := db.QueryRow(sqlStatement, k.Id, k.UserId, k.MobileNo, k.CallBackName, k.ApplicantId, k.SourcedBy, k.VKYCAuditStatus, k.ProductType, k.AuditRejectReason).Scan(&id)

	if err != nil {
		return err
	}

	return nil
}

// create kyc audit data record in DB
func GetRecordByUserId(userId string) (*KycAuditData, error) {
	db := config.GetDB()

	var auditData KycAuditData

	sqlStatement := `
		SELECT id, user_id, mobile_no, callback_name, applicant_id, sourced_by, vkyc_audit_status, product_type, audit_reject_reason
		FROM kyc_audit_data WHERE user_id=$1`

	err := db.QueryRow(sqlStatement, userId).Scan(&auditData.Id, &auditData.UserId, &auditData.MobileNo, &auditData.CallBackName, &auditData.ApplicantId, &auditData.SourcedBy, &auditData.VKYCAuditStatus, &auditData.ProductType, &auditData.AuditRejectReason)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrNoDataFound
		}
		return nil, err
	}

	return &auditData, nil
}

// find the present data in DB by mobile number
func (k *KycAuditData) CheckMobileNumberExists(db *sql.DB, mobileNumber string) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM kyc_audit_data WHERE mobile_no = $1)"

	var exists bool
	err := db.QueryRow(query, mobileNumber).Scan(&exists)
	if err != nil {
		log.Printf("Error checking mobile number existence: %v", err)
		return false, err
	}

	return exists, nil
}
