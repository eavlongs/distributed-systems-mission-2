package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/eavlongs/file_sync/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

var jwtKey = getJwtSecret()

type Claims struct {
	ID           uint   `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	DepartmentID string `json:"department_id"`
	jwt.RegisteredClaims
}

func getJwtSecret() []byte {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println(err)
		panic(err)
	}

	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		log.Fatalln("JWT secret not found, make sure JWT_SECRET is set in .env file")
	}

	return []byte(secret)
}

func GenerateJWT(user *models.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ParseJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, err
	}
	return claims, nil
}
