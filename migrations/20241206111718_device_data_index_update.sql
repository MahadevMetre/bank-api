-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_device_data_device_id ON device_data (device_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_device_data_device_id;
-- +goose StatementEnd
