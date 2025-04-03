package models

import (
	"database/sql"
	"time"
)

type CustomerOnboardingUserDetails struct {
	ID          int            `json:"id"`
	CustomerID  string         `json:"customer_id"`
	FirstName   sql.NullString `json:"first_name"`
	LastName    sql.NullString `json:"last_name"`
	MiddleName  sql.NullString `json:"middle_name"`
	Gender      sql.NullString `json:"gender"`
	EmailId     sql.NullString `json:"email_id"`
	DateOfBirth sql.NullString `json:"dob"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type CustomerDemographicDataDetails struct {
	Id          int            `json:"id"`
	CustomerID  string         `json:"customer_id"`
	HouseNumber sql.NullString `json:"house_number"`
	StreetName  sql.NullString `json:"street_name"`
	Locality    sql.NullString `json:"locality"`
	Landmark    sql.NullString `json:"landmark"`
	City        sql.NullString `json:"city"`
	State       sql.NullString `json:"state"`
	Country     sql.NullString `json:"country"`
	ZipCode     sql.NullString `json:"zip_code"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
