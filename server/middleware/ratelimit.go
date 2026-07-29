package middleware

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// rateLimiter 简单的 IP 限流器
type rateLimiter struct {
	mu       sync.Mutex
	visitors map[string]*visitor
	rate     int           // 每个时间窗口允许的最大请求数
	window   time.Duration // 时间窗口
}

type visitor struct {
	count    int
	lastSeen time.Time
}

var loginLimiter = &rateLimiter{
	visitors: make(map[string]*visitor),
	rate:     10,              // 每个窗口最多 10 次请求
	window:   5 * time.Minute, // 5 分钟窗口
}

// RateLimit 请求限流中间件
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		loginLimiter.mu.Lock()

		v, exists := loginLimiter.visitors[ip]
		if !exists {
			loginLimiter.visitors[ip] = &visitor{count: 1, lastSeen: time.Now()}
			loginLimiter.mu.Unlock()
			c.Next()
			return
		}

		// 如果窗口已过期，重置计数
		if time.Since(v.lastSeen) > loginLimiter.window {
			v.count = 1
			v.lastSeen = time.Now()
			loginLimiter.mu.Unlock()
			c.Next()
			return
		}

		// 检查是否超过限制
		if v.count >= loginLimiter.rate {
			loginLimiter.mu.Unlock()
			log.Printf("[限流] IP %s 超过请求限制 (%d次/%v)", ip, loginLimiter.rate, loginLimiter.window)
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    429,
				"message": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}

		v.count++
		v.lastSeen = time.Now()
		loginLimiter.mu.Unlock()

		c.Next()
	}
}

// CleanupExpiredVisitors 定期清理过期的访问者记录
func init() {
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			loginLimiter.mu.Lock()
			now := time.Now()
			for ip, v := range loginLimiter.visitors {
				if now.Sub(v.lastSeen) > loginLimiter.window {
					delete(loginLimiter.visitors, ip)
				}
			}
			loginLimiter.mu.Unlock()
		}
	}()
}
