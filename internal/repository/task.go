package repository

import (
	"github.com/miaomiaopu/ipmp-server/internal/model"
	"gorm.io/gorm"
)

type TaskRepository struct{ db *gorm.DB }

func NewTaskRepository(db *gorm.DB) *TaskRepository { return &TaskRepository{db: db} }

func (r *TaskRepository) FindByID(id string) (*model.Task, error) {
	var t model.Task
	err := r.db.Preload("Assignee").Preload("Project").Preload("Customer").Where("id = ?", id).First(&t).Error
	return &t, err
}

func (r *TaskRepository) List(page, pageSize int, taskType string, projectID, customerID, assigneeID *string, status, keyword string) ([]model.Task, int64, error) {
	var tasks []model.Task
	var total int64
	q := r.db.Model(&model.Task{})
	if taskType != "" {
		q = q.Where("task_type = ?", taskType)
	}
	if projectID != nil && *projectID != "" {
		q = q.Where("project_id = ?", *projectID)
	}
	if customerID != nil && *customerID != "" {
		q = q.Where("customer_id = ?", *customerID)
	}
	if assigneeID != nil && *assigneeID != "" {
		q = q.Where("assignee_id = ?", *assigneeID)
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
	err := q.Preload("Assignee").Preload("Project").Preload("Customer").
		Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&tasks).Error
	return tasks, total, err
}

func (r *TaskRepository) Create(t *model.Task) error  { return r.db.Create(t).Error }
func (r *TaskRepository) Update(t *model.Task) error  { return r.db.Save(t).Error }
func (r *TaskRepository) SoftDelete(id string) error  { return r.db.Where("id = ?", id).Delete(&model.Task{}).Error }

func (r *TaskRepository) ForceDelete(id string) error { return ForceDelete(r.db, &model.Task{}, id) }
func (r *TaskRepository) Restore(id string) error { return Restore(r.db, &model.Task{}, id) }
