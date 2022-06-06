// TODO: update error message
package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/argon2"
)

var (
	ErrInvalidInput   error = errors.New("Email or Password cannot be empty")
	ErrInternalServer error = errors.New("There was an error, try again after sometime")
	ErrEmailExists    error = errors.New("Email is already in use, try logging in...")
)

type PasswordConfig struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
}

type Response struct {
	Success      bool   `json:"success,omitempty"`
	StatusCode   int    `json:"statusCode,omitempty"`
	Body         string `json:"body,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

func Main(input map[string]interface{}) *Response {
	// create mongo connection
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("DATABASE_URL")))
	if err != nil {
		return &Response{
			StatusCode:   http.StatusInternalServerError,
			Success:      false,
			ErrorMessage: ErrInternalServer.Error(),
		}
	}
	defer client.Disconnect(ctx)
	db := client.Database("auth")
	coll := db.Collection("users")

	// check if the input is valid
	email, ok := input["email"].(string)
	if !ok {
		return &Response{
			StatusCode:   http.StatusBadRequest,
			Success:      false,
			ErrorMessage: ErrInvalidInput.Error(),
		}
	}
	password, ok := input["password"].(string)
	if !ok {
		return &Response{
			StatusCode:   http.StatusBadRequest,
			Success:      false,
			ErrorMessage: ErrInvalidInput.Error(),
		}
	}

	// check if the email already exists
	filter := bson.M{"email": email}
	singleRes := coll.FindOne(ctx, filter)
	if singleRes.Err() != mongo.ErrNoDocuments {
		return &Response{
			StatusCode:   http.StatusBadRequest,
			Success:      false,
			ErrorMessage: ErrEmailExists.Error(),
		}
	}

	// hash the password before saving
	pc := &PasswordConfig{
		time:    1,
		memory:  64 * 1024,
		threads: 4,
		keyLen:  32,
	}
	hashedPass, err := GeneratePassword(pc, password)
	if err != nil {
		return &Response{
			StatusCode:   http.StatusInternalServerError,
			Success:      false,
			ErrorMessage: ErrInternalServer.Error(),
		}
	}

	// saved the username and hashed password
	user := bson.D{
		{Key: "email", Value: email},
		{Key: "password", Value: hashedPass},
	}
	insertOneRes, err := coll.InsertOne(ctx, user)
	if err != nil {
		return &Response{
			StatusCode:   http.StatusInternalServerError,
			Success:      false,
			ErrorMessage: ErrInternalServer.Error(),
		}
	}
	insertedId := fmt.Sprint(insertOneRes.InsertedID.(primitive.ObjectID))
	// return the InsertedID
	return &Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Body:       insertedId,
	}
}

func GeneratePassword(c *PasswordConfig, password string) (string, error) {

	// Generate a Salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, c.time, c.memory, c.threads, c.keyLen)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	full := fmt.Sprintf(format, argon2.Version, c.memory, c.time, c.threads, b64Salt, b64Hash)
	return full, nil
}
