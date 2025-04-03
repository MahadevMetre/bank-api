-- +goose Up
-- +goose StatementBegin
ALTER TABLE debit_card_data
ADD COLUMN if not exists is_permanently_blocked BOOLEAN DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE debit_card_data
DROP COLUMN if exists is_permanently_blocked;
-- +goose StatementEnd
