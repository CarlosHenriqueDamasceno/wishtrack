package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"github.com/CarlosHenriqueDamasceno/wishtrack/cmd/api/server"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/content"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/user"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func run(conf *server.Config, logger *slog.Logger) error {
	userRepository := user.NewDatabaseRepository(conf.Database.Conn)
	contentRepository := content.NewDatabaseRepository(conf.Database.Conn)
	jwtAuth := user.NewJwtAuthenticator(conf.Auth.Key, conf.Auth.Iss, conf.Auth.Aud, conf.Auth.Exp)
	userService := user.NewService(userRepository, jwtAuth)
	contentService := content.NewService(contentRepository)
	api := server.NewApi(http.NewServeMux(), conf, logger, userService, contentService)
	server := http.Server{
		Addr:    conf.Address,
		Handler: api,
	}
	return server.ListenAndServe()
}

// @title						WishTrack API
// @version					1.0
// @description				A app to manage all the things you wanna see (or had saw).
// @contact.name				API Support
// @contact.email				carlos@wishtrack.io
// @license.name				Apache 2.0
// @license.url				http://www.apache.org/licenses/LICENSE-2.0.html
// @host						localhost:8080
// @BasePath					/api/v1
// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file")
	}

	config := server.LoadEnv()

	conn, err := sql.Open("mysql", config.Database.Dsn)
	if err != nil {
		logger.Error("Database connection fail", "error", err.Error())
	}
	defer conn.Close()

	config.SetDatabaseConnection(conn)

	logger.Info("Database connection established.")

	if err := run(config, logger); err != nil {
		logger.Error("Server error", "error", err.Error())
	}
}
