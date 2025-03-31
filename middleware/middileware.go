package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

// JWT Middleware for authentication
func AuthMiddleware(next http.Handler) http.Handler {
	// returning HandlerFunc type is an adapter to allow the use of ordinary functions as HTTP handlers.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// getting the token string from header
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Replace returns a copy of the string s with the first n non-overlapping instances of old replaced by new
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		// Parse, validate, and return a token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte("secret"), nil // returning the byte values of token and nil error
		})

		// checking any error caused
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// validating the token
		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Claims that uses the map[string]interface{} for JSON decoding This is the default claims type if you don't supply one
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// giving value into userid variable
		userID, ok := claims["user_id"].(string)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// returns a copy of parent in which the value associated with key is val.
		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
