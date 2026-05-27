package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miaomiaopu/ipmp-server/internal/pkg/response"
)

// MethodAllow 仅放行 GET 和 POST 方法，其余返回 405
func MethodAllow() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodPost {
			response.MethodNotAllowed(c)
			return
		}
		c.Next()
	}
}
