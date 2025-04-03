-- +goose Up
-- +goose StatementBegin
ALTER TABLE beneficiaries
    ADD COLUMN is_active BOOLEAN DEFAULT TRUE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE beneficiaries
    DROP COLUMN is_active;
-- +goose StatementEnd
