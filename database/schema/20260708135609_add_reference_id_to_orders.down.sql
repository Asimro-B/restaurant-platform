-- +migrate Down

-- Write your rollback here
ALTER TABLE orders DROP COLUMN IF EXISTS reference_id;