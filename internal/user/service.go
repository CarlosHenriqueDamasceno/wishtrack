package user

import (
	"context"
	"errors"
)

var ErrUserNotFound = errors.New("user not found")

type Service interface {
	Register(ctx context.Context, input *RegisterInput) (*RegisterOutput, error)
	Login(ctx context.Context, input *LoginInput) (*LoginOutput, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}
