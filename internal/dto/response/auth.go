package response

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	User         UserBrief `json:"user"`
}

// UserBrief 用户简要信息
type UserBrief struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Role        string `json:"role"`
}

// AIKeyPreview AI Key 掩码预览
type AIKeyPreview struct {
	Configured bool   `json:"configured"`
	Provider   string `json:"provider"`
	Model      string `json:"model"`
	KeyPreview string `json:"key_preview"`
	IsActive   bool   `json:"is_active"`
}
