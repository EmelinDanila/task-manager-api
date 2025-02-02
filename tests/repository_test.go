package tests

import (
	"testing"

	"github.com/EmelinDanila/task-manager-api/models"
	"github.com/EmelinDanila/task-manager-api/repository"
	"github.com/EmelinDanila/task-manager-api/tests/testutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestTaskRepository(t *testing.T) {
	// Initialize the test database
	db := testutils.SetupTestDB(t)
	defer testutils.TeardownTestDB(db)

	// Create the repository
	repo := repository.NewTaskRepository(db.GetDB())

	t.Run("Create Task", func(t *testing.T) {
		task := &models.Task{
			Title:       "Test Task",
			Description: "This is a test task",
			Status:      "Pending",
		}

		// Create the task using the repository
		err := repo.Create(task)
		assert.NoError(t, err)
		assert.NotZero(t, task.ID)

		// Fetch the task from the database and check if it's correctly saved
		var fetchedTask models.Task
		err = db.GetDB().First(&fetchedTask, task.ID).Error
		assert.NoError(t, err)
		assert.Equal(t, "Test Task", fetchedTask.Title)
		defer db.GetDB().Delete(&fetchedTask)
	})

	t.Run("Get Task by ID", func(t *testing.T) {
		task := &models.Task{
			Title:       "Test Task",
			Description: "This is a test task",
			Status:      "Pending",
		}
		db.GetDB().Create(task)

		// Retrieve the task by ID and check if it matches
		fetchedTask, err := repo.GetByID(task.ID)
		assert.NoError(t, err)
		assert.Equal(t, task.Title, fetchedTask.Title)
		assert.Equal(t, task.ID, fetchedTask.ID)
		defer db.GetDB().Delete(&fetchedTask)
	})

	t.Run("Get All Tasks", func(t *testing.T) {
		tasks := []models.Task{
			{Title: "Task 1", Description: "First task", Status: "Pending"},
			{Title: "Task 2", Description: "Second task", Status: "Completed"},
		}

		// Save tasks to the database
		for _, task := range tasks {
			db.GetDB().Create(&task)
		}

		// Retrieve all tasks and check if the count matches
		allTasks, err := repo.GetAll()
		assert.NoError(t, err)
		assert.Len(t, allTasks, len(tasks))

		// Clean up the tasks
		for _, task := range tasks {
			defer db.GetDB().Delete(&task)
		}
	})

	t.Run("Update Task", func(t *testing.T) {
		task := &models.Task{
			Title:       "Old Task Title",
			Description: "Old description",
			Status:      "Pending",
		}
		db.GetDB().Create(task)

		// Update the task's title and description
		task.Title = "Updated Task Title"
		task.Description = "Updated description"
		err := repo.Update(task)
		assert.NoError(t, err)

		// Fetch the updated task and check the new values
		var updatedTask models.Task
		db.GetDB().First(&updatedTask, task.ID)
		assert.Equal(t, "Updated Task Title", updatedTask.Title)
		assert.Equal(t, "Updated description", updatedTask.Description)
		defer db.GetDB().Delete(&updatedTask)
	})

	t.Run("Delete Task", func(t *testing.T) {
		task := &models.Task{
			Title:       "Task to be deleted",
			Description: "This task will be deleted",
			Status:      "Pending",
		}
		db.GetDB().Create(task)

		// Delete the task
		err := repo.Delete(task.ID)
		assert.NoError(t, err)

		// Check if the task was deleted from the database
		var deletedTask models.Task
		err = db.GetDB().First(&deletedTask, task.ID).Error
		assert.Error(t, err) // We expect an error because the task should be deleted
		assert.Equal(t, gorm.ErrRecordNotFound, err)
	})
}
