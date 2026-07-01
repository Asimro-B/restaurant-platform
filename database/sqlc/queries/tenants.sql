-- name: CreateTenant :one
INSERT INTO tenants (
    name, slug, status
) VALUES (
    $1, $2, $3
)
RETURNING *;