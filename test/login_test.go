package test

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CarlosHenriqueDamasceno/wishtrack/cmd/api/server"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/user"
	"github.com/stretchr/testify/suite"
)

const baseUrl = "/api/v1/login"

type LoginTestSuite struct {
	suite.Suite
	conn           *sql.DB
	userRepository user.Repository
	userService    user.Service
	server         *server.Api
}

func (suite *LoginTestSuite) SetupTest() {
	conn, err := SetupDatabase()
	suite.Assert().Nil(err, "Fail to connect to database")

	suite.conn = conn

	suite.userRepository = user.NewDatabaseRepository(suite.conn)
	suite.userService = user.NewService(suite.userRepository)
	suite.server = server.NewApi(
		http.NewServeMux(),
		server.Config{},
		slog.Default(),
		suite.userService,
	)
}

func (suite *LoginTestSuite) TearDownTest() {
	suite.conn.Close()
}

func (suite *LoginTestSuite) mockUser(email, password string) *user.RegisterOutput {
	input := user.RegisterInput{
		Name:     "Carlos",
		Email:    email,
		Password: password,
	}

	out, err := suite.userService.Register(context.Background(), &input)
	suite.Assert().Nil(err, "Should register user")
	return out
}

//TODO: Test e-mail not found (find a way to treat not found errors in database repository)

func (suite *LoginTestSuite) TestShouldLogin() {
	email := "carlos@wishtrack.com"
	password := "12346578"
	suite.mockUser(email, password)

	input := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    email,
		Password: password,
	}

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, baseUrl, PrepareBody(input, &suite.Suite))

	suite.server.ServeHTTP(recorder, req)
	suite.Assert().Equal(http.StatusOK, recorder.Result().StatusCode)

	res := struct {
		Token string `json:"token"`
	}{}

	err := json.NewDecoder(recorder.Body).Decode(&res)
	suite.Assert().Nil(err, "Should be able to parse response")

	suite.Assert().NotZero(res.Token)
}

func TestLoginTestSuite(t *testing.T) {
	suite.Run(t, new(LoginTestSuite))
}
