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

// ClearTestDB clears the test database
func ClearTestDB(db config.Database) {
	var tables []string
	db.GetDB().Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'").Scan(&tables)

	for _, table := range tables {
		if err := db.GetDB().Exec("TRUNCATE TABLE " + table + " RESTART IDENTITY CASCADE").Error; err != nil {
			log.Fatalf("Failed to clear table %s: %v", table, err)
		}
	}
}

// TeardownTestDB cleans up the test database
func TeardownTestDB(db config.Database) {
	var tables []string
	db.GetDB().Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'").Scan(&tables)

	for _, table := range tables {
		if err := db.GetDB().Migrator().DropTable(table); err != nil {
			log.Printf("Failed to drop table %s: %v", table, err)
		}
	}

	// Close the database connection
	err := db.Close()
	if err != nil {
		log.Printf("Error closing database connection: %v", err)
	}
}
