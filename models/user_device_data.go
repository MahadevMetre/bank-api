package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"bankapi/constants"

	"bitbucket.org/paydoh/paydoh-commons/types"
	"github.com/google/uuid"
)

// id	int4
// user_data_id	int4
// mobile_number	varchar
// signing_key	varchar
// user_created_at	timestamptz
// user_updated_at	timestamptz
// user_id	varchar
// device_id	varchar
// sim_vendor_id	varchar
// device_token	varchar
// os	varchar
// package_id	varchar
// created_at	timestamptz
// updated_at	timestamptz

type UserDeviceData struct {
	Id           int64          `json:"id"`
	MobileNumber string         `json:"mobile_number"`
	SigningKey   string         `json:"signing_key"`
	UserId       string         `json:"user_id"`
	ApplicantId  string         `json:"applicant_id"`
	DeviceId     string         `json:"device_id"`
	SimVendorId  sql.NullString `json:"sim_vendor_id"`
	DeviceToken  string         `json:"device_token"`
	DeviceIp     string         `json:"device_ip,omitempty"`
	OS           string         `json:"os,omitempty"`
	OSVersion    string         `json:"os_version,omitempty"`
	LatLong      string         `json:"lat_long,omitempty"`
	PackageId    string         `json:"package_id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type UserDevicePersonal struct {
	Id           uuid.UUID            `json:"id"`
	UserId       string               `json:"user_id"`
	MobileNumber string               `json:"mobile_number"`
	ApplicantId  string               `json:"applicant_id"`
	FirstName    string               `json:"first_name"`
	LastName     string               `json:"last_name"`
	UPIID        types.NullableString `json:"upi_id"`
	DeviceToken  string               `json:"device_token"`
	OS           string               `json:"os,omitempty"`
	OSVersion    string               `json:"os_version,omitempty"`
	PackageId    string               `json:"package_id"`
	CreatedAt    time.Time            `json:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at"`
}

func NewUserDeviceData() *UserDeviceData {
	return &UserDeviceData{}
}

func NewUserDevicePersonal() *UserDevicePersonal {
	return &UserDevicePersonal{}
}

