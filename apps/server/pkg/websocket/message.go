package websocket

import "encoding/json"

// WebSocket 消息协议类型常量
const (
	// 客户端 -> 服务端
	MsgTypeJoinRoom   = "join_room"   // 加入房间
	MsgTypeLeaveRoom  = "leave_room"  // 离开房间
	MsgTypeSendMessage = "send_message" // 发送消息
	MsgTypeTyping     = "typing"       // 输入指示器
	MsgTypePing       = "ping"         // 心跳

	// 服务端 -> 客户端
	MsgTypeNewMessage    = "new_message"    // 新消息广播
	MsgTypeUserJoined    = "user_joined"    // 用户加入通知
	MsgTypeUserLeft      = "user_left"      // 用户离开通知
	MsgTypeOnlineUsers   = "online_users"   // 在线用户列表
	MsgTypeUserTyping    = "user_typing"    // 他人正在输入
	MsgTypePong          = "pong"           // 心跳响应
	MsgTypeError         = "error"          // 错误消息
	MsgTypeRoomHistory   = "room_history"   // 房间历史消息
)

// WSMessage WebSocket 消息协议
type WSMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// JoinRoomPayload 加入房间载荷
type JoinRoomPayload struct {
	RoomID uint `json:"room_id"`
}

// LeaveRoomPayload 离开房间载荷
type LeaveRoomPayload struct {
	RoomID uint `json:"room_id"`
}

// SendMessagePayload 发送消息载荷
type SendMessagePayload struct {
	RoomID  uint   `json:"room_id"`
	Content string `json:"content"`
	MsgType int8   `json:"msg_type"` // 1=文本, 2=图片, 3=系统
}

// TypingPayload 输入指示器载荷
type TypingPayload struct {
	RoomID uint `json:"room_id"`
}

// NewMessagePayload 新消息广播载荷
type NewMessagePayload struct {
	ID        uint   `json:"id"`
	RoomID    uint   `json:"room_id"`
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	Content   string `json:"content"`
	MsgType   int8   `json:"msg_type"`
	CreatedAt string `json:"created_at"`
}

// UserEventPayload 用户事件载荷（加入/离开）
type UserEventPayload struct {
	RoomID   uint   `json:"room_id"`
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
}

// OnlineUsersPayload 在线用户列表载荷
type OnlineUsersPayload struct {
	RoomID uint               `json:"room_id"`
	Users  []OnlineUserSimple `json:"users"`
	Count  int                `json:"count"`
}

// OnlineUserSimple 在线用户简要信息
type OnlineUserSimple struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
}

// UserTypingPayload 他人正在输入载荷
type UserTypingPayload struct {
	RoomID   uint   `json:"room_id"`
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
}

// ErrorPayload 错误载荷
type ErrorPayload struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// RoomHistoryPayload 房间历史消息载荷
type RoomHistoryPayload struct {
	RoomID   uint               `json:"room_id"`
	Messages []HistoryMessage   `json:"messages"`
}

// HistoryMessage 历史消息条目
type HistoryMessage struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	Content   string `json:"content"`
	MsgType   int8   `json:"msg_type"`
	CreatedAt string `json:"created_at"`
}
