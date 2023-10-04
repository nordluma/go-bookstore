package util

import (
	"errors"
	"log"
	"net/http"
)

var (
	ErrBadRequest       = errors.New("Bad request")
	ErrInternal         = errors.New("Internal error")
	ErrInvalidAPICall   = errors.New("Invalid API call")
	ErrNotAuthenticated = errors.New("Not authenticated")
	ErrResourceNotFound = errors.New("Resource not found")

	MapErrorTypeToHTTPStatus = mapErrorTypeToHTTPStatus
	IsError                  = isError
	NewError                 = newError
)

const (
	ErrorCodeInternal           = 0
	ErrorCodeInvalidJSONBody    = 30
	ErrorCodeInvalidCredentials = 201
	ErrorCodeEntityNotFound     = 404
	ErrorCodeValidation         = 500
)

type ErrorResponse struct {
	ErrorCode int
	Cause     string
}

type serverError struct {
	code      int
	cause     string
	errorType error
}

func (e serverError) Error() string {
	return e.cause
}

// Map our error types to HTTP status codes
func mapErrorTypeToHTTPStatus(err error) int {
	switch err {
	case ErrBadRequest:
		return http.StatusBadRequest
	case ErrInternal:
		return http.StatusInternalServerError
	case ErrInvalidAPICall, ErrResourceNotFound:
		return http.StatusNotFound
	case ErrNotAuthenticated:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

// Return underlying error type
func isError(errorType error) (bool, int, string, error) {
	err, IsError := errorType.(serverError)
	if !IsError {
		return false, 0, "", errorType
	}

	return true, err.code, err.cause, err.errorType
}

// Create new error
func newError(cause string, code int, errorType, err error) error {
	if err != nil {
		log.Printf("error: %v: %v", cause, err)
	} else {
		log.Printf("error: %v", cause)
	}

	return serverError{code, cause, errorType}
}
