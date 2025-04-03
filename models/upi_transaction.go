package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type UPITransaction struct {
	ID            uuid.UUID
	UserID        string
	TransactionID string
	CredType      string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// CreateUPITransaction inserts a new UPI transaction into the database
func CreateUPITransaction(db *sql.DB, ut *UPITransaction) error {
	query := `
		INSERT INTO upi_transactions (id, user_id, transaction_id, cred_type)
		VALUES ($1, $2, $3, $4)
		RETURNING created_at, updated_at
	`
	ut.ID = uuid.New()
	err := db.QueryRow(query, ut.ID, ut.UserID, ut.TransactionID, ut.CredType).
		Scan(&ut.CreatedAt, &ut.UpdatedAt)
	return err
}

// GetUPITransaction retrieves a UPI transaction by ID
func GetUPITransaction(db *sql.DB, id uuid.UUID) (*UPITransaction, error) {
	query := `
		SELECT id, user_id, transaction_id, cred_type, created_at, updated_at
		FROM upi_transactions
		WHERE id = $1
	`
	ut := &UPITransaction{}
	err := db.QueryRow(query, id).
		Scan(&ut.ID, &ut.UserID, &ut.TransactionID, &ut.CredType, &ut.CreatedAt, &ut.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return ut, nil
}

// UpdateUPITransaction updates an existing UPI transaction
func UpdateUPITransaction(db *sql.DB, ut *UPITransaction) error {
	query := `
		UPDATE upi_transactions
		SET user_id = $2, transaction_id = $3, cred_type = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at
	`
	err := db.QueryRow(query, ut.ID, ut.UserID, ut.TransactionID, ut.CredType).
		Scan(&ut.UpdatedAt)
	return err
}

// DeleteUPITransaction deletes a UPI transaction by ID
func DeleteUPITransaction(db *sql.DB, id uuid.UUID) error {
	query := "DELETE FROM upi_transactions WHERE id = $1"
	_, err := db.Exec(query, id)
	return err
}

// ListUPITransactions retrieves all UPI transactions
func ListUPITransactions(db *sql.DB) ([]*UPITransaction, error) {
	query := `
		SELECT id, user_id, transaction_id, cred_type, created_at, updated_at
		FROM upi_transactions
		ORDER BY created_at DESC
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*UPITransaction
	for rows.Next() {
		ut := &UPITransaction{}
		err := rows.Scan(&ut.ID, &ut.UserID, &ut.TransactionID, &ut.CredType, &ut.CreatedAt, &ut.UpdatedAt)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, ut)
	}
	return transactions, rows.Err()
}
