package service

import (
	"errors"

	"github.com/miaomiaopu/ipmp-server/internal/model"
	"github.com/miaomiaopu/ipmp-server/internal/pkg/jwt"
	"github.com/miaomiaopu/ipmp-server/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserInactive       = errors.New("user is inactive")
)

type AuthService struct {
	userRepo   *repository.UserRepository
	jwtManager *jwt.Manager
}

func NewAuthService(userRepo *repository.UserRepository, jwtManager *jwt.Manager) *AuthService {
	return &AuthService{userRepo: userRepo, jwtManager: jwtManager}
}

// Login 用户登录，返回 access + refresh token
func (s *AuthService) Login(username, password string) (*model.User, string, string, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", "", ErrInvalidCredentials
		}
		return nil, "", "", err
	}
	if user.Status != "active" {
		return nil, "", "", ErrUserInactive
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", "", ErrInvalidCredentials
	}
	accessToken, err := s.jwtManager.GenerateAccess(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, "", "", err
	}
	refreshToken, err := s.jwtManager.GenerateRefresh(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, "", "", err
	}
	return user, accessToken, refreshToken, nil
}

// RefreshToken 刷新 access token
func (s *AuthService) RefreshToken(refreshTokenStr string) (*model.User, string, string, error) {
	claims, err := s.jwtManager.ParseRefresh(refreshTokenStr)
	if err != nil {
		return nil, "", "", ErrInvalidCredentials
	}
	user, err := s.userRepo.FindByID(claims.UserID)
	if err != nil {
		return nil, "", "", ErrInvalidCredentials
	}
	if user.Status != "active" {
		return nil, "", "", ErrUserInactive
	}
	accessToken, err := s.jwtManager.GenerateAccess(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, "", "", err
	}
	newRefreshToken, err := s.jwtManager.GenerateRefresh(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, "", "", err
	}
	return user, accessToken, newRefreshToken, nil
}

// GetCurrentUser 获取当前用户信息
func (s *AuthService) GetCurrentUser(userID string) (*model.User, error) {
	return s.userRepo.FindByID(userID)
}

// HashPassword bcrypt 哈希密码 (cost=12)
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(hash), err
}
