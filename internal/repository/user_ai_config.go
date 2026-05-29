package repository

import (
	"github.com/miaomiaopu/ipmp-server/internal/model"
	"gorm.io/gorm"
)

type UserAIConfigRepository struct{ db *gorm.DB }

func NewUserAIConfigRepository(db *gorm.DB) *UserAIConfigRepository {
	return &UserAIConfigRepository{db: db}
}

func (r *UserAIConfigRepository) FindByUserID(userID string) (*model.UserAIConfig, error) {
	var cfg model.UserAIConfig
	err := r.db.Where("user_id = ?", userID).First(&cfg).Error
	return &cfg, err
}

func (r *UserAIConfigRepository) CreateOrUpdate(cfg *model.UserAIConfig) error {
	var existing model.UserAIConfig
	err := r.db.Where("user_id = ?", cfg.UserID).First(&existing).Error
	if err == nil {
		cfg.ID = existing.ID
		cfg.CreatedAt = existing.CreatedAt
	}
	return r.db.Save(cfg).Error
}

func (r *UserAIConfigRepository) DeleteByUserID(userID string) error {
	return r.db.Where("user_id = ?", userID).Delete(&model.UserAIConfig{}).Error
}
