-- +goose Up
-- +goose StatementBegin
CREATE TABLE address_update_data (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id VARCHAR(25),
    Req_Ref_No VARCHAR(50),
    status BOOLEAN DEFAULT FALSE,
    communication_address JSONB DEFAULT NULL,
    current_status VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS address_update_data;
-- +goose StatementEnd