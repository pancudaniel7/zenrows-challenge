CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS zenrows."user"
(
    id            UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    username      TEXT UNIQUE NOT NULL CHECK (char_length(trim(username)) BETWEEN 3 AND 64),
    password_hash TEXT        NOT NULL CHECK (char_length(password_hash) >= 20),
    created_at    TIMESTAMP   NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS zenrows.device_template
(
    id              UUID PRIMARY KEY   DEFAULT gen_random_uuid(),
    name            TEXT      NOT NULL CHECK (char_length(trim(name)) BETWEEN 1 AND 100),
    device_type     TEXT      NOT NULL CHECK (device_type IN ('desktop', 'mobile')),
    width           INT CHECK (width IS NULL OR width > 0),
    height          INT CHECK (height IS NULL OR height > 0),
    user_agent      TEXT      NOT NULL CHECK (char_length(trim(user_agent)) > 0),
    country_code    CHAR(2) CHECK (country_code IS NULL OR country_code ~ '^[A-Z]{2}$'),
    default_headers JSONB CHECK (default_headers IS NULL OR jsonb_typeof(default_headers) = 'object'),
    created_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS zenrows.device_profile
(
    id             UUID PRIMARY KEY   DEFAULT gen_random_uuid(),
    user_id        UUID      NOT NULL REFERENCES zenrows."user" (id) ON DELETE CASCADE,
    template_id    UUID REFERENCES zenrows.device_template (id),
    name           TEXT      NOT NULL CHECK (char_length(trim(name)) BETWEEN 1 AND 100),
    device_type    TEXT      NOT NULL CHECK (device_type IN ('desktop', 'mobile')),
    width          INT CHECK (width IS NULL OR width > 0),
    height         INT CHECK (height IS NULL OR height > 0),
    user_agent     TEXT CHECK (user_agent IS NULL OR char_length(trim(user_agent)) > 0),
    country_code   CHAR(2) CHECK (country_code IS NULL OR country_code ~ '^[A-Z]{2}$'),
    custom_headers JSONB CHECK (custom_headers IS NULL OR jsonb_typeof(custom_headers) = 'object'),
    created_at     TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, name)
);
