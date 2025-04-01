package models

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID                  uint64     `json:"id" gorm:"primaryKey"`
	Username            string     `json:"username" gorm:"unique;not null"`
	Email               string     `json:"email" gorm:"unique;not null"`
	Password            string     `json:"Password" gorm:"not null"`
	FullName            string     `json:"full_name"`
	PhoneNumber         string     `json:"phone_number"`
	ProfilePictureURL   string     `json:"profile_picture_url"`
	CreatedAt           time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt           time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
	LastLogin           time.Time  `json:"last_login"`
	Status              string     `json:"status" gorm:"default:'active'"`
	Role                string     `json:"role" gorm:"default:'user'"`
	TwoFactorEnabled    bool       `json:"two_factor_enabled" gorm:"default:false"`
	FailedLoginAttempts int        `json:"failed_login_attempts" gorm:"default:0"`
	LockStatus          bool       `json:"lock_status" gorm:"default:false"`
	DeletedAt           *time.Time `json:"deleted_at"`
}

// type UserLogin struct {
// 	ID       uint64 `gorm:"primaryKey" json:"user_id"`
// 	Username string `gorm:"size:100;not null" json:"username"`
// 	Password string `gorm:"unique;not null" json:"password"`
// }

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
