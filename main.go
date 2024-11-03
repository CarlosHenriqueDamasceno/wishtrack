package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/CarlosHenriqueDamasceno/wishtrack/api"
	"github.com/CarlosHenriqueDamasceno/wishtrack/user"
	_ "github.com/mattn/go-sqlite3"
)

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
