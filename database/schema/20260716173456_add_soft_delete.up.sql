-- +migrate Up

-- Write your migration here

ALTER TABLE orders ADD COLUMN deleted_at TIMESTAMPTZ;
ALTER TABLE payments ADD COLUMN deleted_at TIMESTAMPTZ;
ALTER TABLE users ADD COLUMN deleted_at TIMESTAMPTZ;
ALTER TABLE menus ADD COLUMN deleted_at TIMESTAMPTZ;
ALTER TABLE menu_categories ADD COLUMN deleted_at TIMESTAMPTZ;
ALTER TABLE menu_items ADD COLUMN deleted_at TIMESTAMPTZ;
