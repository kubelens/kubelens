package io

import (
	"net/http"

	"github.com/gorilla/websocket"
)

func setup() {
	upgrader = websocket.Upgrader{
		ReadBufferSize:  maxMessageSize * 2,
		WriteBufferSize: maxMessageSize * 2,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
}
