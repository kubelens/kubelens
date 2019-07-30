/*
MIT License

Copyright (c) 2019 The KubeLens Authors

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

package v1

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kubelens/kubelens/api/config"
	"github.com/kubelens/kubelens/api/k8v1"
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
	k8Client k8v1.Clienter
}

// New creates a new request instance
func New(k8Client k8v1.Clienter) Requestor {
	return &request{
		k8Client,
	}
}

// Health checks the health of the API. Should try
// to run commands to ensure proper permissions.
func (rq request) Health(w http.ResponseWriter, r *http.Request) {
	SanityCheck := rq.k8Client.SanityCheck()
	if !SanityCheck {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"I'm not feeling too well"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(`{"message":"I'm good, thanks for checking."}`)))
}

// Register registers all routes with the v1 sub router.
func (rq request) Register(router *mux.Router) {
	router.HandleFunc(config.C.HealthRoute, rq.Health).Methods("GET")

	rv1 := router.PathPrefix("/v1").Subrouter()
	// /apps
	rv1.HandleFunc("/apps/{name}", rq.Apps).Methods("GET")

	// /pods
	rv1.HandleFunc("/pods/{name}", rq.PodDetail).Methods("GET")

	// /services
	rv1.HandleFunc("/services", rq.Services).Methods("GET")
	rv1.HandleFunc("/services/{name}", rq.Services).Methods("GET")

	// /logs
	rv1.HandleFunc("/logs/{pod}", rq.Logs).Methods("GET")
}
