package test

import (
	"database/sql"

	"github.com/stretchr/testify/suite"
)

func AssertDatabaseCount(conn *sql.DB, suite *suite.Suite, expectedCount int, table string, column string) {
	var count int
	err := conn.QueryRow("SELECT COUNT(" + column + ") FROM " + table).Scan(&count)
	suite.Nil(err, "Fail to count rows")
	suite.Equal(expectedCount, count)
}
