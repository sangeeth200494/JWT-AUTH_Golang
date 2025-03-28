package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Global database variable
var db *gorm.DB

func DBConnection() (error, db *gorm.DB) {
	godotenv.Load()
	// Load from environment variables
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"), os.Getenv("DB_PORT"),
	)

	// Connect to database
	//var err error
	db1, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}

	fmt.Println("Connected using environment variables!")

	errr := db1.Exec(`CREATE TABLE IF NOT EXISTS users (id TEXT PRIMARY KEY,username TEXT,password TEXT);`)
	if err != nil {
		fmt.Printf("error in creating unexciting users table in database: %v", errr.Error)
	}

	// Ensure database connection is closed when main function exits
	//defer DBC()
	return db1, nil
}

func DBC() {
	sqlDB, err := db.DB()
	if err != nil {
		log.Println("Error getting database instance:", err)
		return
	}

	err = sqlDB.Close()
	if err != nil {
		log.Println("Error closing database connection:", err)
		return
	}

	fmt.Println("Database connection closed successfully!")
}
