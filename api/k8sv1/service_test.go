package k8sv1

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
		LabelSelector: map[string]string{"app": "test"},
	})

	assert.Nil(t, err)
}

func TestGetServiceOverviewsDefaultMissingNamespace(t *testing.T) {
	c := setupClient("test", "test", false, false)

	_, err := c.ServiceOverviews(ServiceOptions{
		UserRole:      &rbacfakes.RoleAssignment{},
		Logger:        &logfakes.Logger{},
		Namespace:     "",
		LabelSelector: map[string]string{"app": "test"},
	})

	assert.Nil(t, err)
}

func TestGetServiceOverviewsDefaultMissingLabelSearch(t *testing.T) {
	c := setupClient("test", "test", false, false)

	_, err := c.ServiceOverviews(ServiceOptions{
		UserRole:      &rbacfakes.RoleAssignment{},
		Logger:        &logfakes.Logger{},
		Namespace:     "fake",
		LabelSelector: map[string]string{},
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
	assert.True(t, len(r) > 0)
}

func TestServiceAppInfosDefaultSuccess(t *testing.T) {
	c := setupClient("testns", "test", false, false)

	ls := map[string]string{}
	ls[AppNameLabel] = FriendlyAppName

	d, err := c.ServiceAppInfos(ServiceOptions{
		UserRole:  &rbacfakes.RoleAssignment{},
		Logger:    &logfakes.Logger{},
		Namespace: "testns",
	})

	assert.Nil(t, err)
	assert.Len(t, d, 1)
	assert.Equal(t, d[0].Namespace, "testns")
}

func TestGetServiceAppInfosDefaultFail(t *testing.T) {
	c := setupClient("testns", "test", true, true)

	_, err := c.ServiceAppInfos(ServiceOptions{
		UserRole:  &rbacfakes.RoleAssignment{},
		Logger:    &logfakes.Logger{},
		Namespace: "testns",
	})

	assert.NotNil(t, err)
}
