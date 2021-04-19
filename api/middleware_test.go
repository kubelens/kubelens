package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	iofakes "github.com/kubelens/kubelens/api/io/fakes"
	k8fakes "github.com/kubelens/kubelens/api/k8sv1/fakes"
	logfakes "github.com/kubelens/kubelens/api/log/fakes"

	"github.com/kubelens/kubelens/api/config"
	klog "github.com/kubelens/kubelens/api/log"
	"github.com/stretchr/testify/assert"
)

func TestSetMiddleware(t *testing.T) {
	config.C.EnableAuth = false
	req := httptest.NewRequest("GET", "/io/", nil)
	w := httptest.NewRecorder()

	dctx := klog.NewContext(req.Context(), "", &logfakes.Logger{})
	req = req.WithContext(dctx)

	tmw := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/io/", r.URL.Path)
	})

	mw := setMiddleware(&iofakes.SocketFactory{}, &k8fakes.K8sV1{}, tmw)

	mw.ServeHTTP(w, req)
}

func TestSetMiddlewareIncorrectSocketMethod(t *testing.T) {
	config.C.EnableAuth = false
	req := httptest.NewRequest("POST", "/io/", nil)
	w := httptest.NewRecorder()

	dctx := klog.NewContext(req.Context(), "", &logfakes.Logger{})
	req = req.WithContext(dctx)

	tmw := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/io/", r.URL.Path)
	})

	wsh := websocketHandler(&iofakes.SocketFactory{}, &k8fakes.K8sV1{}, tmw)
	wsh.ServeHTTP(w, req)
}
