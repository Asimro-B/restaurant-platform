-- +migrate Up

-- Write your migration here
ALTER TABLE orders ADD COLUMN IF NOT EXISTS reference_id VARCHAR(100) UNIQUE;
