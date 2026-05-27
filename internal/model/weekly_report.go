package model

import "time"

// WeeklyReport 周报
// report_type: personal(个人周报) / project(项目周报)
// personal 周报: UserID 必填，ProjectID 为空
// project 周报: ProjectID 必填，汇总该项目所有成员的工时
// ai_raw_content 存储 AI 生成的原始内容，content 为编辑后的最终内容
type WeeklyReport struct {
	BaseModel
	UserID       string    `gorm:"size:36;index" json:"user_id"`
	WeekStart    time.Time `json:"week_start"`
	WeekEnd      time.Time `json:"week_end"`
	ReportType   string    `gorm:"size:32;default:personal" json:"report_type"`
	ProjectID    *string   `gorm:"size:36" json:"project_id"`
	Content      string    `gorm:"type:text" json:"content"`
	AIRawContent *string   `gorm:"type:text" json:"ai_raw_content,omitempty"`
	Status       string    `gorm:"size:32;default:draft" json:"status"`

	// UserID → users
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	// ProjectID → projects (仅 project 周报)
	Project *Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}

func (WeeklyReport) TableName() string { return "weekly_reports" }

const (
	ReportTypePersonal = "personal"
	ReportTypeProject  = "project"

	ReportStatusDraft = "draft"
	ReportStatusFinal = "final"
)
