package request

type WeeklyReportQuery struct {
	Page     int `form:"page"`
	PageSize int `form:"page_size"`
}

func (q *WeeklyReportQuery) Normalize() (int, int) {
	p, ps := q.Page, q.PageSize
	if p <= 0 {
		p = 1
	}
	if ps <= 0 || ps > 100 {
		ps = 20
	}
	return p, ps
}

type GenerateWeeklyReportRequest struct {
	ReportType string  `json:"report_type" binding:"omitempty,oneof=personal project"`
	ProjectID  *string `json:"project_id"`
	WeekStart  string  `json:"week_start" binding:"required"`
	WeekEnd    string  `json:"week_end" binding:"required"`
	UseAI      bool    `json:"use_ai"`
}

type UpdateWeeklyReportRequest struct {
	Content string  `json:"content" binding:"required"`
	Status  *string `json:"status" binding:"omitempty,oneof=draft submitted reviewed"`
}
