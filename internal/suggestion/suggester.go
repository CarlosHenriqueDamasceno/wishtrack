package suggestion

import (
	"context"
	"fmt"
	"net/http"

	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/content"
	"github.com/google/uuid"
)

var (
	ErrUnauthorized        = fmt.Errorf("could not unauthorize in provider")
	ErrFailContactProvider = fmt.Errorf("failed to contact provider")
	ErrGenericRequestError = fmt.Errorf("generic request error")
)

type SuggesterType string

const (
	TMDB                       SuggesterType = "tmdb"
	DefaultNumberOfSuggestions int           = 5
)

type Suggestion struct {
	Name     string   `json:"name"`
	Category string   `json:"category"`
	Genres   []string `json:"genres"`
	Summary  string   `json:"summary"`
}

type Suggester interface {
	Suggest(context.Context, int, uuid.UUID) ([]Suggestion, error)
	Name() string
	Type() SuggesterType
}

func NewSuggester(
	suggesterType SuggesterType,
	client *http.Client,
	baseURL string,
	repository content.Repository,
) (Suggester, error) {
	switch suggesterType {
	case TMDB:
		return NewTMDBSuggester(client, baseURL), nil
	default:
		return nil, fmt.Errorf("unsupported suggester type: %s", suggesterType)
	}
}
