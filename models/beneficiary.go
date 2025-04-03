package models

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"bankapi/constants"
	"bankapi/responses"

	"bitbucket.org/paydoh/paydoh-commons/types"
	"github.com/google/uuid"
)

// For send beneficiary nickname to mobile team we are using this struct to show proper benficiary nickname
// Scenario 13 If user select inactive beneficiary , data should be pre populated
type BeneficiaryDetail struct {
	BenfMob          string `json:"beneficiary_mobile_number"`
	BenfName         string `json:"beneficiary_name"`
	BenfId           string `json:"beneficiary_id"`
	BenfNickName     string `json:"beneficiary_nickname"`
	BenfAcctNo       string `json:"beneficiary_account_number"`
	BenfIFSC         string `json:"beneficiary_ifsc"`
	BenfAcctType     string `json:"beneficiary_account_type"`
	PaymentMode      string `json:"payment_mode"`
	BenfStatus       string `json:"beneficiary_status"`
	BenfActivateTime string `json:"beneficiary_activate_time"`
	TxnIdentifier    string `json:"txn_identifier"`
}

// For send same transaction while beneficary operation for the same record we are using this struct TxnIdentifier
// Scenario 9-OTP Recived and OTP verification request  initiated and receive success response
type BeneficiaryDTO struct {
	Id               uuid.UUID            `json:"id"`
	BenfId           string               `json:"beneficiary_id"`
	UserId           string               `json:"user_id"`
	BenfName         string               `json:"beneficiary_name"`
	BenfNickName     types.NullableString `json:"beneficiary_nickname"`
	TxnIdentifier    string               `json:"txn_identifier"`
	BenfMobNo        string               `json:"beneficiary_mobile_number"`
	BenfAccountNo    string               `json:"beneficiary_account_number"`
	BenfIfsc         string               `json:"beneficiary_ifsc"`
	BenfAcctType     string               `json:"beneficiary_account_type"`
	PaymentMode      string               `json:"payment_mode"`
	BenfStatus       string               `json:"beneficiary_status"`
	BenfActivateTime string               `json:"beneficiary_activate_time"`
	ActivatedDtTime  string               `json:"activatedDtTime"`
	IsActive         bool                 `json:"is_active"`
	CreatedAt        time.Time            `json:"created_at"`
	UpdatedAt        time.Time            `json:"updated_at"`
}

func NewBeneficiaryDetail() *BeneficiaryDetail {
	return &BeneficiaryDetail{}
}

func NewBeneficiaryDTO() *BeneficiaryDTO {
	return &BeneficiaryDTO{}
}

// Scenario 13: If the user selects an inactive beneficiary, the data should be pre-populated.
// Here, we are using two different data sources:
// 1. Bank response - to fetch the latest beneficiary details from the bank. BankBeneficiaryFetch func is used for this.
// 2. Database response - to retrieve any existing beneficiary records stored in our system. DbBeneficiaryFetch func is used for this.
// The combined data is then sent to the mobile team for display.
func (b *BeneficiaryDetail) Bind(response *responses.FetchBeneficiaryResponse, data []BeneficiaryDTO) ([]BeneficiaryDetail, error) {
	bankBeneficiaries, err := b.BankBeneficiaryFetch(response)
	if err != nil {
		return nil, err
	}

	inactiveBeneficiaries, err := b.DbBeneficiaryFetch(data)
	if err != nil {
		return nil, err
	}

	for key, ben := range inactiveBeneficiaries {
		if _, exists := bankBeneficiaries[key]; !exists {
			bankBeneficiaries[key] = ben
		}
	}

	beneficiaries := make([]BeneficiaryDetail, 0, len(bankBeneficiaries))
	for _, ben := range bankBeneficiaries {
		beneficiaries = append(beneficiaries, *ben)
	}

	if len(beneficiaries) == 0 {
		return nil, constants.ErrNoDataFound
	}

	return beneficiaries, nil
}

