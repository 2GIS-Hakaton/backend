package database

import (
	"fmt"
	"os"

	"github.com/dimmy-kor/audioguid/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect подключается к PostgreSQL
func Connect() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgresql://postgres:postgres@localhost:30101/audioguid?sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

// Migrate выполняет миграции БД
func Migrate(db *gorm.DB) error {
	// Пытаемся включить PostGIS extension (не критично если не получится)
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS postgis").Error; err != nil {
		fmt.Printf("Warning: Could not create PostGIS extension: %v\n", err)
		fmt.Println("Continuing without PostGIS (using simple distance calculations)...")
	}

	// Автомиграция моделей
	if err := db.AutoMigrate(
		&models.POI{},
		&models.Route{},
		&models.Waypoint{},
		&models.Content{},
	); err != nil {
		return fmt.Errorf("failed to migrate models: %w", err)
	}

	return nil
}
