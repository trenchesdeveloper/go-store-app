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
SET
first_name = COALESCE($2, first_name),
last_name = COALESCE($3, last_name),
email = COALESCE($4, email),
phone = COALESCE($5, phone),
updated_at = NOW()
WHERE id = $1
RETURNING id, first_name, last_name, email, phone, code, expiry, verified, user_type, created_at, updated_at;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, first_name, last_name, email, password, phone, code, expiry, verified, user_type, created_at, updated_at
FROM users
WHERE email = $1
LIMIT 1;

-- name: UpdateUserCodeAndExpiry :one
UPDATE users
SET code = $2, expiry = $3, updated_at = NOW()
WHERE id = $1
    RETURNING id, first_name, last_name, email, password, phone, code, expiry, verified, user_type, created_at, updated_at;

-- name: UpdateUserVerified :one
UPDATE users
SET verified = $2, updated_at = NOW()
WHERE id = $1
    RETURNING id, first_name, last_name, email, password, phone, code, expiry, verified, user_type, created_at, updated_at;

-- name: UpdateUserToSeller :one
UPDATE users
SET user_type =$2, first_name = $3, last_name = $4, phone = $5, updated_at = NOW()
WHERE id = $1
    RETURNING id, first_name, last_name, email, password, phone, code, expiry, verified, user_type, created_at, updated_at;