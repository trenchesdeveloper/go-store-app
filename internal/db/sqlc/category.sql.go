// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: category.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createCategory = `-- name: CreateCategory :one
INSERT INTO categories (name, parent_id, image_url, display_order)
VALUES ($1, $2, $3, $4)
RETURNING id, name, parent_id, image_url, display_order, created_at, updated_at
`

type CreateCategoryParams struct {
	Name         string      `json:"name"`
	ParentID     pgtype.Int4 `json:"parent_id"`
	ImageUrl     pgtype.Text `json:"image_url"`
	DisplayOrder pgtype.Int4 `json:"display_order"`
}

func (q *Queries) CreateCategory(ctx context.Context, arg CreateCategoryParams) (Category, error) {
	row := q.db.QueryRow(ctx, createCategory,
		arg.Name,
		arg.ParentID,
		arg.ImageUrl,
		arg.DisplayOrder,
	)
	var i Category
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ParentID,
		&i.ImageUrl,
		&i.DisplayOrder,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteCategory = `-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = $1
`

func (q *Queries) DeleteCategory(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteCategory, id)
	return err
}

const getCategory = `-- name: GetCategory :one
SELECT id, name, parent_id, image_url, display_order, created_at, updated_at
FROM categories
WHERE id = $1
`

func (q *Queries) GetCategory(ctx context.Context, id int32) (Category, error) {
	row := q.db.QueryRow(ctx, getCategory, id)
	var i Category
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ParentID,
		&i.ImageUrl,
		&i.DisplayOrder,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listCategories = `-- name: ListCategories :many
SELECT id, name, parent_id, image_url, display_order, created_at, updated_at
FROM categories
ORDER BY id
LIMIT $1 OFFSET $2
`

type ListCategoriesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListCategories(ctx context.Context, arg ListCategoriesParams) ([]Category, error) {
	rows, err := q.db.Query(ctx, listCategories, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Category{}
	for rows.Next() {
		var i Category
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.ParentID,
			&i.ImageUrl,
			&i.DisplayOrder,
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

const updateCategory = `-- name: UpdateCategory :one
UPDATE categories
SET name = $2, parent_id = $3, image_url = $4, display_order = $5, updated_at = NOW()
WHERE id = $1
RETURNING id, name, parent_id, image_url, display_order, created_at, updated_at
`

type UpdateCategoryParams struct {
	ID           int32       `json:"id"`
	Name         string      `json:"name"`
	ParentID     pgtype.Int4 `json:"parent_id"`
	ImageUrl     pgtype.Text `json:"image_url"`
	DisplayOrder pgtype.Int4 `json:"display_order"`
}

func (q *Queries) UpdateCategory(ctx context.Context, arg UpdateCategoryParams) (Category, error) {
	row := q.db.QueryRow(ctx, updateCategory,
		arg.ID,
		arg.Name,
		arg.ParentID,
		arg.ImageUrl,
		arg.DisplayOrder,
	)
	var i Category
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.ParentID,
		&i.ImageUrl,
		&i.DisplayOrder,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}