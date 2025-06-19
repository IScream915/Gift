package repo

import (
	"context"
	"fmt"
	"gift/repo/models"
	"github.com/redis/go-redis/v9"
)

type InventoryRds interface {
	WriteToRedis(ctx context.Context, inventories []*models.Inventory) error
	Get(ctx context.Context, inventoryId uint64) error
	Delete(ctx context.Context, inventoryId uint64) error
}

func NewInventoryRds(rdsClient redis.UniversalClient) InventoryRds {
	return &inventoryRds{
		rdsClient: rdsClient,
	}
}

type inventoryRds struct {
	rdsClient redis.UniversalClient
}

func (obj *inventoryRds) WriteToRedis(ctx context.Context, inventories []*models.Inventory) error {
	pipe := obj.rdsClient.Pipeline()
	for _, inventory := range inventories {
		key := fmt.Sprintf("inventory:%d", inventory.ID)

		// 这里用Hash存储库存和更新时间
		pipe.HSet(ctx, key, map[string]interface{}{
			"count":      inventory.Count,
			"updated_at": inventory.UpdatedAt.Unix(),
		})
	}
	_, err := pipe.Exec(ctx)
	return err
}

func (obj *inventoryRds) Get(ctx context.Context, inventoryId uint64) error {
	//TODO implement me
	panic("implement me")
}

func (obj *inventoryRds) Delete(ctx context.Context, inventoryId uint64) error {
	//TODO implement me
	panic("implement me")
}
