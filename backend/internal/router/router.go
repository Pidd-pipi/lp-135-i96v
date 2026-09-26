package router

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/givetrack/givetrack/internal/config"
	"github.com/givetrack/givetrack/internal/constants"
	"github.com/givetrack/givetrack/internal/handler"
	"github.com/givetrack/givetrack/internal/middleware"
	"github.com/givetrack/givetrack/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// Setup 组装路由。
func Setup(
	db *gorm.DB,
	authSvc *service.AuthService,
	projectSvc *service.ProjectService,
	donationSvc *service.DonationService,
	rankingSvc *service.RankingService,
	adminSvc *service.AdminService,
	cfg *config.Config,
	logger *slog.Logger,
) *gin.Engine {
	r := gin.New()
	r.Use(middleware.RequestID(), middleware.RequestLog(logger), middleware.Recovery(logger), middleware.CORS(cfg.CORSOrigins()))

	healthHandler := handler.NewHealthHandler(db)
	authHandler := handler.NewAuthHandler(authSvc)
	projectHandler := handler.NewProjectHandler(projectSvc)
	donationHandler := handler.NewDonationHandler(donationSvc)
	rankingHandler := handler.NewRankingHandler(rankingSvc)
	adminHandler := handler.NewAdminHandler(adminSvc)

	r.GET("/healthz", healthHandler.Healthz)
	r.GET("/readyz", healthHandler.Readyz)

	v1 := r.Group("/api/" + constants.APIVersion)
	v1.GET("/healthz", healthHandler.Healthz)
	v1.GET("/readyz", healthHandler.Readyz)

	auth := v1.Group("/auth")
	{
		auth.POST("/register", middleware.RateLimit(cfg.AuthRateLimit, cfg.AuthRateWindowSecs), authHandler.Register)
		auth.POST("/login", middleware.RateLimit(cfg.AuthRateLimit, cfg.AuthRateWindowSecs), authHandler.Login)
		auth.GET("/me", middleware.Auth(authSvc), authHandler.Me)
		auth.PUT("/profile", middleware.Auth(authSvc), authHandler.UpdateProfile)
	}

	projects := v1.Group("/projects")
	{
		projects.GET("", projectHandler.List)
		projects.GET("/org/my", middleware.Auth(authSvc), middleware.RequireRole(constants.RoleOrg), projectHandler.MyProjects)
		projects.POST("", middleware.Auth(authSvc), middleware.RequireRole(constants.RoleOrg), projectHandler.Create)
		projects.GET("/:id", projectHandler.GetDetail)
		projects.GET("/:id/updates", projectHandler.Updates)
		projects.POST("/:id/updates", middleware.Auth(authSvc), middleware.RequireRole(constants.RoleOrg), projectHandler.CreateUpdate)
	}

	donations := v1.Group("/donations")
	{
		donations.POST("", middleware.Auth(authSvc), donationHandler.Create)
		donations.GET("/my", middleware.Auth(authSvc), donationHandler.My)
		donations.GET("/:id/certificate", middleware.Auth(authSvc), donationHandler.Certificate)
	}

	ranking := v1.Group("/ranking")
	{
		ranking.GET("/donation", rankingHandler.DonationRanking)
		ranking.GET("/service", rankingHandler.ServiceRanking)
		ranking.GET("/stats", rankingHandler.Stats)
	}

	admin := v1.Group("/admin", middleware.Auth(authSvc), middleware.RequireRole(constants.RoleAdmin))
	{
		admin.GET("/projects/pending", adminHandler.PendingProjects)
		admin.POST("/projects/:id/review", adminHandler.ReviewProject)
		admin.GET("/organizations/pending", adminHandler.PendingOrganizations)
		admin.POST("/organizations/:id/review", adminHandler.ReviewOrganization)
	}

	r.GET("/swagger/doc.json", healthHandler.SwaggerJSON)
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))
	return r
}
