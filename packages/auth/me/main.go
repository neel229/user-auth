package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

func Main(input map[string]interface{}) *Response {
	// create mongo connection
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	client, err := NewMongoClient(ctx)
	if err != nil {
		log.Printf("error starting mongo connection: %v\n", err)
	}
	defer client.Disconnect(ctx)

	// get the value of jwt token and verify it
	jwt, ok := input["__ow_headers"].(map[string]interface{})["authorization"].(string)
	if !ok {
		return &Response{
			StatusCode: http.StatusUnauthorized,
		}
	}
	claims, err := VerifyToken(jwt)
	if err != nil {
		return &Response{
			StatusCode: http.StatusUnauthorized,
		}
	}
	if len(claims.UserID) == 0 {
		return &Response{
			StatusCode: http.StatusUnauthorized,
		}
	}

	// fetch account details
	var user Creds
	id, _ := primitive.ObjectIDFromHex(claims.UserID)
	filter := bson.M{"_id": id}
	singleRes := client.Database("auth").Collection("users").FindOne(ctx, filter)
	if err := singleRes.Decode(&user); err != nil {
		log.Printf("error decoding result: %v\n", err)
		return &Response{
			StatusCode: http.StatusInternalServerError,
		}
	}
	return &Response{
		StatusCode: http.StatusOK,
		Body:       user.UserName,
	}
}
