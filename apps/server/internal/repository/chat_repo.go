package repository

import (
	"fullstack-engineering-lab/server/internal/model"

	"gorm.io/gorm"
)

// ChatRepository 聊天数据层接口
type ChatRepository interface {
	// 房间操作
	CreateRoom(room *model.ChatRoom) error
	FindRoomByID(id uint) (*model.ChatRoom, error)
	FindAllRooms() ([]model.ChatRoom, error)
	FindRoomsByUserID(userID uint) ([]model.ChatRoom, error)

	// 消息操作
	CreateMessage(msg *model.ChatMessage) error
	FindMessagesByRoomID(roomID uint, limit int, beforeID uint) ([]model.ChatMessage, bool, error)

	// 用户信息
	FindUsernameByID(userID uint) (string, error)
}

type chatRepository struct {
	db *gorm.DB
}

// NewChatRepository 创建聊天数据层实例
func NewChatRepository(db *gorm.DB) ChatRepository {
	return &chatRepository{db: db}
}

// CreateRoom 创建聊天室
func (r *chatRepository) CreateRoom(room *model.ChatRoom) error {
	return r.db.Create(room).Error
}

// FindRoomByID 根据 ID 查找聊天室
func (r *chatRepository) FindRoomByID(id uint) (*model.ChatRoom, error) {
	var room model.ChatRoom
	err := r.db.First(&room, id).Error
	if err != nil {
		return nil, err
	}
	return &room, nil
}

// FindAllRooms 查找所有活跃聊天室
func (r *chatRepository) FindAllRooms() ([]model.ChatRoom, error) {
	var rooms []model.ChatRoom
	err := r.db.Where("status = ?", 1).Order("id ASC").Find(&rooms).Error
	return rooms, err
}

// FindRoomsByUserID 查找用户创建的房间
func (r *chatRepository) FindRoomsByUserID(userID uint) ([]model.ChatRoom, error) {
	var rooms []model.ChatRoom
	err := r.db.Where("creator_id = ? AND status = ?", userID, 1).Find(&rooms).Error
	return rooms, err
}

// CreateMessage 创建消息
func (r *chatRepository) CreateMessage(msg *model.ChatMessage) error {
	return r.db.Create(msg).Error
}

// FindMessagesByRoomID 获取房间历史消息（游标分页）
func (r *chatRepository) FindMessagesByRoomID(roomID uint, limit int, beforeID uint) ([]model.ChatMessage, bool, error) {
	var messages []model.ChatMessage
	query := r.db.Where("room_id = ?", roomID)

	if beforeID > 0 {
		query = query.Where("id < ?", beforeID)
	}

	// 多查一条以判断是否还有更多
	err := query.Order("id DESC").Limit(limit + 1).Find(&messages).Error
	if err != nil {
		return nil, false, err
	}

	hasMore := len(messages) > limit
	if hasMore {
		messages = messages[:limit]
	}

	// 反转为时间正序
	for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
		messages[i], messages[j] = messages[j], messages[i]
	}

	return messages, hasMore, nil
}

// FindUsernameByID 根据用户 ID 查找用户名
func (r *chatRepository) FindUsernameByID(userID uint) (string, error) {
	var user model.User
	err := r.db.Select("username").First(&user, userID).Error
	if err != nil {
		return "", err
	}
	return user.Username, nil
}
