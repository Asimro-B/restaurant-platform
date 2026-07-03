-- name: CreateTenant :one
INSERT INTO tenants (
    name, slug, status
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: ListTenants :many
SELECT *
FROM tenants
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountTenants :one
SELECT COUNT(*)
FROM tenants
WHERE deleted_at IS NULL;

-- name: GetTenantByID :one
SELECT *
FROM tenants
WHERE id = $1
  AND deleted_at IS NULL;

-- name: GetTenantBySlug :one
SELECT *
FROM tenants
WHERE slug = $1
  AND deleted_at IS NULL;

-- name: UpdateTenant :one
UPDATE tenants
SET
    name = $2,
    slug = $3,
    status = $4,
    updated_at = NOW()
WHERE id = $1
  AND deleted_at IS NULL
RETURNING *;

-- name: RestoreTenant :one
UPDATE tenants
SET
    deleted_at = NULL,
    updated_at = NOW()
WHERE id = $1
  AND deleted_at IS NOT NULL
RETURNING *;

-- name: DeleteTenant :one
UPDATE tenants
SET
    deleted_at = NOW(),
    updated_at = NOW()
WHERE id = $1
  AND deleted_at IS NULL
RETURNING *;
