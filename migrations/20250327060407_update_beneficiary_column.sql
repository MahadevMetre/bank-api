-- +goose Up
-- +goose StatementBegin
-- Update the beneficiaries table to make the benf_account column form null to not null
ALTER TABLE beneficiaries
    ALTER COLUMN benf_account SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE beneficiaries
    ALTER COLUMN benf_account DROP NOT NULL;
-- +goose StatementEnd
