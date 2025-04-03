-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE if not exists device_data (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id VARCHAR(50) NOT NULL,
    device_id TEXT NOT NULL,
    sim_vendor_id TEXT NULL,
    device_token TEXT NOT NULL,
    os VARCHAR(20) NOT NULL,
    package_id VARCHAR(100) NOT NULL,
    os_version VARCHAR(20) NOT NULL,
    device_ip VARCHAR(50) NOT NULL,
    lat_long VARCHAR(100),
    server_id TEXT,
    client_id TEXT,
		created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE INDEX if not EXISTS idx_device_data_user_id ON device_data(user_id);

-- +goose Down
DROP INDEX IF EXISTS idx_device_data_user_id;
DROP TABLE device_data;
