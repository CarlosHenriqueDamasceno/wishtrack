package suggestion

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

const (
	SuggestionsUrl = "3/trending/movie/day?language=pt-BR"
	GenresUrl      = "3/genre/movie/list?language=pt-BR"
)

type TMDBSuggester struct {
	client  *http.Client
	baseUrl string
}

type tmdbSuggestion struct {
	Title     string `json:"title"`
	Overview  string `json:"overview"`
	MediaType string `json:"media_type"`
	GenreIds  []int  `json:"genre_ids"`
}

type tmdbResponse struct {
	Page         int              `json:"page"`
	Results      []tmdbSuggestion `json:"results"`
	TotalPages   int              `json:"total_pages"`
	TotalResults int              `json:"total_results"`
}

type tmdbGenre struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func NewTMDBSuggester(client *http.Client, baseUrl string) Suggester {
	return &TMDBSuggester{
		client:  client,
		baseUrl: baseUrl,
	}
}

func (s *TMDBSuggester) Suggest(_ context.Context, numberOfSuggestions int, _ uuid.UUID) ([]Suggestion, error) {
	res, err := s.doRequest(SuggestionsUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var response tmdbResponse

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, err
	}

	tmdbGenres, err := s.fetchGenres()
	if err != nil {
		return nil, err
	}

	suggestions := make([]Suggestion, len(response.Results))
	for i, result := range response.Results {
		genres := s.parseGenres(result.GenreIds, tmdbGenres)
		suggestions[i] = Suggestion{
			Name:     result.Title,
			Category: result.MediaType,
			Genres:   genres,
			Summary:  result.Overview,
		}
	}

	return suggestions[:numberOfSuggestions], nil
}

func (s *TMDBSuggester) Name() string {
	return "The movie database"
}

func (s *TMDBSuggester) Type() SuggesterType {
	return TMDB
}

func (s *TMDBSuggester) fetchGenres() ([]tmdbGenre, error) {
	res, err := s.doRequest(GenresUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result struct {
		Genres []tmdbGenre `json:"genres"`
	}

	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return result.Genres, nil
}

func (s *TMDBSuggester) doRequest(url string) (*http.Response, error) {
	actualUrl := fmt.Sprintf("%s/%s", s.baseUrl, url)
	res, err := s.client.Get(actualUrl)
	if err != nil {
		return nil, fmt.Errorf(err.Error()+". error: %w", ErrFailContactProvider)
	}

	if res.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}

	if res.StatusCode != http.StatusOK {
		err := fmt.Errorf("failed to fetch data from TMDB. status code: %d", res.StatusCode)
		return nil, fmt.Errorf(err.Error()+". error: %w", ErrGenericRequestError)
	}

	return res, nil
}

func (s *TMDBSuggester) parseGenres(genreIds []int, tmdbGenres []tmdbGenre) []string {
	genres := make([]string, len(genreIds))
	for i, genreId := range genreIds {
		for _, genre := range tmdbGenres {
			if genre.Id == genreId {
				genres[i] = genre.Name
			}
		}
	}

	return genres
}
