package k8v1

import (
	"testing"

	rbacfakes "github.com/kubelens/kubelens/api/auth/fakes"
	logfakes "github.com/kubelens/kubelens/api/log/fakes"
	"github.com/stretchr/testify/assert"
)

func TestGetServiceOverviewsDefaultSuccess(t *testing.T) {
	c := setupClient("test", "test", false, false)

	_, err := c.ServiceOverviews(ServiceOptions{
		UserRole:      &rbacfakes.RoleAssignment{},
		Logger:        &logfakes.Logger{},
		Namespace:     "test",
		LabelSelector: "app=test",
	})

	assert.Nil(t, err)
}

func TestGetServiceOverviewsDefaultMissingNamespace(t *testing.T) {
	c := setupClient("test", "test", false, false)

	_, err := c.ServiceOverviews(ServiceOptions{
		UserRole:      &rbacfakes.RoleAssignment{},
		Logger:        &logfakes.Logger{},
		Namespace:     "",
		LabelSelector: "app=test",
	})

	assert.Nil(t, err)
}

func TestGetServiceOverviewsDefaultMissingLabelSearch(t *testing.T) {
	c := setupClient("test", "test", false, false)

	_, err := c.ServiceOverviews(ServiceOptions{
		UserRole:      &rbacfakes.RoleAssignment{},
		Logger:        &logfakes.Logger{},
		Namespace:     "fake",
		LabelSelector: "",
	})

	assert.Nil(t, err)
}

func TestGetServiceOverviewsCondensed(t *testing.T) {
	c := setupClient("test", "test", false, false)

	r, err := c.ServiceOverviews(ServiceOptions{
		UserRole:  &rbacfakes.RoleAssignment{},
		Logger:    &logfakes.Logger{},
		Namespace: "test",
		Detailed:  false,
	})

	assert.Nil(t, err)
	assert.Nil(t, r[0].Spec)
}

func TestGetServiceOverviewsDetailed(t *testing.T) {
	c := setupClient("test", "test", false, false)

	r, err := c.ServiceOverviews(ServiceOptions{
		UserRole:  &rbacfakes.RoleAssignment{},
		Logger:    &logfakes.Logger{},
		Namespace: "test",
		Detailed:  true,
	})

	assert.Nil(t, err)
	assert.NotNil(t, r[0].Spec)
	assert.NotNil(t, r[0].ConfigMaps)
	assert.Len(t, *r[0].ConfigMaps, 1)
	assert.Len(t, r[0].DeploymentOverviews, 1)
}
