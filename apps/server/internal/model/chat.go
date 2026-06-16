package model

import "gorm.io/gorm"

// ChatRoom 聊天室模型
type ChatRoom struct {
	gorm.Model
	Name        string `gorm:"type:varchar(128);not null"`
	Description string `gorm:"type:varchar(512);default:''"`
	Type        int8   `gorm:"type:tinyint;not null;default:1"` // 1=群聊, 2=私聊
	CreatorID   uint   `gorm:"not null"`
	MaxMembers  int    `gorm:"not null;default:500"`
	Status      int8   `gorm:"type:tinyint;not null;default:1"` // 0=关闭, 1=活跃
}

func (ChatRoom) TableName() string {
	return "chat_rooms"
}

// ChatMessage 聊天消息模型
type ChatMessage struct {
	gorm.Model
	RoomID  uint   `gorm:"not null;index"`
	UserID  uint   `gorm:"not null;index"`
	Content string `gorm:"type:text;not null"`
	MsgType int8   `gorm:"type:tinyint;not null;default:1"` // 1=文本, 2=图片, 3=系统
}

func (ChatMessage) TableName() string {
	return "chat_messages"
}
