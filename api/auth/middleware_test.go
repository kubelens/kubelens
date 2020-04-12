package auth

import (
	"bytes"
	"crypto/rsa"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/kubelens/kubelens/api/config"
	klog "github.com/kubelens/kubelens/api/log"
	logfakes "github.com/kubelens/kubelens/api/log/fakes"
	"github.com/kubelens/kubelens/api/testdata"
	"github.com/stretchr/testify/assert"
)

const jwkdata = `{"keys":[{"alg":"RS256","kty":"RSA","use":"sig","x5c":["MIIC9jCCAd6gAwIBAgIJfrWuWelwcfdiMA0GCSqGSIb3DQEBBQUAMCIxIDAeBgNVBAMTF3Rlc3Qtc3NvLmNocm9iaW5zb24uY29tMB4XDTE2MDUwNjE1MzM1OFoXDTMwMDExMzE2MzM1OFowIjEgMB4GA1UEAxMXdGVzdC1zc28uY2hyb2JpbnNvbi5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCtZo+k6Rz9D4yk01OVXhx05/8/UxlCyXxt182biWhoTMYNJCKXvHKEzKkvjh11/MZV/NShpOyombN1Tzg2ckABOvT23ZrhzkhN6eD2qSE0LUj9ZkNKfvouAIuMkn3NOWdjhkudPIRXfULEG0rBvCHkSIVeiLSNQ2uHW7oTpbG5EfnbZAzoOHwFIAMmd9IDv0yZOmZzcbCsRHJWQEbw2uWAKG8JY6Re5F3ftHW+3k7ROzxM5JFtk9g9jzeGwOeTUq3P6NbmDeViNxUfkxbVdQDPhNjLTjXnJ6gZGC23YgjPcX6kYMt6L+sWJZ8fBHoJw7HFkSSZR1VOCr7DFVk+AnnlAgMBAAGjLzAtMAwGA1UdEwQFMAMBAf8wHQYDVR0OBBYEFEaZUSzk2/EPQsdIRPuUQFu5LYOrMA0GCSqGSIb3DQEBBQUAA4IBAQBNqEzHRc8mZeCpTKQIvuQBXDHAL/2NOxuY/WF14GLWCj88+b9G/n7qePj4yUH9k7rRgPeg7rYUts/b5Vyo4fiHZb0LwwJPOaaNKykKizoaQnBmVyZ/gbaE9M13k172Kq3KLJwY/EHbOHi4rJi4ll3ncJpFyPMTEsncGUazwmvQuux+lCPDHx9mPFpRIZD2Cwm4NaXxwjNY0aqBqqXn6G6nCqPw39+kwj6X/hFkvcjSRMcV/2LBiJv+3e3KFnhq1+dhFotalTbdylFC7iudHmJzoWkw+g5JTZaB4FUDtdoH3NNqMBmHEjlQS4DdSrYAAb9En2uEmPVZSAwY/FtAgHXm"],"n":"rWaPpOkc_Q-MpNNTlV4cdOf_P1MZQsl8bdfNm4loaEzGDSQil7xyhMypL44ddfzGVfzUoaTsqJmzdU84NnJAATr09t2a4c5ITeng9qkhNC1I_WZDSn76LgCLjJJ9zTlnY4ZLnTyEV31CxBtKwbwh5EiFXoi0jUNrh1u6E6WxuRH522QM6Dh8BSADJnfSA79MmTpmc3GwrERyVkBG8NrlgChvCWOkXuRd37R1vt5O0Ts8TOSRbZPYPY83hsDnk1Ktz-jW5g3lYjcVH5MW1XUAz4TYy0415yeoGRgtt2IIz3F-pGDLei_rFiWfHwR6CcOxxZEkmUdVTgq-wxVZPgJ55Q","e":"AQAB","kid":"MDhDRUVDOEEyMkY5MEZBNjc5QTBGNzU5MDM0MTExRkQzMjBENTAyNg","x5t":"MDhDRUVDOEEyMkY5MEZBNjc5QTBGNzU5MDM0MTExRkQzMjBENTAyNg"}]}`

