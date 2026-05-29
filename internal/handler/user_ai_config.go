package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/miaomiaopu/ipmp-server/internal/dto/request"
	"github.com/miaomiaopu/ipmp-server/internal/model"
	"github.com/miaomiaopu/ipmp-server/internal/pkg/response"
	"github.com/miaomiaopu/ipmp-server/internal/service"
)

// maskKey 掩码显示 API Key (sk-****xyz)
func maskKey(key string) string {
	if len(key) <= 8 {
		return "****"
	}
	return key[:3] + "****" + key[len(key)-4:]
}

type UserAIConfigHandler struct{ svc *service.UserAIConfigService }

func NewUserAIConfigHandler(svc *service.UserAIConfigService) *UserAIConfigHandler {
	return &UserAIConfigHandler{svc: svc}
}

// Get GET /ai-config — 获取当前用户的 AI 配置
func (h *UserAIConfigHandler) Get(c *gin.Context) {
	userID, _ := c.Get("user_id")
	cfg, err := h.svc.Get(userID.(string))
	if err != nil {
		response.Success(c, gin.H{"configured": false, "provider": "deepseek", "model": "deepseek-v4-flash"})
		return
	}
	response.Success(c, gin.H{
		"configured":  true,
		"provider":    cfg.Provider,
		"model":       cfg.Model,
		"base_url":    cfg.BaseURL,
		"key_preview": maskKey(string(cfg.APIKey)),
		"is_active":   cfg.IsActive,
	})
}

// Update POST /ai-config/update — 更新 AI 配置
func (h *UserAIConfigHandler) Update(c *gin.Context) {
	var req request.UpdateAIConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	cfg := &model.UserAIConfig{
		UserID:   userID.(string),
		Provider: req.Provider,
		APIKey:   model.EncryptedField(req.APIKey),
		Model:    req.Model,
		BaseURL:  req.BaseURL,
		IsActive: true,
	}
	if cfg.Provider == "" {
		cfg.Provider = "deepseek"
	}
	if err := h.svc.Save(cfg); err != nil {
		response.InternalError(c, "failed to save ai config")
		return
	}
	response.Success(c, gin.H{"message": "saved"})
}

// Delete POST /ai-config/delete
func (h *UserAIConfigHandler) Delete(c *gin.Context) {
	userID, _ := c.Get("user_id")
	if err := h.svc.Delete(userID.(string)); err != nil {
		response.InternalError(c, "failed to delete ai config")
		return
	}
	response.Success(c, nil)
}
