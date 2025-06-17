package models

const TableNameInventories = "t_inventories"

var InventoryColumn = inventoryColumn{
	ID:          "id",
	Name:        "name",
	Description: "description",
	Picture:     "picture",
	Price:       "price",
	Count:       "count",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}

type Inventory struct {
	BaseModel
	Name        string `gorm:"column:name;type:varchar(50);uniqueIndex:uniq_name;not null;comment:物品名称" json:"name"`                 // 物品名称
	Description string `gorm:"column:description;type:varchar(255);default:'default desc';not null;comment:物品描述" json:"description"` // 物品描述
	Picture     string `gorm:"column:picture;type:varchar(255);not null;comment:物品图片保存路径" json:"picture"`                            // 物品图片保存路径
	Price       uint64 `gorm:"column:price;type:int unsigned;not null;comment:物品价格" json:"price"`                                    // 物品价格
	Count       uint64 `gorm:"column:count;type:int unsigned;not null;comment:库存数量" json:"count"`                                    // 库存数量
}

func (*Inventory) TableName() string { return TableNameInventories }

type inventoryColumn struct {
	ID          string
	Name        string
	Description string
	Picture     string
	Price       string
	Count       string
	CreatedAt   string
	UpdatedAt   string
	DeletedAt   string
}
