package repo

import (
	"context"
	"gift/pkg/gormutil"
	"gift/repo/models"
)

type InventoryQueryOption struct {
	gormutil.Option
	name   string
	desc   string
	price  uint64
	count  uint64
	offset int
}

type InventoryQueryFunc func(opt *InventoryQueryOption)

func WithName(name string) InventoryQueryFunc {
	return func(opt *InventoryQueryOption) {
		opt.name = name
	}
}

func WithDesc(desc string) InventoryQueryFunc {
	return func(opt *InventoryQueryOption) {
		opt.desc = desc
	}
}

func WithPrice(price uint64) InventoryQueryFunc {
	return func(opt *InventoryQueryOption) {
		opt.price = price
	}
}

func WithCount(count uint64) InventoryQueryFunc {
	return func(opt *InventoryQueryOption) {
		opt.count = count
	}
}

func WithPagination(page, pageSize int) InventoryQueryFunc {
	return func(opt *InventoryQueryOption) {
		opt.Page, opt.PageSize = gormutil.RepairPaging(page, pageSize)
	}
}

func WithOffset(offset int) InventoryQueryFunc {
	return func(opt *InventoryQueryOption) {
		opt.offset = offset
	}
}

func WithSort(sort string) InventoryQueryFunc {
	return func(opt *InventoryQueryOption) {
		opt.Sort = sort
	}
}

func (obj *inventory) FindInventoryList(ctx context.Context, opts ...InventoryQueryFunc) ([]*models.Inventory, uint64, error) {
	option := InventoryQueryOption{}
	for _, opt := range opts {
		opt(&option)
	}
	flag := true

	query := obj.client.WithContext(ctx).Table(models.TableNameInventories)

	if option.name != "" {
		query = query.Where("name LIKE ?", "%"+option.name+"%")
	}
	if option.desc != "" {
		query = query.Where("description LIKE ?", "%"+option.desc+"%")
	}
	if option.price > 0 {
		query = query.Where("price = ?", option.price)
	}
	if option.count > 0 {
		query = query.Where("count = ?", option.count)
	}

	if option.offset >= 0 {
		// 手动指定排序
		query = query.Offset(option.offset).Limit(gormutil.DefaultPageSize)
		flag = false
	}

	if flag {
		// 自动指定, 分页和排序
		query = option.PagingAndSort(query)
	}

	records := make([]*models.Inventory, 0)
	if err := query.Find(&records).Debug().Error; err != nil {
		return nil, 0, err
	}

	return records, uint64(len(records)), nil
}
