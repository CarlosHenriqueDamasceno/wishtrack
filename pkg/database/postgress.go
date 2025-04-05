package database

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

var ErrRecordNotFound = errors.New("record not found")

func ParseDatabaseError(err error) error {
	if postgressError, ok := err.(*pq.Error); ok {
		return errors.New(postgressError.Message)
	}

	if errors.Is(err, sql.ErrNoRows) {
		return ErrRecordNotFound
	}

	return err
}

func New(addr string, maxOpenConns, maxIdleConns int, maxIdleTime time.Duration) (*sql.DB, error) {
	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxIdleTime(maxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
