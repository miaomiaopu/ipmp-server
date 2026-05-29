package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/miaomiaopu/ipmp-server/internal/dto/request"
	"github.com/miaomiaopu/ipmp-server/internal/model"
	"github.com/miaomiaopu/ipmp-server/internal/pkg/response"
	"github.com/miaomiaopu/ipmp-server/internal/service"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// List GET /users
func (h *UserHandler) List(c *gin.Context) {
	var query request.UserQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	page, pageSize := query.Normalize()
	users, total, err := h.svc.List(page, pageSize, query.Keyword)
	if err != nil {
		response.InternalError(c, "failed to list users")
		return
	}
	response.SuccessWithPagination(c, users, page, pageSize, total)
}

// GetByID GET /users/:id
func (h *UserHandler) GetByID(c *gin.Context) {
	user, err := h.svc.GetByID(c.Param("id"))
	if err != nil {
		if err == service.ErrUserNotFound {
			response.NotFound(c, "user not found")
			return
		}
		response.InternalError(c, "failed to get user")
		return
	}
	response.Success(c, user)
}

// Create POST /users
func (h *UserHandler) Create(c *gin.Context) {
	var req request.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	user := &model.User{
		Username:    req.Username,
		DisplayName: req.DisplayName,
		Email:       model.EncryptedField(req.Email),
		Phone:       model.EncryptedField(req.Phone),
		Role:        req.Role,
		Status:      "active",
	}
	if user.Role == "" {
		user.Role = "user"
	}
	created, password, err := h.svc.Create(user, req.Password)
	if err != nil {
		if err == service.ErrUserExists {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, "failed to create user")
		return
	}
	result := gin.H{"id": created.ID, "username": created.Username, "display_name": created.DisplayName, "role": created.Role}
	if password != "" {
		result["initial_password"] = password // 仅在创建时返回一次
	}
	response.Success(c, result)
}

// Update POST /users/:id/update
func (h *UserHandler) Update(c *gin.Context) {
	var req request.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	updates := make(map[string]interface{})
	if req.DisplayName != nil {
		updates["display_name"] = *req.DisplayName
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.Phone != nil {
		updates["phone"] = *req.Phone
	}
	if req.Role != nil {
		updates["role"] = *req.Role
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if len(updates) == 0 {
		response.BadRequest(c, "no fields to update")
		return
	}
	if err := h.svc.Update(c.Param("id"), updates); err != nil {
		if err == service.ErrUserNotFound {
			response.NotFound(c, "user not found")
			return
		}
		response.InternalError(c, "failed to update user")
		return
	}
	response.Success(c, nil)
}

// Delete POST /users/:id/delete
func (h *UserHandler) Delete(c *gin.Context) {
	if err := h.svc.Delete(c.Param("id")); err != nil {
		if err == service.ErrUserNotFound {
			response.NotFound(c, "user not found")
			return
		}
		response.InternalError(c, "failed to delete user")
		return
	}
	response.Success(c, nil)
}

// ResetPassword POST /users/:id/reset-password
func (h *UserHandler) ResetPassword(c *gin.Context) {
	var req request.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	password, err := h.svc.ResetPassword(c.Param("id"), req.NewPassword)
	if err != nil {
		if err == service.ErrUserNotFound {
			response.NotFound(c, "user not found")
			return
		}
		response.InternalError(c, "failed to reset password")
		return
	}
	response.Success(c, gin.H{"new_password": password})
}

// ChangePassword POST /auth/change-password
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req request.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, _ := c.Get("user_id")
	if err := h.svc.ChangePassword(userID.(string), req.OldPassword, req.NewPassword); err != nil {
		if err == service.ErrWrongPassword {
			response.BadRequest(c, "wrong password")
			return
		}
		response.InternalError(c, "failed to change password")
		return
	}
	response.Success(c, nil)
}
