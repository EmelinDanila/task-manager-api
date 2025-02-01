package tests

import (
	"testing"

	"github.com/EmelinDanila/task-manager-api/config"
	"github.com/EmelinDanila/task-manager-api/models"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) config.Database {
	// Load environment variables
	config.LoadEnvVars()

	// Connect to the database through the Database interface
	db, err := config.ConnectDatabase()
	if err != nil {
		t.Fatalf("Could not connect to the database: %v", err)
	}

	// Create the table for the Task model
	err = db.GetDB().AutoMigrate(&models.Task{})
	if err != nil {
		t.Fatalf("Could not migrate database: %v", err)
	}

	return db
}

func teardownTestDB(db config.Database) {
	// Drop the table after the tests
	db.GetDB().Migrator().DropTable(&models.Task{})
	db.Close()
}

func TestTaskModel(t *testing.T) {
	// Set up the test database
	db := setupTestDB(t)
	defer teardownTestDB(db)

	// Create a test task
	task := models.Task{
		Title:       "Test Task",
		Description: "This is a test task",
		Status:      "In Progress",
		UserID:      1, // Example user ID
	}

	// Save the task to the database
	result := db.GetDB().Create(&task)

	// Check that there are no errors during saving
	assert.NoError(t, result.Error)

	// Check that the task was saved in the database
	var savedTask models.Task
	db.GetDB().First(&savedTask, task.ID)
	assert.Equal(t, task.Title, savedTask.Title)
	assert.Equal(t, task.Description, savedTask.Description)
	assert.Equal(t, task.Status, savedTask.Status)
}
