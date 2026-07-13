-- +migrate Down

-- Write your rollback here
-- drop to original constraints without cascade
ALTER TABLE order_items DROP CONSTRAINT order_items_menu_item_id_fkey;

ALTER TABLE order_items
ADD CONSTRAINT order_items_menu_item_id_fkey
FOREIGN KEY (menu_item_id) REFERENCES menu_items(id);