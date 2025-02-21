package user

import (
	"context"

	"github.com/google/uuid"
)

func (service *service) Find(ctx context.Context, id uuid.UUID) (*User, error) {
	return service.repository.Find(ctx, id)
}
