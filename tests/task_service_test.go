package tests

import (
	"testing"

	"github.com/EmelinDanila/task-manager-api/models"
	"github.com/EmelinDanila/task-manager-api/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTaskRepository is a mock implementation of TaskRepository
type MockTaskRepository struct {
	mock.Mock
}

// GetByID implements repository.TaskRepository.
func (m *MockTaskRepository) GetByID(id uint) (*models.Task, error) {
	panic("unimplemented")
}

func (m *MockTaskRepository) Create(task *models.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskRepository) GetByIDAndUserID(taskID, userID uint, task *models.Task) error {
	args := m.Called(taskID, userID, task)
	return args.Error(0)
}

func (m *MockTaskRepository) GetByUserID(userID uint, tasks *[]models.Task) error {
	args := m.Called(userID, tasks)
	return args.Error(0)
}

func (m *MockTaskRepository) GetAll() ([]models.Task, error) {
	args := m.Called()
	return args.Get(0).([]models.Task), args.Error(1)
}

func (m *MockTaskRepository) Update(task *models.Task) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockTaskRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func Test2CreateTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	taskService := services.NewTaskService(mockRepo)

	task := &models.Task{Title: "Test Task", UserID: 1}
	mockRepo.On("Create", task).Return(nil)

	err := taskService.CreateTask(task)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetTaskByID(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	taskService := services.NewTaskService(mockRepo)

	task := &models.Task{ID: 1, Title: "Test Task", UserID: 1}
	mockRepo.On("GetByIDAndUserID", task.ID, task.UserID, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		*(args.Get(2).(*models.Task)) = *task
	})

	result, err := taskService.GetTaskByID(task.ID, task.UserID)

	assert.NoError(t, err)
	assert.Equal(t, task, result)
}

func TestUpdateTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	taskService := services.NewTaskService(mockRepo)

	task := &models.Task{ID: 1, Title: "Updated Task", UserID: 1}
	mockRepo.On("GetByIDAndUserID", task.ID, task.UserID, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		*(args.Get(2).(*models.Task)) = *task
	})
	mockRepo.On("Update", task).Return(nil)

	err := taskService.UpdateTask(task, task.UserID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteTask(t *testing.T) {
	mockRepo := new(MockTaskRepository)
	taskService := services.NewTaskService(mockRepo)

	task := &models.Task{ID: 1, UserID: 1}
	mockRepo.On("GetByIDAndUserID", task.ID, task.UserID, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		*(args.Get(2).(*models.Task)) = *task
	})
	mockRepo.On("Delete", task.ID).Return(nil)

	err := taskService.DeleteTask(task.ID, task.UserID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
