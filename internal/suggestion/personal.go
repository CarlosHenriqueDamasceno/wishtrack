package suggestion

import (
	"context"

	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/content"
	"github.com/google/uuid"
)

type PersonalSuggester struct {
	repository content.Repository
}

func NewPersonalSuggester(repository content.Repository) *PersonalSuggester {
	return &PersonalSuggester{
		repository: repository,
	}
}

func (s *PersonalSuggester) Suggest(ctx context.Context, numberOfSuggestions int, userId uuid.UUID) ([]Suggestion, error) {
	contents, err := s.repository.Feed(ctx, userId, numberOfSuggestions)
	if err != nil {
		return nil, err
	}

	suggestions := make([]Suggestion, numberOfSuggestions)
	for i, content := range contents {
		suggestions[i] = Suggestion{
			Name:     string(content.Name),
			Category: content.Category,
			Genres:   content.Genres,
			Summary:  content.Summary,
		}
	}
	return suggestions, nil
}

func (s *PersonalSuggester) Name() string {
	return "Personal"
}

func (s *PersonalSuggester) Type() SuggesterType {
	return PERSONAL
}
