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

// @Summary Create a new task
// @Description Create a new task for the authenticated user
// @Tags tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body object{title=string,description=string,status=string} true "Task data"
// @Success 201 {object} models.TaskResponse "Task created successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid request data"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /tasks [post]
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

// @Summary Get a task by ID
// @Description Retrieves a specific task by ID for the authenticated user
// @Tags tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Task ID"
// @Success 200 {object} models.TaskResponse "Task found"
// @Failure 400 {object} models.ErrorResponse "Invalid task ID"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden: You cannot access this task"
// @Failure 404 {object} models.ErrorResponse "Task not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /tasks/{id} [get]
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

// @Summary Get all tasks for the authenticated user
// @Tags tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} models.TaskListResponse "List of tasks for the authenticated user"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /tasks [get]
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

// @Summary Update an existing task
// @Description Update a task only if the authenticated user is the owner of the task
// @Tags tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Task ID"
// @Param request body object{title=string,description=string,status=string} true "Updated task data"
// @Success 200 {object} models.TaskResponse "Task updated successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid task ID or request data"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden: You cannot update another user's task"
// @Failure 404 {object} models.ErrorResponse "Task not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /tasks/{id} [put]
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

// @Summary Delete a task
// @Description Delete a task only if the authenticated user is the owner of the task
// @Tags tasks
// @Security ApiKeyAuth
// @Param id path int true "Task ID"
// @Success 204 "Task deleted successfully"
// @Failure 400 {object} models.ErrorResponse "Invalid task ID"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 403 {object} models.ErrorResponse "Forbidden: You cannot delete another user's task"
// @Failure 404 {object} models.ErrorResponse "Task not found"
// @Failure 500 {object} models.ErrorResponse "Internal server error"
// @Router /tasks/{id} [delete]
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
