package models

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID       uint64 `gorm:"primaryKey" json:"user_id"`
	Username string `gorm:"size:100;not null" json:"username"`
	Password string `gorm:"unique;not null" json:"password"`
}

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details"`
}

func (e *APIResponse) Error() string {
	return fmt.Sprintf("Error%d:%s", e.Code, e.Message)
}

func HashPassword(password string) (string, error) {
	// GenerateFromPassword returns the bcrypt hash of the password at the given cost.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Error generate password: %s", err.Error()) // returning the empty string and caused error
	}
	return string(hashedPassword), nil // returning the hashed password and nil error
}

func ValidatePasswords(providedPassword, storedHashedPassword string) error {
	// comparing the both password
	errr := bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(providedPassword))
	if errr != nil {
		return fmt.Errorf("error in checking hashed password: %s", errr.Error())
	}
	return nil // returning nil error
}

func CheckPasswordExistence(username, newPassword string, db *gorm.DB) (bool, error) {
	var user User

	// GenerateFromPassword returns the bcrypt hash of the password at the given cost.
	bytes, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return false, fmt.Errorf("Error generate password: %s", err.Error())
	}

	// Fetch user from the database
	errr := db.Where("username = ? AND password = ?", username, bytes).First(&user).Error
	if errr != nil {
		if errors.Is(errr, gorm.ErrRecordNotFound) {
			return false, errors.New("user not found") // returning error user not found
		}
		return false, fmt.Errorf("error in retrieving user by username and password %s", errr.Error()) // returning error
	}
	return true, nil //returning the password existence
}
