package content

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	WriteDown(ctx context.Context, input *WriteDownInput) (*WriteDownOutput, error)
	Feed(ctx context.Context, id uuid.UUID) (*FeedOutput, error)
	Edit(ctx context.Context, input *EditContentInput) (*EditContentOutput, error)
	Rate(ctx context.Context, input *RateContentInput) error
	Find(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*FindContentOutput, error)
	Delete(ctx context.Context, id, userID uuid.UUID) error
}

type service struct {
	repository Repository
}

func NewService(repo Repository) Service {
	return &service{
		repository: repo,
	}
}
