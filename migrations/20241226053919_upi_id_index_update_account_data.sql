-- +goose Up
-- +goose StatementBegin
CREATE INDEX if not EXISTS idx_account_data_upi_id ON account_data(upi_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_account_data_upi_id;
-- +goose StatementEnd
