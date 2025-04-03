package models

import (
	"bankapi/constants"
	"bankapi/requests"
	"database/sql"
	"encoding/json"
	"time"

	"bitbucket.org/paydoh/paydoh-commons/types"
	"github.com/google/uuid"
)

type AddressUpdate struct {
	Id                   uuid.UUID            `json:"id"`
	UserId               string               `json:"user_id"`
	ReqRefNo             string               `json:"Req_Ref_No"`
	Status               string               `json:"status"`
	CommunicationAddress types.NullableString `json:"communication_address"`
	CurrentStatus        types.NullableString `json:"current_status"`
	CreatedAt            time.Time            `json:"created_at"`
	UpdatedAt            time.Time            `json:"updated_at"`
}

func NewAddressUpdate() *AddressUpdate {
	return &AddressUpdate{}
}

func InsertAddressUpdate(db *sql.DB, userId, ReqRefNo string, communicationAddress *requests.CommunicationAddress) error {

	var err error
	var jsonRawMessage []byte

	if communicationAddress != nil {
		jsonRawMessage, err = json.Marshal(communicationAddress)
		if err != nil {
			return err
		}
		if _, err := db.Exec("INSERT INTO address_update_data (user_id, Req_Ref_No,communication_address) VALUES ($1, $2, $3)", userId, ReqRefNo, jsonRawMessage); err != nil {
			return err
		}
	}

	return nil
}

func GetAddressUpdateDataByUserId(db *sql.DB, userId string) ([]AddressUpdate, error) {
	updateDatas := make([]AddressUpdate, 0)
	rows, err := db.Query("SELECT id, user_id, Req_Ref_No, status, created_at, updated_at,communication_address,current_status FROM address_update_data WHERE status = false AND user_id = $1", userId)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrNoDataFound
		}
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		data := NewAddressUpdate()

		if err := rows.Scan(
			&data.Id,
			&data.UserId,
			&data.ReqRefNo,
			&data.Status,
			&data.CreatedAt,
			&data.UpdatedAt,
			&data.CommunicationAddress,
			&data.CurrentStatus,
		); err != nil {
			return nil, err
		}
		updateDatas = append(updateDatas, *data)
	}

	return updateDatas, nil
}

func GetAddressUpdateRequestListByUserId(db *sql.DB, userId string) ([]AddressUpdate, error) {
	updateDatas := make([]AddressUpdate, 0)
	rows, err := db.Query("SELECT id, user_id, Req_Ref_No, status, created_at, updated_at,communication_address,current_status FROM address_update_data WHERE user_id = $1", userId)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrNoDataFound
		}
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		data := NewAddressUpdate()

		if err := rows.Scan(
			&data.Id,
			&data.UserId,
			&data.ReqRefNo,
			&data.Status,
			&data.CreatedAt,
			&data.UpdatedAt,
			&data.CommunicationAddress,
			&data.CurrentStatus,
		); err != nil {
			return nil, err
		}
		updateDatas = append(updateDatas, *data)
	}

	return updateDatas, nil
}

func GetAddressUpdateRequestList(db *sql.DB) ([]AddressUpdate, error) {
	updateDatas := make([]AddressUpdate, 0)
	rows, err := db.Query("SELECT id, user_id, Req_Ref_No, status, created_at, updated_at,communication_address,current_status FROM address_update_data")

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrNoDataFound
		}
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		data := NewAddressUpdate()

		if err := rows.Scan(
			&data.Id,
			&data.UserId,
			&data.ReqRefNo,
			&data.Status,
			&data.CreatedAt,
			&data.UpdatedAt,
			&data.CommunicationAddress,
			&data.CurrentStatus,
		); err != nil {
			return nil, err
		}
		updateDatas = append(updateDatas, *data)
	}

	return updateDatas, nil
}
