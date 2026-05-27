package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserAIConfig 用户独立 AI 配置
// 每个用户可配置自己的 AI Provider + API Key
// APIKey 使用 EncryptedField 加密存储，JSON 序列化时完全隐藏(json:"-")
// 服务端调用 AI 时解密用户 Key，Key 不出服务器
type UserAIConfig struct {
	ID        string         `gorm:"primaryKey;size:36" json:"id"`
	UserID    string         `gorm:"uniqueIndex;size:36" json:"user_id"`
	Provider  string         `gorm:"size:32;default:deepseek" json:"provider"`
	APIKey    EncryptedField `gorm:"type:text" json:"-"`
	Model     string         `gorm:"size:128" json:"model"`
	BaseURL   string         `gorm:"size:512" json:"base_url"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (UserAIConfig) TableName() string { return "user_ai_configs" }

func (c *UserAIConfig) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}
