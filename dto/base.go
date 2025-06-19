package dto

type Pagination struct {
	Page     int `form:"page" json:"page"`           // 页码
	PageSize int `form:"page_size" json:"page_size"` // 每页大小
	Total    int `json:"count"`                      // 总数
}
