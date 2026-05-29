package service

import (
	"github.com/miaomiaopu/ipmp-server/internal/model"
	"github.com/miaomiaopu/ipmp-server/internal/repository"
)

type UserAIConfigService struct{ repo *repository.UserAIConfigRepository }

func NewUserAIConfigService(repo *repository.UserAIConfigRepository) *UserAIConfigService {
	return &UserAIConfigService{repo: repo}
}

func (s *UserAIConfigService) Get(userID string) (*model.UserAIConfig, error) {
	return s.repo.FindByUserID(userID)
}

func (s *UserAIConfigService) Save(cfg *model.UserAIConfig) error {
	return s.repo.CreateOrUpdate(cfg)
}

func (s *UserAIConfigService) Delete(userID string) error {
	return s.repo.DeleteByUserID(userID)
}