func a0Reset() {
	config.Set("../testdata/mock_config.json")
}

func TestSetAuthMiddlewareNoProvider(t *testing.T) {
	config.C.EnableAuth = false

	r := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	dctx := klog.NewContext(r.Context(), "/", &logfakes.Logger{})
	r = r.WithContext(dctx)

	mh2 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	p := func() {
		SetMiddleware(mh2).ServeHTTP(w, r)
	}

	assert.NotPanics(t, p)
}

func TestSetAuthMiddlewareOAuthProvider(t *testing.T) {
	a0Reset()

	mh := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	p := func() {
		SetMiddleware(mh)
	}

	assert.NotPanics(t, p)
}

func TestAuthMW(t *testing.T) {
	a0Reset()

	mh := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		r = httptest.NewRequest("GET", "/", nil)

		r.Header.Add("Authorization", "Bearer THIS_IS_A_KEY")
		dctx := klog.NewContext(r.Context(), "/", &logfakes.Logger{})
		r = r.WithContext(dctx)

		mh2 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			s := "good"
			w.Write([]byte(s))
		})

		a := authMiddleware(mh2)

		a.ServeHTTP(w, r)
	})

	s := httptest.NewServer(mh)

	defer s.Close()

	res, err := http.Get(s.URL)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, 401, res.StatusCode)
}

func TestAuthMWUnauthorized(t *testing.T) {
	a0Reset()

	mh := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		r = httptest.NewRequest("GET", "/", nil)
		r.Header.Add("Authorization", "Bearer THIS_IS_A_KEY")

		dctx := klog.NewContext(r.Context(), "/", &logfakes.Logger{})
		r = r.WithContext(dctx)

		r.Header.Add("test", "bad")

		mh2 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Fail(t, "Expected 401 Unautorized")
		})

		a := authMiddleware(mh2)

		a.ServeHTTP(w, r)
	})

	s := httptest.NewServer(mh)

	defer s.Close()

	res, err := http.Get(s.URL)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, 401, res.StatusCode)
}

func TestAuthMWHealthRouteAllowed(t *testing.T) {
	a0Reset()

	r := httptest.NewRequest("GET", "/health", nil)
	r.Header.Add("Authorization", "Bearer THIS_IS_A_KEY")
	w := httptest.NewRecorder()

	dctx := klog.NewContext(r.Context(), "/", &logfakes.Logger{})
	r = r.WithContext(dctx)

	r.Header.Add("test", "bad")

	mh2 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/health", r.URL.Path)
	})

	authMiddleware(mh2).ServeHTTP(w, r)
}

func TestAuthMWSocket(t *testing.T) {
	a0Reset()

	mh := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		r = httptest.NewRequest("GET", "/io/blah/blah?blah=blah&key=THIS_IS_A_KEY", nil)

		dctx := klog.NewContext(r.Context(), "/", &logfakes.Logger{})
		r = r.WithContext(dctx)

		mh2 := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.URL.Query().Get("key")
			if len(key) == 0 {
				w.WriteHeader(http.StatusNotImplemented)
				w.Write([]byte(`{"message":"missing authorization key, e.g. &key={key}"`))
				assert.Fail(t, "missing key")
				return
			}
			assert.Equal(t, "THIS_IS_A_KEY", r.URL.Query().Get("key"))
		})

		a := authMiddleware(mh2)

		a.ServeHTTP(w, r)
	})

	s := httptest.NewServer(mh)

	defer s.Close()

	res, err := http.Get(s.URL)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, 401, res.StatusCode)
}

