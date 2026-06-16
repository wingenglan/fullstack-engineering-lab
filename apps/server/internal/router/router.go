package router

import (
	"fullstack-engineering-lab/server/internal/config"
	"fullstack-engineering-lab/server/internal/handler"
	"fullstack-engineering-lab/server/internal/middleware"
	"fullstack-engineering-lab/server/internal/repository"
	"fullstack-engineering-lab/server/internal/service"
	ws "fullstack-engineering-lab/server/pkg/websocket"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Setup(cfg *config.Config, db *gorm.DB, rdb *redis.Client) (*gin.Engine, *ws.Hub) {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORS())

	// WebSocket Hub
	hub := ws.NewHub()

	// 数据层
	userRepo := repository.NewUserRepository(db)
	chatRepo := repository.NewChatRepository(db)

	// 服务层
	authService := service.NewAuthService(userRepo, rdb, &cfg.JWT)
	lockService := service.NewLockService(rdb)
	chatService := service.NewChatService(chatRepo, hub)

	// 处理器
	authHandler := handler.NewAuthHandler(authService)
	lockHandler := handler.NewLockHandler(lockService)
	healthHandler := handler.NewHealthHandler(db, rdb)
	chatHandler := handler.NewChatHandler(chatService)

	// 路由注册
	api := r.Group("/api/v1")
	{
		// 公开路由
		api.GET("/health", healthHandler.Check)
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)

		// 需认证路由
		auth := api.Group("/auth")
		auth.Use(middleware.JWTAuth(cfg.JWT.Secret, rdb))
		{
			auth.GET("/profile", authHandler.Profile)
			auth.POST("/logout", authHandler.Logout)
		}

		// Redis 分布式锁演示路由（公开访问）
		lock := api.Group("/lock")
		{
			lock.POST("/acquire", lockHandler.Acquire)
			lock.POST("/release", lockHandler.Release)
			lock.POST("/status", lockHandler.Status)
			lock.POST("/contention", lockHandler.ContentionDemo)
		}

		// WebSocket 实时通讯路由
		chat := api.Group("/chat")
		{
			// 公开接口（获取房间列表、消息历史）
			chat.GET("/rooms", chatHandler.GetRooms)
			chat.GET("/rooms/:id", chatHandler.GetRoomInfo)
			chat.GET("/rooms/:id/online", chatHandler.GetOnlineUsers)
			chat.GET("/messages", chatHandler.GetMessageHistory)

			// 需认证接口
			chatAuth := chat.Group("")
			chatAuth.Use(middleware.JWTAuth(cfg.JWT.Secret, rdb))
			{
				chatAuth.POST("/rooms", chatHandler.CreateRoom)
			}

			// WebSocket 连接（使用 WSAuth 从 URL 参数读取 Token）
			chat.GET("/ws", middleware.WSAuth(cfg.JWT.Secret, rdb), chatHandler.WS)
		}
	}

	return r, hub
}
