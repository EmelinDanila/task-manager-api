package tests

import (
	"testing"

	"github.com/EmelinDanila/task-manager-api/models"
	"github.com/EmelinDanila/task-manager-api/repository"
	"github.com/EmelinDanila/task-manager-api/tests/testutils"

	"github.com/stretchr/testify/assert"
)

// TestUserRepository groups related tests for user repository functions
func TestUserRepository(t *testing.T) {
	// Setup the test database
	db := testutils.SetupTestDB(t)
	defer testutils.TeardownTestDB(db)

	// Create repository instance
	repo := repository.NewUserRepository(db.GetDB())

	// Test case for user creation
	t.Run("CreateUser", func(t *testing.T) {
		user := &models.User{
			Email:    "user@example.com",
			Password: "password123",
		}

		// Attempt to create a new user
		err := repo.CreateUser(user)
		assert.NoError(t, err, "Expected no error while creating a user")

		// Verify that the user was successfully created
		foundUser, err := repo.FindByEmail(user.Email)
		assert.NoError(t, err, "Expected no error while retrieving the user")
		assert.NotNil(t, foundUser, "User should not be nil")
		assert.Equal(t, user.Email, foundUser.Email, "Emails should match")

		t.Log(foundUser)
	})

	// Test case for finding a user by email
	t.Run("FindByEmail", func(t *testing.T) {
		user := &models.User{
			Email:    "user@example.com",
			Password: "password123",
		}

		// Create a user directly in the database
		repo.CreateUser(user)

		// Attempt to find the user by email
		foundUser, err := repo.FindByEmail(user.Email)
		assert.NoError(t, err, "Expected no error while retrieving the user")
		assert.NotNil(t, foundUser, "User should not be nil")
		assert.Equal(t, user.Email, foundUser.Email, "Emails should match")

		t.Log(foundUser)
	})
}
