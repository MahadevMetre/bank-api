-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_onboarding_status
ADD COLUMN if not exists is_account_audit_complete BOOLEAN DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_onboarding_status
DROP COLUMN if exists is_account_audit_complete;
-- +goose StatementEnd
