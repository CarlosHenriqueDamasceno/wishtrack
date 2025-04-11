package suggestion

import (
	"fmt"
	"net/http"
)

var (
	ErrUnauthorized         = fmt.Errorf("unauthorized")
	ErrNotEnoughSuggestions = fmt.Errorf("not enough suggestions")
	ErrFailContactProvider  = fmt.Errorf("failed to contact provider")
	ErrGenericRequestError  = fmt.Errorf("generic request error")
)

type SuggesterType string

const (
	TMDB SuggesterType = "tmdb"
)

type Suggestion struct {
	Name     string
	Category string
	Genres   []string
	Summary  string
}

type Suggester interface {
	Suggest(int) ([]Suggestion, error)
}

func NewSuggester(suggesterType SuggesterType, client *http.Client, baseURL string) (Suggester, error) {
	switch suggesterType {
	case TMDB:
		return NewTMDBSuggester(client, baseURL), nil
	default:
		return nil, fmt.Errorf("unsupported suggester type: %s", suggesterType)
	}
}
