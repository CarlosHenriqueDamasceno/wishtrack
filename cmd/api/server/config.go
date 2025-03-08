package server

import (
	"database/sql"
	"time"

	"github.com/CarlosHenriqueDamasceno/wishtrack/pkg/env"
	"github.com/google/uuid"
)

type Config struct {
	Address  string
	Env      string
	Auth     AuthConfig
	Database DatabaseConfig
}

type DatabaseConfig struct {
	Dsn          string
	Conn         *sql.DB
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
}

type AuthConfig struct {
	Key string
	Iss string
	Aud string
	Exp time.Duration
}

func LoadEnv() *Config {
	return &Config{
		Address: env.GetString("ADDR", ":8080"),
		Env:     env.GetString("ENV", "local"),
		Database: DatabaseConfig{
			Dsn:          env.GetString("DB_DSN", "sqlite"),
			MaxOpenConns: env.GetInt("DB_MAX_CONN", 30),
			MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONN", 30),
			MaxIdleTime:  env.GetDuration("DB_MAX_IDLE_TIME", time.Minute*15),
		},
		Auth: AuthConfig{
			Key: env.GetString("AUTH_KEY", uuid.NewString()),
			Iss: env.GetString("AUTH_ISS", "wishtrack"),
			Aud: env.GetString("AUTH_AUD", "wishtrack"),
			Exp: env.GetDuration("AUTH_EXP", time.Minute*30),
		},
	}
}

func (c *Config) SetDatabaseConnection(conn *sql.DB) {
	c.Database.Conn = conn
}
