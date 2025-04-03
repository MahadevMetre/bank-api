-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS mpin_attempt_data (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),  
    user_id VARCHAR(50) NOT NULL,
    attempts INT NOT NULL DEFAULT 0,
    last_attempt TIMESTAMPTZ NULL, 
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS mpin_attempt_data;  -- Drops the table if it exists during rollback
-- +goose StatementEnd
