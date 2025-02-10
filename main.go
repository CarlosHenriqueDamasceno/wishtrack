package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/CarlosHenriqueDamasceno/wishtrack/cmd/api"
	"github.com/CarlosHenriqueDamasceno/wishtrack/internal/user"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	Address string
	Db      *sql.DB
}

func run(conf Config) error {
	userRepository := user.NewDatabaseRepository(conf.Db)
	userService := user.NewService(userRepository)
	api := api.NewApi(http.NewServeMux(), userService)
	server := http.Server{
		Addr:    conf.Address,
		Handler: api,
	}
	return server.ListenAndServe()
}

// @title                     WishTrack API
// @version                   1.0
// @description               A app to manage all the things you want see or have saw.
// @contact.name              API Support
// @contact.email             carlos@wishtrack.io
// @license.name              Apache 2.0
// @license.url               http://www.apache.org/licenses/LICENSE-2.0.html
// @host                      localhost:8080
// @BasePath                  /api/v1
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	conn, err := sql.Open("sqlite3", "etc/database.sqlite")
	if err != nil {
		log.Fatalf("Database connection fail: %s", err.Error())
	}
	defer conn.Close()

	if err := run(Config{
		Address: ":8080",
		Db:      conn,
	}); err != nil {
		log.Fatalf("Server error: %s", err.Error())
	}
}
