package models

const TableNameOrder = "t_orders"

var OrderColumns = orderColumn{
	ID:          "id",
	UserId:      "user_id",
	InventoryId: "inventory_id",
	Count:       "count",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}

type Order struct {
	BaseModel
	UserId      uint64 `gorm:"column:user_id;type:int unsigned;not null;comment:用户ID" json:"user_id"`
	InventoryId uint64 `gorm:"column:inventory_id;type:int unsigned;not null;comment:物品ID" json:"inventory_id"`
	Count       uint64 `gorm:"column:count;type:int unsigned;not null;comment:下单数量" json:"count"`
}

func (*Order) TableName() string { return TableNameOrder }

type orderColumn struct {
	ID          string
	UserId      string
	InventoryId string
	Count       string
	CreatedAt   string
	UpdatedAt   string
	DeletedAt   string
}
