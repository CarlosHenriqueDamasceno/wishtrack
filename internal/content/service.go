package content

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	WriteDown(context.Context, *WriteDownInput) (*WriteDownOutput, error)
	Feed(context.Context, uuid.UUID) (*FeedOutput, error)
}

type service struct {
	repository Repository
}

func NewService(repo Repository) Service {
	return &service{
		repository: repo,
	}
}