// Senerio 12-OTP Recived and OTP verification request initiated and receive fail response back Office Time Out
// Senerio 12-OTP Recived and OTP verification request initiated and receive fail response Technical Error
// Here, we are retrieving the beneficiary details from the bank's response
// to ensure we fetch the most recent and accurate data.
func (b *BeneficiaryDetail) BankBeneficiaryFetch(response *responses.FetchBeneficiaryResponse) (map[string]*BeneficiaryDetail, error) {
	beneficiaryMap := make(map[string]*BeneficiaryDetail)

	for _, beneficiary := range response.BeneficiaryDetails {
		ben := NewBeneficiaryDetail()

		ben.BenfAcctNo = maskAccountNumber(beneficiary.BenfAcctNo)
		ben.BenfActivateTime = beneficiary.BenfActivateTime
		ben.BenfAcctType = beneficiary.BenfAcctType
		ben.BenfId = beneficiary.BenfID
		ben.BenfIFSC = beneficiary.BenfIFSC
		ben.BenfName = beneficiary.BenfName
		ben.BenfStatus = beneficiary.BenfStatus
		ben.BenfMob = beneficiary.BenfMob
		ben.PaymentMode = beneficiary.PaymentMode
		ben.BenfNickName = beneficiary.BenfID
		ben.TxnIdentifier = ""

		key := ben.BenfAcctNo + "_" + ben.BenfMob + "_" + ben.BenfName + "_" + ben.BenfIFSC + "_" + ben.BenfId
		beneficiaryMap[key] = ben
	}

	return beneficiaryMap, nil
}

// Scenario 13: If the user selects an inactive beneficiary, the data should be pre-populated.
// Here, we are fetching the data only from the database (DBFetch)
// since inactive beneficiaries are stored in our system.
// The retrieved data is then sent to the mobile team for display.
func (b *BeneficiaryDetail) DbBeneficiaryFetch(data []BeneficiaryDTO) (map[string]*BeneficiaryDetail, error) {
	beneficiaryMap := make(map[string]*BeneficiaryDetail)

	for _, beneficiary := range data {
		if !beneficiary.IsActive {
			ben := NewBeneficiaryDetail()
			ben.BenfAcctNo = maskAccountNumber(beneficiary.BenfAccountNo)
			ben.BenfActivateTime = beneficiary.ActivatedDtTime
			ben.BenfAcctType = beneficiary.BenfAcctType
			ben.BenfId = beneficiary.BenfId
			ben.BenfIFSC = beneficiary.BenfIfsc
			ben.BenfName = beneficiary.BenfName
			ben.BenfStatus = "inactive"
			ben.BenfMob = beneficiary.BenfMobNo
			ben.PaymentMode = beneficiary.PaymentMode
			ben.BenfNickName = beneficiary.BenfNickName.String
			ben.TxnIdentifier = beneficiary.TxnIdentifier

			key := ben.BenfAcctNo + "_" + ben.BenfMob + "_" + ben.BenfName + "_" + ben.BenfIFSC + "_" + ben.BenfId
			beneficiaryMap[key] = ben
		}
	}

	return beneficiaryMap, nil
}

// Helper FetchBind function is used to wrap the data in to the structure.
// i'm using this one in the VerifyAndConfirmBeneficiary function
func (b *BeneficiaryDetail) FetchBind(response *responses.FetchBeneficiaryResponse) ([]BeneficiaryDetail, error) {

	beneficiaries := make([]BeneficiaryDetail, 0)

	for _, beneficiary := range response.BeneficiaryDetails {
		ben := NewBeneficiaryDetail()
		ben.BenfAcctNo = beneficiary.BenfAcctNo
		ben.BenfActivateTime = beneficiary.BenfActivateTime
		ben.BenfAcctType = beneficiary.BenfAcctType
		ben.BenfId = beneficiary.BenfID
		ben.BenfIFSC = beneficiary.BenfIFSC
		ben.BenfName = beneficiary.BenfName
		ben.BenfStatus = beneficiary.BenfStatus
		ben.BenfMob = beneficiary.BenfMob
		ben.PaymentMode = beneficiary.PaymentMode
		beneficiaries = append(beneficiaries, *ben)
	}

	if len(beneficiaries) == 0 {
		return nil, constants.ErrNoDataFound
	}

	return beneficiaries, nil
}