func GetUserDeviceDataByUserId(db *sql.DB, userId string) (*UserDeviceData, error) {
	userDeviceData := NewUserDeviceData()

	row := db.QueryRow(
		`SELECT
			id,
			user_id,
			mobile_number,
			signing_key,
			device_id,
			sim_vendor_id,
			device_token,
			device_ip,
			os,
			os_version,
			lat_long,
			package_id,
			created_at,
			updated_at
			FROM user_device_data
			WHERE user_id=$1`,
		userId,
	)

	if err := row.Scan(
		&userDeviceData.Id,
		&userDeviceData.UserId,
		&userDeviceData.MobileNumber,
		&userDeviceData.SigningKey,
		&userDeviceData.DeviceId,
		&userDeviceData.SimVendorId,
		&userDeviceData.DeviceToken,
		&userDeviceData.DeviceIp,
		&userDeviceData.OS,
		&userDeviceData.OSVersion,
		&userDeviceData.LatLong,
		&userDeviceData.PackageId,
		&userDeviceData.CreatedAt,
		&userDeviceData.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return userDeviceData, nil
}

func GetUserByMobileNumber(db *sql.DB, mobileNumber string) (*UserDeviceData, error) {
	userDeviceData := NewUserDeviceData()

	row := db.QueryRow(
		`SELECT
			id,
			user_id,
			mobile_number,
			signing_key,
			device_id,
			sim_vendor_id,
			device_token,
			device_ip,
			os,
			os_version,
			lat_long,
			package_id,
			created_at,
			updated_at
			FROM user_device_data
			WHERE mobile_number=$1`,
		mobileNumber,
	)

	if err := row.Scan(
		&userDeviceData.Id,
		&userDeviceData.UserId,
		&userDeviceData.MobileNumber,
		&userDeviceData.SigningKey,
		&userDeviceData.DeviceId,
		&userDeviceData.SimVendorId,
		&userDeviceData.DeviceToken,
		&userDeviceData.DeviceIp,
		&userDeviceData.OS,
		&userDeviceData.OSVersion,
		&userDeviceData.LatLong,
		&userDeviceData.PackageId,
		&userDeviceData.CreatedAt,
		&userDeviceData.UpdatedAt,
	); err != nil {
		return nil, err
	}

	return userDeviceData, nil
}

// func UpdateDeviceData(db *sql.DB, userId, deviceIp, os, osVersion, latLong string) (*UserDeviceData, error) {

// 	userDeviceData := NewUserDeviceData()

// 	stmt, err := db.Prepare("SELECT * FROM update_device_data_v2($1, $2, $3, $4, $5)")

// 	if err != nil {
// 		return nil, err
// 	}

// 	defer stmt.Close()

// 	row := stmt.QueryRow(userId, deviceIp, os, osVersion, latLong)

// 	if err := row.Scan(
// 		&userDeviceData.Id,
// 		&userDeviceData.UserId,
// 		&userDeviceData.MobileNumber,
// 		&userDeviceData.SigningKey,
// 		&userDeviceData.DeviceId,
// 		&userDeviceData.SimVendorId,
// 		&userDeviceData.DeviceToken,
// 		&userDeviceData.DeviceIp,
// 		&userDeviceData.OS,
// 		&userDeviceData.OSVersion,
// 		&userDeviceData.LatLong,
// 		&userDeviceData.PackageId,
// 		&userDeviceData.CreatedAt,
// 		&userDeviceData.UpdatedAt,
// 	); err != nil {

// 		if err == sql.ErrNoRows {
// 			return nil, constants.ErrUserNotFound
// 		}
// 		return nil, err
// 	}

// 	return userDeviceData, nil
// }

func UpdateDeviceData(db *sql.DB, userID, deviceIP, os, osVersion, latLong string) (*UserDeviceData, error) {
	query := `
		WITH updated AS (
			UPDATE device_data
			SET
				os = COALESCE(NULLIF($3, ''), os),
				os_version = COALESCE(NULLIF($4, ''), os_version),
				device_ip = COALESCE(NULLIF($2, ''), device_ip),
				lat_long = COALESCE(NULLIF($5, ''), lat_long),
				updated_at = CURRENT_TIMESTAMP
			WHERE is_active = true AND user_id = $1
			RETURNING *
		)
		SELECT
			u.user_id, u.applicant_id, u.mobile_number, u.signing_key,
			ud.device_id, ud.sim_vendor_id, ud.device_token, ud.device_ip,
			ud.os, ud.os_version, ud.lat_long, ud.package_id, u.created_at, ud.updated_at
		FROM user_data u
		JOIN updated ud ON u.user_id = ud.user_id
	`
	var deviceData UserDeviceData

	err := db.QueryRow(query, userID, deviceIP, os, osVersion, latLong).Scan(
		&deviceData.UserId, &deviceData.ApplicantId, &deviceData.MobileNumber, &deviceData.SigningKey,
		&deviceData.DeviceId, &deviceData.SimVendorId, &deviceData.DeviceToken, &deviceData.DeviceIp,
		&deviceData.OS, &deviceData.OSVersion, &deviceData.LatLong, &deviceData.PackageId,
		&deviceData.CreatedAt, &deviceData.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no device data found for user ID %s", userID)
		}
		return nil, fmt.Errorf("failed to update device data: %v", err)
	}

	return &deviceData, nil
}

