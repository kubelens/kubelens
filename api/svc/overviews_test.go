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

func TestGetOverviews(t *testing.T) {
	h := getSvc()
	req := httptest.NewRequest("GET", `/overviews"`, nil)
	w := httptest.NewRecorder()

	dctx := klog.NewContext(req.Context(), "", &logfakes.Logger{})
	req = req.WithContext(dctx)

	h.Overviews(w, req)

	resp := w.Result()

	defer resp.Body.Close()

	resBody, _ := ioutil.ReadAll(resp.Body)

	var b []k8sv1.Overview
	err := json.Unmarshal(resBody, &b)

	if err != nil {
		assert.Fail(t, err.Error())
		return
	}
	assert.Equal(t, 200, resp.StatusCode)
	assert.True(t, len(b) > 0)
}

func TestGetAppOverviewDefault(t *testing.T) {
	h := getSvc()
	req := httptest.NewRequest("GET", `/overviews/appname?namespace=default"`, nil)
	w := httptest.NewRecorder()

	dctx := klog.NewContext(req.Context(), "", &logfakes.Logger{})
	req = req.WithContext(dctx)

	h.Overview(w, req)

	resp := w.Result()

	defer resp.Body.Close()

	resBody, _ := ioutil.ReadAll(resp.Body)

	var b k8sv1.Overview
	err := json.Unmarshal(resBody, &b)

	if err != nil {
		assert.Fail(t, err.Error())
		return
	}
	assert.Equal(t, 200, resp.StatusCode)
	assert.True(t, len(b.DaemonSets) > 0)
	assert.True(t, len(b.Deployments) > 0)
	assert.True(t, len(b.Jobs) > 0)
	assert.True(t, len(b.Pods) > 0)
	assert.True(t, len(b.ReplicaSets) > 0)
	assert.True(t, len(b.Services) > 0)
}
