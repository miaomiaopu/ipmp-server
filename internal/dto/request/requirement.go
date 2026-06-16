package request

type CreateRequirementRequest struct {
	ReqType         string  `json:"req_type" binding:"required,oneof=project after_sales"`
	Title           string  `json:"title" binding:"required,max=256"`
	Description     string  `json:"description"`
	ProjectID       *string `json:"project_id"`
	CustomerID      *string `json:"customer_id"`
	RequirementCode string  `json:"requirement_code" binding:"omitempty,max=32"`
	Priority        string  `json:"priority" binding:"omitempty,oneof=low medium high"`
	ScheduledDate   *string `json:"scheduled_date"`
}

type UpdateRequirementRequest struct {
	Title           *string `json:"title"`
	Description     *string `json:"description"`
	ProjectID       *string `json:"project_id"`
	CustomerID      *string `json:"customer_id"`
	RequirementCode *string `json:"requirement_code" binding:"omitempty,max=32"`
	Priority        *string `json:"priority" binding:"omitempty,oneof=low medium high"`
	Status          *string `json:"status" binding:"omitempty,oneof=pending testing closed tested online"`
	ScheduledDate   *string `json:"scheduled_date"`
}

type RequirementQuery struct {
	ReqType       string  `form:"req_type"`
	ProjectID     *string `form:"project_id"`
	CustomerID    *string `form:"customer_id"`
	ScheduledDate string  `form:"scheduled_date"`
	Status        string  `form:"status"`
	Keyword       string  `form:"keyword"`
	Page          int     `form:"page"`
	PageSize      int     `form:"page_size"`
}

func (q *RequirementQuery) Normalize() (int, int) {
	p, ps := q.Page, q.PageSize
	if p <= 0 {
		p = 1
	}
	if ps <= 0 || ps > 100 {
		ps = 20
	}
	return p, ps
}
