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
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/user"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/suite"
)

const (
	loginBaseUrl    = "/api/v1/users/login"
	secret          = "ff1feb0beced46fc8ae6f662f7846eb2"
	aud             = "wishtrack"
	tokenExpiration = time.Minute
)

type LoginTestSuite struct {
	DatabaseSuite
	userRepository user.Repository
	userService    user.Service
	server         *server.Api
	auth           user.Authenticator
}

func (suite *LoginTestSuite) SetupSuite() {
	suite.SetupDatabase()
}

func (suite *LoginTestSuite) SetupTest() {
	suite.userRepository = user.NewDatabaseRepository(suite.conn)

	suite.auth = user.NewJwtAuthenticator(
		secret,
		aud,
		aud,
		time.Minute,
	)

	suite.userService = user.NewService(suite.userRepository, suite.auth)

	suite.server = server.NewApi(
		http.NewServeMux(),
		&server.Config{},
		slog.Default(),
		suite.userService,
		nil,
	)
}

func (suite *LoginTestSuite) TearDownTest() {
	suite.ClearDatabase()
}

func (suite *LoginTestSuite) TearDownSuite() {
	suite.DestroyDatabase(context.Background())
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

func (suite *LoginTestSuite) TestShouldFailToLoginWithInvalidEmail() {
	input := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    "carlos@wishtrack.com",
		Password: "12334545654",
	}

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, loginBaseUrl, PrepareBody(input, &suite.Suite))

	suite.server.ServeHTTP(recorder, req)
	suite.Assert().Equal(http.StatusUnauthorized, recorder.Result().StatusCode)

	resp := &struct {
		Error string `json:"error"`
	}{}

	err := json.NewDecoder(recorder.Result().Body).Decode(resp)
	suite.Assert().Nil(err, "fail to parse response")

	expectedMessage := "the provided e-mail or password are incorrect"
	suite.Assert().Equal(expectedMessage, resp.Error)
}

func (suite *LoginTestSuite) TestShouldFailToLoginWithWrongPassword() {
	email := "carlos@wishtrack.com"
	password := "12346578"
	suite.mockUser(email, password)

	input := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    email,
		Password: "2342342342342",
	}

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, loginBaseUrl, PrepareBody(input, &suite.Suite))

	suite.server.ServeHTTP(recorder, req)
	suite.Assert().Equal(http.StatusUnauthorized, recorder.Result().StatusCode)

	resp := &struct {
		Error string `json:"error"`
	}{}

	err := json.NewDecoder(recorder.Result().Body).Decode(resp)
	suite.Assert().Nil(err, "fail to parse response")

	expectedMessage := "the provided e-mail or password are incorrect"
	suite.Assert().Equal(expectedMessage, resp.Error)
}

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
	req := httptest.NewRequest(http.MethodPost, loginBaseUrl, PrepareBody(input, &suite.Suite))

	suite.server.ServeHTTP(recorder, req)
	suite.Assert().Equal(http.StatusOK, recorder.Result().StatusCode)

	res := struct {
		Token string `json:"token"`
	}{}

	err := json.NewDecoder(recorder.Body).Decode(&res)
	suite.Assert().Nil(err, "Should be able to parse response")

	suite.Assert().NotZero(res.Token)

	jwtToken := parseToken(&suite.Suite, res.Token)

	exp, err := jwtToken.Claims.GetExpirationTime()
	suite.Assert().Nil(err)
	suite.Assert().Equal(
		exp.Time.Format(time.DateTime),
		time.Now().Add(tokenExpiration).Format(time.DateTime),
	)
}

func parseToken(suite *suite.Suite, token string) *jwt.Token {
	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}

		return []byte(secret), nil
	},
		jwt.WithExpirationRequired(),
		jwt.WithAudience(aud),
		jwt.WithIssuer(aud),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
	)

	suite.Assert().Nil(err)
	return jwtToken
}

func TestLoginTestSuite(t *testing.T) {
	suite.Run(t, new(LoginTestSuite))
}
