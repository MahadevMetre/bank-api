package models

import (
	"bankapi/constants"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/paydoh/paydoh-commons/types"
	"github.com/google/uuid"
)

type PaymentMode string

const (
	PaymentModeIFT  PaymentMode = "IFT"
	PaymentModeIMPS PaymentMode = "IMPS"
	PaymentModeNEFT PaymentMode = "NEFT"
	PaymentModeUPI  PaymentMode = "UPI"
)

type Transaction struct {
	ID              uuid.UUID            `db:"id" json:"id"`
	UserID          string               `db:"user_id" json:"user_id"`
	TransactionID   string               `db:"transaction_id" json:"transaction_id"`
	TransactionDesc types.NullableString `db:"transaction_desc" json:"transaction_desc"`
	BeneficiaryID   types.NullableString `db:"beneficiary_id" json:"beneficiary_id"`
	PaymentMode     PaymentMode          `db:"payment_mode" json:"payment_mode"`
	Amount          types.NullableString `db:"amount" json:"amount"`
	UTRRefNumber    types.NullableString `db:"utr_ref_number" json:"utr_ref_number"`
	OTPStatus       types.NullableString `db:"otp_status" json:"otp_status"`
	CBSStatus       types.NullableString `db:"cbs_status" json:"cbs_status"`
	UPIPayeeAddr    types.NullableString `db:"upi_payee_addr" json:"upi_payee_addr"`
	UPIPayeeName    types.NullableString `db:"upi_payee_addr" json:"upi_payee_name"`
	CreatedAt       time.Time            `db:"created_at" json:"-"`
	UpdatedAt       time.Time            `db:"updated_at" json:"-"`
	TransactionType string               `db:"transaction_type" json:"transaction_type"`
}

type TransactionDTO struct {
	ID                 uuid.UUID            `json:"id"`
	UserDetails        UserDetails          `json:"user"`
	TransactionID      string               `json:"transaction_id"`
	TransactionDesc    types.NullableString `json:"transaction_desc"`
	BeneficiaryID      types.NullableString `json:"beneficiary_id"`
	PaymentMode        string               `json:"payment_mode"`
	Amount             string               `json:"amount"`
	UTRRefNumber       types.NullableString `json:"utr_ref_number"`
	UPIPayeeAddr       types.NullableString `json:"upi_payee_addr"`
	UPIPayeeName       types.NullableString `json:"upi_payee_name"`
	BeneficiaryName    types.NullableString `json:"beneficiary_name"`
	BeneficiaryIFSC    types.NullableString `json:"beneficiary_ifsc"`
	BeneficiaryAccount types.NullableString `json:"beneficiary_account"`
}

type UserDetails struct {
	UserID        string               `json:"user_id"`
	FirstName     types.NullableString `json:"first_name"`
	MiddleName    types.NullableString `json:"middle_name"`
	LastName      types.NullableString `json:"last_name"`
	UpiID         types.NullableString `json:"upi_id"`
	AccountNumber types.NullableString `json:"account_number"`
}

type TransactionFilter struct {
	StartDate       *time.Time
	EndDate         *time.Time
	PaymentMode     string
	CBSStatus       string
	Page            int
	ItemsPerPage    int
	TransactionType string // "SENT", "RECEIVED", or "ALL"
}

type FilterParams struct {
	UserID       string
	UtrRefNumber string
	CodeDRCR     string
	Amount       string
	TxnDate      string
}

