package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username     string `gorm:"type:varchar(64);uniqueIndex;not null"`
	Email        string `gorm:"type:varchar(128);uniqueIndex;not null"`
	PasswordHash string `gorm:"type:varchar(255);not null"`
	Nickname     string `gorm:"type:varchar(64);default:''"`
	Status       int8   `gorm:"type:tinyint;not null;default:1"`
}

func (User) TableName() string {
	return "users"
}
