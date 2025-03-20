package content

import (
	"context"

	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/query"
	"github.com/google/uuid"
)

type Service interface {
	WriteDown(ctx context.Context, input *WriteDownInput) (*WriteDownOutput, error)
	Feed(ctx context.Context, id uuid.UUID) (*FeedOutput, error)
	Edit(ctx context.Context, input *EditContentInput) (*EditContentOutput, error)
	Rate(ctx context.Context, input *RateContentInput) error
	Find(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*FindContentOutput, error)
	Delete(ctx context.Context, id, userID uuid.UUID) error
	List(ctx context.Context, userID uuid.UUID, pagination query.PaginationInput, filters ContentListFilters) (*ListContentsOutput, error)
}

type service struct {
	repository Repository
}

func NewService(repo Repository) Service {
	return &service{
		repository: repo,
	}
}
