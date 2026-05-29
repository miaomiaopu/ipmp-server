package request

type CreateRequirementRequest struct {
	ReqType     string  `json:"req_type" binding:"required,oneof=project after_sales"`
	Title       string  `json:"title" binding:"required,max=256"`
	Description string  `json:"description"`
	ProjectID   *string `json:"project_id"`
	CustomerID  *string `json:"customer_id"`
	Priority    string  `json:"priority" binding:"omitempty,oneof=low medium high urgent"`
	Submitter   string  `json:"submitter"`
}

type UpdateRequirementRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Priority    *string `json:"priority" binding:"omitempty,oneof=low medium high urgent"`
	Status      *string `json:"status" binding:"omitempty,oneof=pending approved rejected in_progress done"`
	Submitter   *string `json:"submitter"`
}

type RequirementQuery struct {
	ReqType    string  `form:"req_type"`
	ProjectID  *string `form:"project_id"`
	CustomerID *string `form:"customer_id"`
	Status     string  `form:"status"`
	Keyword    string  `form:"keyword"`
	Page       int     `form:"page"`
	PageSize   int     `form:"page_size"`
}

func (q *RequirementQuery) Normalize() (int, int) {
	p, ps := q.Page, q.PageSize
	if p <= 0 { p = 1 }
	if ps <= 0 || ps > 100 { ps = 20 }
	return p, ps
}
