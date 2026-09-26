package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/givetrack/givetrack/internal/constants"
	"github.com/givetrack/givetrack/internal/util"
)

type rateBucket struct {
	windowStart time.Time
	count       int
}

// RateLimit 基于客户端 IP 的固定窗口限流。
func RateLimit(limit, windowSecs int) gin.HandlerFunc {
	if limit <= 0 {
		limit = 10
	}
	if windowSecs <= 0 {
		windowSecs = 60
	}
	var mu sync.Mutex
	buckets := map[string]*rateBucket{}
	window := time.Duration(windowSecs) * time.Second
	return func(c *gin.Context) {
		key := c.ClientIP() + "|" + c.FullPath()
		now := time.Now()
		mu.Lock()
		if len(buckets) >= 4096 {
			for k, bucket := range buckets {
				if now.Sub(bucket.windowStart) >= window {
					delete(buckets, k)
				}
			}
		}
		b := buckets[key]
		if b == nil || now.Sub(b.windowStart) >= window {
			b = &rateBucket{windowStart: now, count: 0}
			buckets[key] = b
		}
		b.count++
		exceeded := b.count > limit
		mu.Unlock()
		if exceeded {
			util.Fail(c, http.StatusTooManyRequests, constants.CodeTooManyRequests, "too many requests, please slow down")
			c.Abort()
			return
		}
		c.Next()
	}
}
