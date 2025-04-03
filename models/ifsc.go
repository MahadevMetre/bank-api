package models

import (
	"bankapi/constants"
	"bankapi/responses"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IFSCDto struct {
	Id            int64     `json:"Id"`
	BankName      string    `json:"bank_name"`
	IfscCode      string    `json:"ifsc_code"`
	BranchName    string    `json:"branch_name"`
	BranchCity    string    `json:"branch_city"`
	BranchState   string    `json:"branch_state"`
	BranchCountry string    `json:"branch_country"`
	PaymentMode   string    `json:"payment_mode"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type LastSyncIfscData struct {
	Id         primitive.ObjectID `bson:"_id" json:"id"`
	LastSynced string             `bson:"last_synced" json:"last_synced"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time          `bson:"updated_at" json:"updated_at"`
}

type IFSCDtos []IFSCDto

func NewIFSCDto() *IFSCDto {
	return &IFSCDto{}
}

func NewLastSyncIfscData() *LastSyncIfscData {
	return &LastSyncIfscData{}
}

func NewIfscDtos() *IFSCDtos {
	return &IFSCDtos{}
}

func (d *IFSCDto) Bind(bankDetail *responses.BanksIFSCDtl) error {
	d.BankName = bankDetail.BankName
	d.IfscCode = bankDetail.IFSCCode
	d.BranchName = bankDetail.BranchName
	d.BranchCity = bankDetail.BranchCity
	d.BranchState = bankDetail.BranchState
	d.BranchCountry = string(bankDetail.BranchCountry)
	d.PaymentMode = string(bankDetail.PaymentMode)
	return nil
}

func (d *IFSCDtos) Bind(response *responses.IfscDataResponse) error {
	for _, bankDetails := range response.BanksIFSCDtls {
		ifscData := NewIFSCDto()
		ifscData.BankName = bankDetails.BankName
		ifscData.IfscCode = bankDetails.IFSCCode
		ifscData.BranchName = bankDetails.BranchName
		ifscData.BranchCity = bankDetails.BranchCity
		ifscData.BranchState = bankDetails.BranchState
		ifscData.BranchCountry = string(bankDetails.BranchCountry)
		ifscData.PaymentMode = string(bankDetails.PaymentMode)
		*d = append(*d, *ifscData)
	}

	return nil
}

func (IFSCDto *IFSCDto) Marshal() ([]byte, error) {
	return json.Marshal(IFSCDto)
}

func (IFSCDto *IFSCDto) Unmarshal(data []byte) error {
	return json.Unmarshal(data, IFSCDto)
}

func InsertIFSCData(db *sql.DB, data *IFSCDto) error {
	if _, err := db.Exec("CALL insert_ifsc_data($1)", data); err != nil {
		return err
	}

	return nil
}

func InsertManyIfscData(db *sql.DB, data []IFSCDto) error {
	if _, err := db.Exec("CALL insert_many_ifsc_data($1)", data); err != nil {
		return err
	}
	return nil
}

func GetIFSCData(db *sql.DB, ifscCode string) (*IFSCDto, error) {

	IFSCDto := NewIFSCDto()

	row := db.QueryRow("SELECT id, bank_name, ifsc_code, branch_name, branch_city, branch_state, branch_country, payment_mode, created_at, updated_at FROM ifsc_data WHERE ifsc_code = $1", ifscCode)

	if err := row.Scan(
		&IFSCDto.Id,
		&IFSCDto.BankName,
		&IFSCDto.IfscCode,
		&IFSCDto.BranchName,
		&IFSCDto.BranchCity,
		&IFSCDto.BranchState,
		&IFSCDto.BranchCountry,
		&IFSCDto.PaymentMode,
		&IFSCDto.CreatedAt,
		&IFSCDto.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrNoDataFound
		}
		return nil, err
	}

	return IFSCDto, nil
}

func GetDistinctBankNames(db *sql.DB) ([]string, error) {

	rows, err := db.Query("SELECT DISTINCT bank_name FROM ifsc_data")

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, constants.ErrNoDataFound
		}
		return nil, err
	}

	defer rows.Close()

	bankNames := make([]string, 0)

	for rows.Next() {
		var bankName string
		if err := rows.Scan(
			&bankName,
		); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, constants.ErrNoDataFound
			}
			return nil, err
		}

		bankNames = append(bankNames, bankName)
	}

	if len(bankNames) == 0 {
		return nil, constants.ErrNoDataFound
	}

	return bankNames, nil

}

func GetIfscByBankName(db *sql.DB, bankName string) ([]IFSCDto, error) {
	rows, err := db.Query("SELECT id, bank_name, ifsc_code, branch_name, branch_city, branch_state, branch_country, payment_mode, created_at, updated_at FROM ifsc_data WHERE bank_name = $1", bankName)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, constants.ErrNoDataFound
		}
		return nil, err
	}

	defer rows.Close()

	IFSCDtos := make([]IFSCDto, 0)

	for rows.Next() {
		IFSCDto := NewIFSCDto()
		if err := rows.Scan(
			&IFSCDto.Id,
			&IFSCDto.BankName,
			&IFSCDto.IfscCode,
			&IFSCDto.BranchName,
			&IFSCDto.BranchCity,
			&IFSCDto.BranchState,
			&IFSCDto.BranchCountry,
			&IFSCDto.PaymentMode,
			&IFSCDto.CreatedAt,
			&IFSCDto.UpdatedAt,
		); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, constants.ErrNoDataFound
			}
			return nil, err
		}

		IFSCDtos = append(IFSCDtos, *IFSCDto)
	}

	if len(IFSCDtos) == 0 {
		return nil, constants.ErrNoDataFound
	}

	return IFSCDtos, nil
}

func GetAllIFSCData(db *sql.DB) ([]IFSCDto, error) {
	rows, err := db.Query("SELECT id, bank_name, ifsc_code, branch_name, branch_city, branch_state, branch_country, payment_mode, created_at, updated_at FROM ifsc_data")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, constants.ErrNoDataFound
		}
		return nil, err
	}

	defer rows.Close()

	IFSCDtos := make([]IFSCDto, 0)

	for rows.Next() {
		IFSCDto := NewIFSCDto()
		if err := rows.Scan(
			&IFSCDto.Id,
			&IFSCDto.BankName,
			&IFSCDto.IfscCode,
			&IFSCDto.BranchName,
			&IFSCDto.BranchCity,
			&IFSCDto.BranchState,
			&IFSCDto.BranchCountry,
			&IFSCDto.PaymentMode,
			&IFSCDto.CreatedAt,
			&IFSCDto.UpdatedAt,
		); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, constants.ErrNoDataFound
			}
			return nil, err
		}

		IFSCDtos = append(IFSCDtos, *IFSCDto)

	}

	if len(IFSCDtos) == 0 {
		return nil, constants.ErrNoDataFound
	}

	return IFSCDtos, nil

}
