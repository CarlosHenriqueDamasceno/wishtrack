package content

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type output struct {
	ID        uuid.UUID `json:"id"`
	Name      Name      `json:"name"`
	Category  string    `json:"category"`
	Genres    []string  `json:"genres"`
	Summary   string    `json:"summary"`
	WishLevel WishLevel `json:"wish_level"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FeedOutput []output

func (service *service) Feed(ctx context.Context, id uuid.UUID) (*FeedOutput, error) {
	contents, err := service.repository.Feed(ctx, id)
	if err != nil {
		return nil, err
	}

	return contentsToOutput(contents), nil
}

func contentsToOutput(contents []*Content) *FeedOutput {
	var out FeedOutput

	for _, content := range contents {
		o := output{
			ID:        content.ID,
			Name:      content.Name,
			Category:  content.Category,
			Genres:    content.Genres,
			Summary:   content.Summary,
			WishLevel: content.WishLevel,
			CreatedAt: content.CreatedAt,
			UpdatedAt: content.UpdatedAt,
		}

		out = append(out, o)
	}

	return &out
}
