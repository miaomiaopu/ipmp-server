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
	err := preloadWorkLogRefs(r.db).Where("id = ?", id).First(&w).Error
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
	err := preloadWorkLogRefs(q).
		Order("log_date DESC").Offset(offset).Limit(pageSize).Find(&items).Error
	return items, total, err
}

// Stats 工时聚合统计
func (r *WorkLogRepository) Stats(userID *string, startDate, endDate, groupBy string) ([]model.WorkLog, error) {
	var items []model.WorkLog
	q := preloadWorkLogRefs(r.db)
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

func (r *WorkLogRepository) FindForReport(userID, projectID *string, startDate, endDate string) ([]model.WorkLog, error) {
	var items []model.WorkLog
	q := preloadWorkLogRefs(r.db)
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
	err := q.Order("log_date ASC").Find(&items).Error
	return items, err
}

func (r *WorkLogRepository) FindTaskForWorkLog(taskID string) (*model.Task, error) {
	var task model.Task
	err := r.db.Preload("Project", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).
		Preload("Customer", func(db *gorm.DB) *gorm.DB { return db.Unscoped() }).
		Where("id = ?", taskID).First(&task).Error
	return &task, err
}

func (r *WorkLogRepository) Create(w *model.WorkLog) error { return r.db.Create(w).Error }
func (r *WorkLogRepository) Update(w *model.WorkLog) error { return r.db.Save(w).Error }
func (r *WorkLogRepository) SoftDelete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.WorkLog{}).Error
}

func (r *WorkLogRepository) ForceDelete(id string) error {
	return ForceDelete(r.db, &model.WorkLog{}, id)
}
func (r *WorkLogRepository) Restore(id string) error { return Restore(r.db, &model.WorkLog{}, id) }

func preloadWorkLogRefs(db *gorm.DB) *gorm.DB {
	return db.Preload("Task", func(tx *gorm.DB) *gorm.DB { return tx.Unscoped() }).
		Preload("Project", func(tx *gorm.DB) *gorm.DB { return tx.Unscoped() }).
		Preload("Customer", func(tx *gorm.DB) *gorm.DB { return tx.Unscoped() })
}
