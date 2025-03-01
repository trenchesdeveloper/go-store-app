-- name: CreateCart :one
INSERT INTO cart (user_id, seller_id, product_id, image_url, price, name, quantity)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, user_id, seller_id, product_id, image_url, price, name, quantity, created_at, updated_at;

-- name: FindCartItems :many
SELECT id, user_id, seller_id, product_id, image_url, price, name, quantity, created_at, updated_at
FROM cart
WHERE user_id = $1;

-- name: DeleteCartById :exec
DELETE FROM cart
WHERE id = $1;

-- name: DeleteCartItems :exec
DELETE FROM cart
WHERE user_id = $1;

-- name: UpdateCart :one
UPDATE cart
SET
    user_id = COALESCE($1, user_id),
    seller_id = COALESCE($2, seller_id),
    product_id = COALESCE($3, product_id),
    image_url = COALESCE($4, image_url),
    price = COALESCE($5, price),
    name = COALESCE($6, name),
    quantity = COALESCE($7, quantity),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $8
RETURNING id, user_id, seller_id, product_id, image_url, price, name, quantity, created_at, updated_at;

-- name: GetCartById :one
SELECT id, user_id, seller_id, product_id, image_url, price, name, quantity, created_at, updated_at
FROM cart
WHERE id = $1;

-- name: FindCartItem :one
SELECT id, user_id, seller_id, product_id, image_url, price, name, quantity, created_at, updated_at
FROM cart
WHERE user_id = $1 AND product_id = $2;