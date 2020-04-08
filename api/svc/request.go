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

package svc

import (
	"net/http"

	"github.com/gorilla/mux"
	k8sv1 "github.com/kubelens/kubelens/api/k8sv1"
)

// Requestor interfaces request http handler functions
type Requestor interface {
	Register(router *mux.Router)
	Health(w http.ResponseWriter, r *http.Request)
	Apps(w http.ResponseWriter, r *http.Request)
	PodDetail(w http.ResponseWriter, r *http.Request)
	Services(w http.ResponseWriter, r *http.Request)
	Logs(w http.ResponseWriter, r *http.Request)
}

// request registers route handlers and dependencies.
type request struct {
	k8Client k8sv1.Clienter
}

// New creates a new request instance
func New(k8Client k8sv1.Clienter) Requestor {
	return &request{
		k8Client,
	}
}

// Health checks the health of the API. Should try
// to run commands to ensure proper permissions.
func (rq request) Health(w http.ResponseWriter, r *http.Request) {
	err := rq.k8Client.SanityCheck()
	if err != nil {
		w.WriteHeader(err.Code)
		w.Write([]byte(err.Message))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// Ready is just a means to indicate the API is up and running and ready for traffic.
func (rq request) Ready(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

// Register registers all routes with the v1 sub router.
func (rq request) Register(router *mux.Router) {
	router.HandleFunc("/ready", rq.Ready).Methods("GET")
	router.HandleFunc("/health", rq.Health).Methods("GET")

	// /apps
	router.HandleFunc("/apps", rq.Apps).Methods("GET")
	router.HandleFunc("/apps/{name}", rq.AppOverview).Methods("GET")

	// /pods
	router.HandleFunc("/pods/{name}", rq.PodDetail).Methods("GET")

	// /services
	router.HandleFunc("/services", rq.Services).Methods("GET")

	// /logs
	router.HandleFunc("/logs/{pod}", rq.Logs).Methods("GET")
}
