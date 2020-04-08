package k8sv1

import (
	"testing"

	rbacfakes "github.com/kubelens/kubelens/api/auth/fakes"
	logfakes "github.com/kubelens/kubelens/api/log/fakes"
	"github.com/stretchr/testify/assert"
)

func TestGetDeploymentOverviewsDefaultSuccess(t *testing.T) {
	c := setupClient("testns", "cmvalue", false, false)

	ls := map[string]string{}
	ls[AppNameLabel] = FriendlyAppName

	d, err := c.DeploymentOverviews(DeploymentOptions{
		UserRole:  &rbacfakes.RoleAssignment{},
		Logger:    &logfakes.Logger{},
		Namespace: "testns",
		// just use the config map labelselctor for the service for ease.
		// set in wrapper_test.go
		LabelSelector: ls,
	})

	assert.Nil(t, err)
	assert.Len(t, d, 1)
	assert.Equal(t, d[0].Namespace, "testns")
}

func TestGetDeploymentOverviewsDefaultFail(t *testing.T) {
	c := setupClient("testns", "test", true, true)

	_, err := c.DeploymentOverviews(DeploymentOptions{
		UserRole:      &rbacfakes.RoleAssignment{},
		Logger:        &logfakes.Logger{},
		Namespace:     "testns",
		LabelSelector: map[string]string{"random": "labelvalue"},
	})

	assert.NotNil(t, err)
}
