package repository

import (
	"github.com/EmelinDanila/task-manager-api/models"
	"gorm.io/gorm"
)

// TaskRepository defines the interface for interacting with tasks in the database
type TaskRepository interface {
	Create(task *models.Task) error
	GetByID(id uint) (*models.Task, error)
	GetAll() ([]models.Task, error)
	Update(task *models.Task) error
	Delete(id uint) error
}

type taskRepository struct {
	db *gorm.DB
}

// NewTaskRepository initializes a new instance of TaskRepository
func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

// Create adds a new task to the database
func (r *taskRepository) Create(task *models.Task) error {
	return r.db.Create(task).Error
}

// GetByID retrieves a task by its ID
func (r *taskRepository) GetByID(id uint) (*models.Task, error) {
	var task models.Task
	if err := r.db.First(&task, id).Error; err != nil {
		return nil, err
	}
	return &task, nil
}

// GetAll retrieves all tasks from the database
func (r *taskRepository) GetAll() ([]models.Task, error) {
	var tasks []models.Task
	if err := r.db.Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

// Update modifies an existing task in the database
func (r *taskRepository) Update(task *models.Task) error {
	return r.db.Save(task).Error
}

// Delete removes a task from the database by its ID
func (r *taskRepository) Delete(id uint) error {
	return r.db.Delete(&models.Task{}, id).Error
}
