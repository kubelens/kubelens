package k8sv1

import (
	"context"
	"testing"

	rbacfakes "github.com/kubelens/kubelens/api/auth/fakes"
	logfakes "github.com/kubelens/kubelens/api/log/fakes"
	"github.com/stretchr/testify/assert"
)

func TestGetServicesDefaultSuccess(t *testing.T) {
	c := setupClient("test", "svc1", false, false)

	s, err := c.Services(ServiceOptions{
		UserRole:   &rbacfakes.RoleAssignment{},
		Logger:     &logfakes.Logger{},
		Namespace:  "test",
		LinkedName: "svc1",
		Labels:     map[string]string{"app": "test"},
		Context:    context.Background(),
	})

	assert.Nil(t, err)
	assert.NotNil(t, s)
	assert.True(t, len(s) > 0)
	assert.Equal(t, "svc1", s[0].LinkedName)
}

func TestGetServiceDefaultSuccess(t *testing.T) {
	c := setupClient("testns", "svc2", false, false)

	r, err := c.Service(ServiceOptions{
		UserRole:   &rbacfakes.RoleAssignment{},
		Logger:     &logfakes.Logger{},
		Namespace:  "testns",
		Name:       "svc2",
		LinkedName: "whatever",
		Context:    context.Background(),
	})

	assert.Nil(t, err)
	assert.NotNil(t, r)
	assert.Equal(t, r.Name, "svc2")
}

func TestGetServiceDefaultFail(t *testing.T) {
	c := setupClient("testns", "svc3", true, true)

	_, err := c.Service(ServiceOptions{
		UserRole:  &rbacfakes.RoleAssignment{},
		Logger:    &logfakes.Logger{},
		Namespace: "testns",
		Name:      "svc3",
		Context:   context.Background(),
	})

	assert.NotNil(t, err)
}
