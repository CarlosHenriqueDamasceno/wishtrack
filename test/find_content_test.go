package test

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/CarlosHenriqueDamasceno/wishtrack/cmd/api/server"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/content"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/user"
	"github.com/stretchr/testify/suite"
)

const findBaseUrl = "/api/v1/contents"

type FindContentTestSuite struct {
	LoggedRequestBaseSuite
}

func (suite *FindContentTestSuite) SetupTest() {
	conn, err := SetupDatabase()
	suite.Assert().Nil(err, "Fail to connect to database")

	suite.conn = conn
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

func (suite *FindContentTestSuite) TearDownTest() {
	suite.conn.Close()
}

func (suite *FindContentTestSuite) mockContent() *content.WriteDownOutput {
	ctx := context.Background()

	movie := &content.WriteDownInput{
		Name:      "The Lord of the Rings: The Return of the King",
		Category:  "book",
		Genres:    []string{"fantasy", "adventure"},
		Summary:   "The third movie from the series The Lord of The Rings",
		WishLevel: 5,
		UserID:    suite.user.ID,
	}

	out, err := suite.contentService.WriteDown(ctx, movie)
	suite.Assert().Nil(err)
	return out
}

func (suite *FindContentTestSuite) TestFindAContent() {
	c := suite.mockContent()

	recorder := httptest.NewRecorder()
	url := fmt.Sprintf("%s/%s", findBaseUrl, c.ID.String())
	req := httptest.NewRequest(http.MethodGet, url, nil)
	suite.mockToken(DefaultUserEmail, DefaultPassword, req)

	suite.server.ServeHTTP(recorder, req)
	suite.Assert().Equal(http.StatusOK, recorder.Result().StatusCode)

	responseBody := struct {
		ID        string         `json:"id"`
		Name      string         `json:"name"`
		Category  string         `json:"category"`
		Genres    content.Genres `json:"genres"`
		Summary   string         `json:"summary"`
		WishLevel int            `json:"wish_level"`
		Rate      *int           `json:"rate"`
		Comment   *string        `json:"comment"`
		CreatedAt time.Time      `json:"created_at"`
		UpdatedAt time.Time      `json:"updated_at"`
	}{}

	err := json.NewDecoder(recorder.Result().Body).Decode(&responseBody)
	suite.Assert().Nil(err, "Fail to unmarshal response: %s")

	suite.Assert().Equal(c.ID.String(), responseBody.ID)
	suite.Assert().Equal(string(c.Name), responseBody.Name)
	suite.Assert().Equal(c.Category, responseBody.Category)
	suite.Assert().Equal(c.Genres, responseBody.Genres)
	suite.Assert().Equal(c.Summary, responseBody.Summary)
	suite.Assert().Equal(int(c.WishLevel), responseBody.WishLevel)
	suite.Assert().Nil(responseBody.Rate)
	suite.Assert().Nil(responseBody.Comment)
	suite.Assert().Equal(c.CreatedAt, responseBody.CreatedAt)
	suite.Assert().Equal(c.UpdatedAt, responseBody.UpdatedAt)
}

func TestFindContentTestSuite(t *testing.T) {
	suite.Run(t, new(FindContentTestSuite))
}
