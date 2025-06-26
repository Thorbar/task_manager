package auth

import (
	"errors"
	"os"
	"time"

	"task-manager/backend-go/internal/i18n"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

var jwtKey []byte

// Initialize JWT secret key from environment variables
func init() {
	_ = godotenv.Load()
	jwtKey = []byte(os.Getenv("JWT_SECRET"))
}

// GenerateJWT creates a signed JWT token containing the user ID and expiration
func GenerateJWT(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// ParseToken validates the JWT token string and extracts the user ID
func ParseToken(tokenStr string) (int, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return 0, errors.New(i18n.T("error.invalid_token"))
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New(i18n.T("error.invalid_token_claims"))
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New(i18n.T("error.invalid_token_claims"))
	}

	return int(userIDFloat), nil
}
