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
	Overview(w http.ResponseWriter, r *http.Request)
	Overviews(w http.ResponseWriter, r *http.Request)
	Deployment(w http.ResponseWriter, r *http.Request)
	Deployments(w http.ResponseWriter, r *http.Request)
	DaemonSet(w http.ResponseWriter, r *http.Request)
	DaemonSets(w http.ResponseWriter, r *http.Request)
	Job(w http.ResponseWriter, r *http.Request)
	Jobs(w http.ResponseWriter, r *http.Request)
	Pod(w http.ResponseWriter, r *http.Request)
	Pods(w http.ResponseWriter, r *http.Request)
	ReplicaSet(w http.ResponseWriter, r *http.Request)
	ReplicaSets(w http.ResponseWriter, r *http.Request)
	Service(w http.ResponseWriter, r *http.Request)
	Services(w http.ResponseWriter, r *http.Request)
	Logs(w http.ResponseWriter, r *http.Request)
}

// Req .
type Req struct {
	// the name of the application
	Appname string `json:"appname,omitempty"`
	// the name of a pod
	PodName string `json:"podName,omitempty"`
	// the name of a container in the pod
	ContainerName string `json:"containerName,omitempty"`
	// the namespace to filter
	Namespace string `json:"namespace,omitempty"`
	// the label selector to search, example: ?labels="app=some-app,app.kubernetes.io/name=app"
	LinkedName string `json:"linkedName"`
	// the number of lines to grab from the output
	Tail int `json:"tail,omitempty"`
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

	// /overviews
	router.HandleFunc("/overviews", rq.Overviews).Methods("GET")
	router.HandleFunc("/overviews/{linkedName}", rq.Overview).Methods("GET")

	// /daemonsets
	router.HandleFunc("/daemonsets", rq.DaemonSets).Methods("GET")
	router.HandleFunc("/daemonsets/{name}", rq.DaemonSet).Methods("GET")

	// /deployments
	router.HandleFunc("/deployments", rq.Deployments).Methods("GET")
	router.HandleFunc("/deployments/{name}", rq.Deployment).Methods("GET")

	// /jobs
	router.HandleFunc("/jobs", rq.Jobs).Methods("GET")
	router.HandleFunc("/jobs/{name}", rq.Job).Methods("GET")

	// /pods
	router.HandleFunc("/pods", rq.Pods).Methods("GET")
	router.HandleFunc("/pods/{name}", rq.Pod).Methods("GET")

	// /replicasets
	router.HandleFunc("/replicasets", rq.ReplicaSets).Methods("GET")
	router.HandleFunc("/replicasets/{name}", rq.ReplicaSet).Methods("GET")

	// /services
	router.HandleFunc("/services", rq.Services).Methods("GET")
	router.HandleFunc("/services/{name}", rq.Service).Methods("GET")

	// /logs
	router.HandleFunc("/logs/{pod}", rq.Logs).Methods("GET")
}
