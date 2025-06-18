package services

import (
	"context"
	"gift/assemble"
	"gift/dto"
	"gift/repo"
	"gift/repo/models"
	"gorm.io/gorm"
)

type Inventory interface {
	Create(ctx context.Context, req *dto.CreateInventoryReq) error
	Update(ctx context.Context, req *dto.UpdateInventoryReq) error
	Delete(ctx context.Context, req *dto.DeleteInventoryReq) error
	GetInventories(ctx context.Context, req *dto.GetInventoriesReq) (*dto.GetInventoriesResp, error)
}

func NewInventory(repo repo.Inventory) Inventory {
	return &inventory{
		repo: repo,
	}
}

type inventory struct {
	repo repo.Inventory
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
	newInventory := &models.Inventory{
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
		repo.WithPagination(int(req.Page), int(req.PageSize)),
	)
	if err != nil {
		return nil, err
	}

	inventoryList := assemble.Model2Info(list)

	resp := &dto.GetInventoriesResp{
		InventoryList: inventoryList,
		Pagination: dto.Pagination{
			Page:     req.Page,
			PageSize: req.PageSize,
			Total:    int(total),
		},
	}

	return resp, nil
}
