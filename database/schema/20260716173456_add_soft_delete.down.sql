-- +migrate Down

-- Write your rollback here

ALTER TABLE orders DROP COLUMN deleted_at;
ALTER TABLE payments DROP COLUMN deleted_at;
ALTER TABLE users DROP COLUMN deleted_at;
ALTER TABLE menus DROP COLUMN deleted_at;
ALTER TABLE menu_categories DROP COLUMN deleted_at;
ALTER TABLE menu_items DROP COLUMN deleted_at;