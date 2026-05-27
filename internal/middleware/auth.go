package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/miaomiaopu/ipmp-server/internal/pkg/jwt"
	"github.com/miaomiaopu/ipmp-server/internal/pkg/response"
)

// AuthRequired JWT 认证中间件，从 Authorization: Bearer <token> 解析
func AuthRequired(jwtManager *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			response.Unauthorized(c, "missing or invalid authorization header")
			return
		}
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := jwtManager.ParseAccess(tokenStr)
		if err != nil {
			if err == jwt.ErrTokenExpired {
				response.Unauthorized(c, "token expired")
			} else {
				response.Unauthorized(c, "invalid token")
			}
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// RequireRole 角色校验中间件，需在 AuthRequired 之后使用
func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		roleStr, ok := role.(string)
		if !ok || !allowed[roleStr] {
			response.Forbidden(c, "insufficient permissions")
			return
		}
		c.Next()
	}
}
