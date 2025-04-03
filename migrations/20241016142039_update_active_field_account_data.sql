-- +goose Up
ALTER TABLE account_data
ADD COLUMN if not exists is_active BOOLEAN DEFAULT FALSE;

-- +goose Down
ALTER TABLE account_data
DROP COLUMN if exists is_active;