package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/miaomiaopu/ipmp-server/internal/dto/request"
	dtoResp "github.com/miaomiaopu/ipmp-server/internal/dto/response"
	"github.com/miaomiaopu/ipmp-server/internal/pkg/response"
	"github.com/miaomiaopu/ipmp-server/internal/service"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Login POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	user, accessToken, refreshToken, err := h.svc.Login(req.Username, req.Password)
	if err != nil {
		if err == service.ErrInvalidCredentials || err == service.ErrUserInactive {
			response.Unauthorized(c, err.Error())
			return
		}
		response.InternalError(c, "login failed")
		return
	}
	response.Success(c, dtoResp.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    86400,
		User: dtoResp.UserBrief{
			ID:          user.ID,
			Username:    user.Username,
			DisplayName: user.DisplayName,
			Role:        user.Role,
		},
	})
}

// Refresh POST /auth/refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	var req request.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	user, accessToken, refreshToken, err := h.svc.RefreshToken(req.RefreshToken)
	if err != nil {
		response.Unauthorized(c, "invalid refresh token")
		return
	}
	response.Success(c, dtoResp.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    86400,
		User: dtoResp.UserBrief{
			ID:          user.ID,
			Username:    user.Username,
			DisplayName: user.DisplayName,
			Role:        user.Role,
		},
	})
}

// Logout POST /auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	response.Success(c, gin.H{"message": "logged out"})
}

// Me GET /auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := c.Get("user_id")
	user, err := h.svc.GetCurrentUser(userID.(string))
	if err != nil {
		response.NotFound(c, "user not found")
		return
	}
	response.Success(c, dtoResp.UserBrief{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Role:        user.Role,
	})
}
