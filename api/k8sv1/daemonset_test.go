package k8sv1

import (
	"testing"

	rbacfakes "github.com/kubelens/kubelens/api/auth/fakes"
	logfakes "github.com/kubelens/kubelens/api/log/fakes"
	"github.com/stretchr/testify/assert"
)

func TestDaemonSetOverviewsDefaultSuccess(t *testing.T) {
	c := setupClient("testns", "test", false, false)

	ls := map[string]string{}
	ls[AppNameLabel] = FriendlyAppName

	d, err := c.DaemonSetOverviews(DaemonSetOptions{
		UserRole:      &rbacfakes.RoleAssignment{},
		Logger:        &logfakes.Logger{},
		Namespace:     "testns",
		LabelSelector: ls,
	})

	assert.Nil(t, err)
	assert.Len(t, d, 1)
	assert.Equal(t, d[0].Namespace, "testns")
}

func TestGetDaemonSetOverviewsDefaultFail(t *testing.T) {
	c := setupClient("testns", "test", true, true)

	_, err := c.DeploymentOverviews(DeploymentOptions{
		UserRole:      &rbacfakes.RoleAssignment{},
		Logger:        &logfakes.Logger{},
		Namespace:     "testns",
		LabelSelector: map[string]string{"random": "labelvalue"},
	})

	assert.NotNil(t, err)
}

func TestDaemonSetAppInfosDefaultSuccess(t *testing.T) {
	c := setupClient("testns", "test", false, false)

	ls := map[string]string{}
	ls[AppNameLabel] = FriendlyAppName

	d, err := c.DaemonSetAppInfos(DaemonSetOptions{
		UserRole:  &rbacfakes.RoleAssignment{},
		Logger:    &logfakes.Logger{},
		Namespace: "testns",
	})

	assert.Nil(t, err)
	assert.Len(t, d, 1)
	assert.Equal(t, d[0].Namespace, "testns")
}

func TestGetDaemonSetAppInfosDefaultFail(t *testing.T) {
	c := setupClient("testns", "test", true, true)

	_, err := c.DaemonSetAppInfos(DaemonSetOptions{
		UserRole:  &rbacfakes.RoleAssignment{},
		Logger:    &logfakes.Logger{},
		Namespace: "testns",
	})

	assert.NotNil(t, err)
}
