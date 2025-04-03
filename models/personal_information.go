package models

import (
	"bankapi/constants"
	"bankapi/requests"
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"bitbucket.org/paydoh/paydoh-commons/types"
	"github.com/google/uuid"
)

type PersonalInformation struct {
	Id                      uuid.UUID            `json:"id"`
	UserId                  string               `json:"user_id"`
	FirstName               string               `json:"first_name"`
	MiddleName              string               `json:"middle_name"`
	LastName                string               `json:"last_name"`
	Gender                  string               `json:"gender"`
	Email                   string               `json:"email"`
	PinCode                 types.NullableString `json:"pin_code"`
	DateOfBirth             string               `json:"date_of_birth"`
	IsEmailVerified         bool                 `json:"is_email_verified"`
	IsAccountDetailMailSent bool                 `json:"is_account_detail_email_sent"`
}

func NewPersonalInformation() *PersonalInformation {
	return &PersonalInformation{}
}

func (p *PersonalInformation) Bind(request *requests.PersonalInformationRequest, userId string) {
	p.UserId = userId
	p.FirstName = request.FirstName
	p.MiddleName = request.MiddleName
	p.LastName = request.LastName
	p.Gender = request.Gender
	p.Email = request.Email
	p.DateOfBirth = request.DateOfBirth
}

func (p *PersonalInformation) BindUpdate(request *requests.PersonalInformationRequest, existingPersonalInformation *PersonalInformation) error {
	if request.FirstName != existingPersonalInformation.FirstName {
		p.FirstName = request.FirstName
	}

	if request.MiddleName != existingPersonalInformation.MiddleName {
		p.MiddleName = request.MiddleName
	}

	if request.LastName != existingPersonalInformation.LastName {
		p.LastName = request.LastName
	}

	if request.Gender != existingPersonalInformation.Gender {
		p.Gender = request.Gender
	}

	if request.Email != existingPersonalInformation.Email {
		p.Email = request.Email
	}

	if request.DateOfBirth != existingPersonalInformation.DateOfBirth {
		p.DateOfBirth = request.DateOfBirth
	}

	if p.FirstName == existingPersonalInformation.FirstName &&
		p.MiddleName == existingPersonalInformation.MiddleName &&
		p.LastName == existingPersonalInformation.LastName &&
		p.Gender == existingPersonalInformation.Gender &&
		p.Email == existingPersonalInformation.Email &&
		p.DateOfBirth == existingPersonalInformation.DateOfBirth {
		return errors.New("no data to update")
	}

	return nil
}

func (p *PersonalInformation) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

func (p *PersonalInformation) UnMarshal(data []byte) error {
	return json.Unmarshal(data, p)
}

func InsertPersonalInformation(db *sql.DB, personalInformation *PersonalInformation) error {
	_, err := db.Exec(
		"INSERT INTO personal_information (user_id, first_name, middle_name, last_name, gender, email, date_of_birth) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		personalInformation.UserId,
		personalInformation.FirstName,
		personalInformation.MiddleName,
		personalInformation.LastName,
		personalInformation.Gender,
		personalInformation.Email,
		personalInformation.DateOfBirth,
	)

	if err != nil {
		return err
	}

	return nil
}

func InsertEmailInPersonalInfo(db *sql.DB, personalInformation *PersonalInformation) error {
	_, err := db.Exec(
		"INSERT INTO personal_information (user_id, first_name, middle_name, last_name, gender, email, date_of_birth,is_email_verified) VALUES ($1, $2, $3, $4, $5, $6, $7,$8)",
		personalInformation.UserId,
		personalInformation.FirstName,
		personalInformation.MiddleName,
		personalInformation.LastName,
		personalInformation.Gender,
		personalInformation.Email,
		personalInformation.DateOfBirth,
		true,
	)

	if err != nil {
		return err
	}

	return nil
}

