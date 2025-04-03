-- +goose Up
-- +goose StatementBegin
CREATE TABLE if not exists transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id VARCHAR(50) NOT NULL,
    transaction_id VARCHAR(60) NOT NULL,
    transaction_desc TEXT DEFAULT NULL,
    beneficiary_id VARCHAR(60) DEFAULT NULL,
    payment_mode VARCHAR(10)  NOT NULL,
    amount VARCHAR(20) DEFAULT NULL,
    utr_ref_number VARCHAR(120) DEFAULT NULL,
    otp_status VARCHAR(20) DEFAULT NULL,
    cbs_status VARCHAR(20) DEFAULT NULL,
    upi_payee_addr VARCHAR(50) DEFAULT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transactions;
-- +goose StatementEnd
