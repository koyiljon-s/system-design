package jwt

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// GenerateToken creates a new JWT token for a user ID
func GenerateToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID.String(), 
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenString string) (uuid.UUID, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, jwt.ErrTokenInvalidClaims
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return uuid.Nil, errors.New("invalid user_id in token")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}
