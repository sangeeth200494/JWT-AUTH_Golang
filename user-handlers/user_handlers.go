package userhandlers

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/sangeeth200494/JWT-AUTH_Golang/database"
	"github.com/sangeeth200494/JWT-AUTH_Golang/helpers"
	"github.com/sangeeth200494/JWT-AUTH_Golang/models"
)

func createUserR(user models.User) error {
	godotenv.Load()
	db, err := database.DBConnection()
	if err != nil {
		return fmt.Errorf("error in connecting database: %s", err.Error)
	}
	//defer db.Close()

	// Hash the password using bcrypt
	hashedPassword, err := helpers.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("error hashing password: %s", err.Error)
	}

	err = db.Exec("INSERT INTO users (id, username, password) VALUES ($1, $2, $3)", user.ID, user.Username, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("error inserting user: %s", err.Error)
	}
	return nil
}

func GetUserByUsername(username string, password string) (*models.User, error) {
	var user models.User

	db, err := database.DBConnection()
	if err != nil {
		log.Println("error in connecting database", err.Error)
		return nil, err.Error
	}
	//defer database.Close(database.DBC())

	//var db *gorm.DB
	row := db.Raw("SELECT id, username, password FROM users WHERE username = ? AND password =?", username, password)
	errr := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return User{}, err
	}
	// Check for errors
	if row.Error != nil {
		log.Println("Error fetching user:", row.Error)
		return nil, row.Error
	}

	// Check if the user exists
	if user.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}

	return &user, nil
}
