package database

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/givetrack/givetrack/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect 连接 MySQL（带启动期重试）并执行自动迁移。
func Connect(dsn string, maxOpen, maxIdle, connMaxLifetime, retryCount, retryInterval int) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	for i := 0; i <= retryCount; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Warn)})
		if err == nil {
			if sqlDB, sqlErr := db.DB(); sqlErr == nil && sqlDB.Ping() == nil {
				break
			}
		}
		if i == retryCount {
			return nil, fmt.Errorf("open mysql after %d retries: %w", retryCount, err)
		}
		slog.Warn("database not ready, retrying", "attempt", i+1, "err", err)
		time.Sleep(time.Duration(retryInterval) * time.Second)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db: %w", err)
	}
	sqlDB.SetMaxOpenConns(maxOpen)
	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetConnMaxLifetime(time.Duration(connMaxLifetime) * time.Minute)

	if err := db.AutoMigrate(
		&model.User{},
		&model.Organization{},
		&model.Project{},
		&model.ProjectUpdate{},
		&model.Donation{},
		&model.AdminReview{},
		&model.VolunteerService{},
	); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}
	return db, nil
}
