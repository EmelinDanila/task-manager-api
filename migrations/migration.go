package migrations

import (
	"fmt"
	"log"

	"github.com/EmelinDanila/task-manager-api/models"
	"gorm.io/gorm"
)

// Migrate выполняет все миграции
func Migrate(db *gorm.DB) {
	// Автоматически создает таблицы на основе моделей
	if err := db.AutoMigrate(&models.User{}, &models.Task{}); err != nil { // Проверяем ошибку непосредственно
		log.Fatalf("Migration failed: %v", err)
	}
	fmt.Println("Database migration completed successfully!")
}
