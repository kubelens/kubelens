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

	"github.com/kubelens/kubelens/api/auth/rbac"

	"github.com/creack/httpreq"
	klog "github.com/kubelens/kubelens/api/log"
)

// Logs Retrieves logs for a pod.
func (h request) Logs(w http.ResponseWriter, r *http.Request) {
	l := klog.MustFromContext(r.Context())
	ra := rbac.MustFromContext(r.Context())

	// "/v1/logs/{pod}" = []string{"", apps", "name"}
	podname := strings.Split(r.URL.Path, "/")[2]

	// get query params
	var data Req
	if err := httpreq.NewParsingMapPre(3).
		ToString("namespace", &data.Namespace).
		ToInt("tail", &data.Tail).
		ToString("containerName", &data.ContainerName).
		Parse(r.URL.Query()); err != nil {
		l.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var tl int64 = 100

	if data.Tail > 0 {
		tl = int64(data.Tail)
	}

	logs, apiErr := h.k8Client.Logs(k8sv1.LogOptions{
		UserRole:      ra,
		Logger:        l,
		Namespace:     data.Namespace,
		PodName:       podname,
		ContainerName: data.ContainerName,
		Tail:          tl,
		Follow:        false,
	})

	if apiErr != nil {
		l.Error(apiErr)
		http.Error(w, apiErr.Message, apiErr.Code)
		return
	}

	res, err := json.Marshal(logs)

	if err != nil {
		l.Error(err)
		e := errs.SerializationError(err.Error())
		http.Error(w, e.Message, e.Code)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
