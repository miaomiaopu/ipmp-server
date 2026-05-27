package model

import "time"

// WorkLog 工时记录
// 记录用户每日在各任务上的工时消耗
//
// project_id / customer_id 为冗余字段：
//   录入时从关联 Task 推导填充，避免报表聚合时多表 JOIN
//   写入时同步维护，换取工时统计报表的读性能
type WorkLog struct {
	BaseModel
	TaskID      *string   `gorm:"size:36" json:"task_id"`
	UserID      string    `gorm:"size:36;index" json:"user_id"`
	ProjectID   *string   `gorm:"size:36" json:"project_id"`
	CustomerID  *string   `gorm:"size:36" json:"customer_id"`
	LogDate     time.Time `json:"log_date"`
	Hours       float64   `gorm:"type:decimal(5,2)" json:"hours"`
	Description string    `gorm:"type:text" json:"description"`

	// 关联
	Task     *Task     `gorm:"foreignKey:TaskID" json:"task,omitempty"`
	User     *User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Project  *Project  `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	Customer *Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
}

func (WorkLog) TableName() string { return "work_logs" }
