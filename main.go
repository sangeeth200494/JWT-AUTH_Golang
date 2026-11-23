package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sangeeth200494/JWT-AUTH_Golang/helpers"
	"github.com/sangeeth200494/JWT-AUTH_Golang/login"
	"github.com/sangeeth200494/JWT-AUTH_Golang/middleware"
	userhandlers "github.com/sangeeth200494/JWT-AUTH_Golang/user-handlers"
)

// corsMiddleware adds CORS headers to allow cross-origin requests
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Cookie")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// demoTokenHandler provides a convenience endpoint to test JWT authentication even when PostgreSQL is offline
func demoTokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	now := time.Now()
	token, err := helpers.CreateToken(1, "demo_user", now, now, now, "active", "admin")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"code":    500,
			"message": "failed to generate demo token",
			"details": err.Error(),
		})
		return
	}

	// Set HttpOnly cookie for session testing
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt_token",
		Value:    token,
		Path:     "/",
		Expires:  now.Add(24 * time.Hour),
		MaxAge:   86400,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   r.TLS != nil,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"message": "demo token generated successfully",
		"token":   token,
		"user": map[string]interface{}{
			"id":       1,
			"username": "demo_user",
			"role":     "admin",
			"status":   "active",
		},
	})
}

func main() {
	router := mux.NewRouter()

	// Apply CORS middleware globally
	router.Use(corsMiddleware)

	// Public API Endpoints
	router.HandleFunc("/Home", login.Home).Methods("GET")
	router.HandleFunc("/login", login.LoginHandler).Methods("POST")
	router.HandleFunc("/logout", login.LogoutHandler).Methods("POST", "GET")
	router.HandleFunc("/api/v1/auth/logout", login.LogoutHandler).Methods("POST", "GET")
	router.HandleFunc("/users", userhandlers.RegisterUser).Methods("POST")
	router.HandleFunc("/users/{id:[0-9]+}", userhandlers.RetrieveUserByID).Methods("GET")
	router.HandleFunc("/demo-token", demoTokenHandler).Methods("GET", "POST")

	// Protected API Routes (Protected by AuthMiddleware)
	privateRouter := router.PathPrefix("/").Subrouter()
	privateRouter.Use(middleware.AuthMiddleware)
	privateRouter.HandleFunc("/protected", login.ProtectedHandler).Methods("GET")
	privateRouter.HandleFunc("/api/v1/auth/session", login.SessionHandler).Methods("GET")
	privateRouter.HandleFunc("/auth/session", login.SessionHandler).Methods("GET")

	// Static UI Server
	staticFileServer := http.FileServer(http.Dir("./static"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticFileServer))

	// Direct UI Page Routes for Browser
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	}).Methods("GET")

	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/login.html")
	}).Methods("GET")

	router.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/dashboard.html")
	}).Methods("GET")

	router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/register.html")
	}).Methods("GET")

	log.Println(":Listening on :4000")
	err := http.ListenAndServe("localhost:4000", router)
	if err != nil {
		fmt.Println("Could not start the server", err)
	}
}