func TestKeyLookupUnexpecedSigningMethod(t *testing.T) {

	HTTPClient = testdata.NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Status:     http.StatusText(200),
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(`{"keys":[{"alg":"RS256","kty":"RSA","use":"sig","x5c":["MIIC9jCCAd6gAwIBAgIJfrWuWelwcfdiMA0GCSqGSIb3DQEBBQUAMCIxIDAeBgNVBAMTF3Rlc3Qtc3NvLmNocm9iaW5zb24uY29tMB4XDTE2MDUwNjE1MzM1OFoXDTMwMDExMzE2MzM1OFowIjEgMB4GA1UEAxMXdGVzdC1zc28uY2hyb2JpbnNvbi5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCtZo+k6Rz9D4yk01OVXhx05/8/UxlCyXxt182biWhoTMYNJCKXvHKEzKkvjh11/MZV/NShpOyombN1Tzg2ckABOvT23ZrhzkhN6eD2qSE0LUj9ZkNKfvouAIuMkn3NOWdjhkudPIRXfULEG0rBvCHkSIVeiLSNQ2uHW7oTpbG5EfnbZAzoOHwFIAMmd9IDv0yZOmZzcbCsRHJWQEbw2uWAKG8JY6Re5F3ftHW+3k7ROzxM5JFtk9g9jzeGwOeTUq3P6NbmDeViNxUfkxbVdQDPhNjLTjXnJ6gZGC23YgjPcX6kYMt6L+sWJZ8fBHoJw7HFkSSZR1VOCr7DFVk+AnnlAgMBAAGjLzAtMAwGA1UdEwQFMAMBAf8wHQYDVR0OBBYEFEaZUSzk2/EPQsdIRPuUQFu5LYOrMA0GCSqGSIb3DQEBBQUAA4IBAQBNqEzHRc8mZeCpTKQIvuQBXDHAL/2NOxuY/WF14GLWCj88+b9G/n7qePj4yUH9k7rRgPeg7rYUts/b5Vyo4fiHZb0LwwJPOaaNKykKizoaQnBmVyZ/gbaE9M13k172Kq3KLJwY/EHbOHi4rJi4ll3ncJpFyPMTEsncGUazwmvQuux+lCPDHx9mPFpRIZD2Cwm4NaXxwjNY0aqBqqXn6G6nCqPw39+kwj6X/hFkvcjSRMcV/2LBiJv+3e3KFnhq1+dhFotalTbdylFC7iudHmJzoWkw+g5JTZaB4FUDtdoH3NNqMBmHEjlQS4DdSrYAAb9En2uEmPVZSAwY/FtAgHXm"],"n":"rWaPpOkc_Q-MpNNTlV4cdOf_P1MZQsl8bdfNm4loaEzGDSQil7xyhMypL44ddfzGVfzUoaTsqJmzdU84NnJAATr09t2a4c5ITeng9qkhNC1I_WZDSn76LgCLjJJ9zTlnY4ZLnTyEV31CxBtKwbwh5EiFXoi0jUNrh1u6E6WxuRH522QM6Dh8BSADJnfSA79MmTpmc3GwrERyVkBG8NrlgChvCWOkXuRd37R1vt5O0Ts8TOSRbZPYPY83hsDnk1Ktz-jW5g3lYjcVH5MW1XUAz4TYy0415yeoGRgtt2IIz3F-pGDLei_rFiWfHwR6CcOxxZEkmUdVTgq-wxVZPgJ55Q","e":"AQAB","kid":"MDhDRUVDOEEyMkY5MEZBNjc5QTBGNzU5MDM0MTExRkQzMjBENTAyNg","x5t":"MDhDRUVDOEEyMkY5MEZBNjc5QTBGNzU5MDM0MTExRkQzMjBENTAyNg"}]}`)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	token := &jwt.Token{}

	token.Header = make(map[string]interface{})
	token.Header["alg"] = "HS256"
	token.Header["kid"] = "MDhDRUVDOEEyMkY5MEZBNjc5QTBGNzU5MDM0MTExRkQzMjBENTAyNg"

	token.Method = jwt.SigningMethodHS256

	_, err := keyLookup(token)

	assert.Equal(t, "Unexpected signing method: HS256", err.Error())
}

