package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/givetrack/givetrack/internal/util"
	"gorm.io/gorm"
)

// HealthHandler 健康检查。
type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// Healthz 存活检查。
func (h *HealthHandler) Healthz(c *gin.Context) {
	util.OK(c, gin.H{
		"status":  "ok",
		"time":    time.Now().Format(time.RFC3339),
		"service": "givetrack",
	})
}

// Readyz 就绪检查（含 DB ping）。
func (h *HealthHandler) Readyz(c *gin.Context) {
	dbStatus := "up"
	if sqlDB, err := h.db.DB(); err != nil || sqlDB.Ping() != nil {
		dbStatus = "down"
	}
	status := "ok"
	if dbStatus != "up" {
		status = "degraded"
	}
	util.OK(c, gin.H{
		"status":  status,
		"db":      dbStatus,
		"time":    time.Now().Format(time.RFC3339),
		"service": "givetrack",
	})
}
