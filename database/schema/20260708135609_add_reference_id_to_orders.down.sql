-- +migrate Down

-- Write your rollback here
ALTER TABLE orders DROP COLUMN reference_id;