-- +migrate Up

-- Write your migration here
-- drop existing foreign key constaints
ALTER TABLE menu_items DROP CONSTRAINT menu_items_menu_id_fkey;

-- re-add with cascade
ALTER TABLE menu_items
ADD CONSTRAINT menu_items_menu_id_fkey
FOREIGN KEY (menu_id) REFERENCES menus(id) ON DELETE CASCADE;
