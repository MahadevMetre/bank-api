-- +goose Up
-- +goose StatementBegin
ALTER TABLE debit_card_data
ADD COLUMN if not exists delivery_status VARCHAR(255) DEFAULT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE debit_card_data DROP COLUMN if exists delivery_status;
-- +goose StatementEnd
