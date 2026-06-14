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

// ===== Redis 分布式锁 DTO =====

type LockAcquireRequest struct {
	Resource string `json:"resource" binding:"required"`
	TTL      int    `json:"ttl" binding:"required,min=1,max=300"`
}

type LockReleaseRequest struct {
	Resource string `json:"resource" binding:"required"`
}

type LockStatusRequest struct {
	Resource string `json:"resource" binding:"required"`
}

type ContentionDemoRequest struct {
	Resource   string `json:"resource" binding:"required"`
	TTL        int    `json:"ttl" binding:"required,min=1,max=300"`
	Goroutines int    `json:"goroutines" binding:"required,min=2,max=20"`
	HoldMs     int    `json:"hold_ms" binding:"required,min=100,max=10000"`
}

type LockStatusResponse struct {
	Resource string `json:"resource"`
	Locked   bool   `json:"locked"`
	Owner    string `json:"owner,omitempty"`
	TTL      int64  `json:"ttl_ms"`
}

type LockAcquireResponse struct {
	Resource string `json:"resource"`
	Owner    string `json:"owner"`
	TTL      int    `json:"ttl"`
}

type ContentionResult struct {
	GoroutineID int    `json:"goroutine_id"`
	Acquired    bool   `json:"acquired"`
	WaitMs      int64  `json:"wait_ms"`
	Message     string `json:"message"`
}

type ContentionDemoResponse struct {
	Resource  string              `json:"resource"`
	Results   []ContentionResult  `json:"results"`
	Summary   ContentionSummary   `json:"summary"`
}

type ContentionSummary struct {
	Total      int `json:"total"`
	Succeeded  int `json:"succeeded"`
	Failed     int `json:"failed"`
}
