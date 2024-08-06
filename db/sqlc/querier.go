// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"context"
)

type Querier interface {
	CreateInventory(ctx context.Context, arg CreateInventoryParams) (Inventory, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteInventory(ctx context.Context, id int64) error
	GetInventory(ctx context.Context, id int64) (Inventory, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	ListInventories(ctx context.Context, arg ListInventoriesParams) ([]Inventory, error)
	UpdateInventory(ctx context.Context, arg UpdateInventoryParams) (Inventory, error)
}

var _ Querier = (*Queries)(nil)
