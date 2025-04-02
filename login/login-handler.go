package login

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/sangeeth200494/JWT-AUTH_Golang/database"
	"github.com/sangeeth200494/JWT-AUTH_Golang/helpers"
	"github.com/sangeeth200494/JWT-AUTH_Golang/models"
	userhandlers "github.com/sangeeth200494/JWT-AUTH_Golang/user-handlers"
	"gorm.io/gorm"
)

// sample function
func Home(w http.ResponseWriter, r *http.Request) {
	// It informs the client (browser, Postman, frontend app, etc.) that the response body will be in JSON format.
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(&models.APIResponse{Code: http.StatusOK, Message: "Welcome Home..!",
		Details: "You Are Entered Into A Server"})
}

// function for login the user and user will get a token after successful login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// It informs the client (browser, Postman, frontend app, etc.) that the response body will be in JSON format.
	w.Header().Set("Content-Type", "application/json")

	// binding the request body
	var u models.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		json.NewEncoder(w).Encode(&models.APIResponse{Code: 500, Message: "invalid request payload", Details: err.Error()})
		return
	}

	// parsing db connection into a variable
	db, err := database.DBConnection()
	if err != nil {
		json.NewEncoder(w).Encode(&models.APIResponse{Code: 500, Message: "database connection failed", Details: err.Error()})
		return
	}
	defer database.DBC(db)

	// calling function to check the existence of a user by username and password
	user, err := userhandlers.GetUserByUsername(u.Username, u.Password, db)
	if err != nil {
		// comparing the error with gorm error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			json.NewEncoder(w).Encode(&models.APIResponse{Code: 401, Message: "invalid username or password"})
			return
		} else {
			json.NewEncoder(w).Encode(&models.APIResponse{Code: 500, Message: "server error", Details: err.Error()})
			return
		}
	}

	// generate token
	token, err := helpers.CreateToken(user.ID, user.Username, user.CreatedAt, user.UpdatedAt, user.LastLogin, user.Status, user.Role)
	if err != nil {
		json.NewEncoder(w).Encode(&models.APIResponse{Code: 500, Message: "server error", Details: err.Error()})
		return
	}
	// returns a new encoder that writes to w.
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	// It informs the client (browser, Postman, frontend app, etc.) that the response body will be in JSON format.
	w.Header().Set("Content-Type", "application/json")

	// getting token from header
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		json.NewEncoder(w).Encode(&models.APIResponse{Code: 401, Message: "Missing authorization header", Details: tokenString})
		return
	}
	tokenString = tokenString[len("Bearer "):]

	username, _ := helpers.ExtractUsernameFromToken(tokenString)

	// validating the token
	errr := helpers.VerifyToken(tokenString)
	if errr != nil {
		json.NewEncoder(w).Encode(&models.APIResponse{Code: 401, Message: "Invalid token in verification", Details: errr.Error()})
		return
	}

	// success response
	json.NewEncoder(w).Encode(&models.APIResponse{Code: http.StatusOK, Message: "Welcome to the the protected area", Details: username})
}

func VerifyUser(UserName string, db *gorm.DB) (*models.User, error) {
	var user models.User
	// querying to get a user name from db
	res := db.Raw("SELECT username FROM users WHERE username = ?", UserName).Scan(&user)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("username doesn't exist")
	}

	// checking any error caused
	if res.Error != nil {
		return nil, fmt.Errorf("error in checking user existence: %s", res.Error)
	}
	return &user, nil // returning the specified user and nil error
}