func maskAccountNumber(accountNo string) string {
	if len(accountNo) <= 4 {
		return accountNo
	}
	return strings.Repeat("X", len(accountNo)-4) + accountNo[len(accountNo)-4:]
}

// For send same transaction while beneficary operation for the same record we are using this struct TxnIdentifier
// Senerio 9-OTP Recived and OTP verification request  initiated and receive success response
func (dto *BeneficiaryDTO) BindData(
	beneficiaryId,
	userId,
	beneficiaryName,
	beneficiaryNickName,
	beneficiaryMobileNumber,
	beneficiaryAccount,
	beneficiaryIFSC,
	beneficiaryAccountType,
	paymentMode, activatedDtTime, txnId string) error {

	dto.BenfId = beneficiaryId
	dto.UserId = userId
	dto.BenfName = beneficiaryName

	if beneficiaryNickName == "" {
		dto.BenfNickName = types.NewNullableString(nil)
	} else {
		dto.BenfNickName = types.FromString(beneficiaryNickName)
	}

	dto.BenfMobNo = beneficiaryMobileNumber
	dto.BenfAccountNo = beneficiaryAccount
	dto.BenfIfsc = beneficiaryIFSC
	dto.BenfAcctType = beneficiaryAccountType
	dto.PaymentMode = paymentMode
	dto.ActivatedDtTime = activatedDtTime
	dto.TxnIdentifier = txnId
	return nil
}

func (b *BeneficiaryDTO) Marshal() ([]byte, error) {
	return json.Marshal(b)
}

func (b *BeneficiaryDTO) Unmarshal(data []byte) error {
	return json.Unmarshal(data, b)
}

// For send same transaction while beneficary operation for the same record we are using this struct TxnIdentifier
// Senerio 9-OTP Recived and OTP verification request  initiated and receive success response
func InsertBeneficiaryDetails(db *sql.DB, data *BeneficiaryDTO) (string, error) {
	var id string
	maskedAccount := strings.Repeat("X", len(data.BenfAccountNo)-4) + data.BenfAccountNo[len(data.BenfAccountNo)-4:]

	err := db.QueryRow(
		"INSERT INTO beneficiaries (benf_id, user_id, benf_name, benf_nickname, benf_mobile_number, benf_ifsc, benf_acct_type, payment_mode, benf_activated_time, is_active, benf_account, txn_identifier) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id",
		data.BenfId,
		data.UserId,
		data.BenfName,
		data.BenfNickName,
		data.BenfMobNo,
		data.BenfIfsc,
		data.BenfAcctType,
		data.PaymentMode,
		data.ActivatedDtTime,
		data.IsActive,
		maskedAccount,
		data.TxnIdentifier,
	).Scan(&id)

	if err != nil {
		return "", err
	}

	return id, nil
}

