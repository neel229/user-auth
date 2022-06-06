package main

import "errors"

var (
	ErrInvalidInput   error = errors.New("Email or Password cannot be empty")
	ErrInternalServer error = errors.New("There was an error, try again after sometime")
	ErrEmailExists    error = errors.New("Email is already in use, try logging in...")
	ErrPassInvalid    error = errors.New("Password doesn't match, try again...")
)
