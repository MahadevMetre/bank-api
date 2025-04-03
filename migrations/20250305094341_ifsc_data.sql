-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS ifsc_data (
    id BIGSERIAL PRIMARY KEY,
    ifsc_code VARCHAR(12) NULL,
    bank_name VARCHAR(255) NULL,
    branch_name VARCHAR(255) NULL,
    branch_city VARCHAR(255) NULL,
    branch_state VARCHAR(255) NULL,
    branch_country VARCHAR(255) NULL,
    payment_mode VARCHAR(255) NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ifsc_data;
-- +goose StatementEnd