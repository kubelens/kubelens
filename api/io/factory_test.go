package io

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	rbacfakes "github.com/kubelens/kubelens/api/auth/fakes"
	"github.com/kubelens/kubelens/api/auth/rbac"
	"github.com/kubelens/kubelens/api/config"
	k8fakes "github.com/kubelens/kubelens/api/k8v1/fakes"
	klog "github.com/kubelens/kubelens/api/log"
	logfakes "github.com/kubelens/kubelens/api/log/fakes"
)

func TestFactoryWriteReadSuccess(t *testing.T) {
	config.Set("../config/config.json")

	wsFactory := New()

	go wsFactory.Run()

	setup()

	// Create test server with the echo handler.
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dctx := klog.NewContext(r.Context(), "", &logfakes.Logger{})
		r = r.WithContext(dctx)

		ctx := rbac.NewContext(r.Context(), rbacfakes.RoleAssignment{})
		r = r.WithContext(ctx)

		wsFactory.Register(&k8fakes.K8V1{}, w, r)
	}))

	defer s.Close()

	s.URL = fmt.Sprintf("%s/io/test-pod/logs?namespace=default", s.URL)
	// Convert http://127.0.0.1 to ws://127.0.0.1
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	expected := "message"
	// Send message to server, read response and check to see if it's what we expect.
	for i := 0; i < 10; i++ {
		if err := ws.WriteMessage(logStream, []byte(expected)); err != nil {
			t.Fatalf("%v", err)
		}
		_, p, err := ws.ReadMessage()
		if err != nil {
			t.Fatalf("%v", err)
		}

		if !strings.Contains(string(p), expected) {
			t.Fatal(fmt.Sprintf("%s != %s", string(p), expected))
		}
	}
}

func TestFactoryWriteReadForbiddenHost(t *testing.T) {
	config.Set("../config/config.json")
	config.C.AllowedHosts = []string{"test"}
	wsFactory := New()

	go wsFactory.Run()

	upgrader = websocket.Upgrader{
		ReadBufferSize:  maxMessageSize * 2,
		WriteBufferSize: maxMessageSize * 2,
		CheckOrigin: func(r *http.Request) bool {
			allowed := false
			for _, host := range config.C.AllowedHosts {
				if strings.Compare(r.Host, host) == 0 {
					allowed = true
					break
				}
			}
			return allowed
		},
	}

	// Create test server with the echo handler.
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dctx := klog.NewContext(r.Context(), "", &logfakes.Logger{})
		r = r.WithContext(dctx)

		ctx := rbac.NewContext(r.Context(), rbacfakes.RoleAssignment{})
		r = r.WithContext(ctx)

		wsFactory.Register(&k8fakes.K8V1{}, w, r)
	}))

	defer s.Close()

	s.URL = fmt.Sprintf("%s/io/test-pod/logs?namespace=default", s.URL)
	// Convert http://127.0.0.1 to ws://127.0.0.1
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server
	_, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err == nil {
		t.Fatal("expected error")
	}
	// defer ws.Close()

	assert.Equal(t, "websocket: bad handshake", err.Error())
}