// Scenario 2 - If Beneficiary data sent to KVB and success res received and left the half way without completion, retried with new request with change in value
// after success response data update in DB
func UpdateBeneficiary(db *sql.DB, data *BeneficiaryDTO) error {
	var updates []string
	var args []interface{}
	argPosition := 1

	if data.BenfName != "" {
		updates = append(updates, fmt.Sprintf("benf_name = $%d", argPosition))
		args = append(args, data.BenfName)
		argPosition++
	}

	if data.BenfId != "" {
		updates = append(updates, fmt.Sprintf("benf_id = $%d", argPosition))
		args = append(args, data.BenfId)
		argPosition++
	}

	if data.BenfNickName.String != "" {
		updates = append(updates, fmt.Sprintf("benf_nickname = $%d", argPosition))
		args = append(args, data.BenfNickName)
		argPosition++
	}
	if data.BenfMobNo != "" {
		updates = append(updates, fmt.Sprintf("benf_mobile_number = $%d", argPosition))
		args = append(args, data.BenfMobNo)
		argPosition++
	}
	if data.BenfIfsc != "" {
		updates = append(updates, fmt.Sprintf("benf_ifsc = $%d", argPosition))
		args = append(args, data.BenfIfsc)
		argPosition++
	}
	if data.BenfAcctType != "" {
		updates = append(updates, fmt.Sprintf("benf_acct_type = $%d", argPosition))
		args = append(args, data.BenfAcctType)
		argPosition++
	}
	if data.PaymentMode != "" {
		updates = append(updates, fmt.Sprintf("payment_mode = $%d", argPosition))
		args = append(args, data.PaymentMode)
		argPosition++
	}

	if data.BenfAccountNo != "" {
		updates = append(updates, fmt.Sprintf("benf_account = $%d", argPosition))
		args = append(args, data.BenfAccountNo)
		argPosition++
	}

	if data.ActivatedDtTime != "" {
		updates = append(updates, fmt.Sprintf("benf_activated_time = $%d", argPosition))
		args = append(args, data.ActivatedDtTime)
		argPosition++
	}
	if data.IsActive {
		updates = append(updates, fmt.Sprintf("is_active = $%d", argPosition))
		args = append(args, data.IsActive)
		argPosition++
	}

	updates = append(updates, "updated_at = now()")

	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf("UPDATE beneficiaries SET %s WHERE benf_id = $%d AND user_id = $%d",
		strings.Join(updates, ", "),
		argPosition, argPosition+1)

	args = append(args, data.BenfId, data.UserId)

	result, err := db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error updating beneficiary: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("beneficiary with id %s not found", data.BenfId)
	}

	return nil
}

// This Is Unused Code
// func UpdateBeneficiaryStatus(beneficiaryId, activationTime string) error {
// 	db := config.GetDB()
// 	if _, err := db.Exec("UPDATE beneficiaries SET benf_activated_time = $1, updated_at = now() WHERE benf_id = $2", activationTime, beneficiaryId); err != nil {
// 		return err
// 	}

// 	return nil
// }

// For send beneficiary's name to mobile team we are using this struct to show proper benficiary nick name
// Scenario 13 If user select inactive beneficiary , data should be pre populated
func FindBeneficiariesByUserId(db *sql.DB, userId string) ([]BeneficiaryDTO, error) {

	beneficiaries := make([]BeneficiaryDTO, 0)

	rows, err := db.Query("SELECT id, benf_id, user_id, benf_name, benf_nickname, benf_mobile_number, benf_ifsc, benf_acct_type, payment_mode, benf_activated_time, created_at, updated_at, is_active, benf_account, txn_identifier FROM beneficiaries WHERE user_id = $1", userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		ben := NewBeneficiaryDTO()
		err := rows.Scan(
			&ben.Id,
			&ben.BenfId,
			&ben.UserId,
			&ben.BenfName,
			&ben.BenfNickName,
			&ben.BenfMobNo,
			&ben.BenfIfsc,
			&ben.BenfAcctType,
			&ben.PaymentMode,
			&ben.BenfActivateTime,
			&ben.CreatedAt,
			&ben.UpdatedAt,
			&ben.IsActive,
			&ben.BenfAccountNo,
			&ben.TxnIdentifier,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, constants.ErrNoDataFound
			}
			return nil, err
		}
		beneficiaries = append(beneficiaries, *ben)
	}

	return beneficiaries, nil
}

func BeneficiaryByBenfId(db *sql.DB, benfId string) (*BeneficiaryDTO, error) {
	ben := NewBeneficiaryDTO()

	row := db.QueryRow("SELECT id, benf_id, user_id, benf_name, benf_nickname, benf_mobile_number, benf_ifsc, benf_acct_type, payment_mode, benf_activated_time, created_at, updated_at, is_active FROM beneficiaries WHERE benf_id = $1", benfId)

	if err := row.Scan(
		&ben.Id,
		&ben.BenfId,
		&ben.UserId,
		&ben.BenfName,
		&ben.BenfNickName,
		&ben.BenfMobNo,
		&ben.BenfIfsc,
		&ben.BenfAcctType,
		&ben.PaymentMode,
		&ben.BenfActivateTime,
		&ben.CreatedAt,
		&ben.UpdatedAt,
		&ben.IsActive,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, constants.ErrNoDataFound
		}

		return nil, err
	}

	return ben, nil
}

