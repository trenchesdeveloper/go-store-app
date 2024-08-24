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
SET
    name = COALESCE($2, name),
    description = COALESCE($3, description),
    category_id = COALESCE($4, category_id),
    image_url = COALESCE($5, image_url),
    price = COALESCE($6, price),
    user_id = COALESCE($7, user_id),
    stock = COALESCE($8, stock),
    updated_at = now()
WHERE id = $1
RETURNING id, name, description, category_id, image_url, price, user_id, stock, created_at, updated_at;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE id = $1;

-- name: FindSellerProducts :many
SELECT id, name, description, category_id, image_url, price, user_id, stock, created_at, updated_at
FROM products
WHERE user_id = $1;

-- name: FindProductByCategory :many
SELECT id, name, description, category_id, image_url, price, user_id, stock, created_at, updated_at
FROM products
WHERE category_id = $1;