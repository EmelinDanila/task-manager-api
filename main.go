package main

import (
	"log"

	"github.com/EmelinDanila/task-manager-api/config"
	"github.com/EmelinDanila/task-manager-api/migrations"
	"github.com/EmelinDanila/task-manager-api/routes"
	"github.com/gin-gonic/gin"
)

// @title Task Manager API
// @version 1.0
// @description Task Manager API is a RESTful API for managing tasks. The application will allow users to create, update, delete, and view tasks. It will use Go for the backend, PostgreSQL for the database, and JWT for authentication.
// @termsOfService http://swagger.io/terms/

// @contact.name Danila Emelin
// @contact.url https://github.com/EmelinDanila
// @contact.email d.emelin.qa@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost: 8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
// @description Use 'Bearer' followed by your JWT token. Example: "Bearer your_token_here"

func main() {
	config.LoadEnvVars()

	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatal("Could not connect to the database")
	}

	sqlDB, _ := db.GetDB().DB()
	defer sqlDB.Close()

	migrations.Migrate(db.GetDB())

	router := gin.Default()

	routes.SetupRoutes(router, db.GetDB())

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
