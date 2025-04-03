-- +goose Up
-- +goose StatementBegin
-- add txn_identifier column in beneficiaries table to reuse the transaction identifier for the beneficiary transaction
ALTER TABLE beneficiaries
ADD COLUMN IF NOT EXISTS txn_identifier VARCHAR(50) DEFAULT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE beneficiaries
DROP COLUMN IF EXISTS txn_identifier;
-- +goose StatementEnd
