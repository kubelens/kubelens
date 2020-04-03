package k8v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	rbacfakes "github.com/kubelens/kubelens/api/auth/fakes"
	logfakes "github.com/kubelens/kubelens/api/log/fakes"
)

func TestGetLogsDefault(t *testing.T) {
	ns := "fake"
	n := "test"

	c := setupClient(ns, n, false, false)

	p := func() {
		c.Logs(LogOptions{
			UserRole:  &rbacfakes.RoleAssignment{},
			Logger:    &logfakes.Logger{},
			Namespace: ns,
			PodName:   n,
			Tail:      100,
			Follow:    false,
		})
	}

	assert.Panics(t, p)
}

func TestGetLogsMissingName(t *testing.T) {
	ns := "fake"
	n := "test"

	c := setupClient(ns, n, false, false)

	_, err := c.Logs(LogOptions{
		UserRole:  &rbacfakes.RoleAssignment{},
		Logger:    &logfakes.Logger{},
		Namespace: ns,
		PodName:   "",
		Tail:      100,
		Follow:    false,
	})

	assert.Equal(t, "\nBad Request: podname must be provided when getting logs\n", err.Message)
	assert.Equal(t, 400, err.Code)
}

func TestGetLogsMissingNamespace(t *testing.T) {
	ns := "fake"
	n := "test"

	c := setupClient(ns, n, false, false)

	_, err := c.Logs(LogOptions{
		UserRole:  &rbacfakes.RoleAssignment{},
		Logger:    &logfakes.Logger{},
		Namespace: "",
		PodName:   n,
		Tail:      100,
		Follow:    false,
	})

	assert.Equal(t, "\nBad Request: namespace must be provided when getting logs\n", err.Message)
	assert.Equal(t, 400, err.Code)
}

func TestGetLogsClientError(t *testing.T) {
	ns := "fake"
	n := "test"

	c := setupClient(ns, n, true, false)

	_, err := c.Logs(LogOptions{
		UserRole:  &rbacfakes.RoleAssignment{},
		Logger:    &logfakes.Logger{},
		Namespace: ns,
		PodName:   "test-123",
		Tail:      100,
		Follow:    false,
	})

	assert.Equal(t, "\nInternal Server Error: GetClientSet Test Error\n", err.Message)
	assert.Equal(t, 500, err.Code)
}

// func TestGetLogsGetListError(t *testing.T) {
// 	ns := "fake"
// 	n := "test"

// 	c := setupClient(ns, n, false, true)

// 	_, err := c.Logs(LogOptions{
// 		UserRole:  &rbacfakes.RoleAssignment{},
// 		Logger:    &logfakes.Logger{},
// 		Namespace: ns,
// 		PodName:   "test-123",
// 		Tail:      100,
// 		Follow:    false,
// 	})

// 	assert.Equal(t, "\nInternal Server Error: GetClientSet Test Error\n", err.Message)
// 	assert.Equal(t, 500, err.Code)
// }
