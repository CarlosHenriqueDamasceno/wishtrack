package content

import (
	"context"
	"errors"
	"time"

	"github.com/CarlosHenriqueDamasceno/wishtrack/pkg/database"
	"github.com/CarlosHenriqueDamasceno/wishtrack/pkg/validation"
	"github.com/google/uuid"
)

var ErrContentNotFound = errors.New("content not found")

type EditContentInput struct {
	ID        uuid.UUID `json:"-"`
	Name      string    `json:"name"`
	Category  string    `json:"category"`
	Genres    Genres    `json:"genres"`
	Summary   string    `json:"summary"`
	WishLevel int       `json:"wish_level"`
	UserID    uuid.UUID `json:"-"`
}

type EditContentOutput struct {
	ID        uuid.UUID `json:"id"`
	Name      Name      `json:"name"`
	Category  string    `json:"category"`
	Genres    Genres    `json:"genres"`
	Summary   string    `json:"summary"`
	WishLevel WishLevel `json:"wish_level"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *service) Edit(ctx context.Context, input *EditContentInput) (*EditContentOutput, error) {
	content, err := s.repository.Find(ctx, input.ID)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrRecordNotFound):
			return nil, ErrContentNotFound
		default:
			return nil, err
		}
	}

	if content.UserID != input.UserID {
		return nil, ErrContentNotFound
	}

	content.Name = Name(input.Name)
	content.Category = input.Category
	content.Genres = input.Genres
	content.Summary = input.Summary
	content.WishLevel = WishLevel(input.WishLevel)

	err = s.validateEdit(content)
	if err != nil {
		return nil, err
	}

	err = s.repository.Update(ctx, content)
	if err != nil {
		return nil, err
	}

	return (*EditContentOutput)(outputFromContent(content)), nil
}

func (s *service) validateEdit(content *Content) error {
	var errors validation.ErrorCollection

	if err := content.Name.validate(); err != nil {
		errors.Append(err)
	}

	if err := content.WishLevel.validate(); err != nil {
		errors.Append(err)
	}

	if errors.HasError() {
		return errors
	}

	return nil
}
