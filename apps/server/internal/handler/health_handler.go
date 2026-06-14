package handler

import (
	"context"
	"net/http"

	"fullstack-engineering-lab/server/internal/response"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewHealthHandler(db *gorm.DB, rdb *redis.Client) *HealthHandler {
	return &HealthHandler{db: db, rdb: rdb}
}

func (h *HealthHandler) Check(c *gin.Context) {
	data := gin.H{
		"status": "ok",
		"db":     "ok",
		"redis":  "ok",
	}

	// 检查数据库
	sqlDB, err := h.db.DB()
	if err != nil || sqlDB.Ping() != nil {
		data["db"] = "error"
		data["status"] = "degraded"
	}

	// 检查 Redis
	if h.rdb != nil {
		if err := h.rdb.Ping(context.Background()).Err(); err != nil {
			data["redis"] = "error"
			data["status"] = "degraded"
		}
	}

	code := http.StatusOK
	if data["status"] != "ok" {
		code = http.StatusServiceUnavailable
	}

	c.JSON(code, response.Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}
