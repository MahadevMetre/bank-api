package models

import (
	"database/sql"
	"time"

	"bankapi/constants"
	"bankapi/requests"
)

type PaymentStatus struct {
	Id                int64     `json:"id"`
	UserId            string    `json:"user_id"`
	TransactionStatus string    `json:"txn_status"`
	TransactionId     string    `json:"txn_id"`
	ReceiptId         string    `json:"receipt_id"`
	Remarks           string    `json:"payment_remarks"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func NewPaymentStatus() *PaymentStatus {
	return &PaymentStatus{}
}

func (s *PaymentStatus) Bind(request *requests.AddPaymentStatusRequest, userId string, remarks string) {
	s.UserId = userId
	s.TransactionStatus = request.TransactionStatus
	s.TransactionId = request.TransactionId
	s.ReceiptId = request.ReceiptId
	s.Remarks = remarks
}

func InsertPaymentStatus(db *sql.DB, paymentStatus *PaymentStatus) error {
	if _, err := db.Exec(
		`INSERT into payment_status(
			user_id,
			txn_status,
			txn_id,
			receipt_id,
			payment_remarks
		) VALUES ($1, $2, $3, $4, $5)`,
		paymentStatus.UserId,
		paymentStatus.TransactionStatus,
		paymentStatus.TransactionId,
		paymentStatus.ReceiptId,
		paymentStatus.Remarks,
	); err != nil {
		return err
	}
	return nil
}

func GetPaymentStatusByUserId(db *sql.DB, userId string) (*PaymentStatus, error) {
	paymentStatus := NewPaymentStatus()
	row := db.QueryRow(
		`SELECT
			id,
			user_id,
			txn_status,
			txn_id,
			receipt_id,
			payment_remarks,
			created_at,
			updated_at
		FROM payment_status
		WHERE 
			user_id = $1`,
		userId,
	)

	if err := row.Scan(
		&paymentStatus.Id,
		&paymentStatus.UserId,
		&paymentStatus.TransactionStatus,
		&paymentStatus.TransactionId,
		&paymentStatus.ReceiptId,
		&paymentStatus.Remarks,
		&paymentStatus.CreatedAt,
		&paymentStatus.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrNoDataFound
		}

		return nil, err
	}

	return paymentStatus, nil
}
