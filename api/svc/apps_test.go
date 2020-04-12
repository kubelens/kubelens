package svc

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	rbacfakes "github.com/kubelens/kubelens/api/auth/fakes"
	"github.com/kubelens/kubelens/api/auth/rbac"
	k8sv1 "github.com/kubelens/kubelens/api/k8sv1"
	klog "github.com/kubelens/kubelens/api/log"
	logfakes "github.com/kubelens/kubelens/api/log/fakes"
	"github.com/stretchr/testify/assert"
)

func TestGetAppsDefault(t *testing.T) {
	h := getSvc()
	req := httptest.NewRequest("GET", `/apps?labelSelector="app.kubernetes.io/name=test,component=api"`, nil)
	w := httptest.NewRecorder()

	dctx := klog.NewContext(req.Context(), "", &logfakes.Logger{})
	req = req.WithContext(dctx)

	ctx := rbac.NewContext(req.Context(), rbacfakes.RoleAssignment{})
	req = req.WithContext(ctx)

	h.Apps(w, req)

	resp := w.Result()

	defer resp.Body.Close()

	resBody, _ := ioutil.ReadAll(resp.Body)

	var b []App
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
	req := httptest.NewRequest("GET", `/apps/test?labelSelector="app.kubernetes.io/name=test,component=api"`, nil)
	w := httptest.NewRecorder()

	dctx := klog.NewContext(req.Context(), "", &logfakes.Logger{})
	req = req.WithContext(dctx)

	ctx := rbac.NewContext(req.Context(), rbacfakes.RoleAssignment{})
	req = req.WithContext(ctx)

	h.AppOverview(w, req)

	resp := w.Result()

	defer resp.Body.Close()

	resBody, _ := ioutil.ReadAll(resp.Body)

	var b k8sv1.AppOverview
	err := json.Unmarshal(resBody, &b)

	if err != nil {
		assert.Fail(t, err.Error())
		return
	}
	assert.Equal(t, 200, resp.StatusCode)
	assert.True(t, len(b.ServiceOverviews) > 0)
	assert.True(t, len(b.PodOverviews.Name) > 0)
}

func TestGetAppOverviewDefaultMissingLabelSelector(t *testing.T) {
	h := getSvc()
	req := httptest.NewRequest("GET", `/apps/test?labelSelector=`, nil)
	w := httptest.NewRecorder()

	dctx := klog.NewContext(req.Context(), "", &logfakes.Logger{})
	req = req.WithContext(dctx)

	ctx := rbac.NewContext(req.Context(), rbacfakes.RoleAssignment{})
	req = req.WithContext(ctx)

	h.AppOverview(w, req)

	resp := w.Result()

	defer resp.Body.Close()

	resBody, _ := ioutil.ReadAll(resp.Body)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Equal(t, "missing labelSelector\n", string(resBody))
}
