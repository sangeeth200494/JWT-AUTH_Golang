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

	response := map[string]string{"message": "Welcome Home..!"}
	json.NewEncoder(w).Encode(response)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var u models.User
	err := json.NewDecoder(r.Body).Decode(&u)
	//fmt.Printf("The user request value %v", u)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	user, err := userhandlers.GetUserByUsername(u.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "invalid username or password", http.StatusUnauthorized)
		} else {
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password)); err != nil {
		http.Error(w, "invalid username or password", http.StatusUnauthorized)
		return
	}

	token, err := helpers.CreateToken(user.ID, user.Username)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
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
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Missing authorization header")
		return
	}
	tokenString = tokenString[len("Bearer "):]

	err := helpers.VerifyToken(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Invalid token")
		return
	}

	fmt.Fprint(w, "Welcome to the the protected area")

}
