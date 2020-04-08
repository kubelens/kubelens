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

	"github.com/creack/httpreq"
	"github.com/kubelens/kubelens/api/auth/rbac"
	"github.com/kubelens/kubelens/api/errs"
	k8sv1 "github.com/kubelens/kubelens/api/k8sv1"
	klog "github.com/kubelens/kubelens/api/log"
)

// Services retrieves a list of services for an application
func (h request) Services(w http.ResponseWriter, r *http.Request) {
	l := klog.MustFromContext(r.Context())
	ra := rbac.MustFromContext(r.Context())

	// get query params
	var data Req
	if err := httpreq.NewParsingMapPre(1).
		ToString("namespace", &data.Namespace).
		ToString("labelSelector", &data.LabelSelector).
		ToBool("detailed", &data.Detailed).
		Parse(r.URL.Query()); err != nil {
		l.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	so := k8sv1.ServiceOptions{
		UserRole:      ra,
		Logger:        l,
		Namespace:     data.Namespace,
		Detailed:      data.Detailed,
		LabelSelector: data.LabelSelectorMap(),
	}

	services, apiErr := h.k8Client.ServiceOverviews(so)

	if apiErr != nil {
		l.Error(apiErr)
		http.Error(w, apiErr.Message, apiErr.Code)
		return
	}

	res, err := json.Marshal(services)

	if err != nil {
		e := errs.SerializationError(err.Error())
		http.Error(w, e.Message, e.Code)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(res)
}
