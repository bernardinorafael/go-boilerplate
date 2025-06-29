package httputil

import "errors"

// TODO: use sentinel errors to send custom errors on decode body erros

var (
	ErrUnknownRequestBodyKey = errors.New("request body contains unknown key")
	ErrEmptyRequestBody      = errors.New("body cannot be empty")
)
