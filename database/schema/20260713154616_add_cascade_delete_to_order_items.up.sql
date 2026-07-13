-- +migrate Up

-- Write your migration here
-- drop existing foreign key constaints
ALTER TABLE order_items DROP CONSTRAINT order_items_menu_item_id_fkey;

-- re-add with cascade
ALTER TABLE order_items
ADD CONSTRAINT order_items_menu_item_id_fkey
FOREIGN KEY (menu_item_id) REFERENCES menu_items(id) ON DELETE CASCADE;
