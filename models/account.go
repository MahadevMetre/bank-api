package models

import (
	"database/sql"
	"encoding/json"
	jsonEncoder "encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"bankapi/config"
	"bankapi/constants"
	"bankapi/requests"
	"bankapi/responses"

	"bitbucket.org/paydoh/paydoh-commons/types"
	"github.com/google/uuid"
)

type Account struct {
	Id                   uuid.UUID            `json:"id"`
	UserId               string               `json:"user_id"`
	AccountNumber        string               `json:"account_number"`
	ApplicationID        string               `json:"application_id"`
	CustomerId           string               `json:"customer_id"`
	UpiId                sql.NullString       `json:"upi_id"`
	MotherMaidenName     types.NullableString `json:"mother_maiden_name"`
	IsAddrSameAsAadhaar  bool                 `json:"is_addr_same_as_aadhaar"`
	CommunicationAddress types.NullableString `json:"communication_address,omitempty"`
	CreatedAt            time.Time            `json:"created_at"`
	UpdatedAt            time.Time            `json:"updated_at"`
}

type UserPersonalInformationAccountData struct {
	Id            int64     `json:"id"`
	UserId        string    `json:"user_id"`
	MobileNumber  string    `json:"mobile_number"`
	AccountNumber string    `json:"account_number"`
	CustomerId    string    `json:"customer_id"`
	FirstName     string    `json:"first_name"`
	MiddleName    string    `json:"middle_name"`
	LastName      string    `json:"last_name"`
	Email         string    `json:"email"`
	Gender        string    `json:"gender"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type GetAccountDetailsResponseV2 struct {
	Id                      uuid.UUID              `json:"id"`
	UserId                  string                 `json:"user_id"`
	AccountNumber           sql.NullString         `json:"account_number"`
	Status                  sql.NullString         `json:"status"`
	CustomerId              sql.NullString         `json:"customer_id"`
	UpiID                   sql.NullString         `json:"upi_id"`
	FirstName               string                 `json:"first_name"`
	MiddleName              string                 `json:"middle_name"`
	LastName                string                 `json:"last_name"`
	IsEmailVerified         bool                   `json:"is_email_verified"`
	IsAccountDetailMailSent bool                   `json:"is_account_detail_email_sent"`
	ProfilePic              jsonEncoder.RawMessage `json:"profile_pic"`
	Email                   string                 `json:"email"`
	Gender                  string                 `json:"gender"`
	MaritalStatus           types.NullableString   `json:"-"`
	ProfessionCode          types.NullableString   `json:"-"`
	CustomerEducation       types.NullableString   `json:"-"`
	AnnualTurnOver          types.NullableString   `json:"-"`
	MotherMaidenName        types.NullableString   `json:"-"`
	IsAddrSameAsAdhaar      bool                   `json:"-"`
}

type GetAccountDetailsResponse struct {
	Id            uuid.UUID      `json:"id"`
	UserId        string         `json:"user_id"`
	AccountNumber string         `json:"account_number"`
	Status        sql.NullString `json:"status"`
	CustomerId    string         `json:"customer_id"`
	FirstName     string         `json:"first_name"`
	MiddleName    string         `json:"middle_name"`
	LastName      string         `json:"last_name"`
}

type AccountDataUpdate struct {
	AccountNumber        string                         `json:"account_number,omitempty"`
	CustomerId           string                         `json:"customer_id"`
	ServiceName          string                         `json:"service_name"`
	ApplicationId        string                         `json:"application_id"`
	SourcedBy            string                         `json:"sourced_by"`
	ProductType          string                         `json:"product_type"`
	ApplicantID          string                         `json:"applicant_id"`
	CallbackName         string                         `json:"callback_name"`
	Status               string                         `json:"status"`
	UpiId                string                         `json:"upi_id"`
	IsActive             bool                           `json:"is_active"`
	IsAddrSameAsAdhaar   bool                           `json:"is_addr_same_as_aadhaar"`
	CommunicationAddress *requests.CommunicationAddress `json:"communication_address,omitempty"`
	MaritalStatus        string                         `json:"marital_status"`
	ProfessionCode       string                         `json:"profession_code"`
	MotherMaidenName     string                         `json:"mother_maiden_name"`
	AnnualTurnOver       string                         `json:"annual_turn_over"`
	CustomerEducation    string                         `json:"customer_education"`
}

func NewAccount() *Account {
	return &Account{}
}

func NewUserPersonalInformationAccountData() *UserPersonalInformationAccountData {
	return &UserPersonalInformationAccountData{}
}

func (account *Account) Bind(cbcStatus *responses.CbsStatus) error {
	accountNumber := strconv.Itoa(int(cbcStatus.AccountNo))
	account.UserId = cbcStatus.ApplicantId
	account.AccountNumber = accountNumber
	account.CustomerId = cbcStatus.CustomerID

	return nil
}

func (account *Account) Marshal() ([]byte, error) {
	return json.Marshal(account)
}

func (account *Account) Unmarshal(data []byte) error {
	return json.Unmarshal(data, account)
}

func (account *UserPersonalInformationAccountData) Marshal() ([]byte, error) {
	return json.Marshal(account)
}

func (account *UserPersonalInformationAccountData) Unmarshal(data []byte) error {
	return json.Unmarshal(data, account)
}

func InsertNewAccount(db *sql.DB, userId, applicationId, serviceName string, communicationAddress *requests.CommunicationAddress, isAddrSameAsAdhaar bool, motherMaidenName, annualTurnOver, educationQualification, professionCode, maritalStatus string) error {
	var err error
	var jsonRawMessage []byte

	if communicationAddress != nil {
		jsonRawMessage, err = json.Marshal(communicationAddress)
		if err != nil {
			return err
		}
		_, err = db.Exec("INSERT INTO account_data (user_id, application_id, service_name, communication_address, is_addr_same_as_aadhaar, mother_maiden_name, annual_turn_over, education_qualification, profession_code, marital_status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", userId, applicationId, serviceName, jsonRawMessage, isAddrSameAsAdhaar, motherMaidenName, annualTurnOver, educationQualification, professionCode, maritalStatus)
	} else {
		_, err = db.Exec("INSERT INTO account_data (user_id, application_id, service_name, is_addr_same_as_aadhaar, mother_maiden_name, annual_turn_over, education_qualification, profession_code, marital_status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)", userId, applicationId, serviceName, isAddrSameAsAdhaar, motherMaidenName, annualTurnOver, educationQualification, professionCode, maritalStatus)
	}

	return err
}

func GetAccountDataByUserId(db *sql.DB, userId string) (*Account, error) {
	accountData := NewAccount()
	row := db.QueryRow(
		"SELECT id, user_id, account_number, customer_id, upi_id, created_at, updated_at FROM account_data WHERE user_id=$1",
		userId,
	)

	if err := row.Scan(
		&accountData.Id,
		&accountData.UserId,
		&accountData.AccountNumber,
		&accountData.CustomerId,
		&accountData.UpiId,
		&accountData.CreatedAt,
		&accountData.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, constants.ErrNoDataFound
		}

		return nil, err
	}

	return accountData, nil
}

func GetAccountDataByUserIdV2(userId string) (*Account, error) {
	accountData := NewAccount()
	row := config.GetDB().QueryRow(
		"SELECT id, user_id, account_number, customer_id, upi_id, mother_maiden_name, is_addr_same_as_aadhaar, communication_address, created_at, updated_at FROM account_data WHERE user_id=$1",
		userId,
	)

	if err := row.Scan(
		&accountData.Id,
		&accountData.UserId,
		&accountData.AccountNumber,
		&accountData.CustomerId,
		&accountData.UpiId,
		&accountData.MotherMaidenName,
		&accountData.IsAddrSameAsAadhaar,
		&accountData.CommunicationAddress,
		&accountData.CreatedAt,
		&accountData.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, constants.ErrNoDataFound
		}

		return nil, err
	}

	if accountData.CommunicationAddress.String == "" {
		accountData.CommunicationAddress = types.NewNullableString(nil)
	} else {
		accountData.CommunicationAddress = types.FromString(accountData.CommunicationAddress.String)
	}

	return accountData, nil
}

// get account data by application id
func GetAccountDataByApplicationId(db *sql.DB, applicationId string) (*Account, error) {
	accountData := NewAccount()
	row := db.QueryRow(
		`SELECT
			id,
			user_id,
			application_id
			FROM account_data
			WHERE application_id=$1`,
		applicationId,
	)

	if err := row.Scan(
		&accountData.Id,
		&accountData.UserId,
		&accountData.ApplicationID,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrUserNotFound
		}
		return nil, err
	}

	return accountData, nil
}

func FindUserPersonalInformationAccountDataByMobileNumber(db *sql.DB, mobileNumber string) (*UserPersonalInformationAccountData, error) {
	userPersonalInformationAccountData := NewUserPersonalInformationAccountData()

	row := db.QueryRow("SELECT id, user_id, mobile_number, account_number, customer_id, first_name, middle_name, last_name, email, gender, created_at, updated_at FROM user_personal_information WHERE mobile_number = $1", mobileNumber)

	if err := row.Scan(
		&userPersonalInformationAccountData.Id,
		&userPersonalInformationAccountData.UserId,
		&userPersonalInformationAccountData.MobileNumber,
		&userPersonalInformationAccountData.AccountNumber,
		&userPersonalInformationAccountData.CustomerId,
		&userPersonalInformationAccountData.FirstName,
		&userPersonalInformationAccountData.MiddleName,
		&userPersonalInformationAccountData.LastName,
		&userPersonalInformationAccountData.Email,
		&userPersonalInformationAccountData.Gender,
		&userPersonalInformationAccountData.CreatedAt,
		&userPersonalInformationAccountData.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrNoDataFound
		}
		return nil, err
	}

	return userPersonalInformationAccountData, nil
}

// func InsertNewAccount(db *sql.DB, userId string, accountNumber string, customerId string, applicationId string) error {
// 	if _, err := db.Exec("INSERT INTO account_data (user_id, account_number, customer_id) VALUES ($1, $2, $3)", userId, accountNumber, customerId); err != nil {
// 		return err
// 	}

// 	return nil
// }

// updata user account
func UpdateAccount(db *sql.DB, updateData *AccountDataUpdate, applicationId string) error {
	var clauses []string
	var params []interface{}
	paramCount := 1

	// Check and add each field to the update
	if updateData.AccountNumber != "" {
		clauses = append(clauses, fmt.Sprintf("account_number = $%d", paramCount))
		params = append(params, updateData.AccountNumber)
		paramCount++
	}
	if updateData.CustomerId != "" {
		clauses = append(clauses, fmt.Sprintf("customer_id = $%d", paramCount))
		params = append(params, updateData.CustomerId)
		paramCount++
	}
	if updateData.ServiceName != "" {
		clauses = append(clauses, fmt.Sprintf("service_name = $%d", paramCount))
		params = append(params, updateData.ServiceName)
		paramCount++
	}

	if updateData.SourcedBy != "" {
		clauses = append(clauses, fmt.Sprintf("sourced_by = $%d", paramCount))
		params = append(params, updateData.SourcedBy)
		paramCount++
	}
	if updateData.ProductType != "" {
		clauses = append(clauses, fmt.Sprintf("product_type = $%d", paramCount))
		params = append(params, updateData.ProductType)
		paramCount++
	}

	if updateData.Status != "" {
		clauses = append(clauses, fmt.Sprintf("status = $%d", paramCount))
		params = append(params, updateData.Status)
		paramCount++
	}

	if updateData.CallbackName != "" {
		clauses = append(clauses, fmt.Sprintf("callback_name = $%d", paramCount))
		params = append(params, updateData.CallbackName)
		paramCount++
	}

	// Check if there are any updates to make
	if len(clauses) == 0 {
		return fmt.Errorf("no updates provided")
	}

	clauses = append(clauses, "updated_at = now()")

	query := fmt.Sprintf("UPDATE account_data SET %s WHERE application_id = $%d", strings.Join(clauses, ", "), paramCount)
	params = append(params, applicationId)

	_, err := db.Exec(query, params...)
	if err != nil {
		return fmt.Errorf("failed to update account_data: %v", err)
	}

	return nil
}

func UpdateAccountByUserId(updateData *AccountDataUpdate, userId string) error {
	var clauses []string
	var params []interface{}
	paramCount := 1

	// Check and add each field to the update
	if updateData.AccountNumber != "" {
		clauses = append(clauses, fmt.Sprintf("account_number = $%d", paramCount))
		params = append(params, updateData.AccountNumber)
		paramCount++
	}
	if updateData.CustomerId != "" {
		clauses = append(clauses, fmt.Sprintf("customer_id = $%d", paramCount))
		params = append(params, updateData.CustomerId)
		paramCount++
	}
	if updateData.ServiceName != "" {
		clauses = append(clauses, fmt.Sprintf("service_name = $%d", paramCount))
		params = append(params, updateData.ServiceName)
		paramCount++
	}

	if updateData.SourcedBy != "" {
		clauses = append(clauses, fmt.Sprintf("sourced_by = $%d", paramCount))
		params = append(params, updateData.SourcedBy)
		paramCount++
	}
	if updateData.ProductType != "" {
		clauses = append(clauses, fmt.Sprintf("product_type = $%d", paramCount))
		params = append(params, updateData.ProductType)
		paramCount++
	}

	if updateData.Status != "" {
		clauses = append(clauses, fmt.Sprintf("status = $%d", paramCount))
		params = append(params, updateData.Status)
		paramCount++
	}

	if updateData.CallbackName != "" {
		clauses = append(clauses, fmt.Sprintf("callback_name = $%d", paramCount))
		params = append(params, updateData.CallbackName)
		paramCount++
	}

	if updateData.UpiId != "" {
		clauses = append(clauses, fmt.Sprintf("upi_id = $%d", paramCount))
		params = append(params, updateData.UpiId)
		paramCount++
	}

	if updateData.ApplicationId != "" {
		clauses = append(clauses, fmt.Sprintf("application_id = $%d", paramCount))
		params = append(params, updateData.ApplicationId)
		paramCount++
	}

	if updateData.IsActive {
		clauses = append(clauses, fmt.Sprintf("is_active = $%d", paramCount))
		params = append(params, updateData.IsActive)
		paramCount++
	}

	if !updateData.IsAddrSameAsAdhaar && updateData.CommunicationAddress != nil {
		jsonRawMessage, err := json.Marshal(updateData.CommunicationAddress)
		if err != nil {
			return err
		}
		clauses = append(clauses, fmt.Sprintf("communication_address = $%d, is_addr_same_as_aadhaar = false", paramCount))
		params = append(params, jsonRawMessage)
		paramCount++
	}

	if updateData.IsAddrSameAsAdhaar {
		clauses = append(clauses, fmt.Sprintf("is_addr_same_as_aadhaar = $%d, communication_address = null", paramCount))
		params = append(params, updateData.IsAddrSameAsAdhaar)
		paramCount++
	}

	if updateData.MotherMaidenName != "" {
		clauses = append(clauses, fmt.Sprintf("mother_maiden_name = $%d", paramCount))
		params = append(params, updateData.MotherMaidenName)
		paramCount++
	}

	if updateData.AnnualTurnOver != "" {
		clauses = append(clauses, fmt.Sprintf("annual_turn_over = $%d", paramCount))
		params = append(params, updateData.AnnualTurnOver)
		paramCount++
	}

	if updateData.CustomerEducation != "" {
		clauses = append(clauses, fmt.Sprintf("education_qualification = $%d", paramCount))
		params = append(params, updateData.CustomerEducation)
		paramCount++
	}

	if updateData.ProfessionCode != "" {
		clauses = append(clauses, fmt.Sprintf("profession_code = $%d", paramCount))
		params = append(params, updateData.ProfessionCode)
		paramCount++
	}

	if updateData.MaritalStatus != "" {
		clauses = append(clauses, fmt.Sprintf("marital_status = $%d", paramCount))
		params = append(params, updateData.MaritalStatus)
		paramCount++
	}

	// Check if there are any updates to make
	if len(clauses) == 0 {
		return fmt.Errorf("no updates provided")
	}

	clauses = append(clauses, "updated_at = now()")

	query := fmt.Sprintf("UPDATE account_data SET %s WHERE user_id = $%d", strings.Join(clauses, ", "), paramCount)
	params = append(params, userId)

	_, err := config.GetDB().Exec(query, params...)
	if err != nil {
		return fmt.Errorf("failed to update account_data: %v", err)
	}

	return nil
}

// check if user account already present in DB
func CheckAccountDataAvailability(db *sql.DB, userID string) error {
	query := `SELECT COUNT(*) FROM account_data WHERE user_id = $1 and LOWER(status) = 'success'`

	var count int
	err := db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return fmt.Errorf("error querying account data: %v", err)
	}

	if count > 0 {
		return fmt.Errorf("account data available for the user ID: %s", userID)
	}

	return nil
}

func IsAccountAlreadyUpdated(db *sql.DB, applicationID string) bool {
	query := `SELECT COUNT(*) FROM account_data WHERE application_id = $1 AND LOWER(status) != 'failure'`

	var count int
	err := db.QueryRow(query, applicationID).Scan(&count)
	if err != nil {
		return false
	}

	if count > 0 {
		return true
	}

	return false
}

func GetAccountDetailsV2(db *sql.DB, userId string) (*GetAccountDetailsResponseV2, error) {
	var accountData GetAccountDetailsResponseV2
	row := db.QueryRow(
		`SELECT
        account_data.id,
        account_data.user_id,
        account_data.account_number,
        account_data.status,
        account_data.customer_id,
        account_data.upi_id,
		personal_information.first_name,
		personal_information.middle_name,
		personal_information.last_name,
		COALESCE(personal_information.profile_pic, '{}'::jsonb),
		personal_information.is_email_verified,
		personal_information.is_account_detail_email_sent,
		personal_information.email,
		personal_information.gender,
		account_data.is_addr_same_as_aadhaar,
		account_data.mother_maiden_name,
		account_data.annual_turn_over,
		account_data.education_qualification,
		account_data.profession_code,
		account_data.marital_status
    FROM
        account_data
    JOIN
		personal_information
    ON
        account_data.user_id = personal_information.user_id
    WHERE
        account_data.user_id=$1`,
		userId,
	)

	if err := row.Scan(
		&accountData.Id,
		&accountData.UserId,
		&accountData.AccountNumber,
		&accountData.Status,
		&accountData.CustomerId,
		&accountData.UpiID,
		&accountData.FirstName,
		&accountData.MiddleName,
		&accountData.LastName,
		&accountData.ProfilePic,
		&accountData.IsEmailVerified,
		&accountData.IsAccountDetailMailSent,
		&accountData.Email,
		&accountData.Gender,
		&accountData.IsAddrSameAsAdhaar,
		&accountData.MotherMaidenName,
		&accountData.AnnualTurnOver,
		&accountData.CustomerEducation,
		&accountData.ProfessionCode,
		&accountData.MaritalStatus,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, constants.ErrNoDataFound
		}

		return nil, err
	}

	return &accountData, nil
}

func GetAccountDetails(db *sql.DB, userId string) (*GetAccountDetailsResponse, error) {
	var accountData GetAccountDetailsResponse
	row := db.QueryRow(
		`SELECT
        account_data.id,
        account_data.user_id,
        account_data.account_number,
        account_data.status,
        account_data.customer_id,
		personal_information.first_name,
		personal_information.middle_name,
		personal_information.last_name
    FROM
        account_data
    JOIN
			personal_information
    ON
        account_data.user_id = personal_information.user_id
    WHERE
        account_data.user_id=$1 AND account_number is NOT NULL`,
		userId,
	)

	if err := row.Scan(
		&accountData.Id,
		&accountData.UserId,
		&accountData.AccountNumber,
		&accountData.Status,
		&accountData.CustomerId,
		&accountData.FirstName,
		&accountData.MiddleName,
		&accountData.LastName,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, constants.ErrNoDataFound
		}

		return nil, err
	}

	return &accountData, nil
}

// check if user is eligible to create account (checking previous steps completd for the user)
func IsEligibleForAccountCreate(db *sql.DB, userId string) error {
	// kyc_consent should be present in DB
	_, err := FindKycConsentByUserId(db, userId)
	if err != nil {
		return errors.New("kyc consent data not found for this user")
	}

	// kyc_update_data status should be 'success', should be present in DB
	_, err = GetKycUpdateDataByUserId(db, userId)
	if err != nil {
		return errors.New("kyc update data not found for this user")
	}

	// if kycUpdateData.Status != "SUCCESS" {
	// 	return errors.New("Kyc update has failed for this account")
	// }

	return nil
}
