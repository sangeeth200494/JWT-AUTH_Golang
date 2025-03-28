package helpers

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = []byte("secret-key")

func CreateToken(userID int64, username string) (string, error) {
	fmt.Println(1111111111111111)
	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Println(22222222222222222)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		fmt.Println(333333333333)
		return "", err
	}
	fmt.Println(4444444444444444)
	return tokenString, nil
}

func VerifyToken(tokenString string) error {
	fmt.Println("AAAAAAAAAAAAA")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		fmt.Println("BBBBBBBBBBBBBBB")
		return secretKey, nil
	})

	if err != nil {
		fmt.Println("CCCCCCCCCCCCCC")
		return err
	}

	if !token.Valid {
		fmt.Println("DDDDDDDDDDDDD")
		return fmt.Errorf("invalid token")
	}
	fmt.Println("EEEEEEEEE")
	return nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("error hashing password: %s", err.Error())
	}
	return string(hashedPassword), nil
}
