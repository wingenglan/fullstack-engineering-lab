package tcp

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// ChatRoom 聊天室
type ChatRoom struct {
	Name       string
	members    map[string]*ChatMember // nickname → member
	messages   []ChatMessage          // 消息历史环形缓冲区
	maxHistory int
	mu         sync.RWMutex
	createdAt  time.Time
}

// ChatMember 聊天室成员（与 TCP 连接的映射）
type ChatMember struct {
	Nickname  string
	Room      string
	Conn      net.Conn
	JoinedAt  time.Time
	MsgCount  int64
}

// ChatMessage 聊天消息
type ChatMessage struct {
	Room     string `json:"room"`
	Nickname string `json:"nickname"`
	Content  string `json:"content"`
	Type     string `json:"type"` // message / system_join / system_leave
	Time     string `json:"time"`
}

// RoomManager 聊天室管理器
type RoomManager struct {
	rooms      map[string]*ChatRoom
	connections map[net.Conn]*ChatMember // 所有活跃连接的上下文
	mu         sync.RWMutex
}

// NewRoomManager 创建聊天室管理器
func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms:       make(map[string]*ChatRoom),
		connections: make(map[net.Conn]*ChatMember),
	}
}

// GetOrCreateRoom 获取或创建聊天室
func (m *RoomManager) GetOrCreateRoom(name string) *ChatRoom {
	m.mu.Lock()
	defer m.mu.Unlock()

	if room, ok := m.rooms[name]; ok {
		return room
	}

	room := &ChatRoom{
		Name:       name,
		members:    make(map[string]*ChatMember),
		messages:   make([]ChatMessage, 0, 200),
		maxHistory: 200,
		createdAt:  time.Now(),
	}
	m.rooms[name] = room
	return room
}

// JoinRoom 成员加入聊天室
func (m *RoomManager) JoinRoom(conn net.Conn, nickname, roomName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 先退出旧房间
	if member, ok := m.connections[conn]; ok && member.Room != "" {
		m.leaveRoomLocked(conn, member)
	}

	room := m.getOrCreateRoomLocked(roomName)

	// 检查昵称在该房间是否已存在
	if _, exists := room.members[nickname]; exists {
		return fmt.Errorf("昵称 %s 在该聊天室已被占用", nickname)
	}

	member := &ChatMember{
		Nickname: nickname,
		Room:     roomName,
		Conn:     conn,
		JoinedAt: time.Now(),
	}
	m.connections[conn] = member
	room.members[nickname] = member

	// 广播加入消息
	msg := ChatMessage{
		Room:     roomName,
		Nickname: "System",
		Content:  fmt.Sprintf("%s 加入了聊天室", nickname),
		Type:     "system_join",
		Time:     time.Now().Format(time.RFC3339),
	}
	m.broadcastToRoomLocked(room, msg, "")

	return nil
}

// LeaveRoom 成员离开当前聊天室
func (m *RoomManager) LeaveRoom(conn net.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()

	member, ok := m.connections[conn]
	if !ok || member.Room == "" {
		return
	}

	m.leaveRoomLocked(conn, member)
}

// SendMessage 发送消息到成员所在的聊天室
func (m *RoomManager) SendMessage(conn net.Conn, content string) {
	m.mu.RLock()
	member, ok := m.connections[conn]
	if !ok || member.Room == "" {
		m.mu.RUnlock()
		conn.Write([]byte("[Error] 请先使用 JOIN <聊天室名> 加入聊天室\n"))
		return
	}
	roomName := member.Room
	nickname := member.Nickname
	room, hasRoom := m.rooms[roomName]
	m.mu.RUnlock()

	if !hasRoom {
		return
	}

	msg := ChatMessage{
		Room:     roomName,
		Nickname: nickname,
		Content:  content,
		Type:     "message",
		Time:     time.Now().Format(time.RFC3339),
	}

	room.mu.Lock()
	room.appendMessageLocked(msg)
	room.mu.Unlock()

	m.mu.RLock()
	m.broadcastToRoomLocked(room, msg, "")
	m.mu.RUnlock()
}

// GetMember 获取连接的成员信息
func (m *RoomManager) GetMember(conn net.Conn) *ChatMember {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.connections[conn]
}

