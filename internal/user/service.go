package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrUserNotFound = errors.New("user not found")

type Service interface {
	Find(ctx context.Context, id uuid.UUID) (*User, error)
	Register(ctx context.Context, input *RegisterInput) (*RegisterOutput, error)
	Login(ctx context.Context, input *LoginInput) (*LoginOutput, error)
	ValidateToken(ctx context.Context, token string) (*User, error)
}

type service struct {
	repository    Repository
	authenticator Authenticator
}

func NewService(r Repository, auth Authenticator) Service {
	return &service{
		repository:    r,
		authenticator: auth,
	}
}
