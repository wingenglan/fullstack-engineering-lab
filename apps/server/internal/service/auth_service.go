package service

import (
	"context"
	"errors"
	"time"

	"fullstack-engineering-lab/server/internal/config"
	"fullstack-engineering-lab/server/internal/model"
	"fullstack-engineering-lab/server/internal/repository"
	"fullstack-engineering-lab/server/pkg/jwt"
	"fullstack-engineering-lab/server/pkg/password"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo repository.UserRepository
	rdb      *redis.Client
	jwtCfg   *config.JWTConfig
}

func NewAuthService(userRepo repository.UserRepository, rdb *redis.Client, jwtCfg *config.JWTConfig) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		rdb:      rdb,
		jwtCfg:   jwtCfg,
	}
}

func (s *AuthService) Register(req *model.RegisterRequest) error {
	// Check duplicate username
	if _, err := s.userRepo.FindByUsername(req.Username); err == nil {
		return errors.New("username already exists")
	}

	// Check duplicate email
	if _, err := s.userRepo.FindByEmail(req.Email); err == nil {
		return errors.New("email already exists")
	}

	// Hash password
	hash, err := password.Hash(req.Password)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hash,
		Status:       1,
	}

	return s.userRepo.Create(user)
}

func (s *AuthService) Login(req *model.LoginRequest) (*model.LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid username or password")
		}
		return nil, err
	}

	if !password.Check(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid username or password")
	}

	token, err := jwt.GenerateToken(user.ID, s.jwtCfg)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &model.LoginResponse{
		AccessToken: token,
		ExpiresIn:   s.jwtCfg.ExpireMinutes * 60,
	}, nil
}

func (s *AuthService) GetProfile(userID uint) (*model.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	resp := model.ToUserResponse(user)
	return &resp, nil
}

func (s *AuthService) Logout(token string) error {
	if s.rdb == nil {
		return nil // Graceful degradation
	}

	// Store token in blacklist with TTL
	ttl := time.Duration(s.jwtCfg.ExpireMinutes) * time.Minute
	return s.rdb.Set(context.Background(), "blacklist:"+token, "1", ttl).Err()
}