// SetNickname 为连接设置昵称（加入聊天室之前）
func (m *RoomManager) SetNickname(conn net.Conn, nickname string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查昵称是否已被其他连接使用
	for _, member := range m.connections {
		if member.Nickname == nickname {
			return fmt.Errorf("昵称 %s 已被使用", nickname)
		}
	}

	if existing, ok := m.connections[conn]; ok {
		existing.Nickname = nickname
	} else {
		m.connections[conn] = &ChatMember{
			Nickname: nickname,
			Conn:     conn,
			JoinedAt: time.Now(),
		}
	}
	return nil
}

// RemoveConnection 移除连接（断开时调用）
func (m *RoomManager) RemoveConnection(conn net.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()

	member, ok := m.connections[conn]
	if !ok {
		return
	}

	if member.Room != "" {
		m.leaveRoomLocked(conn, member)
	}

	delete(m.connections, conn)
}

// GetRoomInfo 获取聊天室信息
func (m *RoomManager) GetRoomInfo(name string) map[string]any {
	m.mu.RLock()
	defer m.mu.RUnlock()

	room, ok := m.rooms[name]
	if !ok {
		return nil
	}

	room.mu.RLock()
	defer room.mu.RUnlock()

	users := make([]string, 0, len(room.members))
	for nick := range room.members {
		users = append(users, nick)
	}

	return map[string]any{
		"name":       name,
		"user_count": len(room.members),
		"users":      users,
		"created_at": room.createdAt.Format(time.RFC3339),
	}
}

// ListRooms 列出所有聊天室
func (m *RoomManager) ListRooms() []map[string]any {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]map[string]any, 0, len(m.rooms))
	for name, room := range m.rooms {
		room.mu.RLock()
		result = append(result, map[string]any{
			"name":       name,
			"user_count": len(room.members),
			"created_at": room.createdAt.Format(time.RFC3339),
		})
		room.mu.RUnlock()
	}
	return result
}

// GetRoomMessages 获取聊天室最近消息
func (m *RoomManager) GetRoomMessages(name string, limit int) []ChatMessage {
	m.mu.RLock()
	room, ok := m.rooms[name]
	m.mu.RUnlock()

	if !ok {
		return []ChatMessage{}
	}

	room.mu.RLock()
	defer room.mu.RUnlock()

	if limit <= 0 || limit > len(room.messages) {
		result := make([]ChatMessage, len(room.messages))
		copy(result, room.messages)
		return result
	}
	result := make([]ChatMessage, limit)
	copy(result, room.messages[len(room.messages)-limit:])
	return result
}

// --- 内部方法（调用者需持有锁） ---

func (m *RoomManager) getOrCreateRoomLocked(name string) *ChatRoom {
	if room, ok := m.rooms[name]; ok {
		return room
	}
	room := &ChatRoom{
		Name:       name,
		members:    make(map[string]*ChatMember),
		messages:   make([]ChatMessage, 0, 200),
		maxHistory: 200,
		createdAt:  time.Now(),
	}
	m.rooms[name] = room
	return room
}

func (m *RoomManager) leaveRoomLocked(conn net.Conn, member *ChatMember) {
	room, ok := m.rooms[member.Room]
	if !ok {
		return
	}

	delete(room.members, member.Nickname)

	msg := ChatMessage{
		Room:     member.Room,
		Nickname: "System",
		Content:  fmt.Sprintf("%s 离开了聊天室", member.Nickname),
		Type:     "system_leave",
		Time:     time.Now().Format(time.RFC3339),
	}
	m.broadcastToRoomLocked(room, msg, "")

	// 清理空房间
	if len(room.members) == 0 {
		delete(m.rooms, member.Room)
	}

	member.Room = ""
}

func (m *RoomManager) broadcastToRoomLocked(room *ChatRoom, msg ChatMessage, excludeNickname string) {
	formatted := formatChatMessage(msg)
	for nickname, member := range room.members {
		if nickname == excludeNickname {
			continue
		}
		member.Conn.Write([]byte(formatted + "\n"))
	}
}

func (room *ChatRoom) appendMessageLocked(msg ChatMessage) {
	if len(room.messages) >= room.maxHistory {
		room.messages = room.messages[1:]
	}
	room.messages = append(room.messages, msg)
}

func formatChatMessage(msg ChatMessage) string {
	switch msg.Type {
	case "system_join", "system_leave":
		return fmt.Sprintf("[%s] %s", msg.Nickname, msg.Content)
	default:
		return fmt.Sprintf("[%s] %s: %s", msg.Room, msg.Nickname, msg.Content)
	}
}
