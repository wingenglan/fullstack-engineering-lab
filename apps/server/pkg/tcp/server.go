package tcp

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
)

// Server TCP 服务器（集成聊天室功能）
type Server struct {
	listener   net.Listener
	cfg        *ServerConfig
	startedAt  time.Time
	stopCh     chan struct{}
	connCount  int64
	roomMgr    *RoomManager
	mu         sync.Mutex
}

// ServerConfig TCP 服务器配置
type ServerConfig struct {
	Host            string
	Port            int
	ReadTimeoutSec  int
	WriteTimeoutSec int
}

// NewServer 创建 TCP 服务器
func NewServer(cfg *ServerConfig) *Server {
	return &Server{
		cfg:       cfg,
		stopCh:    make(chan struct{}),
		startedAt: time.Now(),
		roomMgr:   NewRoomManager(),
	}
}

// Start 启动 TCP 服务器
func (s *Server) Start(ctx context.Context) error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("监听 %s 失败: %w", addr, err)
	}
	s.listener = ln
	zap.L().Info("TCP 服务器已启动", zap.String("addr", addr))

	go s.acceptLoop()
	return nil
}

// Stop 停止 TCP 服务器
func (s *Server) Stop() {
	close(s.stopCh)
	if s.listener != nil {
		s.listener.Close()
	}
	zap.L().Info("TCP 服务器已停止")
}

// Addr 返回监听地址
func (s *Server) Addr() string {
	if s.listener == nil {
		return ""
	}
	return s.listener.Addr().String()
}

// Uptime 运行时长
func (s *Server) Uptime() time.Duration {
	return time.Since(s.startedAt)
}

// RoomManager 返回聊天室管理器
func (s *Server) RoomManager() *RoomManager {
	return s.roomMgr
}

func (s *Server) acceptLoop() {
	for {
		select {
		case <-s.stopCh:
			return
		default:
		}

		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.stopCh:
				return
			default:
				zap.L().Error("TCP Accept 错误", zap.Error(err))
				continue
			}
		}

		s.mu.Lock()
		s.connCount++
		s.mu.Unlock()

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer func() {
		s.roomMgr.RemoveConnection(conn)
		conn.Close()
		s.mu.Lock()
		s.connCount--
		s.mu.Unlock()
	}()

	remoteAddr := conn.RemoteAddr().String()
	zap.L().Debug("TCP 新连接", zap.String("addr", remoteAddr))

	// 欢迎消息
	conn.Write([]byte("欢迎来到 TCP 聊天服务器！输入 HELP 查看基础命令，输入 NICK <昵称> 进入聊天模式\n"))

	buf := make([]byte, 4096)
	for {
		select {
		case <-s.stopCh:
			return
		default:
		}

		if s.cfg.ReadTimeoutSec > 0 {
			conn.SetReadDeadline(time.Now().Add(time.Duration(s.cfg.ReadTimeoutSec) * time.Second))
		}

		n, err := conn.Read(buf)
		if err != nil {
			return
		}

		line := trimNewlines(string(buf[:n]))
		if line == "" {
			continue
		}

		cmd, args := ParseCommand(line)
		reply := s.handleCommand(conn, cmd, args)

		if s.cfg.WriteTimeoutSec > 0 {
			conn.SetWriteDeadline(time.Now().Add(time.Duration(s.cfg.WriteTimeoutSec) * time.Second))
		}

		if reply != "" {
			conn.Write([]byte(reply + "\n"))
		}
	}
}

// handleCommand 处理命令：先尝试聊天命令，再回退到基础命令
func (s *Server) handleCommand(conn net.Conn, cmd, args string) string {
	switch cmd {
	case "NICK":
		if args == "" {
			return "用法: NICK <昵称>"
		}
		if err := s.roomMgr.SetNickname(conn, args); err != nil {
			return fmt.Sprintf("错误: %s", err.Error())
		}
		return fmt.Sprintf("昵称设置为: %s。输入 JOIN <聊天室名> 加入聊天室，或 HELP 查看更多命令", args)

	case "JOIN":
		if args == "" {
			return "用法: JOIN <聊天室名>"
		}
		member := s.roomMgr.GetMember(conn)
		if member == nil || member.Nickname == "" {
			return "请先使用 NICK <昵称> 设置昵称"
		}
		if err := s.roomMgr.JoinRoom(conn, member.Nickname, args); err != nil {
			return fmt.Sprintf("加入失败: %s", err.Error())
		}
		return fmt.Sprintf("你已加入聊天室: %s", args)

	case "LEAVE":
		s.roomMgr.LeaveRoom(conn)
		return "你已离开聊天室"

	case "MSG":
		if args == "" {
			return "用法: MSG <消息内容>"
		}
		member := s.roomMgr.GetMember(conn)
		if member == nil || member.Room == "" {
			return "请先使用 JOIN <聊天室名> 加入聊天室"
		}
		s.roomMgr.SendMessage(conn, args)
		return "" // 消息已广播到房间，无需回显

	case "/USERS":
		member := s.roomMgr.GetMember(conn)
		if member == nil || member.Room == "" {
			return "你不在任何聊天室中"
		}
		info := s.roomMgr.GetRoomInfo(member.Room)
		if info == nil {
			return "聊天室不存在"
		}
		users := info["users"].([]string)
		return fmt.Sprintf("聊天室 [%s] 在线用户 (%d): %s", member.Room, len(users), strings.Join(users, ", "))

	case "/ROOMS":
		rooms := s.roomMgr.ListRooms()
		if len(rooms) == 0 {
			return "暂无活跃聊天室"
		}
		lines := make([]string, 0, len(rooms))
		for _, r := range rooms {
			lines = append(lines, fmt.Sprintf("  %s (%d人)", r["name"], r["user_count"]))
		}
		return "聊天室列表:\n" + strings.Join(lines, "\n")

	case "QUIT":
		return "BYE"

	default:
		// 回退到基础协议命令
		result := Handle(cmd, args)
		return result.Response
	}
}

func trimNewlines(s string) string {
	s = strings.TrimSuffix(s, "\r\n")
	s = strings.TrimSuffix(s, "\n")
	s = strings.TrimSuffix(s, "\r")
	return s
}
