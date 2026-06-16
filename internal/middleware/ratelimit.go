package middleware

import (
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/miaomiaopu/ipmp-server/internal/pkg/response"
)

// RateLimiter 基于内存的简单限流器
type RateLimiter struct {
	mu       sync.Mutex
	window   time.Duration
	maxReq   int
	requests map[string][]time.Time
}

// NewRateLimiter window 内最多 maxReq 次请求
func NewRateLimiter(window time.Duration, maxReq int) *RateLimiter {
	rl := &RateLimiter{
		window:   window,
		maxReq:   maxReq,
		requests: make(map[string][]time.Time),
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, times := range rl.requests {
			valid := times[:0]
			for _, t := range times {
				if now.Sub(t) < rl.window {
					valid = append(valid, t)
				}
			}
			if len(valid) == 0 {
				delete(rl.requests, key)
			} else {
				rl.requests[key] = valid
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	times := rl.requests[key]
	valid := times[:0]
	for _, t := range times {
		if now.Sub(t) < rl.window {
			valid = append(valid, t)
		}
	}
	rl.requests[key] = valid

	if len(valid) >= rl.maxReq {
		return false
	}
	rl.requests[key] = append(rl.requests[key], now)
	return true
}

// LoginRateLimit 登录限流：每 IP 每分钟最多 5 次
func LoginRateLimit() gin.HandlerFunc {
	limiter := NewRateLimiter(time.Minute, 5)
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.Allow(ip) {
			response.TooManyRequests(c, "too many login attempts, please try again later")
			return
		}
		c.Next()
	}
}

// GlobalRateLimit 全局限流：每 IP 每秒最多 100 次
func GlobalRateLimit() gin.HandlerFunc {
	limiter := NewRateLimiter(time.Second, 100)
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.Allow(ip) {
			response.TooManyRequests(c, "too many requests")
			return
		}
		c.Next()
	}
}

// AIRateLimit AI 生成限流：每用户每分钟最多 10 次，未认证时按 IP 兜底。
func AIRateLimit() gin.HandlerFunc {
	limiter := NewRateLimiter(time.Minute, 10)
	return func(c *gin.Context) {
		key := c.GetString("user_id")
		if key == "" {
			key = c.ClientIP()
		}
		if !limiter.Allow(key) {
			response.TooManyRequests(c, "too many ai requests")
			return
		}
		c.Next()
	}
}
