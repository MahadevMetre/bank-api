-- +goose Up
-- +goose StatementBegin

ALTER TABLE transactions
ADD COLUMN if not exists upi_payee_name VARCHAR(250) DEFAULT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE transactions
DROP COLUMN if exists upi_payee_name;
-- +goose StatementEnd
