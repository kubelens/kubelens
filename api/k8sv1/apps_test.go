package k8sv1

import (
	"testing"

	rbacfakes "github.com/kubelens/kubelens/api/auth/fakes"
	logfakes "github.com/kubelens/kubelens/api/log/fakes"
	"github.com/stretchr/testify/assert"
)

func TestAppsDefault(t *testing.T) {
	c := setupClient("default", "test", false, false)

	r, err := c.Apps(AppOptions{
		UserRole: rbacfakes.RoleAssignment{},
		Logger:   &logfakes.Logger{},
	})

	assert.Nil(t, err)
	assert.True(t, len(r) > 0)
}

func TestAppOverviewDefault(t *testing.T) {
	c := setupClient("default", "test", false, false)

	r, err := c.AppOverview(AppOverviewOptions{
		UserRole:      rbacfakes.RoleAssignment{},
		Logger:        &logfakes.Logger{},
		LabelSelector: map[string]string{"app": "test"},
		Namespace:     "default",
		AppName:       "test",
	})

	assert.Nil(t, err)
	assert.True(t, len(r.PodOverviews.Name) > 0)
}
