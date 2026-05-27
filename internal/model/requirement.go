package model

// Requirement 统一需求模型，通过 ReqType 区分需求类别：
//   - project:     项目需求，ProjectID 必填
//   - after_sales: 售后需求，CustomerID 必填
//
// Service 层负责根据 ReqType 校验必填约束
type Requirement struct {
	BaseModel
	ReqType     string  `gorm:"size:32;default:project" json:"req_type"`
	Title       string  `gorm:"size:256" json:"title"`
	Description string  `gorm:"type:text" json:"description"`
	ProjectID   *string `gorm:"size:36" json:"project_id"`
	CustomerID  *string `gorm:"size:36" json:"customer_id"`
	Priority    string  `gorm:"size:32;default:medium" json:"priority"`
	Status      string  `gorm:"size:32;default:pending" json:"status"`
	Submitter   string  `gorm:"size:128" json:"submitter"`

	// ProjectID → projects (仅 project 需求)
	Project *Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	// CustomerID → customers (仅 after_sales 需求)
	Customer *Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
}

func (Requirement) TableName() string { return "requirements" }

const (
	ReqTypeProject    = "project"
	ReqTypeAfterSales = "after_sales"

	ReqStatusPending  = "pending"
	ReqStatusApproved = "approved"
	ReqStatusRejected = "rejected"
)