func GetBeneficiariesIn(db *sql.DB, benfIds []any) ([]BeneficiaryDTO, error) {
	beneficiaries := make([]BeneficiaryDTO, 0)

	placeHolders := make([]string, len(benfIds))
	for i := range benfIds {
		placeHolders[i] = fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf("SELECT id, benf_id, user_id, benf_name, benf_nickname, benf_mob_no, benf_account_no, benf_ifsc, benf_acct_type, payment_mode, benf_status, benf_activate_time, created_at, updated_at FROM beneficiary_details WHERE benf_id IN (%s)", strings.Join(placeHolders, ","))

	stmt, err := db.Prepare(query)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(benfIds...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		ben := NewBeneficiaryDTO()
		err := rows.Scan(
			&ben.Id,
			&ben.BenfId,
			&ben.UserId,
			&ben.BenfName,
			&ben.BenfNickName,
			&ben.BenfMobNo,
			&ben.BenfAccountNo,
			&ben.BenfIfsc,
			&ben.BenfAcctType,
			&ben.PaymentMode,
			&ben.BenfStatus,
			&ben.BenfActivateTime,
			&ben.CreatedAt,
			&ben.UpdatedAt,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, constants.ErrNoDataFound
			}
			return nil, err
		}
		beneficiaries = append(beneficiaries, *ben)
	}
	return beneficiaries, nil
}

func FindBeneficiaryByNameAndIfscCode(db *sql.DB, userId, name, ifscCode, nickName, mobileNumber, benfAcctNo string) (*BeneficiaryDTO, error) {

	ben := NewBeneficiaryDTO()
	maskedAccount := strings.Repeat("X", len(benfAcctNo)-4) + benfAcctNo[len(benfAcctNo)-4:]
	row := db.QueryRow(
		`SELECT
            id,
            benf_id,
            user_id,
            benf_name,
            benf_nickname,
            benf_mobile_number,
            benf_ifsc,
            benf_acct_type,
            payment_mode,
            benf_activated_time,
            created_at,
            updated_at,
		benf_account
        FROM beneficiaries
        WHERE user_id = $1
          AND LOWER(benf_name) = LOWER($2)
          AND benf_ifsc = $3
          AND LOWER(benf_nickname) = LOWER($4)
          AND benf_mobile_number = $5 AND benf_account = $6`,
		userId, name, ifscCode, nickName, mobileNumber, maskedAccount,
	)

	err := row.Scan(
		&ben.Id,
		&ben.BenfId,
		&ben.UserId,
		&ben.BenfName,
		&ben.BenfNickName,
		&ben.BenfMobNo,
		&ben.BenfIfsc,
		&ben.BenfAcctType,
		&ben.PaymentMode,
		&ben.BenfActivateTime,
		&ben.CreatedAt,
		&ben.UpdatedAt,
		&ben.BenfAccountNo,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrNoDataFound
		}
		return nil, err
	}

	return ben, nil
}

func FindBeneficiaryByNameAndIfscCodeV2(db *sql.DB, userId, name, ifscCode, nickName, mobileNumber string) (*BeneficiaryDTO, error) {

	ben := NewBeneficiaryDTO()

	row := db.QueryRow(
		`SELECT
            id,
            benf_id,
            user_id,
            benf_name,
            benf_nickname,
            benf_mobile_number,
            benf_ifsc,
            benf_acct_type,
            payment_mode,
            benf_activated_time,
            created_at,
            updated_at
        FROM beneficiaries
        WHERE user_id = $1
          AND LOWER(benf_name) = LOWER($2)
          AND benf_ifsc = $3
          AND LOWER(benf_nickname) = LOWER($4)
          AND benf_mobile_number = $5`,
		userId, name, ifscCode, nickName, mobileNumber,
	)

	err := row.Scan(
		&ben.Id,
		&ben.BenfId,
		&ben.UserId,
		&ben.BenfName,
		&ben.BenfNickName,
		&ben.BenfMobNo,
		&ben.BenfIfsc,
		&ben.BenfAcctType,
		&ben.PaymentMode,
		&ben.BenfActivateTime,
		&ben.CreatedAt,
		&ben.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrNoDataFound
		}
		return nil, err
	}

	return ben, nil
}

