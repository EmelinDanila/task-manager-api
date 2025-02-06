package services

import (
	"errors"

	"github.com/EmelinDanila/task-manager-api/models"
	"github.com/EmelinDanila/task-manager-api/repository"
)

// TaskService defines the interface for working with tasks.
type TaskService interface {
	CreateTask(task *models.Task) error        // Create a new task.
	GetTaskByID(id uint) (*models.Task, error) // Get a task by its ID.
	GetAllTasks() ([]models.Task, error)       // Retrieve all tasks.
	UpdateTask(task *models.Task) error        // Update an existing task.
	DeleteTask(id uint) error                  // Delete a task by its ID.
}

type taskService struct {
	repo repository.TaskRepository // Repository for interacting with the database.
}

// NewTaskService creates a new instance of TaskService.
func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskService{repo: repo}
}

// CreateTask adds a new task to the database.
func (s *taskService) CreateTask(task *models.Task) error {
	// Ensure task title is provided.
	if task.Title == "" {
		return errors.New("task title cannot be empty")
	}
	// Use repository to create the task.
	return s.repo.Create(task)
}

// GetTaskByID fetches a task from the database by its ID.
func (s *taskService) GetTaskByID(id uint) (*models.Task, error) {
	// Retrieve the task using the repository.
	return s.repo.GetByID(id)
}

// GetAllTasks retrieves all tasks from the database.
func (s *taskService) GetAllTasks() ([]models.Task, error) {
	// Use repository to fetch all tasks.
	return s.repo.GetAll()
}

// UpdateTask updates the details of an existing task.
func (s *taskService) UpdateTask(task *models.Task) error {
	// Ensure task has a valid ID.
	if task.ID == 0 {
		return errors.New("task ID cannot be zero")
	}
	// Update the task using the repository.
	return s.repo.Update(task)
}

// DeleteTask deletes a task from the database by its ID.
func (s *taskService) DeleteTask(id uint) error {
	// Delete the task using the repository.
	return s.repo.Delete(id)
}
