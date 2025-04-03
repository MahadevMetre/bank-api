-- +goose Up
-- +goose StatementBegin
CREATE TABLE if not exists nominees (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    account_data_id UUID,
    nom_name VARCHAR(100),
    nom_applicant_id VARCHAR(50),
    nom_req_type VARCHAR(50),
    nom_cbs_status VARCHAR(50),
    nom_updated_time VARCHAR(50),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS nominees;
-- +goose StatementEnd
