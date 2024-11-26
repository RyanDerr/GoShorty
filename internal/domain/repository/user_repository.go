package repository

import (
	"fmt"
	"html"
	"strings"

	"github.com/RyanDerr/GoShorty/internal/domain/entity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

type IUserRepository interface {
	Save(*entity.User) (*entity.User, error)
	BeforeSave(*entity.User) error
	UsernameExists(string) (bool, error)
	GetUserByUsername(string) (*entity.User, error)
	GetUserById(uint) (*entity.User, error)
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Save(user *entity.User) (*entity.User, error) {
	if err := r.BeforeSave(user); err != nil {
		return nil, err
	}

	err := r.db.Create(user).Error
	if err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) BeforeSave(user *entity.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = string(hash)
	user.Username = html.EscapeString(strings.TrimSpace(user.Username))
	return nil
}

func (r *UserRepository) UsernameExists(username string) (bool, error) {
	var user entity.User
	err := r.db.Where("username = ?", username).First(&user).Error

	if err == gorm.ErrRecordNotFound {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("failed to check if username exists: %w", err)
	}

	return true, nil
}

func (r *UserRepository) GetUserByUsername(username string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("username = ?", username).First(&user).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetUserById(id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.Preload("Entries").Where("id = ?", id).First(&user).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}
