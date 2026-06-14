package router

import (
	"fullstack-engineering-lab/server/internal/config"
	"fullstack-engineering-lab/server/internal/handler"
	"fullstack-engineering-lab/server/internal/middleware"
	"fullstack-engineering-lab/server/internal/repository"
	"fullstack-engineering-lab/server/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Setup(cfg *config.Config, db *gorm.DB, rdb *redis.Client) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	// Repositories
	userRepo := repository.NewUserRepository(db)

	// Services
	authService := service.NewAuthService(userRepo, rdb, &cfg.JWT)

	// Handlers
	authHandler := handler.NewAuthHandler(authService)
	healthHandler := handler.NewHealthHandler(db, rdb)

	// Routes
	api := r.Group("/api/v1")
	{
		// Public routes
		api.GET("/health", healthHandler.Check)
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)

		// Protected routes
		auth := api.Group("/auth")
		auth.Use(middleware.JWTAuth(cfg.JWT.Secret, rdb))
		{
			auth.GET("/profile", authHandler.Profile)
			auth.POST("/logout", authHandler.Logout)
		}
	}

	return r
}
