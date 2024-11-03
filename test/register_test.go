package test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/CarlosHenriqueDamasceno/wishtrack/api"
	"github.com/CarlosHenriqueDamasceno/wishtrack/user"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type RegisterTestSuite struct {
	suite.Suite
	userRepository user.Repository
	server         *api.ApiServer
}

func (suite *RegisterTestSuite) SetupTest() {
	connection, err := sql.Open("sqlite3", ":memory:")
	suite.Assert().Nil(err, "Fail to connect to database")

	connection.Exec("CREATE TABLE users (id varchar(255), name varchar(255), email varchar(255), password varchar(255), created_at timestamp)")

	suite.userRepository = user.NewDatabaseRepository(connection)
	suite.server = api.NewApiServer(http.NewServeMux(), suite.userRepository)
}

func (suite *RegisterTestSuite) TestShouldRegister() {
	input := struct {
		Name     string
		Email    string
		Password string
	}{
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

	savedUser, err := suite.userRepository.Find(uuid.MustParse(responseBody.ID))
	suite.Assert().Nil(err, "The user must be saved at this point: %s")

	suite.Assert().Equal(input.Name, savedUser.Name)
	suite.Assert().Equal(input.Email, savedUser.Email)
	suite.Assert().True(savedUser.VerifyPassword(input.Password))
}

func TestRegisterTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterTestSuite))
}
