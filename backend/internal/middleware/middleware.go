package middleware

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/givetrack/givetrack/internal/service"
	"github.com/givetrack/givetrack/internal/util"
	"github.com/google/uuid"
)

// RequestID 请求 ID 中间件。
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader("X-Request-ID")
		if rid == "" {
			rid = uuid.NewString()
		}
		c.Set("request_id", rid)
		c.Header("X-Request-ID", rid)
		c.Next()
	}
}

// RequestLog 结构化请求日志中间件。
func RequestLog(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.Info("http request",
			"request_id", c.GetString("request_id"),
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
			"client_ip", c.ClientIP(),
		)
	}
}

// Recovery 错误恢复中间件。
func Recovery(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("panic recovered", "panic", r, "path", c.Request.URL.Path)
				c.AbortWithStatusJSON(500, gin.H{"code": 50000, "message": "internal server error"})
			}
		}()
		c.Next()
	}
}

// Auth JWT 鉴权中间件。
func Auth(authSvc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			util.Fail(c, http.StatusUnauthorized, 40100, "missing or invalid authorization header")
			c.Abort()
			return
		}
		claims, err := authSvc.ParseToken(strings.TrimPrefix(header, "Bearer "))
		if err != nil {
			util.Fail(c, http.StatusUnauthorized, 40100, "invalid or expired token")
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Set("username", claims.Subject)
		c.Next()
	}
}

// RequireRole RBAC 角色校验中间件。
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		for _, r := range roles {
			if role == r {
				c.Next()
				return
			}
		}
		util.Fail(c, http.StatusForbidden, 40300, "forbidden: insufficient role")
		c.Abort()
	}
}
