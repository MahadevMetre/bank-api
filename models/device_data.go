package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"bankapi/config"
	"bankapi/constants"
	"bankapi/requests"
	"bankapi/security"

	"bitbucket.org/paydoh/paydoh-commons/types"
	"github.com/google/uuid"
)

type NullString struct {
	sql.NullString
}

func (s NullString) MarshalJSON() ([]byte, error) {
	if !s.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(s.String)
}

func (s *NullString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		s.String, s.Valid = "", false
		return nil
	}
	s.String, s.Valid = string(data), true
	return nil
}

type DeviceData struct {
	Id            uuid.UUID            `json:"id"`
	UserId        string               `json:"user_id"`
	DeviceId      string               `json:"device_id"`
	DeviceIp      sql.NullString       `json:"device_ip,omitempty"`
	IsActive      bool                 `json:"is_active"`
	SimVendorId   sql.NullString       `json:"sim_vendor_id"`
	SimOperator   types.NullableString `json:"sim_operator"`
	IsSimVerified bool                 `json:"is_sim_verified"`
	PackageId     string               `json:"package_id"`
	OS            sql.NullString       `json:"os,omitempty"`
	OSVersion     sql.NullString       `json:"os_version,omitempty"`
	DeviceToken   sql.NullString       `json:"device_token,omitempty"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
}

func NewDeviceData() *DeviceData {
	return &DeviceData{}
}

func (d *DeviceData) Create(userId, deviceIp, os, osVersion, key string, request *requests.DataRequest) error {
	d.UserId = userId
	encryptedDeviceId, err := security.Encrypt([]byte(request.DeviceId), []byte(key))

	if err != nil {
		return err
	}

	d.DeviceId = encryptedDeviceId

	encryptedSimVendorId, err := security.Encrypt([]byte(request.SimVendorId), []byte(key))

	if err != nil {
		return err
	}

	d.DeviceIp = sql.NullString{
		String: deviceIp,
		Valid:  deviceIp != "",
	}

	d.SimVendorId = sql.NullString{
		String: encryptedSimVendorId,
		Valid:  encryptedSimVendorId != "",
	}
	d.OS = sql.NullString{
		String: os,
		Valid:  os != "",
	}
	d.OSVersion = sql.NullString{
		String: osVersion,
		Valid:  osVersion != "",
	}
	d.PackageId = request.PackageId
	d.DeviceToken = sql.NullString{
		String: request.DeviceToken,
		Valid:  request.DeviceToken != "",
	}

	return nil
}

func InsertDevice(db *sql.DB, deviceData *DeviceData) error {
	_, err := db.Exec(
		"INSERT INTO device_data (user_id, device_id, sim_vendor_id, device_ip, os, os_version, package_id, device_token) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		deviceData.UserId,
		deviceData.DeviceId,
		deviceData.SimVendorId,
		deviceData.DeviceIp,
		deviceData.OS,
		deviceData.OSVersion,
		deviceData.PackageId,
		deviceData.DeviceToken,
	)

	if err != nil {
		return err
	}

	return nil
}

func UpdateDevice(db *sql.DB, deviceData *DeviceData, userId string) error {
	var columns []string
	var params []interface{}

	if deviceData.DeviceId != "" {
		columns = append(columns, "device_id = $"+strconv.Itoa(len(columns)+1))
		params = append(params, deviceData.DeviceId)
	}

	if deviceData.DeviceIp.Valid {
		columns = append(columns, "device_ip = $"+strconv.Itoa(len(columns)+1))
		params = append(params, deviceData.DeviceIp.String)
	}

	if deviceData.SimVendorId.Valid {
		columns = append(columns, "sim_vendor_id = $"+strconv.Itoa(len(columns)+1))
		params = append(params, deviceData.SimVendorId)
	}

	if deviceData.OS.Valid {
		columns = append(columns, "os = $"+strconv.Itoa(len(columns)+1))
		params = append(params, deviceData.OS.String)
	}

	if deviceData.OSVersion.Valid {
		columns = append(columns, "os_version = $"+strconv.Itoa(len(columns)+1))
		params = append(params, deviceData.OSVersion.String)
	}

	if deviceData.PackageId != "" {
		columns = append(columns, "package_id = $"+strconv.Itoa(len(columns)+1))
		params = append(params, deviceData.PackageId)
	}

	if deviceData.DeviceToken.String != "" {
		columns = append(columns, "device_token = $"+strconv.Itoa(len(columns)+1))
		params = append(params, deviceData.DeviceToken.String)
	}

	if deviceData.IsSimVerified {
		columns = append(columns, "is_sim_verified = $"+strconv.Itoa(len(columns)+1))
		params = append(params, deviceData.IsSimVerified)
	}

	if deviceData.SimOperator.Valid {
		columns = append(columns, "sim_operator = $"+strconv.Itoa(len(columns)+1))
		params = append(params, deviceData.SimOperator.String)
	}

	if deviceData.IsActive {
		columns = append(columns, "is_active = $"+strconv.Itoa(len(columns)+1))
		params = append(params, deviceData.IsActive)
	}

	if len(columns) > 0 {
		columns = append(columns, "updated_at = NOW()")
		query := "UPDATE device_data SET " + strings.Join(columns, ", ") + " WHERE user_id =$" + strconv.Itoa(len(params)+1)
		params = append(params, userId)

		fmt.Println("QUERY ", query)
		_, err := db.Exec(query, params...)
		if err != nil {
			return err
		}

		return nil
	}

	return nil
}

func FindOneDeviceByUserID(db *sql.DB, userId string) (*DeviceData, error) {
	deviceData := NewDeviceData()

	row := db.QueryRow(
		"SELECT id, user_id, device_id, sim_vendor_id, device_ip, os, os_version, package_id, device_token, created_at, updated_at, is_sim_verified FROM device_data WHERE user_id = $1", userId)

	if err := row.Scan(
		&deviceData.Id,
		&deviceData.UserId,
		&deviceData.DeviceId,
		&deviceData.SimVendorId,
		&deviceData.DeviceIp,
		&deviceData.OS,
		&deviceData.OSVersion,
		&deviceData.PackageId,
		&deviceData.DeviceToken,
		&deviceData.CreatedAt,
		&deviceData.UpdatedAt,
		&deviceData.IsSimVerified,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrDeviceNotFound
		}
		return nil, err
	}

	return deviceData, nil
}

func FindOneDeviceByUserIDV2(userId string) (*DeviceData, error) {
	deviceData := NewDeviceData()

	row := config.GetDB().QueryRow(
		"SELECT id, user_id, device_id, sim_vendor_id, device_ip, os, os_version, package_id, device_token, created_at, updated_at, is_sim_verified FROM device_data WHERE user_id = $1", userId)

	if err := row.Scan(
		&deviceData.Id,
		&deviceData.UserId,
		&deviceData.DeviceId,
		&deviceData.SimVendorId,
		&deviceData.DeviceIp,
		&deviceData.OS,
		&deviceData.OSVersion,
		&deviceData.PackageId,
		&deviceData.DeviceToken,
		&deviceData.CreatedAt,
		&deviceData.UpdatedAt,
		&deviceData.IsSimVerified,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrDeviceNotFound
		}
		return nil, err
	}

	return deviceData, nil
}

func GetSingleDeviceDataByUserID(userId string) (*DeviceData, error) {
	deviceData := NewDeviceData()

	row := config.GetDB().QueryRow(
		"SELECT id, user_id, device_id, sim_vendor_id, device_ip, os, os_version, package_id, device_token, created_at, updated_at, is_sim_verified FROM device_data WHERE is_active = true AND user_id = $1", userId)

	if err := row.Scan(
		&deviceData.Id,
		&deviceData.UserId,
		&deviceData.DeviceId,
		&deviceData.SimVendorId,
		&deviceData.DeviceIp,
		&deviceData.OS,
		&deviceData.OSVersion,
		&deviceData.PackageId,
		&deviceData.DeviceToken,
		&deviceData.CreatedAt,
		&deviceData.UpdatedAt,
		&deviceData.IsSimVerified,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrDeviceNotFound
		}
		return nil, err
	}

	return deviceData, nil
}

func GetSingleDeviceDataByUserIDV2(userId string) (*DeviceData, error) {
	deviceData := NewDeviceData()

	row := config.GetDB().QueryRow(
		"SELECT id, user_id, device_id, sim_vendor_id, device_ip, os, os_version, package_id, device_token, created_at, updated_at, is_sim_verified FROM device_data WHERE user_id = $1", userId)

	if err := row.Scan(
		&deviceData.Id,
		&deviceData.UserId,
		&deviceData.DeviceId,
		&deviceData.SimVendorId,
		&deviceData.DeviceIp,
		&deviceData.OS,
		&deviceData.OSVersion,
		&deviceData.PackageId,
		&deviceData.DeviceToken,
		&deviceData.CreatedAt,
		&deviceData.UpdatedAt,
		&deviceData.IsSimVerified,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, constants.ErrDeviceNotFound
		}
		return nil, err
	}

	return deviceData, nil
}

func UpdateIsActiveAndSimVerifiedStatus(isActive, isVerified bool, userId string) error {
	_, err := config.GetDB().Exec(
		"UPDATE device_data SET is_active = $1, is_sim_verified = $2 WHERE user_id = $3",
		isActive, isVerified, userId,
	)
	if err != nil {
		return err
	}

	return nil
}
