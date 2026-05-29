package request

type CreateTaskRequest struct {
	TaskType       string  `json:"task_type" binding:"required,oneof=project customer daily"`
	Title          string  `json:"title" binding:"required,max=256"`
	Description    string  `json:"description"`
	ProjectID      *string `json:"project_id"`
	CustomerID     *string `json:"customer_id"`
	AssigneeID     *string `json:"assignee_id"`
	Priority       string  `json:"priority" binding:"omitempty,oneof=low medium high urgent"`
	DueDate        *string `json:"due_date"`
	EstimatedHours float64 `json:"estimated_hours"`
}

type UpdateTaskRequest struct {
	Title          *string  `json:"title"`
	Description    *string  `json:"description"`
	ProjectID      *string  `json:"project_id"`
	CustomerID     *string  `json:"customer_id"`
	AssigneeID     *string  `json:"assignee_id"`
	Status         *string  `json:"status" binding:"omitempty,oneof=todo in_progress done"`
	Priority       *string  `json:"priority" binding:"omitempty,oneof=low medium high urgent"`
	DueDate        *string  `json:"due_date"`
	EstimatedHours *float64 `json:"estimated_hours"`
}

type UpdateTaskStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=todo in_progress done"`
}

type TaskQuery struct {
	TaskType   string  `form:"task_type"`
	ProjectID  *string `form:"project_id"`
	CustomerID *string `form:"customer_id"`
	AssigneeID *string `form:"assignee_id"`
	Status     string  `form:"status"`
	Keyword    string  `form:"keyword"`
	Page       int     `form:"page"`
	PageSize   int     `form:"page_size"`
}

func (q *TaskQuery) Normalize() (page, pageSize int) {
	page, pageSize = q.Page, q.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	return
}
