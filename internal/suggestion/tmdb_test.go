package suggestion_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/suggestion"
	httphelper "github.com/CarlosHenriqueDamasceno/wishtrack/pkg/http"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const fakeToken = "fake_token"

type Suggestion struct {
	Name     string
	Category string
	Genres   []string
	Summary  string
}

type TMDBSuggestion struct {
	BackdropPath     string  `json:"backdrop_path"`
	Id               int     `json:"id"`
	Title            string  `json:"title"`
	OriginalTitle    string  `json:"original_title"`
	Overview         string  `json:"overview"`
	PosterPath       string  `json:"poster_path"`
	MediaType        string  `json:"media_type"`
	Adult            bool    `json:"adult"`
	OriginalLanguage string  `json:"original_language"`
	GenreIds         []int   `json:"genre_ids"`
	Popularity       float64 `json:"popularity"`
	ReleaseDate      string  `json:"release_date"`
	Video            bool    `json:"video"`
	VoteAverage      float64 `json:"vote_average"`
	VoteCount        int     `json:"vote_count"`
}

type TMDBResponse struct {
	Page         int              `json:"page"`
	Results      []TMDBSuggestion `json:"results"`
	TotalPages   int              `json:"total_pages"`
	TotalResults int              `json:"total_results"`
}

var tmdbResponses = TMDBResponse{
	Page: 1,
	Results: []TMDBSuggestion{
		{
			BackdropPath:     "/fake_backdrop_path_1.jpg",
			Id:               1,
			Title:            "Lord of the Rings",
			OriginalTitle:    "Lord of the Rings",
			Overview:         "An action-packed adventure.",
			PosterPath:       "/fake_poster_path_1.jpg",
			MediaType:        "Movie",
			Adult:            false,
			OriginalLanguage: "en",
			GenreIds:         []int{28, 12}, // Action, Adventure
			Popularity:       9.8,
			ReleaseDate:      "2001-12-19",
			Video:            false,
			VoteAverage:      8.7,
			VoteCount:        15000,
		},
		{
			BackdropPath:     "/fake_backdrop_path_2.jpg",
			Id:               2,
			Title:            "Saving Private Ryan",
			OriginalTitle:    "Saving Private Ryan",
			Overview:         "A gripping war drama.",
			PosterPath:       "/fake_poster_path_2.jpg",
			MediaType:        "Movie",
			Adult:            false,
			OriginalLanguage: "en",
			GenreIds:         []int{18, 10752}, // Drama, War
			Popularity:       8.5,
			ReleaseDate:      "1998-07-24",
			Video:            false,
			VoteAverage:      8.6,
			VoteCount:        12000,
		},
		{
			BackdropPath:     "/fake_backdrop_path_3.jpg",
			Id:               3,
			Title:            "Pulp Fiction",
			OriginalTitle:    "Pulp Fiction",
			Overview:         "A classic crime drama.",
			PosterPath:       "/fake_poster_path_3.jpg",
			MediaType:        "Movie",
			Adult:            false,
			OriginalLanguage: "en",
			GenreIds:         []int{80, 18}, // Crime, Drama
			Popularity:       9.0,
			ReleaseDate:      "1994-10-14",
			Video:            false,
			VoteAverage:      8.9,
			VoteCount:        14000,
		},
		{
			BackdropPath:     "/fake_backdrop_path_4.jpg",
			Id:               4,
			Title:            "The Shining",
			OriginalTitle:    "The Shining",
			Overview:         "A mysterious horror film.",
			PosterPath:       "/fake_poster_path_4.jpg",
			MediaType:        "Movie",
			Adult:            false,
			OriginalLanguage: "en",
			GenreIds:         []int{27, 9648}, // Horror, Mystery
			Popularity:       7.8,
			ReleaseDate:      "1980-05-23",
			Video:            false,
			VoteAverage:      8.4,
			VoteCount:        11000,
		},
		{
			BackdropPath:     "/fake_backdrop_path_5.jpg",
			Id:               5,
			Title:            "The Holiday",
			OriginalTitle:    "The Holiday",
			Overview:         "A romantic comedy.",
			PosterPath:       "/fake_poster_path_5.jpg",
			MediaType:        "Movie",
			Adult:            false,
			OriginalLanguage: "en",
			GenreIds:         []int{10749, 35}, // Romance, Comedy
			Popularity:       7.5,
			ReleaseDate:      "2006-12-08",
			Video:            false,
			VoteAverage:      7.2,
			VoteCount:        9000,
		},
	},
	TotalPages:   1,
	TotalResults: 5,
}

