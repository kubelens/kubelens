package k8v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	rbacfakes "github.com/kubelens/kubelens/api/auth/fakes"
	logfakes "github.com/kubelens/kubelens/api/log/fakes"
)

func TestAppOverviewDefault(t *testing.T) {
	c := setupClient("default", "test", false, false)

	r, err := c.AppOverview(AppOverviewOptions{
		UserRole:  rbacfakes.RoleAssignment{},
		Logger:    &logfakes.Logger{},
		LabelKey:  "app",
		Namespace: "default",
		AppName:   "test",
	})

	assert.Nil(t, err)
	assert.True(t, len(r.PodOverviews.Name.Value) > 0)
}
