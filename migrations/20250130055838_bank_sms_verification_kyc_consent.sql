-- +goose Up
-- +goose StatementBegin
ALTER TABLE kyc_consent ADD COLUMN bank_sms_verification_status BOOLEAN DEFAULT false;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE kyc_consent DROP COLUMN bank_sms_verification_status;
-- +goose StatementEnd
