package websocket

import (
	"net/http"

	gorillaWs "github.com/gorilla/websocket"
)

// WebSocket 连接升级器配置
var upgrader = gorillaWs.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 允许所有来源（开发环境）
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Upgrade HTTP 连接升级为 WebSocket
func Upgrade(w http.ResponseWriter, r *http.Request) (*gorillaWs.Conn, error) {
	return upgrader.Upgrade(w, r, nil)
}
