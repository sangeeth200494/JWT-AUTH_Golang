package userhandlers

import (
	"encoding/json"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/sangeeth200494/JWT-AUTH_Golang/database"
	"github.com/sangeeth200494/JWT-AUTH_Golang/models"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	// It informs the client (browser, Postman, frontend app, etc.) that the response body will be in JSON format.
	w.Header().Set("Content-Type", "application/json")
	var user models.User
	godotenv.Load()

	//binding the req body
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		json.NewEncoder(w).Encode(&models.APIResponse{Code: 400, Message: "invalid request body", Details: err.Error()})
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
	json.NewEncoder(w).Encode(&models.APIResponse{Code: 201, Message: "user registered successfully", Details: user.ID})
}
