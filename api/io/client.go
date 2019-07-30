package io

import (
	"bytes"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kubelens/kubelens/api/config"
)

const (
	// writeWait - Time allowed to write a message to the peer.
	writeWait time.Duration = 10 * time.Second

	// pongWait - Time allowed to read the next pong message from the peer.
	pongWait time.Duration = 60 * time.Second

	// pingPeriod - Send pings to peer with this period. Must be less than pongWait.
	pingPeriod time.Duration = (pongWait * 9) / 10

	// maxMessageSize - Maximum message size allowed from peer.
	maxMessageSize int = 2048
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
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

// client is a middleman between the websocket connection and the factory.
type client struct {
	factory *Factory
	// Conn is the websocket connection.
	conn *websocket.Conn
	// Buffered channel of outbound messages.
	send chan []byte
	// pod is the pod name to key with the connection
	pod string
}

// readPump pumps messages from the websocket connection to the Factory.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *client) readPump() {
	defer func() {
		c.factory.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(int64(maxMessageSize))
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	message := make([]byte, upgrader.ReadBufferSize)

	for {
		_, mb, err := c.conn.NextReader()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		_, err = mb.Read(message)
		// _, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		c.factory.broadcast <- message
	}
}

// writePump pumps messages from the Factory to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
