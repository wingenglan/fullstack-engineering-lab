package tcp

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

// SessionEvent 会话事件
type SessionEvent struct {
	Type       string `json:"type"` // connected / message / system / disconnected / command_result
	SessionID  string `json:"session_id"`
	Nickname   string `json:"nickname,omitempty"`
	Room       string `json:"room,omitempty"`
	Content    string `json:"content,omitempty"`
	Command    string `json:"command,omitempty"`
	Response   string `json:"response,omitempty"`
	DurationMs int64  `json:"duration_ms,omitempty"`
	Timestamp  string `json:"timestamp"`
}

// Session TCP 会话（支持基础协议和聊天模式）
type Session struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	conn      net.Conn
	lastActAt time.Time
	mu        sync.Mutex // 保护 conn 写入
	cmdCount  int64
	eventCh   chan SessionEvent
	closed    atomic.Bool
	readerStarted atomic.Bool
}

// NewSession 创建会话（建立 TCP 连接）
func NewSession(serverAddr string, readTimeout, writeTimeout time.Duration) (*Session, error) {
	conn, err := net.DialTimeout("tcp", serverAddr, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("连接 TCP 服务器失败: %w", err)
	}

	if writeTimeout > 0 {
		conn.SetWriteDeadline(time.Now().Add(writeTimeout))
	}

	sess := &Session{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		conn:      conn,
		lastActAt: time.Now(),
		eventCh:   make(chan SessionEvent, 256),
	}

	sess.emitEvent(SessionEvent{
		Type:      "connected",
		SessionID: sess.ID,
		Timestamp: time.Now().Format(time.RFC3339),
	})

	return sess, nil
}

// StartChatReader 启动聊天模式的异步读取（用于聊天室场景）
func (s *Session) StartChatReader() {
	if !s.readerStarted.CompareAndSwap(false, true) {
		return
	}

	go func() {
		reader := bufio.NewReader(s.conn)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				s.emitEvent(SessionEvent{
					Type:      "disconnected",
					SessionID: s.ID,
					Timestamp: time.Now().Format(time.RFC3339),
				})
				s.Close()
				return
			}

			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}

			s.mu.Lock()
			s.lastActAt = time.Now()
			s.mu.Unlock()
			atomic.AddInt64(&s.cmdCount, 1)

			evt := parseChatLine(s.ID, line)
			s.emitEvent(evt)
		}
	}()
}

// SendRaw 发送原始命令（不等待响应，用于聊天模式）
func (s *Session) SendRaw(cmd string) error {
	if s.closed.Load() {
		return fmt.Errorf("会话已关闭")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.conn.Write([]byte(cmd + "\n"))
	return err
}

// SendCommand 发送命令并同步读取响应（用于基础协议模式）
func (s *Session) SendCommand(cmd string, readTimeout, writeTimeout time.Duration) (response string, durationMs int64, err error) {
	if s.closed.Load() {
		return "", 0, fmt.Errorf("会话已关闭")
	}

	start := time.Now()

	s.mu.Lock()
	defer s.mu.Unlock()

	if writeTimeout > 0 {
		s.conn.SetWriteDeadline(time.Now().Add(writeTimeout))
	}
	_, err = s.conn.Write([]byte(cmd + "\n"))
	if err != nil {
		return "", 0, fmt.Errorf("发送命令失败: %w", err)
	}

	if readTimeout > 0 {
		s.conn.SetReadDeadline(time.Now().Add(readTimeout))
	}

	// 使用 bufio 读取一行
	reader := bufio.NewReader(s.conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", 0, fmt.Errorf("读取响应失败: %w", err)
	}

	response = strings.TrimSpace(line)
	durationMs = time.Since(start).Milliseconds()

	s.lastActAt = time.Now()
	atomic.AddInt64(&s.cmdCount, 1)

	s.emitEvent(SessionEvent{
		Type:       "command_result",
		SessionID:  s.ID,
		Command:    cmd,
		Response:   response,
		DurationMs: durationMs,
		Timestamp:  time.Now().Format(time.RFC3339),
	})

	return response, durationMs, nil
}

// Close 关闭会话
func (s *Session) Close() {
	if s.closed.Swap(true) {
		return
	}

	s.mu.Lock()
	s.conn.Write([]byte("QUIT\n"))
	s.conn.Close()
	s.mu.Unlock()

	s.emitEvent(SessionEvent{
		Type:      "disconnected",
		SessionID: s.ID,
		Timestamp: time.Now().Format(time.RFC3339),
	})

	close(s.eventCh)
}

// EventChannel 返回事件通道（只读）
func (s *Session) EventChannel() <-chan SessionEvent {
	return s.eventCh
}

// IsAlive 检查会话是否存活
func (s *Session) IsAlive() bool {
	return !s.closed.Load()
}

// Stats 返回会话统计
func (s *Session) Stats() *SessionStats {
	s.mu.Lock()
	lastAct := s.lastActAt
	s.mu.Unlock()

	return &SessionStats{
		SessionID:    s.ID,
		RemoteAddr:   s.conn.RemoteAddr().String(),
		CreatedAt:    s.CreatedAt.Format(time.RFC3339),
		LastActAt:    lastAct.Format(time.RFC3339),
		DurationSec:  int(time.Since(s.CreatedAt).Seconds()),
		CommandCount: atomic.LoadInt64(&s.cmdCount),
		IsAlive:      !s.closed.Load(),
	}
}

// SessionStats 会话统计信息
type SessionStats struct {
	SessionID    string `json:"session_id"`
	RemoteAddr   string `json:"remote_addr"`
	CreatedAt    string `json:"created_at"`
	LastActAt    string `json:"last_act_at"`
	DurationSec  int    `json:"duration_sec"`
	CommandCount int64  `json:"command_count"`
	IsAlive      bool   `json:"is_alive"`
}

func (s *Session) emitEvent(evt SessionEvent) {
	select {
	case s.eventCh <- evt:
	default:
	}
}

// parseChatLine 解析聊天消息行
func parseChatLine(sessionID, line string) SessionEvent {
	ts := time.Now().Format(time.RFC3339)

	// [System] xxx 加入了聊天室
	// [System] xxx 离开了聊天室
	if strings.HasPrefix(line, "[System]") {
		return SessionEvent{
			Type:      "system",
			SessionID: sessionID,
			Content:   strings.TrimPrefix(line, "[System] "),
			Timestamp: ts,
		}
	}

	// [RoomName] Nickname: message
	if strings.HasPrefix(line, "[") {
		// 提取房间名
		closeBracket := strings.Index(line, "]")
		if closeBracket > 1 {
			room := line[1:closeBracket]
			rest := strings.TrimSpace(line[closeBracket+1:])
			// 分离昵称和消息
			colonIdx := strings.Index(rest, ": ")
			if colonIdx > 0 {
				nickname := rest[:colonIdx]
				content := rest[colonIdx+2:]
				return SessionEvent{
					Type: "message",
					SessionID: sessionID,
					Room:      room,
					Nickname:  nickname,
					Content:   content,
					Timestamp: ts,
				}
			}
		}
	}

	// 默认当作普通消息
	return SessionEvent{
		Type:      "message",
		SessionID: sessionID,
		Content:   line,
		Timestamp: ts,
	}
}
