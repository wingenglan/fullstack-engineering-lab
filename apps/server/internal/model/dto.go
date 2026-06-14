package model

import "time"

type RegisterRequest struct {
	Username string `json:"username" binding:"required" validate:"min=3,max=64"`
	Email    string `json:"email" binding:"required,email" validate:"email"`
	Password string `json:"password" binding:"required" validate:"min=6,max=128"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Nickname  string    `json:"nickname"`
	Status    int8      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func ToUserResponse(u *User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Nickname:  u.Nickname,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
	}
}
