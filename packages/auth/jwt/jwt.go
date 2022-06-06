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

// VerityToken verifies the jwt token
func VerifyToken(tokenString string) (bool, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			log.Println("invalid signature header")
			return false, err
		}
		log.Printf("error parsing token: %v\n", err)
		return false, err
	}
	if !token.Valid {
		return false, nil
	}
	return true, nil
}
