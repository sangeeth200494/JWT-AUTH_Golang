package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/sangeeth200494/JWT-AUTH_Golang/models"
)

// ContextKey defines a custom type for context keys to prevent collisions
type ContextKey string

const (
	UserIDKey   ContextKey = "user_id"
	UsernameKey ContextKey = "username"
	RoleKey     ContextKey = "role"
	StatusKey   ContextKey = "status"
	ExpKey      ContextKey = "exp"
	RawTokenKey ContextKey = "raw_token"
)

func init() {
	godotenv.Load()
}

// AuthMiddleware validates JWT authentication from Authorization header OR secure HttpOnly cookie
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := ""

		// 1. Check Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			tokenString = strings.TrimSpace(tokenString)
		}

		// 2. Fall back to secure HttpOnly cookie
		if tokenString == "" {
			if cookie, err := r.Cookie("jwt_token"); err == nil && cookie.Value != "" {
				tokenString = strings.TrimSpace(cookie.Value)
			}
		}

		// If still missing, reject with 401
		if tokenString == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(&models.APIResponse{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized: No authentication token provided",
				Details: "Authentication required via Bearer token or session cookie",
			})
			return
		}

		// Parse and validate the JWT token using jwt/v5
		secret := []byte(os.Getenv("JWT_SECRET"))
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secret, nil
		})

		if err != nil || !token.Valid {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(&models.APIResponse{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized: Invalid or expired token",
				Details: "Token validation failed",
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(&models.APIResponse{
				Code:    http.StatusUnauthorized,
				Message: "Unauthorized: Invalid token claims",
			})
			return
		}

		// Check expiration
		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(&models.APIResponse{
					Code:    http.StatusUnauthorized,
					Message: "Unauthorized: Token has expired",
					Details: "Session expired, please sign in again",
				})
				return
			}
		}

		// Extract values safely
		userID, _ := claims["user_id"].(float64)
		username, _ := claims["username"].(string)
		createdAt, _ := claims["created_at"].(string)
		updatedAt, _ := claims["updated_at"].(string)
		lastLogin, _ := claims["last_login"].(string)
		status, _ := claims["status"].(string)
		role, _ := claims["role"].(string)

		// Chain context values
		ctx := context.WithValue(r.Context(), UserIDKey, uint64(userID))
		ctx = context.WithValue(ctx, UsernameKey, username)
		ctx = context.WithValue(ctx, RoleKey, role)
		ctx = context.WithValue(ctx, StatusKey, status)
		ctx = context.WithValue(ctx, RawTokenKey, tokenString)

		// Backward-compatible string keys
		ctx = context.WithValue(ctx, ContextKey("user_id"), userID)
		ctx = context.WithValue(ctx, ContextKey("username"), username)
		ctx = context.WithValue(ctx, ContextKey("created_at"), createdAt)
		ctx = context.WithValue(ctx, ContextKey("updated_at"), updatedAt)
		ctx = context.WithValue(ctx, ContextKey("last_login"), lastLogin)
		ctx = context.WithValue(ctx, ContextKey("status"), status)
		ctx = context.WithValue(ctx, ContextKey("role"), role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
