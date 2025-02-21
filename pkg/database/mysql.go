package database

import (
	"errors"

	"github.com/go-sql-driver/mysql"
)

func WrapMysqlError(err error) error {
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		return errors.New(mysqlErr.Message)
	}
	return err
}
