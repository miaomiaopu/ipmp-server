package request

// Pagination 分页请求参数
type Pagination struct {
	Page     int    `form:"page" json:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" json:"page_size" binding:"omitempty,min=1,max=100"`
	Keyword  string `form:"keyword" json:"keyword"`
	Status   string `form:"status" json:"status"`
}

func (p *Pagination) Normalize() (page, pageSize int) {
	page = p.Page
	if page < 1 {
		page = 1
	}
	pageSize = p.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return
}
