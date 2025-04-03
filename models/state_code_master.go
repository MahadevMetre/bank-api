package models

import (
	"bankapi/config"
	"bankapi/constants"
	"database/sql"

	"github.com/google/uuid"
)

type StateCodeMaster struct {
	Id        uuid.UUID `json:"id"`
	StateName string    `json:"state_name"`
	StateCode string    `json:"state_code"`
}

func GetStateCodeByName(name string) (*StateCodeMaster, error) {
	db := config.GetDB()
	var stateData StateCodeMaster
	row := db.QueryRow("SELECT id, state_name, state_code FROM state_code_master WHERE LOWER(state_name) = $1", name)

	if err := row.Scan(
		&stateData.Id,
		&stateData.StateName,
		&stateData.StateCode,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrNoDataFound
		}

		return nil, err
	}

	return &stateData, nil
}
