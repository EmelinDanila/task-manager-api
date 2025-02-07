package controllers

import (
	"net/http"
	"strconv"

	"github.com/EmelinDanila/task-manager-api/models"
	"github.com/EmelinDanila/task-manager-api/services"
	"github.com/gin-gonic/gin"
)

// TaskController handles HTTP requests for task management
type TaskController struct {
	service services.TaskService
}

// NewTaskController sets up the task routes
func NewTaskController(router *gin.Engine, service services.TaskService) {
	controller := &TaskController{service: service}

	// Define task-related routes and their handlers
	router.POST("/tasks", controller.CreateTask)
	router.GET("/tasks/:id", controller.GetTaskByID)
	router.GET("/tasks", controller.GetAllTasks)
	router.PUT("/tasks/:id", controller.UpdateTask)
	router.DELETE("/tasks/:id", controller.DeleteTask)
}

// CreateTask adds a new task to the database
func (c *TaskController) CreateTask(ctx *gin.Context) {
	var task models.Task
	// Bind JSON request body to task object
	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the task using the service
	if err := c.service.CreateTask(&task); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the created task with HTTP status 201
	ctx.JSON(http.StatusCreated, task)
}

// GetTaskByID retrieves a task by its ID
func (c *TaskController) GetTaskByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// Fetch task by ID using the service
	task, err := c.service.GetTaskByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Return the task with HTTP status 200
	ctx.JSON(http.StatusOK, task)
}

// GetAllTasks retrieves all tasks from the database
func (c *TaskController) GetAllTasks(ctx *gin.Context) {
	tasks, err := c.service.GetAllTasks()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the list of tasks with HTTP status 200
	ctx.JSON(http.StatusOK, tasks)
}

// UpdateTask updates an existing task
func (c *TaskController) UpdateTask(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var task models.Task
	// Bind the request body to the task object
	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.ID = uint(id)
	// Update the task using the service
	if err := c.service.UpdateTask(&task); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the updated task with HTTP status 200
	ctx.JSON(http.StatusOK, task)
}

// DeleteTask deletes a task by its ID
func (c *TaskController) DeleteTask(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// Delete the task using the service
	if err := c.service.DeleteTask(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return HTTP status 204 (No Content)
	ctx.JSON(http.StatusNoContent, nil)
}
