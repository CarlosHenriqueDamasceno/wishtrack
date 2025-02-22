package content

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(context.Context, *Content) error
	Find(context.Context, uuid.UUID) (*Content, error)
	Feed(context.Context, uuid.UUID) ([]*Content, error)
}
