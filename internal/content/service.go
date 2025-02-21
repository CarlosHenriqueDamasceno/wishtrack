package content

import (
	"context"
)

type Service interface {
	WriteDown(context.Context, *WriteDownInput) (*WriteDownOutput, error)
}

type service struct {
	repository Repository
}

func NewService(repo Repository) Service {
	return &service{
		repository: repo,
	}
}