func TestKeyLookup(t *testing.T) {

	HTTPClient = testdata.NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Status:     http.StatusText(200),
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(`{"keys":[{"alg":"RS256","kty":"RSA","use":"sig","x5c":["MIIC9jCCAd6gAwIBAgIJfrWuWelwcfdiMA0GCSqGSIb3DQEBBQUAMCIxIDAeBgNVBAMTF3Rlc3Qtc3NvLmNocm9iaW5zb24uY29tMB4XDTE2MDUwNjE1MzM1OFoXDTMwMDExMzE2MzM1OFowIjEgMB4GA1UEAxMXdGVzdC1zc28uY2hyb2JpbnNvbi5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCtZo+k6Rz9D4yk01OVXhx05/8/UxlCyXxt182biWhoTMYNJCKXvHKEzKkvjh11/MZV/NShpOyombN1Tzg2ckABOvT23ZrhzkhN6eD2qSE0LUj9ZkNKfvouAIuMkn3NOWdjhkudPIRXfULEG0rBvCHkSIVeiLSNQ2uHW7oTpbG5EfnbZAzoOHwFIAMmd9IDv0yZOmZzcbCsRHJWQEbw2uWAKG8JY6Re5F3ftHW+3k7ROzxM5JFtk9g9jzeGwOeTUq3P6NbmDeViNxUfkxbVdQDPhNjLTjXnJ6gZGC23YgjPcX6kYMt6L+sWJZ8fBHoJw7HFkSSZR1VOCr7DFVk+AnnlAgMBAAGjLzAtMAwGA1UdEwQFMAMBAf8wHQYDVR0OBBYEFEaZUSzk2/EPQsdIRPuUQFu5LYOrMA0GCSqGSIb3DQEBBQUAA4IBAQBNqEzHRc8mZeCpTKQIvuQBXDHAL/2NOxuY/WF14GLWCj88+b9G/n7qePj4yUH9k7rRgPeg7rYUts/b5Vyo4fiHZb0LwwJPOaaNKykKizoaQnBmVyZ/gbaE9M13k172Kq3KLJwY/EHbOHi4rJi4ll3ncJpFyPMTEsncGUazwmvQuux+lCPDHx9mPFpRIZD2Cwm4NaXxwjNY0aqBqqXn6G6nCqPw39+kwj6X/hFkvcjSRMcV/2LBiJv+3e3KFnhq1+dhFotalTbdylFC7iudHmJzoWkw+g5JTZaB4FUDtdoH3NNqMBmHEjlQS4DdSrYAAb9En2uEmPVZSAwY/FtAgHXm"],"n":"rWaPpOkc_Q-MpNNTlV4cdOf_P1MZQsl8bdfNm4loaEzGDSQil7xyhMypL44ddfzGVfzUoaTsqJmzdU84NnJAATr09t2a4c5ITeng9qkhNC1I_WZDSn76LgCLjJJ9zTlnY4ZLnTyEV31CxBtKwbwh5EiFXoi0jUNrh1u6E6WxuRH522QM6Dh8BSADJnfSA79MmTpmc3GwrERyVkBG8NrlgChvCWOkXuRd37R1vt5O0Ts8TOSRbZPYPY83hsDnk1Ktz-jW5g3lYjcVH5MW1XUAz4TYy0415yeoGRgtt2IIz3F-pGDLei_rFiWfHwR6CcOxxZEkmUdVTgq-wxVZPgJ55Q","e":"AQAB","kid":"MDhDRUVDOEEyMkY5MEZBNjc5QTBGNzU5MDM0MTExRkQzMjBENTAyNg","x5t":"MDhDRUVDOEEyMkY5MEZBNjc5QTBGNzU5MDM0MTExRkQzMjBENTAyNg"}]}`)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	token := &jwt.Token{}

	token.Header = make(map[string]interface{})
	token.Header["alg"] = "RS256"
	token.Header["kid"] = "MDhDRUVDOEEyMkY5MEZBNjc5QTBGNzU5MDM0MTExRkQzMjBENTAyNg"

	token.Method = jwt.SigningMethodRS256

	parsed, err := keyLookup(token)

	assert.Nil(t, err)
	assert.NotNil(t, parsed.(*rsa.PublicKey))
}

