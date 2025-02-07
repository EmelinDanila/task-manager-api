package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/EmelinDanila/task-manager-api/controllers"
	"github.com/EmelinDanila/task-manager-api/models"
	"github.com/EmelinDanila/task-manager-api/repository"
	"github.com/EmelinDanila/task-manager-api/services"
	"github.com/EmelinDanila/task-manager-api/tests/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTaskController(t *testing.T) {
	// Save the old environment variable to restore later
	oldEnv := os.Getenv("GO-ENV")
	// Set up a test database
	db := testutils.SetupTestDB(t)
	defer testutils.TeardownTestDB(db)

	// Create necessary repository, service, and router
	repo := repository.NewTaskRepository(db.GetDB())
	service := services.NewTaskService(repo)
	router := gin.Default()
	controllers.NewTaskController(router, service)

	// Test case for creating a new task
	t.Run("CreateTask", func(t *testing.T) {
		testutils.ClearTestDB(db)

		task := models.Task{
			Title:       "New Task",
			Description: "This is a new task",
			Status:      "Pending",
		}

		// Create HTTP request with JSON body
		body, _ := json.Marshal(task)
		req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		// Recorder to capture response
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		// Assert the status code is 201 (Created)
		assert.Equal(t, http.StatusCreated, resp.Code)

		var createdTask models.Task
		// Unmarshal response body into task object
		json.Unmarshal(resp.Body.Bytes(), &createdTask)
		// Assert that the task title is the same as the one sent
		assert.Equal(t, "New Task", createdTask.Title)
	})

	// Test case for fetching a task by its ID
	t.Run("GetTaskByID", func(t *testing.T) {
		testutils.ClearTestDB(db)

		// Create a new task
		task := &models.Task{
			Title: "Test Task",
		}
		err := service.CreateTask(task)
		assert.NoError(t, err)

		// Make a GET request to fetch the task by ID
		req, _ := http.NewRequest("GET", fmt.Sprintf("/tasks/%d", task.ID), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		// Assert the response status is 200 (OK)
		assert.Equal(t, http.StatusOK, resp.Code)

		var fetchedTask models.Task
		// Unmarshal response body into task object
		json.Unmarshal(resp.Body.Bytes(), &fetchedTask)
		// Assert that the fetched task ID is the same as the created task ID
		assert.Equal(t, task.ID, fetchedTask.ID)
	})

	// Test case for fetching all tasks
	t.Run("GetAllTasks", func(t *testing.T) {
		testutils.ClearTestDB(db)

		// Create two tasks
		service.CreateTask(&models.Task{Title: "Task 1"})
		service.CreateTask(&models.Task{Title: "Task 2"})

		// Make a GET request to fetch all tasks
		req, _ := http.NewRequest("GET", "/tasks", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		// Assert the response status is 200 (OK)
		assert.Equal(t, http.StatusOK, resp.Code)

		var tasks []models.Task
		// Unmarshal response body into tasks array
		json.Unmarshal(resp.Body.Bytes(), &tasks)
		// Assert that two tasks were returned
		assert.Len(t, tasks, 2)
	})

	// Test case for updating an existing task
	t.Run("UpdateTask", func(t *testing.T) {
		testutils.ClearTestDB(db)

		// Create a new task
		task := &models.Task{Title: "Old Task"}
		err := service.CreateTask(task)
		assert.NoError(t, err)

		// Update the task title
		task.Title = "Updated Task"
		body, _ := json.Marshal(task)
		req, _ := http.NewRequest("PUT", fmt.Sprintf("/tasks/%d", task.ID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		// Make a PUT request to update the task
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		// Assert the response status is 200 (OK)
		assert.Equal(t, http.StatusOK, resp.Code)

		var updatedTask models.Task
		// Unmarshal response body into task object
		json.Unmarshal(resp.Body.Bytes(), &updatedTask)
		// Assert that the task title was updated
		assert.Equal(t, "Updated Task", updatedTask.Title)
	})

	// Test case for deleting a task
	t.Run("DeleteTask", func(t *testing.T) {
		testutils.ClearTestDB(db)

		// Create a new task
		task := &models.Task{Title: "Task to Delete"}
		err := service.CreateTask(task)
		assert.NoError(t, err)

		// Make a DELETE request to delete the task by ID
		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/tasks/%d", task.ID), nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		// Assert the response status is 204 (No Content)
		assert.Equal(t, http.StatusNoContent, resp.Code)
	})

	// Restore the old environment variable
	os.Setenv("GO-ENV", oldEnv)
}
