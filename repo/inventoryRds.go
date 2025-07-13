package repo

import (
	"context"
	"errors"
	"fmt"
	"gift/repo/models"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type InventoryRds interface {
	Save(ctx context.Context, inventories []*models.Inventory) error
	Get(ctx context.Context, inventoryId uint64) error
	Delete(ctx context.Context, inventoryId uint64) error
	StockDeduct(ctx context.Context, inventoryId uint64) error
}

func NewInventoryRds(rdsClient redis.UniversalClient) InventoryRds {
	return &inventoryRds{
		rdsClient: rdsClient,
	}
}

type inventoryRds struct {
	rdsClient redis.UniversalClient
}

func (obj *inventoryRds) Save(ctx context.Context, inventories []*models.Inventory) error {
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

func (obj *inventoryRds) StockDeduct(ctx context.Context, inventoryId uint64) error {
	var deductStock = redis.NewScript(`
  		local key   = KEYS[1]
  		local delta = tonumber(ARGV[1])
  		local now   = ARGV[2]

		local cnt   = tonumber(redis.call("HGET", key, "count") or "0")
  		if cnt < delta then
    		return -1
  		end

  		redis.call("HINCRBY", key, "count", -delta)
  		redis.call("HSET",   key, "updated_at", now)
  		return cnt - delta
	`)

	key := fmt.Sprintf("inventory:%d", inventoryId)
	now := strconv.FormatInt(time.Now().Unix(), 10)
	// Run 脚本：KEYS = [key]，ARGV = [qty, now]
	res, err := deductStock.Run(ctx, obj.rdsClient, []string{key}, 1, now).Result()
	if err != nil {
		return err
	}

	// 这里res经过lua脚本之后会返回-1或者剩余的合法的库存量
	remain := res.(int64)
	// 当库存不够的时候remain为-1
	if remain < 0 {
		return errors.New("库存不足")
	}
	return nil
}
