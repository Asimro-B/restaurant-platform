-- +migrate Down

-- Write your rollback here
DROP TABLE IF EXISTS reservations CASCADE;