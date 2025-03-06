package content

import (
	"context"
	"errors"
	"time"

	"github.com/CarlosHenriqueDamasceno/wishtrack/pkg/database"
	"github.com/google/uuid"
)

type FindContentOutput struct {
	ID        uuid.UUID `json:"id"`
	Name      Name      `json:"name"`
	Category  string    `json:"category"`
	Genres    []string  `json:"genres"`
	Summary   string    `json:"summary"`
	WishLevel WishLevel `json:"wish_level"`
	Rate      *Rate     `json:"rate"`
	Comment   *string   `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func findOutputFromContent(content *Content) *FindContentOutput {
	return &FindContentOutput{
		ID:        content.ID,
		Name:      content.Name,
		Category:  content.Category,
		Genres:    content.Genres,
		Summary:   content.Summary,
		WishLevel: content.WishLevel,
		Rate:      content.Rate,
		Comment:   content.Comment,
		CreatedAt: content.CreatedAt,
		UpdatedAt: content.UpdatedAt,
	}
}

func (s *service) Find(ctx context.Context, id, userID uuid.UUID) (*FindContentOutput, error) {
	content, err := s.repository.Find(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrRecordNotFound):
			return nil, ErrContentNotFound
		default:
			return nil, err
		}
	}

	if content.UserID != userID {
		return nil, ErrContentNotFound
	}

	return findOutputFromContent(content), nil
}
