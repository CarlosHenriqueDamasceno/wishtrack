package test

import (
	"context"
	"database/sql"
	"net/http"
	"path/filepath"
	"time"

	"github.com/CarlosHenriqueDamasceno/wishtrack/cmd/api/server"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/content"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/user"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/golang-migrate/migrate/v4"
	migratePg "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/lib/pq"
)

const (
	DefaultUserEmail = "carlos@wishtrack.com"
	DefaultPassword  = "12345678"

	DatabaseImage = "postgres:17.4-alpine"
	DatabaseUser  = "postgres"
	DatabasePass  = "postgres"
	DatabaseName  = "wishtrack"
)

var DatabasePath = filepath.Join("..", "testdata", "init-db.sql")

type DatabaseSuite struct {
	suite.Suite
	conn      *sql.DB
	container *postgres.PostgresContainer
}

func (suite *DatabaseSuite) SetupDatabase() {
	ctx := context.Background()

	suite.createDatabase(ctx)
	suite.conn = suite.startConnection()
	suite.migrate()
}

func (suite *DatabaseSuite) createDatabase(ctx context.Context) {
	timeout := 3 * time.Minute
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	pgContainer, err := postgres.Run(ctx,
		DatabaseImage,
		postgres.WithInitScripts(),
		postgres.WithDatabase(DatabaseName),
		postgres.WithUsername(DatabaseUser),
		postgres.WithPassword(DatabasePass),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)

	suite.Assert().Nil(err)

	suite.container = pgContainer
}

func (suite *DatabaseSuite) startConnection() *sql.DB {
	conn, err := sql.Open("postgres", suite.dsn(context.Background()))
	suite.Assert().Nil(err, "Should connect to database")
	return conn
}

func (suite *DatabaseSuite) destroyDatabase(ctx context.Context) {
	if err := suite.container.Terminate(ctx); err != nil {
		suite.Assert().Nil(err)
	}
}

func (suite *DatabaseSuite) migrate() {
	driver, err := migratePg.WithInstance(suite.conn, &migratePg.Config{})
	suite.Assert().Nil(err)

	m, err := migrate.NewWithDatabaseInstance(
		"file://../etc/migrations",
		"wishtrack",
		driver,
	)
	suite.Assert().Nil(err)

	err = m.Up()
	suite.Assert().Nil(err)
}

func (suite *DatabaseSuite) dsn(ctx context.Context) string {
	dsn, err := suite.container.ConnectionString(ctx, "sslmode=disable", "TimeZone=UTC")
	suite.Assert().Nil(err)
	return dsn
}

type LoggedRequestBaseSuite struct {
	DatabaseSuite
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
