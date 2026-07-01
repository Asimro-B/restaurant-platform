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
