package models

import (
	"bankapi/config"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type UserClient struct {
	Id         uuid.UUID      `json:"id"`
	UserId     string         `json:"user_id"`
	ClientId   sql.NullString `json:"client_id"`
	ServerId   sql.NullString `json:"server_id"`
	TransId    sql.NullString `json:"trans_id"`
	LoginRefId sql.NullString `json:"login_ref_id"`
	IsDeleted  bool           `json:"is_deleted"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

func NewUserClient() *UserClient {
	return &UserClient{}
}

func SaveTransIdAndClientIdByUserId(db *sql.DB, userId, transId, clientId string) (*UserClient, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	if userId == "" || transId == "" || clientId == "" {
		return nil, fmt.Errorf("userId, transId, and clientId must not be empty")
	}

	userClient, err := FindOneClientIDByUserId(db, userId)
	if err != nil && err.Error() != "user not found" {
		return userClient, err
	}

	if userClient != nil {
		_, err = db.Exec(
			`UPDATE user_client_data SET trans_id=$1 WHERE user_id=$2`,
			transId, userId,
		)
		if err != nil {
			return nil, err
		}
		return userClient, nil
	}

	userClient = &UserClient{}

	err = db.QueryRow(
		`INSERT INTO user_client_data (user_id, trans_id, client_id)
			 VALUES ($1, $2, $3)
			 RETURNING id, user_id, client_id, server_id, trans_id, created_at, updated_at`,
		userId, transId, clientId,
	).Scan(
		&userClient.Id,
		&userClient.UserId,
		&userClient.ClientId,
		&userClient.ServerId,
		&userClient.TransId,
		&userClient.CreatedAt,
		&userClient.UpdatedAt,
	)

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			// This assumes you're using PostgreSQL and the pq driver
			return nil, fmt.Errorf("database error: %s, Detail: %s, Where: %s, Code: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code)
		}
		return nil, fmt.Errorf("error inserting data: %v", err)
	}

	return userClient, nil
}

func FindOneTransIdByUserId(db *sql.DB, userId string) (*UserClient, error) {
	userClient := NewUserClient()

	row := db.QueryRow(
		`SELECT
			id,
			user_id,
			client_id,
			server_id,
			trans_id,
			created_at,
			updated_at
			FROM user_client_data
			WHERE user_id=$1 
			AND 
			is_deleted=false`,
		userId,
	)

	if err := row.Scan(
		&userClient.Id,
		&userClient.UserId,
		&userClient.ClientId,
		&userClient.ServerId,
		&userClient.TransId,
		&userClient.CreatedAt,
		&userClient.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return userClient, nil
}

func FindOneClientIDByUserId(db *sql.DB, userId string) (*UserClient, error) {
	userClient := NewUserClient()

	row := db.QueryRow(
		`SELECT
			id,
			user_id,
			client_id,
			server_id,
			trans_id,
			created_at,
			updated_at,
			login_ref_id
			FROM user_client_data
			WHERE user_id=$1
			AND 
			is_deleted=false`,
		userId,
	)

	if err := row.Scan(
		&userClient.Id,
		&userClient.UserId,
		&userClient.ClientId,
		&userClient.ServerId,
		&userClient.TransId,
		&userClient.CreatedAt,
		&userClient.UpdatedAt,
		&userClient.LoginRefId,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return userClient, nil
}

func SaveServerIdByUserId(db *sql.DB, userId, serverId string, loginRefID string, transId, clientId string) (*UserClient, error) {
	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	if userId == "" || serverId == "" {
		return nil, fmt.Errorf("userId and serverId must not be empty")
	}

	userClient := NewUserClient()

	err := db.QueryRow(
		`UPDATE user_client_data
		 SET server_id = $2, login_ref_id = $3, trans_id = $4, client_id = $5, updated_at = NOW()
		 WHERE user_id = $1 AND is_deleted=false
		 RETURNING id, user_id, client_id, server_id, trans_id, created_at, updated_at`,
		userId, serverId, loginRefID, transId, clientId,
	).Scan(
		&userClient.Id,
		&userClient.UserId,
		&userClient.ClientId,
		&userClient.ServerId,
		&userClient.TransId,
		&userClient.CreatedAt,
		&userClient.UpdatedAt,
	)

	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			// This assumes you're using PostgreSQL and the pq driver
			return nil, fmt.Errorf("database error: %s, Detail: %s, Where: %s, Code: %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code)
		}
		return nil, fmt.Errorf("error updating data: %v", err)
	}

	return userClient, nil
}

func FindOneServerIDByUserId(db *sql.DB, userId string) (*UserClient, error) {
	userClient := NewUserClient()

	row := db.QueryRow(
		`SELECT
			id,
			user_id,
			client_id,
			server_id,
			trans_id,
			created_at,
			updated_at
			FROM user_client_data
			WHERE user_id=$1
			AND 
			is_deleted=false`,
		userId,
	)

	if err := row.Scan(
		&userClient.Id,
		&userClient.UserId,
		&userClient.ClientId,
		&userClient.ServerId,
		&userClient.TransId,
		&userClient.CreatedAt,
		&userClient.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return userClient, nil
}

func DeleteClientIDByUserId(userId string) (bool, error) {
	db := config.GetDB()

	query := "UPDATE user_client_data SET is_deleted = TRUE WHERE user_id = $1 AND is_deleted = FALSE"
	_, err := db.Exec(query, userId)
	if err != nil {
		return false, err
	}

	return true, nil
}
