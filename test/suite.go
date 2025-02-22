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

type LoggedRequestBaseSuite struct {
	suite.Suite
	conn           *sql.DB
	server         *server.Api
	userService    user.Service
	contentService content.Service
}

func (suite *LoggedRequestBaseSuite) mockToken(req *http.Request) {
	userInput := &user.RegisterInput{
		Name:     "Carlos",
		Email:    "carlos@wishtrack.com",
		Password: "12345678",
	}

	loginInput := &user.LoginInput{
		Email:    userInput.Email,
		Password: userInput.Password,
	}

	ctx := context.Background()
	_, err := suite.userService.Register(ctx, userInput)
	suite.Assert().Nil(err, "Should register user")

	token, err := suite.userService.Login(ctx, loginInput)
	suite.Assert().Nil(err, "Should login")

	req.Header.Add("Authorization", "Bearer "+token.Token)
}
