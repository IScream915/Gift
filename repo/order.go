package repo

import (
	"context"
	"gift/repo/models"
	"gorm.io/gorm"
)

type Order interface {
	Create(ctx context.Context, order *models.Order) error
	Update(ctx context.Context, order *models.Order) error
	DeleteByIds(ctx context.Context, ids []uint64) error
	FindById(ctx context.Context, id uint64) (*models.Order, error)
	FindAll(ctx context.Context) ([]*models.Order, error)
	Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error
	WithTx(tx *gorm.DB) Order
}

func NewOrder(client *gorm.DB) Order {
	return &order{
		client: client,
	}
}

type order struct {
	client *gorm.DB
}

func (obj *order) Create(ctx context.Context, order *models.Order) error {
	return obj.client.WithContext(ctx).Create(order).Error
}

func (obj *order) Update(ctx context.Context, order *models.Order) error {
	return obj.client.WithContext(ctx).Model(&models.Order{}).Where("id = ?", order.ID).Updates(order).Error
}

func (obj *order) DeleteByIds(ctx context.Context, ids []uint64) error {
	return obj.client.WithContext(ctx).Delete(&models.Order{}, ids).Error
}

func (obj *order) FindById(ctx context.Context, id uint64) (*models.Order, error) {
	record := &models.Order{}
	err := obj.client.WithContext(ctx).Where("id = ?", id).First(record).Error
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (obj *order) FindAll(ctx context.Context) ([]*models.Order, error) {
	records := make([]*models.Order, 0)
	err := obj.client.WithContext(ctx).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (obj *order) Transaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	tx := obj.client.WithContext(ctx).Begin()
	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (obj *order) WithTx(tx *gorm.DB) Order {
	return &order{client: tx}
}
