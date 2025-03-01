// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: cart.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createCart = `-- name: CreateCart :one
INSERT INTO cart (user_id, seller_id, product_id, image_url, price, name, quantity)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, user_id, seller_id, product_id, image_url, price, name, quantity, created_at, updated_at
`

type CreateCartParams struct {
	UserID    int32          `json:"user_id"`
	SellerID  int32          `json:"seller_id"`
	ProductID int32          `json:"product_id"`
	ImageUrl  string         `json:"image_url"`
	Price     pgtype.Numeric `json:"price"`
	Name      string         `json:"name"`
	Quantity  int32          `json:"quantity"`
}

func (q *Queries) CreateCart(ctx context.Context, arg CreateCartParams) (Cart, error) {
	row := q.db.QueryRow(ctx, createCart,
		arg.UserID,
		arg.SellerID,
		arg.ProductID,
		arg.ImageUrl,
		arg.Price,
		arg.Name,
		arg.Quantity,
	)
	var i Cart
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.SellerID,
		&i.ProductID,
		&i.ImageUrl,
		&i.Price,
		&i.Name,
		&i.Quantity,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteCartById = `-- name: DeleteCartById :exec
DELETE FROM cart
WHERE id = $1
`

func (q *Queries) DeleteCartById(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteCartById, id)
	return err
}

const deleteCartItems = `-- name: DeleteCartItems :exec
DELETE FROM cart
WHERE user_id = $1
`

func (q *Queries) DeleteCartItems(ctx context.Context, userID int32) error {
	_, err := q.db.Exec(ctx, deleteCartItems, userID)
	return err
}

const findCartItem = `-- name: FindCartItem :one
SELECT id, user_id, seller_id, product_id, image_url, price, name, quantity, created_at, updated_at
FROM cart
WHERE user_id = $1 AND product_id = $2
`

type FindCartItemParams struct {
	UserID    int32 `json:"user_id"`
	ProductID int32 `json:"product_id"`
}

func (q *Queries) FindCartItem(ctx context.Context, arg FindCartItemParams) (Cart, error) {
	row := q.db.QueryRow(ctx, findCartItem, arg.UserID, arg.ProductID)
	var i Cart
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.SellerID,
		&i.ProductID,
		&i.ImageUrl,
		&i.Price,
		&i.Name,
		&i.Quantity,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const findCartItems = `-- name: FindCartItems :many
SELECT id, user_id, seller_id, product_id, image_url, price, name, quantity, created_at, updated_at
FROM cart
WHERE user_id = $1
`

func (q *Queries) FindCartItems(ctx context.Context, userID int32) ([]Cart, error) {
	rows, err := q.db.Query(ctx, findCartItems, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Cart{}
	for rows.Next() {
		var i Cart
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.SellerID,
			&i.ProductID,
			&i.ImageUrl,
			&i.Price,
			&i.Name,
			&i.Quantity,
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

const getCartById = `-- name: GetCartById :one
SELECT id, user_id, seller_id, product_id, image_url, price, name, quantity, created_at, updated_at
FROM cart
WHERE id = $1
`

func (q *Queries) GetCartById(ctx context.Context, id int32) (Cart, error) {
	row := q.db.QueryRow(ctx, getCartById, id)
	var i Cart
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.SellerID,
		&i.ProductID,
		&i.ImageUrl,
		&i.Price,
		&i.Name,
		&i.Quantity,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateCart = `-- name: UpdateCart :one
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
RETURNING id, user_id, seller_id, product_id, image_url, price, name, quantity, created_at, updated_at
`

type UpdateCartParams struct {
	UserID    int32          `json:"user_id"`
	SellerID  int32          `json:"seller_id"`
	ProductID int32          `json:"product_id"`
	ImageUrl  string         `json:"image_url"`
	Price     pgtype.Numeric `json:"price"`
	Name      string         `json:"name"`
	Quantity  int32          `json:"quantity"`
	ID        int32          `json:"id"`
}

func (q *Queries) UpdateCart(ctx context.Context, arg UpdateCartParams) (Cart, error) {
	row := q.db.QueryRow(ctx, updateCart,
		arg.UserID,
		arg.SellerID,
		arg.ProductID,
		arg.ImageUrl,
		arg.Price,
		arg.Name,
		arg.Quantity,
		arg.ID,
	)
	var i Cart
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.SellerID,
		&i.ProductID,
		&i.ImageUrl,
		&i.Price,
		&i.Name,
		&i.Quantity,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
