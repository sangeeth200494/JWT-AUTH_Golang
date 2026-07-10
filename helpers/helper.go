package helpers

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/sangeeth200494/JWT-AUTH_Golang/models"
	"gorm.io/gorm"
)

// Load environment variables from .env file
func init() {
	godotenv.Load()
}

// var secretKey = []byte("secret-key")
// var secretKey = []byte(os.Getenv("JWT_SECRET"))

func CreateToken(userID uint64, username string, createdAt time.Time, updatedAt time.Time, lastLogin time.Time, status string, role string) (string, error) {
	// adding claims into generating token
	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["username"] = username
	claims["created_at"] = createdAt
	claims["updated_at"] = updatedAt
	claims["last_login"] = lastLogin
	claims["status"] = status
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	// creates and returns a complete, signed JWT
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil // returning the token string and nill error
}

func VerifyToken(tokenString string) error {

	// Remove "Bearer
	TokenString := strings.TrimPrefix(tokenString, "Bearer ")

	// Parse the token, validate the signature and retrieve the token claims
	token, err := jwt.Parse(TokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure that the token's signing method is correct (HS512)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	// Handle error while parsing token (invalid format, signature verification failed, etc.)
	if err != nil {
		return fmt.Errorf("failed to parse token: %v", err)
	}

	// Validate the token's validity
	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	// Additional check for token expiration (optional, depending on your needs)
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if exp, ok := claims["exp"].(float64); ok {
			if exp < float64(time.Now().Unix()) {
				return fmt.Errorf("token has expired")
			}
		}
	}
	return nil // Return nil if the token is valid
}

func ExtractUsernameFromToken(tokenString string) (interface{}, error) {
	// Remove "Bearer
	TokenString := strings.TrimPrefix(tokenString, "Bearer ")

	// Parse the token
	token, err := jwt.Parse(TokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
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

	// Optionally check if the token has expired
	if exp, ok := claims["exp"].(float64); ok {
		if int64(exp) < time.Now().Unix() {
			return "", fmt.Errorf("token has expired")
		}
	}

	// Extract username
	username, exists := claims["username"].(string)
	if !exists {
		return "", fmt.Errorf("username claim not found") // returning the empty string and caused error
	}
	return username, nil // returning the username and nil error
}

/////////////////////////////--------Get User helper function--------/////////////////////////////

func GetUserByUserEmail(email string, password string, db *gorm.DB) (*models.User, error) {
	var user models.User
	fmt.Println("aaaaaa")

	// parsing stored hashed password of a user by given username
	StoredHashed, err := GetStoredPassword(db, email)
	if err != nil {
		return nil, fmt.Errorf("error getting stored hashed password from db: %s", err.Error())
	}
	fmt.Println("bbbbbb")

	// validating the password with user input password
	errr := models.ValidatePasswords(password, StoredHashed)
	if errr != nil {
		return nil, fmt.Errorf("error in validating user: %s", errr.Error())
	}
	fmt.Println("cccccc")

	//var db *gorm.DB
	// row := db.Raw("SELECT id, username, email, password FROM users WHERE email = ?", email).Scan(&user)
	// if row.Error != nil {
	// 	fmt.Println("error is error")
	// 	return nil, fmt.Errorf("error in retrieving user details: %s", row.Error)
	// }
	// fmt.Println("ddddddd")

	// // checking the any error is caused
	// if row.Error != nil {
	// 	return nil, fmt.Errorf("user not found: %s", row.Error)
	// }
	fmt.Println("eeeeeee")
	return &user, nil // returning the user
}

func GetStoredPassword(db *gorm.DB, email string) (string, error) {
	var user models.User
	fmt.Println("email is: ", email)
	// retrieving user from database using given email
	result := db.Where("email = ?", email).First(&user)
	fmt.Println("result is :", result)
	if result.Error != nil {
		return "", result.Error // Return error if user not found
	}
	return user.Password, nil // Return the stored hashed password
}