func UpdateDeviceDataV2(db *sql.DB, userID, deviceIP, os, osVersion, latLong string) (*UserDeviceData, error) {
	query := `
		WITH updated AS (
			UPDATE device_data
			SET
				os = COALESCE(NULLIF($3, ''), os),
				os_version = COALESCE(NULLIF($4, ''), os_version),
				device_ip = COALESCE(NULLIF($2, ''), device_ip),
				lat_long = COALESCE(NULLIF($5, ''), lat_long),
				updated_at = CURRENT_TIMESTAMP
			WHERE is_active = false AND user_id = $1
			RETURNING *
		)
		SELECT
			u.user_id, u.applicant_id, u.mobile_number, u.signing_key,
			ud.device_id, ud.sim_vendor_id, ud.device_token, ud.device_ip,
			ud.os, ud.os_version, ud.lat_long, ud.package_id, u.created_at, ud.updated_at
		FROM user_data u
		JOIN updated ud ON u.user_id = ud.user_id
	`
	var deviceData UserDeviceData

	err := db.QueryRow(query, userID, deviceIP, os, osVersion, latLong).Scan(
		&deviceData.UserId, &deviceData.ApplicantId, &deviceData.MobileNumber, &deviceData.SigningKey,
		&deviceData.DeviceId, &deviceData.SimVendorId, &deviceData.DeviceToken, &deviceData.DeviceIp,
		&deviceData.OS, &deviceData.OSVersion, &deviceData.LatLong, &deviceData.PackageId,
		&deviceData.CreatedAt, &deviceData.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("no device data found for user ID %s", userID)
		}
		return nil, fmt.Errorf("failed to update device data: %v", err)
	}

	return &deviceData, nil
}

func FindOneUserPersonal(db *sql.DB, userId string) (*UserDevicePersonal, error) {
	userDataPersonal := NewUserDevicePersonal()
	row := db.QueryRow(
		`SELECT
			personal_information.id,
			personal_information.user_id,
			personal_information.first_name,
			personal_information.last_name,
			personal_information.created_at,
			personal_information.updated_at,
			device_data.device_token,
			device_data.os,
			device_data.os_version,
			account_data.upi_id,
			user_data.mobile_number
		FROM personal_information
		JOIN device_data
			ON personal_information.user_id = device_data.user_id
		JOIN account_data
			ON personal_information.user_id = account_data.user_id
		JOIN user_data
			ON personal_information.user_id = user_data.user_id
		WHERE user_data.user_id = $1`, userId)

	if err := row.Scan(
		&userDataPersonal.Id,
		&userDataPersonal.UserId,
		&userDataPersonal.FirstName,
		&userDataPersonal.LastName,
		&userDataPersonal.CreatedAt,
		&userDataPersonal.UpdatedAt,
		&userDataPersonal.DeviceToken,
		&userDataPersonal.OS,
		&userDataPersonal.OSVersion,
		&userDataPersonal.UPIID,
		&userDataPersonal.MobileNumber,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrUserNotFound
		}
		return nil, err
	}

	return userDataPersonal, nil
}

func FindOneUserPersonalMobileNumber(db *sql.DB, mobileNumber string) (*UserDevicePersonal, error) {

	userDataPersonal := NewUserDevicePersonal()

	row := db.QueryRow(
		`SELECT
			personal_information.id,
			personal_information.user_id,
			personal_information.first_name,
			personal_information.last_name,
			personal_information.created_at,
			personal_information.updated_at,
			device_data.device_token,
			device_data.os,
			device_data.os_version,
			user_data.mobile_number,
			user_data.applicant_id
		FROM personal_information
		JOIN device_data
			ON personal_information.user_id = device_data.user_id
		JOIN user_data
			ON personal_information.user_id = user_data.user_id
		WHERE user_data.mobile_number = $1`, mobileNumber)

	if err := row.Scan(
		&userDataPersonal.Id,
		&userDataPersonal.UserId,
		&userDataPersonal.FirstName,
		&userDataPersonal.LastName,
		&userDataPersonal.CreatedAt,
		&userDataPersonal.UpdatedAt,
		&userDataPersonal.DeviceToken,
		&userDataPersonal.OS,
		&userDataPersonal.OSVersion,
		&userDataPersonal.MobileNumber,
		&userDataPersonal.ApplicantId,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrUserNotFound
		}
		return nil, err
	}

	return userDataPersonal, nil

}
