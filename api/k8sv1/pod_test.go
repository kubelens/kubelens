package k8sv1

import (
	"net/http"
	"testing"

	rbacfakes "github.com/kubelens/kubelens/api/auth/fakes"
	"github.com/kubelens/kubelens/api/config"
	logfakes "github.com/kubelens/kubelens/api/log/fakes"
	"github.com/stretchr/testify/assert"
)

func TestPodDetailDefault(t *testing.T) {
	c := setupClient("default", "test", false, false)

	r, err := c.PodDetail(PodDetailOptions{
		UserRole:  rbacfakes.RoleAssignment{},
		Logger:    &logfakes.Logger{},
		Name:      "test",
		Namespace: "default",
	})

	assert.Nil(t, err)
	assert.Equal(t, "test", r.Name)
}

func TestPodDetailForbidden(t *testing.T) {
	c := setupClient("default", "test", true, false)

	_, err := c.PodDetail(PodDetailOptions{
		UserRole:  rbacfakes.RoleAssignment{},
		Logger:    &logfakes.Logger{},
		Name:      "test",
		Namespace: "default",
	})

	assert.Equal(t, http.StatusInternalServerError, err.Code)
	assert.Equal(t, "\nInternal Server Error: GetClientSet Test Error\n", err.Message)
}

func TestPodOverviewDefault(t *testing.T) {
	c := setupClient("default", "test", false, false)

	r, err := c.PodOverview(PodOverviewOptions{
		UserRole:      rbacfakes.RoleAssignment{},
		Logger:        &logfakes.Logger{},
		LabelSelector: map[string]string{"app": "test"},
		Namespace:     "default",
		AppName:       "test",
	})

	assert.Nil(t, err)
	assert.True(t, len(r.PodInfo) > 0)
}

func TestPodOverviewDefaultWithFilters(t *testing.T) {
	c := setupClient("default", "test", false, false)

	lbl := make(map[string]string)
	lbl["app"] = "test"
	lbl[AppNameLabel] = FriendlyAppName

	r, err := c.PodOverview(PodOverviewOptions{
		UserRole:      rbacfakes.RoleAssignment{},
		Logger:        &logfakes.Logger{},
		LabelSelector: lbl,
		Namespace:     "default",
		AppName:       "test",
	})

	assert.Nil(t, err)
	assert.True(t, len(r.PodInfo) > 0)
}

func TestGetPodOverviewByName(t *testing.T) {
	n := "test"

	c := setupClient("blah", n, false, false)

	r, err := c.PodOverview(PodOverviewOptions{
		UserRole:      rbacfakes.RoleAssignment{},
		Logger:        &logfakes.Logger{},
		LabelSelector: map[string]string{"app": "test"},
		Namespace:     "blah",
		AppName:       "test",
	})

	assert.Nil(t, err)
	assert.Equal(t, n, r.Name)
}

func TestGetPodOverviewByNamespaceAndName(t *testing.T) {
	config.Set("../config/config.json")
	ns := "somens"
	n := "some-projects-1234-dev"

	c := setupClient(ns, n, false, false)

	r, err := c.PodOverview(PodOverviewOptions{
		UserRole:      rbacfakes.RoleAssignment{},
		Logger:        &logfakes.Logger{},
		LabelSelector: map[string]string{"app": n},
		Namespace:     ns,
		AppName:       n,
	})

	assert.Nil(t, err)
	assert.Equal(t, ns, r.Namespace)
	assert.Equal(t, n, r.Name)
}
