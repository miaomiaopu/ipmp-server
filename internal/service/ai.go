package service

import (
	"errors"
	"strings"

	"github.com/miaomiaopu/ipmp-server/internal/config"
	"github.com/miaomiaopu/ipmp-server/internal/repository"
	"gorm.io/gorm"
)

type AIService struct {
	cfg        config.AIConfig
	configRepo *repository.UserAIConfigRepository
}

func NewAIService(cfg config.AIConfig, configRepo *repository.UserAIConfigRepository) *AIService {
	return &AIService{cfg: cfg, configRepo: configRepo}
}

func (s *AIService) GenerateReport(userID, seed string) (map[string]string, error) {
	provider, modelName := s.resolveProvider(userID)
	content := "AI mock 周报草稿"
	if strings.TrimSpace(seed) != "" {
		content += "\n\n" + strings.TrimSpace(seed)
	}
	return map[string]string{
		"provider": provider,
		"model":    modelName,
		"content":  content,
	}, nil
}

func (s *AIService) Summarize(userID, content string) (map[string]string, error) {
	provider, modelName := s.resolveProvider(userID)
	text := strings.TrimSpace(content)
	if len([]rune(text)) > 120 {
		text = string([]rune(text)[:120])
	}
	if text == "" {
		text = "暂无内容"
	}
	return map[string]string{
		"provider": provider,
		"model":    modelName,
		"summary":  "AI mock 摘要: " + text,
	}, nil
}

func (s *AIService) resolveProvider(userID string) (string, string) {
	if s.configRepo != nil {
		cfg, err := s.configRepo.FindByUserID(userID)
		if err == nil && cfg != nil && cfg.IsActive {
			return cfg.Provider, cfg.Model
		}
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return s.cfg.Provider, ""
		}
	}
	switch s.cfg.Provider {
	case "deepseek":
		return s.cfg.Provider, s.cfg.DeepSeek.Model
	case "openai":
		return s.cfg.Provider, s.cfg.OpenAI.Model
	case "claude":
		return s.cfg.Provider, s.cfg.Claude.Model
	default:
		return "mock", "mock"
	}
}
