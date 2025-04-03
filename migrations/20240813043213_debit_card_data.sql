-- +goose Up
-- +goose StatementBegin
CREATE TABLE if not exists debit_card_data (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    txn_identifier VARCHAR(25),
    user_id VARCHAR(50) NOT NULL,
    cid VARCHAR(20),
    proxy_number VARCHAR(20),
    enrollment_id VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS debit_card_data;
-- +goose StatementEnd
