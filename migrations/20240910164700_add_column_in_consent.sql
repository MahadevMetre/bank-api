-- +goose Up
-- +goose StatementBegin
ALTER TABLE kyc_consent
ADD COLUMN if not exists virtual_card_consent BOOLEAN DEFAULT FALSE,
ADD COLUMN if not exists physical_card_consent BOOLEAN  DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE kyc_consent
DROP COLUMN if exists virtual_card_consent,
DROP COLUMN if exists physical_card_consent;
-- +goose StatementEnd
