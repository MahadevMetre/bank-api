-- +goose Up
-- +goose StatementBegin
ALTER TABLE user_client_data
ADD COLUMN if not exists login_ref_id VARCHAR(10) DEFAULT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE user_client_data
DROP COLUMN if exists login_ref_id;
-- +goose StatementEnd
