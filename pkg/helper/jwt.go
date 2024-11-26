package helper

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/RyanDerr/GoShorty/internal/domain/entity"
	"github.com/RyanDerr/GoShorty/internal/domain/service"
	"github.com/RyanDerr/GoShorty/pkg/mapper"
	"github.com/RyanDerr/GoShorty/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const (
	expTime       = 15 * time.Minute
	jwtSigningKey = "JWT_PRIVATE_KEY"
)

type SigningKey struct {
	PrivateKey []byte
}

func getSigningKey() (*SigningKey, error) {
	prvKey := []byte(os.Getenv(jwtSigningKey))
	if len(prvKey) == 0 {
		log.Println("JWT signing key not found")
		return nil, fmt.Errorf("JWT signing key not found: %s", jwtSigningKey)
	}

	return &SigningKey{
		PrivateKey: prvKey,
	}, nil
}

func GenerateJWT(user *entity.User) (string, error) {
	key, err := getSigningKey()
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  user.ID,
		"iat": time.Now().Unix(),
		"eat": time.Now().Add(expTime).Unix(),
	})

	return token.SignedString(key.PrivateKey)
}

func ValidateJWT(context *gin.Context) error {
	key, err := getSigningKey()
	if err != nil {
		return err
	}

	token, err := getToken(context, key)
	if err != nil {
		return err
	}
	_, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return nil
	}
	return fmt.Errorf("invalid token")
}

func CurrentUser(ctx *gin.Context, userSvc service.IUserService) (*response.UserResponse, error) {
	key, err := getSigningKey()
	if err != nil {
		return nil, err
	}

	err = ValidateJWT(ctx)
	if err != nil {
		return nil, err
	}

	token, err := getToken(ctx, key)
	if err != nil {
		return nil, err
	}

	claims, _ := token.Claims.(jwt.MapClaims)
	id := uint(claims["id"].(float64))

	user, _, err := userSvc.GetUserById(id)

	if err != nil {
		return nil, err
	}

	return mapper.MapUserEntityToResponse(user), nil
}

func getToken(ctx *gin.Context, key *SigningKey) (*jwt.Token, error) {
	tokenString := getTokenFromRequest(ctx)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return key.PrivateKey, nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

func getTokenFromRequest(ctx *gin.Context) string {
	bearerToken := ctx.Request.Header.Get("Authorization")
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}
