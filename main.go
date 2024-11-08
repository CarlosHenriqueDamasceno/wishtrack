package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/CarlosHenriqueDamasceno/wishtrack/api"
	"github.com/CarlosHenriqueDamasceno/wishtrack/user"
	_ "github.com/mattn/go-sqlite3"
)

// @title           WishTrack API
// @version         1.0
// @description     A app to manage all the things you want see or have saw.

// @contact.name   API Support
// @contact.email  carlos@wishtrack.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/

func main() {
	conn, err := sql.Open("sqlite3", "etc/database.sqlite")
	if err != nil {
		log.Fatalf("Database connection fail: %s", err.Error())
	}
	userRepository := user.NewDatabaseRepository(conn)
	server := api.NewApiServer(http.NewServeMux(), userRepository)
	err = http.ListenAndServe(":8080", server)
	if err != nil {
		log.Fatalf("Server error: %s", err.Error())
	}
}
