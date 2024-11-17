package jtwtoken

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	id         = "userId"
	exp        = "expTime"
	authorized = "authorized"
	jwtSecret  = "JWT_SIGNATURE"
)

type TokenDetail struct {
	AccessToken    string
	ExpirationTime time.Time
}

type AccessDetails struct {
	UserId     string
	Authorized bool
}

func CreateToken(userId string) (*TokenDetail, error) {
	expTime := time.Now().Add(time.Minute * 15)

	atClaims := &jwt.MapClaims{
		id:         userId,
		exp:        expTime,
		authorized: true,
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)

	tkn, err := at.SignedString([]byte(os.Getenv(jwtSecret)))

	if err != nil {
		return nil, err
	}

	return &TokenDetail{
		AccessToken:    tkn,
		ExpirationTime: expTime,
	}, nil
}

func TokenValid(r *http.Request) error {
	token, err := verifyToken(r)

	if err != nil {
		return err
	}

	if !token.Valid {
		return err
	}

	return nil
}

func ExtractMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := verifyToken(r)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		authorized, ok := claims[authorized].(bool)

		if !ok {
			return nil, err
		}

		userId := claims[id].(string)

		return &AccessDetails{
			UserId:     userId,
			Authorized: authorized,
		}, nil
	}

	return nil, err
}

func extractToken(r *http.Request) string {
	token := r.Header.Get("Authorization")
	strArr := strings.Split(token, " ")

	if len(strArr) == 2 {
		return strArr[1]
	}

	return ""
}

func verifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := extractToken(r)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(os.Getenv(jwtSecret)), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

