package dto

type CreateInventoryReq struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Picture     string `json:"picture" binding:"required"`
	Price       uint64 `json:"price" binding:"required"`
	Count       uint64 `json:"count" binding:"required"`
}

type UpdateInventoryReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Picture     string `json:"picture"`
	Price       uint64 `json:"price"`
	Count       uint64 `json:"count"`
}

type DeleteInventoryReq struct {
	ID uint64 `json:"id" binding:"required"`
}
