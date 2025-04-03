-- +goose Up
-- +goose StatementBegin
ALTER TABLE debit_card_data
ADD COLUMN physical_debitcard_txnid VARCHAR(25),
ADD COLUMN is_virtual_generated BOOLEAN DEFAULT FALSE,
ADD COLUMN is_physical_generated BOOLEAN DEFAULT FALSE,
ADD COLUMN public_key TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE debit_card_data
DROP COLUMN IF EXISTS physical_debitcard_txnid,
DROP COLUMN IF EXISTS is_virtual_generated,
DROP COLUMN IF EXISTS is_physical_generated,
DROP COLUMN IF EXISTS public_key;
-- +goose StatementEnd