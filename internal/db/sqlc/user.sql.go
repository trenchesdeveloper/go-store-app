// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: user.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (first_name, last_name, email, password, phone, code, expiry, verified, user_type, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
RETURNING id, first_name, last_name, email, phone, code, expiry, verified, user_type, created_at, updated_at
`

type CreateUserParams struct {
	FirstName string           `json:"first_name"`
	LastName  string           `json:"last_name"`
	Email     string           `json:"email"`
	Password  string           `json:"password"`
	Phone     pgtype.Text      `json:"phone"`
	Code      pgtype.Text      `json:"code"`
	Expiry    pgtype.Timestamp `json:"expiry"`
	Verified  bool             `json:"verified"`
	UserType  UserType         `json:"user_type"`
}

type CreateUserRow struct {
	ID        int32            `json:"id"`
	FirstName string           `json:"first_name"`
	LastName  string           `json:"last_name"`
	Email     string           `json:"email"`
	Phone     pgtype.Text      `json:"phone"`
	Code      pgtype.Text      `json:"code"`
	Expiry    pgtype.Timestamp `json:"expiry"`
	Verified  bool             `json:"verified"`
	UserType  UserType         `json:"user_type"`
	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (CreateUserRow, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.FirstName,
		arg.LastName,
		arg.Email,
		arg.Password,
		arg.Phone,
		arg.Code,
		arg.Expiry,
		arg.Verified,
		arg.UserType,
	)
	var i CreateUserRow
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.Phone,
		&i.Code,
		&i.Expiry,
		&i.Verified,
		&i.UserType,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteUser, id)
	return err
}

const getUser = `-- name: GetUser :one
SELECT id, first_name, last_name, email, password, phone, code, expiry, verified, user_type, created_at, updated_at
FROM users
WHERE id = $1
`

func (q *Queries) GetUser(ctx context.Context, id int32) (User, error) {
	row := q.db.QueryRow(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.Password,
		&i.Phone,
		&i.Code,
		&i.Expiry,
		&i.Verified,
		&i.UserType,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, first_name, last_name, email, password, phone, code, expiry, verified, user_type, created_at, updated_at
FROM users
WHERE email = $1
LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.Password,
		&i.Phone,
		&i.Code,
		&i.Expiry,
		&i.Verified,
		&i.UserType,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, first_name, last_name, email, password, phone, code, expiry, verified, user_type, created_at, updated_at
FROM users
ORDER BY id
LIMIT $1 OFFSET $2
`

type ListUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error) {
	rows, err := q.db.Query(ctx, listUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Password,
			&i.Phone,
			&i.Code,
			&i.Expiry,
			&i.Verified,
			&i.UserType,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUser = `-- name: UpdateUser :one
UPDATE users
SET first_name = $2, last_name = $3, email = $4, password = $5, phone = $6, code = $7, expiry = $8, verified = $9, user_type = $10, updated_at = NOW()
WHERE id = $1
RETURNING id, first_name, last_name, email, password, phone, code, expiry, verified, user_type, created_at, updated_at
`

type UpdateUserParams struct {
	ID        int32            `json:"id"`
	FirstName string           `json:"first_name"`
	LastName  string           `json:"last_name"`
	Email     string           `json:"email"`
	Password  string           `json:"password"`
	Phone     pgtype.Text      `json:"phone"`
	Code      pgtype.Text      `json:"code"`
	Expiry    pgtype.Timestamp `json:"expiry"`
	Verified  bool             `json:"verified"`
	UserType  UserType         `json:"user_type"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUser,
		arg.ID,
		arg.FirstName,
		arg.LastName,
		arg.Email,
		arg.Password,
		arg.Phone,
		arg.Code,
		arg.Expiry,
		arg.Verified,
		arg.UserType,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.Password,
		&i.Phone,
		&i.Code,
		&i.Expiry,
		&i.Verified,
		&i.UserType,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateUserCodeAndExpiry = `-- name: UpdateUserCodeAndExpiry :one
UPDATE users
SET code = $2, expiry = $3, updated_at = NOW()
WHERE id = $1
    RETURNING id, first_name, last_name, email, password, phone, code, expiry, verified, user_type, created_at, updated_at
`

type UpdateUserCodeAndExpiryParams struct {
	ID     int32            `json:"id"`
	Code   pgtype.Text      `json:"code"`
	Expiry pgtype.Timestamp `json:"expiry"`
}

func (q *Queries) UpdateUserCodeAndExpiry(ctx context.Context, arg UpdateUserCodeAndExpiryParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUserCodeAndExpiry, arg.ID, arg.Code, arg.Expiry)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.Password,
		&i.Phone,
		&i.Code,
		&i.Expiry,
		&i.Verified,
		&i.UserType,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateUserVerified = `-- name: UpdateUserVerified :one
UPDATE users
SET verified = $2, updated_at = NOW()
WHERE id = $1
    RETURNING id, first_name, last_name, email, password, phone, code, expiry, verified, user_type, created_at, updated_at
`

type UpdateUserVerifiedParams struct {
	ID       int32 `json:"id"`
	Verified bool  `json:"verified"`
}

func (q *Queries) UpdateUserVerified(ctx context.Context, arg UpdateUserVerifiedParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUserVerified, arg.ID, arg.Verified)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.Password,
		&i.Phone,
		&i.Code,
		&i.Expiry,
		&i.Verified,
		&i.UserType,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}