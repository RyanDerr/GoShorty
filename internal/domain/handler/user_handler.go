package handler

import (
	"net/http"

	"github.com/RyanDerr/GoShorty/internal/domain/service"
	"github.com/RyanDerr/GoShorty/pkg/helper"
	"github.com/RyanDerr/GoShorty/pkg/mapper"
	"github.com/RyanDerr/GoShorty/pkg/request"
	"github.com/RyanDerr/GoShorty/pkg/response"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.IUserService
}

func NewUserHandler(userService service.IUserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// RegisterUser godoc
//
//	@Summary		Register a new user
//	@Description	Register a new user with a username and password
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		request.UserAuthInput		true	"User credentials"
//	@Success		201		{object}	response.UserResponse		"User registered successfully"
//	@Failure		400		{object}	response.ResponseErrorModel	"Bad Request"
//	@Failure		409		{object}	response.ResponseErrorModel	"Conflict"
//	@Failure		500		{object}	response.ResponseErrorModel	"Internal Server Error"
//	@Router			/auth/register [post]
func (h *UserHandler) RegisterUser(ctx *gin.Context) {
	var input request.UserAuthInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		response.ResponseError(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	err := input.Validate()
	if err != nil {
		response.ResponseError(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	user, status, err := h.userService.CreateUser(ctx, mapper.MapUserAuthInputToEntity(&input))
	if err != nil {
		response.ResponseError(ctx, err.Error(), status)
		return
	}

	response.ResponseCreatedWithData(ctx, mapper.MapUserEntityToResponse(user))
}

// LoginUser godoc
//
//	@Summary		Login a user
//	@Description	Authenticate a user and return a JWT token
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		request.UserAuthInput		true	"User credentials"
//	@Success		200		{object}	response.JwtResponse		"JWT token"
//	@Failure		400		{object}	response.ResponseErrorModel	"Bad Request"
//	@Failure		401		{object}	response.ResponseErrorModel	"Unauthorized"
//	@Failure		404		{object}	response.ResponseErrorModel	"Not Found"
//	@Failure		500		{object}	response.ResponseErrorModel	"Internal Server Error"
//	@Router			/auth/login [post]
func (h *UserHandler) LoginUser(ctx *gin.Context) {
	var input request.UserAuthInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		response.ResponseError(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	user, status, err := h.userService.GetUserByUsername(input.Username)
	if err != nil {
		response.ResponseError(ctx, err.Error(), status)
		return
	}

	status, err = h.userService.ValidatePassword(input.Username, input.Password)
	if err != nil {
		response.ResponseError(ctx, err.Error(), status)
		return
	}

	jwt, err := helper.GenerateJWT(user)
	if err != nil {
		response.ResponseError(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	response.ResponseOKWithData(ctx, mapper.MapSignedJwtToResponse(jwt))
}
