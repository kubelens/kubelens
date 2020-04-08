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

package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/handlers"
	kauth "github.com/kubelens/kubelens/api/auth"
	"github.com/kubelens/kubelens/api/config"
	"github.com/kubelens/kubelens/api/io"
	k8sv1 "github.com/kubelens/kubelens/api/k8sv1"
	klog "github.com/kubelens/kubelens/api/log"
	"github.com/sirupsen/logrus"
	"github.com/unrolled/secure"
)

// compression represents speed vs compression ratio.
// Lowest = 1 (fastest, least compression)
// Hightest = 9 (slowest, most compression)
const compression int = 1

func setMiddleware(wsFactory io.SocketFactory, k8Client k8sv1.Clienter, next http.Handler) http.Handler {
	logger := klog.NewMiddleware(logrus.New(), "kubelens-api")

	// any route handler registered after this point will have auth
	amw := kauth.SetMiddleware(websocketHandler(wsFactory, k8Client, next))

	secOpts := secure.Options{
		HostsProxyHeaders:     []string{"X-Forwarded-Host"},
		SSLRedirect:           config.C.EnableTLS,
		SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
		STSSeconds:            315360000,
		STSIncludeSubdomains:  config.C.EnableTLS,
		STSPreload:            config.C.EnableTLS,
		FrameDeny:             config.C.EnableTLS,
		ContentTypeNosniff:    config.C.EnableTLS,
		BrowserXssFilter:      config.C.EnableTLS,
		ContentSecurityPolicy: config.C.ContentSecurityPolicy,
	}

	if len(config.C.PublicKeyHPKP) > 0 {
		secOpts.PublicKey = config.C.PublicKeyHPKP
	}

	securemw := secure.New(secOpts)

	// set logger before route handlers and security middleware for logging access
	// gzip compression at level 1 will be fast, but not as compressed as up to 9.
	// since I'm not sure how to handle pagination yet, speed might need to take a hit for
	// larger datasets, but starting low for now to get some benefit.
	return handlers.CompressHandlerLevel(securemw.Handler(logger.Set(amw)), compression)
}

func websocketHandler(wsFactory io.SocketFactory, k8Client k8sv1.Clienter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/io/") {
			if r.Method != "GET" {
				http.Error(w, fmt.Sprintf("%s - Websocket connection must be GET.", http.StatusText(http.StatusForbidden)), http.StatusForbidden)
				return
			}
			// upgrade connection
			w.Header().Set("Connection", "keep-alive")
			wsFactory.Register(k8Client, w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
