-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS kyc_audit_data (
    id UUID PRIMARY KEY,
    user_id VARCHAR(100),
    mobile_no VARCHAR(15),
    callback_name VARCHAR(100),
    applicant_id VARCHAR(50),
    sourced_by VARCHAR(100),
    vkyc_audit_status VARCHAR(50),
    product_type VARCHAR(50),
    audit_reject_reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE INDEX if not exists idx_user_id ON kyc_audit_data (user_id);
CREATE INDEX if not exists idx_mobile_no ON kyc_audit_data (mobile_no);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS kyc_audit_data;
-- +goose StatementEnd