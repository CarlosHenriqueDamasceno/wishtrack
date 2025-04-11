package test

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/CarlosHenriqueDamasceno/wishtrack/cmd/api/server"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/content"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/suggestion"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/user"
	httputils "github.com/CarlosHenriqueDamasceno/wishtrack/pkg/http"
	"github.com/stretchr/testify/suite"
)

const (
	suggestionsBaseUrl = "/api/v1/suggestions"
	tmdbBaseUrl        = "https://api.themoviedb.org"
	tmdbApiKey         = "eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOiJhMTVkYzFjMDgzYTkxYmFjM2Y2YjAxMDY4MDcyOGI5OSIsIm5iZiI6MTc0NDMzMDIxOC4zNTgsInN1YiI6IjY3Zjg1ZGVhZDRjNDQ0YTFjYzk5OTg2MSIsInNjb3BlcyI6WyJhcGlfcmVhZCJdLCJ2ZXJzaW9uIjoxfQ.Hw0Rn85k2w9GPfmn7UBlip_IhfPIk-o5leM6alir8b8"
)

type SuggestionsTestSuite struct {
	LoggedRequestBaseSuite
}

func (suite *SuggestionsTestSuite) SetupSuite() {
	suite.SetupDatabase()
}

func (suite *SuggestionsTestSuite) SetupTest() {
	userRepository := user.NewDatabaseRepository(suite.conn)

	auth := user.NewJwtAuthenticator(
		secret,
		aud,
		aud,
		time.Minute,
	)

	contentRepository := content.NewDatabaseRepository(suite.conn)
	suite.UserService = user.NewService(userRepository, auth)
	suite.ContentService = content.NewService(contentRepository)

	client := &http.Client{
		Transport: httputils.NewAuthenticatedTransport(tmdbApiKey),
	}

	tdmb, err := suggestion.NewSuggester(suggestion.TMDB, client, tmdbBaseUrl, nil)
	suite.Assert().Nil(err)

	personal, err := suggestion.NewSuggester(suggestion.PERSONAL, client, "", contentRepository)
	suite.Assert().Nil(err)

	suggesters := []suggestion.Suggester{tdmb, personal}

	suite.SuggestionService = suggestion.NewService(suggesters)

	suite.server = server.NewApi(
		http.NewServeMux(),
		&server.Config{},
		slog.Default(),
		suite.UserService,
		suite.ContentService,
		suite.SuggestionService,
	)

	suite.MockUser(DefaultUserEmail, DefaultPassword)
}

func (suite *SuggestionsTestSuite) TearDownTest() {
	suite.ClearDatabase()
}

func (suite *SuggestionsTestSuite) TearDownSuite() {
	suite.DestroyDatabase(context.Background())
}

type ProviderResponse struct {
	Name        string       `json:"name"`
	Suggestions []Suggestion `json:"suggestions"`
}

type Suggestion struct {
	Name     string
	Category string
	Genres   []string
	Summary  string
}

func (suite *SuggestionsTestSuite) mockContents() []*content.WriteDownOutput {
	ctx := context.Background()
	var output []*content.WriteDownOutput

	inputs := []*content.WriteDownInput{
		{
			Name:      "The Lord of the Rings: The Return of the King",
			Category:  "book",
			Genres:    []string{"fantasy", "adventure"},
			Summary:   "The third movie from the series The Lord of The Rings",
			WishLevel: 5,
			UserID:    suite.User.ID,
		},
		{
			Name:     "Saving private Ryan",
			Category: "movie",
			Genres:   []string{"war", "historical"},
			Summary: `Inspired by the books of Stephen E. Ambrose and accounts of multiple soldiers in a single
			family, such as the Niland brothers, being killed in action`,
			WishLevel: 4,
			UserID:    suite.User.ID,
		},
		{
			Name:      "The Shawshank Redemption",
			Category:  "movie",
			Genres:    []string{"drama", "crime"},
			Summary:   "adasdadsa",
			WishLevel: 3,
			UserID:    suite.User.ID,
		},
		{
			Name:      "The Godfather",
			Category:  "movie",
			Genres:    []string{"crime", "drama"},
			Summary:   "adasdadsa",
			WishLevel: 2,
			UserID:    suite.User.ID,
		},
		{
			Name:      "The Dark Knight",
			Category:  "movie",
			Genres:    []string{"action", "crime"},
			Summary:   "adasdadsa",
			WishLevel: 1,
			UserID:    suite.User.ID,
		},
	}

	for _, movie := range inputs {
		out, err := suite.ContentService.WriteDown(ctx, movie)
		suite.Assert().Nil(err)
		output = append(output, out)
	}

	return output
}

func (suite *SuggestionsTestSuite) TestItShouldGetSuggestions() {
	suggestions := suite.mockContents()

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, suggestionsBaseUrl, nil)
	suite.MockToken(DefaultUserEmail, DefaultPassword, req)

	suite.server.ServeHTTP(recorder, req)
	suite.Assert().Equal(http.StatusOK, recorder.Result().StatusCode)

	response := make(map[string]ProviderResponse)

	err := json.NewDecoder(recorder.Body).Decode(&response)
	suite.Assert().Nil(err)
	suite.Assert().NotEmpty(response)

	for providerId, provider := range response {
		if providerId == string(suggestion.PERSONAL) {
			suite.Assert().NotEmpty(provider.Suggestions)
			suite.Assert().Len(provider.Suggestions, len(suggestions))
			for i, suggestion := range provider.Suggestions {
				suite.Assert().Equal(string(suggestions[i].Name), suggestion.Name)
				suite.Assert().Equal(suggestions[i].Category, suggestion.Category)
				suite.Assert().Equal(suggestions[i].Summary, suggestion.Summary)
				suite.Assert().Equal(len(suggestions[i].Genres), len(suggestion.Genres))
				suite.Assert().Equal([]string(suggestions[i].Genres), suggestion.Genres)
			}
		} else {
			suite.Assert().Equal(string(suggestion.TMDB), providerId)
		}
	}
}

func (suite *SuggestionsTestSuite) TestItShouldGetTheExactlyNumberOfSuggestions() {
	suite.mockContents()

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, suggestionsBaseUrl, nil)

	newQuery := req.URL.Query()
	newQuery.Add("quantity", "2")
	req.URL.RawQuery = newQuery.Encode()

	suite.MockToken(DefaultUserEmail, DefaultPassword, req)

	suite.server.ServeHTTP(recorder, req)
	suite.Assert().Equal(http.StatusOK, recorder.Result().StatusCode)

	var response []ProviderResponse

	err := json.NewDecoder(recorder.Body).Decode(&response)
	suite.Assert().Nil(err)
	suite.Assert().Len(response[0].Suggestions, 2)
}

func TestSuggestionsTestSuite(t *testing.T) {
	suite.Run(t, new(SuggestionsTestSuite))
}
