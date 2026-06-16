package model

// Customer 客户信息
// ContactPhone/ContactEmail/Address 使用 EncryptedField 加密存储
// Status: active(正常) / inactive(停用)
type Customer struct {
	BaseModel
	CustomerCode  string         `gorm:"size:32;index" json:"customer_code"`
	Name          string         `gorm:"size:256" json:"name"`
	ContactPerson string         `gorm:"size:128" json:"contact_person"`
	ContactPhone  EncryptedField `gorm:"type:text" json:"contact_phone"`
	ContactEmail  EncryptedField `gorm:"type:text" json:"contact_email"`
	Address       EncryptedField `gorm:"type:text" json:"address"`
	Notes         string         `gorm:"type:text" json:"notes"`
	Status        string         `gorm:"size:32;default:active" json:"status"`
}

func (Customer) TableName() string { return "customers" }
