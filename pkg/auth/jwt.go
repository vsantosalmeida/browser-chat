package auth

import (
	"time"

	"github.com/vsantosalmeida/browser-chat/entity"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

type contextKey string

const (
	secret         = "none"
	expirationTime = 86400 // 1 day
	UserContextKey = contextKey("user")
)

type Claims struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func (c *Claims) GetId() int {
	return c.ID
}

func (c *Claims) GetUsername() string {
	return c.Username
}

// CreateJWTToken
func CreateJWTToken(user *entity.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ID":        user.ID,
		"Username":  user.Username,
		"ExpiresAt": time.Now().Unix() + expirationTime,
	})
	tokenString, err := token.SignedString([]byte(secret))

	return tokenString, err
}

// ValidateJWTToken
func ValidateJWTToken(tokenString string) (entity.AuthenticatedUser, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected method")
		}

		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("user not validated")
}
