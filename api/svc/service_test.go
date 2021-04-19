package svc

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	k8sv1 "github.com/kubelens/kubelens/api/k8sv1"
	klog "github.com/kubelens/kubelens/api/log"
	logfakes "github.com/kubelens/kubelens/api/log/fakes"
	"github.com/stretchr/testify/assert"
)

func TestGetServicesDefault(t *testing.T) {
	h := getSvc()
	req := httptest.NewRequest("GET", `/services?detailed=false&labelSelector="app.kubernetes.io/name=service-name,component=api"`, nil)
	w := httptest.NewRecorder()

	dctx := klog.NewContext(req.Context(), "", &logfakes.Logger{})
	req = req.WithContext(dctx)

	h.Services(w, req)

	resp := w.Result()

	defer resp.Body.Close()

	resBody, _ := ioutil.ReadAll(resp.Body)

	var b []k8sv1.ServiceOverview
	err := json.Unmarshal(resBody, &b)

	if err != nil {
		assert.Fail(t, err.Error())
		return
	}
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "-service", b[0].Name)
}

func TestGetServicesError(t *testing.T) {
	h := getSvc()
	req := httptest.NewRequest("GET", "/services?namespace=bad", nil)
	w := httptest.NewRecorder()

	dctx := klog.NewContext(req.Context(), "", &logfakes.Logger{})
	req = req.WithContext(dctx)

	h.Services(w, req)

	resp := w.Result()

	defer resp.Body.Close()

	assert.Equal(t, 500, resp.StatusCode)
}
