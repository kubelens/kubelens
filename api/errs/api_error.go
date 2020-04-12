package errs

import (
	"fmt"
	"net/http"
	"strings"
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

// ListToInternalServerError concats a list of APIErrors and returns them as a single pipe delimited error.
func ListToInternalServerError(list []*APIError) *APIError {
	var errstr string
	for i, e := range list {
		errstr += fmt.Sprintf("Error %d: %s\n", (i + 1), strings.ReplaceAll(e.Message, "\n", ""))
	}
	return &APIError{
		Code:    http.StatusInternalServerError,
		Message: fmt.Sprintf("%s\n", errstr),
	}
}
