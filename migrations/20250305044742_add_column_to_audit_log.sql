-- +goose Up
-- +goose StatementBegin
ALTER TABLE audit_logs
    ADD COLUMN device_id VARCHAR(120) DEFAULT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE audit_logs
    DROP COLUMN device_id;
-- +goose StatementEnd