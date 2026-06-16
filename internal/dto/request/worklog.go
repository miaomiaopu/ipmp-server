package request

type CreateWorkLogRequest struct {
	TaskID      *string `json:"task_id"`
	LogDate     string  `json:"log_date" binding:"required"`
	Hours       float64 `json:"hours" binding:"required,gte=0.5,lte=24"`
	Description string  `json:"description"`
}

type UpdateWorkLogRequest struct {
	LogDate     *string  `json:"log_date"`
	Hours       *float64 `json:"hours" binding:"omitempty,gte=0.5,lte=24"`
	Description *string  `json:"description"`
}

type WorkLogQuery struct {
	UserID    *string `form:"user_id"`
	ProjectID *string `form:"project_id"`
	StartDate string  `form:"start_date"`
	EndDate   string  `form:"end_date"`
	Page      int     `form:"page"`
	PageSize  int     `form:"page_size"`
}

func (q *WorkLogQuery) Normalize() (int, int) {
	p, ps := q.Page, q.PageSize
	if p <= 0 {
		p = 1
	}
	if ps <= 0 || ps > 100 {
		ps = 20
	}
	return p, ps
}
