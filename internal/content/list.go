package content

import (
	"context"

	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/query"
	"github.com/google/uuid"
)

type ListContentsOutput struct {
	Data []*FindContentOutput `json:"data"`
	query.PaginationOutput
}

type ContentListFilters struct {
	Watched   *bool
	Category  *string
	Genres    *[]string
	Name      *string
	Summary   *string
	WishLevel *int
}

func (s *service) List(ctx context.Context, userID uuid.UUID, pagination query.PaginationInput, filters ContentListFilters) (*ListContentsOutput, error) {
	contents, total, err := s.repository.List(ctx, userID, pagination, filters)
	if err != nil {
		return nil, err
	}

	var data []*FindContentOutput

	for _, content := range contents {
		data = append(data, findOutputFromContent(content))
	}

	return &ListContentsOutput{
		Data: data,
		PaginationOutput: query.PaginationOutput{
			Total: total,
			Limit: pagination.Limit,
			Page:  pagination.Page,
		},
	}, nil
}
