-- +goose Up
-- +goose StatementBegin
ALTER TABLE personal_information ADD COLUMN if not exists pin_code VARCHAR(6) DEFAULT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE personal_information DROP COLUMN if exists pin_code;
-- +goose StatementEnd