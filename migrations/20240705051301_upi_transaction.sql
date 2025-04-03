-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS upi_transactions (
    id UUID PRIMARY KEY,
    user_id VARCHAR(50),
    transaction_id VARCHAR(50),
    cred_type VARCHAR(50),
		created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS upi_transactions;
-- +goose StatementEnd