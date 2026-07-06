-- name: CreateOrder :one
INSERT INTO orders (tenant_id, table_id, user_id, notes, total_amount, status)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListOrders :many
SELECT * FROM orders
WHERE tenant_id = $1
  AND table_id = $2
  AND user_id = $3
ORDER BY created_at DESC
LIMIT $4 OFFSET $5;

-- name: CountOrders :one
SELECT COUNT(*) FROM orders
WHERE tenant_id = $1
  AND table_id = $2
  AND user_id=$3;

-- name: GetOrderByID :one
SELECT * FROM orders
WHERE id = $1 AND tenant_id = $2 AND table_id = $3 AND user_id = $4;

-- name: UpdateOrderStatus :one
UPDATE orders
SET name = status = $1, updated_at = now()
WHERE id = $2 AND tenant_id = $3 AND table_id = $4 AND user_id = $5
RETURNING *;

-- name: DeleteOrder :exec
DELETE FROM orders
WHERE id = $1 AND tenant_id = $2 AND table_id = $3 AND user_id = $4;