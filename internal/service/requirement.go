package service

import (
	"errors"
	"time"

	"github.com/miaomiaopu/ipmp-server/internal/model"
	"github.com/miaomiaopu/ipmp-server/internal/repository"
	"gorm.io/gorm"
)

var ErrRequirementNotFound = errors.New("requirement not found")

type RequirementService struct{ repo *repository.RequirementRepository }

func NewRequirementService(repo *repository.RequirementRepository) *RequirementService {
	return &RequirementService{repo: repo}
}

func (s *RequirementService) List(page, pageSize int, reqType string, projectID, customerID *string, status, keyword string) ([]model.Requirement, int64, error) {
	return s.repo.List(page, pageSize, reqType, projectID, customerID, status, keyword)
}

func (s *RequirementService) GetByID(id string) (*model.Requirement, error) {
	m, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRequirementNotFound
		}
		return nil, err
	}
	return m, nil
}

func (s *RequirementService) Create(m *model.Requirement) error { return s.repo.Create(m) }

func (s *RequirementService) Update(id string, u map[string]interface{}) error {
	m, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRequirementNotFound
		}
		return err
	}
	for k, v := range u {
		switch k {
		case "title": m.Title = v.(string)
		case "description": m.Description = v.(string)
		case "priority": m.Priority = v.(string)
		case "status": m.Status = v.(string)
		}
	}
	m.UpdatedAt = time.Now()
	return s.repo.Update(m)
}

func (s *RequirementService) Delete(id string) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrRequirementNotFound
		}
		return err
	}
	return s.repo.SoftDelete(id)
}

func (s *RequirementService) ForceDelete(id string) error { return s.repo.ForceDelete(id) }
func (s *RequirementService) Restore(id string) error { return s.repo.Restore(id) }
