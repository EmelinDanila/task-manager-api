package services

import (
	"errors"

	"github.com/EmelinDanila/task-manager-api/models"
	"github.com/EmelinDanila/task-manager-api/repository"
	"gorm.io/gorm"
)

// TaskService defines the interface for working with tasks.
type TaskService interface {
	CreateTask(task *models.Task) error
	GetTaskByID(id, userID uint) (*models.Task, error)
	GetUserTasks(userID uint) ([]models.Task, error)
	UpdateTask(task *models.Task, userID uint) error
	DeleteTask(id, userID uint) error
}

type taskService struct {
	repo repository.TaskRepository
}

// NewTaskService creates a new instance of TaskService.
func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskService{repo: repo}
}

// CreateTask ensures task belongs to a user before saving.
func (s *taskService) CreateTask(task *models.Task) error {
	if task.Title == "" {
		return errors.New("task title cannot be empty")
	}
	return s.repo.Create(task)
}

// GetTaskByID ensures user can only retrieve their own tasks.
func (s *taskService) GetTaskByID(taskID, userID uint) (*models.Task, error) {
	task := &models.Task{}
	err := s.repo.GetByIDAndUserID(taskID, userID, task)

	// If the task is not found, we return "task not found" (404)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("task not found")
	}

	// If another error occurs, we return it
	if err != nil {
		return nil, err
	}

	return task, nil
}

// GetUserTasks ensures user only sees their own tasks.
func (s *taskService) GetUserTasks(userID uint) ([]models.Task, error) {
	var tasks []models.Task
	if err := s.repo.GetByUserID(userID, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

// UpdateTask checks if the user owns the task before updating.
func (s *taskService) UpdateTask(task *models.Task, userID uint) error {
	existingTask, err := s.GetTaskByID(task.ID, userID)
	if err != nil {
		return err // Уже содержит "task not found"
	}

	// Проверяем, владеет ли пользователь задачей
	if existingTask.UserID != userID {
		return errors.New("forbidden") // 403 Forbidden
	}

	// Обновляем только разрешенные поля
	existingTask.Title = task.Title
	existingTask.Description = task.Description
	existingTask.Status = task.Status

	return s.repo.Update(existingTask)
}

// DeleteTask ensures only the owner can delete a task.
func (s *taskService) DeleteTask(id, userID uint) error {
	task, err := s.GetTaskByID(id, userID)
	if err != nil {
		return err // Уже содержит "task not found"
	}

	if task.UserID != userID {
		return errors.New("forbidden") // 403
	}

	return s.repo.Delete(task.ID)
}
