package database

import (
	"log"

	"github.com/zhakazx/cleanshort/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Connect establishes a connection to the PostgreSQL database
func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return db, nil
}

// Migrate runs the database migrations
func Migrate(db *gorm.DB) error {
	log.Println("Running database migrations...")

	// Enable UUID extension
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return err
	}

	// Auto-migrate models
	err := db.AutoMigrate(
		&models.User{},
		&models.Link{},
		&models.RefreshToken{},
	)
	if err != nil {
		return err
	}

	// Create additional indexes
	if err := createIndexes(db); err != nil {
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// createIndexes creates additional database indexes for performance
func createIndexes(db *gorm.DB) error {
	// Index for links table
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_links_user_code ON links(user_id, short_code)").Error; err != nil {
		return err
	}

	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_links_short_code ON links(short_code)").Error; err != nil {
		return err
	}

	// Index for refresh tokens
	if err := db.Exec("CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user ON refresh_tokens(user_id)").Error; err != nil {
		return err
	}

	return nil
}