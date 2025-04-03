-- +goose Up
-- +goose StatementBegin
ALTER TABLE personal_information
ADD COLUMN if not exists is_account_detail_email_sent BOOLEAN DEFAULT FALSE,
ADD COLUMN if not exists is_email_verified BOOLEAN DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE personal_information
DROP COLUMN if exists is_email_verified,
DROP COLUMN IF EXISTS is_account_detail_email_sent;
-- +goose StatementEnd
