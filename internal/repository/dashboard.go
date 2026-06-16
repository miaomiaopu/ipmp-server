package repository

import (
	"time"

	"github.com/miaomiaopu/ipmp-server/internal/model"
	"gorm.io/gorm"
)

type DashboardRepository struct{ db *gorm.DB }

func NewDashboardRepository(db *gorm.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}

func (r *DashboardRepository) CountUsers(status string) (int64, error) {
	var count int64
	q := r.db.Model(&model.User{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	return count, q.Count(&count).Error
}

func (r *DashboardRepository) CountCustomers(status string) (int64, error) {
	var count int64
	q := r.db.Model(&model.Customer{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	return count, q.Count(&count).Error
}

func (r *DashboardRepository) CountProjects(status string) (int64, error) {
	var count int64
	q := r.db.Model(&model.Project{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	return count, q.Count(&count).Error
}

func (r *DashboardRepository) CountTasks(statusNot string) (int64, error) {
	var count int64
	q := r.db.Model(&model.Task{})
	if statusNot != "" {
		q = q.Where("status <> ?", statusNot)
	}
	return count, q.Count(&count).Error
}

func (r *DashboardRepository) CountRequirements(statusNot string) (int64, error) {
	var count int64
	q := r.db.Model(&model.Requirement{})
	if statusNot != "" {
		q = q.Where("status <> ?", statusNot)
	}
	return count, q.Count(&count).Error
}

func (r *DashboardRepository) SumWorkHours(userID *string, start, end time.Time) (float64, error) {
	var total float64
	q := r.db.Model(&model.WorkLog{}).Where("log_date BETWEEN ? AND ?", start, end)
	if userID != nil && *userID != "" {
		q = q.Where("user_id = ?", *userID)
	}
	err := q.Select("COALESCE(SUM(hours), 0)").Scan(&total).Error
	return total, err
}
