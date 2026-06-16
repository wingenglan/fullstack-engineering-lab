package websocket

import (
	"encoding/json"
	"time"

	gorillaWs "github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	// 写超时
	writeWait = 10 * time.Second

	// 读 pong 超时
	pongWait = 60 * time.Second

	// 心跳发送间隔（必须小于 pongWait）
	pingPeriod = (pongWait * 9) / 10

	// 最大消息大小
	maxMessageSize = 64 * 1024 // 64KB
)

// Client WebSocket 客户端连接
type Client struct {
	Hub      *Hub
	Conn     *gorillaWs.Conn
	send     chan []byte
	UserID   uint
	Username string
}

// NewClient 创建客户端实例
func NewClient(hub *Hub, conn *gorillaWs.Conn, userID uint, username string) *Client {
	return &Client{
		Hub:      hub,
		Conn:     conn,
		send:     make(chan []byte, 256),
		UserID:   userID,
		Username: username,
	}
}

// ReadPump 读取 WebSocket 消息
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if gorillaWs.IsUnexpectedCloseError(err, gorillaWs.CloseGoingAway, gorillaWs.CloseNormalClosure) {
				zap.L().Warn("WebSocket 异常关闭", zap.Error(err))
			}
			break
		}

		var wsMsg WSMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			zap.L().Warn("消息解析失败", zap.Error(err))
			continue
		}

		c.handleMessage(&wsMsg)
	}
}

// WritePump 写入 WebSocket 消息
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub 关闭了通道
				c.Conn.WriteMessage(gorillaWs.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteMessage(gorillaWs.TextMessage, message); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(gorillaWs.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage 处理接收到的消息
func (c *Client) handleMessage(msg *WSMessage) {
	switch msg.Type {
	case MsgTypeJoinRoom:
		var p JoinRoomPayload
		if err := json.Unmarshal(msg.Payload, &p); err != nil {
			return
		}
		c.Hub.JoinRoom(c, p.RoomID)

	case MsgTypeLeaveRoom:
		var p LeaveRoomPayload
		if err := json.Unmarshal(msg.Payload, &p); err != nil {
			return
		}
		c.Hub.LeaveRoom(c, p.RoomID)

	case MsgTypeSendMessage:
		c.Hub.HandleSendMessage(c, msg.Payload)

	case MsgTypeTyping:
		c.Hub.HandleTyping(c, msg.Payload)

	case MsgTypePing:
		pong, _ := json.Marshal(WSMessage{Type: MsgTypePong})
		select {
		case c.send <- pong:
		default:
		}

	default:
		zap.L().Debug("未知消息类型", zap.String("type", msg.Type))
	}
}
