-- +goose Up

ALTER TABLE user_client_data
ADD COLUMN if not exists is_deleted BOOLEAN DEFAULT FALSE;

-- +goose Down

ALTER TABLE user_client_data
DROP COLUMN if exists is_deleted;
