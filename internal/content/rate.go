package content

import (
	"context"
	"errors"

	"github.com/CarlosHenriqueDamasceno/wishtrack/pkg/database"
	"github.com/CarlosHenriqueDamasceno/wishtrack/pkg/validation"
	"github.com/google/uuid"
)

type RateContentInput struct {
	ID      uuid.UUID `json:"-"`
	UserID  uuid.UUID `json:"-"`
	Rate    int       `json:"rate"`
	Comment string    `json:"comment"`
}

func (s *service) Rate(ctx context.Context, input *RateContentInput) error {
	content, err := s.repository.Find(ctx, input.ID)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrRecordNotFound):
			return ErrContentNotFound
		default:
			return err
		}
	}

	if content.UserID != input.UserID {
		return ErrContentNotFound
	}

	content.Rate = (*Rate)(&input.Rate)
	content.Comment = &input.Comment

	err = s.validateRate(content)
	if err != nil {
		return err
	}

	return s.repository.Update(ctx, content)
}

func (s *service) validateRate(content *Content) error {
	var errors validation.ErrorCollection

	if err := content.Rate.validate(); err != nil {
		errors.Append(err)
	}

	if errors.HasError() {
		return errors
	}

	return nil
}
