-- +goose Up
-- +goose StatementBegin
ALTER TABLE device_data
ADD COLUMN if not exists is_active BOOLEAN DEFAULT TRUE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE device_data
DROP COLUMN if exists is_active;
-- +goose StatementEnd
