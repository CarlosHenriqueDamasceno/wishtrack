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
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/user"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

const writeDownBaseUrl = "/api/v1/contents/write-down"

type WriteDownTestSuite struct {
	LoggedRequestBaseSuite
}

func (suite *WriteDownTestSuite) SetupSuite() {
	suite.SetupDatabase()
}

func (suite *WriteDownTestSuite) SetupTest() {
	userRepository := user.NewDatabaseRepository(suite.conn)
	contentRepository := content.NewDatabaseRepository(suite.conn)

	auth := user.NewJwtAuthenticator(
		secret,
		aud,
		aud,
		time.Minute,
	)

	suite.userService = user.NewService(userRepository, auth)
	suite.contentService = content.NewService(contentRepository)

	suite.server = server.NewApi(
		http.NewServeMux(),
		&server.Config{},
		slog.Default(),
		suite.userService,
		suite.contentService,
	)

	suite.mockUser(DefaultUserEmail, DefaultPassword)
}

func (suite *WriteDownTestSuite) TearDownTest() {
	suite.ClearDatabase()
}

func (suite *WriteDownTestSuite) TearDownSuite() {
	suite.DestroyDatabase(context.Background())
}

func (suite *WriteDownTestSuite) TestShouldFailToWriteDownWithoutName() {
	input := struct {
		Name      string   `json:"name"`
		Category  string   `json:"category"`
		Genres    []string `json:"genres"`
		Summary   string   `json:"summary"`
		WishLevel int      `json:"wish_level"`
	}{
		Name:      "",
		Category:  "movie",
		Genres:    []string{"fantasy", "adventure"},
		Summary:   "The third movie from the series The Lord of The Rings",
		WishLevel: 5,
	}

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, writeDownBaseUrl, PrepareBody(input, &suite.Suite))
	suite.mockToken(DefaultUserEmail, DefaultPassword, req)

	suite.server.ServeHTTP(recorder, req)
	suite.Assert().Equal(http.StatusUnprocessableEntity, recorder.Result().StatusCode)

	resp := &struct {
		Errors map[string][]string `json:"errors"`
	}{}

	err := json.NewDecoder(recorder.Result().Body).Decode(resp)
	suite.Assert().Nil(err, "fail to parse response")

	expectedErrors := map[string][]string{"name": {"field \"name\" is required"}}
	suite.Assert().Equal(expectedErrors, resp.Errors)

	suite.AssertDatabaseCount(0, "contents", "id")
}

func (suite *WriteDownTestSuite) TestShouldFailToWriteDownInvalidName() {
	input := struct {
		Name      string   `json:"name"`
		Category  string   `json:"category"`
		Genres    []string `json:"genres"`
		Summary   string   `json:"summary"`
		WishLevel int      `json:"wish_level"`
	}{
		Name:      "as",
		Category:  "movie",
		Genres:    []string{"fantasy", "adventure"},
		Summary:   "The third movie from the series The Lord of The Rings",
		WishLevel: 5,
	}

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, writeDownBaseUrl, PrepareBody(input, &suite.Suite))
	suite.mockToken(DefaultUserEmail, DefaultPassword, req)

	suite.server.ServeHTTP(recorder, req)
	suite.Assert().Equal(http.StatusUnprocessableEntity, recorder.Result().StatusCode)

	resp := &struct {
		Errors map[string][]string `json:"errors"`
	}{}

	err := json.NewDecoder(recorder.Result().Body).Decode(resp)
	suite.Assert().Nil(err, "fail to parse response")

	expectedErrors := map[string][]string{"name": {"field \"name\" must be at least 3 characters long"}}
	suite.Assert().Equal(expectedErrors, resp.Errors)

	suite.AssertDatabaseCount(0, "contents", "id")
}

func (suite *WriteDownTestSuite) TestShouldFailToWriteDownWithInvalidWishLevel() {
	input := struct {
		Name      string   `json:"name"`
		Category  string   `json:"category"`
		Genres    []string `json:"genres"`
		Summary   string   `json:"summary"`
		WishLevel int      `json:"wish_level"`
	}{
		Name:     "The Lord of the Rings: The Return of the King",
		Category: "movie",
		Genres:   []string{"fantasy", "adventure"},
		Summary:  "The third movie from the series The Lord of The Rings",
	}

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, writeDownBaseUrl, PrepareBody(input, &suite.Suite))
	suite.mockToken(DefaultUserEmail, DefaultPassword, req)

	suite.server.ServeHTTP(recorder, req)
	suite.Assert().Equal(http.StatusUnprocessableEntity, recorder.Result().StatusCode)

	resp := &struct {
		Errors map[string][]string `json:"errors"`
	}{}

	err := json.NewDecoder(recorder.Result().Body).Decode(resp)
	suite.Assert().Nil(err, "fail to parse response")

	expectedErrors := map[string][]string{"wish_level": {"field \"wish_level\" must be between 1 and 5"}}
	suite.Assert().Equal(expectedErrors, resp.Errors)

	suite.AssertDatabaseCount(0, "contents", "id")
}

func (suite *WriteDownTestSuite) TestShouldWriteDown() {
	input := struct {
		Name      string   `json:"name"`
		Category  string   `json:"category"`
		Genres    []string `json:"genres"`
		Summary   string   `json:"summary"`
		WishLevel int      `json:"wish_level"`
	}{
		Name:      "The Lord of the Rings: The Return of the King",
		Category:  "movie",
		Genres:    []string{"fantasy", "adventure"},
		Summary:   "The third movie from the series The Lord of The Rings",
		WishLevel: 5,
	}

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, writeDownBaseUrl, PrepareBody(input, &suite.Suite))
	suite.mockToken(DefaultUserEmail, DefaultPassword, req)

	suite.server.ServeHTTP(recorder, req)
	suite.Assert().Equal(http.StatusCreated, recorder.Result().StatusCode)

	responseBody := struct {
		ID        string    `json:"id"`
		Name      string    `json:"name"`
		Category  string    `json:"category"`
		Genres    []string  `json:"genres"`
		Summary   string    `json:"summary"`
		WishLevel int       `json:"wish_level"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}{}

	err := json.NewDecoder(recorder.Result().Body).Decode(&responseBody)
	suite.Assert().Nil(err, "Fail to unmarshal response: %s")

	err = uuid.Validate(responseBody.ID)
	suite.Assert().Nil(err, "Result ID is invalid: %s")

	suite.Assert().NotZero(responseBody.CreatedAt, "Created at must be defined")
	suite.Assert().NotZero(responseBody.UpdatedAt, "Updated at must be defined")
	suite.Assert().Equal(input.Genres, responseBody.Genres)
}

func TestWriteDownTestSuite(t *testing.T) {
	suite.Run(t, new(WriteDownTestSuite))
}
