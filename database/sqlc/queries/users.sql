-- name: CreateUser :one
INSERT INTO users (
    tenant_id,
    email,
    password_hash,
    role,
    first_name,
    last_name,
    location,
    mobile_phone,
    phone
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email=$1 AND tenant_id=$2;

-- name: ListUsers :many
SELECT *
FROM users
WHERE tenant_id = sqlc.arg(tenant_id)
  AND (
    sqlc.arg(search)::text = ''
    OR email ILIKE '%' || sqlc.arg(search)::text || '%'
    OR COALESCE(first_name, '') ILIKE '%' || sqlc.arg(search)::text || '%'
    OR COALESCE(last_name, '') ILIKE '%' || sqlc.arg(search)::text || '%'
    OR COALESCE(location, '') ILIKE '%' || sqlc.arg(search)::text || '%'
    OR COALESCE(mobile_phone, '') ILIKE '%' || sqlc.arg(search)::text || '%'
    OR COALESCE(phone, '') ILIKE '%' || sqlc.arg(search)::text || '%'
  )
  AND (
    sqlc.arg(role)::text = ''
    OR role = sqlc.arg(role)::text
  )
ORDER BY
    CASE WHEN sqlc.arg(sort_by)::text = 'email' AND sqlc.arg(sort_order)::text = 'asc' THEN email END ASC,
    CASE WHEN sqlc.arg(sort_by)::text = 'email' AND sqlc.arg(sort_order)::text = 'desc' THEN email END DESC,
    CASE WHEN sqlc.arg(sort_by)::text = 'role' AND sqlc.arg(sort_order)::text = 'asc' THEN role END ASC,
    CASE WHEN sqlc.arg(sort_by)::text = 'role' AND sqlc.arg(sort_order)::text = 'desc' THEN role END DESC,
    CASE WHEN sqlc.arg(sort_by)::text = 'first_name' AND sqlc.arg(sort_order)::text = 'asc' THEN first_name END ASC,
    CASE WHEN sqlc.arg(sort_by)::text = 'first_name' AND sqlc.arg(sort_order)::text = 'desc' THEN first_name END DESC,
    CASE WHEN sqlc.arg(sort_by)::text = 'created_at' AND sqlc.arg(sort_order)::text = 'asc' THEN created_at END ASC,
    CASE WHEN sqlc.arg(sort_by)::text = 'created_at' AND sqlc.arg(sort_order)::text = 'desc' THEN created_at END DESC,
    CASE WHEN sqlc.arg(sort_by)::text = 'updated_at' AND sqlc.arg(sort_order)::text = 'asc' THEN updated_at END ASC,
    CASE WHEN sqlc.arg(sort_by)::text = 'updated_at' AND sqlc.arg(sort_order)::text = 'desc' THEN updated_at END DESC,
    created_at DESC,
    id DESC
LIMIT sqlc.arg(limit_count) OFFSET sqlc.arg(offset_count);

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1
  AND tenant_id = $2;

-- name: UpdateUser :one
UPDATE users
SET
    email = $3,
    password_hash = $4,
    role = $5,
    first_name = $6,
    last_name = $7,
    location = $8,
    mobile_phone = $9,
    phone = $10,
    updated_at = NOW()
WHERE id = $1
  AND tenant_id = $2
RETURNING *;

-- name: DeleteUser :one
DELETE FROM users
WHERE id = $1
  AND tenant_id = $2
RETURNING *;
