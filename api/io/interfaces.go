package io

import (
	"net/http"

	"github.com/kubelens/kubelens/api/k8v1"
)

// SocketFactory interfaces Factory to run websockets
type SocketFactory interface {
	// Run starts the websocket
	Run()
	// Register registers the new connection with the factory
	Register(k8Client k8v1.Clienter, w http.ResponseWriter, r *http.Request)
}
