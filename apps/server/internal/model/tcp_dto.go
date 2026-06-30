package model

// TCPSessionResponse 创建会话响应
type TCPSessionResponse struct {
	SessionID string `json:"session_id"`
	CreatedAt string `json:"created_at"`
}

// TCPCommandRequest 发送命令请求
type TCPCommandRequest struct {
	Command string `json:"command" binding:"required"`
}

// TCPCommandResponse 命令执行响应
type TCPCommandResponse struct {
	SessionID  string `json:"session_id"`
	Command    string `json:"command"`
	Response   string `json:"response"`
	DurationMs int64  `json:"duration_ms"`
	Timestamp  string `json:"timestamp"`
}

// TCPChatSessionRequest 创建聊天会话请求
type TCPChatSessionRequest struct {
	Nickname string `json:"nickname" binding:"required"`
	Room     string `json:"room" binding:"required"`
}

// TCPChatMessageRequest 发送聊天消息请求
type TCPChatMessageRequest struct {
	Content string `json:"content" binding:"required"`
}

// TCPSessionInfo 会话信息
type TCPSessionInfo struct {
	SessionID    string `json:"session_id"`
	RemoteAddr   string `json:"remote_addr"`
	CreatedAt    string `json:"created_at"`
	LastActAt    string `json:"last_act_at"`
	DurationSec  int    `json:"duration_sec"`
	CommandCount int64  `json:"command_count"`
	IsAlive      bool   `json:"is_alive"`
}

// TCPStatsResponse 服务器统计响应
type TCPStatsResponse struct {
	ServerAddr     string `json:"server_addr"`
	ActiveSessions int    `json:"active_sessions"`
	MaxSessions    int    `json:"max_sessions"`
	UptimeSec      int64  `json:"uptime_sec"`
}
