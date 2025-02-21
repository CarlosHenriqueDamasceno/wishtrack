package user

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Authenticator interface {
	CreateToken(*User) (string, error)
	ValidateToken(string) (userId uuid.UUID, err error)
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

func (a *JwtAuthenticator) ValidateToken(token string) (userId uuid.UUID, err error) {
	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}

		return []byte(a.key), nil
	},
		jwt.WithExpirationRequired(),
		jwt.WithAudience(a.aud),
		jwt.WithIssuer(a.iss),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)
	if err != nil {
		return [16]byte{}, err
	}

	id, err := jwtToken.Claims.GetSubject()
	if err != nil {
		return [16]byte{}, err
	}

	uuid, err := uuid.Parse(id)
	if err != nil {
		return [16]byte{}, errors.New("invalid subject")
	}

	return uuid, nil
}
