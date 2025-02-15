package server

import "database/sql"

type Config struct {
	Address string
	Db      *sql.DB
}
