package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Creds struct {
	UserId   primitive.ObjectID `json:"_id" bson:"_id"`
	UserName string             `json:"username" bson:"username"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password" bson:"password"`
}

type Response struct {
	StatusCode int    `json:"statusCode,omitempty"`
	Body       string `json:"body,omitempty"`
}

// Main is the entrypoint of jwt auth function
func Main(input map[string]interface{}) *Response {
	// create mongo connection
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	client, err := NewMongoClient(ctx)
	if err != nil {
		log.Printf("error starting mongo connection: %v\n", err)
	}
	defer client.Disconnect(ctx)

	log.Println("input: ", input)

	// extract and check if the input is valid
	email, ok := input["email"].(string)
	if !ok {
		return &Response{
			StatusCode: http.StatusBadRequest,
			Body:       ErrInvalidInput.Error(),
		}
	}
	password, ok := input["password"].(string)
	if !ok {
		return &Response{
			StatusCode: http.StatusBadRequest,
			Body:       ErrInvalidInput.Error(),
		}
	}

	// check for login/register and call
	// corresponding function
	_, ok = input["login"].(string)
	if !ok {
		_, ok := input["register"].(string)
		if !ok {
			return &Response{
				StatusCode: http.StatusNotFound,
			}
		}
		username, ok := input["username"].(string)
		if !ok {
			return &Response{
				StatusCode: http.StatusBadRequest,
				Body:       ErrInvalidInput.Error(),
			}
		}
		user := Creds{
			Email:    email,
			Password: password,
			UserName: username,
		}
		token, err := Register(ctx, client, user)
		if err != nil {
			if err == ErrEmailExists {
				return &Response{
					StatusCode: http.StatusBadRequest,
					Body:       ErrEmailExists.Error(),
				}
			}
			return &Response{
				StatusCode: http.StatusInternalServerError,
				Body:       ErrInternalServer.Error(),
			}
		}
		return &Response{
			StatusCode: http.StatusOK,
			Body:       token,
		}
	}
	// call login function
	user := Creds{
		Email:    email,
		Password: password,
	}
	token, err := Login(ctx, client, user)
	if err != nil {
		return &Response{
			StatusCode: http.StatusInternalServerError,
			Body:       ErrInternalServer.Error(),
		}
	}
	return &Response{
		StatusCode: http.StatusOK,
		Body:       token,
	}
}
