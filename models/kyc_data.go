package models

import (
	"bankapi/constants"
	"bankapi/requests"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type KycUpdateData struct {
	Id        uuid.UUID `json:"id"`
	UserId    string    `json:"user_id"`
	Status    string    `json:"status"`
	Acom      int       `json:"acom"`
	Astat     int       `json:"astat"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewKycUpdateData() *KycUpdateData {
	return &KycUpdateData{}
}

func (kyc *KycUpdateData) Bind(request *requests.KycDataUpdateRequest) error {
	kyc.UserId = request.UserId
	kyc.Status = request.Status
	kyc.Astat = 0
	kyc.Acom = 0

	return nil
}

func InsertKycUpdateData(db *sql.DB, model *KycUpdateData) error {
	if _, err := db.Exec(
		"INSERT INTO kyc_update_data (user_id, status, acom, astat) VALUES ($1, $2, $3, $4)",
		model.UserId,
		model.Status,
		model.Acom,
		model.Astat,
	); err != nil {
		return err
	}

	return nil
}

func UpdateKycUpdateData(db *sql.DB, model *KycUpdateData) error {
	if _, err := db.Exec(
		"UPDATE kyc_update_data SET status = $1, acom = $2, astat = $3 WHERE user_id = $4",
		model.Status,
		model.Acom,
		model.Astat,
		model.UserId,
	); err != nil {
		return err
	}

	return nil
}

func FindOneKycUpdateData(db *sql.DB, userId string) (*KycUpdateData, error) {
	kycUpdateData := NewKycUpdateData()
	row := db.QueryRow("SELECT id, user_id, status, acom, astat, created_at, updated_at FROM kyc_update_data WHERE user_id = $1", userId)

	if err := row.Scan(
		&kycUpdateData.Id,
		&kycUpdateData.UserId,
		&kycUpdateData.Status,
		&kycUpdateData.Acom,
		&kycUpdateData.Astat,
		&kycUpdateData.CreatedAt,
		&kycUpdateData.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrNoDataFound
		}
		return nil, err
	}

	return kycUpdateData, nil
}

func FindOneAndUpdateKycUpdateData(db *sql.DB, userId string, status string, acom string, astat string) (*KycUpdateData, error) {
	kycUpdateData := NewKycUpdateData()
	row := db.QueryRow("SELECT id, user_id, status, acom, astat, created_at, updated_at FROM kyc_update_data WHERE user_id = $1", userId)

	if err := row.Scan(
		&kycUpdateData.Id,
		&kycUpdateData.UserId,
		&kycUpdateData.Status,
		&kycUpdateData.Acom,
		&kycUpdateData.Astat,
		&kycUpdateData.CreatedAt,
		&kycUpdateData.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrNoDataFound
		}
		return nil, err
	}

	kycUpdateData.Status = status
	kycUpdateData.Acom = 0
	kycUpdateData.Astat = 0

	if _, err := db.Exec("UPDATE kyc_update_data SET status = $1, acom = $2, astat = $3, updated_at = NOW() WHERE user_id = $4", status, acom, astat, userId); err != nil {
		return nil, err
	}

	return kycUpdateData, nil
}

func GetKycUpdateDataByUserId(db *sql.DB, userId string) (*KycUpdateData, error) {
	kycUpdateData := NewKycUpdateData()
	row := db.QueryRow("SELECT id, user_id, status, acom, astat, created_at, updated_at FROM kyc_update_data WHERE user_id = $1", userId)

	if err := row.Scan(
		&kycUpdateData.Id,
		&kycUpdateData.UserId,
		&kycUpdateData.Status,
		&kycUpdateData.Acom,
		&kycUpdateData.Astat,
		&kycUpdateData.CreatedAt,
		&kycUpdateData.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrNoDataFound
		}
		return nil, err
	}

	return kycUpdateData, nil
}
