-- name: CreateAddress :one
INSERT INTO address (user_id, address_line1,address_line2, city, state, post_code, country)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, user_id, address_line1, address_line2, city, state, post_code, country, created_at, updated_at;

-- name: GetAddress :one
SELECT id, user_id, address_line1, address_line2, city, state, post_code, country, created_at, updated_at
FROM address
WHERE id = $1;

-- name: ListAddresses :many
SELECT id, user_id, address_line1, address_line2, city, state, post_code, country, created_at, updated_at
FROM address
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateAddress :one
UPDATE address
SET
    user_id = COALESCE($2, user_id),
    address_line1 = COALESCE($3, address_line1),
    address_line2 = COALESCE($4, address_line2),
    city = COALESCE($5, city),
    state = COALESCE($6, state),
    post_code = COALESCE($7, post_code),
    country = COALESCE($8, country),
    updated_at = CURRENT_TIMESTAMP
WHERE user_id = $1
RETURNING id, user_id, address_line1, address_line2, city, state, post_code, country, created_at, updated_at;

-- name: DeleteAddress :exec
DELETE FROM address
WHERE id = $1;

-- name: FindAddressByUser :many
SELECT id, user_id, address_line1, address_line2, city, state, post_code, country, created_at, updated_at
FROM address
WHERE user_id = $1;

