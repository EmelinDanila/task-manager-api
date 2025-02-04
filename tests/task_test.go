package tests

import (
	"os"
	"testing"

	"github.com/EmelinDanila/task-manager-api/models"
	"github.com/EmelinDanila/task-manager-api/tests/testutils"
	"github.com/stretchr/testify/assert"
)

func TestTaskModel(t *testing.T) {
	// Save the current environment
	oldEnv := os.Getenv("GO_ENV")

	// Set up the test database
	db := testutils.SetupTestDB(t)
	defer testutils.TeardownTestDB(db)

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

	// Restore the original environment
	os.Setenv("GO_ENV", oldEnv)
}
