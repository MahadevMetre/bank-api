-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE if not exists kyc_consent (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id VARCHAR(50) NOT NULL,
    indian_resident BOOLEAN NOT NULL,
    politically_exposed_person BOOLEAN NOT NULL,
    aadhar_consent BOOLEAN NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE INDEX if not exists idx_kyc_consent_user_id ON kyc_consent (user_id);

-- +goose Down
DROP INDEX IF EXISTS idx_kyc_consent_user_id;
DROP TABLE if exists kyc_consent;
