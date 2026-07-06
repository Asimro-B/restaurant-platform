-- +migrate Up

-- Write your migration here
CREATE TABLE IF NOT EXISTS tables (
    id          BIGSERIAL PRIMARY KEY,
    tenant_id   BIGINT NOT NULL REFERENCES tenants(id),
    name        VARCHAR(100) NOT NULL,
    capacity    INT NOT NULL DEFAULT 2,
    status      VARCHAR(50) NOT NULL DEFAULT 'available',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
