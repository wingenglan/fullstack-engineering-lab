package websocket

import (
	"encoding/json"
	"sync"
	"time"

	"go.uber.org/zap"
)

// ClientInfo 客户端关联的用户信息
type ClientInfo struct {
	UserID   uint
	Username string
}

// RoomInfo 房间内的客户端集合
type RoomInfo struct {
	Clients map[*Client]bool
	mu      sync.RWMutex
}

// Hub WebSocket 连接管理中心
type Hub struct {
	// 已注册的客户端集合
	clients map[*Client]bool

	// 房间映射：roomID -> RoomInfo
	rooms map[uint]*RoomInfo

	// 用户到客户端的映射：userID -> []*Client（同一用户可能多标签页）
	userClients map[uint][]*Client

	// 注册通道
	register chan *Client

	// 注销通道
	unregister chan *Client

	// 广播通道（指定房间）
	broadcast chan *RoomMessage

	// 互斥锁
	mu sync.RWMutex

	// 消息存储回调：由 service 层注入
	onSaveMessage func(roomID, userID uint, content string, msgType int8) (uint, error)

	// 查询用户信息回调
	onGetUsername func(userID uint) string

	// 查询房间历史回调
	onGetHistory func(roomID uint, limit int) ([]HistoryMessage, error)
}

// RoomMessage 房间消息
type RoomMessage struct {
	RoomID  uint
	Payload []byte
	Exclude *Client // 排除发送者自身
}

// NewHub 创建 Hub 实例
func NewHub() *Hub {
	return &Hub{
		clients:     make(map[*Client]bool),
		rooms:       make(map[uint]*RoomInfo),
		userClients: make(map[uint][]*Client),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		broadcast:   make(chan *RoomMessage, 256),
	}
}

// SetSaveMessageCallback 设置消息持久化回调
func (h *Hub) SetSaveMessageCallback(fn func(roomID, userID uint, content string, msgType int8) (uint, error)) {
	h.onSaveMessage = fn
}

// SetGetUsernameCallback 设置查询用户名回调
func (h *Hub) SetGetUsernameCallback(fn func(userID uint) string) {
	h.onGetUsername = fn
}

// SetGetHistoryCallback 设置查询历史消息回调
func (h *Hub) SetGetHistoryCallback(fn func(roomID uint, limit int) ([]HistoryMessage, error)) {
	h.onGetHistory = fn
}

// Run 启动 Hub 主循环
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.handleRegister(client)
		case client := <-h.unregister:
			h.handleUnregister(client)
		case msg := <-h.broadcast:
			h.handleBroadcast(msg)
		}
	}
}

// Register 注册客户端到 Hub（公开方法，供 handler 调用）
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// handleRegister 处理客户端注册
func (h *Hub) handleRegister(client *Client) {
	h.mu.Lock()
	h.clients[client] = true
	h.userClients[client.UserID] = append(h.userClients[client.UserID], client)
	h.mu.Unlock()

	zap.L().Info("WebSocket 客户端已连接",
		zap.Uint("userID", client.UserID),
		zap.String("username", client.Username),
	)
}

// handleUnregister 处理客户端注销
func (h *Hub) handleUnregister(client *Client) {
	// 收集需要通知的房间，避免在持锁状态下发 broadcast
	var notifyRooms []uint

	h.mu.Lock()
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)

		// 从用户映射中移除
		clients := h.userClients[client.UserID]
		for i, c := range clients {
			if c == client {
				h.userClients[client.UserID] = append(clients[:i], clients[i+1:]...)
				break
			}
		}
		if len(h.userClients[client.UserID]) == 0 {
			delete(h.userClients, client.UserID)
		}

		// 从所有房间中移除，记录需要通知的房间
		for roomID, room := range h.rooms {
			room.mu.Lock()
			if room.Clients[client] {
				delete(room.Clients, client)
				notifyRooms = append(notifyRooms, roomID)
			}
			room.mu.Unlock()
		}
	}
	h.mu.Unlock()

	// 释放锁后再发送广播通知
	for _, roomID := range notifyRooms {
		h.broadcastUserLeft(roomID, client)
		h.broadcastOnlineUsers(roomID)
	}

	zap.L().Info("WebSocket 客户端已断开",
		zap.Uint("userID", client.UserID),
		zap.String("username", client.Username),
	)
}

