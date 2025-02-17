package controllers

import (
	"net/http"
	"regexp"

	"github.com/EmelinDanila/task-manager-api/models"
	"github.com/EmelinDanila/task-manager-api/repository"
	"github.com/EmelinDanila/task-manager-api/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

// AuthController handles authentication-related operations.
type AuthController struct {
	authService services.AuthService      // Service for authentication operations
	userRepo    repository.UserRepository // Repository for user data access
	validate    *validator.Validate       // Validator for request data
}

// NewAuthController creates a new instance of AuthController.
func NewAuthController(authService services.AuthService, userRepo repository.UserRepository) *AuthController {
	return &AuthController{
		authService: authService,
		userRepo:    userRepo,
		validate:    validator.New(),
	}
}

// RegisterUser handles user registration.
func (ac *AuthController) RegisterUser(c *gin.Context) {
	var requestData struct {
		Email    string `json:"email" validate:"required,email"`    // User's email
		Password string `json:"password" validate:"required,min=8"` // User's password
	}

	// Bind and validate the request data
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Validate email format
	if !isValidEmail(requestData.Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format"})
		return
	}

	// Validate password strength
	if !isValidPassword(requestData.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long and contain a number, a special character, and an uppercase letter"})
		return
	}

	// Check if the user already exists
	existingUser, _ := ac.userRepo.FindByEmail(requestData.Email)
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Create a new user
	user := &models.User{
		Email:    requestData.Email,
		Password: requestData.Password,
	}
	if err := ac.userRepo.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// LoginUser handles user login.
func (ac *AuthController) LoginUser(c *gin.Context) {
	var loginData struct {
		Email    string `json:"email" validate:"required,email"`    // User's email
		Password string `json:"password" validate:"required,min=8"` // User's password
	}

	// Bind and validate the request data
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Check if email and password are provided
	if loginData.Email == "" || loginData.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
		return
	}

	// Find the user by email
	user, err := ac.userRepo.FindByEmail(loginData.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Check if the user exists
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Verify the password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate a JWT token
	token, err := ac.authService.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	// Return the token in the response
	c.JSON(http.StatusOK, gin.H{"token": token})
}

// isValidEmail checks if the email is in a valid format.
func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// isValidPassword checks if the password meets the strength requirements.
func isValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasUpper, hasNumber, hasSpecial := false, false, false
	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= '0' && char <= '9':
			hasNumber = true
		case (char >= '!' && char <= '/') || (char >= ':' && char <= '@') || (char >= '[' && char <= '`') || (char >= '{' && char <= '~'):
			hasSpecial = true
		}
	}
	return hasUpper && hasNumber && hasSpecial
}
