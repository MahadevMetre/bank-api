-- +goose Up
CREATE TABLE if not exists beneficiaries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    benf_id VARCHAR(50) NOT NULL,
    user_id VARCHAR(100) NOT NULL,
    benf_name VARCHAR(100) NOT NULL,
    benf_nickname VARCHAR(50),
    benf_mobile_number VARCHAR(12) NOT NULL,
    benf_ifsc VARCHAR(11) NOT NULL,
    benf_acct_type VARCHAR(50),
    payment_mode VARCHAR(50),
    benf_activated_time VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS beneficiaries;
