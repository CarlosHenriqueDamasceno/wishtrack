package content

import (
	"context"
	"errors"

	"github.com/CarlosHenriqueDamasceno/wishtrack/pkg/database"
	"github.com/google/uuid"
)

func (s *service) Delete(ctx context.Context, id, userID uuid.UUID) error {
	content, err := s.repository.Find(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrRecordNotFound):
			return ErrContentNotFound
		default:
			return err
		}
	}

	if content.UserID != userID {
		return ErrContentNotFound
	}

	return s.repository.Delete(ctx, id)
}
