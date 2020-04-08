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
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// C holds variables that can be used application wide. These fields
// are set from environment variables. These variables are meant for
// dynamic values per environment. See conf.Set() for defaults.
var config c

type c struct {
	ServerPort int    `json:"serverPort"`
	EnableTLS  bool   `json:"enableTLS"`
	TLSCert    string `json:"tlsCert"`
	TLSKey     string `json:"tlsKey"`
}

// ths could be cleaner, but it'll probably never be looked at again anyway
func main() {
	confPath := flag.String("confpath", "", `config files path: usage -confpath="./"`)
	sitePath := flag.String("sitepath", "", `path to website files: usage -sitepath="./"`)
	flag.Parse()

	if confPath == nil || len(*confPath) == 0 {
		panic(fmt.Errorf("missing config files path"))
	}

	if sitePath == nil || len(*sitePath) == 0 {
		panic(fmt.Errorf("missing path to website files"))
	}

	file := fmt.Sprintf("%s/server.json", *confPath)
	if _, err := os.Stat(file); os.IsNotExist(err) {
		panic(err)
	}

	b, _ := ioutil.ReadFile(file)

	err := json.Unmarshal(b, &config)

	if err != nil {
		panic(err)
	}

	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		cf := fmt.Sprintf("%s/config.json", *confPath)

		if _, err := os.Stat(cf); os.IsNotExist(err) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		b, err := ioutil.ReadFile(cf)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(b)
	})

	// html related files should be relative to the executing path of this
	http.Handle("/", http.FileServer(http.Dir(*sitePath)))

	host, _ := os.Hostname()

	address := fmt.Sprintf(":%d", config.ServerPort)

	fmt.Printf("\nhost %s listing on %s\n", host, address)

	if config.EnableTLS {
		panic(http.ListenAndServeTLS(address, config.TLSCert, config.TLSKey, nil))
	} else {
		panic(http.ListenAndServe(address, nil))
	}
}
