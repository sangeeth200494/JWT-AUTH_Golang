package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

// Define a custom type for context keys to prevent collisions
type contextKey string

func init() {
	godotenv.Load()
}

// var secretKey = []byte("secret-key")
var secretKey = []byte(os.Getenv("JWT_SECRET"))

// JWT Middleware for authentication
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from the Authorization header
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
			return
		}

		// Remove "Bearer " prefix
		TokenString := strings.TrimPrefix(tokenString, "Bearer ")

		// Parse and validate the JWT token
		token, err := jwt.Parse(TokenString, func(token *jwt.Token) (interface{}, error) {
			// Ensure that the token's signing method is correct (HS512)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized: Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Extract claims from token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Unauthorized: Invalid claims", http.StatusUnauthorized)
			return
		}

		// Extract values safely
		userID, _ := claims["user_id"].(string)
		userName, _ := claims["username"].(string)
		createdAt, _ := claims["created_at"].(string)
		updatedAt, _ := claims["updated_at"].(string)
		lastLogin, _ := claims["last_login"].(string) // Fixed inconsistent naming
		status, _ := claims["status"].(string)
		role, _ := claims["role"].(string)
		exP, _ := claims["exp"].(string)

		// Chain context values properly
		ctx := context.WithValue(r.Context(), contextKey("user_id"), userID)
		ctx = context.WithValue(ctx, contextKey("username"), userName)
		ctx = context.WithValue(ctx, contextKey("created_at"), createdAt)
		ctx = context.WithValue(ctx, contextKey("updated_at"), updatedAt)
		ctx = context.WithValue(ctx, contextKey("last_login"), lastLogin)
		ctx = context.WithValue(ctx, contextKey("status"), status)
		ctx = context.WithValue(ctx, contextKey("role"), role)
		ctx = context.WithValue(ctx, contextKey("exp"), exP)

		// Pass the request with the modified context to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
