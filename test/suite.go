package test

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/CarlosHenriqueDamasceno/wishtrack/cmd/api/server"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/content"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/user"
	"github.com/stretchr/testify/suite"
)

const (
	DefaultUserEmail = "carlos@wishtrack.com"
	DefaultPassword  = "12345678"
)

type LoggedRequestBaseSuite struct {
	suite.Suite
	conn           *sql.DB
	server         *server.Api
	userService    user.Service
	contentService content.Service
	user           *user.RegisterOutput
}

func (suite *LoggedRequestBaseSuite) mockUser(email, password string) {

	userInput := &user.RegisterInput{
		Name:     "Carlos",
		Email:    email,
		Password: password,
	}

	ctx := context.Background()
	registeredUser, err := suite.userService.Register(ctx, userInput)
	suite.Assert().Nil(err, "Should register user")
	suite.user = registeredUser
}

func (suite *LoggedRequestBaseSuite) mockToken(email, password string, req *http.Request) {
	loginInput := &user.LoginInput{
		Email:    email,
		Password: password,
	}

	token, err := suite.userService.Login(context.Background(), loginInput)
	suite.Assert().Nil(err, "Should login")

	req.Header.Add("Authorization", "Bearer "+token.Token)
}
