package svc

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/kubelens/kubelens/api/config"
	"github.com/kubelens/kubelens/api/k8sv1/fakes"
	"github.com/stretchr/testify/assert"
)

func getSvc() *request {
	return &request{&fakes.K8sV1{}}
}

func TestRegister(t *testing.T) {
	rc := mux.NewRouter()

	rq := New(&fakes.K8sV1{})

	p := func() {
		rq.Register(rc)
	}

	assert.NotPanics(t, p)
}

func TestHealthHandler(t *testing.T) {
	config.Set("../config/config.json")
	rc := mux.NewRouter()

	rq := New(&fakes.K8sV1{})
	rq.Register(rc)

	ts := httptest.NewServer(rc)
	defer ts.Close()

	// GET /health
	res, err := http.Get((ts.URL + "/health"))
	if err != nil {
		assert.Fail(t, err.Error())
	}
	r1, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		assert.Fail(t, err.Error())
	}
	assert.Equal(t, http.StatusText(http.StatusOK), string(r1))
}
