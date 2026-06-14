package middleware

import (
	"net/http"
	"strings"

	"fullstack-engineering-lab/server/internal/response"
	appJwt "fullstack-engineering-lab/server/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func JWTAuth(jwtSecret string, rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, response.CodeTokenExpired, "missing authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, response.CodeTokenExpired, "invalid authorization format")
			c.Abort()
			return
		}

		tokenString := parts[1]

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
