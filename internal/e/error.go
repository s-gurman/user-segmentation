package e

import (
	"fmt"
	"net/http"
)

type CustomError struct {
	message string
	source  string
	code    int
}

// Implementing error interface for type conversion.
func (e CustomError) Error() string {
	return fmt.Sprintf("%s err: %s", e.source, e.message)
}

func (e CustomError) Message() string {
	return e.message
}

func (e CustomError) Code() int {
	return e.code
}

func NewBadRequest(msg, from string) CustomError {
	return CustomError{
		message: msg,
		source:  from,
		code:    http.StatusBadRequest,
	}
}

func NewNotFound(msg, from string) CustomError {
	return CustomError{
		message: msg,
		source:  from,
		code:    http.StatusNotFound,
	}
}
