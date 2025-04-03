-- +goose Up
-- +goose StatementBegin
ALTER TABLE device_data
ADD COLUMN if not exists sim_operator VARCHAR(50);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE device_data
DROP COLUMN if exists sim_operator;
-- +goose StatementEnd
