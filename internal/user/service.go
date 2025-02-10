package user

import "context"

type Service interface {
	Register(ctx context.Context, input *RegisterInput) (*RegisterOutput, error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}
