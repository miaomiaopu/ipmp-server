package model

import "time"

// Project 项目
// 关联一个客户(customers)和一个项目经理(users)
// Status: planning(立项) → in_progress(实施) → completed(竣工) / suspended(挂起)
type Project struct {
	BaseModel
	ProjectCode    string     `gorm:"uniqueIndex;size:32" json:"project_code"`
	Name           string     `gorm:"size:256" json:"name"`
	CustomerID     *string    `gorm:"size:36" json:"customer_id"`
	ManagerID      *string    `gorm:"size:36" json:"manager_id"`
	StartDate      *time.Time `json:"start_date"`
	GoLiveDate     *time.Time `json:"go_live_date"`
	CompletionDate *time.Time `json:"completion_date"`
	Status         string     `gorm:"size:32;default:planning" json:"status"`
	Description    string     `gorm:"type:text" json:"description"`

	// CustomerID → customers 表
	Customer *Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	// ManagerID → users 表
	Manager *User `gorm:"foreignKey:ManagerID" json:"manager,omitempty"`
}

func (Project) TableName() string { return "projects" }

const (
	ProjectStatusPlanning   = "planning"
	ProjectStatusInProgress = "in_progress"
	ProjectStatusCompleted  = "completed"
	ProjectStatusSuspended  = "suspended"
)
