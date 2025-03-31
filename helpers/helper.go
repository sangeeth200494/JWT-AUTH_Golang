package helpers

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("secret-key")

func CreateToken(userID uint64, username string) (string, error) {
	// adding claims into generating token
	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// creates and returns a complete, signed JWT
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil // returning the token string and nill error
}

func VerifyToken(tokenString string) error {
	// Parse parses, validates, verifies the signature and returns the parsed token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	// checking any error caused
	if err != nil {
		return err
	}

	// validating the token
	if !token.Valid {
		return fmt.Errorf("invalid token")
	}
	return nil // returning nil error
}

func ExtractUsernameFromToken(tokenString string) (string, error) {
	// Remove "Bearer
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	// If parsing fails, return an error
	if err != nil {
		return "", fmt.Errorf("error parsing token: %v", err)
	}

	// Extract claims and validate
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	// Extract username
	username, exists := claims["username"].(string)
	if !exists {
		return "", fmt.Errorf("username claim not found") // returning the empty string and caused error
	}
	return username, nil // returning the username and nil error
}
