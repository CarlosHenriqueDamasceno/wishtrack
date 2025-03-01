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
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

const editBaseUrl = "/api/v1/contents"

type ContentEditTestSuite struct {
	LoggedRequestBaseSuite
}

func (suite *ContentEditTestSuite) SetupTest() {
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

func (suite *ContentEditTestSuite) TearDownTest() {
	suite.conn.Close()
}

func (suite *ContentEditTestSuite) mockContent() *content.WriteDownOutput {
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

func (suite *ContentEditTestSuite) TestContentNotFound() {
	input := struct {
		Name      string   `json:"name"`
		Category  string   `json:"category"`
		Genres    []string `json:"genres"`
		Summary   string   `json:"summary"`
		WishLevel int      `json:"wish_level"`
	}{
		Name:      "The Lord of the Rings: The Return of the King Edited",
		Category:  "book",
		Genres:    []string{"fantasy"},
		Summary:   "",
		WishLevel: 5,
	}

	recorder := httptest.NewRecorder()
	url := fmt.Sprintf("%s/%s", editBaseUrl, uuid.New())
	req := httptest.NewRequest(http.MethodPut, url, PrepareBody(input, &suite.Suite))
	suite.mockToken(DefaultUserEmail, DefaultPassword, req)

	suite.server.ServeHTTP(recorder, req)
	suite.Assert().Equal(http.StatusNotFound, recorder.Result().StatusCode)
}

func (suite *ContentEditTestSuite) TestShouldEditAContent() {
	content := suite.mockContent()

	input := struct {
		Name      string   `json:"name"`
		Category  string   `json:"category"`
		Genres    []string `json:"genres"`
		Summary   string   `json:"summary"`
		WishLevel int      `json:"wish_level"`
	}{
		Name:      "The Lord of the Rings: The Return of the King Edited",
		Category:  "book",
		Genres:    content.Genres,
		Summary:   content.Summary,
		WishLevel: int(content.WishLevel),
	}

	recorder := httptest.NewRecorder()
	url := fmt.Sprintf("%s/%s", editBaseUrl, content.ID.String())
	req := httptest.NewRequest(http.MethodPut, url, PrepareBody(input, &suite.Suite))
	suite.mockToken(DefaultUserEmail, DefaultPassword, req)

	suite.server.ServeHTTP(recorder, req)
	suite.Assert().Equal(http.StatusOK, recorder.Result().StatusCode)

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

	suite.Assert().Equal(content.CreatedAt, responseBody.CreatedAt)
	suite.Assert().NotEqual(content.UpdatedAt, responseBody.UpdatedAt)
	suite.Assert().Equal(input.Genres, responseBody.Genres)
	suite.Assert().Equal(input.Summary, responseBody.Summary)
	suite.Assert().Equal(input.WishLevel, responseBody.WishLevel)
}

func TestEditContentTestSuite(t *testing.T) {
	suite.Run(t, new(ContentEditTestSuite))
}