func InsertTransaction(db *sql.DB, tx *Transaction) error {
	columns := make([]string, 0)
	values := make([]interface{}, 0)
	placeholders := make([]string, 0)
	valueCount := 1

	columns = append(columns, "user_id", "transaction_id")
	values = append(values, tx.UserID, tx.TransactionID)
	placeholders = append(placeholders, fmt.Sprintf("$%d", valueCount), fmt.Sprintf("$%d", valueCount+1))
	valueCount += 2

	if tx.TransactionDesc.Valid {
		columns = append(columns, "transaction_desc")
		values = append(values, tx.TransactionDesc.String)
		placeholders = append(placeholders, fmt.Sprintf("$%d", valueCount))
		valueCount++
	}

	if tx.BeneficiaryID.Valid {
		columns = append(columns, "beneficiary_id")
		values = append(values, tx.BeneficiaryID.String)
		placeholders = append(placeholders, fmt.Sprintf("$%d", valueCount))
		valueCount++
	}

	if tx.PaymentMode != "" {
		columns = append(columns, "payment_mode")
		values = append(values, tx.PaymentMode)
		placeholders = append(placeholders, fmt.Sprintf("$%d", valueCount))
		valueCount++
	}

	if tx.Amount.Valid {
		columns = append(columns, "amount")
		values = append(values, tx.Amount)
		placeholders = append(placeholders, fmt.Sprintf("$%d", valueCount))
		valueCount++
	}

	if tx.UPIPayeeAddr.Valid {
		columns = append(columns, "upi_payee_addr")
		values = append(values, tx.UPIPayeeAddr.String)
		placeholders = append(placeholders, fmt.Sprintf("$%d", valueCount))
		valueCount++
	}

	if tx.OTPStatus.Valid {
		columns = append(columns, "otp_status")
		values = append(values, tx.OTPStatus.String)
		placeholders = append(placeholders, fmt.Sprintf("$%d", valueCount))
		valueCount++
	}

	if tx.CBSStatus.Valid {
		columns = append(columns, "cbs_status")
		values = append(values, tx.CBSStatus.String)
		placeholders = append(placeholders, fmt.Sprintf("$%d", valueCount))
		valueCount++
	}

	if tx.UTRRefNumber.Valid {
		columns = append(columns, "utr_ref_number")
		values = append(values, tx.UTRRefNumber.String)
		placeholders = append(placeholders, fmt.Sprintf("$%d", valueCount))
		valueCount++
	}

	query := fmt.Sprintf(`
		INSERT INTO transactions (%s)
		VALUES (%s)
		RETURNING id, created_at, updated_at`,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))

	err := db.QueryRow(query, values...).Scan(&tx.ID, &tx.CreatedAt, &tx.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert transaction: %w", err)
	}

	return nil
}

func UpdateTransactionByTransID(db *sql.DB, tx *Transaction) error {
	updates := make([]string, 0)
	values := make([]interface{}, 0)
	valueCount := 1

	if tx.TransactionDesc.Valid {
		updates = append(updates, fmt.Sprintf("transaction_desc = $%d", valueCount))
		values = append(values, tx.TransactionDesc.String)
		valueCount++
	}

	if tx.BeneficiaryID.Valid {
		updates = append(updates, fmt.Sprintf("beneficiary_id = $%d", valueCount))
		values = append(values, tx.BeneficiaryID.String)
		valueCount++
	}

	if tx.PaymentMode != "" {
		updates = append(updates, fmt.Sprintf("payment_mode = $%d", valueCount))
		values = append(values, tx.PaymentMode)
		valueCount++
	}

	if tx.Amount.Valid {
		updates = append(updates, fmt.Sprintf("amount = $%d", valueCount))
		values = append(values, tx.Amount)
		valueCount++
	}

	if tx.UPIPayeeAddr.Valid {
		updates = append(updates, fmt.Sprintf("upi_payee_addr = $%d", valueCount))
		values = append(values, tx.UPIPayeeAddr.String)
		valueCount++
	}

	if tx.UPIPayeeName.Valid {
		updates = append(updates, fmt.Sprintf("upi_payee_name = $%d", valueCount))
		values = append(values, tx.UPIPayeeName.String)
		valueCount++
	}

	if tx.OTPStatus.Valid {
		updates = append(updates, fmt.Sprintf("otp_status = $%d", valueCount))
		values = append(values, tx.OTPStatus.String)
		valueCount++
	}

	if tx.CBSStatus.Valid {
		updates = append(updates, fmt.Sprintf("cbs_status = $%d", valueCount))
		values = append(values, tx.CBSStatus.String)
		valueCount++
	}

	if tx.UTRRefNumber.Valid {
		updates = append(updates, fmt.Sprintf("utr_ref_number = $%d", valueCount))
		values = append(values, tx.UTRRefNumber.String)
		valueCount++
	}

	if len(updates) == 0 {
		return nil
	}

	values = append(values, tx.TransactionID)

	query := fmt.Sprintf(`
		UPDATE transactions
		SET %s, updated_at = CURRENT_TIMESTAMP
		WHERE transaction_id = $%d`,
		strings.Join(updates, ", "),
		valueCount)

	result, err := db.Exec(query, values...)
	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("no transaction found with ID %s", tx.ID)
	}

	return nil
}

