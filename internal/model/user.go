package model

// User 系统用户
// 角色分 admin(管理员) / user(普通用户)
// Email/Phone 使用 EncryptedField 加密存储，JSON 自动掩码
// PasswordHash 为 bcrypt 哈希，JSON 序列化时隐藏
type User struct {
	BaseModel
	Username     string         `gorm:"uniqueIndex;size:64" json:"username"`
	PasswordHash string         `gorm:"size:256" json:"-"`
	DisplayName  string         `gorm:"size:128" json:"display_name"`
	Email        EncryptedField `gorm:"type:text" json:"email"`
	Phone        EncryptedField `gorm:"type:text" json:"phone"`
	Role         string         `gorm:"size:32;default:user" json:"role"`
	Status       string         `gorm:"size:32;default:active" json:"status"`
}

func (User) TableName() string { return "users" }

const (
	RoleAdmin = "admin"
	RoleUser  = "user"

	UserStatusActive   = "active"
	UserStatusDisabled = "disabled"
)