type TMDBGenre struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

var tmdbGenres = []TMDBGenre{
	{Id: 28, Name: "Action"},
	{Id: 12, Name: "Adventure"},
	{Id: 18, Name: "Drama"},
	{Id: 10752, Name: "War"},
	{Id: 80, Name: "Crime"},
	{Id: 27, Name: "Horror"},
	{Id: 9648, Name: "Mystery"},
	{Id: 10749, Name: "Romance"},
	{Id: 35, Name: "Comedy"},
}

func mockTMDBService() (*http.Client, *httptest.Server) {
	client := &http.Client{
		Transport: httphelper.NewAuthenticatedTransport(fakeToken),
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Checks
		if r.URL.Query().Get("language") != "pt-BR" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if r.Header.Get("Authorization") != "Bearer "+fakeToken {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Mock responses
		if r.URL.Path == "/3/trending/movie/day" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(tmdbResponses)
			return
		}

		if r.URL.Path == "/3/genre/movie/list" {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(struct {
				Genres []TMDBGenre `json:"genres"`
			}{
				Genres: tmdbGenres,
			})
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}))

	return client, server
}

func TestItShouldProvideSuggestionsFromTMDB(t *testing.T) {
	client, server := mockTMDBService()
	defer server.Close()

	suggester, err := suggestion.NewSuggester(suggestion.TMDB, client, server.URL, nil)
	assert.Nil(t, err)

	result, err := suggester.Suggest(t.Context(), len(tmdbResponses.Results), uuid.UUID{})
	assert.Nil(t, err)

	assert.Equal(t, len(tmdbResponses.Results), len(result))
	for i, suggestion := range result {
		assert.Equal(t, tmdbResponses.Results[i].Title, suggestion.Name)
		assert.Equal(t, tmdbResponses.Results[i].MediaType, suggestion.Category)
		assert.Equal(t, tmdbResponses.Results[i].Overview, suggestion.Summary)
		assert.Equal(t, len(tmdbResponses.Results[i].GenreIds), len(suggestion.Genres))
	}
}

func TestShouldReturnAuthorizationErrorIfTokenIsInvalid(t *testing.T) {
	client, server := mockTMDBService()
	defer server.Close()

	client.Transport = httphelper.NewAuthenticatedTransport("invalid_token")

	suggester, err := suggestion.NewSuggester(suggestion.TMDB, client, server.URL, nil)
	assert.Nil(t, err)

	_, err = suggester.Suggest(t.Context(), len(tmdbResponses.Results), uuid.UUID{})
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, suggestion.ErrUnauthorized)
}

func TestShouldReturnProviderConnectionError(t *testing.T) {
	client, server := mockTMDBService()
	server.Close()

	suggester, err := suggestion.NewSuggester(suggestion.TMDB, client, server.URL, nil)
	assert.Nil(t, err)

	_, err = suggester.Suggest(t.Context(), len(tmdbResponses.Results), uuid.UUID{})
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, suggestion.ErrFailContactProvider)
}

func TestShouldReturnGenericRequestError(t *testing.T) {
	client, server := mockTMDBService()
	defer server.Close()

	server.Config.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	suggester, err := suggestion.NewSuggester(suggestion.TMDB, client, server.URL, nil)
	assert.Nil(t, err)

	_, err = suggester.Suggest(t.Context(), len(tmdbResponses.Results), uuid.UUID{})
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, suggestion.ErrGenericRequestError)
}
