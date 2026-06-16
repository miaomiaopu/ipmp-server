package request

// CreateProjectRequest 创建项目请求
type CreateProjectRequest struct {
	ProjectCode    string  `json:"project_code" binding:"required,min=1,max=32"`
	Name           string  `json:"name" binding:"required,min=1,max=256"`
	CustomerID     *string `json:"customer_id" binding:"required,uuid"`
	ManagerID      *string `json:"manager_id" binding:"omitempty,uuid"`
	StartDate      *string `json:"start_date" binding:"omitempty"`
	GoLiveDate     *string `json:"go_live_date" binding:"omitempty"`
	CompletionDate *string `json:"completion_date" binding:"omitempty"`
	Status         string  `json:"status" binding:"omitempty,oneof=planning in_progress online completed"`
	Description    string  `json:"description" binding:"omitempty,max=4096"`
}

// UpdateProjectRequest 更新项目请求
type UpdateProjectRequest struct {
	Name           *string `json:"name" binding:"omitempty,min=1,max=256"`
	ProjectCode    *string `json:"project_code" binding:"omitempty,min=1,max=32"`
	CustomerID     *string `json:"customer_id" binding:"omitempty,uuid"`
	ManagerID      *string `json:"manager_id" binding:"omitempty,uuid"`
	StartDate      *string `json:"start_date" binding:"omitempty"`
	GoLiveDate     *string `json:"go_live_date" binding:"omitempty"`
	CompletionDate *string `json:"completion_date" binding:"omitempty"`
	Status         *string `json:"status" binding:"omitempty,oneof=planning in_progress online completed"`
	Description    *string `json:"description" binding:"omitempty,max=4096"`
}

// ProjectQuery 项目列表查询参数
type ProjectQuery struct {
	CustomerID *string `form:"customer_id"`
	Status     string  `form:"status"`
	Page       int     `form:"page"`
	PageSize   int     `form:"page_size"`
	Keyword    string  `form:"keyword"`
}

func (q *ProjectQuery) Normalize() (page, pageSize int) {
	page = q.Page
	if page < 1 {
		page = 1
	}
	pageSize = q.PageSize
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return
}
