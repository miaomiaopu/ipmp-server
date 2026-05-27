package request

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required,min=2,max=64"`
	Password string `json:"password" binding:"required,min=6,max=128"`
}

// RefreshRequest 刷新Token请求
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
