package tests

import (
	"os"
	"testing"

	"github.com/EmelinDanila/task-manager-api/models"
	"github.com/EmelinDanila/task-manager-api/repository"
	"github.com/EmelinDanila/task-manager-api/services"
	"github.com/EmelinDanila/task-manager-api/tests/testutils"
	"github.com/stretchr/testify/assert"
)

func TestTaskService(t *testing.T) {
	// Save the current environment variable for later restoration
	oldEnv := os.Getenv("GO-ENV")

	// Initialize test dependencies, including setting up the test database
	db := testutils.SetupTestDB(t)
	defer testutils.TeardownTestDB(db) // Ensure database teardown after tests

	// Create repository and service for testing
	repo := repository.NewTaskRepository(db.GetDB())
	service := services.NewTaskService(repo)

	// Test: Creating a task
	t.Run("CreateTask", func(t *testing.T) {
		testutils.ClearTestDB(db) // Clear the database before starting the test

		task := &models.Task{
			Title:       "Test Task",
			Description: "Test Description",
			Status:      "Pending",
		}

		err := service.CreateTask(task)
		assert.NoError(t, err, "Should create task without errors")
		assert.NotZero(t, task.ID, "Task ID should be set after creation")
	})

	// Test: Fetching a task by ID
	t.Run("GetTaskByID", func(t *testing.T) {
		testutils.ClearTestDB(db) // Clear the database before starting the test

		task := &models.Task{
			Title:       "Test Task",
			Description: "Test Description",
			Status:      "Pending",
		}

		err := service.CreateTask(task)
		assert.NoError(t, err)

		fetchedTask, err := service.GetTaskByID(task.ID)
		assert.NoError(t, err, "Should fetch task without errors")
		assert.Equal(t, task.ID, fetchedTask.ID, "Fetched task ID should match")
	})

	// Test: Fetching all tasks
	t.Run("GetAllTasks", func(t *testing.T) {
		testutils.ClearTestDB(db) // Clear the database before starting the test

		err := service.CreateTask(&models.Task{Title: "Task 1"})
		assert.NoError(t, err)

		err = service.CreateTask(&models.Task{Title: "Task 2"})
		assert.NoError(t, err)

		tasks, err := service.GetAllTasks()
		assert.NoError(t, err, "Should fetch tasks without errors")
		assert.Len(t, tasks, 2, "Should return 2 tasks")
	})

	// Test: Updating a task
	t.Run("UpdateTask", func(t *testing.T) {
		testutils.ClearTestDB(db) // Clear the database before starting the test

		task := &models.Task{
			Title:       "Old Task",
			Description: "Old Description",
			Status:      "Pending",
		}

		err := service.CreateTask(task)
		assert.NoError(t, err)

		task.Title = "Updated Task"
		err = service.UpdateTask(task)
		assert.NoError(t, err)

		updatedTask, err := service.GetTaskByID(task.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Task", updatedTask.Title, "Task title should be updated")
	})

	// Test: Deleting a task
	t.Run("DeleteTask", func(t *testing.T) {
		testutils.ClearTestDB(db) // Clear the database before starting the test

		task := &models.Task{
			Title:       "Task to Delete",
			Description: "Description",
			Status:      "Pending",
		}

		err := service.CreateTask(task)
		assert.NoError(t, err)

		err = service.DeleteTask(task.ID)
		assert.NoError(t, err)

		_, err = service.GetTaskByID(task.ID)
		assert.Error(t, err, "Fetching deleted task should return an error")
	})

	// Restore the original environment variable
	os.Setenv("GO-ENV", oldEnv)
}
