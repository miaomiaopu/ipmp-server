package middleware

import (
	"bytes"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
)

// sensitiveFields 日志中需要脱敏的字段
var sensitiveFields = []string{
	"password", "api_key", "token", "secret",
	"access_token", "refresh_token",
	"contact_phone", "contact_email", "address",
	"email", "phone",
}

// SensitiveLog 从请求体中过滤敏感字段写入日志上下文
func SensitiveLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body != nil && c.Request.ContentLength > 0 && c.Request.ContentLength < 65536 {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				_ = c.Request.Body.Close()
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				// 简化的日志脱敏标记
				c.Set("body_sensitive", containsSensitiveKeys(string(bodyBytes)))
			}
		}
		c.Next()
	}
}

func containsSensitiveKeys(body string) bool {
	lower := strings.ToLower(body)
	for _, field := range sensitiveFields {
		if strings.Contains(lower, strings.ToLower(field)) {
			return true
		}
	}
	return false
}