// handleBroadcast 处理房间消息广播
func (h *Hub) handleBroadcast(msg *RoomMessage) {
	h.mu.RLock()
	room, ok := h.rooms[msg.RoomID]
	h.mu.RUnlock()
	if !ok {
		return
	}

	room.mu.RLock()
	defer room.mu.RUnlock()

	var slowClients []*Client

	for client := range room.Clients {
		if client == msg.Exclude {
			continue
		}
		select {
		case client.send <- msg.Payload:
		default:
			// 发送缓冲已满，标记待清理
			slowClients = append(slowClients, client)
		}
	}

	// 异步关闭慢客户端，避免在持锁状态下发 unregister
	for _, client := range slowClients {
		go func(c *Client) {
			h.unregister <- c
		}(client)
	}
}

// JoinRoom 加入房间
func (h *Hub) JoinRoom(client *Client, roomID uint) {
	h.mu.Lock()
	room, ok := h.rooms[roomID]
	if !ok {
		room = &RoomInfo{Clients: make(map[*Client]bool)}
		h.rooms[roomID] = room
	}
	h.mu.Unlock()

	room.mu.Lock()
	room.Clients[client] = true
	room.mu.Unlock()

	// 发送历史消息
	if h.onGetHistory != nil {
		messages, err := h.onGetHistory(roomID, 50)
		if err == nil && len(messages) > 0 {
			payload, _ := json.Marshal(RoomHistoryPayload{
				RoomID:   roomID,
				Messages: messages,
			})
			msg, _ := json.Marshal(WSMessage{
				Type:    MsgTypeRoomHistory,
				Payload: payload,
			})
			client.send <- msg
		}
	}

	// 广播用户加入
	h.broadcastUserJoined(roomID, client)

	// 更新在线用户列表
	h.broadcastOnlineUsers(roomID)
}

// LeaveRoom 离开房间
func (h *Hub) LeaveRoom(client *Client, roomID uint) {
	h.mu.RLock()
	room, ok := h.rooms[roomID]
	h.mu.RUnlock()
	if !ok {
		return
	}

	room.mu.Lock()
	delete(room.Clients, client)
	isEmpty := len(room.Clients) == 0
	room.mu.Unlock()

	// 广播用户离开
	h.broadcastUserLeft(roomID, client)

	// 更新在线用户列表
	h.broadcastOnlineUsers(roomID)

	// 清理空房间
	if isEmpty {
		h.mu.Lock()
		// 再次检查（竞态条件防护）
		room.mu.RLock()
		if len(room.Clients) == 0 {
			delete(h.rooms, roomID)
		}
		room.mu.RUnlock()
		h.mu.Unlock()
	}
}

// HandleSendMessage 处理发送消息
func (h *Hub) HandleSendMessage(client *Client, payload []byte) {
	var p SendMessagePayload
	if err := json.Unmarshal(payload, &p); err != nil {
		h.sendError(client, 400, "消息格式无效")
		return
	}

	if p.Content == "" {
		h.sendError(client, 400, "消息内容不能为空")
		return
	}

	if len(p.Content) > 5000 {
		h.sendError(client, 400, "消息内容不能超过 5000 字符")
		return
	}

	if p.MsgType == 0 {
		p.MsgType = 1 // 默认文本消息
	}

	// 持久化消息
	var msgID uint
	if h.onSaveMessage != nil {
		var err error
		msgID, err = h.onSaveMessage(p.RoomID, client.UserID, p.Content, p.MsgType)
		if err != nil {
			h.sendError(client, 500, "消息保存失败")
			return
		}
	}

	// 广播新消息
	newMsg := NewMessagePayload{
		ID:        msgID,
		RoomID:    p.RoomID,
		UserID:    client.UserID,
		Username:  client.Username,
		Content:   p.Content,
		MsgType:   p.MsgType,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	payloadBytes, _ := json.Marshal(newMsg)
	wsMsg, _ := json.Marshal(WSMessage{
		Type:    MsgTypeNewMessage,
		Payload: payloadBytes,
	})

	h.broadcast <- &RoomMessage{
		RoomID:  p.RoomID,
		Payload: wsMsg,
	}
}

// HandleTyping 处理输入指示器
func (h *Hub) HandleTyping(client *Client, payload []byte) {
	var p TypingPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return
	}

	typingPayload, _ := json.Marshal(UserTypingPayload{
		RoomID:   p.RoomID,
		UserID:   client.UserID,
		Username: client.Username,
	})

	wsMsg, _ := json.Marshal(WSMessage{
		Type:    MsgTypeUserTyping,
		Payload: typingPayload,
	})

	h.broadcast <- &RoomMessage{
		RoomID:  p.RoomID,
		Payload: wsMsg,
		Exclude: client,
	}
}

