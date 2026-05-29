package request

// CreateUserRequest 创建用户
type CreateUserRequest struct {
	Username    string `json:"username" binding:"required,min=2,max=64"`
	Password    string `json:"password" binding:"omitempty,min=6,max=128"`
	DisplayName string `json:"display_name" binding:"omitempty,max=128"`
	Email       string `json:"email" binding:"omitempty,max=256"`
	Phone       string `json:"phone" binding:"omitempty,max=64"`
	Role        string `json:"role" binding:"omitempty,oneof=admin manager user"`
}

// UpdateUserRequest 更新用户
type UpdateUserRequest struct {
	DisplayName *string `json:"display_name"`
	Email       *string `json:"email"`
	Phone       *string `json:"phone"`
	Role        *string `json:"role" binding:"omitempty,oneof=admin manager user"`
	Status      *string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// ChangePasswordRequest 修改密码
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=128"`
}

// ResetPasswordRequest 管理员重置用户密码
type ResetPasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required,min=6,max=128"`
}

// UserQuery 用户列表查询
type UserQuery struct {
	Keyword  string `form:"keyword"`
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
}

func (q *UserQuery) Normalize() (page, pageSize int) {
	page, pageSize = q.Page, q.PageSize
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	return
}
