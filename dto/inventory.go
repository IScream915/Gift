package dto

type InventoryInfo struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Picture     string `json:"picture"`
	Price       uint64 `json:"price"`
	Count       uint64 `json:"count"`
	CreatedAt   uint64 `json:"created_at"`
	UpdatedAt   uint64 `json:"updated_at"`
}
type CreateInventoryReq struct {
	Name        string `json:"name" binding:"required"`    // 物品名称
	Description string `json:"description"`                // 物品描述
	Picture     string `json:"picture" binding:"required"` // 物品图片路径
	Price       uint64 `json:"price" binding:"required"`   // 物品价格
	Count       uint64 `json:"count" binding:"required"`   // 物品库存
}

type UpdateInventoryReq struct {
	ID          uint64 `json:"id" binding:"required"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Picture     string `json:"picture"`
	Price       uint64 `json:"price"`
	Count       uint64 `json:"count"`
}

type DeleteInventoryReq struct {
	ID uint64 `json:"id" binding:"required"`
}

type GetInventoriesReq struct {
	Name        string `form:"name"`        // 物品名称
	Description string `form:"description"` // 物品描述
	Price       uint64 `form:"price"`       // 物品价格
	Count       uint64 `form:"count"`       // 物品库存
	Pagination         // 分页需求
}

type GetInventoriesResp struct {
	InventoryList []*InventoryInfo `json:"inventory_list"`
	Pagination    Pagination       `json:"pagination"`
}

type SecKillReq struct {
	ID     uint64 `json:"id" binding:"required"`      // 物品ID
	UserID uint64 `json:"user_id" binding:"required"` // 用户ID
}
