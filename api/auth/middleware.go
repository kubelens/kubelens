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

package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
	jwtverifier "github.com/okta/okta-jwt-verifier-golang"

	"github.com/kubelens/kubelens/api/config"
	klog "github.com/kubelens/kubelens/api/log"
)

type jwkKey struct {
	Alg string   `json:"alg"`
	Kty string   `json:"kty"`
	Use string   `json:"use"`
	X5c []string `json:"x5c"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	Kid string   `json:"kid"`
	X5t string   `json:"x5t"`
}

type jwk struct {
	Keys []jwkKey `json:"keys"`
}

// HTTPClient can be an actual client or mocked version for testing.
// Is there a better way to handle this?
var HTTPClient *http.Client

// SetMiddleware adds the auth provider middleware given a supported provider
func SetMiddleware(next http.Handler) http.HandlerFunc {
	if config.C.EnableAuth {
		// really used for testing
		if HTTPClient == nil {
			HTTPClient = &http.Client{}
		}
		return authMiddleware(next).ServeHTTP
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

// authMiddleware should work for any OCID compliant OAuth providers.
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := klog.MustFromContext(r.Context())

		if r.URL.Path == "/health" || r.URL.Path == "/ready" {
			next.ServeHTTP(w, r)
			return
		}

		// Support either JWT by query string or header for WebSockets.
		// Vanilla JavaScript's WebSocket has support for basic auth via URI connection string
		// so this is a work around to still allow Bearer auth.
		if strings.HasPrefix(r.URL.Path, "/io/") && len(r.Header.Get("Authorization")) == 0 {
			key := r.URL.Query().Get("key")
			// convert key to auth header
			r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", key))
			// remove key's value from query string for less confusion.
			// unsure if the query string key itself can be removed but not that big of a deal.
			r.URL.Query().Del("key")
		}

		// need to skip claims if OPTIONS
		if !strings.EqualFold(r.Method, http.MethodOptions) {
			// -H "Authorization: Bearer <JWT>"
			authBearer := r.Header.Get("Authorization")
			// "Bearer " is 7 characters
			if len(authBearer) < 7 {
				l.Errorf("Missing Authoriztion Header: %s - %+v", r.URL.RequestURI(), r.Header)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
				return
			}
			// Authorization: Bearer ... so it's 7 characters to start of jwt
			requestJWT := authBearer[7:]

			// check for Okta specific provider
			if strings.Contains(strings.ToLower(config.C.OAuthJWTIssuer), "okta") {
				if _, err := oktaAuthorization(l, requestJWT); err != nil {
					l.Errorf("ERROR: %s : %s - %+v", err.Error(), r.URL.RequestURI(), r.Header)
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
					return
				}

			} else {
				// support for generic jwt verification.
				if _, err := jwt.ParseWithClaims(requestJWT, &jwt.StandardClaims{}, keyLookup); err != nil {
					l.Errorf("ERROR: %s : %s - %+v", err.Error(), r.URL.RequestURI(), r.Header)
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
					return
				}
			}

			// log user request
			go l.Infof("%s - %v", r.URL.RequestURI(), time.Now())
		}

		next.ServeHTTP(w, r)
	})
}

func oktaAuthorization(l klog.Logger, requestJWT string) (*jwtverifier.Jwt, error) {
	toValidate := map[string]string{}

	toValidate["aud"] = config.C.OAuthAudience
	toValidate["cid"] = config.C.OAuthClientID

	jwtVerifierSetup := jwtverifier.JwtVerifier{
		Issuer:           config.C.OAuthJWTIssuer,
		ClaimsToValidate: toValidate,
	}

	verifier := jwtVerifierSetup.New()

	return verifier.VerifyAccessToken(requestJWT)
}

func keyLookup(token *jwt.Token) (interface{}, error) {
	// validate the alg
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	}

	// Look up key
	key, err := getPemCert(token.Header["kid"].(string))
	if err != nil {
		return nil, err
	}

	// Unpack key from PEM encoded PKCS8
	return jwt.ParseRSAPublicKeyFromPEM(key)
}

func getPemCert(kid string) (cert []byte, err error) {
	// get well-known jwks
	resp, err := HTTPClient.Get(config.C.OAuthJWK)

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks jwk
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	// get cert value from well-known
	for k := range jwks.Keys {
		if kid == jwks.Keys[k].Kid {
			cert = []byte(fmt.Sprintf("-----BEGIN CERTIFICATE-----\n%s\n-----END CERTIFICATE-----", jwks.Keys[k].X5c[0]))
		}
	}

	if len(cert) == 0 {
		return cert, fmt.Errorf("unable to find appropriate key: kid: %s - jwks: %+v", kid, jwks)
	}

	return cert, nil
}
