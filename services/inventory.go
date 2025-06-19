package services

import (
	"context"
	"gift/assemble"
	"gift/dto"
	"gift/pkg/gormutil"
	"gift/repo"
	"gift/repo/models"
	"gorm.io/gorm"
	"sync"
)

type Inventory interface {
	Create(ctx context.Context, req *dto.CreateInventoryReq) error
	Update(ctx context.Context, req *dto.UpdateInventoryReq) error
	Delete(ctx context.Context, req *dto.DeleteInventoryReq) error
	GetInventories(ctx context.Context, req *dto.GetInventoriesReq) (*dto.GetInventoriesResp, error)
	LoadInventories(ctx context.Context) error
}

func NewInventory(repo repo.Inventory, rdsRepo repo.InventoryRds) Inventory {
	return &inventory{
		repo:    repo,
		rdsRepo: rdsRepo,
	}
}

type inventory struct {
	repo    repo.Inventory
	rdsRepo repo.InventoryRds
}

func (obj *inventory) Create(ctx context.Context, req *dto.CreateInventoryReq) error {
	newInventory := &models.Inventory{
		Name:        req.Name,
		Description: req.Description,
		Picture:     req.Picture,
		Price:       req.Price,
		Count:       req.Count,
	}

	// 启用事务进行更新
	txErr := obj.repo.Transaction(ctx, func(tx *gorm.DB) error {
		if err := obj.repo.WithTx(tx).Create(ctx, newInventory); err != nil {
			return err
		}
		return nil
	})
	if txErr != nil {
		return txErr
	}
	return nil
}

func (obj *inventory) Update(ctx context.Context, req *dto.UpdateInventoryReq) error {
	baseModel := models.BaseModel{
		ID: req.ID,
	}
	newInventory := &models.Inventory{
		BaseModel:   baseModel,
		Name:        req.Name,
		Description: req.Description,
		Picture:     req.Picture,
		Price:       req.Price,
		Count:       req.Count,
	}

	// 启用事务进行更新
	txErr := obj.repo.Transaction(ctx, func(tx *gorm.DB) error {
		if err := obj.repo.WithTx(tx).Update(ctx, newInventory); err != nil {
			return err
		}
		return nil
	})
	if txErr != nil {
		return txErr
	}
	return nil
}

func (obj *inventory) Delete(ctx context.Context, req *dto.DeleteInventoryReq) error {
	txErr := obj.repo.Transaction(ctx, func(tx *gorm.DB) error {
		if err := obj.repo.WithTx(tx).DeleteByIds(ctx, []uint64{req.ID}); err != nil {
			return err
		}
		return nil
	})
	if txErr != nil {
		return txErr
	}
	return nil
}

func (obj *inventory) GetInventories(ctx context.Context, req *dto.GetInventoriesReq) (*dto.GetInventoriesResp, error) {
	list, total, err := obj.repo.FindInventoryList(ctx,
		repo.WithName(req.Name),
		repo.WithDesc(req.Description),
		repo.WithPrice(req.Price),
		repo.WithCount(req.Count),
		repo.WithPagination(req.Page, req.PageSize),
	)
	if err != nil {
		return nil, err
	}

	inventoryList := assemble.Model2Info(list)

	page, pageSize := gormutil.RepairPaging(req.Page, req.PageSize)
	resp := &dto.GetInventoriesResp{
		InventoryList: inventoryList,
		Pagination: dto.Pagination{
			Page:     page,
			PageSize: pageSize,
			Total:    int(total),
		},
	}

	return resp, nil
}

// LoadInventories 数据预热, 将物品数据从mysql中载入redis中, 为之后的高并发需求做准备
func (obj *inventory) LoadInventories(ctx context.Context) error {
	maxWorks := InventoryLoadWaxworks
	wg := sync.WaitGroup{}
	jobs := make(chan []*models.Inventory, maxWorks)
	errChan := make(chan error, maxWorks)

	// 启用协程, 将inventory数据载入到redis中
	for i := 0; i < maxWorks; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for inventories := range jobs {
				if err := obj.rdsRepo.Save(ctx, inventories); err != nil {
					errChan <- err
					return
				}
			}
		}()
	}

	// 从mysql读取数据, 发送到jobs通道
	offset := 1
	for {
		inventories, _, err := obj.repo.FindInventoryList(ctx)
		if err != nil {
			close(jobs)
			return err
		}

		if len(inventories) == 0 {
			break
		}

		jobs <- inventories
		offset += gormutil.DefaultPageSize
	}
	close(jobs)

	wg.Wait()

	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}
