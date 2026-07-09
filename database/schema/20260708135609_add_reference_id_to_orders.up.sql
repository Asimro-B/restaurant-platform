-- +migrate Up

-- Write your migration here
ALTER TABLE orders ADD COLUMN reference_id VARCHAR(100) UNIQUE;
