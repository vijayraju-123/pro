package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin" // Add this import for gin.Context
)

var jwtSecret = []byte("your-secret-key") // Replace with a secure secret from config

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

func GenerateJWT(userID uint) (string, error) {
	claims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}

func GetUserIDFromToken(c *gin.Context) uint {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	return userID.(uint)
}

func HashPassword(password string) string {
	// Implement password hashing (e.g., using bcrypt)
	// Example: return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// For now, return a placeholder
	return password // Replace with actual hashing logic
}