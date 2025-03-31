package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DBConnection() (db *gorm.DB, err error) {
	godotenv.Load()
	// Load from environment variables
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"), os.Getenv("DB_PORT"),
	)

	// Connect to database
	db, errr := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if errr != nil {
		log.Fatal("Failed to connect:", errr.Error())
	}

	fmt.Println("Connected using environment variables!")

	res := db.Exec(`CREATE TABLE IF NOT EXISTS users (id TEXT PRIMARY KEY,username TEXT,password TEXT);`)
	if res.Error != nil {
		fmt.Printf("error in creating unexciting users table in database: %v", res.Error)
	}
	return db, nil
}

func DBC(db *gorm.DB) {
	sqlDB, err := db.DB() // Get the raw SQL DB connection
	if err != nil {
		log.Println("Error getting SQL database:", err)
		return
	}

	err = sqlDB.Close() // Close the connection
	if err != nil {
		log.Println("Error closing database connection:", err)
		return
	}

	fmt.Println("Database connection closed successfully!")
}
