package userhandlers

import (
	"fmt"
	"log"

	"github.com/sangeeth200494/JWT-AUTH_Golang/models"
	"gorm.io/gorm"
)

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	var db *gorm.DB
	row := db.Raw("SELECT id, username, password FROM users WHERE username = ?", username)
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
