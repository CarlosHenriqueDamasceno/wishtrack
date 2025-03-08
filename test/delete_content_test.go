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

const deleteBaseUrl = "/api/v1/contents"

type DeleteContentTestSuite struct {
	LoggedRequestBaseSuite
}

func (suite *DeleteContentTestSuite) SetupTest() {
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

func (suite *DeleteContentTestSuite) TearDownTest() {
	suite.conn.Close()
}

func (suite *DeleteContentTestSuite) mockContent() *content.WriteDownOutput {
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

func TestDeleteContentTestSuite(t *testing.T) {
	suite.Run(t, new(DeleteContentTestSuite))
}

func (suite *DeleteContentTestSuite) TestDeleteAContent() {
	ct := suite.mockContent()

	recorder := httptest.NewRecorder()
	url := fmt.Sprintf("%s/%s", deleteBaseUrl, ct.ID.String())
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	suite.mockToken(DefaultUserEmail, DefaultPassword, req)

	suite.server.ServeHTTP(recorder, req)
	suite.Assert().Equal(http.StatusNoContent, recorder.Result().StatusCode)

	_, err := suite.contentService.Find(context.Background(), ct.ID, suite.user.ID)
	suite.Assert().ErrorIs(err, content.ErrContentNotFound)
}
