-- +goose Up
-- +goose StatementBegin
ALTER TABLE account_data
ADD COLUMN if not exists profession_code VARCHAR(10) DEFAULT NULL,
ADD COLUMN if not exists annual_turn_over VARCHAR(100) DEFAULT NULL,
ADD COLUMN if not exists marital_status VARCHAR(20) DEFAULT NULL,
ADD COLUMN if not exists mother_maiden_name VARCHAR(120) DEFAULT NULL,
ADD COLUMN if not exists communication_address JSONB DEFAULT NULL,
ADD COLUMN if not exists education_qualification VARCHAR(10) DEFAULT NULL,
ADD COLUMN if not exists is_addr_same_as_aadhaar BOOLEAN DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE account_data
DROP COLUMN if exists mother_maiden_name,
DROP COLUMN if exists profession_code,
DROP COLUMN if exists annual_turn_over,
DROP COLUMN if exists marital_status,
DROP COLUMN if exists communication_address,
DROP COLUMN if exists education_qualification,
DROP COLUMN if exists is_addr_same_as_aadhaar;
-- +goose StatementEnd
