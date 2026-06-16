package handler

import (
	"net/http"

	"fullstack-engineering-lab/server/internal/response"
	"fullstack-engineering-lab/server/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	ws "fullstack-engineering-lab/server/pkg/websocket"
)

// handshake 执行 WebSocket 连接握手
func handshake(userID uint, username string, chatService *service.ChatService, c *gin.Context) {
	hub := chatService.GetHub()

	// 如果中间件未提供用户名，通过 service 查询
	if username == "" {
		username = chatService.GetUsername(userID)
	}

	conn, err := ws.Upgrade(c.Writer, c.Request)
	if err != nil {
		zap.L().Error("WebSocket 升级失败", zap.Error(err))
		response.Error(c, http.StatusInternalServerError, response.CodeChatError, "WebSocket 连接失败")
		return
	}

	client := ws.NewClient(hub, conn, userID, username)

	// 注册客户端到 Hub
	hub.Register(client)

	// 启动读写协程
	go client.WritePump()
	go client.ReadPump()
}
