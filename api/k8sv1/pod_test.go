package k8v1

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	rbacfakes "github.com/kubelens/kubelens/api/auth/fakes"
	"github.com/kubelens/kubelens/api/config"
	logfakes "github.com/kubelens/kubelens/api/log/fakes"
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
		UserRole:        rbacfakes.RoleAssignment{},
		Logger:          &logfakes.Logger{},
		AppNameLabelKey: "app",
		Namespace:       "default",
		AppName:         "test",
	})

	assert.Nil(t, err)
	assert.True(t, len(r.PodDetails) > 0)
}

func TestPodOverviewDefaultWithFilters(t *testing.T) {
	c := setupClient("default", "test", false, false)

	r, err := c.PodOverview(PodOverviewOptions{
		UserRole:        rbacfakes.RoleAssignment{},
		Logger:          &logfakes.Logger{},
		AppNameLabelKey: "app",
		Namespace:       "default",
		AppName:         "test",
	})

	assert.Nil(t, err)
	assert.True(t, len(r.PodDetails) > 0)
}

func TestGetPodOverviewByName(t *testing.T) {
	n := "test"

	c := setupClient("blah", n, false, false)

	r, err := c.PodOverview(PodOverviewOptions{
		UserRole:        rbacfakes.RoleAssignment{},
		Logger:          &logfakes.Logger{},
		AppNameLabelKey: "app",
		Namespace:       "blah",
		AppName:         "test",
	})

	assert.Nil(t, err)
	assert.Equal(t, n, r.Name.Value)
}

func TestGetPodOverviewByNamespaceAndName(t *testing.T) {
	config.Set("../config/config.json")
	ns := "app"
	n := "some-projects-1234-dev"

	c := setupClient(ns, n, false, false)

	r, err := c.PodOverview(PodOverviewOptions{
		UserRole:        rbacfakes.RoleAssignment{},
		Logger:          &logfakes.Logger{},
		AppNameLabelKey: "app",
		Namespace:       ns,
		AppName:         n,
	})

	assert.Nil(t, err)
	assert.Equal(t, ns, r.Name.LabelKey)
	assert.Equal(t, n, r.Name.Value)
}
