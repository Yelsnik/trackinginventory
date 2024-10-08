// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: inventory.sql

package db

import (
	"context"
)

const createInventory = `-- name: CreateInventory :one
INSERT INTO inventory (
  item,
  serial_number,
  price,
  owner
) VALUES (
  $1, $2, $3, $4
)
RETURNING id, item, serial_number, price, owner, created_at
`

type CreateInventoryParams struct {
	Item         string `json:"item"`
	SerialNumber string `json:"serial_number"`
	Price        int64  `json:"price"`
	Owner        int64  `json:"owner"`
}

func (q *Queries) CreateInventory(ctx context.Context, arg CreateInventoryParams) (Inventory, error) {
	row := q.db.QueryRowContext(ctx, createInventory,
		arg.Item,
		arg.SerialNumber,
		arg.Price,
		arg.Owner,
	)
	var i Inventory
	err := row.Scan(
		&i.ID,
		&i.Item,
		&i.SerialNumber,
		&i.Price,
		&i.Owner,
		&i.CreatedAt,
	)
	return i, err
}

const deleteInventory = `-- name: DeleteInventory :exec
DELETE FROM inventory
WHERE id = $1
`

func (q *Queries) DeleteInventory(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteInventory, id)
	return err
}

const getInventory = `-- name: GetInventory :one
SELECT id, item, serial_number, price, owner, created_at FROM inventory
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetInventory(ctx context.Context, id int64) (Inventory, error) {
	row := q.db.QueryRowContext(ctx, getInventory, id)
	var i Inventory
	err := row.Scan(
		&i.ID,
		&i.Item,
		&i.SerialNumber,
		&i.Price,
		&i.Owner,
		&i.CreatedAt,
	)
	return i, err
}

const listInventories = `-- name: ListInventories :many
SELECT id, item, serial_number, price, owner, created_at FROM inventory
WHERE owner = $3
ORDER BY  id
LIMIT $1
OFFSET $2
`

type ListInventoriesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
	Owner  int64 `json:"owner"`
}

func (q *Queries) ListInventories(ctx context.Context, arg ListInventoriesParams) ([]Inventory, error) {
	rows, err := q.db.QueryContext(ctx, listInventories, arg.Limit, arg.Offset, arg.Owner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Inventory
	for rows.Next() {
		var i Inventory
		if err := rows.Scan(
			&i.ID,
			&i.Item,
			&i.SerialNumber,
			&i.Price,
			&i.Owner,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateInventory = `-- name: UpdateInventory :one
UPDATE inventory
  set item = $1,
   serial_number = $2
WHERE id = $3
RETURNING id, item, serial_number, price, owner, created_at
`

type UpdateInventoryParams struct {
	Item         string `json:"item"`
	SerialNumber string `json:"serial_number"`
	ID           int64  `json:"id"`
}

func (q *Queries) UpdateInventory(ctx context.Context, arg UpdateInventoryParams) (Inventory, error) {
	row := q.db.QueryRowContext(ctx, updateInventory, arg.Item, arg.SerialNumber, arg.ID)
	var i Inventory
	err := row.Scan(
		&i.ID,
		&i.Item,
		&i.SerialNumber,
		&i.Price,
		&i.Owner,
		&i.CreatedAt,
	)
	return i, err
}
