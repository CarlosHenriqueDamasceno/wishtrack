package user

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Authenticator interface {
	CreateToken(*User) (string, error)
}

type JwtAuthenticator struct {
	key string
	iss string
	aud string
	exp time.Duration
}

func NewJwtAuthenticator(key, iss, aud string, exp time.Duration) Authenticator {
	return &JwtAuthenticator{
		key: key,
		iss: iss,
		aud: aud,
		exp: exp,
	}
}

func (a *JwtAuthenticator) CreateToken(user *User) (string, error) {
	claims := jwt.MapClaims{
		"iss": a.iss,
		"aud": a.aud,
		"sub": user.ID.String(),
		"exp": time.Now().Add(a.exp).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(a.key))
}
