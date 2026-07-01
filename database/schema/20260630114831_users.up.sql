-- +migrate Up

-- Write your migration here
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'owner',
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    location VARCHAR(255),
    mobile_phone VARCHAR(255),
    phone VARCHAR(255),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), --Record creation time
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), -- Last update time

    -- Constraints
    UNIQUE (tenant_id, email)
)

