package userhandlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sangeeth200494/JWT-AUTH_Golang/database"
	"github.com/sangeeth200494/JWT-AUTH_Golang/models"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user models.User
	godotenv.Load()

	// binding the req body
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		json.NewEncoder(w).Encode(&models.APIResponse{Code: 400, Message: "invalid request body", Details: err.Error()})
		return
	}

	if strings.TrimSpace(user.Username) == "" || strings.TrimSpace(user.Email) == "" || strings.TrimSpace(user.Password) == "" {
		json.NewEncoder(w).Encode(&models.APIResponse{Code: 400, Message: "validation failed", Details: "Username, email and password are required."})
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
		json.NewEncoder(w).Encode(&models.APIResponse{Code: 500, Message: "error hashing password", Details: errr.Error()})
		return
	}
	user.Password = HashedPASSWORD

	// creating or inserting user details into db
	result := db.Create(&user)
	if result.Error != nil {
		rawErr := result.Error.Error()
		friendlyMsg := rawErr
		if strings.Contains(rawErr, "users_username_key") {
			friendlyMsg = "Username '" + user.Username + "' is already registered. Please choose another username."
		} else if strings.Contains(rawErr, "users_email_key") {
			friendlyMsg = "Email '" + user.Email + "' is already registered. Please sign in or use another email."
		}
		json.NewEncoder(w).Encode(&models.APIResponse{Code: 400, Message: "registering user failed", Details: friendlyMsg})
		return
	}

	// success response with inserted user_id
	json.NewEncoder(w).Encode(&models.APIResponse{Code: 201, Message: "user registered successfully", Details: user.ID})
}
