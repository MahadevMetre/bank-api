-- +goose Up
ALTER TABLE account_data
ADD COLUMN if not exists upi_id VARCHAR(50) DEFAULT NULL;

-- +goose Down
ALTER TABLE account_data
DROP COLUMN if exists upi_id;