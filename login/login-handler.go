package login

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sangeeth200494/JWT-AUTH_Golang/helpers"
	"github.com/sangeeth200494/JWT-AUTH_Golang/models"
	userhandlers "github.com/sangeeth200494/JWT-AUTH_Golang/user-handlers"
	"golang.org/x/crypto/bcrypt"
)

func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(&models.APIResponse{Code: http.StatusOK, Message: "Welcome Home..!",
		Details: "You Are Entered Into A Server"})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var u models.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		json.NewEncoder(w).Encode(&models.APIResponse{Code: 500, Message: "invalid request payload", Details: err.Error()})
		return
	}

	user, err := userhandlers.GetUserByUsername(u.Username, u.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			json.NewEncoder(w).Encode(&models.APIResponse{Code: 401, Message: "invalid username or password"})
			return
		} else {
			json.NewEncoder(w).Encode(&models.APIResponse{Code: 500, Message: "server error", Details: err.Error()})
			return
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password)); err != nil {
		json.NewEncoder(w).Encode(&models.APIResponse{Code: 401, Message: "invalid username or password", Details: err.Error()})
		return
	}

	token, err := helpers.CreateToken(user.ID, user.Username)
	if err != nil {
		json.NewEncoder(w).Encode(&models.APIResponse{Code: 500, Message: "server error", Details: err.Error()})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		json.NewEncoder(w).Encode(&models.APIResponse{Code: 401, Message: "Missing authorization header", Details: tokenString})
		// w.WriteHeader(http.StatusUnauthorized)
		// fmt.Fprint(w, "Missing authorization header")
		return
	}
	tokenString = tokenString[len("Bearer "):]

	err := helpers.VerifyToken(tokenString)
	if err != nil {
		json.NewEncoder(w).Encode(&models.APIResponse{Code: 401, Message: "Invalid token", Details: err.Error()})
		// w.WriteHeader(http.StatusUnauthorized)
		// fmt.Fprint(w, "Invalid token")
		return
	}
	json.NewEncoder(w).Encode(&models.APIResponse{Code: http.StatusOK, Message: "Welcome to the the protected area"})
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %s", err.Error())
	}
	return string(hashedPassword), nil
}

// CheckPassword compares a hashed password with a plain password
func CheckPassword(hashedPassword, plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil // Returns true if the password matches
}
