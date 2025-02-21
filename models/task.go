package models

import (
	"time"

	"gorm.io/gorm"
)

// Task represents the task model
// @Description Task model containing task details.
// @property ID uint "Unique identifier for the task"
// @property Title string "Title of the task"
// @property Description string "Detailed description of the task"
// @property Status string "Current status of the task (Pending, In Progress, Completed)"
// @property UserID uint "ID of the user associated with the task"
// @property CreatedAt time.Time "Timestamp when the task was created"
// @property UpdatedAt time.Time "Timestamp when the task was last updated"
type Task struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"not null" json:"title"`
	Description string         `json:"description"`
	Status      string         `gorm:"default:'Pending'" json:"status"` // Possible values: Pending, In Progress, Completed
	UserID      uint           `json:"user_id"`                         // Relationship with user (if provided)
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"` // Field for soft delete
}

// TableName allows setting the table name for the Task model
// @Description Customizes the table name for tasks in the database.
func (Task) TableName() string {
	return "tasks"
}

// BeforeCreate sets default values before creating a task
// @Description Ensures the task status is set to 'Pending' if not provided.
func (t *Task) BeforeCreate(tx *gorm.DB) (err error) {
	if t.Status == "" {
		t.Status = "Pending"
	}
	return
}
