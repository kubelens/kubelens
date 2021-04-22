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
	"encoding/json"
	"net/http"
	"strings"

	"github.com/kubelens/kubelens/api/errs"
	k8sv1 "github.com/kubelens/kubelens/api/k8sv1"

	"github.com/creack/httpreq"
	klog "github.com/kubelens/kubelens/api/log"
)

/*
Name string `json:"name"`
	// namespace to filter on
	Namespace string `json:"namespace"`
	// the labels to match kinds
	Labels map[string]string `json:"labels"`
	// logger instance
	Logger klog.Logger
*/

// Overviews .
func (h request) Overviews(w http.ResponseWriter, r *http.Request) {
	l := klog.MustFromContext(r.Context())

	overviews, apiErr := h.k8Client.Overviews(k8sv1.OverviewOptions{
		Logger:  l,
		Context: r.Context(),
	})

	if apiErr != nil {
		l.Error(apiErr)
		http.Error(w, apiErr.Message, apiErr.Code)
		return
	}

	res, err := json.Marshal(overviews)

	if err != nil {
		e := errs.SerializationError(err.Error())
		http.Error(w, e.Message, e.Code)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

// Overview retrieves an overview of a given name.
func (h request) Overview(w http.ResponseWriter, r *http.Request) {
	l := klog.MustFromContext(r.Context())

	var linkedName string

	// "/overviews/{linkedName}" = []string{"", "overviews", "name"}
	if params := strings.Split(r.URL.Path, "/"); len(params) >= 2 {
		linkedName = params[2]
	}

	// get query params
	var data Req
	if err := httpreq.NewParsingMapPre(1).
		ToString("namespace", &data.Namespace).
		Parse(r.URL.Query()); err != nil {
		l.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	overviews, apiErr := h.k8Client.Overview(k8sv1.OverviewOptions{
		Logger:     l,
		Context:    r.Context(),
		Namespace:  data.Namespace,
		LinkedName: linkedName,
	})

	if apiErr != nil {
		l.Error(apiErr)
		http.Error(w, apiErr.Message, apiErr.Code)
		return
	}

	res, err := json.Marshal(overviews)

	if err != nil {
		e := errs.SerializationError(err.Error())
		http.Error(w, e.Message, e.Code)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
