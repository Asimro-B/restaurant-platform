-- name: CreateTable :one
INSERT INTO tables (tenant_id, name, capacity, status)
VALUES ($1, $2, $3, 'available')
RETURNING *;

-- name: ListTables :many
SELECT * FROM tables
WHERE tenant_id = $1
  AND ($2::text = '' OR status = $2)
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: CountTables :one
SELECT COUNT(*) FROM tables
WHERE tenant_id = $1
  AND ($2::text = '' OR status = $2);

-- name: GetTableByID :one
SELECT * FROM tables
WHERE id = $1 AND tenant_id = $2;

-- name: UpdateTable :one
UPDATE tables
SET name = $1, capacity = $2, updated_at = now()
WHERE id = $3 AND tenant_id = $4
RETURNING *;

-- name: UpdateTableStatus :one
UPDATE tables
SET status = $1, updated_at = now()
WHERE id = $2 AND tenant_id = $3
RETURNING *;

-- name: DeleteTable :exec
DELETE FROM tables
WHERE id = $1 AND tenant_id = $2;