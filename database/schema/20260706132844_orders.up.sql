-- +migrate Up

-- Write your migration here
-- migration: create_orders
CREATE TABLE IF NOT EXISTS orders (
    id            BIGSERIAL PRIMARY KEY,
    tenant_id     BIGINT NOT NULL REFERENCES tenants(id),
    table_id      BIGINT NOT NULL REFERENCES tables(id),
    user_id       BIGINT NOT NULL REFERENCES users(id),  -- waiter who placed it
    status        VARCHAR(50) NOT NULL DEFAULT 'created',
    notes         TEXT,
    total_amount  NUMERIC(10,2) NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- migration: create_order_items
CREATE TABLE IF NOT EXISTS order_items (
    id              BIGSERIAL PRIMARY KEY,
    tenant_id       BIGINT NOT NULL REFERENCES tenants(id),
    order_id        BIGINT NOT NULL REFERENCES orders(id),
    menu_item_id    BIGINT NOT NULL REFERENCES menu_items(id),
    quantity        INT NOT NULL DEFAULT 1,
    unit_price      NUMERIC(10,2) NOT NULL,  -- snapshot of price at order time
    total_price     NUMERIC(10,2) NOT NULL,  -- quantity * unit_price
    notes           TEXT,                    -- "no onion", "extra spicy"
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
