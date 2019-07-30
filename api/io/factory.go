package io

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/kubelens/kubelens/api/auth/rbac"
	"github.com/kubelens/kubelens/api/k8v1"
	klog "github.com/kubelens/kubelens/api/log"
)

const (
	// writeDeadline - Time allowed to write a message to the peer.
	writeDeadline time.Duration = 10 * time.Second
	logStream     int           = 1
)

// Factory maintains the set of active clients
type Factory struct {
	// Clients is a map of registered clients.
	clients map[*client]bool
	// Inbound messages from the clients.
	broadcast chan []byte
	// Register requests from the clients.
	register chan *client
	// Unregister requests from clients.
	unregister chan *client
}

// New initializes the socket factory
func New() *Factory {
	return &Factory{
		broadcast:  make(chan []byte),
		register:   make(chan *client),
		unregister: make(chan *client),
		clients:    make(map[*client]bool),
	}
}

// Run starts the socket factory to handle client connections
// This function is available via NewFactory
// since Run() should only be called when the api starts.
func (f *Factory) Run() {
	for {
		select {
		case client := <-f.register:
			f.clients[client] = true
		case client := <-f.unregister:
			if _, ok := f.clients[client]; ok {
				close(client.send)
				delete(f.clients, client)
			}
		case message := <-f.broadcast:
			for client := range f.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(f.clients, client)
				}
			}
		}
	}
}

// Register handles websocket requests from the peer.
func (f *Factory) Register(k8Client k8v1.Clienter, w http.ResponseWriter, r *http.Request) {
	l := klog.MustFromContext(r.Context())
	ra := rbac.MustFromContext(r.Context())

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		l.Errorf("WebSocket Upgrader Error : %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
		return
	}

	// "/io/{pod}/logs?namespace=ns" = []string{"", "io", "pod", "logs"}
	p := strings.Split(r.URL.Path, "/")
	ns := r.URL.Query().Get("namespace")

	fmt.Printf("\nLOG STREAM:\nURL: %s\n\n", r.URL)

	if len(ns) == 0 {
		l.Error(`WebSocket Validation Error : Query string param "namespace" must be provided.`)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(http.StatusText(http.StatusBadRequest)))
		return
	}
	pod := p[2]

	c := &client{
		factory: f,
		conn:    conn,
		send:    make(chan []byte, 256),
		pod:     pod,
	}

	c.factory.register <- c

	// start the pumps
	go c.readPump()
	go c.writePump()

	// get stream
	stream, apiErr := k8Client.ReadLogs(k8v1.LogOptions{
		UserRole:  ra,
		PodName:   pod,
		Namespace: ns,
		Tail:      1,
		Follow:    true,
	})

	if apiErr != nil {
		l.Errorf("WebSocket LogStream Error : %d - %s", apiErr.Code, apiErr.Message)
		klog.Trace()
		return
	}

	f.streamLogs(l, stream, pod)
}

// streamLogs infinitely loops for as long as the reader is open
func (f *Factory) streamLogs(l klog.Logger, stream io.ReadCloser, podName string) {
	// close stream after connection drops
	defer stream.Close()

	reader := bufio.NewReaderSize(stream, 16)

	lastLine := ""

	// infinitely read stream until the connection drops
	for {
		data, isPrefix, err := reader.ReadLine()
		if err != nil {
			l.Errorf("WebSocket Error : %s", err.Error())
			klog.Trace()
			return
		}

		lines := strings.Split(string(data), "\r")

		length := len(lines)

		if len(lastLine) > 0 {
			lines[0] = lastLine + lines[0]
			lastLine = ""
		}

		if isPrefix {
			lastLine = lines[length-1]
			lines = lines[:(length - 1)]
		}

		// for every line, go through every connected client listening to the given pod logs.
		for _, line := range lines {
			for client := range f.clients {
				if client.pod == podName {
					client.conn.SetWriteDeadline(time.Now().Add(writeDeadline))
					client.conn.WriteMessage(logStream, []byte(line))
				}
			}
		}
	}
}