// GetOnlineUsers 获取房间在线用户列表
func (h *Hub) GetOnlineUsers(roomID uint) []OnlineUserSimple {
	h.mu.RLock()
	room, ok := h.rooms[roomID]
	h.mu.RUnlock()

	if !ok {
		return nil
	}

	room.mu.RLock()
	defer room.mu.RUnlock()

	users := make([]OnlineUserSimple, 0, len(room.Clients))
	seen := make(map[uint]bool)
	for client := range room.Clients {
		if !seen[client.UserID] {
			users = append(users, OnlineUserSimple{
				UserID:   client.UserID,
				Username: client.Username,
			})
			seen[client.UserID] = true
		}
	}
	return users
}

// GetOnlineUserCount 获取指定房间在线用户数
func (h *Hub) GetOnlineUserCount(roomID uint) int {
	return len(h.GetOnlineUsers(roomID))
}

// GetTotalOnlineCount 获取总在线连接数
func (h *Hub) GetTotalOnlineCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// broadcastUserJoined 广播用户加入事件
func (h *Hub) broadcastUserJoined(roomID uint, client *Client) {
	payload, _ := json.Marshal(UserEventPayload{
		RoomID:   roomID,
		UserID:   client.UserID,
		Username: client.Username,
	})
	wsMsg, _ := json.Marshal(WSMessage{
		Type:    MsgTypeUserJoined,
		Payload: payload,
	})
	h.broadcast <- &RoomMessage{
		RoomID:  roomID,
		Payload: wsMsg,
		Exclude: client,
	}
}

// broadcastUserLeft 广播用户离开事件
func (h *Hub) broadcastUserLeft(roomID uint, client *Client) {
	payload, _ := json.Marshal(UserEventPayload{
		RoomID:   roomID,
		UserID:   client.UserID,
		Username: client.Username,
	})
	wsMsg, _ := json.Marshal(WSMessage{
		Type:    MsgTypeUserLeft,
		Payload: payload,
	})
	h.broadcast <- &RoomMessage{
		RoomID:  roomID,
		Payload: wsMsg,
	}
}

// broadcastOnlineUsers 广播在线用户列表更新
func (h *Hub) broadcastOnlineUsers(roomID uint) {
	users := h.GetOnlineUsers(roomID)
	payload, _ := json.Marshal(OnlineUsersPayload{
		RoomID: roomID,
		Users:  users,
		Count:  len(users),
	})
	wsMsg, _ := json.Marshal(WSMessage{
		Type:    MsgTypeOnlineUsers,
		Payload: payload,
	})
	h.broadcast <- &RoomMessage{
		RoomID:  roomID,
		Payload: wsMsg,
	}
}

// sendError 向客户端发送错误消息
func (h *Hub) sendError(client *Client, code int, message string) {
	payload, _ := json.Marshal(ErrorPayload{
		Code:    code,
		Message: message,
	})
	wsMsg, _ := json.Marshal(WSMessage{
		Type:    MsgTypeError,
		Payload: payload,
	})
	select {
	case client.send <- wsMsg:
	default:
	}
}
