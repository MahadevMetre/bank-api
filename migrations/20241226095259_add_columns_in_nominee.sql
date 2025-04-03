-- +goose Up
-- +goose StatementBegin
ALTER TABLE nominees
ADD COLUMN if not exists date_of_birth TEXT DEFAULT NULL,
ADD COLUMN if not exists relation TEXT DEFAULT NULL,
ADD COLUMN if not exists address_1 TEXT DEFAULT NULL,
ADD COLUMN if not exists address_2 TEXT DEFAULT NULL,
ADD COLUMN if not exists address_3 TEXT DEFAULT NULL,
ADD COLUMN if not exists city TEXT DEFAULT NULL,
ADD COLUMN if not exists nominee_mobile_number VARCHAR(50) DEFAULT NULL,
ADD COLUMN if not exists is_otp_sent BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS txn_identifier VARCHAR(50) DEFAULT NULL,
ADD COLUMN if not exists pincode TEXT DEFAULT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE nominees 
DROP COLUMN if exists date_of_birth,
DROP COLUMN if exists relation,
DROP COLUMN if exists address_1,
DROP COLUMN if exists address_2,
DROP COLUMN if exists address_3,
DROP COLUMN if exists city,
DROP COLUMN if exists nominee_mobile_number,
DROP COLUMN if exists is_otp_sent,
DROP COLUMN IF EXISTS txn_identifier,
DROP COLUMN if exists pincode;
-- +goose StatementEnd