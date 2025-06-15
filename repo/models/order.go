package models

const TableNameOrder = "t_orders"

var OrderColumns = orderColumn{
	ID:        "id",
	UserId:    "user_id",
	GiftId:    "gift_id",
	Count:     "count",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	DeletedAt: "deleted_at",
}

type Order struct {
	BaseModel
	UserId uint64 `gorm:"column:user_id;type:int unsigned;not null;comment:用户ID" json:"user_id"`
	GiftId uint64 `gorm:"column:gift_id;type:int unsigned;not null;comment:物品ID" json:"gift_id"`
	Count  uint64 `gorm:"column:count;type:int unsigned;not null;comment:下单数量" json:"count"`
}

func (*Order) TableName() string { return TableNameOrder }

type orderColumn struct {
	ID        string
	UserId    string
	GiftId    string
	Count     string
	CreatedAt string
	UpdatedAt string
	DeletedAt string
}
