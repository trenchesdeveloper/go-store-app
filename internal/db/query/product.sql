-- name: CreateProduct :one
INSERT INTO products (name, description, category_id, image_url, price, user_id, stock)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, name, description, category_id, image_url, price, user_id, stock, created_at, updated_at;

-- name: GetProductByID :one
SELECT id, name, description, category_id, image_url, price, user_id, stock, created_at, updated_at
FROM products
WHERE id = $1;

-- name: ListProducts :many
SELECT id, name, description, category_id, image_url, price, user_id, stock, created_at, updated_at
FROM products
ORDER BY id ASC;

-- name: UpdateProduct :one
UPDATE products
SET name = $2, description = $3, category_id = $4, image_url = $5, price = $6, user_id = $7, stock = $8, updated_at = now()
WHERE id = $1
RETURNING id, name, description, category_id, image_url, price, user_id, stock, created_at, updated_at;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;