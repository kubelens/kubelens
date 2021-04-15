package k8sv1

import (
	"context"
	"testing"

	rbacfakes "github.com/kubelens/kubelens/api/auth/fakes"
	logfakes "github.com/kubelens/kubelens/api/log/fakes"
	"github.com/stretchr/testify/assert"
)

func TestOverviewsSuccess(t *testing.T) {
	c := setupClient("testns", "ov1", false, false)

	d, err := c.Overviews(OverviewOptions{
		UserRole: &rbacfakes.RoleAssignment{},
		Logger:   &logfakes.Logger{},
		Context:  context.Background(),
	})

	assert.Nil(t, err)
	assert.True(t, len(d) > 0)
	assert.Equal(t, "testns", d[0].Namespace)
}

func TestOverviewSuccess(t *testing.T) {
	c := setupClient("testns", "ov2", false, false)

	d, err := c.Overview(OverviewOptions{
		Namespace:  "testns",
		LinkedName: "ov2",
		UserRole:   &rbacfakes.RoleAssignment{},
		Logger:     &logfakes.Logger{},
		Context:    context.Background(),
	})

	assert.Nil(t, err)
	assert.NotNil(t, d)
}
