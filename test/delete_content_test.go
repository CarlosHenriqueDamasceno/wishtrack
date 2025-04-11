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

func (suite *DeleteContentTestSuite) SetupSuite() {
	suite.SetupDatabase()
}

func (suite *DeleteContentTestSuite) SetupTest() {
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

func (suite *DeleteContentTestSuite) TearDownTest() {
	suite.ClearDatabase()
}

func (suite *DeleteContentTestSuite) TearDownSuite() {
	suite.DestroyDatabase(context.Background())
}
func (suite *DeleteContentTestSuite) mockContent() *content.WriteDownOutput {
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

func TestDeleteContentTestSuite(t *testing.T) {
	suite.Run(t, new(DeleteContentTestSuite))
}

func (suite *DeleteContentTestSuite) TestDeleteAContent() {
	ct := suite.mockContent()

	recorder := httptest.NewRecorder()
	url := fmt.Sprintf("%s/%s", deleteBaseUrl, ct.ID.String())
	req := httptest.NewRequest(http.MethodDelete, url, nil)
	suite.MockToken(DefaultUserEmail, DefaultPassword, req)

	suite.server.ServeHTTP(recorder, req)
	suite.Assert().Equal(http.StatusNoContent, recorder.Result().StatusCode)

	_, err := suite.ContentService.Find(context.Background(), ct.ID, suite.User.ID)
	suite.Assert().ErrorIs(err, content.ErrContentNotFound)
}
