package repository

import (
	"github.com/EmelinDanila/task-manager-api/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserRepository defines the interface for user-related database operations
type UserRepository interface {
	FindByEmail(email string) (*models.User, error) // Find a user by email
	CreateUser(user *models.User) error             // Create a new user
}

// userRepository implements the UserRepository interface
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// FindByEmail searches for a user by their email address
func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	// Query the database for a user with the given email
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		// If no user is found, return nil without an error
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		// For any other error, return it
		return nil, err
	}
	// Return the found user
	return &user, nil
}

// CreateUser adds a new user to the database
func (r *userRepository) CreateUser(user *models.User) error {
	// Hash the password before saving
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	// Replace the plain text password with the hashed version
	user.Password = string(hashedPassword)

	// Insert the new user into the database
	return r.db.Create(user).Error
}