func TestGetPemCert(t *testing.T) {

	HTTPClient = testdata.NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Status:     http.StatusText(200),
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(`{"keys":[{"alg":"RS256","kty":"RSA","use":"sig","x5c":["MIIC9jCCAd6gAwIBAgIJfrWuWelwcfdiMA0GCSqGSIb3DQEBBQUAMCIxIDAeBgNVBAMTF3Rlc3Qtc3NvLmNocm9iaW5zb24uY29tMB4XDTE2MDUwNjE1MzM1OFoXDTMwMDExMzE2MzM1OFowIjEgMB4GA1UEAxMXdGVzdC1zc28uY2hyb2JpbnNvbi5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCtZo+k6Rz9D4yk01OVXhx05/8/UxlCyXxt182biWhoTMYNJCKXvHKEzKkvjh11/MZV/NShpOyombN1Tzg2ckABOvT23ZrhzkhN6eD2qSE0LUj9ZkNKfvouAIuMkn3NOWdjhkudPIRXfULEG0rBvCHkSIVeiLSNQ2uHW7oTpbG5EfnbZAzoOHwFIAMmd9IDv0yZOmZzcbCsRHJWQEbw2uWAKG8JY6Re5F3ftHW+3k7ROzxM5JFtk9g9jzeGwOeTUq3P6NbmDeViNxUfkxbVdQDPhNjLTjXnJ6gZGC23YgjPcX6kYMt6L+sWJZ8fBHoJw7HFkSSZR1VOCr7DFVk+AnnlAgMBAAGjLzAtMAwGA1UdEwQFMAMBAf8wHQYDVR0OBBYEFEaZUSzk2/EPQsdIRPuUQFu5LYOrMA0GCSqGSIb3DQEBBQUAA4IBAQBNqEzHRc8mZeCpTKQIvuQBXDHAL/2NOxuY/WF14GLWCj88+b9G/n7qePj4yUH9k7rRgPeg7rYUts/b5Vyo4fiHZb0LwwJPOaaNKykKizoaQnBmVyZ/gbaE9M13k172Kq3KLJwY/EHbOHi4rJi4ll3ncJpFyPMTEsncGUazwmvQuux+lCPDHx9mPFpRIZD2Cwm4NaXxwjNY0aqBqqXn6G6nCqPw39+kwj6X/hFkvcjSRMcV/2LBiJv+3e3KFnhq1+dhFotalTbdylFC7iudHmJzoWkw+g5JTZaB4FUDtdoH3NNqMBmHEjlQS4DdSrYAAb9En2uEmPVZSAwY/FtAgHXm"],"n":"rWaPpOkc_Q-MpNNTlV4cdOf_P1MZQsl8bdfNm4loaEzGDSQil7xyhMypL44ddfzGVfzUoaTsqJmzdU84NnJAATr09t2a4c5ITeng9qkhNC1I_WZDSn76LgCLjJJ9zTlnY4ZLnTyEV31CxBtKwbwh5EiFXoi0jUNrh1u6E6WxuRH522QM6Dh8BSADJnfSA79MmTpmc3GwrERyVkBG8NrlgChvCWOkXuRd37R1vt5O0Ts8TOSRbZPYPY83hsDnk1Ktz-jW5g3lYjcVH5MW1XUAz4TYy0415yeoGRgtt2IIz3F-pGDLei_rFiWfHwR6CcOxxZEkmUdVTgq-wxVZPgJ55Q","e":"AQAB","kid":"MDhDRUVDOEEyMkY5MEZBNjc5QTBGNzU5MDM0MTExRkQzMjBENTAyNg","x5t":"MDhDRUVDOEEyMkY5MEZBNjc5QTBGNzU5MDM0MTExRkQzMjBENTAyNg"}]}`)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	b, err := getPemCert("MDhDRUVDOEEyMkY5MEZBNjc5QTBGNzU5MDM0MTExRkQzMjBENTAyNg")

	assert.Nil(t, err)
	assert.Equal(t, "-----BEGIN CERTIFICATE-----\nMIIC9jCCAd6gAwIBAgIJfrWuWelwcfdiMA0GCSqGSIb3DQEBBQUAMCIxIDAeBgNVBAMTF3Rlc3Qtc3NvLmNocm9iaW5zb24uY29tMB4XDTE2MDUwNjE1MzM1OFoXDTMwMDExMzE2MzM1OFowIjEgMB4GA1UEAxMXdGVzdC1zc28uY2hyb2JpbnNvbi5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCtZo+k6Rz9D4yk01OVXhx05/8/UxlCyXxt182biWhoTMYNJCKXvHKEzKkvjh11/MZV/NShpOyombN1Tzg2ckABOvT23ZrhzkhN6eD2qSE0LUj9ZkNKfvouAIuMkn3NOWdjhkudPIRXfULEG0rBvCHkSIVeiLSNQ2uHW7oTpbG5EfnbZAzoOHwFIAMmd9IDv0yZOmZzcbCsRHJWQEbw2uWAKG8JY6Re5F3ftHW+3k7ROzxM5JFtk9g9jzeGwOeTUq3P6NbmDeViNxUfkxbVdQDPhNjLTjXnJ6gZGC23YgjPcX6kYMt6L+sWJZ8fBHoJw7HFkSSZR1VOCr7DFVk+AnnlAgMBAAGjLzAtMAwGA1UdEwQFMAMBAf8wHQYDVR0OBBYEFEaZUSzk2/EPQsdIRPuUQFu5LYOrMA0GCSqGSIb3DQEBBQUAA4IBAQBNqEzHRc8mZeCpTKQIvuQBXDHAL/2NOxuY/WF14GLWCj88+b9G/n7qePj4yUH9k7rRgPeg7rYUts/b5Vyo4fiHZb0LwwJPOaaNKykKizoaQnBmVyZ/gbaE9M13k172Kq3KLJwY/EHbOHi4rJi4ll3ncJpFyPMTEsncGUazwmvQuux+lCPDHx9mPFpRIZD2Cwm4NaXxwjNY0aqBqqXn6G6nCqPw39+kwj6X/hFkvcjSRMcV/2LBiJv+3e3KFnhq1+dhFotalTbdylFC7iudHmJzoWkw+g5JTZaB4FUDtdoH3NNqMBmHEjlQS4DdSrYAAb9En2uEmPVZSAwY/FtAgHXm\n-----END CERTIFICATE-----", string(b))
}

