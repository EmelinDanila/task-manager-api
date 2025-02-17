package controllers

import (
	"net/http"
	"strconv"

	"github.com/EmelinDanila/task-manager-api/middleware"
	"github.com/EmelinDanila/task-manager-api/models"
	"github.com/EmelinDanila/task-manager-api/services"
	"github.com/gin-gonic/gin"
)

// TaskController handles HTTP requests for task management
type TaskController struct {
	Service services.TaskService
}

// NewTaskController sets up the task routes
func NewTaskController(router *gin.Engine, service services.TaskService) {
	controller := &TaskController{Service: service}

	router.POST("/tasks", controller.CreateTask)
	router.GET("/tasks/:id", controller.GetTaskByID)
	router.GET("/tasks", controller.GetAllTasks)
	router.PUT("/tasks/:id", controller.UpdateTask)
	router.DELETE("/tasks/:id", controller.DeleteTask)
}

// CreateTask ensures the task is assigned to the correct user.
func (c *TaskController) CreateTask(ctx *gin.Context) {
	userID, exists := middleware.GetUserID(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var task models.Task
	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.UserID = userID

	if err := c.Service.CreateTask(&task); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, task)
}

// GetTaskByID ensures user can only access their own tasks.
func (c *TaskController) GetTaskByID(ctx *gin.Context) {
	userID, exists := middleware.GetUserID(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	task, err := c.Service.GetTaskByID(uint(id), userID)
	if err != nil {
		if err.Error() == "task not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else if err.Error() == "forbidden" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, task)
}

// GetAllTasks returns only the tasks for the logged-in user.
func (c *TaskController) GetAllTasks(ctx *gin.Context) {
	userID, exists := middleware.GetUserID(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	tasks, err := c.Service.GetUserTasks(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, tasks)
}

// UpdateTask ensures only task owners can update.
func (c *TaskController) UpdateTask(ctx *gin.Context) {
	userID, exists := middleware.GetUserID(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var task models.Task
	if err := ctx.ShouldBindJSON(&task); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task.ID = uint(id)

	if err := c.Service.UpdateTask(&task, userID); err != nil {
		if err.Error() == "forbidden" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "You cannot update another user's task"})
		} else if err.Error() == "task not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, task)
}

// DeleteTask ensures only task owners can delete.
func (c *TaskController) DeleteTask(ctx *gin.Context) {
	userID, exists := middleware.GetUserID(ctx)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	if err := c.Service.DeleteTask(uint(id), userID); err != nil {
		if err.Error() == "forbidden" {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "You cannot delete another user's task"})
		} else if err.Error() == "task not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
