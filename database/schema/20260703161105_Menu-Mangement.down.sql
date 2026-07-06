-- +migrate Down

-- Write your rollback here
DROP TABLE IF EXISTS menus CASCADE;

DROP TABLE IF EXISTS menu_categories CASCADE;

DROP TABLE IF EXISTS menu_items CASCADE;