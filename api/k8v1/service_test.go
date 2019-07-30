package k8v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	rbacfakes "github.com/kubelens/kubelens/api/auth/fakes"
	logfakes "github.com/kubelens/kubelens/api/log/fakes"
)

func TestGetServiceOverviewsDefaultSuccess(t *testing.T) {
	c := setupClient("test", "test", false, false)

	_, err := c.ServiceOverviews(ServiceOptions{
		UserRole:    &rbacfakes.RoleAssignment{},
		Logger:      &logfakes.Logger{},
		Namespace:   "test",
		LabelSearch: "app=test",
	})

	assert.Nil(t, err)
}

func TestGetServiceOverviewsDefaultMissingNamespace(t *testing.T) {
	c := setupClient("test", "test", false, false)

	_, err := c.ServiceOverviews(ServiceOptions{
		UserRole:    &rbacfakes.RoleAssignment{},
		Logger:      &logfakes.Logger{},
		Namespace:   "",
		LabelSearch: "app=test",
	})

	assert.Nil(t, err)
}

func TestGetServiceOverviewsDefaultMissingLabelSearch(t *testing.T) {
	c := setupClient("test", "test", false, false)

	_, err := c.ServiceOverviews(ServiceOptions{
		UserRole:    &rbacfakes.RoleAssignment{},
		Logger:      &logfakes.Logger{},
		Namespace:   "fake",
		LabelSearch: "",
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
}
