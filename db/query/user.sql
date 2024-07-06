-- name: CreateUser :one
INSERT INTO users (first_name, last_name, email, password, phone, code, expiry, verified, user_type, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
RETURNING id, first_name, last_name, email, phone, code, expiry, verified, user_type, created_at, updated_at;

-- name: GetUser :one
SELECT id, first_name, last_name, email, password, phone, code, expiry, verified, user_type, created_at, updated_at
FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT id, first_name, last_name, email, password, phone, code, expiry, verified, user_type, created_at, updated_at
FROM users
ORDER BY id
LIMIT $1 OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET first_name = $2, last_name = $3, email = $4, password = $5, phone = $6, code = $7, expiry = $8, verified = $9, user_type = $10, updated_at = NOW()
WHERE id = $1
RETURNING id, first_name, last_name, email, password, phone, code, expiry, verified, user_type, created_at, updated_at;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;