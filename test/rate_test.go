package test

import (
	"context"
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

const rateBaseUrl = "/api/v1/contents"

type RateContentTestSuite struct {
	LoggedRequestBaseSuite
	contentRepository content.Repository
}

func (suite *RateContentTestSuite) SetupSuite() {
	suite.SetupDatabase()
}

func (suite *RateContentTestSuite) SetupTest() {
	userRepository := user.NewDatabaseRepository(suite.conn)
	suite.contentRepository = content.NewDatabaseRepository(suite.conn)

	auth := user.NewJwtAuthenticator(
		secret,
		aud,
		aud,
		time.Minute,
	)

	suite.UserService = user.NewService(userRepository, auth)
	suite.ContentService = content.NewService(suite.contentRepository)

	suite.server = server.NewApi(
		http.NewServeMux(),
		&server.Config{},
		slog.Default(),
		suite.UserService,
		suite.ContentService,
		nil,
	)

	suite.MockUser(DefaultUserEmail, DefaultPassword)
}

func (suite *RateContentTestSuite) TearDownTest() {
	suite.ClearDatabase()
}

func (suite *RateContentTestSuite) TearDownSuite() {
	suite.DestroyDatabase(context.Background())
}

func (suite *RateContentTestSuite) mockContent() *content.WriteDownOutput {
	ctx := context.Background()

	movie := &content.WriteDownInput{
		Name:      "The Lord of the Rings: The Return of the King",
		Category:  "book",
		Genres:    []string{"fantasy", "adventure"},
		Summary:   "The third movie from the series The Lord of The Rings",
		WishLevel: 5,
		UserID:    suite.User.ID,
	}

	out, err := suite.ContentService.WriteDown(ctx, movie)
	suite.Assert().Nil(err)
	return out
}

func (suite *RateContentTestSuite) TestShouldRateAContent() {
	content := suite.mockContent()

	input := struct {
		Rate    int    `json:"rate"`
		Comment string `json:"comment"`
	}{
		Rate:    5,
		Comment: "master piece",
	}

	recorder := httptest.NewRecorder()
	url := fmt.Sprintf("%s/%s/rate", rateBaseUrl, content.ID.String())
	req := httptest.NewRequest(http.MethodPost, url, PrepareBody(input, &suite.Suite))
	suite.MockToken(DefaultUserEmail, DefaultPassword, req)

	suite.server.ServeHTTP(recorder, req)
	suite.Assert().Equal(http.StatusNoContent, recorder.Result().StatusCode)

	persistedContent, err := suite.contentRepository.Find(context.Background(), content.ID)
	suite.Assert().Nil(err)
	suite.Assert().Equal(input.Rate, (int)(*persistedContent.Rate))
	suite.Assert().Equal(input.Comment, *persistedContent.Comment)
}

func TestRateContentTestSuite(t *testing.T) {
	suite.Run(t, new(RateContentTestSuite))
}
