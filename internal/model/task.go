package model

import "time"

// Task 统一任务模型，通过 TaskType 区分任务类别：
//   - project:  项目任务，ProjectID 必填，CustomerID 可从关联项目推导
//   - customer: 客户任务，CustomerID 必填，ProjectID 为空
//   - daily:    日常任务，ProjectID/CustomerID 均为空
//
// 关联字段均为可空(*string)，Service 层负责根据 TaskType 校验必填约束
type Task struct {
	BaseModel
	TaskType       string     `gorm:"size:32;default:project" json:"task_type"`
	Title          string     `gorm:"size:256" json:"title"`
	Description    string     `gorm:"type:text" json:"description"`
	ProjectID      *string    `gorm:"size:36" json:"project_id"`
	CustomerID     *string    `gorm:"size:36" json:"customer_id"`
	Status         string     `gorm:"size:32;default:in_progress" json:"status"`
	Priority       string     `gorm:"size:32;default:medium" json:"priority"`
	DueDate        *time.Time `json:"due_date"`

	// ProjectID → projects (仅 project 任务)
	Project *Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	// CustomerID → customers (仅 customer 任务)
	Customer *Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
}

func (Task) TableName() string { return "tasks" }

const (
	TaskTypeProject  = "project"
	TaskTypeCustomer = "customer"
	TaskTypeDaily    = "daily"

	TaskStatusTodo       = "todo"
	TaskStatusInProgress = "in_progress"
	TaskStatusDone       = "done"

	TaskPriorityLow    = "low"
	TaskPriorityMedium = "medium"
	TaskPriorityHigh   = "high"
)
