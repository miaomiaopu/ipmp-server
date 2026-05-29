package service

import (
	"errors"
	"time"

	"github.com/miaomiaopu/ipmp-server/internal/dto/request"
	"github.com/miaomiaopu/ipmp-server/internal/model"
	"github.com/miaomiaopu/ipmp-server/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrTaskNotFound    = errors.New("task not found")
	ErrInvalidTaskType = errors.New("task type requires project_id (project) or customer_id (customer)")
)

type TaskService struct{ repo *repository.TaskRepository }

func NewTaskService(repo *repository.TaskRepository) *TaskService { return &TaskService{repo: repo} }

func (s *TaskService) List(page, pageSize int, query request.TaskQuery) ([]model.Task, int64, error) {
	return s.repo.List(page, pageSize, query.TaskType, query.ProjectID, query.CustomerID, query.AssigneeID, query.Status, query.Keyword)
}

func (s *TaskService) GetByID(id string) (*model.Task, error) {
	t, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	return t, nil
}

func (s *TaskService) Create(task *model.Task) error {
	switch task.TaskType {
	case model.TaskTypeProject:
		if task.ProjectID == nil || *task.ProjectID == "" {
			return ErrInvalidTaskType
		}
	case model.TaskTypeCustomer:
		if task.CustomerID == nil || *task.CustomerID == "" {
			return ErrInvalidTaskType
		}
	}
	return s.repo.Create(task)
}

func (s *TaskService) Update(id string, updates map[string]interface{}) error {
	t, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTaskNotFound
		}
		return err
	}
	for k, v := range updates {
		switch k {
		case "title":
			t.Title = v.(string)
		case "description":
			t.Description = v.(string)
		case "project_id":
			if s, ok := v.(string); ok {
				t.ProjectID = &s
			}
		case "customer_id":
			if s, ok := v.(string); ok {
				t.CustomerID = &s
			}
		case "assignee_id":
			if s, ok := v.(string); ok {
				t.AssigneeID = &s
			}
		case "status":
			t.Status = v.(string)
		case "priority":
			t.Priority = v.(string)
		case "estimated_hours":
			t.EstimatedHours = v.(float64)
		case "due_date":
			if s, ok := v.(string); ok && s != "" {
				dt, _ := time.Parse("2006-01-02", s)
				t.DueDate = &dt
			}
		}
	}
	t.UpdatedAt = time.Now()
	return s.repo.Update(t)
}

func (s *TaskService) UpdateStatus(id, status string) error {
	t, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTaskNotFound
		}
		return err
	}
	t.Status = status
	return s.repo.Update(t)
}

func (s *TaskService) Delete(id string) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTaskNotFound
		}
		return err
	}
	return s.repo.SoftDelete(id)
}
