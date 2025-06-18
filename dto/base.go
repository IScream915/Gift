package dto

type Pagination struct {
	Page     uint64 `form:"page" json:"page" binding:"required"`           // 页码
	PageSize uint64 `form:"page_size" json:"page_size" binding:"required"` // 每页大小
	Total    int    `json:"count"`                                         // 总数
}
