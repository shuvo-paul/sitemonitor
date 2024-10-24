package config

import (
	"testing"

	"github.com/joho/godotenv" // Import the godotenv package
	_ "github.com/lib/pq"      // Import the database driver
)

func TestInitDatabase(t *testing.T) {
	// Load environment variables from .env.test file
	if err := godotenv.Load("../.env.test"); err != nil {
		t.Fatalf("Error loading .env.test file: %v", err)
	}

	// Call the function to test
	InitDatabase()

	// Check if the DB variable is not nil
	if DB == nil {
		t.Fatal("Expected DB to be initialized, but it was nil")
	}

	// Check if the database connection is valid
	if err := DB.Ping(); err != nil {
		t.Fatalf("Expected to connect to the database, but got error: %v", err)
	}

	// Clean up
	DB.Close()
}
