-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE if not exists user_data (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL,
    mobile_number VARCHAR(20) NOT NULL,
    applicant_id VARCHAR(50) UNIQUE,
    signing_key VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE INDEX if not exists idx_user_data_mobile_number ON user_data (mobile_number);
CREATE INDEX if not exists idx_user_data_applicant_id ON user_data (applicant_id);

-- +goose Down
DROP INDEX IF EXISTS idx_user_data_mobile_number;
DROP INDEX IF EXISTS idx_user_data_applicant_id;