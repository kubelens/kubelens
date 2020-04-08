/*
MIT License

Copyright (c) 2020 The KubeLens Authors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package log

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type contextKey string

var loggerKey = contextKey("applogger")
var loginst *logrus.Logger
var appname string

// Middleware holds functions for setting a logger instance to
// Context within http.Handler middleware
type Middleware struct{}

// NewMiddleware returns an instance of logger
func NewMiddleware(l *logrus.Logger, a string) *Middleware {
	loginst = l
	appname = a
	return &Middleware{}
}

// Set adds the logger to the request context
func (*Middleware) Set(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dctx := NewContext(r.Context(), r.URL.RequestURI(), &logger{
			log: loginst,
			fields: logrus.Fields{
				"timestamp":  time.Now(),
				"request_id": uuid.New().String(),
				"appname":    appname,
			},
		})

		r = r.WithContext(dctx)

		h(w, r)
	}
}

// NewContext adds logger to context
func NewContext(ctx context.Context, uri string, l Logger) context.Context {
	l.Info(uri)
	return context.WithValue(ctx, loggerKey, l)
}

// MustFromContext retrieves a logger instance from context and panics if not found
func MustFromContext(ctx context.Context) Logger {
	instance, ok := ctx.Value(loggerKey).(Logger)
	if !ok {
		panic("logger was not found in context")
	}
	return instance
}
