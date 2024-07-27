-- name: CreateBankAccount :one
INSERT INTO bank_account (user_id, bank_account, swift_code, payment_type, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW())
RETURNING id, user_id, bank_account, swift_code, payment_type, created_at, updated_at;