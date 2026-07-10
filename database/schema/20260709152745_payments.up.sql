CREATE TABLE IF NOT EXISTS payments (
    id              BIGSERIAL PRIMARY KEY,
    tenant_id       BIGINT NOT NULL REFERENCES tenants(id),
    order_id        BIGINT NOT NULL REFERENCES orders(id),
    amount          NUMERIC(10,2) NOT NULL,
    payment_method  VARCHAR(50) NOT NULL, -- cash | card | telebirr
    status          VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending | completed | failed
    reference       VARCHAR(255),         -- external reference (TeleBirr transaction ID, card ref)
    notes           TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
