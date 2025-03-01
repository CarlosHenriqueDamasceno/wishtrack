package database

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
)

var ErrRecordNotFound = errors.New("record not found")

func ParseDatabaseError(err error) error {
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		return errors.New(mysqlErr.Message)
	}

	if errors.Is(err, sql.ErrNoRows) {
		return ErrRecordNotFound
	}

	return err
}
