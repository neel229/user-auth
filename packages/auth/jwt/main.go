// TODO: Add JWT authentication
package main

import (
	"context"
	"log"
	"time"
)

type Input struct {
	Email    string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
}

type Response struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

// Main is the entrypoint of jwt auth function
func Main(input map[string]interface{}) {
	// create mongo connection
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	client, err := NewMongoClient(ctx)
	if err != nil {
		log.Printf("error starting mongo connection: %v\n", err)
	}
	defer client.Disconnect(ctx)

	// extract and check if the input is valid
	email, ok := input["email"].(string)
	if !ok {
		return
	}
	password, ok := input["password"].(string)
	if !ok {
		return
	}

	user := Input{
		Email:    email,
		Password: password,
	}

	// check for login/register and call
	// corresponding function
	_, ok = input["login"].(string)
	if !ok {
		_, ok := input["register"].(string)
		if !ok {
			return
		}
		_, err := Register(ctx, client, user)
		if err != nil {
			return
		}
	}
	// call login function
	_, err = Login(ctx, client, user)
	if err != nil {
		return
	}
}
