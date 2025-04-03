package services

import (
	"bankapi/config"
	"bankapi/constants"
	"bankapi/utils"
	"context"
	"fmt"
	"time"

	"bitbucket.org/paydoh/paydoh-commons/pkg/json"
	"bitbucket.org/paydoh/paydoh-commons/pkg/task"
	"bitbucket.org/paydoh/paydoh-commons/services"
	"github.com/hibiken/asynq"
)

type AuditLogService interface {
	Save(ctx context.Context, auditLog *AuditLog) error
	AuditLogHandler(ctx context.Context, t *asynq.Task) error
}

type auditService struct {
	logger       *services.LoggerService
	taskEnqueuer task.TaskEnqueuer
}

type AuditLog struct {
	TransactionID  string    `json:"transaction_id"`
	RequestID      string    `json:"request_id"`
	UserID         string    `json:"user_id"`
	ApplicantID    string    `json:"applicant_id"`
	SourceIP       string    `json:"source_ip"`
	DeviceOS       string    `json:"device_os"`
	AppVersion     string    `json:"app_version"`
	RequestURL     string    `json:"request_url"`
	HTTPMethod     string    `json:"http_method"`
	RequestBody    string    `json:"request_body"`
	ResponseStatus int       `json:"response_status"`
	Action         string    `json:"action"`
	CreatedAt      time.Time `json:"created_at"`
	DeviceID       string    `json:"device_id"`
}

func NewAuditLogService(logger *services.LoggerService, taskEnqueuer task.TaskEnqueuer) AuditLogService {
	return &auditService{
		logger:       logger,
		taskEnqueuer: taskEnqueuer,
	}
}

func (a *auditService) insert(ctx context.Context, log *AuditLog) error {
	db := config.GetDB()
	query := `
		INSERT INTO audit_logs (
			transaction_id, request_id, user_id, applicant_id, source_ip,
			device_os, app_version, request_url, http_method, request_body,
			response_status, action, created_at, device_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13,$14)`

	_, err := db.ExecContext(ctx, query,
		log.TransactionID, log.RequestID, log.UserID, log.ApplicantID, log.SourceIP,
		log.DeviceOS, log.AppVersion, log.RequestURL, log.HTTPMethod, log.RequestBody,
		log.ResponseStatus, log.Action, log.CreatedAt, log.DeviceID)

	if err != nil {
		return fmt.Errorf("failed to insert audit log: %w", err)
	}

	return nil
}

func (a *auditService) Save(ctx context.Context, auditLog *AuditLog) error {
	auditLogData := AuditLog{
		UserID:         auditLog.UserID,
		ApplicantID:    auditLog.ApplicantID,
		TransactionID:  auditLog.TransactionID,
		RequestID:      utils.GetRequestIDFromContext(ctx),
		SourceIP:       utils.GetSourceIPFromContext(ctx),
		DeviceOS:       utils.GetDeviceOSFromContext(ctx),
		AppVersion:     utils.GetAppVersionFromContext(ctx),
		DeviceID:       utils.GetDeviceIdFromContext(ctx),
		RequestURL:     auditLog.RequestURL,
		HTTPMethod:     auditLog.HTTPMethod,
		RequestBody:    auditLog.RequestBody,
		ResponseStatus: auditLog.ResponseStatus,
		Action:         auditLog.Action,
		CreatedAt:      time.Now(),
	}

	if utils.GetRequestIDFromContext(ctx) == "" {
		auditLogData.RequestID = auditLog.RequestID
	}

	if utils.GetSourceIPFromContext(ctx) == "" {
		auditLogData.SourceIP = auditLog.SourceIP
	}

	if utils.GetDeviceOSFromContext(ctx) == "" {
		auditLogData.DeviceOS = auditLog.DeviceOS
	}

	if utils.GetAppVersionFromContext(ctx) == "" {
		auditLogData.AppVersion = auditLog.AppVersion
	}

	if utils.GetDeviceIdFromContext(ctx) == "" {
		auditLogData.DeviceID = auditLog.DeviceID
	}

	if _, _, err := a.taskEnqueuer.EnqueueNow(constants.AuditLogType, auditLogData, "default"); err != nil {
		a.logger.Logger.Error("Failed to save audit log")
		return err
	}
	return nil
}

func (a *auditService) AuditLogHandler(ctx context.Context, t *asynq.Task) error {
	var auditLog AuditLog
	if err := json.Unmarshal(t.Payload(), &auditLog); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if err := a.insert(ctx, &auditLog); err != nil {
		a.logger.Logger.Error("Failed to save audit log")
		return fmt.Errorf("failed to save audit log: %w", err)
	}

	a.logger.Logger.Info("Audit log saved successfully")
	return nil
}
