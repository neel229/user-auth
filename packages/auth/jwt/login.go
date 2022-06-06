package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Login checks if the credentials of the user are correct
func Login(ctx context.Context, mc *mongo.Client, input Creds) (string, error) {
	filter := bson.M{"email": input.Email}
	singleRes := mc.Database("auth").Collection("users").FindOne(ctx, filter)
	if singleRes.Err() != nil {
		if singleRes.Err() == mongo.ErrNoDocuments {
			return "", mongo.ErrNoDocuments
		}
		log.Printf("error fetching document: %v\n", singleRes.Err())
		return "", singleRes.Err()
	}
	var user Creds
	if err := singleRes.Decode(&user); err != nil {
		log.Println("error decoding user")
		return "", err
	}
	log.Println(user)
	match, err := Verify(input.Password, user.Password)
	if err != nil {
		log.Println(err)
		return "", ErrPassInvalid
	}
	if !match {
		return "", ErrPassInvalid
	}
	userId := user.UserId.String()
	token, err := CreateToken(userId)
	if err != nil {
		return "", err
	}
	return token, nil
}
