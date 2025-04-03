-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE if not exists account_data (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id VARCHAR(50) NOT NULL,
    account_number VARCHAR(50),
    customer_id VARCHAR(50),
    application_id VARCHAR(50),
    service_name VARCHAR(100),
    sourced_by VARCHAR(100),
    product_type VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE INDEX idx_account_data_user_id ON account_data(user_id);
CREATE INDEX idx_account_data_application_id ON account_data(application_id);

-- +goose Down
DROP INDEX IF EXISTS idx_account_data_user_id;
DROP INDEX IF EXISTS idx_account_data_application_id;
DROP TABLE IF EXISTS account_data;
