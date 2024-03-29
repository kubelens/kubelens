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

func TestGetDeployments(t *testing.T) {
	h := getSvc()
	req := httptest.NewRequest("GET", `/deployments?namespace=test&linkedName=appname"`, nil)
	w := httptest.NewRecorder()

	dctx := klog.NewContext(req.Context(), "", &logfakes.Logger{})
	req = req.WithContext(dctx)

	h.Deployments(w, req)

	resp := w.Result()

	defer resp.Body.Close()

	resBody, _ := ioutil.ReadAll(resp.Body)

	var b []k8sv1.DeploymentOverview
	err := json.Unmarshal(resBody, &b)

	if err != nil {
		assert.Fail(t, err.Error())
		return
	}
	assert.Equal(t, 200, resp.StatusCode)
	assert.True(t, len(b) > 0)
}

func TestGetDeployment(t *testing.T) {
	h := getSvc()
	req := httptest.NewRequest("GET", "/deployments/test?namespace=test", nil)
	w := httptest.NewRecorder()

	dctx := klog.NewContext(req.Context(), "", &logfakes.Logger{})
	req = req.WithContext(dctx)

	h.Deployment(w, req)

	resp := w.Result()

	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)
}
