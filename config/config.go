package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Database interface for working with the database
type Database interface {
	Close() error
	GetDB() *gorm.DB
}

// RealDB implementation of the Database interface for a real database
type RealDB struct {
	db *gorm.DB
}

// Close closes the connection to the database
func (r *RealDB) Close() error {
	sqlDB, err := r.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// GetDB returns the *gorm.DB instance
func (r *RealDB) GetDB() *gorm.DB {
	return r.db
}

// LoadEnvVars loads environment variables from the .env file
func LoadEnvVars() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// getDatabaseName returns the correct database name based on the GO_ENV value
func getDatabaseName() string {
	goEnv := os.Getenv("GO_ENV")
	if goEnv == "test" {
		return os.Getenv("TEST_DB_NAME") // Use test database name
	}
	return os.Getenv("DB_NAME") // Use default database name for other environments
}

// ConnectDatabase connects to the database and returns a Database interface
func ConnectDatabase() (Database, error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		getDatabaseName(),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &RealDB{db: db}, nil
}
