-- +goose Up
-- +goose StatementBegin
CREATE TABLE if not exists payment_beneficiary_data (
    id UUID PRIMARY KEY,
    application_id VARCHAR(50) NOT NULL,
    service_name VARCHAR(50),
    product_type VARCHAR(20),
    sourced_by VARCHAR(20),
    callback_name VARCHAR(50),
    cbs_status JSONB,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX if not exists idx_application_id ON payment_beneficiary_data (application_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE if exists payment_beneficiary_data;
-- +goose StatementEnd