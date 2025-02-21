package routes

import (
	"github.com/EmelinDanila/task-manager-api/controllers"
	"github.com/EmelinDanila/task-manager-api/docs"
	"github.com/EmelinDanila/task-manager-api/middleware"
	"github.com/EmelinDanila/task-manager-api/repository"
	"github.com/EmelinDanila/task-manager-api/services"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// Swagger documentation
	docs.SwaggerInfo.Title = "Task Manager API"
	docs.SwaggerInfo.Description = "This is a task manager API."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Auth service and controller setup
	authService := services.NewAuthService()
	userRepo := repository.NewUserRepository(db)
	authController := controllers.NewAuthController(authService, userRepo)

	// Auth routes
	router.POST("/register", authController.RegisterUser)
	router.POST("/login", authController.LoginUser)

	// Protected group for authenticated routes
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware(authService))
	{
		// Profile route
		protected.GET("/profile", func(c *gin.Context) {
			userID, exists := middleware.GetUserID(c)
			if !exists {
				c.JSON(401, gin.H{"error": "Unauthorized"})
				return
			}
			c.JSON(200, gin.H{"message": "Welcome!", "userID": userID})
		})

		// Task routes
		taskRepo := repository.NewTaskRepository(db)
		taskService := services.NewTaskService(taskRepo)
		taskController := controllers.TaskController{Service: taskService}
		// Create a task
		protected.POST("/tasks", taskController.CreateTask)

		// Get all tasks
		protected.GET("/tasks", taskController.GetAllTasks)

		// Get task by ID
		protected.GET("/tasks/:id", taskController.GetTaskByID)

		// Update task
		protected.PUT("/tasks/:id", taskController.UpdateTask)

		// Delete task
		protected.DELETE("/tasks/:id", taskController.DeleteTask)
	}
}
