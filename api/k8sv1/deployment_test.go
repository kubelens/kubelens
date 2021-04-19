package k8sv1

import (
	"context"
	"testing"

	logfakes "github.com/kubelens/kubelens/api/log/fakes"
	"github.com/stretchr/testify/assert"
)

func TestGetDeploymentsDefaultSuccess(t *testing.T) {
	c := setupClient("testns", "dpl1", false, false)

	d, err := c.Deployments(DeploymentOptions{
		Logger:    &logfakes.Logger{},
		Namespace: "testns",
		// just use the config map labelselctor for the service for ease.
		// set in wrapper_test.go
		LinkedName: "dpl1",
		Context:    context.Background(),
	})

	assert.Nil(t, err)
	assert.Len(t, d, 1)
	assert.Equal(t, d[0].Namespace, "testns")
}

func TestGetDeploymentDefaultFail(t *testing.T) {
	c := setupClient("testns", "dpl2", true, true)

	_, err := c.Deployments(DeploymentOptions{
		Logger:     &logfakes.Logger{},
		Namespace:  "testns",
		LinkedName: "dpl2",
		Context:    context.Background(),
	})

	assert.NotNil(t, err)
}
