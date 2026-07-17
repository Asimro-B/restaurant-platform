-- name: CreateMenu :one
INSERT INTO menus (
    tenant_id,
    name,
    description,
    is_active
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: ListMenus :many
SELECT *
FROM menus
WHERE tenant_id = $1
  AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetMenuByID :one
SELECT *
FROM menus
WHERE id = $1
  AND tenant_id = $2
  AND deleted_at IS NULL;

-- name: UpdateMenu :one
UPDATE menus
SET
    name = $3,
    description = $4,
    is_active = $5,
    updated_at = NOW()
WHERE id = $1
  AND tenant_id = $2
RETURNING *;

-- name: DeleteMenu :exec
UPDATE menus
SET deleted_at = NOW()
WHERE id = $1
  AND tenant_id = $2;

-- name: RestoreMenu :exec
UPDATE menus
SET deleted_at = NULL
WHERE id = $1
  AND tenant_id = $2;

-- name: CountMenus :one
SELECT COUNT(*)
FROM menus
WHERE tenant_id = $1
  AND deleted_at IS NULL;

-- name: CreateMenuCategory :one
INSERT INTO menu_categories (
  tenant_id,
  menu_id,
  name,
  description,
  sort_order,
  is_active
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: ListMenuCategories :many
SELECT *
FROM menu_categories
WHERE tenant_id = $1
  AND menu_id = $2
  AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: GetMenuCategoryByID :one
SELECT *
FROM menu_categories
WHERE tenant_id = $1
  AND menu_id = $2
  AND id = $3
  AND deleted_at IS NULL;

-- name: UpdateMenuCategory :one
UPDATE menu_categories
SET
  name = $4,
  description = $5,
  sort_order = $6,
  is_active = $7,
  updated_at = NOW()
WHERE id = $1
  AND menu_id = $2
  AND tenant_id = $3
RETURNING *;

-- name: DeleteMenuCategory :exec
UPDATE menu_categories
SET deleted_at = NOW()
WHERE id = $1
  AND menu_id = $2
  AND tenant_id = $3;

-- name: RestoreMenuCategory :exec
UPDATE menu_categories
SET deleted_at = NULL
WHERE id = $1
  AND menu_id = $2
  AND tenant_id = $3;

-- name: CountMenuCategories :one
SELECT COUNT(*)
FROM menu_categories
WHERE menu_id = $1
  AND tenant_id = $2
  AND deleted_at IS NULL;

-- name: CreateMenuItem :one
INSERT INTO menu_items (
  tenant_id,
  category_id,
  menu_id,
  name,
  description,
  price,
  is_available
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: ListMenuItems :many
SELECT *
FROM menu_items
WHERE tenant_id = $1
  AND category_id = $2
  AND menu_id = $3
  AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $4 OFFSET $5;

-- name: GetMenuItemByID :one
SELECT * FROM menu_items
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL;

-- name: UpdateMenuItem :one
UPDATE menu_items
SET
  name = $5,
  description = $6,
  price = $7,
  is_available = $8,
  updated_at = NOW()
WHERE id = $1
  AND category_id = $2
  AND menu_id = $3
  AND tenant_id = $4
RETURNING *;

-- name: DeleteMenuItem :exec
UPDATE menu_items
SET deleted_at = NOW()
WHERE id = $1
  AND category_id = $2
  AND menu_id = $3
  AND tenant_id = $4;

-- name: RestoreMenuItem :exec
UPDATE menu_items
SET deleted_at = NULL
WHERE id = $1
  AND category_id = $2
  AND menu_id = $3
  AND tenant_id = $4;

-- name: CountMenuItems :one
SELECT COUNT(*)
FROM menu_items
WHERE category_id = $1
  AND menu_id = $2
  AND tenant_id = $3
  AND deleted_at IS NULL;