-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE if not exists kyc_update_data (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    user_id VARCHAR(50) NOT NULL,
    astat INT NOT NULL,
    acom INT NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE INDEX if not exists idx_kyc_update_data_user_id ON kyc_update_data(user_id);

-- +goose Down
DROP TABLE if exists kyc_update_data;

DROP INDEX IF EXISTS idx_kyc_update_data_user_id;