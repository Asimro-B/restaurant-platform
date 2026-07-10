-- name: CreatePayment :one
INSERT INTO payments(tenant_id, order_id, amount, payment_method, reference, notes)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdatePaymentStatus :one
UPDATE payments
SET status = $1, updated_at = now()
WHERE id = $2 AND tenant_id = $3
RETURNING *;

-- name: GetOrderWithItems :many
SELECT
    o.id as order_id,
    o.reference_id,
    o.table_id,
    o.status,
    o.total_amount,
    o.created_at,
    oi.id as order_item_id,
    oi.quantity,
    oi.unit_price,
    oi.total_price,
    oi.notes as item_notes,
    mi.name as menu_item_name
FROM orders o
JOIN order_items oi ON oi.order_id = o.id
JOIN menu_items mi ON mi.id = oi.menu_item_id
WHERE o.reference_id = $1 AND o.tenant_id = $2;