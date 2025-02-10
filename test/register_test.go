package test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/CarlosHenriqueDamasceno/wishtrack/cmd/api"
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
	server         *api.Api
}

func (suite *RegisterTestSuite) SetupTest() {
	conn, err := sql.Open("sqlite3", ":memory:")
	suite.Assert().Nil(err, "Fail to connect to database")

	suite.conn = conn
	suite.conn.Exec(
		`CREATE TABLE users (
			id varchar(255),
			name varchar(255),
			email varchar(255) UNIQUE,
			password varchar(255),
			created_at timestamp default (datetime('now','localtime'))
		)`,
	)

	suite.userRepository = user.NewDatabaseRepository(suite.conn)
	suite.userService = user.NewService(suite.userRepository)
	suite.server = api.NewApi(http.NewServeMux(), suite.userService)
}

func (suite *RegisterTestSuite) TestShouldFailToRegisterWithInvalidEmail() {
	input := user.RegisterInput{
		Name:     "Carlos",
		Email:    "carlos",
		Password: "senha",
	}

	body, err := json.Marshal(input)
	suite.Assert().Nil(err, "Body should be serialized")

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewReader(body))

	suite.server.ServeHTTP(recorder, req)
	suite.Assert().Equal(http.StatusUnprocessableEntity, recorder.Result().StatusCode, "Response status code should be 422 unprocessable entity")

	resp := &struct {
		Errors map[string][]string `json:"errors"`
	}{}

	err = json.NewDecoder(recorder.Result().Body).Decode(resp)
	suite.Assert().Nil(err, "fail to parse response")

	expectedErrors := map[string][]string{"email": {"field \"email\" must be a valid e-mail address"}}
	suite.Assert().Equal(expectedErrors, resp.Errors)

	res, err := suite.conn.Exec("select count(id) from users where email = ?", input.Email)
	suite.Assert().NoError(err, "fail to fetch database")
	rows, err := res.RowsAffected()
	suite.Assert().NoError(err, "fail to fetch rows")
	suite.Assert().Equal(int64(0), rows)
}

func (suite *RegisterTestSuite) TestShouldFailToRegisterWithEmailInUse() {
	input := user.RegisterInput{
		Name:     "Carlos",
		Email:    "carlos@teste.com",
		Password: "senha",
	}

	err := suite.userRepository.Create(context.Background(), &user.User{
		Name:     input.Name,
		Email:    user.Email(input.Email),
		Password: user.Password(input.Password),
	})
	if err != nil {
		suite.Failf("Fake user should be saved, but got error: %s", err.Error())
	}

	body, err := json.Marshal(input)
	suite.Assert().Nil(err, "Body should be serialized")

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewReader(body))

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

	res, err := suite.conn.Exec("select count(id) from users where email = ?", input.Email)
	suite.Assert().NoError(err, "fail to fetch database")

	rows, err := res.RowsAffected()
	suite.Assert().NoError(err, "fail to fetch rows")
	suite.Assert().Equal(int64(1), rows)
}

func (suite *RegisterTestSuite) TestShouldRegister() {
	input := user.RegisterInput{
		Name:     "Carlos",
		Email:    "carlos@wishtrack.com",
		Password: "password",
	}

	body, err := json.Marshal(input)
	suite.Assert().Nil(err, "Body should be serialized")

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/register", bytes.NewReader(body))

	suite.server.ServeHTTP(recorder, req)
	suite.Assert().Equal(http.StatusCreated, recorder.Result().StatusCode, "Response status code should be 201 created")

	responseBody := struct {
		ID        string    `json:"id"`
		Name      string    `json:"name"`
		Email     string    `json:"email"`
		CreatedAt time.Time `json:"created_at"`
	}{}

	err = json.NewDecoder(recorder.Result().Body).Decode(&responseBody)
	suite.Assert().Nil(err, "Fail to unmarshal response: %s")

	err = uuid.Validate(responseBody.ID)
	suite.Assert().Nil(err, "Result ID is invalid: %s")

	suite.Assert().NotZero(responseBody.CreatedAt, "Created at must be defined")

	savedUser, err := suite.userRepository.Find(context.Background(), uuid.MustParse(responseBody.ID))
	suite.Assert().Nil(err, "The user must be saved at this point: %s")

	suite.Assert().Equal(input.Name, savedUser.Name)
	suite.Assert().Equal(input.Email, string(savedUser.Email))
	suite.Assert().True(savedUser.Password.VerifyPassword(input.Password))
}

func TestRegisterTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterTestSuite))
}
