-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE if not exists user_client_data (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id VARCHAR(100) NOT NULL,
    client_id VARCHAR(100),
    server_id VARCHAR(100),
    trans_id VARCHAR(100) DEFAULT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE INDEX if not EXISTS idx_user_client_data_user_id ON user_client_data(user_id);

-- +goose Down
DROP TABLE IF EXISTS user_client_data;