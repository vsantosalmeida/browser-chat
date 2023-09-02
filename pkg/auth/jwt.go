package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/vsantosalmeida/browser-chat/entity"
)

const secret = "none"
const expirationTime = 86400 // 1 day

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
