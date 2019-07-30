package v1

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/kubelens/kubelens/api/config"
	"github.com/kubelens/kubelens/api/k8v1/fakes"
)

func getSvc() *request {
	return &request{&fakes.K8V1{}}
}

func TestRegister(t *testing.T) {
	rc := mux.NewRouter()

	rq := New(&fakes.K8V1{})

	p := func() {
		rq.Register(rc)
	}

	assert.NotPanics(t, p)
}

func TestHealthHandler(t *testing.T) {
	config.Set("../config/config.json")
	rc := mux.NewRouter()

	rq := New(&fakes.K8V1{})
	rq.Register(rc)

	ts := httptest.NewServer(rc)
	defer ts.Close()

	// GET /health
	res, err := http.Get((ts.URL + config.C.HealthRoute))
	if err != nil {
		assert.Fail(t, err.Error())
	}
	r1, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	var m map[string]interface{}
	err = json.Unmarshal(r1, &m)

	if err != nil {
		assert.Fail(t, err.Error())
	}
	assert.Equal(t, "I'm good, thanks for checking.", m["message"])
}
