package service

import (
	"errors"
	"fmt"

	"fullstack-engineering-lab/server/internal/model"
	"fullstack-engineering-lab/server/internal/repository"
	ws "fullstack-engineering-lab/server/pkg/websocket"

	"gorm.io/gorm"
)

// ChatService 聊天业务逻辑
type ChatService struct {
	chatRepo repository.ChatRepository
	hub      *ws.Hub
}

// NewChatService 创建聊天服务实例
func NewChatService(chatRepo repository.ChatRepository, hub *ws.Hub) *ChatService {
	svc := &ChatService{
		chatRepo: chatRepo,
		hub:      hub,
	}

	// 注入 Hub 的回调函数
	hub.SetSaveMessageCallback(svc.SaveMessage)
	hub.SetGetUsernameCallback(svc.GetUsername)
	hub.SetGetHistoryCallback(svc.GetRoomHistory)

	return svc
}

// GetRooms 获取所有聊天室列表（含在线人数）
func (s *ChatService) GetRooms() ([]model.RoomResponse, error) {
	rooms, err := s.chatRepo.FindAllRooms()
	if err != nil {
		return nil, err
	}

	result := make([]model.RoomResponse, 0, len(rooms))
	for _, room := range rooms {
		result = append(result, model.RoomResponse{
			ID:          room.ID,
			Name:        room.Name,
			Description: room.Description,
			Type:        room.Type,
			CreatorID:   room.CreatorID,
			MemberCount: s.hub.GetOnlineUserCount(room.ID),
			Status:      room.Status,
			CreatedAt:   room.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return result, nil
}

// CreateRoom 创建聊天室
func (s *ChatService) CreateRoom(req *model.CreateRoomRequest, userID uint) (*model.RoomResponse, error) {
	room := &model.ChatRoom{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		CreatorID:   userID,
		MaxMembers:  500,
		Status:      1,
	}

	if err := s.chatRepo.CreateRoom(room); err != nil {
		return nil, fmt.Errorf("创建聊天室失败: %w", err)
	}

	return &model.RoomResponse{
		ID:          room.ID,
		Name:        room.Name,
		Description: room.Description,
		Type:        room.Type,
		CreatorID:   room.CreatorID,
		MemberCount: 0,
		Status:      room.Status,
		CreatedAt:   room.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// GetRoomInfo 获取聊天室详情
func (s *ChatService) GetRoomInfo(roomID uint) (*model.RoomResponse, error) {
	room, err := s.chatRepo.FindRoomByID(roomID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("聊天室不存在")
		}
		return nil, err
	}

	return &model.RoomResponse{
		ID:          room.ID,
		Name:        room.Name,
		Description: room.Description,
		Type:        room.Type,
		CreatorID:   room.CreatorID,
		MemberCount: s.hub.GetOnlineUserCount(room.ID),
		Status:      room.Status,
		CreatedAt:   room.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

// GetMessageHistory 获取消息历史
func (s *ChatService) GetMessageHistory(req *model.MessageHistoryRequest) (*model.MessageHistoryResponse, error) {
	limit := req.Limit
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	// 验证房间是否存在
	_, err := s.chatRepo.FindRoomByID(req.RoomID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("聊天室不存在")
		}
		return nil, err
	}

	messages, hasMore, err := s.chatRepo.FindMessagesByRoomID(req.RoomID, limit, req.Before)
	if err != nil {
		return nil, err
	}

	// 查询用户名并组装响应
	result := make([]model.MessageResponse, 0, len(messages))
	for _, msg := range messages {
		username, _ := s.chatRepo.FindUsernameByID(msg.UserID)
		if username == "" {
			username = fmt.Sprintf("用户%d", msg.UserID)
		}
		result = append(result, model.MessageResponse{
			ID:        msg.ID,
			RoomID:    msg.RoomID,
			UserID:    msg.UserID,
			Username:  username,
			Content:   msg.Content,
			MsgType:   msg.MsgType,
			CreatedAt: msg.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &model.MessageHistoryResponse{
		Messages: result,
		HasMore:  hasMore,
	}, nil
}

// GetOnlineUsers 获取房间在线用户
func (s *ChatService) GetOnlineUsers(roomID uint) []ws.OnlineUserSimple {
	return s.hub.GetOnlineUsers(roomID)
}

// GetHub 获取 Hub 实例（供 Handler 使用）
func (s *ChatService) GetHub() *ws.Hub {
	return s.hub
}

// SaveMessage 持久化消息（Hub 回调）
func (s *ChatService) SaveMessage(roomID, userID uint, content string, msgType int8) (uint, error) {
	msg := &model.ChatMessage{
		RoomID:  roomID,
		UserID:  userID,
		Content: content,
		MsgType: msgType,
	}
	if err := s.chatRepo.CreateMessage(msg); err != nil {
		return 0, err
	}
	return msg.ID, nil
}

// GetUsername 查询用户名（Hub 回调）
func (s *ChatService) GetUsername(userID uint) string {
	username, _ := s.chatRepo.FindUsernameByID(userID)
	if username == "" {
		return fmt.Sprintf("用户%d", userID)
	}
	return username
}

// GetRoomHistory 获取房间历史消息（Hub 回调）
func (s *ChatService) GetRoomHistory(roomID uint, limit int) ([]ws.HistoryMessage, error) {
	messages, _, err := s.chatRepo.FindMessagesByRoomID(roomID, limit, 0)
	if err != nil {
		return nil, err
	}

	result := make([]ws.HistoryMessage, 0, len(messages))
	for _, msg := range messages {
		username, _ := s.chatRepo.FindUsernameByID(msg.UserID)
		if username == "" {
			username = fmt.Sprintf("用户%d", msg.UserID)
		}
		result = append(result, ws.HistoryMessage{
			ID:        msg.ID,
			UserID:    msg.UserID,
			Username:  username,
			Content:   msg.Content,
			MsgType:   msg.MsgType,
			CreatedAt: msg.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return result, nil
}
