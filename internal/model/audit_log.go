package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// AuditLog 审计日志
// 记录所有写操作和敏感读取，用于安全审计
// Detail 为 JSON 格式，存储操作的详细上下文（变更前后值、请求参数摘要等）
// UserID 可为空（未登录的失败操作也记录）
type AuditLog struct {
	ID         string          `gorm:"primaryKey;size:36" json:"id"`
	UserID     *string         `gorm:"size:36" json:"user_id"`
	Action     string          `gorm:"size:64" json:"action"`
	Resource   string          `gorm:"size:64" json:"resource"`
	ResourceID *string         `gorm:"size:64" json:"resource_id"`
	Detail     json.RawMessage `gorm:"type:json" json:"detail,omitempty"`
	IPAddress  string          `gorm:"size:64" json:"ip_address"`
	CreatedAt  time.Time       `json:"created_at"`
}

func (AuditLog) TableName() string { return "audit_logs" }

func (a *AuditLog) BeforeCreate(tx *gorm.DB) error {
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	return nil
}
