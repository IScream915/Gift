package gormutil

import "gorm.io/gorm"

const (
	DefaultPage     = 1
	DefaultPageSize = 10
)

type Option struct {
	Page     int
	PageSize int
	Sort     string
}

// PagingAndSort 处理分页和排序
func (obj *Option) PagingAndSort(query *gorm.DB) *gorm.DB {
	// 分页
	query = query.Offset((obj.Page - 1) * obj.PageSize).Limit(obj.PageSize)

	// 排序
	if obj.Sort != "" {
		query = query.Order(obj.Sort)
	}

	return query
}

// RepairPaging 修复分页参数
func RepairPaging(page, pageSize int) (int, int) {
	if page <= 0 {
		page = DefaultPage
	}

	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}

	return page, pageSize
}
