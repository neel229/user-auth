package main

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Register creates a new user entry in the db
func Register(ctx context.Context, mc *mongo.Client, input Input) (string, error) {
	// check if the email already exists
	filter := bson.M{"email": input.Email}
	singleRes := mc.Database("auth").Collection("users").FindOne(ctx, filter)
	if singleRes.Err() != mongo.ErrNoDocuments {
		log.Printf("error fetching document: %v\n", singleRes.Err())
		return "", singleRes.Err()
	}

	if singleRes.Err() == mongo.ErrNoDocuments {
		// hash the password before storing
		hashedPass, err := Hash(input.Password)
		if err != nil {
			return "", err
		}

		// create a new record of the user
		user := bson.D{
			{Key: "email", Value: input.Email},
			{Key: "password", Value: hashedPass},
		}
		insertOneRes, err := mc.Database("auth").Collection("users").InsertOne(ctx, user)
		if err != nil {
			return "", err
		}
		_ = insertOneRes
	}
	// TODO: generate jwt token
	return "", nil
}
