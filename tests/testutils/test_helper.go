package testutils

import (
	"log"
	"os"
	"testing"

	"github.com/EmelinDanila/task-manager-api/config"

	"github.com/EmelinDanila/task-manager-api/models"
)

// SetupTestDB initializes the test database
func SetupTestDB(t *testing.T) config.Database {
	// Load environment variables
	os.Setenv("GO_ENV", "test")
	config.LoadEnvVars()

	// Connect to the database
	db, err := config.ConnectDatabase()
	if err != nil {
		t.Fatalf("Could not connect to the database: %v", err)
	}

	// Migrate the Task model
	err = db.GetDB().AutoMigrate(&models.Task{})
	if err != nil {
		t.Fatalf("Could not migrate database: %v", err)
	}

	return db
}

// TeardownTestDB cleans up the test database
func TeardownTestDB(db config.Database) {
	// Drop the Task table after tests
	db.GetDB().Migrator().DropTable(&models.Task{})

	// Close the database connection
	err := db.Close()
	if err != nil {
		log.Printf("Error closing database connection: %v", err)
	}
}
