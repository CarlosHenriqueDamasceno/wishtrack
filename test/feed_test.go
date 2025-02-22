package test

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/CarlosHenriqueDamasceno/wishtrack/cmd/api/server"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/content"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/user"
)

const feedBaseUrl = "/api/v1/contents/feed"

type FeedTestSuite struct {
	LoggedRequestBaseSuite
}

func (suite *FeedTestSuite) SetupTest() {
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
}

func (suite *FeedTestSuite) TearDownTest() {
	suite.conn.Close()
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
	}

	lessWished := &content.WriteDownInput{
		Name:     "Saving private Ryan",
		Category: "movie",
		Genres:   []string{"war", "historical"},
		Summary: `Inspired by the books of Stephen E. Ambrose and accounts of multiple soldiers in a single
		family, such as the Niland brothers, being killed in action`,
		WishLevel: 2,
	}

	out, err := suite.contentService.WriteDown(ctx, lessWished)
	suite.Assert().Nil(err)
	output = append(output, out)

	out, err = suite.contentService.WriteDown(ctx, mostWished)
	suite.Assert().Nil(err)
	output = append(output, out)

	return output
}

func (suite *FeedTestSuite) TestGetFeed() {
	contents := suite.mockContents()

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, feedBaseUrl, nil)
	suite.mockToken(req)

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

	suite.Assert().Equal(contents[1].ID, responseBody[0].ID)
}
