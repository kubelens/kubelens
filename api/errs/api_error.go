package errs

import (
	"fmt"
	"net/http"
)

// APIError represents an Kubelens API error
type APIError struct {
	// http code
	Code int `json:"code"`
	// error message
	Message string `json:"message"`
}

// Unauthorized returns 401/Unauthorized status/message
func Unauthorized() *APIError {
	return &APIError{
		Code:    http.StatusUnauthorized,
		Message: http.StatusText(http.StatusUnauthorized),
	}
}

// Forbidden returns 403/Forbidden status/message
func Forbidden() *APIError {
	return &APIError{
		Code:    http.StatusForbidden,
		Message: http.StatusText(http.StatusForbidden),
	}
}

// InternalServerError returns 500/Message status/message
func InternalServerError(err string) *APIError {
	return &APIError{
		Code:    http.StatusInternalServerError,
		Message: fmt.Sprintf("\n%s: %s\n", http.StatusText(http.StatusInternalServerError), err),
	}
}

// ValidationError returns 400/Message status/message
func ValidationError(err string) *APIError {
	return &APIError{
		Code:    http.StatusBadRequest,
		Message: fmt.Sprintf("\n%s: %s\n", http.StatusText(http.StatusBadRequest), err),
	}
}

// SerializationError returns 500/Message status/message
func SerializationError(err string) *APIError {
	return &APIError{
		Code:    http.StatusInternalServerError,
		Message: fmt.Sprintf("\n%s: %s\n", http.StatusText(http.StatusInternalServerError), err),
	}
}
