-- +goose Up
-- +goose StatementBegin
CREATE TABLE if not exists audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transaction_id VARCHAR(255) NOT NULL,
    request_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    applicant_id VARCHAR(255),
    source_ip VARCHAR(45),
    device_os VARCHAR(100),
    app_version VARCHAR(50),
    request_url TEXT,
    http_method VARCHAR(10),
    request_body TEXT,
    response_status INT,
    action VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS audit_logs;
-- +goose StatementEnd
