package repository

import (
	"github.com/miaomiaopu/ipmp-server/internal/model"
	"gorm.io/gorm"
)

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) FindByID(id string) (*model.Project, error) {
	var project model.Project
	err := r.db.Preload("Customer").Preload("Manager").Where("id = ?", id).First(&project).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) FindByCode(code string) (*model.Project, error) {
	var project model.Project
	err := r.db.Where("project_code = ?", code).First(&project).Error
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) List(page, pageSize int, keyword, status string, customerID *string) ([]model.Project, int64, error) {
	var projects []model.Project
	var total int64
	query := r.db.Model(&model.Project{})
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("name LIKE ? OR project_code LIKE ?", like, like)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if customerID != nil && *customerID != "" {
		query = query.Where("customer_id = ?", *customerID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err := query.Preload("Customer").Preload("Manager").
		Order(`CASE status WHEN 'planning' THEN go_live_date WHEN 'online' THEN completion_date ELSE '9999-12-31' END ASC`).
		Offset(offset).Limit(pageSize).Find(&projects).Error
	return projects, total, err
}

func (r *ProjectRepository) Create(project *model.Project) error {
	return r.db.Create(project).Error
}

func (r *ProjectRepository) Update(project *model.Project) error {
	return r.db.Save(project).Error
}

func (r *ProjectRepository) SoftDelete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.Project{}).Error
}

// CountTasks 统计项目关联的任务数
func (r *ProjectRepository) CountTasks(projectID string) (int64, error) {
	var count int64
	err := r.db.Model(&model.Task{}).Where("project_id = ?", projectID).Count(&count).Error
	return count, err
}

// CountRequirements 统计项目关联的需求数
func (r *ProjectRepository) CountRequirements(projectID string) (int64, error) {
	var count int64
	err := r.db.Model(&model.Requirement{}).Where("project_id = ?", projectID).Count(&count).Error
	return count, err
}

func (r *ProjectRepository) ForceDelete(id string) error { return ForceDelete(r.db, &model.Project{}, id) }
func (r *ProjectRepository) Restore(id string) error { return Restore(r.db, &model.Project{}, id) }
