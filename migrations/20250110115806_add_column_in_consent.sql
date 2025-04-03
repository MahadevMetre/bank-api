-- +goose Up
-- +goose StatementBegin
ALTER TABLE kyc_consent
ADD COLUMN if not exists aadhar2_consent BOOLEAN DEFAULT FALSE,
ADD COLUMN if not exists nomination_consent BOOLEAN DEFAULT FALSE,
ADD COLUMN if not exists location_consent BOOLEAN DEFAULT FALSE,
ADD COLUMN if not exists privacy_policy_consent BOOLEAN DEFAULT FALSE,
ADD COLUMN if not exists terms_and_condition BOOLEAN DEFAULT FALSE,
Add COLUMN if not exists address_change_consent BOOLEAN DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE kyc_consent
DROP COLUMN if exists aadhar2_consent,
DROP COLUMN if exists nomination_consent,
DROP COLUMN if exists location_consent,
DROP COLUMN if exists privacy_policy_consent,
DROP COLUMN if exists terms_and_condition,
DROP COLUMN if exists address_change_consent;
-- +goose StatementEnd