package svc

import (
	"encoding/json"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	rbacfakes "github.com/kubelens/kubelens/api/auth/fakes"
	"github.com/kubelens/kubelens/api/auth/rbac"
	k8sv1 "github.com/kubelens/kubelens/api/k8sv1"
	klog "github.com/kubelens/kubelens/api/log"
	logfakes "github.com/kubelens/kubelens/api/log/fakes"
	"github.com/stretchr/testify/assert"
)

func TestGetJobs(t *testing.T) {
	h := getSvc()
	req := httptest.NewRequest("GET", `/jobs?namespace=test&linkedName=appname"`, nil)
	w := httptest.NewRecorder()

	dctx := klog.NewContext(req.Context(), "", &logfakes.Logger{})
	req = req.WithContext(dctx)

	ctx := rbac.NewContext(req.Context(), rbacfakes.RoleAssignment{})
	req = req.WithContext(ctx)

	h.Jobs(w, req)

	resp := w.Result()

	defer resp.Body.Close()

	resBody, _ := ioutil.ReadAll(resp.Body)

	var b []k8sv1.JobOverview
	err := json.Unmarshal(resBody, &b)

	if err != nil {
		assert.Fail(t, err.Error())
		return
	}
	assert.Equal(t, 200, resp.StatusCode)
	assert.True(t, len(b) > 0)
}

func TestGetJob(t *testing.T) {
	h := getSvc()
	req := httptest.NewRequest("GET", "/jobs/test?namespace=test", nil)
	w := httptest.NewRecorder()

	dctx := klog.NewContext(req.Context(), "", &logfakes.Logger{})
	req = req.WithContext(dctx)

	ctx := rbac.NewContext(req.Context(), rbacfakes.RoleAssignment{})
	req = req.WithContext(ctx)

	h.Job(w, req)

	resp := w.Result()

	defer resp.Body.Close()

	assert.Equal(t, 200, resp.StatusCode)
}