package repository

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/shuvo-paul/sitemonitor/internal/database/migrations"
	"github.com/shuvo-paul/sitemonitor/tests"
)

var db *sql.DB

func TestMain(m *testing.M) {
	var err error
	db, err = tests.SetupTestDB()
	if err != nil {
		panic(fmt.Sprintf("Error setting up test database: %v", err))
	}
	migrations.SetupMigration(db)

	// Run the tests
	code := m.Run()

	// Cleanup (optional)
	if db != nil {
		db.Close()
	}

	// Exit with the test result code
	os.Exit(code)
}