func UpdatePersonalInformation(db *sql.DB, update *PersonalInformation, userId string) error {
	var clauses []string
	var params []interface{}

	if update.FirstName != "" {
		clauses = append(clauses, "first_name = $"+strconv.Itoa(len(clauses)+1))
		params = append(params, update.FirstName)
	}

	if update.MiddleName != "" {
		clauses = append(clauses, "middle_name = $"+strconv.Itoa(len(clauses)+1))
		params = append(params, update.MiddleName)
	}

	if update.LastName != "" {
		clauses = append(clauses, "last_name = $"+strconv.Itoa(len(clauses)+1))
		params = append(params, update.LastName)
	}

	if update.Gender != "" {
		clauses = append(clauses, "gender = $"+strconv.Itoa(len(clauses)+1))
		params = append(params, update.Gender)
	}

	if update.Email != "" {
		clauses = append(clauses, "email = $"+strconv.Itoa(len(clauses)+1))
		params = append(params, update.Email)
	}

	if update.PinCode.Valid {
		clauses = append(clauses, "pin_code = $"+strconv.Itoa(len(clauses)+1))
		params = append(params, update.PinCode.String)
	}

	if update.DateOfBirth != "" {
		clauses = append(clauses, "date_of_birth = $"+strconv.Itoa(len(clauses)+1))
		params = append(params, update.DateOfBirth)
	}

	if update.IsEmailVerified {
		clauses = append(clauses, "is_email_verified = $"+strconv.Itoa(len(clauses)+1))
		params = append(params, true)
	}

	if update.IsAccountDetailMailSent {
		clauses = append(clauses, "is_account_detail_email_sent = $"+strconv.Itoa(len(clauses)+1))
		params = append(params, true)
	}

	if len(clauses) > 0 {

		clauses = append(clauses, "updated_at = NOW()")
		query := "UPDATE personal_information SET " + strings.Join(clauses, ", ") + " WHERE user_id =$" + strconv.Itoa(len(params)+1)
		params = append(params, userId)

		if _, err := db.Exec(query, params...); err != nil {
			return err
		}

		return nil
	}

	return errors.New("nothing to update")
}

func GetEmailVerificationStatus(db *sql.DB, userId string) (*PersonalInformation, error) {
	personalInformation := NewPersonalInformation()
	row := db.QueryRow("SELECT id, user_id, email, is_email_verified FROM personal_information WHERE user_id = $1 AND is_email_verified=true", userId)

	if err := row.Scan(
		&personalInformation.Id,
		&personalInformation.UserId,
		&personalInformation.Email,
		&personalInformation.IsEmailVerified,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrNoDataFound
		}

		return nil, err
	}

	return personalInformation, nil
}

func GetPersonalInformation(db *sql.DB, userId string) (*PersonalInformation, error) {
	personalInformation := NewPersonalInformation()
	row := db.QueryRow("SELECT id, user_id, first_name, middle_name, last_name, gender, email, pin_code, date_of_birth, is_email_verified,is_account_detail_email_sent FROM personal_information WHERE user_id = $1", userId)

	if err := row.Scan(
		&personalInformation.Id,
		&personalInformation.UserId,
		&personalInformation.FirstName,
		&personalInformation.MiddleName,
		&personalInformation.LastName,
		&personalInformation.Gender,
		&personalInformation.Email,
		&personalInformation.PinCode,
		&personalInformation.DateOfBirth,
		&personalInformation.IsEmailVerified,
		&personalInformation.IsAccountDetailMailSent,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrNoDataFound
		}

		return nil, err
	}

	if personalInformation.PinCode.String == "" {
		personalInformation.PinCode = types.NewNullableString(nil)
	} else {
		personalInformation.PinCode = types.FromString(personalInformation.PinCode.String)
	}

	return personalInformation, nil
}

func AddAccountDetailSentOnMailStatus(db *sql.DB, userId string) error {
	query := "UPDATE personal_information SET updated_at = NOW(), is_account_detail_email_sent = $1 WHERE user_id = $2"

	if _, err := db.Exec(query, true, userId); err != nil {
		return err
	}

	return nil
}
