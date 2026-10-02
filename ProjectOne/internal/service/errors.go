package service

import "errors"

var ErrUserNotFound = errors.New("user to id not found")
var ErrInvalidInput = errors.New("invalid input")
var ErrEmailAlreadyExists = errors.New("this email already exists")
var ErrInvalidCredentials = errors.New("invalid email or password")
