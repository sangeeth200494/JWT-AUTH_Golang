package userhandlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/sangeeth200494/JWT-AUTH_Golang/database"
	"github.com/sangeeth200494/JWT-AUTH_Golang/models"
	"gorm.io/gorm"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	// It informs the client (browser, Postman, frontend app, etc.) that the response body will be in JSON format.
	w.Header().Set("Content-Type", "application/json")
	var user models.User
	godotenv.Load()

	//binding the req body
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		json.NewEncoder(w).Encode(&models.APIResponse{Code: 500, Message: "invalid request body", Details: err.Error()})
		return
	}

	// connecting database
	db, err := database.DBConnection()
	if err != nil {
		json.NewEncoder(w).Encode(&models.APIResponse{Code: 500, Message: "database connection failed", Details: err.Error()})
		return
	}
	defer database.DBC(db)

	// Hash the password using bcrypt
	HashedPASSWORD, errr := models.HashPassword(user.Password)
	if errr != nil {
		json.NewEncoder(w).Encode(&models.APIResponse{Code: 500, Message: "error hashing password: %s", Details: errr.Error()})
		return
	}
	user.Password = HashedPASSWORD

	// creating or inserting user details into db
	result := db.Create(&user)
	if result.Error != nil {
		json.NewEncoder(w).Encode(&models.APIResponse{Code: 400, Message: "registering user failed", Details: result.Error})
		return
	}
	// success response with inserted user_id
	json.NewEncoder(w).Encode(&models.APIResponse{Code: http.StatusCreated, Message: "user registered successfully", Details: user.ID})
}

func GetUserByUsername(username string, password string, db *gorm.DB) (*models.User, error) {
	var user models.User

	// parsing stored hashed password of a user by given username
	StoredHashed, err := GetStoredPassword(db, username)
	if err != nil {
		return nil, fmt.Errorf("error getting stored hashed password from db: %s", err.Error())
	}

	// validating the password with user input password
	errr := models.ValidatePasswords(password, StoredHashed)
	if errr != nil {
		return nil, fmt.Errorf("error in registering user: %s", errr.Error())
	}

	//var db *gorm.DB
	row := db.Raw("SELECT id, username, password FROM users WHERE username = ?", username).Scan(&user)
	if row.Error != nil {
		return nil, fmt.Errorf("error in retrieving user details: %s", row.Error)
	}

	// checking the any error is caused
	if row.Error != nil {
		return nil, fmt.Errorf("user not found: %s", row.Error)
	}
	return &user, nil // returning the user
}

func GetStoredPassword(db *gorm.DB, username string) (string, error) {
	var user models.User
	// retrieving user from database using given username
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return "", result.Error // Return error if user not found
	}
	return user.Password, nil // Return the stored hashed password
}
