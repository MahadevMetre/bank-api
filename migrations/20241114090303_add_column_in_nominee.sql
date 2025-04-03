-- +goose Up
-- +goose StatementBegin
ALTER TABLE nominees
ADD COLUMN if not exists is_verified BOOLEAN DEFAULT FALSE,
ADD COLUMN if not exists user_id VARCHAR(50) DEFAULT NULL,
ADD COLUMN if not exists is_active BOOLEAN DEFAULT FALSE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE nominees
DROP COLUMN IF EXISTS is_verified,
DROP COLUMN IF EXISTS user_id,
DROP COLUMN IF EXISTS is_active;
-- +goose StatementEnd
