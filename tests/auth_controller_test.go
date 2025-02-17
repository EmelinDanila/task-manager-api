package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/EmelinDanila/task-manager-api/controllers"
	"github.com/EmelinDanila/task-manager-api/models"
	"github.com/EmelinDanila/task-manager-api/repository"
	"github.com/EmelinDanila/task-manager-api/services"
	"github.com/EmelinDanila/task-manager-api/tests/testutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Test for user registration
func TestAuthController_RegisterUser(t *testing.T) {
	// Initializing the test database
	db := testutils.SetupTestDB(t)
	defer testutils.TeardownTestDB(db)

	// Creating required services and repositories
	authService := services.NewAuthService()
	repo := repository.NewUserRepository(db.GetDB())
	controller := controllers.NewAuthController(authService, repo)

	// Test data for user registration
	registerData := `{
		"email": "testuser@example.com",
		"password": "Password123!"
	}`

	// Request to register user
	req, _ := http.NewRequest("POST", "/register", bytes.NewBufferString(registerData))
	req.Header.Set("Content-Type", "application/json")

	// Creating a test recorder
	w := httptest.NewRecorder()

	// Registering route in Gin
	r := gin.Default()
	r.POST("/register", controller.RegisterUser)

	// Executing the request
	r.ServeHTTP(w, req)

	// Checking the response status
	assert.Equal(t, http.StatusCreated, w.Code)

	// Checking the response content
	assert.Contains(t, w.Body.String(), "User registered successfully")

	// Checking the registration attempt with existing email
	req, _ = http.NewRequest("POST", "/register", bytes.NewBufferString(registerData))
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Checking the error response for existing user
	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), "User already exists")
}

// Test for user login
func TestAuthController_LoginUser(t *testing.T) {
	// Initializing the test database
	db := testutils.SetupTestDB(t)
	defer testutils.TeardownTestDB(db)

	// Creating required services and repositories
	authService := services.NewAuthService()
	repo := repository.NewUserRepository(db.GetDB())
	controller := controllers.NewAuthController(authService, repo)

	// Creating a user for testing
	user := &models.User{
		Email:    "loginuser@example.com",
		Password: "Password123!",
	}

	// Saving the user to the database
	err := repo.CreateUser(user)
	if err != nil {
		t.Fatalf("Could not create user: %v", err)
	}

	// Test data for user login
	loginData := `{
		"email": "loginuser@example.com",
		"password": "Password123!"
	}`

	// Request to login user
	req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(loginData))
	req.Header.Set("Content-Type", "application/json")

	// Creating a test recorder
	w := httptest.NewRecorder()

	// Registering route in Gin
	r := gin.Default()
	r.POST("/login", controller.LoginUser)

	// Executing the request
	r.ServeHTTP(w, req)

	// Checking the response status
	assert.Equal(t, http.StatusOK, w.Code)

	// Checking the presence of token in the response
	assert.Contains(t, w.Body.String(), "token")
}

// Test for user login with invalid credentials
func TestAuthController_LoginUser_InvalidCredentials(t *testing.T) {
	// Initializing the test database
	db := testutils.SetupTestDB(t)
	defer testutils.TeardownTestDB(db)

	// Creating required services and repositories
	authService := services.NewAuthService()
	repo := repository.NewUserRepository(db.GetDB())
	controller := controllers.NewAuthController(authService, repo)

	// Test data for login with incorrect credentials
	loginData := `{
		"email": "wronguser@example.com",
		"password": "WrongPassword123!"
	}`

	// Request to login user
	req, _ := http.NewRequest("POST", "/login", bytes.NewBufferString(loginData))
	req.Header.Set("Content-Type", "application/json")

	// Creating a test recorder
	w := httptest.NewRecorder()

	// Registering route in Gin
	r := gin.Default()
	r.POST("/login", controller.LoginUser)

	// Executing the request
	r.ServeHTTP(w, req)

	// Checking the response status
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	// Checking the error message in the response
	assert.Contains(t, w.Body.String(), "Invalid email or password")
}
