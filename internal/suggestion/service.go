package suggestion

import (
	"context"

	"github.com/google/uuid"
)

type ProviderResponse struct {
	Provider    string       `json:"provider"`
	Suggestions []Suggestion `json:"suggestions"`
}

type Service interface {
	Suggest(ctx context.Context, numberOfSuggestions int, userId uuid.UUID) (map[string]ProviderResponse, error)
}

type service struct {
	Suggesters []Suggester
}

func NewService(suggesters []Suggester) Service {
	return &service{
		Suggesters: suggesters,
	}
}

func (s *service) Suggest(ctx context.Context, numberOfSuggestions int, userId uuid.UUID) (map[string]ProviderResponse, error) {
	result := make(map[string]ProviderResponse)

	for _, suggester := range s.Suggesters {
		suggestions, err := suggester.Suggest(ctx, numberOfSuggestions, userId)
		if err != nil {
			return nil, err
		}

		result[string(suggester.Type())] = ProviderResponse{
			Provider:    suggester.Name(),
			Suggestions: suggestions,
		}
	}

	return result, nil
}
