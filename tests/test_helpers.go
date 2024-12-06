// tests/test_helpers.go
package tests

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func SetupTestDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	return db, err
}
