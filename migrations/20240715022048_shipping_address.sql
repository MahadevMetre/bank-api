-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE if not exists shipping_address (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id VARCHAR(255) NOT NULL,
    document_type VARCHAR(50) NOT NULL,
    document TEXT,
    address_line_1 VARCHAR(255) NOT NULL,
    street_name VARCHAR(255) NOT NULL,
    locality VARCHAR(255) NOT NULL,
    landmark VARCHAR(255),
    city VARCHAR(100) NOT NULL,
    state VARCHAR(100) NOT NULL,
    pin_code VARCHAR(20) NOT NULL,
    country VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE INDEX if not exists idx_shipping_address_user_id ON shipping_address(user_id);

-- +goose Down
DROP INDEX IF EXISTS idx_shipping_address_user_id;
DROP TABLE shipping_address;
