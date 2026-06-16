package repository

import (
	"github.com/miaomiaopu/ipmp-server/internal/model"
	"gorm.io/gorm"
)

type WeeklyReportRepository struct{ db *gorm.DB }

func NewWeeklyReportRepository(db *gorm.DB) *WeeklyReportRepository {
	return &WeeklyReportRepository{db: db}
}

func (r *WeeklyReportRepository) List(page, pageSize int, userID *string) ([]model.WeeklyReport, int64, error) {
	var items []model.WeeklyReport
	var total int64
	q := r.db.Model(&model.WeeklyReport{})
	if userID != nil && *userID != "" {
		q = q.Where("user_id = ?", *userID)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err := q.Preload("User").Preload("Project").
		Order("week_start DESC, created_at DESC").
		Offset(offset).Limit(pageSize).Find(&items).Error
	return items, total, err
}

func (r *WeeklyReportRepository) FindByID(id string) (*model.WeeklyReport, error) {
	var item model.WeeklyReport
	err := r.db.Preload("User").Preload("Project").Where("id = ?", id).First(&item).Error
	return &item, err
}

func (r *WeeklyReportRepository) Create(item *model.WeeklyReport) error {
	return r.db.Create(item).Error
}

func (r *WeeklyReportRepository) Update(item *model.WeeklyReport) error {
	return r.db.Save(item).Error
}
