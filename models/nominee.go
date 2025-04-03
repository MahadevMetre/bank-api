package models

import (
	"bankapi/config"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"bitbucket.org/paydoh/paydoh-commons/types"
	"github.com/google/uuid"
)

type Nominee struct {
	ID                  types.NullableString `json:"id"`
	AccountDataID       uuid.UUID            `json:"account_data_id"`
	UserId              string               `json:"user_id"`
	NomApplicantID      types.NullableString `json:"nom_applicant_id"`
	NomReqType          types.NullableString `json:"nom_req_type"`
	NomCBSStatus        types.NullableString `json:"nom_cbs_status"`
	NomUpdatedTime      types.NullableString `json:"nom_updated_time"`
	TxnIdentifier       types.NullableString `json:"txn_identifier"`
	NomName             types.NullableString `json:"nom_name"`
	IsVerified          bool                 `json:"is_verified"`
	IsActive            bool                 `json:"is_active"`
	CreatedAt           time.Time            `json:"created_at"`
	UpdatedAt           time.Time            `json:"updated_at"`
	DateOfBirth         types.NullableString `json:"date_of_birth"`
	Relation            types.NullableString `json:"relation"`
	Address1            types.NullableString `json:"address_1"`
	Address2            types.NullableString `json:"address_2"`
	Address3            types.NullableString `json:"address_3"`
	City                types.NullableString `json:"city"`
	IsOtpSent           bool                 `json:"is_otp_sent"`
	Pincode             types.NullableString `json:"pincode"`
	NomineeMobileNumber types.NullableString `json:"nominee_mobile_number"`
}

