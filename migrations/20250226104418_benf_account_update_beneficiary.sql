-- +goose Up
-- +goose StatementBegin
ALTER TABLE beneficiaries
    ADD COLUMN benf_account VARCHAR(25) DEFAULT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE beneficiaries
    DROP COLUMN benf_account;
-- +goose StatementEnd
