package models

import (
	"bankapi/constants"
	"database/sql"
)

type State struct {
	StateName string `json:"state_name"`
}

func NewState() *State {
	return &State{}
}

func FindAllStatesDistinct(db *sql.DB) ([]string, error) {
	states := make([]string, 0)

	row, err := db.Query("SELECT  DISTINCT  state_name from state_city_master ORDER BY state_name;")

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrNoDataFound
		}
		return nil, err
	}

	defer row.Close()

	for row.Next() {
		var state string
		if err := row.Scan(
			&state,
		); err != nil {
			return nil, err
		}

		states = append(states, state)
	}

	return states, nil
}
