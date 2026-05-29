package repository

import (
	"time"

	"github.com/miaomiaopu/ipmp-server/internal/model"
	"gorm.io/gorm"
)

type WorkLogRepository struct{ db *gorm.DB }

func NewWorkLogRepository(db *gorm.DB) *WorkLogRepository { return &WorkLogRepository{db: db} }

func (r *WorkLogRepository) FindByID(id string) (*model.WorkLog, error) {
	var w model.WorkLog
	err := r.db.Preload("Task").Preload("Project").Preload("Customer").Where("id = ?", id).First(&w).Error
	return &w, err
}

func (r *WorkLogRepository) List(page, pageSize int, userID, projectID *string, startDate, endDate string) ([]model.WorkLog, int64, error) {
	var items []model.WorkLog
	var total int64
	q := r.db.Model(&model.WorkLog{})
	if userID != nil && *userID != "" {
		q = q.Where("user_id = ?", *userID)
	}
	if projectID != nil && *projectID != "" {
		q = q.Where("project_id = ?", *projectID)
	}
	if startDate != "" {
		t, _ := time.Parse("2006-01-02", startDate)
		q = q.Where("log_date >= ?", t)
	}
	if endDate != "" {
		t, _ := time.Parse("2006-01-02", endDate)
		q = q.Where("log_date <= ?", t)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err := q.Preload("Task").Preload("Project").Preload("Customer").
		Order("log_date DESC").Offset(offset).Limit(pageSize).Find(&items).Error
	return items, total, err
}

// Stats 工时聚合统计
func (r *WorkLogRepository) Stats(userID *string, startDate, endDate, groupBy string) ([]model.WorkLog, error) {
	var items []model.WorkLog
	q := r.db.Preload("Project").Preload("Customer")
	if userID != nil && *userID != "" {
		q = q.Where("user_id = ?", *userID)
	}
	if startDate != "" {
		t, _ := time.Parse("2006-01-02", startDate)
		q = q.Where("log_date >= ?", t)
	}
	if endDate != "" {
		t, _ := time.Parse("2006-01-02", endDate)
		q = q.Where("log_date <= ?", t)
	}
	err := q.Order("log_date DESC").Find(&items).Error
	return items, err
}

func (r *WorkLogRepository) Create(w *model.WorkLog) error { return r.db.Create(w).Error }
func (r *WorkLogRepository) Update(w *model.WorkLog) error { return r.db.Save(w).Error }
func (r *WorkLogRepository) SoftDelete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.WorkLog{}).Error
}
