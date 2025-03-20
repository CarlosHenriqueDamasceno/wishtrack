package content

import (
	"context"

	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/query"
	"github.com/google/uuid"
)

type Repository interface {
	Create(context.Context, *Content) error
	Find(context.Context, uuid.UUID) (*Content, error)
	List(context.Context, uuid.UUID, query.PaginationInput, ContentListFilters) (data []*Content, total uint64, err error)
	Feed(context.Context, uuid.UUID) ([]*Content, error)
	Update(context.Context, *Content) error
	Delete(context.Context, uuid.UUID) error
}
