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
	fmt.Println("A11111111111111")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		fmt.Println("B111111111")
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println("C1111111111111")
			return
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		fmt.Println("D1111111111111")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				fmt.Println("E1111111111111")
				return nil, fmt.Errorf("unexpected signing method")
			}
			fmt.Println("F111111111111")
			return []byte("secret"), nil
		})

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println("G111111111111")
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println("H1111111111111")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println("I11111111111")
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println("J1111111111111")
			return
		}

		fmt.Println("K11111111111111")
		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
