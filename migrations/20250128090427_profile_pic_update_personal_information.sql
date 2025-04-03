-- +goose Up
-- +goose StatementBegin
ALTER TABLE personal_information ADD COLUMN if not exists profile_pic JSONB DEFAULT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE personal_information DROP COLUMN if exists profile_pic;
-- +goose StatementEnd
