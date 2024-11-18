package service

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RyanDerr/GoShorty/internal/domain/entity"
	"github.com/RyanDerr/GoShorty/internal/domain/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo repository.IUserRepository
}

type IUserService interface {
	CreateUser(*gin.Context, *entity.User) (*entity.User, int, error)
	GetUserByUsername(string) (*entity.User, int, error)
	ValidatePassword(string, string) (int, error)
	GetUserById(uint) (*entity.User, int, error)
}

func NewUserService(userRepo repository.IUserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) CreateUser(ctx *gin.Context, user *entity.User) (*entity.User, int, error) {
	exist, err := s.userRepo.UsernameExists(user.Username)

	if err != nil {
		log.Printf("Error checking if username exists: %s", err.Error())
		return nil, http.StatusInternalServerError, err
	}

	if exist {
		log.Printf("Username already exists: %s", user.Username)
		return nil, http.StatusConflict, fmt.Errorf("username already exists: %s", user.Username)
	}

	response, err := s.userRepo.Save(user)

	if err != nil {
		log.Printf("Error saving user: %s", err.Error())
		return nil, http.StatusInternalServerError, err
	}

	return response, http.StatusCreated, nil
}

func (s *UserService) GetUserByUsername(username string) (*entity.User, int, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		log.Printf("Error getting user by username: %s", err.Error())
		return nil, http.StatusInternalServerError, err
	}

	if user == nil {
		log.Printf("User not found: %s", username)
		return nil, http.StatusNotFound, fmt.Errorf("user not found: %s", username)
	}

	return user, http.StatusOK, nil
}

func (s *UserService) ValidatePassword(username, password string) (int, error) {
	user, err := s.userRepo.GetUserByUsername(username)
	if err != nil {
		log.Printf("Error getting user by username: %s", err.Error())
		return http.StatusInternalServerError, err
	}

	if user == nil {
		log.Printf("User not found: %s", username)
		return http.StatusNotFound, fmt.Errorf("user not found: %s", username)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Printf("Error validating password: %s", err.Error())
		return http.StatusUnauthorized, fmt.Errorf("invalid password")
	}

	return http.StatusOK, nil
}

func (s *UserService) GetUserById(id uint) (*entity.User, int, error) {
	user, err := s.userRepo.GetUserById(id)
	if err != nil {
		log.Printf("Error getting users: %s", err.Error())
		return nil, http.StatusInternalServerError, err
	}

	if user == nil {
		log.Printf("User not found: %d", id)
		return nil, http.StatusNotFound, fmt.Errorf("user not found: %d", id)
	}

	return user, http.StatusOK, nil
}