func InsertNominee(nominee Nominee) (string, error) {
	var id uuid.UUID
	query := `INSERT INTO nominees
			  (account_data_id, nom_name, nom_applicant_id, nom_req_type, nom_cbs_status, nom_updated_time, txn_identifier, is_verified, is_active, user_id, date_of_birth, relation, address_1, address_2, address_3, city, is_otp_sent, pincode,nominee_mobile_number)
			  VALUES
			  ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
			  RETURNING id`
	err := config.GetDB().QueryRow(query,
		nominee.AccountDataID,
		nominee.NomName,
		nominee.NomApplicantID,
		nominee.NomReqType,
		nominee.NomCBSStatus,
		nominee.NomUpdatedTime,
		nominee.TxnIdentifier,
		nominee.IsVerified,
		nominee.IsActive,
		nominee.UserId,
		nominee.DateOfBirth,
		nominee.Relation,
		nominee.Address1,
		nominee.Address2,
		nominee.Address3,
		nominee.City,
		nominee.IsOtpSent,
		nominee.Pincode, nominee.NomineeMobileNumber).Scan(&id)
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

func UpdateNomineeByUserId(db *sql.DB, nominee *Nominee) error {
	updates := make([]string, 0)
	values := make([]interface{}, 0)
	valueCount := 1

	if nominee.ID.Valid {
		updates = append(updates, fmt.Sprintf("id = $%d", valueCount))
		values = append(values, nominee.ID) // Convert to UUID safely
		valueCount++
	}

	if nominee.AccountDataID != uuid.Nil {
		updates = append(updates, fmt.Sprintf("account_data_id = $%d", valueCount))
		values = append(values, nominee.AccountDataID) // Convert to UUID safely
		valueCount++
	}

	if nominee.NomineeMobileNumber.Valid {
		updates = append(updates, fmt.Sprintf("nominee_mobile_number = $%d", valueCount))
		values = append(values, nominee.NomineeMobileNumber.String)
		valueCount++
	}

	if nominee.NomName.Valid {
		updates = append(updates, fmt.Sprintf("nom_name = $%d", valueCount))
		values = append(values, nominee.NomName.String)
		valueCount++
	}

	if nominee.NomApplicantID.Valid {
		updates = append(updates, fmt.Sprintf("nom_applicant_id = $%d", valueCount))
		values = append(values, nominee.NomApplicantID.String)
		valueCount++
	}

	if nominee.NomReqType.Valid {
		updates = append(updates, fmt.Sprintf("nom_req_type = $%d", valueCount))
		values = append(values, nominee.NomReqType.String)
		valueCount++
	}

	if nominee.NomCBSStatus.Valid {
		updates = append(updates, fmt.Sprintf("nom_cbs_status = $%d", valueCount))
		values = append(values, nominee.NomCBSStatus.String)
		valueCount++
	}

	if nominee.NomUpdatedTime.Valid {
		updates = append(updates, fmt.Sprintf("nom_updated_time = $%d", valueCount))
		values = append(values, nominee.NomUpdatedTime.String)
		valueCount++
	}

	if nominee.TxnIdentifier.Valid {
		updates = append(updates, fmt.Sprintf("txn_identifier = $%d", valueCount))
		values = append(values, nominee.TxnIdentifier.String)
		valueCount++
	}

	if nominee.DateOfBirth.Valid {
		updates = append(updates, fmt.Sprintf("date_of_birth = $%d", valueCount))
		values = append(values, nominee.DateOfBirth.String)
		valueCount++
	}

	if nominee.Relation.Valid {
		updates = append(updates, fmt.Sprintf("relation = $%d", valueCount))
		values = append(values, nominee.Relation.String)
		valueCount++
	}

	if nominee.Address1.Valid {
		updates = append(updates, fmt.Sprintf("address_1 = $%d", valueCount))
		values = append(values, nominee.Address1.String)
		valueCount++
	}

	if nominee.Address2.Valid {
		updates = append(updates, fmt.Sprintf("address_2 = $%d", valueCount))
		values = append(values, nominee.Address2.String)
		valueCount++
	}

	if nominee.Address3.Valid {
		updates = append(updates, fmt.Sprintf("address_3 = $%d", valueCount))
		values = append(values, nominee.Address3.String)
		valueCount++
	}

	if nominee.City.Valid {
		updates = append(updates, fmt.Sprintf("city = $%d", valueCount))
		values = append(values, nominee.City.String)
		valueCount++
	}

	if nominee.Pincode.Valid {
		updates = append(updates, fmt.Sprintf("pincode = $%d", valueCount))
		values = append(values, nominee.Pincode.String)
		valueCount++
	}

	if nominee.IsVerified || !nominee.IsVerified {
		updates = append(updates, fmt.Sprintf("is_verified = $%d", valueCount))
		values = append(values, nominee.IsVerified)
		valueCount++
	}

	if nominee.IsActive || !nominee.IsActive {
		updates = append(updates, fmt.Sprintf("is_active = $%d", valueCount))
		values = append(values, nominee.IsActive)
		valueCount++
	}

	if nominee.IsOtpSent || !nominee.IsOtpSent {
		updates = append(updates, fmt.Sprintf("is_otp_sent = $%d", valueCount))
		values = append(values, nominee.IsOtpSent)
		valueCount++
	}

	if len(updates) == 0 {
		return nil
	}

	values = append(values, nominee.UserId)

	query := fmt.Sprintf(`
		UPDATE nominees
		SET %s, updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $%d`,
		strings.Join(updates, ", "),
		valueCount)

	result, err := db.Exec(query, values...)
	if err != nil {
		return fmt.Errorf("failed to update nominee: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("no nominee found with user_id %s", nominee.UserId)
	}

	return nil
}

func FindOneNomineeByUserID(db *sql.DB, userId string) (*Nominee, error) {
	nominee := &Nominee{}

	row := db.QueryRow(
		`SELECT id, account_data_id, nom_name, nom_applicant_id, nom_req_type, nom_cbs_status, nom_updated_time,
       		txn_identifier, is_verified, is_active, created_at, updated_at, date_of_birth, relation,
       		address_1, address_2, address_3, city, is_otp_sent, pincode, nominee_mobile_number
		 	FROM nominees WHERE user_id = $1;`,
		userId,
	)

	err := row.Scan(
		&nominee.ID,
		&nominee.AccountDataID,
		&nominee.NomName,
		&nominee.NomApplicantID,
		&nominee.NomReqType,
		&nominee.NomCBSStatus,
		&nominee.NomUpdatedTime,
		&nominee.TxnIdentifier,
		&nominee.IsVerified,
		&nominee.IsActive,
		&nominee.CreatedAt,
		&nominee.UpdatedAt,
		&nominee.DateOfBirth,
		&nominee.Relation,
		&nominee.Address1,
		&nominee.Address2,
		&nominee.Address3,
		&nominee.City,
		&nominee.IsOtpSent,
		&nominee.Pincode,
		&nominee.NomineeMobileNumber,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return nominee, nil
}

func FindOneVerfiedNomineeByUserID(db *sql.DB, userId string) (*Nominee, error) {
	nominee := &Nominee{}

	row := db.QueryRow(
		`SELECT id, account_data_id, nom_name, nom_applicant_id, nom_req_type, nom_cbs_status, nom_updated_time,
       		txn_identifier, is_verified, is_active, created_at, updated_at, date_of_birth, relation,
       		address_1, address_2, address_3, city, is_otp_sent, pincode, nominee_mobile_number
		 	FROM nominees WHERE user_id = $1 AND is_verified = true;`,
		userId,
	)

	err := row.Scan(
		&nominee.ID,
		&nominee.AccountDataID,
		&nominee.NomName,
		&nominee.NomApplicantID,
		&nominee.NomReqType,
		&nominee.NomCBSStatus,
		&nominee.NomUpdatedTime,
		&nominee.TxnIdentifier,
		&nominee.IsVerified,
		&nominee.IsActive,
		&nominee.CreatedAt,
		&nominee.UpdatedAt,
		&nominee.DateOfBirth,
		&nominee.Relation,
		&nominee.Address1,
		&nominee.Address2,
		&nominee.Address3,
		&nominee.City,
		&nominee.IsOtpSent,
		&nominee.Pincode,
		&nominee.NomineeMobileNumber,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return nominee, nil
}

func UpdateNomineeActiveStatus(userId string, status bool) error {
	query := `UPDATE nominees SET is_active = $1 WHERE user_id = $2`
	_, err := config.GetDB().Exec(query, status, userId)
	if err != nil {
		return err
	}
	return nil
}

func GetNomineeByAccountDataID(accountDataID uuid.UUID) (Nominee, error) {
	var nominee Nominee
	query := `SELECT id, account_data_id, nom_name, nom_applicant_id, nom_req_type, nom_cbs_status,txn_identifier, nom_updated_time, created_at, updated_at
              FROM nominees WHERE account_data_id = $1 AND is_active = true`
	err := config.GetDB().QueryRow(query, accountDataID).Scan(&nominee.ID, &nominee.AccountDataID, &nominee.NomName, &nominee.NomApplicantID, &nominee.NomReqType, &nominee.NomCBSStatus, &nominee.TxnIdentifier, &nominee.NomUpdatedTime, &nominee.CreatedAt, &nominee.UpdatedAt)
	if err != nil {
		return nominee, err
	}
	return nominee, nil
}

func UpdateNominee(nominee Nominee) error {
	query := `UPDATE nominees SET nom_name = $1, nom_applicant_id = $2, nom_req_type = $3, nom_cbs_status = $4, nom_updated_time = $5, updated_at = now()
              WHERE id = $6`
	_, err := config.GetDB().Exec(query, nominee.NomName, nominee.NomApplicantID, nominee.NomReqType, nominee.NomCBSStatus, nominee.NomUpdatedTime, nominee.ID)
	return err
}

func DeleteNominee(userId, nomineeInsertId string) error {
	query := `DELETE FROM nominees WHERE user_id = $1 AND id != $2`
	_, err := config.GetDB().Exec(query, userId, nomineeInsertId)
	return err
}
