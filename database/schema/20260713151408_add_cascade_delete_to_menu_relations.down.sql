-- +migrate Down

-- Write your rollback here
-- drop to original constraints without cascade
ALTER TABLE menu_categories DROP CONSTRAINT menu_categories_menu_id_fkey;
ALTER TABLE menu_items DROP CONSTRAINT menu_items_category_id_fkey;

ALTER TABLE menu_categories
ADD CONSTRAINT menu_categories_menu_id_fkey
FOREIGN KEY (menu_id) REFERENCES menus(id);

ALTER TABLE menu_items
ADD CONSTRAINT menu_items_category_id_fkey
FOREIGN KEY (category_id) REFERENCES menu_categories(id);