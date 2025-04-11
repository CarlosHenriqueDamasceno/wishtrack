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
	"github.com/stretchr/testify/suite"
)

const feedBaseUrl = "/api/v1/contents/feed"

type FeedTestSuite struct {
	LoggedRequestBaseSuite
}

func (suite *FeedTestSuite) SetupSuite() {
	suite.SetupDatabase()
}

func (suite *FeedTestSuite) SetupTest() {
	userRepository := user.NewDatabaseRepository(suite.conn)
	contentRepository := content.NewDatabaseRepository(suite.conn)

	auth := user.NewJwtAuthenticator(
		secret,
		aud,
		aud,
		time.Minute,
	)

	suite.UserService = user.NewService(userRepository, auth)
	suite.ContentService = content.NewService(contentRepository)

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

func (suite *FeedTestSuite) TearDownTest() {
	suite.ClearDatabase()
}

func (suite *FeedTestSuite) TearDownSuite() {
	suite.DestroyDatabase(context.Background())
}

func (suite *FeedTestSuite) mockContents() []*content.WriteDownOutput {
	ctx := context.Background()
	var output []*content.WriteDownOutput

	mostWished := &content.WriteDownInput{
		Name:      "The Lord of the Rings: The Return of the King",
		Category:  "book",
		Genres:    []string{"fantasy", "adventure"},
		Summary:   "The third movie from the series The Lord of The Rings",
		WishLevel: 5,
		UserID:    suite.User.ID,
	}

	lessWished := &content.WriteDownInput{
		Name:     "Saving private Ryan",
		Category: "movie",
		Genres:   []string{"war", "historical"},
		Summary: `Inspired by the books of Stephen E. Ambrose and accounts of multiple soldiers in a single
		family, such as the Niland brothers, being killed in action`,
		WishLevel: 2,
		UserID:    suite.User.ID,
	}

	out, err := suite.ContentService.WriteDown(ctx, lessWished)
	suite.Assert().Nil(err)
	output = append(output, out)

	out, err = suite.ContentService.WriteDown(ctx, mostWished)
	suite.Assert().Nil(err)
	output = append(output, out)

	return output
}

func (suite *FeedTestSuite) TestGetFeed() {
	contents := suite.mockContents()

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, feedBaseUrl, nil)
	suite.MockToken(DefaultUserEmail, DefaultPassword, req)

	suite.server.ServeHTTP(recorder, req)
	suite.Assert().Equal(http.StatusOK, recorder.Result().StatusCode)

	responseBody := []struct {
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

	suite.Assert().Equal(len(contents), len(responseBody))

	suite.Assert().Equal(contents[1].ID.String(), responseBody[0].ID)
}

func (suite *FeedTestSuite) TestFeedShouldNotIncludeRatedContents() {
	contents := suite.mockContents()

	input := &content.RateContentInput{
		ID:      contents[0].ID,
		UserID:  suite.User.ID,
		Rate:    5,
		Comment: "Absolute Cinema!!!!",
	}

	err := suite.ContentService.Rate(context.Background(), input)
	suite.Assert().Nil(err)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, feedBaseUrl, nil)
	suite.MockToken(DefaultUserEmail, DefaultPassword, req)

	suite.server.ServeHTTP(recorder, req)
	suite.Assert().Equal(http.StatusOK, recorder.Result().StatusCode)

	responseBody := []struct{}{}

	err = json.NewDecoder(recorder.Result().Body).Decode(&responseBody)
	suite.Assert().Nil(err, "Fail to unmarshal response: %s")

	suite.Assert().Equal(1, len(responseBody))
}

func TestFeedTestSuite(t *testing.T) {
	suite.Run(t, new(FeedTestSuite))
}
