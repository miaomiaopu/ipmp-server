package repository

import (
	"github.com/miaomiaopu/ipmp-server/internal/model"
	"gorm.io/gorm"
)

type RequirementRepository struct{ db *gorm.DB }

func NewRequirementRepository(db *gorm.DB) *RequirementRepository { return &RequirementRepository{db: db} }

func (r *RequirementRepository) FindByID(id string) (*model.Requirement, error) {
	var m model.Requirement
	err := r.db.Preload("Project").Preload("Customer").Where("id = ?", id).First(&m).Error
	return &m, err
}

func (r *RequirementRepository) List(page, pageSize int, reqType string, projectID, customerID *string, status, keyword string) ([]model.Requirement, int64, error) {
	var items []model.Requirement
	var total int64
	q := r.db.Model(&model.Requirement{})
	if reqType != "" {
		q = q.Where("req_type = ?", reqType)
	}
	if projectID != nil && *projectID != "" {
		q = q.Where("project_id = ?", *projectID)
	}
	if customerID != nil && *customerID != "" {
		q = q.Where("customer_id = ?", *customerID)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if keyword != "" {
		q = q.Where("title LIKE ?", "%"+keyword+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err := q.Preload("Project").Preload("Customer").
		Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&items).Error
	return items, total, err
}

func (r *RequirementRepository) Create(m *model.Requirement) error { return r.db.Create(m).Error }
func (r *RequirementRepository) Update(m *model.Requirement) error { return r.db.Save(m).Error }
func (r *RequirementRepository) SoftDelete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.Requirement{}).Error
}
