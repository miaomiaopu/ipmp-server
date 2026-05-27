package service

import (
	"errors"
	"time"

	"github.com/miaomiaopu/ipmp-server/internal/model"
	"github.com/miaomiaopu/ipmp-server/internal/repository"
	"gorm.io/gorm"
)

var ErrProjectNotFound = errors.New("project not found")

type ProjectService struct {
	repo *repository.ProjectRepository
}

func NewProjectService(repo *repository.ProjectRepository) *ProjectService {
	return &ProjectService{repo: repo}
}

func (s *ProjectService) List(page, pageSize int, keyword, status string, customerID *string) ([]model.Project, int64, error) {
	return s.repo.List(page, pageSize, keyword, status, customerID)
}

// ProjectDetail 项目详情含统计
type ProjectDetail struct {
	model.Project
	TaskCount        int64 `json:"task_count"`
	RequirementCount int64 `json:"requirement_count"`
}

func (s *ProjectService) GetByID(id string) (*ProjectDetail, error) {
	project, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}
	taskCount, _ := s.repo.CountTasks(id)
	reqCount, _ := s.repo.CountRequirements(id)
	return &ProjectDetail{
		Project:          *project,
		TaskCount:        taskCount,
		RequirementCount: reqCount,
	}, nil
}

func (s *ProjectService) Create(req *model.Project) error {
	existing, err := s.repo.FindByCode(req.ProjectCode)
	if err == nil && existing != nil {
		return errors.New("project code already exists")
	}
	return s.repo.Create(req)
}

func (s *ProjectService) Update(id string, updates map[string]interface{}) error {
	project, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProjectNotFound
		}
		return err
	}
	for key, value := range updates {
		switch key {
		case "name":
			project.Name = value.(string)
		case "project_code":
			project.ProjectCode = value.(string)
		case "customer_id":
			v := value.(string)
			project.CustomerID = &v
		case "manager_id":
			v := value.(string)
			project.ManagerID = &v
		case "status":
			project.Status = value.(string)
		case "description":
			project.Description = value.(string)
		}
	}
	project.UpdatedAt = time.Now()
	return s.repo.Update(project)
}

func (s *ProjectService) Delete(id string) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProjectNotFound
		}
		return err
	}
	return s.repo.SoftDelete(id)
}
