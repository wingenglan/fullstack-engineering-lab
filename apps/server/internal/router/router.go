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
	redisDataService := service.NewRedisDataService(rdb)

	// 处理器
	authHandler := handler.NewAuthHandler(authService)
	lockHandler := handler.NewLockHandler(lockService)
	redisDataHandler := handler.NewRedisDataHandler(redisDataService)
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

			// Redis 数据类型演示路由（公开访问）—— Hash / Set / ZSet / List / String
			rd := api.Group("/redis-data")
			{
				// Hash: 用户画像缓存
				rd.POST("/hash/field", redisDataHandler.SetField)
				rd.GET("/hash/profile", redisDataHandler.GetProfile)
				rd.POST("/hash/multi-set", redisDataHandler.MultiSet)
				rd.POST("/hash/delete-field", redisDataHandler.DeleteField)

				// Set: 标签/收藏夹管理
				rd.POST("/set/add", redisDataHandler.SetAdd)
				rd.POST("/set/remove", redisDataHandler.SetRemove)
				rd.GET("/set/members", redisDataHandler.SetMembers)
				rd.POST("/set/intersect", redisDataHandler.SetIntersect)
				rd.POST("/set/union", redisDataHandler.SetUnion)
				rd.POST("/set/diff", redisDataHandler.SetDiff)

				// ZSet: 实时排行榜
				rd.POST("/zset/add-score", redisDataHandler.ZSetAddScore)
				rd.POST("/zset/top-n", redisDataHandler.ZSetTopN)
				rd.POST("/zset/rank", redisDataHandler.ZSetRank)

				// List: 最新活动流 / 简易消息队列
				rd.POST("/list/push", redisDataHandler.ListPush)
				rd.POST("/list/pop", redisDataHandler.ListPop)
				rd.GET("/list/range", redisDataHandler.ListRange)

				// String: 验证码存储 / 计数器
				rd.POST("/string/set", redisDataHandler.StringSet)
				rd.GET("/string/get", redisDataHandler.StringGet)
				rd.POST("/string/incr", redisDataHandler.StringIncr)
			}
	}

	return r, hub
}
