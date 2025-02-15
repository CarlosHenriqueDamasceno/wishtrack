package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"github.com/CarlosHenriqueDamasceno/wishtrack/cmd/api/server"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/user"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func run(conf server.Config, logger *slog.Logger) error {
	userRepository := user.NewDatabaseRepository(conf.Db)
	userService := user.NewService(userRepository)
	api := server.NewApi(http.NewServeMux(), conf, logger, userService)
	server := http.Server{
		Addr:    conf.Address,
		Handler: api,
	}
	return server.ListenAndServe()
}

// @title						WishTrack API
// @version					1.0
// @description				A app to manage all the things you want see or have saw.
// @contact.name				API Support
// @contact.email				carlos@wishtrack.io
// @license.name				Apache 2.0
// @license.url				http://www.apache.org/licenses/LICENSE-2.0.html
// @host						localhost:8080
// @BasePath					/api/v1
// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file")
	}

	conn, err := sql.Open("mysql", os.Getenv("DB_DSN"))
	if err != nil {
		logger.Error("Database connection fail", "error", err.Error())
	}
	defer conn.Close()

	logger.Info("Database connection established.")

	config := server.Config{
		Address: ":8080",
		Db:      conn,
	}

	if err := run(config, logger); err != nil {
		logger.Error("Server error", "error", err.Error())
	}
}
