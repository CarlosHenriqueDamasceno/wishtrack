package content

import (
	"context"
	"time"

	"github.com/CarlosHenriqueDamasceno/wishtrack/pkg/validation"
	"github.com/google/uuid"
)

type WriteDownInput struct {
	Name      string `json:"name"`
	Category  string `json:"category"`
	Genres    Genres `json:"genres"`
	Summary   string `json:"summary"`
	WishLevel int    `json:"wish_level"`
	UserID    uuid.UUID
}

type WriteDownOutput struct {
	ID        uuid.UUID `json:"id"`
	Name      Name      `json:"name"`
	Category  string    `json:"category"`
	Genres    Genres    `json:"genres"`
	Summary   string    `json:"summary"`
	WishLevel WishLevel `json:"wish_level"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *service) WriteDown(ctx context.Context, input *WriteDownInput) (*WriteDownOutput, error) {
	content := NewContent(
		input.Name,
		input.Category,
		input.Genres,
		input.Summary,
		input.WishLevel,
		input.UserID,
	)

	err := s.validateWriteDown(content)
	if err != nil {
		return nil, err
	}

	err = s.repository.Create(ctx, content)
	if err != nil {
		return nil, err
	}

	return outputFromContent(content), nil
}

func (s *service) validateWriteDown(content *Content) error {
	var errors validation.ErrorCollection
	if err := content.Name.validate(); err != nil {
		errors.Append(err)
	}

	if genresErrors := content.Genres.validate(); genresErrors.HasError() {
		errors.AppendCollection(&genresErrors)
	}

	if err := content.WishLevel.validate(); err != nil {
		errors.Append(err)
	}

	if errors.HasError() {
		return errors
	}

	return nil
}

func outputFromContent(content *Content) *WriteDownOutput {
	return &WriteDownOutput{
		ID:        content.ID,
		Name:      content.Name,
		Category:  content.Category,
		Genres:    content.Genres,
		Summary:   content.Summary,
		WishLevel: content.WishLevel,
		CreatedAt: content.CreatedAt,
		UpdatedAt: content.UpdatedAt,
	}
}
