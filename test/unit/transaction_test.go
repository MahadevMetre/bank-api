package unittest

import (
	"bankapi/models"
	"database/sql"
	"testing"
	"time"

	"bitbucket.org/paydoh/paydoh-commons/types"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsertTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	tests := []struct {
		name    string
		tx      *models.Transaction
		mock    func(sqlmock.Sqlmock, *models.Transaction)
		wantErr bool
	}{
		{
			name: "successful insert",
			tx: &models.Transaction{
				UserID:        "user123",
				TransactionID: "tx123",
				PaymentMode:   "UPI",
				Amount:        types.FromString("100.00"),
			},
			mock: func(mock sqlmock.Sqlmock, tx *models.Transaction) {
				expectedID := uuid.New()
				expectedTime := time.Now()

				mock.ExpectQuery(`INSERT INTO transactions`).
					WithArgs(
						tx.UserID,
						tx.TransactionID,
						tx.PaymentMode,
						tx.Amount,
					).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
						AddRow(expectedID, expectedTime, expectedTime))
			},
			wantErr: false,
		},
		{
			name: "database error",
			tx: &models.Transaction{
				UserID:        "user123",
				TransactionID: "tx123",
			},
			mock: func(mock sqlmock.Sqlmock, tx *models.Transaction) {
				mock.ExpectQuery(`INSERT INTO transactions`).
					WillReturnError(sql.ErrConnDone)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(mock, tt.tx)

			err := models.InsertTransaction(db, tt.tx)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, tt.tx.ID)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUpdateTransaction(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	tests := []struct {
		name    string
		tx      *models.Transaction
		mock    func(sqlmock.Sqlmock, *models.Transaction)
		wantErr bool
	}{
		{
			name: "successful update",
			tx: &models.Transaction{
				ID:     uuid.New(),
				Amount: types.FromString("200.00"),
				// CBSStatus: types.NullableString{
				// 	String: "SUCCESS",
				// 	Valid:  true,
				// },
			},
			mock: func(mock sqlmock.Sqlmock, tx *models.Transaction) {
				mock.ExpectExec(`UPDATE transactions`).
					WithArgs(tx.Amount, tx.CBSStatus.String, tx.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			wantErr: false,
		},
		{
			name: "no rows affected",
			tx: &models.Transaction{
				ID:     uuid.New(),
				Amount: types.FromString("300.00"),
			},
			mock: func(mock sqlmock.Sqlmock, tx *models.Transaction) {
				mock.ExpectExec(`UPDATE transactions`).
					WithArgs(tx.Amount, tx.ID).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(mock, tt.tx)

			err := models.UpdateTransactionByTransID(db, tt.tx)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestFetchUserTransactions(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	userID := "user123"
	mockTime := time.Now()
	mockID := uuid.New()

	tests := []struct {
		name      string
		userID    string
		filter    *models.TransactionFilter
		mock      func(sqlmock.Sqlmock)
		wantCount int
		wantErr   bool
	}{
		{
			name:   "successful fetch",
			userID: userID,
			filter: &models.TransactionFilter{
				TransactionType: "ALL",
				ItemsPerPage:    10,
				Page:            1,
			},
			mock: func(mock sqlmock.Sqlmock) {
				// Mock count query
				mock.ExpectQuery(`SELECT COUNT`).
					WithArgs(userID).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

				// Mock transactions query
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "transaction_id", "transaction_desc",
					"beneficiary_id", "payment_mode", "amount", "utr_ref_number",
					"otp_status", "cbs_status", "upi_payee_addr", "created_at",
					"updated_at", "transaction_type",
				}).AddRow(
					mockID, userID, "tx123", sql.NullString{String: "Test", Valid: true},
					sql.NullString{}, "UPI", "100.00", sql.NullString{},
					sql.NullString{}, sql.NullString{}, sql.NullString{},
					mockTime, mockTime, "SENT",
				)

				mock.ExpectQuery(`SELECT \* FROM user_transactions`).
					WithArgs(userID, 10, 0).
					WillReturnRows(rows)
			},
			wantCount: 1,
			wantErr:   false,
		},
		{
			name:   "empty user ID",
			userID: "",
			filter: nil,
			mock: func(mock sqlmock.Sqlmock) {
				// No expectations needed as it should fail early
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(mock)

			transactions, err := models.FetchUserTransactions(db, tt.userID, tt.filter)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantCount, len(transactions))
			assert.NoError(t, mock.ExpectationsWereMet())

			if len(transactions) > 0 {
				tx := transactions[0]
				assert.Equal(t, mockID, tx.ID)
				assert.Equal(t, userID, tx.UserID)
				assert.NotEmpty(t, tx.TransactionID)
			}
		})
	}
}
