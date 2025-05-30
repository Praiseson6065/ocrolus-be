package database

import (
	"fmt"
	"log"

	"Praiseson6065/ocrolus-be/config"
	"Praiseson6065/ocrolus-be/models"

	"gorm.io/gorm"

	"gorm.io/driver/postgres"
)

var db *gorm.DB

// ConnectDB establishes a connection to the database using configuration from environment
func ConnectDB() error {
	// Get database configuration from our config package
	dbConfig := config.Config.Database

	// Build connection string
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata",
		dbConfig.Server, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.Port,
	)

	log.Printf("Connecting to database: %s@%s:%s/%s",
		dbConfig.User, dbConfig.Server, dbConfig.Port, dbConfig.DBName)

	// Connect to database
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	db = database
	log.Println("Connected to database successfully")
	return nil
}

// MigrateDB runs database schema migrations using GORM
func MigrateDB() error {
	log.Println("Running database migrations...")

	// Auto migrate all models
	err := db.AutoMigrate(
		&models.User{},
		&models.Article{},
		&models.RecentlyViewedArticle{},
	)

	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

func init() {
	// Connect to the database
	if err := ConnectDB(); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Run migrations
	if err := MigrateDB(); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return db
}
