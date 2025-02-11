package models

import (
	"gorm.io/gorm"
)

// User represents the user model
type User struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"password"`
}
