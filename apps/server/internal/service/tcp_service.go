package service

import (
	"fmt"
	"fullstack-engineering-lab/server/internal/model"
	tcpPkg "fullstack-engineering-lab/server/pkg/tcp"
	"time"
)

// TCPService TCP 业务逻辑
type TCPService struct {
	server       *tcpPkg.Server
	pool         *tcpPkg.SessionPool
	readTimeout  time.Duration
	writeTimeout time.Duration
}

// NewTCPService 创建 TCP 服务
func NewTCPService(server *tcpPkg.Server, pool *tcpPkg.SessionPool, readTimeoutSec, writeTimeoutSec int) *TCPService {
	return &TCPService{
		server:       server,
		pool:         pool,
		readTimeout:  time.Duration(readTimeoutSec) * time.Second,
		writeTimeout: time.Duration(writeTimeoutSec) * time.Second,
	}
}

// --- 基础协议 ---

// CreateSession 创建 TCP 会话（基础协议模式）
func (s *TCPService) CreateSession() (*model.TCPSessionResponse, error) {
	sess, err := tcpPkg.NewSession(s.server.Addr(), s.readTimeout, s.writeTimeout)
	if err != nil {
		return nil, fmt.Errorf("创建 TCP 会话失败: %w", err)
	}

	if err := s.pool.Add(sess); err != nil {
		sess.Close()
		return nil, err
	}

	return &model.TCPSessionResponse{
		SessionID: sess.ID,
		CreatedAt: sess.CreatedAt.Format(time.RFC3339),
	}, nil
}

// CloseSession 关闭会话
func (s *TCPService) CloseSession(id string) error {
	sess, ok := s.pool.Get(id)
	if !ok {
		return fmt.Errorf("会话 %s 不存在", id)
	}

	s.pool.Remove(id)
	sess.Close()
	return nil
}

// ListSessions 列出所有会话
func (s *TCPService) ListSessions() []*tcpPkg.SessionStats {
	return s.pool.List()
}

// SendCommand 发送命令并返回响应
func (s *TCPService) SendCommand(id string, cmd string) (*model.TCPCommandResponse, error) {
	sess, ok := s.pool.Get(id)
	if !ok {
		return nil, fmt.Errorf("会话 %s 不存在", id)
	}

	response, durationMs, err := sess.SendCommand(cmd, s.readTimeout, s.writeTimeout)
	if err != nil {
		return nil, fmt.Errorf("命令执行失败: %w", err)
	}

	return &model.TCPCommandResponse{
		SessionID:  id,
		Command:    cmd,
		Response:   response,
		DurationMs: durationMs,
		Timestamp:  time.Now().Format(time.RFC3339),
	}, nil
}

// StreamSession 获取会话事件流
func (s *TCPService) StreamSession(id string) (<-chan tcpPkg.SessionEvent, error) {
	sess, ok := s.pool.Get(id)
	if !ok {
		return nil, fmt.Errorf("会话 %s 不存在", id)
	}
	return sess.EventChannel(), nil
}

// GetStats 获取服务器统计
func (s *TCPService) GetStats() *model.TCPStatsResponse {
	return &model.TCPStatsResponse{
		ServerAddr:     s.server.Addr(),
		ActiveSessions: s.pool.Count(),
		MaxSessions:    100,
		UptimeSec:      int64(s.server.Uptime().Seconds()),
	}
}

// --- 聊天室 ---

// CreateChatSession 创建聊天会话（聊天模式）
func (s *TCPService) CreateChatSession(nickname, room string) (*model.TCPSessionResponse, error) {
	sess, err := tcpPkg.NewSession(s.server.Addr(), s.readTimeout, s.writeTimeout)
	if err != nil {
		return nil, fmt.Errorf("创建 TCP 会话失败: %w", err)
	}

	// 启动异步读取器
	sess.StartChatReader()

	// 设置昵称
	if err := sess.SendRaw("NICK " + nickname); err != nil {
		sess.Close()
		return nil, fmt.Errorf("设置昵称失败: %w", err)
	}

	// 加入聊天室
	if room != "" {
		if err := sess.SendRaw("JOIN " + room); err != nil {
			sess.Close()
			return nil, fmt.Errorf("加入聊天室失败: %w", err)
		}
	}

	if err := s.pool.Add(sess); err != nil {
		sess.Close()
		return nil, err
	}

	return &model.TCPSessionResponse{
		SessionID: sess.ID,
		CreatedAt: sess.CreatedAt.Format(time.RFC3339),
	}, nil
}

// SendChatMessage 发送聊天消息
func (s *TCPService) SendChatMessage(id string, content string) error {
	sess, ok := s.pool.Get(id)
	if !ok {
		return fmt.Errorf("会话 %s 不存在", id)
	}
	return sess.SendRaw("MSG " + content)
}

// ListChatRooms 列出所有聊天室
func (s *TCPService) ListChatRooms() []map[string]any {
	return s.server.RoomManager().ListRooms()
}

// GetChatRoomInfo 获取聊天室详情
func (s *TCPService) GetChatRoomInfo(name string) map[string]any {
	return s.server.RoomManager().GetRoomInfo(name)
}

// GetChatRoomMessages 获取聊天室最近消息
func (s *TCPService) GetChatRoomMessages(name string, limit int) []tcpPkg.ChatMessage {
	return s.server.RoomManager().GetRoomMessages(name, limit)
}
