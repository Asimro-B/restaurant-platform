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
SELECT * FROM users WHERE email=$1 AND tenant_id=$2 AND deleted_at IS NULL;

-- name: ListUsers :many
SELECT * FROM users
WHERE tenant_id = $1
  AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1
  AND tenant_id = $2
  AND deleted_at IS NULL;

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

-- name: DeleteUser :exec
UPDATE users
SET deleted_at = NOW()
WHERE id = $1
  AND tenant_id = $2;

-- name: RestoreUser :exec
UPDATE users
SET deleted_at = NULL
WHERE id = $1
  AND tenant_id = $2;
