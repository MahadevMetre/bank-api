package models

import (
	"bankapi/constants"
	"bankapi/security"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"bitbucket.org/paydoh/paydoh-commons/types"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserData struct {
	Id           uuid.UUID `json:"id"`
	UserId       string    `json:"user_id"`
	MobileNumber string    `json:"mobile_number"`
	ApplicantId  string    `json:"applicant_id"`
	SigningKey   string    `json:"signing_key"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	UserId       string             `bson:"user_id,unique" json:"user_id,unique"`
	Email        string             `bson:"email,unique" json:"email,unique"`
	FirstName    string             `bson:"first_name" json:"first_name"`
	MiddleName   string             `bson:"middle_name" json:"middle_name"`
	LastName     string             `bson:"last_name" json:"last_name"`
	DateOfBirth  string             `bson:"date_of_birth" json:"date_of_birth"`
	MobileNumber string             `bson:"mobile_number" json:"mobile_number"`
	Referral     string             `bson:"referral_code" json:"referral_code"`
	ReferredBy   string             `bson:"referred_by" json:"referred_by"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updated_at"`
	DeletedAt    *time.Time         `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}

func NewUserData() *UserData {
	return &UserData{}
}

func (u *UserData) Create(mobileNumber string) error {
	userId, err := security.GenerateRandomUUID(15)

	if err != nil {
		return err
	}

	passphrase := security.GenerateRandomPassphrase()

	secretSigningKey, err := security.Encrypt([]byte(passphrase), []byte(constants.AesPassPhrase))

	if err != nil {
		return err
	}

	u.UserId = userId
	u.MobileNumber = mobileNumber
	u.SigningKey = string(secretSigningKey)

	return nil
}

// get user data by applicant id
func GetUserDataByApplicantId(db *sql.DB, applicantId string) (*UserData, error) {
	userData := NewUserData()
	row := db.QueryRow(
		`SELECT
			id,
			user_id,
			applicant_id
			FROM user_data
			WHERE applicant_id=$1`,
		applicantId,
	)

	if err := row.Scan(
		&userData.Id,
		&userData.UserId,
		&userData.ApplicantId,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrUserNotFound
		}
		return nil, err
	}

	return userData, nil
}

// get user data by user id
func GetUserDataByUserId(db *sql.DB, userId string) (*UserData, error) {
	userData := NewUserData()
	row := db.QueryRow(
		`SELECT
			id,
			user_id,
			applicant_id,
			mobile_number,
			signing_key,
			created_at,
			updated_at
		FROM user_data
		WHERE user_id=$1`,
		userId,
	)

	if err := row.Scan(
		&userData.Id,
		&userData.UserId,
		&userData.ApplicantId,
		&userData.MobileNumber,
		&userData.SigningKey,
		&userData.CreatedAt,
		&userData.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrUserNotFound
		}
		return nil, err
	}

	return userData, nil
}

func GetUserDataByMobileNumber(db *sql.DB, mobileNumber string) (*UserData, error) {
	userData := NewUserData()
	row := db.QueryRow(
		`SELECT
			id,
			user_id,
			applicant_id,
			signing_key
			FROM user_data
			WHERE mobile_number=$1`,
		mobileNumber,
	)

	if err := row.Scan(
		&userData.Id,
		&userData.UserId,
		&userData.ApplicantId,
		&userData.SigningKey,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrUserNotFound
		}
		return nil, err
	}

	return userData, nil
}

type UserPersonalInformationAndAccountData struct {
	Id                      uuid.UUID            `json:"id"`
	UserId                  string               `json:"user_id"`
	MobileNumber            string               `json:"mobile_number"`
	AccountNumber           string               `json:"account_number"`
	Applicant_id            string               `json:"applicant_id"`
	CustomerId              string               `json:"customer_id"`
	UpiId                   sql.NullString       `json:"upi_id"`
	FirstName               string               `json:"first_name"`
	MiddleName              string               `json:"middle_name"`
	LastName                string               `json:"last_name"`
	Email                   string               `json:"email"`
	Gender                  string               `json:"gender"`
	DateOfBirth             string               `json:"date_of_birth"`
	IsEmailVerified         bool                 `json:"is_email_verified"`
	IsAccountDetailMailSent bool                 `json:"is_account_detail_email_sent"`
	CommunicationAddress    types.NullableString `json:"communication_address"`
	IsAddrSameAsAdhaar      bool                 `json:"is_addr_same_as_aadhaar"`
	ProfilePic              json.RawMessage      `json:"profile_pic,omitempty"`
}

func GetUserAndAccountDetailByUserID(db *sql.DB, user_id string) (*UserPersonalInformationAndAccountData, error) {
	var userinfo UserPersonalInformationAndAccountData
	var profilePicNull sql.NullString

	row := db.QueryRow(`SELECT
        account_data.id,
        account_data.user_id,
        account_data.account_number,
        user_data.applicant_id,
        account_data.customer_id,
		account_data.upi_id,
		personal_information.first_name,
		personal_information.middle_name,
		personal_information.last_name,
		personal_information.gender,
		personal_information.date_of_birth,
		personal_information.email,
		personal_information.is_email_verified,
		personal_information.is_account_detail_email_sent,
		user_data.mobile_number,
		account_data.communication_address,
		account_data.is_addr_same_as_aadhaar,
		personal_information.profile_pic
    FROM
        account_data
    LEFT JOIN
			personal_information
    ON
        account_data.user_id = personal_information.user_id
   left join
   user_data ON
        account_data.user_id = user_data.user_id
    WHERE
        account_data.user_id=$1 AND account_number is NOT NULL`, user_id)

	if err :=
		row.Scan(
			&userinfo.Id,
			&userinfo.UserId,
			&userinfo.AccountNumber,
			&userinfo.Applicant_id,
			&userinfo.CustomerId,
			&userinfo.UpiId,
			&userinfo.FirstName,
			&userinfo.MiddleName,
			&userinfo.LastName,
			&userinfo.Gender,
			&userinfo.DateOfBirth,
			&userinfo.Email,
			&userinfo.IsEmailVerified,
			&userinfo.IsAccountDetailMailSent,
			&userinfo.MobileNumber,
			&userinfo.CommunicationAddress,
			&userinfo.IsAddrSameAsAdhaar,
			&profilePicNull,
		); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, constants.ErrNoDataFound
		}
		return nil, err
	}

	if profilePicNull.Valid {
		userinfo.ProfilePic = json.RawMessage(profilePicNull.String)
	}

	return &userinfo, nil
}

func GetUserAndAccountDetailByMobileNumber(db *sql.DB, mobileNumber string) (*UserPersonalInformationAndAccountData, error) {
	var userinfo UserPersonalInformationAndAccountData

	query := `
		SELECT
			a.id,
			a.user_id,
			a.account_number,
			u.applicant_id,
			a.customer_id,
			a.upi_id,
			p.first_name,
			p.middle_name,
			p.last_name,
			p.gender,
			p.date_of_birth,
			p.email,
			u.mobile_number
		FROM
			account_data a
		LEFT JOIN
			personal_information p
		ON
			a.user_id = p.user_id
		LEFT JOIN
			user_data u
		ON
			a.user_id = u.user_id
		WHERE
			u.mobile_number = $1 AND a.account_number IS NOT NULL`

	row := db.QueryRow(query, mobileNumber)

	if err :=
		row.Scan(
			&userinfo.Id,
			&userinfo.UserId,
			&userinfo.AccountNumber,
			&userinfo.Applicant_id,
			&userinfo.CustomerId,
			&userinfo.UpiId,
			&userinfo.FirstName,
			&userinfo.MiddleName,
			&userinfo.LastName,
			&userinfo.Gender,
			&userinfo.DateOfBirth,
			&userinfo.Email,
			&userinfo.MobileNumber,
		); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, constants.ErrNoDataFound
		}
		return nil, err
	}

	return &userinfo, nil
}