func FetchUserTransactions(db *sql.DB, userID string, filter *TransactionFilter) ([]Transaction, error) {
	if userID == "" {
		return nil, errors.New("user ID is required")
	}

	query := `
		WITH user_transactions AS (
			-- Transactions where user is the sender
			SELECT
				t.id, t.user_id, t.transaction_id, t.transaction_desc,
				t.beneficiary_id, t.payment_mode, t.amount, t.utr_ref_number,
				t.otp_status, t.cbs_status, t.upi_payee_addr, t.created_at,
				t.updated_at, 'SENT' as transaction_type
			FROM transactions t
			WHERE t.user_id = $1

			UNION ALL

			-- Transactions where user is the beneficiary
			SELECT
				t.id, t.user_id, t.transaction_id, t.transaction_desc,
				t.beneficiary_id, t.payment_mode, t.amount, t.utr_ref_number,
				t.otp_status, t.cbs_status, t.upi_payee_addr, t.created_at,
				t.updated_at, 'RECEIVED' as transaction_type
			FROM transactions t
			WHERE (t.beneficiary_id = $1 OR t.upi_payee_addr = $1)
		)
		SELECT * FROM user_transactions
	`

	args := []interface{}{userID}
	paramCount := 1
	whereClause := ""

	if filter != nil {
		if filter.TransactionType != "" {
			switch filter.TransactionType {
			case "SENT":
				whereClause += " WHERE transaction_type = 'SENT'"
			case "RECEIVED":
				whereClause += " WHERE transaction_type = 'RECEIVED'"
			case "ALL":
			default:
				return nil, errors.New("invalid transaction type")
			}
		}

		if filter.StartDate != nil {
			paramCount++
			if whereClause == "" {
				whereClause = " WHERE"
			} else {
				whereClause += " AND"
			}
			whereClause += fmt.Sprintf(" created_at >= $%d", paramCount)
			args = append(args, filter.StartDate)
		}

		if filter.EndDate != nil {
			paramCount++
			if whereClause == "" {
				whereClause = " WHERE"
			} else {
				whereClause += " AND"
			}
			whereClause += ` created_at <= $` + strconv.Itoa(paramCount)
			args = append(args, filter.EndDate)
		}

		if filter.PaymentMode != "" {
			paramCount++
			if whereClause == "" {
				whereClause = " WHERE"
			} else {
				whereClause += " AND"
			}
			whereClause += ` payment_mode = $` + strconv.Itoa(paramCount)
			args = append(args, filter.PaymentMode)
		}

		if filter.CBSStatus != "" {
			paramCount++
			if whereClause == "" {
				whereClause = " WHERE"
			} else {
				whereClause += " AND"
			}
			whereClause += ` cbs_status = $` + strconv.Itoa(paramCount)
			args = append(args, filter.CBSStatus)
		}
	}

	query += whereClause

	// Add sorting and pagination
	query += ` ORDER BY created_at DESC`
	if filter != nil && filter.ItemsPerPage > 0 {
		offset := (filter.Page - 1) * filter.ItemsPerPage
		query += ` LIMIT $` + strconv.Itoa(paramCount+1) + ` OFFSET $` + strconv.Itoa(paramCount+2)
		args = append(args, filter.ItemsPerPage, offset)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		err := rows.Scan(
			&t.ID,
			&t.UserID,
			&t.TransactionID,
			&t.TransactionDesc,
			&t.BeneficiaryID,
			&t.PaymentMode,
			&t.Amount,
			&t.UTRRefNumber,
			&t.OTPStatus,
			&t.CBSStatus,
			&t.UPIPayeeAddr,
			&t.CreatedAt,
			&t.UpdatedAt,
			&t.TransactionType,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func FetchOneTransaction(db *sql.DB, filter *FilterParams) (*TransactionDTO, error) {
	parsedTime, err := time.Parse("02-01-2006 15:04:05", filter.TxnDate)
	if err != nil {
		return nil, fmt.Errorf("error parsing date: %v", err)
	}

	formattedDate := parsedTime.Format("2006-01-02 15:04:05")

	timestamp, err := time.Parse("2006-01-02 15:04:05", formattedDate)
	if err != nil {
		return nil, err
	}

	txnTime := timestamp.Add(-5*time.Hour - 30*time.Minute)

	query := `
    SELECT
        t.id,
        t.user_id,
        p.first_name,
        p.middle_name,
        p.last_name,
        a.upi_id,
        a.account_number,
        t.transaction_id,
        t.transaction_desc,
        t.beneficiary_id,
        t.payment_mode,
        t.amount,
        t.utr_ref_number,
        t.upi_payee_addr,
        t.upi_payee_name,
        CASE
            WHEN t.beneficiary_id IS NOT NULL AND t.beneficiary_id != '' THEN COALESCE(b.benf_name, '')
            ELSE ''
        END as beneficiary_name,
        CASE
            WHEN t.beneficiary_id IS NOT NULL AND t.beneficiary_id != '' THEN COALESCE(b.benf_ifsc, '')
            ELSE ''
        END as beneficiary_ifsc,
        CASE
            WHEN t.beneficiary_id IS NOT NULL AND t.beneficiary_id != '' THEN COALESCE(b.benf_account, '')
            ELSE ''
        END as beneficiary_account
    FROM transactions t
    LEFT JOIN account_data a ON t.user_id = a.user_id
    LEFT JOIN personal_information p ON t.user_id = p.user_id
    LEFT JOIN beneficiaries b ON t.beneficiary_id IS NOT NULL
        AND t.beneficiary_id != ''
        AND (
			t.beneficiary_id = b.id::text
			OR t.beneficiary_id = b.benf_id::text
		)
    WHERE t.amount = $1`

	var args []interface{} = []interface{}{
		filter.Amount,
	}

	if filter.UtrRefNumber != "" {
		query += " AND t.utr_ref_number = $2"
		args = append(args, filter.UtrRefNumber)
	} else if filter.TxnDate != "" {
		query += " AND LOWER(t.cbs_status) = 'success' AND t.created_at BETWEEN ($2::timestamp - interval '1 minute') AND ($2::timestamp + interval '1 minute')"
		args = append(args, txnTime)
	}

	var transaction TransactionDTO
	err = db.QueryRow(query, args...).Scan(
		&transaction.ID,
		&transaction.UserDetails.UserID,
		&transaction.UserDetails.FirstName,
		&transaction.UserDetails.MiddleName,
		&transaction.UserDetails.LastName,
		&transaction.UserDetails.UpiID,
		&transaction.UserDetails.AccountNumber,
		&transaction.TransactionID,
		&transaction.TransactionDesc,
		&transaction.BeneficiaryID,
		&transaction.PaymentMode,
		&transaction.Amount,
		&transaction.UTRRefNumber,
		&transaction.UPIPayeeAddr,
		&transaction.UPIPayeeName,
		&transaction.BeneficiaryName,
		&transaction.BeneficiaryIFSC,
		&transaction.BeneficiaryAccount,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query error: %v", err)
	}

	return &transaction, nil
}

func GetRecentTransactionUsers(db *sql.DB, userId string) ([]map[string]interface{}, error) {
	query := `
		WITH RankedTransactions AS (
			SELECT t.id, t.upi_payee_addr, t.upi_payee_name, t.created_at,
				   ROW_NUMBER() OVER (PARTITION BY t.upi_payee_addr ORDER BY t.created_at DESC) as rn
			FROM transactions t
			WHERE t.payment_mode = 'UPI'
			AND t.user_id = $1
			AND t.upi_payee_addr IS NOT NULL
			AND t.upi_payee_addr != ''
		)
		SELECT id, upi_payee_addr, upi_payee_name, created_at
		FROM RankedTransactions
		WHERE rn = 1
		ORDER BY created_at DESC
		LIMIT 20
	`
	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var txnID, upiID string
		var createdAt string
		var upiPayeeName sql.NullString
		if err := rows.Scan(&txnID, &upiID, &upiPayeeName, &createdAt); err != nil {
			return nil, err
		}
		user := map[string]interface{}{
			"txn_id":         txnID,
			"upi_payee_addr": upiID,
			"upi_payee_name": upiPayeeName.String,
			"created_at":     createdAt,
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func FindOneTransactionByUserAndTransactionId(db *sql.DB, userId, transactionId string) (*Transaction, error) {
	transaction := Transaction{}

	row := db.QueryRow("SELECT id, transaction_id,amount, utr_ref_number, cbs_status FROM transactions WHERE user_id = $1 AND transaction_id = $2", userId, transactionId)

	if err := row.Scan(
		&transaction.ID,
		&transaction.TransactionID,
		&transaction.Amount,
		&transaction.UTRRefNumber,
		&transaction.CBSStatus,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrNoDataFound
		}
		return nil, err
	}

	return &transaction, nil
}
