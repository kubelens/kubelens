package log

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestMustFromContextPanics(t *testing.T) {
	r := &http.Request{}
	p := func() {
		MustFromContext(r.Context())
	}

	assert.Panics(t, p)
}

func TestNewLogContext(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var l Logger
		p := func() {
			l = MustFromContext(r.Context())
		}

		assert.NotPanics(t, p)

		l.Info("TestNewLogContext Succes")
	})

	ml := &logrus.Logger{}
	lm := NewMiddleware(ml, "testapp")

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	dctx := NewContext(req.Context(), "/", ml)
	req = req.WithContext(dctx)

	lm.Set(h).ServeHTTP(w, req)
}
