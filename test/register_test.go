package test

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/CarlosHenriqueDamasceno/wishtrack/cmd/api/server"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/user"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type RegisterTestSuite struct {
	suite.Suite
	conn           *sql.DB
	userRepository user.Repository
	userService    user.Service
	server         *server.Api
}

func (suite *RegisterTestSuite) SetupTest() {
	conn, err := SetupDatabase()
	suite.Assert().Nil(err, "Fail to connect to database")

	suite.conn = conn

	suite.userRepository = user.NewDatabaseRepository(suite.conn)
	suite.userService = user.NewService(suite.userRepository, nil)
	suite.server = server.NewApi(
		http.NewServeMux(),
		&server.Config{},
		slog.Default(),
		suite.userService,
	)
}

func (suite *RegisterTestSuite) TearDownTest() {
	suite.conn.Close()
}

func (suite *RegisterTestSuite) TestShouldFailToRegisterWithInvalidEmail() {
	input := user.RegisterInput{
		Name:     "Carlos",
		Email:    "carlos",
		Password: "12345678",
	}

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/register", PrepareBody(input, &suite.Suite))

	suite.server.ServeHTTP(recorder, req)

	suite.Assert().Equal(
		http.StatusUnprocessableEntity,
		recorder.Result().StatusCode,
		"Response status code should be 422 unprocessable entity",
	)

	resp := &struct {
		Errors map[string][]string `json:"errors"`
	}{}

	err := json.NewDecoder(recorder.Result().Body).Decode(resp)
	suite.Assert().Nil(err, "fail to parse response")

	expectedErrors := map[string][]string{"email": {"field \"email\" must be a valid e-mail address"}}
	suite.Assert().Equal(expectedErrors, resp.Errors)

	AssertDatabaseCount(suite.conn, &suite.Suite, 0, "users", "id")
}

func (suite *RegisterTestSuite) TestShouldFailToRegisterWithEmailInUse() {
	input := user.RegisterInput{
		Name:     "Carlos",
		Email:    "carlos@teste.com",
		Password: "12345678",
	}

	err := suite.userRepository.Create(context.Background(), &user.User{
		Name:  input.Name,
		Email: user.Email(input.Email),
	})
	if err != nil {
		suite.Failf("Fake user should be saved, but got error: %s", err.Error())
	}

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/register", PrepareBody(input, &suite.Suite))

	suite.server.ServeHTTP(recorder, req)
	suite.Assert().Equal(
		http.StatusUnprocessableEntity,
		recorder.Result().StatusCode,
		"Response status code should be 422 unprocessable entity",
	)

	resp := &struct {
		Errors map[string][]string `json:"errors"`
	}{}

	err = json.NewDecoder(recorder.Result().Body).Decode(resp)
	suite.Assert().Nil(err, "fail to parse response")

	suite.Assert().Equal(http.StatusUnprocessableEntity, recorder.Result().StatusCode)

	expectedErrors := map[string][]string{"email": {"e-mail already in use"}}
	suite.Assert().Equal(expectedErrors, resp.Errors)

	AssertDatabaseCount(suite.conn, &suite.Suite, 1, "users", "id")
}

func (suite *RegisterTestSuite) TestShouldFailToRegisterWithInvalidPassword() {
	input := user.RegisterInput{
		Name:     "Carlos Henrique",
		Email:    "carlos@wishtrack.com",
		Password: "senha",
	}

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/register", PrepareBody(input, &suite.Suite))

	suite.server.ServeHTTP(recorder, req)

	suite.Assert().Equal(
		http.StatusUnprocessableEntity,
		recorder.Result().StatusCode,
		"Response status code should be 422 unprocessable entity",
	)

	resp := &struct {
		Errors map[string][]string `json:"errors"`
	}{}

	err := json.NewDecoder(recorder.Result().Body).Decode(resp)
	suite.Assert().Nil(err, "fail to parse response")

	expectedErrors := map[string][]string{"password": {"field \"password\" must be at least 8 characters long"}}
	suite.Assert().Equal(expectedErrors, resp.Errors)

	AssertDatabaseCount(suite.conn, &suite.Suite, 0, "users", "id")
}

func (suite *RegisterTestSuite) TestShouldRegister() {
	input := user.RegisterInput{
		Name:     "Carlos",
		Email:    "carlos@wishtrack.com",
		Password: "password",
	}

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/register", PrepareBody(input, &suite.Suite))

	suite.server.ServeHTTP(recorder, req)
	suite.Assert().Equal(http.StatusCreated, recorder.Result().StatusCode, "Response status code should be 201 created")

	responseBody := struct {
		ID        string    `json:"id"`
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
		UpdateAt  time.Time `json:"updated_at"`
	}{}

	err := json.NewDecoder(recorder.Result().Body).Decode(&responseBody)
	suite.Assert().Nil(err, "Fail to unmarshal response: %s")

	err = uuid.Validate(responseBody.ID)
	suite.Assert().Nil(err, "Result ID is invalid: %s")

	suite.Assert().NotZero(responseBody.CreatedAt, "Created at must be defined")
	suite.Assert().NotZero(responseBody.UpdateAt, "Updated at must be defined")

	savedUser, err := suite.userRepository.Find(context.Background(), uuid.MustParse(responseBody.ID))
	suite.Assert().Nil(err, "The user must be saved at this point: %s")

	suite.Assert().Equal(input.Name, savedUser.Name)
	suite.Assert().Equal(input.Email, string(savedUser.Email))
	suite.Assert().True(savedUser.VerifyPassword(input.Password))
}

func TestRegisterTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterTestSuite))
}
