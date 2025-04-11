package test

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/CarlosHenriqueDamasceno/wishtrack/cmd/api/server"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/content"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/user"
	"github.com/stretchr/testify/suite"
)

const listBaseUrl = "/api/v1/contents"

type contentResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Category  string    `json:"category"`
	Genres    []string  `json:"genres"`
	Summary   string    `json:"summary"`
	WishLevel int       `json:"wish_level"`
	Rate      *int      `json:"rate"`
	Comment   *string   `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type listResponse struct {
	Data  []contentResponse `json:"data"`
	Total uint64            `json:"total"`
	Page  uint64            `json:"page"`
	Limit uint64            `json:"limit"`
}

type ListTestSuite struct {
	LoggedRequestBaseSuite
}

func (suite *ListTestSuite) SetupSuite() {
	suite.SetupDatabase()
}

func (suite *ListTestSuite) SetupTest() {
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

func (suite *ListTestSuite) TearDownTest() {
	suite.ClearDatabase()
}

func (suite *ListTestSuite) TearDownSuite() {
	suite.DestroyDatabase(context.Background())
}

func (suite *ListTestSuite) mockContents() []*content.WriteDownOutput {
	ctx := context.Background()
	var output []*content.WriteDownOutput

	firstContent := &content.WriteDownInput{
		Name:      "The Lord of the Rings: The Return of the King",
		Category:  "book",
		Genres:    []string{"fantasy", "adventure"},
		Summary:   "The third part from the series The Lord of The Rings",
		WishLevel: 5,
		UserID:    suite.User.ID,
	}

	secondContent := &content.WriteDownInput{
		Name:     "Saving private Ryan",
		Category: "movie",
		Genres:   []string{"war", "historical"},
		Summary: `Inspired by the books of Stephen E. Ambrose and accounts of multiple soldiers in a single
		family, such as the Niland brothers, being killed in action`,
		WishLevel: 2,
		UserID:    suite.User.ID,
	}

	out, err := suite.ContentService.WriteDown(ctx, firstContent)
	suite.Assert().Nil(err)
	output = append(output, out)

	err = suite.ContentService.Rate(ctx, &content.RateContentInput{
		UserID:  firstContent.UserID,
		ID:      out.ID,
		Rate:    5,
		Comment: "Master piece",
	})
	suite.Assert().Nil(err)

	out, err = suite.ContentService.WriteDown(ctx, secondContent)
	suite.Assert().Nil(err)
	output = append(output, out)

	return output
}

func (suite *ListTestSuite) TestShouldGetPaginatedContents() {
	contents := suite.mockContents()

	const expectedLimit = uint64(1)
	const expectedPage = uint64(1)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, listBaseUrl, nil)
	suite.MockToken(DefaultUserEmail, DefaultPassword, req)

	q := req.URL.Query()
	q.Add("limit", strconv.FormatUint(expectedLimit, 10))
	req.URL.RawQuery = q.Encode()

	suite.server.ServeHTTP(recorder, req)
	suite.Assert().Equal(http.StatusOK, recorder.Result().StatusCode)

	var responseBody listResponse

	err := json.NewDecoder(recorder.Result().Body).Decode(&responseBody)
	suite.Assert().Nil(err, "Fail to unmarshal response: %s")

	suite.Assert().Equal(int(expectedLimit), len(responseBody.Data))
	suite.Assert().Equal(expectedLimit, responseBody.Limit)
	suite.Assert().Equal(expectedPage, responseBody.Page)
	suite.Assert().Equal(uint64(len(contents)), responseBody.Total)
}

func (suite *ListTestSuite) TestShouldGetSecondPageOfContents() {
	contents := suite.mockContents()

	const expectedLimit = uint64(1)
	const expectedPage = uint64(2)

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, listBaseUrl, nil)
	suite.MockToken(DefaultUserEmail, DefaultPassword, req)

	q := req.URL.Query()
	q.Add("limit", strconv.FormatUint(expectedLimit, 10))
	q.Add("page", strconv.FormatUint(expectedPage, 10))
	req.URL.RawQuery = q.Encode()

	suite.server.ServeHTTP(recorder, req)
	suite.Assert().Equal(http.StatusOK, recorder.Result().StatusCode)

	var responseBody listResponse

	err := json.NewDecoder(recorder.Result().Body).Decode(&responseBody)
	suite.Assert().Nil(err, "Fail to unmarshal response: %s")

	suite.Assert().Equal(int(expectedLimit), len(responseBody.Data))
	suite.Assert().Equal(expectedLimit, responseBody.Limit)
	suite.Assert().Equal(expectedPage, responseBody.Page)
	suite.Assert().Equal(uint64(len(contents)), responseBody.Total)
	suite.Assert().Equal(contents[1].ID.String(), responseBody.Data[0].ID)
}

func (suite *ListTestSuite) TestShouldFilterContentsByWatched() {
	suite.mockContents()

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, listBaseUrl, nil)
	suite.MockToken(DefaultUserEmail, DefaultPassword, req)

	q := req.URL.Query()
	q.Add("watched", "true")
	req.URL.RawQuery = q.Encode()

	suite.server.ServeHTTP(recorder, req)
	suite.Assert().Equal(http.StatusOK, recorder.Result().StatusCode)

	var responseBody listResponse

	err := json.NewDecoder(recorder.Result().Body).Decode(&responseBody)
	suite.Assert().Nil(err, "Fail to unmarshal response: %s")

	suite.Assert().Equal(1, len(responseBody.Data))
}

func (suite *ListTestSuite) TestShouldFilterContentsByCategory() {
	suite.mockContents()

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, listBaseUrl, nil)
	suite.MockToken(DefaultUserEmail, DefaultPassword, req)

	q := req.URL.Query()
	q.Add("search", "movie")
	req.URL.RawQuery = q.Encode()

	suite.server.ServeHTTP(recorder, req)
	suite.Assert().Equal(http.StatusOK, recorder.Result().StatusCode)

	var responseBody listResponse

	err := json.NewDecoder(recorder.Result().Body).Decode(&responseBody)
	suite.Assert().Nil(err, "Fail to unmarshal response: %s")

	suite.Assert().Equal(1, len(responseBody.Data))
}

func (suite *ListTestSuite) TestShouldFilterContentsByGenres() {
	suite.mockContents()

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, listBaseUrl, nil)
	suite.MockToken(DefaultUserEmail, DefaultPassword, req)

	q := req.URL.Query()
	q.Add("genres", "fantasy")
	req.URL.RawQuery = q.Encode()

	suite.server.ServeHTTP(recorder, req)
	suite.Assert().Equal(http.StatusOK, recorder.Result().StatusCode)

	var responseBody listResponse

	err := json.NewDecoder(recorder.Result().Body).Decode(&responseBody)
	suite.Assert().Nil(err, "Fail to unmarshal response: %s")

	suite.Assert().Equal(1, len(responseBody.Data))
}

func (suite *ListTestSuite) TestShouldFilterContentsByName() {
	suite.mockContents()

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, listBaseUrl, nil)
	suite.MockToken(DefaultUserEmail, DefaultPassword, req)

	q := req.URL.Query()
	q.Add("search", "lord")
	req.URL.RawQuery = q.Encode()

	suite.server.ServeHTTP(recorder, req)
	suite.Assert().Equal(http.StatusOK, recorder.Result().StatusCode)

	var responseBody listResponse

	err := json.NewDecoder(recorder.Result().Body).Decode(&responseBody)
	suite.Assert().Nil(err, "Fail to unmarshal response: %s")

	suite.Assert().Equal(1, len(responseBody.Data))
}

func (suite *ListTestSuite) TestShouldFilterContentsBySummary() {
	suite.mockContents()

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, listBaseUrl, nil)
	suite.MockToken(DefaultUserEmail, DefaultPassword, req)

	q := req.URL.Query()
	q.Add("search", "The third part")
	req.URL.RawQuery = q.Encode()

	suite.server.ServeHTTP(recorder, req)
	suite.Assert().Equal(http.StatusOK, recorder.Result().StatusCode)

	var responseBody listResponse

	err := json.NewDecoder(recorder.Result().Body).Decode(&responseBody)
	suite.Assert().Nil(err, "Fail to unmarshal response: %s")

	suite.Assert().Equal(1, len(responseBody.Data))
}

func (suite *ListTestSuite) TestShouldFilterContentsByMinWishLevel() {
	suite.mockContents()

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, listBaseUrl, nil)
	suite.MockToken(DefaultUserEmail, DefaultPassword, req)

	q := req.URL.Query()
	q.Add("wishLevel", "3")
	req.URL.RawQuery = q.Encode()

	suite.server.ServeHTTP(recorder, req)
	suite.Assert().Equal(http.StatusOK, recorder.Result().StatusCode)

	var responseBody listResponse

	err := json.NewDecoder(recorder.Result().Body).Decode(&responseBody)
	suite.Assert().Nil(err, "Fail to unmarshal response: %s")

	suite.Assert().Equal(1, len(responseBody.Data))
}

func TestListTestSuite(t *testing.T) {
	suite.Run(t, new(ListTestSuite))
}
