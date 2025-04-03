package models

import (
	"bankapi/constants"
	"bankapi/requests"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// MPIN ...
type MPIN struct {
	ID        int64     `json:"id"`
	UserId    string    `json:"user_id"`
	MPIN      string    `json:"mpin"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

func NewMpin() *MPIN {
	return &MPIN{}
}

func (m *MPIN) Bind(request *requests.MpinRequest, userId string) error {
	m.UserId = userId
	encryptedMpin, err := bcrypt.GenerateFromPassword([]byte(request.Mpin), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	m.MPIN = string(encryptedMpin)
	return nil
}

func InsertMpin(db *sql.DB, request *requests.MpinRequest, userId string) error {
	mpin := NewMpin()
	if err := mpin.Bind(request, userId); err != nil {
		return err
	}

	if _, err := db.Exec("INSERT INTO mpin_data (user_id, mpin) VALUES ($1, $2)",
		mpin.UserId,
		mpin.MPIN,
	); err != nil {
		return err
	}

	return nil
}

func FindOneMpinByUserId(db *sql.DB, userId string) (*MPIN, error) {
	mpin := NewMpin()

	row := db.QueryRow("SELECT id, user_id, mpin, created_at, updated_at FROM mpin_data WHERE user_id = $1", userId)

	if err := row.Scan(
		&mpin.ID,
		&mpin.UserId,
		&mpin.MPIN,
		&mpin.CreatedAt,
		&mpin.UpdatedAt,
	); err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			return nil, constants.ErrNoDataFound
		}
	}

	return mpin, nil
}

type MpinAttemptData struct {
	ID          uuid.UUID    `json:"id"`
	UserId      string       `json:"user_id"`
	Attempts    int          `db:"attempts"`
	LastAttempt sql.NullTime `db:"last_attempt"`
	CreatedAt   time.Time    `db:"created_at"`
	UpdatedAt   time.Time    `db:"updated_at"`
}

func NewMpinAttemptData() *MpinAttemptData {
	return &MpinAttemptData{}
}

func GetMpinAttempts(db *sql.DB, userId string) (*MpinAttemptData, error) {
	attemptData := NewMpinAttemptData()

	query := `
		SELECT 
			id, 
			user_id, 
			attempts, 
			last_attempt, 
			created_at, 
			updated_at 
		FROM 
			mpin_attempt_data 
		WHERE 
			user_id = $1
	`

	row := db.QueryRow(query, userId)
	if err := row.Scan(
		&attemptData.ID,
		&attemptData.UserId,
		&attemptData.Attempts,
		&attemptData.LastAttempt,
		&attemptData.CreatedAt,
		&attemptData.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to retrieve mpin attempt data: %w", err)
	}

	return attemptData, nil
}

func CreateMpinAttempts(db *sql.DB, userId string) (*MpinAttemptData, error) {
	attemptData := NewMpinAttemptData()
	attemptData.UserId = userId
	attemptData.Attempts = 0
	attemptData.LastAttempt = sql.NullTime{Time: time.Now(), Valid: true}

	insertQuery := `
		INSERT INTO mpin_attempt_data (
			user_id, 
			attempts, 
			last_attempt, 
			created_at, 
			updated_at
		) VALUES (
			$1, $2, $3, NOW(), NOW()
		)
	`
	if _, err := db.Exec(
		insertQuery,
		userId,
		attemptData.Attempts,
		attemptData.LastAttempt.Time,
	); err != nil {
		return nil, fmt.Errorf("failed to insert mpin attempt data: %w", err)
	}

	return attemptData, nil
}

func UpdateMpinAttempts(db *sql.DB, data *MpinAttemptData) error {
	query := `
		UPDATE mpin_attempt_data 
		SET 
			attempts = $1, 
			last_attempt = $2, 
			updated_at = NOW() 
		WHERE 
			user_id = $3
	`

	_, err := db.Exec(
		query,
		data.Attempts,
		data.LastAttempt.Time,
		data.UserId,
	)
	if err != nil {
		return fmt.Errorf("failed to update mpin attempts for user_id %s: %w", data.UserId, err)
	}

	return nil
}

func ResetMpinAttempts(db *sql.DB, userId string) error {
	query := `
		UPDATE mpin_attempt_data 
		SET 
			attempts = 0, 
			last_attempt = NOW(), 
			updated_at = NOW() 
		WHERE 
			user_id = $1
	`

	_, err := db.Exec(
		query,
		userId,
	)
	if err != nil {
		return fmt.Errorf("failed to reset mpin attempts for user_id %s: %w", userId, err)
	}

	return nil
}

func UpdateMpinAfterReset(db *sql.DB, encryptedMpin, userId string) error {
	query := `
	    UPDATE mpin_data
		SET 
		    mpin = $1, 
			updated_at = CURRENT_TIMESTAMP 
	    WHERE 
		    user_id = $2
	`

	_, err := db.Exec(
		query,
		encryptedMpin,
		userId,
	)
	if err != nil {
		return fmt.Errorf("failed to update MPIN after resetting MPIN for user_id %s: %w", userId, err)
	}

	return nil
}

func UpdateMpinAfterForgottenResetByUserId(db *sql.DB, encryptedMpin string, userId string) error {
	query := `
        UPDATE mpin_data
        SET 
            mpin = $1, 
            updated_at = NOW() 
        WHERE 
            user_id = $2
    `
	_, err := db.Exec(
		query,
		encryptedMpin,
		userId,
	)
	if err != nil {
		return fmt.Errorf("failed to update MPIN for user_id %s: %w", userId, err)
	}

	return nil
}
