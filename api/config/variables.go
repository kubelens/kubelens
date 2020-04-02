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

package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// C holds variables that can be used application wide. These fields
// are set from environment variables. These variables are meant for
// dynamic values per environment. See conf.Set() for defaults.
var C config

type config struct {
	ServerPort            int      `json:"serverPort"`
	AllowedOrigins        []string `json:"allowedOrigins"`
	AllowedMethods        []string `json:"allowedMethods"`
	AllowedHeaders        []string `json:"allowedHeaders"`
	AllowedHosts          []string `json:"allowedHosts"`
	OAuthJWK              string   `json:"oAuthJwk"`
	OAuthAudience         string   `json:"oAuthAudience"`
	OAuthJWTIssuer        string   `json:"oAuthJwtIssuer"`
	EnableAuth            bool     `json:"enableAuth"`
	EnableRBAC            bool     `json:"enableRBAC"`
	RBACClaim             string   `json:"rbacClaim"`
	ServiceNameRegex      string   `json:"serviceNameRegex"`
	ProjectSlugRegex      string   `json:"projectSlugRegex"`
	DeployerLink          string   `json:"deployerLink"`
	HealthRoute           string   `json:"healthRoute"`
	DefaultSearchLabels   []string `json:"defaultSearchLabels"`
	EnableTLS             bool     `json:"enableTLS"`
	TLSCert               string   `json:"tlsCert"`
	TLSKey                string   `json:"tlsKey"`
	ContentSecurityPolicy string   `json:"contentSecurityPolicy"`
	PublicKeyHPKP         string   `json:"publicKeyHPKP"`
	ViewerLabelInclusions []string `json:"viewerLabelInclusions"`
	ViewerLabelExclusions []string `json:"viewerLabelExclusions"`
	AdminEmails           []string `json:"adminEmails"`
}

// Set deserializes a config.json file into the config struct to allow access to
// configuration values.
func Set(file string) {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		panic(err)
	}

	b, _ := ioutil.ReadFile(file)

	err := json.Unmarshal(b, &C)

	if err != nil {
		panic(err)
	}
}
