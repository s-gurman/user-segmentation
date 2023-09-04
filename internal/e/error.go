package e

import (
	"fmt"
	"net/http"
)

type CustomError interface {
	error

	Message() string
	Code() int
}

type customError struct {
	message string
	source  string
	code    int
}

// Implementing also default error interface for type conversion.
func (e customError) Error() string {
	return fmt.Sprintf("%s err: %s", e.source, e.message)
}

func (e customError) Message() string {
	return e.message
}

func (e customError) Code() int {
	return e.code
}

func NewBadRequest(msg, from string) CustomError {
	return customError{
		message: msg,
		source:  from,
		code:    http.StatusBadRequest,
	}
}

func NewNotFound(msg, from string) CustomError {
	return customError{
		message: msg,
		source:  from,
		code:    http.StatusNotFound,
	}
}
