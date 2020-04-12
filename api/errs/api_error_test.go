package errs

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnauthorized(t *testing.T) {
	r := Unauthorized()
	assert.Equal(t, http.StatusUnauthorized, r.Code)
	assert.Equal(t, "Unauthorized", r.Message)
}

func TestForbidden(t *testing.T) {
	r := Forbidden()
	assert.Equal(t, http.StatusForbidden, r.Code)
	assert.Equal(t, "Forbidden", r.Message)
}

func TestInternalServerError(t *testing.T) {
	r := InternalServerError("test")
	assert.Equal(t, http.StatusInternalServerError, r.Code)
	assert.Equal(t, "\nInternal Server Error: test\n", r.Message)
}

func TestValidationError(t *testing.T) {
	r := ValidationError("test")
	assert.Equal(t, http.StatusBadRequest, r.Code)
	assert.Equal(t, "\nBad Request: test\n", r.Message)
}

func TestSerializationError(t *testing.T) {
	r := SerializationError("test")
	assert.Equal(t, http.StatusInternalServerError, r.Code)
	assert.Equal(t, "\nInternal Server Error: test\n", r.Message)
}

func TestListToInternalServerError(t *testing.T) {
	list := []*APIError{
		InternalServerError("test1"),
		InternalServerError("test2"),
	}

	r := ListToInternalServerError(list)

	assert.Equal(t, http.StatusInternalServerError, r.Code)
	assert.Equal(t, "Error 1: Internal Server Error: test1\nError 2: Internal Server Error: test2\n\n", r.Message)
}
