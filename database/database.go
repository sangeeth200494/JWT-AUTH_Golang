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
	// er := godotenv.Load()
	// if er != nil {
	// 	log.Fatal("Error loading .env file:", er.Error())
	// }
	er := godotenv.Load(".env")
	if er != nil {
		log.Fatal(er.Error())
	}
	fmt.Println("1111")
	// Load from environment variables
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"), os.Getenv("DB_PORT"),
	)
	fmt.Println("dsn: ", dsn)
	fmt.Println("22222")

	// Connect to database
	db, errr := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if errr != nil {
		log.Fatal("Failed to connect:", errr.Error())
	}
	fmt.Println("33333")

	fmt.Println("Connected using environment variables!")

	res := db.Exec(`CREATE TABLE IF NOT EXISTS users (id BIGSERIAL PRIMARY KEY,username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,password TEXT NOT NULL,full_name VARCHAR(255),phone_number VARCHAR(50),
    profile_picture_url TEXT,created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP,status VARCHAR(50) DEFAULT 'active',role VARCHAR(50) DEFAULT 'user',two_factor_enabled BOOLEAN DEFAULT FALSE,
    failed_login_attempts INTEGER DEFAULT 0,lock_status BOOLEAN DEFAULT FALSE,deleted_at TIMESTAMP NULL);
	`)
	if res.Error != nil {
		fmt.Printf("error in creating unexciting users table in database: %v", res.Error)
	}
	return db, nil
}

func DBC(db *gorm.DB) {
	// Get the raw SQL DB connection
	sqlDB, err := db.DB()
	if err != nil {
		log.Println("Error getting SQL database:", err)
		return
	}

	// Close the connection
	err = sqlDB.Close()
	if err != nil {
		log.Println("Error closing database connection:", err)
		return
	}
	// response after closing
	fmt.Println("Database connection closed successfully!")
}
