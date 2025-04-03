package requests

import (
	"encoding/json"

	"bitbucket.org/paydoh/paydoh-commons/customvalidation"
)

type StatementRequest struct {
	FromDate string `json:"from_date" validate:"required"`
	ToDate   string `json:"to_date" validate:"required"`
	Type     string `json:"type" validate:"required"`
}

func NewStatementRequest() *StatementRequest {
	return &StatementRequest{}
}

func (r *StatementRequest) Validate(payload string) error {
	if err := json.Unmarshal([]byte(payload), r); err != nil {
		return err
	}

	if err := customvalidation.ValidateStruct(r); err != nil {
		return err
	}

	return nil
}
