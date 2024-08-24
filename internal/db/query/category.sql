
-- name: CreateCategory :one
INSERT INTO categories (name, parent_id, image_url, display_order)
VALUES ($1, $2, $3, $4)
RETURNING id, name, parent_id, image_url, display_order, created_at, updated_at;


-- name: GetCategory :one
SELECT id, name, parent_id, image_url, display_order, created_at, updated_at
FROM categories
WHERE id = $1;


-- name: ListCategories :many
SELECT id, name, parent_id, image_url, display_order, created_at, updated_at
FROM categories
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateCategory :one
UPDATE categories
SET name = $2, parent_id = $3, image_url = $4, display_order = $5, updated_at = NOW()
WHERE id = $1
RETURNING id, name, parent_id, image_url, display_order, created_at, updated_at;


-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = $1;