-- +goose Up
-- +goose StatementBegin
ALTER TABLE device_data
ADD COLUMN if not exists is_sim_verified BOOLEAN DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE device_data
DROP COLUMN if exists is_sim_verified;
-- +goose StatementEnd
