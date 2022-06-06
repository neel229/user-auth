package main

import (
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserID string `json:"userId"`
	jwt.StandardClaims
}

var secret = []byte(os.Getenv("JWT_KEY"))

// CreateToken creates a new token
func CreateToken(userId string) (string, error) {
	expiresIn := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		UserID: userId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresIn.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		log.Printf("error signing jwt token: %v\n", err)
		return "", err
	}
	return tokenString, nil
}
