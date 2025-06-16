package repo

import (
	"context"
	"gift/repo/models"
	"gorm.io/gorm"
)

type Inventory interface {
	Create(ctx context.Context, inventory *models.Inventory) error
	Update(ctx context.Context, inventory *models.Inventory) error
	DeleteByIds(ctx context.Context, ids []uint64) error
	FindById(ctx context.Context, id uint64) (*models.Inventory, error)
	FindAll(ctx context.Context) ([]*models.Inventory, error)
	Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error
	WithTx(tx *gorm.DB) Inventory
}

func NewInventory(client *gorm.DB) Inventory {
	return &inventory{
		client: client,
	}
}

type inventory struct {
	client *gorm.DB
}

func (obj *inventory) Create(ctx context.Context, inventory *models.Inventory) error {
	return obj.client.WithContext(ctx).Create(inventory).Error
}

func (obj *inventory) Update(ctx context.Context, inventory *models.Inventory) error {
	return obj.client.WithContext(ctx).Model(&models.Inventory{}).Where("id = ?", inventory.ID).Updates(inventory).Error
}

func (obj *inventory) DeleteByIds(ctx context.Context, ids []uint64) error {
	// GORM Delete(..., ...) 根据第二个参数是否为主键切片, 生成相应的 WHERE id = ？或 WHERE id IN (?)
	// 即: GORM的Delete函数会自己判断
	return obj.client.WithContext(ctx).Delete(&models.Inventory{}, ids).Error
}

func (obj *inventory) FindById(ctx context.Context, id uint64) (*models.Inventory, error) {
	record := &models.Inventory{}
	err := obj.client.WithContext(ctx).Where("id = ?", id).First(record).Error
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (obj *inventory) FindAll(ctx context.Context) ([]*models.Inventory, error) {
	records := make([]*models.Inventory, 0)
	err := obj.client.WithContext(ctx).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (obj *inventory) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	tx := obj.client.WithContext(ctx).Begin()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (obj *inventory) WithTx(tx *gorm.DB) Inventory {
	return &inventory{client: tx}
}
