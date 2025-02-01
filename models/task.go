package models

import (
	"time"

	"gorm.io/gorm"
)

// Task represents the task model
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
func (Task) TableName() string {
	return "tasks"
}

// BeforeCreate sets default values before creating a task
func (t *Task) BeforeCreate(tx *gorm.DB) (err error) {
	if t.Status == "" {
		t.Status = "Pending"
	}
	return
}
