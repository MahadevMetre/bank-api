-- +goose Up
ALTER TABLE account_data
ADD COLUMN if not exists callback_name VARCHAR(255),
ADD COLUMN if not exists status VARCHAR(50);

-- +goose Down
ALTER TABLE account_data
DROP COLUMN if exists callback_name,
DROP COLUMN if exists status;