package mapper

import (
	"github.com/RyanDerr/GoShorty/internal/domain/entity"
	"github.com/RyanDerr/GoShorty/pkg/request"
	"github.com/RyanDerr/GoShorty/pkg/response"
)

func MapUserAuthInputToEntity(req *request.UserAuthInput) *entity.User {
	return &entity.User{
		Username: req.Username,
		Password: req.Password,
	}
}

func MapUserEntityToResponse(user *entity.User) *response.UserResponse {
	return &response.UserResponse{
		Id:       user.ID,
		Username: user.Username,
	}
}