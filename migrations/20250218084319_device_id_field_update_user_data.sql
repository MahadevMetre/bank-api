-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_data ADD COLUMN if not exists device_id VARCHAR(120) UNIQUE;
CREATE INDEX IF NOT EXISTS idx_user_data_device_id ON user_data(device_id) WHERE device_id IS NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_user_data_device_id;
ALTER TABLE user_data DROP COLUMN if exists device_id;
-- +goose StatementEnd