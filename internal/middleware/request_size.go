package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/miaomiaopu/ipmp-server/internal/pkg/response"
)

// MaxBodySize 限制请求体最大字节数(默认10MB)
func MaxBodySize(maxBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength > maxBytes {
			response.BadRequest(c, "request body too large")
			c.Abort()
			return
		}
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
		c.Next()
	}
}
