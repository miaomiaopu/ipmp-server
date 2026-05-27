package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/miaomiaopu/ipmp-server/internal/model"
	"gorm.io/gorm"
)

// AuditLogger 审计日志中间件，记录所有 POST 写操作到 audit_logs 表
func AuditLogger(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 仅记录写操作（POST），跳过 GET
		if c.Request.Method != "POST" {
			c.Next()
			return
		}

		// 解析资源类型和动作
		resource, action := parseAuditAction(c.Request.URL.Path)

		// 读取请求体（截断至 1024 字符）
		var detail string
		if c.Request.Body != nil && c.Request.ContentLength > 0 {
			bodyBytes, _ := io.ReadAll(c.Request.Body)
			c.Request.Body.Close()
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			if len(bodyBytes) > 1024 {
				bodyBytes = bodyBytes[:1024]
			}
			detail = sanitizeAuditBody(string(bodyBytes))
		}

		// 执行后续处理
		c.Next()

		// 仅记录成功的写操作
		status := c.Writer.Status()
		if status < 200 || status >= 300 {
			return
		}

		// 构建审计日志
		userID, _ := c.Get("user_id")
		var uid *string
		if userID != nil {
			s := userID.(string)
			uid = &s
		}

		var detailJSON json.RawMessage
		if detail != "" {
			detailJSON, _ = json.Marshal(map[string]string{"body": detail})
		}

		log := model.AuditLog{
			ID:         uuid.New().String(),
			UserID:     uid,
			Action:     action,
			Resource:   resource,
			ResourceID: resourceIDFromPath(c.Request.URL.Path),
			Detail:     detailJSON,
			IPAddress:  c.ClientIP(),
			CreatedAt:  time.Now(),
		}
		db.Create(&log)
	}
}

// parseAuditAction 从 URL 路径解析资源和动作
// /api/v1/customers → resource=customer, action=create
// /api/v1/customers/:id/update → resource=customer, action=update
// /api/v1/customers/:id/delete → resource=customer, action=delete
func parseAuditAction(path string) (resource, action string) {
	// 去掉 /api/v1/ 前缀
	p := strings.TrimPrefix(path, "/api/v1/")
	parts := strings.Split(strings.Trim(p, "/"), "/")

	if len(parts) == 0 {
		return "unknown", "unknown"
	}

	// first segment = resource name (in plural)
	resource = strings.TrimSuffix(parts[0], "s") // customers→customer, projects→project

	// action from last segment
	last := parts[len(parts)-1]
	switch last {
	case "update":
		action = "update"
	case "delete":
		action = "delete"
	case "status":
		action = "status_change"
	case "login":
		action = "login"
	case "refresh":
		action = "refresh_token"
	case "logout":
		action = "logout"
	default:
		// POST to collection root = create
		if len(parts) == 1 {
			action = "create"
		} else {
			action = "update"
		}
	}

	return resource, action
}

// resourceIDFromPath 从 URL 中提取 UUID 格式的资源 ID
func resourceIDFromPath(path string) *string {
	p := strings.TrimPrefix(path, "/api/v1/")
	parts := strings.Split(strings.Trim(p, "/"), "/")
	for _, part := range parts {
		// UUID 长度 36 (含 -)
		if len(part) == 36 && strings.Count(part, "-") == 4 {
			return &part
		}
	}
	return nil
}

// sanitizeAuditBody 脱敏请求体中的敏感字段
func sanitizeAuditBody(body string) string {
	// 替换敏感字段值为 [REDACTED]
	result := body
	for _, field := range sensitiveFields {
		// 简单替换 JSON 中的敏感值
		result = maskJSONField(result, field)
	}
	return result
}

func maskJSONField(body, field string) string {
	// 查找 "field": "value" 模式并替换
	// 注：这是简化实现，生产环境建议使用 JSON 解析
	idx := strings.Index(strings.ToLower(body), `"`+field+`"`)
	if idx >= 0 {
		colonIdx := strings.Index(body[idx:], ":")
		if colonIdx >= 0 {
			start := idx + colonIdx + 1
			// 跳过引号和空白
			rest := strings.TrimLeft(body[start:], " \t\n\r\"")
			end := strings.IndexAny(rest, "\",\n\r}")
			if end > 0 {
				body = body[:start] + body[start:start+len(body[start:])-len(rest)] + "[REDACTED]" + body[start+len(body[start:])-len(rest)+end:]
			}
		}
	}
	return body
}
