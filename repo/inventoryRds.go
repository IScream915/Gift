package repo

import (
	"context"
	"gift/repo/models"
	"github.com/redis/go-redis/v9"
)

type InventoryRds interface {
	Save(ctx context.Context, inventory *models.Inventory) error
	Get(ctx context.Context, inventoryId uint64) error
	Delete(ctx context.Context, inventoryId uint64) error
}

func NewInventoryRds(rdsClient *redis.UniversalClient) InventoryRds {
	return &inventoryRds{
		rdsClient: rdsClient,
	}
}

type inventoryRds struct {
	rdsClient *redis.UniversalClient
}

func (obj *inventoryRds) Save(ctx context.Context, inventory *models.Inventory) error {
	//TODO implement me
	panic("implement me")
}

func (obj *inventoryRds) Get(ctx context.Context, inventoryId uint64) error {
	//TODO implement me
	panic("implement me")
}

func (obj *inventoryRds) Delete(ctx context.Context, inventoryId uint64) error {
	//TODO implement me
	panic("implement me")
}
