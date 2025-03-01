// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: bankAccount.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createBankAccount = `-- name: CreateBankAccount :one
INSERT INTO bank_account (user_id, bank_account, swift_code, payment_type, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW())
RETURNING id, user_id, bank_account, swift_code, payment_type, created_at, updated_at
`

type CreateBankAccountParams struct {
	UserID      int64       `json:"user_id"`
	BankAccount int64       `json:"bank_account"`
	SwiftCode   pgtype.Text `json:"swift_code"`
	PaymentType pgtype.Text `json:"payment_type"`
}

func (q *Queries) CreateBankAccount(ctx context.Context, arg CreateBankAccountParams) (BankAccount, error) {
	row := q.db.QueryRow(ctx, createBankAccount,
		arg.UserID,
		arg.BankAccount,
		arg.SwiftCode,
		arg.PaymentType,
	)
	var i BankAccount
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.BankAccount,
		&i.SwiftCode,
		&i.PaymentType,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
