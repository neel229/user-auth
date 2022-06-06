package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	UserID string `json:"userId"`
	jwt.StandardClaims
}

var secret = []byte(os.Getenv("JWT_KEY"))

// VerityToken verifies the jwt token
func VerifyToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			log.Println("invalid signature header")
			return &Claims{}, err
		}
		log.Printf("error parsing token: %v\n", err)
		return &Claims{}, err
	}
	if !token.Valid {
		return &Claims{}, nil
	}
	fmt.Printf("token.Claims: %v\n", token.Claims)
	return claims, nil
}
