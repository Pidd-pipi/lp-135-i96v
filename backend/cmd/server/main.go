package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/givetrack/givetrack/internal/config"
	"github.com/givetrack/givetrack/internal/database"
	"github.com/givetrack/givetrack/internal/repository"
	"github.com/givetrack/givetrack/internal/router"
	"github.com/givetrack/givetrack/internal/service"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		logger.Error("load config failed", "err", err)
		os.Exit(1)
	}

	db, err := database.Connect(cfg.DSN(), cfg.DBMaxOpenConns, cfg.DBMaxIdleConns, cfg.DBConnMaxLifetime, cfg.DBRetryCount, cfg.DBRetryInterval)
	if err != nil {
		logger.Error("connect database failed", "err", err)
		os.Exit(1)
	}
	if err := database.Seed(db); err != nil {
		logger.Error("seed database failed", "err", err)
		os.Exit(1)
	}

	userRepo := repository.NewUserRepository(db)
	orgRepo := repository.NewOrganizationRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	updateRepo := repository.NewProjectUpdateRepository(db)
	donationRepo := repository.NewDonationRepository(db)
	reviewRepo := repository.NewAdminReviewRepository(db)

	authSvc := service.NewAuthService(userRepo, orgRepo, cfg.JWTSecret, cfg.JWTExpire, logger)
	projectSvc := service.NewProjectService(projectRepo, updateRepo, orgRepo, donationRepo, logger)
	donationSvc := service.NewDonationService(db, donationRepo, projectRepo, userRepo, logger)
	rankingSvc := service.NewRankingService(userRepo, logger)
	adminSvc := service.NewAdminService(projectRepo, orgRepo, reviewRepo, logger)

	engine := router.Setup(db, authSvc, projectSvc, donationSvc, rankingSvc, adminSvc, cfg, logger)

	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      engine,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("server starting", "port", cfg.ServerPort, "env", cfg.AppEnv)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server listen failed", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server shutdown failed", "err", err)
	}
	if sqlDB, err := db.DB(); err == nil {
		_ = sqlDB.Close()
	}
	logger.Info("server exited")
}
