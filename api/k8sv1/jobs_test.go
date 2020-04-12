package k8sv1

import (
	"testing"

	rbacfakes "github.com/kubelens/kubelens/api/auth/fakes"
	logfakes "github.com/kubelens/kubelens/api/log/fakes"
	"github.com/stretchr/testify/assert"
)

func TestJobOverviewsDefaultSuccess(t *testing.T) {
	c := setupClient("testns", "test", false, false)

	ls := map[string]string{}
	ls[AppNameLabel] = FriendlyAppName

	d, err := c.JobOverviews(JobOptions{
		UserRole:      &rbacfakes.RoleAssignment{},
		Logger:        &logfakes.Logger{},
		Namespace:     "testns",
		LabelSelector: ls,
	})

	assert.Nil(t, err)
	assert.Len(t, d, 1)
	assert.Equal(t, "testns", d[0].Namespace)
}

func TestGetJobOverviewsDefaultFail(t *testing.T) {
	c := setupClient("testns", "test", true, true)

	_, err := c.DeploymentOverviews(DeploymentOptions{
		UserRole:      &rbacfakes.RoleAssignment{},
		Logger:        &logfakes.Logger{},
		Namespace:     "testns",
		LabelSelector: map[string]string{"random": "labelvalue"},
	})

	assert.NotNil(t, err)
}

func TestJobAppInfosDefaultSuccess(t *testing.T) {
	c := setupClient("testns", "test", false, false)

	ls := map[string]string{}
	ls[AppNameLabel] = FriendlyAppName

	d, err := c.JobAppInfos(JobOptions{
		UserRole:  &rbacfakes.RoleAssignment{},
		Logger:    &logfakes.Logger{},
		Namespace: "testns",
	})

	assert.Nil(t, err)
	assert.Len(t, d, 1)
	assert.Equal(t, d[0].Namespace, "testns")
}

func TestGetJobAppInfosDefaultFail(t *testing.T) {
	c := setupClient("testns", "test", true, true)

	_, err := c.JobAppInfos(JobOptions{
		UserRole:  &rbacfakes.RoleAssignment{},
		Logger:    &logfakes.Logger{},
		Namespace: "testns",
	})

	assert.NotNil(t, err)
}
