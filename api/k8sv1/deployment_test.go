package k8v1

import (
	"testing"

	rbacfakes "github.com/kubelens/kubelens/api/auth/fakes"
	logfakes "github.com/kubelens/kubelens/api/log/fakes"
	"github.com/stretchr/testify/assert"
)

func TestGetDeploymentOverviewsDefaultSuccess(t *testing.T) {
	c := setupClient("test", "test", false, false)

	d, err := c.DeploymentOverviews(DeploymentOptions{
		UserRole:  &rbacfakes.RoleAssignment{},
		Logger:    &logfakes.Logger{},
		Namespace: "test",
		// just use the config map labelselctor for the service for ease.
		// set in wrapper_test.go
		LabelSelector: "cmtest=cmvalue",
	})

	assert.Nil(t, err)
	assert.Len(t, d, 1)
}
