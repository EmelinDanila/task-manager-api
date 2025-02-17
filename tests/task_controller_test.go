package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EmelinDanila/task-manager-api/controllers"
	"github.com/EmelinDanila/task-manager-api/middleware"
	"github.com/EmelinDanila/task-manager-api/models"
	"github.com/EmelinDanila/task-manager-api/repository"
	"github.com/EmelinDanila/task-manager-api/services"
	"github.com/EmelinDanila/task-manager-api/tests/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupTaskControllerTest prepares test dependencies for TaskController.
func setupTaskControllerTest(t *testing.T) (*gin.Engine, services.TaskService, uint, string) {
	gin.SetMode(gin.TestMode)
	db := testutils.SetupTestDB(t)
	testutils.ClearTestDB(db)
	router := gin.Default()

	// Setup dependencies
	taskRepo := repository.NewTaskRepository(db.GetDB())
	taskService := services.NewTaskService(taskRepo)
	taskController := controllers.TaskController{Service: taskService}
	authService := services.NewAuthService()

	// Register protected routes
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware(authService))
	{
		protected.POST("/tasks", taskController.CreateTask)
		protected.GET("/tasks", taskController.GetAllTasks)
		protected.GET("/tasks/:id", taskController.GetTaskByID)
		protected.PUT("/tasks/:id", taskController.UpdateTask)
		protected.DELETE("/tasks/:id", taskController.DeleteTask)
	}

	// Create test user and generate token
	userRepo := repository.NewUserRepository(db.GetDB())
	user := &models.User{Email: "test@example.com", Password: "Password123!"}
	userRepo.CreateUser(user)
	token, _ := authService.GenerateToken(user.ID)

	return router, taskService, user.ID, token
}

// TestTaskCreation verifies task creation.
func TestTaskCreation(t *testing.T) {
	router, _, _, token := setupTaskControllerTest(t)

	taskData := `{"title": "Test Task", "description": "Task Description"}`
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBufferString(taskData))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

// TestFetchingAllTasks verifies getting all tasks.
func TestFetchingAllTasks(t *testing.T) {
	router, _, _, token := setupTaskControllerTest(t)

	req, _ := http.NewRequest("GET", "/tasks", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestFetchingTaskByID verifies getting a task by ID.
func TestFetchingTaskByID(t *testing.T) {
	router, service, userID, token := setupTaskControllerTest(t)

	task := &models.Task{Title: "Sample Task", Description: "Some description", UserID: userID}
	service.CreateTask(task)

	req, _ := http.NewRequest("GET", "/tasks/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestUpdatingTask verifies task updating.
func TestUpdatingTask(t *testing.T) {
	router, service, userID, token := setupTaskControllerTest(t)

	task := &models.Task{Title: "Old Task", Description: "Old Description", UserID: userID}
	service.CreateTask(task)

	updatedTask := `{"title": "Updated Task", "description": "Updated Description"}`
	req, _ := http.NewRequest("PUT", "/tasks/1", bytes.NewBufferString(updatedTask))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestRemovingTask verifies task deletion.
func TestRemovingTask(t *testing.T) {
	router, service, userID, token := setupTaskControllerTest(t)

	task := &models.Task{Title: "Task to Delete", Description: "Will be deleted", UserID: userID}
	service.CreateTask(task)

	req, _ := http.NewRequest("DELETE", "/tasks/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

// TestUnauthorizedAccess verifies that unauthorized users cannot create a task.
func TestUnauthorizedAccess(t *testing.T) {
	router, _, _, _ := setupTaskControllerTest(t) // No token provided

	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBufferString(`{"title": "Unauthorized Task"}`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code) // Should return 401
}

// TestUnauthorizedTaskRetrieval verifies that unauthorized users cannot get tasks.
func TestUnauthorizedTaskRetrieval(t *testing.T) {
	router, _, _, _ := setupTaskControllerTest(t) // No token provided

	req, _ := http.NewRequest("GET", "/tasks", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code) // Should return 401
}

// TestAccessingOthersTask verifies that a user cannot access another user's task.
func TestAccessingOthersTask(t *testing.T) {
	router, service, _, token := setupTaskControllerTest(t)

	// Create a task for another user (userID = 2)
	task := &models.Task{Title: "Task from another user", Description: "Not yours", UserID: 2}
	service.CreateTask(task)

	req, _ := http.NewRequest("GET", "/tasks/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Expect 404 Not Found or 403 Forbidden
	assert.Contains(t, []int{http.StatusNotFound, http.StatusForbidden}, w.Code)
}

// TestUpdatingOthersTask ensures a user cannot update another user's task.
func TestUpdatingOthersTask(t *testing.T) {
	router, service, _, token := setupTaskControllerTest(t)

	// Create a task for another user
	task := &models.Task{Title: "Task from another user", Description: "Not yours", UserID: 2}
	service.CreateTask(task)

	updatedTask := `{"title": "Updated Task", "description": "Should not be allowed"}`
	req, _ := http.NewRequest("PUT", "/tasks/1", bytes.NewBufferString(updatedTask))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Expect 403 Forbidden or 404 Not Found instead of 500
	assert.Contains(t, []int{http.StatusNotFound, http.StatusForbidden}, w.Code)
}

// TestDeletingOthersTask ensures a user cannot delete another user's task.
func TestDeletingOthersTask(t *testing.T) {
	router, service, _, token := setupTaskControllerTest(t)

	// Create a task for another user
	task := &models.Task{Title: "Task from another user", Description: "Not yours", UserID: 2}
	service.CreateTask(task)

	req, _ := http.NewRequest("DELETE", "/tasks/1", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Expect 403 Forbidden or 404 Not Found instead of 500
	assert.Contains(t, []int{http.StatusNotFound, http.StatusForbidden}, w.Code)
}
