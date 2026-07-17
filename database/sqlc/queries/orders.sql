-- name: CreateOrder :one
INSERT INTO orders (tenant_id, table_id, user_id, notes, total_amount, status, reference_id)
VALUES ($1, $2, $3, $4, $5, 'created', $6)
RETURNING *;

-- name: ListOrders :many
SELECT * FROM orders
WHERE tenant_id = $1
  AND ($2::text = '' OR status = $2)
  AND ($3::bigint = 0 OR table_id = $3)
ORDER BY created_at DESC
LIMIT $4 OFFSET $5;

-- name: CountOrders :one
SELECT COUNT(*) FROM orders
WHERE tenant_id = $1
  AND ($2::text = '' OR status = $2)
  AND ($3::bigint = 0 OR table_id = $3);

-- name: GetOrderByID :one
SELECT * FROM orders
WHERE id = $1 AND tenant_id = $2 AND table_id = $3 AND user_id = $4;

-- name: UpdateOrderStatus :one
UPDATE orders
SET status = $1, updated_at = now()
WHERE id = $2 AND tenant_id = $3
RETURNING *;

-- name: DeleteOrder :exec
DELETE FROM orders
WHERE id = $1 AND tenant_id = $2 AND table_id = $3 AND user_id = $4;

-- name: CreateOrderItem :one
INSERT INTO order_items (tenant_id, order_id, menu_item_id, quantity, unit_price, total_price, notes)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetOrderByReferenceID :one
SELECT * FROM orders
WHERE reference_id = $1 AND tenant_id = $2;