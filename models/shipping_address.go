package models

import (
	"bankapi/constants"
	"bankapi/requests"
	"database/sql"
	"strconv"
	"strings"
	"time"
)

type ShippingAddress struct {
	Id           int64          `json:"id"`
	UserId       string         `json:"user_id"`
	DocumentType string         `json:"document_type"`
	Document     string         `json:"document"`
	AddressLine1 string         `json:"address_line_1"`
	StreetName   string         `json:"street_name"`
	Locality     string         `json:"locality"`
	Landmark     sql.NullString `json:"landmark"`
	City         string         `json:"city"`
	State        string         `json:"state"`
	PinCode      string         `json:"pin_code"`
	Country      string         `json:"country"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

func NewShippingAddress() *ShippingAddress {
	return &ShippingAddress{}
}

func (s *ShippingAddress) Bind(request *requests.AddShippingAddress, userId string) {
	s.UserId = userId
	s.DocumentType = request.DocumentType
	s.Document = request.Document
	s.AddressLine1 = request.AddressLine1
	s.StreetName = request.StreetName
	s.Locality = request.Locality
	s.Landmark = sql.NullString{String: request.Landmark, Valid: request.Landmark != ""}
	s.City = request.City
	s.State = request.State
	s.PinCode = request.PinCode
	s.Country = request.Country
}

func InsertShippingAddress(db *sql.DB, s *ShippingAddress) error {
	if _, err := db.Exec("INSERT INTO shipping_address (user_id, document_type, document, address_line_1, street_name, locality, landmark, city, state, pin_code, country) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		s.UserId,
		s.DocumentType,
		s.Document,
		s.AddressLine1,
		s.StreetName,
		s.Locality,
		s.Landmark.String,
		s.City,
		s.State,
		s.PinCode,
		s.Country,
	); err != nil {
		return err
	}

	return nil
}

func FindShippingAddressByUserId(db *sql.DB, userId string) (*ShippingAddress, error) {
	s := NewShippingAddress()
	err := db.QueryRow("SELECT * FROM shipping_address WHERE user_id = $1", userId).Scan(
		&s.Id,
		&s.UserId,
		&s.DocumentType,
		&s.Document,
		&s.AddressLine1,
		&s.StreetName,
		&s.Locality,
		&s.Landmark,
		&s.City,
		&s.State,
		&s.PinCode,
		&s.Country,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrNoDataFound
		}

		return nil, err
	}
	return s, nil
}

func UpdateShippingAddressByUserId(db *sql.DB, request *requests.UpdateShippingAddress, userId string) error {
	var clause []string
	var parameter []interface{}

	if request.DocumentType != "" {
		clause = append(clause, "document_type = $"+strconv.Itoa(len(parameter)+1))
		parameter = append(parameter, request.DocumentType)
	}

	if request.Document != "" {
		clause = append(clause, "document = $"+strconv.Itoa(len(parameter)+1))
		parameter = append(parameter, request.Document)
	}

	if request.AddressLine1 != "" {
		clause = append(clause, "address_line_1 = $"+strconv.Itoa(len(parameter)+1))
		parameter = append(parameter, request.AddressLine1)
	}

	if request.StreetName != "" {
		clause = append(clause, "street_name = $"+strconv.Itoa(len(parameter)+1))
		parameter = append(parameter, request.StreetName)
	}

	if request.Locality != "" {
		clause = append(clause, "locality = $"+strconv.Itoa(len(parameter)+1))
		parameter = append(parameter, request.Locality)
	}

	if request.Landmark != "" {
		clause = append(clause, "landmark = $"+strconv.Itoa(len(parameter)+1))
		parameter = append(parameter, request.Landmark)
	}

	if request.City != "" {
		clause = append(clause, "city = $"+strconv.Itoa(len(parameter)+1))
		parameter = append(parameter, request.City)
	}

	if request.State != "" {
		clause = append(clause, "state = $"+strconv.Itoa(len(parameter)+1))
		parameter = append(parameter, request.State)
	}

	if request.PinCode != "" {
		clause = append(clause, "pin_code = $"+strconv.Itoa(len(parameter)+1))
		parameter = append(parameter, request.PinCode)
	}

	if request.Country != "" {
		clause = append(clause, "country = $"+strconv.Itoa(len(parameter)+1))
		parameter = append(parameter, request.Country)
	}

	if len(clause) > 0 {
		clause = append(clause, "updated_at = now()")
		_, err := db.Exec("UPDATE shipping_address SET "+strings.Join(clause, ", ")+" WHERE user_id ="+userId, parameter...)
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}
