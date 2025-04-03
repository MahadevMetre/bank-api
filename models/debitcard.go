package models

import (
	"bankapi/constants"
	"bankapi/responses"
	"database/sql"
	"errors"
	"fmt"

	"bitbucket.org/paydoh/paydoh-commons/types"
)

type DebitCardData struct {
	TxnIdentifyer           types.NullableString `json:"txn_identifier"`
	UserID                  string               `json:"user_id"`
	CID                     types.NullableString `json:"cid"`
	Proxy_Number            types.NullableString `json:"proxy_number"`
	Enrollment_id           types.NullableString `json:"enrollment_id"`
	DeliveryStatus          types.NullableString `json:"delivery_status"`
	PublicKey               types.NullableString `json:"public_key"`
	IsPermanentlyBlocked    bool                 `json:"is_permanently_blocked"`
	PhysicalDebitCardTxnId  types.NullableString `json:"physical_debitcard_txnid"`
	IsVirtualCardGenerated  bool                 `json:"is_virtual_generated"`
	IsPhysicalCardGenerated bool                 `json:"is_physical_generated"`
}

func NewDebitCardData() *DebitCardData {
	return &DebitCardData{}
}

func (data *DebitCardData) Bind(accountRes *UserPersonalInformationAndAccountData, debitRes *responses.GenerateDebitcardResponse) {
	data.CID.String = accountRes.CustomerId
	data.UserID = accountRes.UserId
	data.Proxy_Number.String = debitRes.ProxyNumber
	data.TxnIdentifyer.String = debitRes.TxnIdentifier
}

func InsertDebitCardData(db *sql.DB, userId, txnIdetifier string) error {
	_, err := db.Exec("INSERT INTO debit_card_data (txn_identifier,user_id) VALUES  ($1,$2)", txnIdetifier, userId)
	if err != nil {
		return err
	}
	return nil
}

func GetDebitCardData(db *sql.DB, user_id string) (*DebitCardData, error) {
	data := NewDebitCardData()
	row := db.QueryRow("SELECT txn_identifier,user_id,cid,proxy_number,enrollment_id,public_key,is_permanently_blocked,physical_debitcard_txnid,is_virtual_generated,is_physical_generated FROM debit_card_data WHERE user_id=$1", user_id)

	if err := row.Scan(
		&data.TxnIdentifyer,
		&data.UserID,
		&data.CID,
		&data.Proxy_Number,
		&data.Enrollment_id,
		&data.PublicKey,
		&data.IsPermanentlyBlocked,
		&data.PhysicalDebitCardTxnId,
		&data.IsVirtualCardGenerated,
		&data.IsPhysicalCardGenerated,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrNoDataFound
		}

		return nil, err
	}
	return data, nil
}

func UpdateEnrollmentID(db *sql.DB, userId string, enrollmentID string) error {
	_, err := db.Exec("UPDATE debit_card_data SET enrollment_id = $1 WHERE user_id = $2", enrollmentID, userId)
	if err != nil {
		return err
	}
	return nil
}

func UpdateDebitCardData(db *sql.DB, data *DebitCardData) error {
	if data == nil {
		return errors.New("debit card data is nil")
	}

	query := "UPDATE debit_card_data SET"
	params := []interface{}{}
	paramCount := 1

	if data.TxnIdentifyer.Valid && data.TxnIdentifyer.String != "" {
		query += fmt.Sprintf(" txn_identifier=$%d,", paramCount)
		params = append(params, data.TxnIdentifyer.String)
		paramCount++
	}

	if data.CID.Valid && data.CID.String != "" {
		query += fmt.Sprintf(" cid=$%d,", paramCount)
		params = append(params, data.CID)
		paramCount++
	}

	if data.Proxy_Number.Valid && data.Proxy_Number.String != "" {
		query += fmt.Sprintf(" proxy_number=$%d,", paramCount)
		params = append(params, data.Proxy_Number)
		paramCount++
	}

	if data.Enrollment_id.Valid && data.Enrollment_id.String != "" {
		query += fmt.Sprintf(" enrollment_id=$%d,", paramCount)
		params = append(params, data.Enrollment_id.String)
		paramCount++
	}

	if data.DeliveryStatus.Valid && data.DeliveryStatus.String != "" {
		query += fmt.Sprintf(" delivery_status=$%d,", paramCount)
		params = append(params, data.DeliveryStatus.String)
		paramCount++
	}

	if data.PublicKey.Valid && data.PublicKey.String != "" {
		query += fmt.Sprintf(" public_key=$%d,", paramCount)
		params = append(params, data.PublicKey.String)
		paramCount++
	}

	if data.IsPermanentlyBlocked {
		query += fmt.Sprintf(" is_permanently_blocked=$%d,", paramCount)
		params = append(params, data.IsPermanentlyBlocked)
		paramCount++
	}

	if data.PhysicalDebitCardTxnId.Valid && data.PhysicalDebitCardTxnId.String != "" {
		query += fmt.Sprintf(" physical_debitcard_txnid=$%d,", paramCount)
		params = append(params, data.PhysicalDebitCardTxnId.String)
		paramCount++
	}

	if data.IsVirtualCardGenerated {
		query += fmt.Sprintf(" is_virtual_generated=$%d,", paramCount)
		params = append(params, data.IsVirtualCardGenerated)
		paramCount++
	}
	if data.IsPhysicalCardGenerated {
		query += fmt.Sprintf(" is_physical_generated=$%d,", paramCount)
		params = append(params, data.IsPhysicalCardGenerated)
		paramCount++
	}

	if len(params) == 0 {
		return errors.New("no fields to update")
	}

	query = query[:len(query)-1] + fmt.Sprintf(" ,updated_at = now() WHERE user_id=$%d", paramCount)
	params = append(params, data.UserID)

	_, err := db.Exec(query, params...)
	return err
}
