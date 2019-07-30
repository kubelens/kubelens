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

package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kubelens/kubelens/api/config"
	"github.com/kubelens/kubelens/api/io"
	"github.com/kubelens/kubelens/api/k8v1"
	v1 "github.com/kubelens/kubelens/api/v1"
)

func serve() {
	wsFactory := io.New()
	hs := createServer(wsFactory)

	// run websocket
	go wsFactory.Run()

	if config.C.EnableTLS {
		panic(hs.ListenAndServeTLS(config.C.TLSCert, config.C.TLSKey))
	} else {
		panic(hs.ListenAndServe())
	}
}

// createServer creates the http server with middleware.
func createServer(wsFactory io.SocketFactory) *http.Server {
	rc := mux.NewRouter()

	// v1 handlers
	k8Client := k8v1.New(k8v1.NewWrapper())
	creq := v1.New(k8Client)
	creq.Register(rc)

	address := fmt.Sprintf(":%v", config.C.ServerPort)

	hs := &http.Server{
		Addr: address,
		// set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler: handlers.CORS(
			handlers.AllowedMethods(config.C.AllowedMethods),
			handlers.AllowedHeaders(config.C.AllowedHeaders),
			handlers.AllowedOrigins(config.C.AllowedOrigins),
		)(setMiddleware(wsFactory, k8Client, rc)),
	}

	hostname, _ := os.Hostname()

	fmt.Printf("\nhost %s - listening on port %s\n", hostname, address)

	return hs
}
