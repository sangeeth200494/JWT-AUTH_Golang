package login

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sangeeth200494/JWT-AUTH_Golang/database"
	"github.com/sangeeth200494/JWT-AUTH_Golang/helpers"
	"github.com/sangeeth200494/JWT-AUTH_Golang/middleware"
	"github.com/sangeeth200494/JWT-AUTH_Golang/models"
	"gorm.io/gorm"
)

// Home is a public health / welcome endpoint
func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&models.APIResponse{
		Code:    http.StatusOK,
		Message: "Welcome Home..!",
		Details: "You Are Entered Into A Server",
	})
}

// LoginHandler authenticates the user, sets a secure HttpOnly cookie, and returns safe user session data
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var u models.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(&models.APIResponse{Code: 400, Message: "invalid request payload", Details: err.Error()})
		return
	}

	// connect to database
	db, err := database.DBConnection()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&models.APIResponse{Code: 500, Message: "database connection failed", Details: err.Error()})
		return
	}
	defer database.DBC(db)

	// check user existence by email and password
	user, err := helpers.GetUserByUserEmail(u.Email, u.Password, db)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(err.Error(), "record not found") || strings.Contains(err.Error(), "hashed password") {
			json.NewEncoder(w).Encode(&models.APIResponse{Code: 401, Message: "invalid email or password", Details: "Invalid email or password"})
			return
		}
		json.NewEncoder(w).Encode(&models.APIResponse{Code: 401, Message: "authentication failed", Details: "Invalid credentials"})
		return
	}

	// generate token
	token, err := helpers.CreateToken(user.ID, user.Username, user.CreatedAt, user.UpdatedAt, user.LastLogin, user.Status, user.Role)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&models.APIResponse{Code: 500, Message: "server error generating token", Details: err.Error()})
		return
	}

	// Set secure HttpOnly cookie for production session management
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt_token",
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		MaxAge:   86400,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   r.TLS != nil,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"message": "login successful",
		"token":   token, // Provided for external API clients; web UI uses HttpOnly cookie
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
			"status":   user.Status,
		},
	})
}

// LogoutHandler clears the HttpOnly authentication cookie
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Clear the jwt_token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   r.TLS != nil,
	})

	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    200,
		"message": "logged out successfully",
	})
}

// SessionHandler returns dynamic, safe session metadata and a masked token representation.
// CRITICAL: Never exposes the full raw JWT token.
func SessionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Extract token from context (set by AuthMiddleware) or cookie or header
	rawToken := ""
	if ctxToken, ok := r.Context().Value(middleware.RawTokenKey).(string); ok && ctxToken != "" {
		rawToken = ctxToken
	}

	if rawToken == "" {
		// Fallback check cookie
		if cookie, err := r.Cookie("jwt_token"); err == nil && cookie.Value != "" {
			rawToken = strings.TrimSpace(cookie.Value)
		}
	}

	if rawToken == "" {
		// Fallback check Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			rawToken = strings.TrimPrefix(authHeader, "Bearer ")
			rawToken = strings.TrimSpace(rawToken)
		}
	}

	if rawToken == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"authenticated":  false,
			"session_status": "inactive",
			"token_status":   "Missing",
			"message":        "No active session found",
		})
		return
	}

	// Verify token validity and claims
	claims, err := helpers.ExtractClaimsFromToken(rawToken)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"authenticated":  false,
			"session_status": "expired",
			"token_status":   "Invalid or Expired",
			"message":        err.Error(),
		})
		return
	}

	// Compute expiration
	var expTime time.Time
	sessionStatus := "active"
	tokenStatus := "Valid"

	if expFloat, ok := claims["exp"].(float64); ok {
		expTime = time.Unix(int64(expFloat), 0).UTC()
		if expTime.Before(time.Now().UTC()) {
			sessionStatus = "expired"
			tokenStatus = "Expired"
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"authenticated":  false,
				"session_status": sessionStatus,
				"token_status":   tokenStatus,
				"expires_at":     expTime.Format(time.RFC3339),
			})
			return
		}
	}

	// Build masked token: e.g. eyJhbGciOiJIUzUx...••••••••••••••••
	maskedToken := helpers.MaskToken(rawToken)

	userID, _ := claims["user_id"].(float64)
	username, _ := claims["username"].(string)
	role, _ := claims["role"].(string)
	status, _ := claims["status"].(string)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"authenticated":  true,
		"token_type":     "Bearer JWT",
		"session_status": sessionStatus,
		"token_status":   tokenStatus,
		"masked_token":   maskedToken, // Safely masked, never raw!
		"expires_at":     expTime.Format(time.RFC3339),
		"user": map[string]interface{}{
			"id":       uint64(userID),
			"username": username,
			"role":     role,
			"status":   status,
		},
	})
}

// ProtectedHandler handles access to the protected area
func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get username from context (set by AuthMiddleware)
	username := ""
	if ctxUser, ok := r.Context().Value(middleware.UsernameKey).(string); ok && ctxUser != "" {
		username = ctxUser
	}

	if username == "" {
		// Fallback to extract from token
		rawToken := ""
		if ctxToken, ok := r.Context().Value(middleware.RawTokenKey).(string); ok {
			rawToken = ctxToken
		}
		if rawToken != "" {
			if u, err := helpers.ExtractUsernameFromToken(rawToken); err == nil {
				username, _ = u.(string)
			}
		}
	}

	if username == "" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(&models.APIResponse{Code: http.StatusUnauthorized, Message: "unauthorized: invalid session", Details: ""})
		return
	}

	json.NewEncoder(w).Encode(&models.APIResponse{
		Code:    http.StatusOK,
		Message: "Welcome to the protected area",
		Details: username,
	})
}

func VerifyUser(UserName string, db *gorm.DB) (*models.User, error) {
	var user models.User
	res := db.Raw("SELECT username FROM users WHERE username = ?", UserName).Scan(&user)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("username doesn't exist")
	}
	if res.Error != nil {
		return nil, fmt.Errorf("error in checking user existence: %s", res.Error)
	}
	return &user, nil
}
