package models

import (
	"bankapi/requests"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type PaymentBeneficiaryData struct {
	Id            uuid.UUID            `json:"Id"`
	ApplicationID string               `json:"application_id"`
	ServiceName   string               `json:"service_name"`
	ProductType   string               `json:"product_type"`
	SourcedBy     string               `json:"sourced_by"`
	CallbackName  string               `json:"callback_name"`
	CbsStatus     []requests.CbsStatus `json:"cbs_status"`
}

func NewPaymentBeneficiaryData(
	data *requests.PaymentCallbackRequestData) *PaymentBeneficiaryData {
	return &PaymentBeneficiaryData{
		Id:            uuid.New(),
		ApplicationID: data.ApplicationId,
		ServiceName:   data.ServiceName,
		ProductType:   data.ProductType,
		SourcedBy:     data.SourcedBy,
		CallbackName:  data.CallBackName,
		CbsStatus:     data.CbsStatus,
	}
}

func (k *PaymentBeneficiaryData) CreateRecord(db *sql.DB) error {
	query := `
        INSERT INTO payment_beneficiary_data (
            id, application_id, service_name, product_type, sourced_by, callback_name, cbs_status
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)
    `

	cbsStatusJSON, err := json.Marshal(k.CbsStatus)
	if err != nil {
		return fmt.Errorf("failed to marshal cbs_status: %w", err)
	}

	_, err = db.Exec(query, k.Id, k.ApplicationID, k.ServiceName, k.ProductType, k.SourcedBy, k.CallbackName, cbsStatusJSON)
	if err != nil {
		return fmt.Errorf("failed to insert record: %w", err)
	}

	return nil
}

func (p *PaymentBeneficiaryData) CheckUserApplicationId(db *sql.DB) error {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM account_data WHERE application_id = $1)", p.ApplicationID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if application_id exists: %w", err)
	}
	if !exists {
		return fmt.Errorf("application_id %s not found", p.ApplicationID)
	}
	return nil
}

func (p *PaymentBeneficiaryData) UpdateRecord(db *sql.DB) error {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM account_data WHERE application_id = $1)", p.ApplicationID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if application_id exists: %w", err)
	}
	if !exists {
		return fmt.Errorf("application_id %s not found", p.ApplicationID)
	}

	// !!updating the record
	query := `
		UPDATE payment_beneficiary_data
		SET service_name = $1,
			product_type = $2,
			sourced_by = $3,
			callback_name = $4,
			cbs_status = $5,
			updated_at = CURRENT_TIMESTAMP
		WHERE application_id = $6
	`

	cbsStatusJSON, err := json.Marshal(p.CbsStatus)
	if err != nil {
		return fmt.Errorf("failed to marshal cbs_status: %w", err)
	}

	_, err = db.Exec(query, p.ServiceName, p.ProductType, p.SourcedBy, p.CallbackName, cbsStatusJSON, p.ApplicationID)
	if err != nil {
		return fmt.Errorf("failed to update record: %w", err)
	}

	return nil
}
