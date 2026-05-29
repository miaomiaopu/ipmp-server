package service

import (
	"errors"
	"time"

	"github.com/miaomiaopu/ipmp-server/internal/model"
	"github.com/miaomiaopu/ipmp-server/internal/pkg/utils"
	"github.com/miaomiaopu/ipmp-server/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrUserExists     = errors.New("username already exists")
	ErrWrongPassword  = errors.New("wrong password")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) List(page, pageSize int, keyword string) ([]model.User, int64, error) {
	return s.repo.List(page, pageSize, keyword)
}

func (s *UserService) GetByID(id string) (*model.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

// Create 创建用户，密码为空时自动生成随机密码
func (s *UserService) Create(req *model.User, password string) (*model.User, string, error) {
	existing, _ := s.repo.FindByUsername(req.Username)
	if existing != nil {
		return nil, "", ErrUserExists
	}
	if password == "" {
		password = utils.GenPassword(12)
	}
	hash, err := HashPassword(password)
	if err != nil {
		return nil, "", err
	}
	req.PasswordHash = hash
	if err := s.repo.Create(req); err != nil {
		return nil, "", err
	}
	return req, password, nil
}

func (s *UserService) Update(id string, updates map[string]interface{}) error {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}
	for key, value := range updates {
		switch key {
		case "display_name":
			user.DisplayName = value.(string)
		case "email":
			user.Email = model.EncryptedField(value.(string))
		case "phone":
			user.Phone = model.EncryptedField(value.(string))
		case "role":
			user.Role = value.(string)
		case "status":
			user.Status = value.(string)
		}
	}
	user.UpdatedAt = time.Now()
	return s.repo.Update(user)
}

func (s *UserService) Delete(id string) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}
	return s.repo.SoftDelete(id)
}

// ChangePassword 用户修改自己的密码
func (s *UserService) ChangePassword(userID, oldPassword, newPassword string) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return ErrUserNotFound
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return ErrWrongPassword
	}
	hash, err := HashPassword(newPassword)
	if err != nil {
		return err
	}
	user.PasswordHash = hash
	return s.repo.Update(user)
}

// ResetPassword 管理员重置用户密码
func (s *UserService) ResetPassword(userID, newPassword string) (string, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return "", ErrUserNotFound
	}
	hash, err := HashPassword(newPassword)
	if err != nil {
		return "", err
	}
	user.PasswordHash = hash
	return newPassword, s.repo.Update(user)
}
