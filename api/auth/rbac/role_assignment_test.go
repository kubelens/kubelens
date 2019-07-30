package rbac

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMustFromContext(t *testing.T) {
	req := httptest.NewRequest("GET", "/fake", nil)

	ctx := NewContext(req.Context(), RoleAssignment{})
	req = req.WithContext(ctx)

	p := func() {
		a := MustFromContext(req.Context())
		assert.NotNil(t, a)
	}

	assert.NotPanics(t, p)
}

func TestMustFromContextPanics(t *testing.T) {
	req := httptest.NewRequest("GET", "/fake", nil)
	p := func() {
		a := MustFromContext(req.Context())
		assert.NotNil(t, a)
	}

	assert.Panics(t, p)
}
