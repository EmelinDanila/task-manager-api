package main

import (
	"log"

	"github.com/EmelinDanila/task-manager-api/config"
	"github.com/EmelinDanila/task-manager-api/migrations"
	"github.com/EmelinDanila/task-manager-api/routes"
	"github.com/gin-gonic/gin"
)

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
