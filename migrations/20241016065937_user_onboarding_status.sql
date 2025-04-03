-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE if not exists onboarding_stages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    stage_order INT NOT NULL,
    stage_name VARCHAR(50) NOT NULL
);

INSERT INTO onboarding_stages (stage_order, stage_name) VALUES
(1, 'SIM_VERIFICATION'),
(2, 'KYC_CONSENT'),
(3, 'DEMOGRAPHIC_FETCH'),
(4, 'ACCOUNT_CREATION'),
(5, 'DEBIT_CARD_CONSENT'),
(6, 'DEBIT_CARD_PAYMENT'),
(7, 'DEBIT_CARD_GENERATION'),
(8, 'UPI_GENERATION'),
(9, 'UPI_PIN_SETUP'),
(10, 'M_PIN_SETUP');


CREATE TABLE if not exists onboarding_stage_steps (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    step_name VARCHAR(50) NOT NULL,
    stage_id UUID NOT NULL,
    step_order INT NOT NULL,
    FOREIGN KEY (stage_id) REFERENCES onboarding_stages(id)
);

CREATE INDEX if not exists idx_onboarding_stage_steps_stage_id ON onboarding_stage_steps(stage_id);
CREATE INDEX if not exists idx_onboarding_stage_steps_step_name ON onboarding_stage_steps(step_name);
CREATE INDEX if not exists idx_onboarding_stage_steps_order ON onboarding_stage_steps(step_order);


INSERT INTO onboarding_stage_steps (step_name, stage_id, step_order) VALUES
('AUTHORIZATION', (SELECT id FROM onboarding_stages WHERE stage_name = 'SIM_VERIFICATION'), 1),
('PERSONAL_DETAILS', (SELECT id FROM onboarding_stages WHERE stage_name = 'SIM_VERIFICATION'), 2),
('AGENT_URL', (SELECT id FROM onboarding_stages WHERE stage_name = 'KYC_CONSENT'), 1),
('KYC_CALLBACK', (SELECT id FROM onboarding_stages WHERE stage_name = 'KYC_CONSENT'), 2),
('ACCOUNT_CALLBACK', (SELECT id FROM onboarding_stages WHERE stage_name = 'ACCOUNT_CREATION'), 1);


CREATE TABLE if not exists user_onboarding_status (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id VARCHAR(50) NOT NULL,
    current_stage_id UUID,
    current_step_id UUID DEFAULT NULL,
    is_sim_verification_complete BOOLEAN NOT NULL DEFAULT FALSE,
    is_kyc_consent_complete BOOLEAN NOT NULL DEFAULT FALSE,
    is_demographic_fetch_complete BOOLEAN NOT NULL DEFAULT FALSE,
    is_account_creation_complete BOOLEAN NOT NULL DEFAULT FALSE,
    is_debit_card_consent_complete BOOLEAN NOT NULL DEFAULT FALSE,
    is_debit_card_payment_complete BOOLEAN NOT NULL DEFAULT FALSE,
    is_debit_card_generation_complete BOOLEAN NOT NULL DEFAULT FALSE,
    is_upi_generation_complete BOOLEAN NOT NULL DEFAULT FALSE,
    is_upi_pin_setup_complete BOOLEAN NOT NULL DEFAULT FALSE,
    is_m_pin_setup_complete BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);


CREATE UNIQUE INDEX if not exists idx_user_onboarding_status_user_id ON user_onboarding_status(user_id);
CREATE INDEX if not exists idx_user_onboarding_status_current_stage_id ON user_onboarding_status(current_stage_id);

-- +goose Down
DROP TABLE IF EXISTS user_onboarding_status;
DROP TABLE IF EXISTS onboarding_stage_steps;
DROP TABLE IF EXISTS onboarding_stages;

