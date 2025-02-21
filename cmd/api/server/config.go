package server

import (
	"database/sql"
	"os"
	"time"
)

type Config struct {
	Address  string
	Auth     AuthConfig
	Database DatabaseConfig
}

type DatabaseConfig struct {
	Dsn  string
	Conn *sql.DB
}

type AuthConfig struct {
	Key string
	Iss string
	Aud string
	Exp time.Duration
}

func LoadEnv() *Config {
	expDuration, err := time.ParseDuration(os.Getenv("AUTH_EXP"))
	if err != nil {
		expDuration = time.Minute * 30
	}

	return &Config{
		Address: os.Getenv("ADDR"),
		Database: DatabaseConfig{
			Dsn: os.Getenv("DB_DSN"),
		},
		Auth: AuthConfig{
			Key: os.Getenv("AUTH_KEY"),
			Iss: os.Getenv("AUTH_ISS"),
			Aud: os.Getenv("AUTH_AUD"),
			Exp: expDuration,
		},
	}
}

func (c *Config) SetDatabaseConnection(conn *sql.DB) {
	c.Database.Conn = conn
}
