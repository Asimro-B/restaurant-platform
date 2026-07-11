-- +migrate Up

-- Write your migration here

CREATE TABLE IF NOT EXISTS reservations (
    id              BIGSERIAL PRIMARY KEY,
    tenant_id       BIGINT NOT NULL REFERENCES tenants(id),
    table_id        BIGINT NOT NULL REFERENCES tables(id),
    customer_name   VARCHAR(255) NOT NULL,
    customer_phone  VARCHAR(50) NOT NULL,
    party_size      INT NOT NULL,
    reserved_at     TIMESTAMPTZ NOT NULL,  -- when the reservation is for
    status          VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending | confirmed | cancelled | completed
    notes           TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);