func TestGetPemCertNoCert(t *testing.T) {

	HTTPClient = testdata.NewTestClient(func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: 200,
			Status:     http.StatusText(200),
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(`{"keys":[{"alg":"RS256","kty":"RSA","use":"sig","x5c":["MIIC9jCCAd6gAwIBAgIJfrWuWelwcfdiMA0GCSqGSIb3DQEBBQUAMCIxIDAeBgNVBAMTF3Rlc3Qtc3NvLmNocm9iaW5zb24uY29tMB4XDTE2MDUwNjE1MzM1OFoXDTMwMDExMzE2MzM1OFowIjEgMB4GA1UEAxMXdGVzdC1zc28uY2hyb2JpbnNvbi5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCtZo+k6Rz9D4yk01OVXhx05/8/UxlCyXxt182biWhoTMYNJCKXvHKEzKkvjh11/MZV/NShpOyombN1Tzg2ckABOvT23ZrhzkhN6eD2qSE0LUj9ZkNKfvouAIuMkn3NOWdjhkudPIRXfULEG0rBvCHkSIVeiLSNQ2uHW7oTpbG5EfnbZAzoOHwFIAMmd9IDv0yZOmZzcbCsRHJWQEbw2uWAKG8JY6Re5F3ftHW+3k7ROzxM5JFtk9g9jzeGwOeTUq3P6NbmDeViNxUfkxbVdQDPhNjLTjXnJ6gZGC23YgjPcX6kYMt6L+sWJZ8fBHoJw7HFkSSZR1VOCr7DFVk+AnnlAgMBAAGjLzAtMAwGA1UdEwQFMAMBAf8wHQYDVR0OBBYEFEaZUSzk2/EPQsdIRPuUQFu5LYOrMA0GCSqGSIb3DQEBBQUAA4IBAQBNqEzHRc8mZeCpTKQIvuQBXDHAL/2NOxuY/WF14GLWCj88+b9G/n7qePj4yUH9k7rRgPeg7rYUts/b5Vyo4fiHZb0LwwJPOaaNKykKizoaQnBmVyZ/gbaE9M13k172Kq3KLJwY/EHbOHi4rJi4ll3ncJpFyPMTEsncGUazwmvQuux+lCPDHx9mPFpRIZD2Cwm4NaXxwjNY0aqBqqXn6G6nCqPw39+kwj6X/hFkvcjSRMcV/2LBiJv+3e3KFnhq1+dhFotalTbdylFC7iudHmJzoWkw+g5JTZaB4FUDtdoH3NNqMBmHEjlQS4DdSrYAAb9En2uEmPVZSAwY/FtAgHXm"],"n":"rWaPpOkc_Q-MpNNTlV4cdOf_P1MZQsl8bdfNm4loaEzGDSQil7xyhMypL44ddfzGVfzUoaTsqJmzdU84NnJAATr09t2a4c5ITeng9qkhNC1I_WZDSn76LgCLjJJ9zTlnY4ZLnTyEV31CxBtKwbwh5EiFXoi0jUNrh1u6E6WxuRH522QM6Dh8BSADJnfSA79MmTpmc3GwrERyVkBG8NrlgChvCWOkXuRd37R1vt5O0Ts8TOSRbZPYPY83hsDnk1Ktz-jW5g3lYjcVH5MW1XUAz4TYy0415yeoGRgtt2IIz3F-pGDLei_rFiWfHwR6CcOxxZEkmUdVTgq-wxVZPgJ55Q","e":"AQAB","kid":"fake","x5t":"MDhDRUVDOEEyMkY5MEZBNjc5QTBGNzU5MDM0MTExRkQzMjBENTAyNg"}]}`)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})

	_, err := getPemCert("MDhDRUVDOEEyMkY5MEZBNjc5QTBGNzU5MDM0MTExRkQzMjBENTAyNg")

	assert.Contains(t, err.Error(), "unable to find appropriate key")

}

func TestOktaAuthorization(t *testing.T) {
	_, err := oktaAuthorization(&logfakes.Logger{}, "FAKE")

	assert.NotNil(t, err)
}

func TestGetOktaRolesViewers(t *testing.T) {
	a0Reset()

	r := getOktaRoles("random@domain.com")

	assert.True(t, r.Viewers)
	assert.False(t, r.Operators)

	assert.EqualValues(t, r.MatchLabels, config.C.ViewerLabelInclusions)
	assert.EqualValues(t, r.Exclusions, config.C.ViewerLabelExclusions)
}

func TestGetOktaRolesOperators(t *testing.T) {
	a0Reset()

	r := getOktaRoles("test-admin@domain.com")

	assert.True(t, r.Viewers)
	assert.True(t, r.Operators)

	assert.Empty(t, r.Exclusions)
	assert.Empty(t, r.MatchLabels)
}
