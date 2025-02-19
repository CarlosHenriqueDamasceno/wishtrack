package test

import (
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func SetupDatabase() (*sql.DB, error) {
	conn, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	driver, err := sqlite3.WithInstance(conn, &sqlite3.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../etc/migrations",
		"sqlite3",
		driver,
	)
	if err != nil {
		return nil, err
	}

	m.Up()

	return conn, nil
}
