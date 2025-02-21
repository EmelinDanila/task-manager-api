package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents the user model
// @Description User model containing authentication details.
// @property ID uint "Unique identifier for the user"
// @property Email string "Email address of the user"
// @property Password string "Encrypted password of the user"
// @property CreatedAt string "Timestamp when the user was created"
// @property UpdatedAt string "Timestamp when the user was last updated"
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Email     string         `json:"email"`
	Password  string         `json:"password"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type UserRegisterRequest struct {
	Email    string `json:"email" example:"user@example.com"`   // Email пользователя
	Password string `json:"password" example:"StrongP@ssword1"` // Пароль пользователя
}
