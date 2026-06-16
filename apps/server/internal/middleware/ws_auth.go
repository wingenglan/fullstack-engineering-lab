package middleware

import (
	"net/http"

	"fullstack-engineering-lab/server/internal/response"
	appJwt "fullstack-engineering-lab/server/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// WSAuth WebSocket 连接认证中间件
// 从 URL 查询参数中读取 Token（浏览器 WebSocket API 不支持自定义 Header）
func WSAuth(jwtSecret string, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Query("token")
		if tokenString == "" {
			response.Error(c, http.StatusUnauthorized, response.CodeTokenExpired, "缺少 token 参数")
			c.Abort()
			return
		}

		// 检查 Redis 黑名单
		if rdb != nil {
			ctx := c.Request.Context()
			exists, err := rdb.Exists(ctx, "blacklist:"+tokenString).Result()
			if err == nil && exists > 0 {
				response.Error(c, http.StatusUnauthorized, response.CodeTokenExpired, "token has been revoked")
				c.Abort()
				return
			}
		}

		claims, err := appJwt.ParseToken(tokenString, jwtSecret)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, response.CodeTokenExpired, "invalid or expired token")
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("token", tokenString)
		c.Next()
	}
}
