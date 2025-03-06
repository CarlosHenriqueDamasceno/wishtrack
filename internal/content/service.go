package content

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	WriteDown(context.Context, *WriteDownInput) (*WriteDownOutput, error)
	Feed(context.Context, uuid.UUID) (*FeedOutput, error)
	Edit(context.Context, *EditContentInput) (*EditContentOutput, error)
	Rate(context.Context, *RateContentInput) error
	Find(context.Context, uuid.UUID, uuid.UUID) (*FindContentOutput, error)
}

type service struct {
	repository Repository
}

func NewService(repo Repository) Service {
	return &service{
		repository: repo,
	}
}