func GetBeneficiaryByID(db *sql.DB, beneficiaryID, benfAcctNo, benfMobileNo, benfIFSCCode string) (*BeneficiaryDTO, error) {
	ben := NewBeneficiaryDTO()

	maskedAccount := strings.Repeat("X", len(benfAcctNo)-4) + benfAcctNo[len(benfAcctNo)-4:]

	query := `SELECT id, benf_id, user_id, benf_name, benf_nickname, benf_mobile_number, benf_ifsc, benf_acct_type, payment_mode,
		benf_activated_time, created_at, updated_at, is_active
		FROM beneficiaries WHERE benf_id = $1 AND benf_account = $2 AND benf_mobile_number = $3 AND benf_ifsc = $4 AND is_active = true`

	err := db.QueryRow(query, beneficiaryID, maskedAccount, benfMobileNo, benfIFSCCode).Scan(
		&ben.Id, &ben.BenfId, &ben.UserId, &ben.BenfName, &ben.BenfNickName,
		&ben.BenfMobNo, &ben.BenfIfsc, &ben.BenfAcctType,
		&ben.PaymentMode, &ben.BenfActivateTime,
		&ben.CreatedAt, &ben.UpdatedAt, &ben.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrNoDataFound
		}
		return nil, err
	}

	return ben, nil
}

// Scenario 2 - If Beneficiary data sent to KVB and success res received and left the half way without completion, retried with new request with change in value
// to check the requests body and DB avail data is same or not we are retrving data
func GetBeneficiaryByTxnId(db *sql.DB, txnIdentifier, userId string) (*BeneficiaryDTO, error) {
	ben := NewBeneficiaryDTO()

	query := `SELECT id, benf_id, user_id, benf_name, benf_nickname, benf_mobile_number, benf_ifsc, benf_acct_type, payment_mode,
		benf_activated_time, created_at, updated_at, is_active,txn_identifier,benf_account
		FROM beneficiaries WHERE txn_identifier = $1 AND user_id = $2`

	err := db.QueryRow(query, txnIdentifier, userId).Scan(
		&ben.Id, &ben.BenfId, &ben.UserId, &ben.BenfName, &ben.BenfNickName,
		&ben.BenfMobNo, &ben.BenfIfsc, &ben.BenfAcctType,
		&ben.PaymentMode, &ben.BenfActivateTime,
		&ben.CreatedAt, &ben.UpdatedAt, &ben.IsActive, &ben.TxnIdentifier, &ben.BenfAccountNo,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrNoDataFound
		}
		return nil, err
	}

	return ben, nil
}

func (dto *BeneficiaryDTO) BankBindData(
	beneficiaryId,
	userId,
	beneficiaryName,
	beneficiaryNickName,
	beneficiaryMobileNumber,
	beneficiaryAccount,
	beneficiaryIFSC,
	beneficiaryAccountType,
	paymentMode,
	activatedDtTime string,
	isActive bool) error {

	dto.BenfId = beneficiaryId
	dto.UserId = userId
	dto.BenfName = beneficiaryName
	dto.BenfNickName.String = beneficiaryNickName
	dto.BenfMobNo = beneficiaryMobileNumber
	dto.BenfAccountNo = beneficiaryAccount
	dto.BenfIfsc = beneficiaryIFSC
	dto.BenfAcctType = beneficiaryAccountType
	dto.PaymentMode = paymentMode
	dto.ActivatedDtTime = activatedDtTime
	dto.IsActive = isActive
	return nil
}
