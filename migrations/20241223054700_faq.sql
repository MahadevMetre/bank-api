-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS faq_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL
);

-- insert default categories
INSERT INTO faq_categories (name) VALUES
    ('Account'),
    ('Fees and charges'),
    ('Split bill and Budget book'),
    ('Rewards'),
    ('Shopping'),
    ('Referral'),
    ('Debit card'),
    ('How to Videos');

CREATE INDEX IF NOT EXISTS idx_faq_categories_name ON faq_categories(name);

CREATE TABLE IF NOT EXISTS faqs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    question VARCHAR(255) NOT NULL,
    answer VARCHAR(255) NOT NULL,
    category_id UUID NOT NULL,
    website_only BOOLEAN NOT NULL DEFAULT TRUE,
    app_only BOOLEAN NOT NULL DEFAULT TRUE,
    video_url TEXT,
    FOREIGN KEY (category_id) REFERENCES faq_categories (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_faqs_category_id ON faqs(category_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS faqs;
DROP TABLE IF EXISTS faq_categories;
-- +goose StatementEnd
