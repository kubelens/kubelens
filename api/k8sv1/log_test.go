package k8sv1

import (
	"context"
	"testing"

	logfakes "github.com/kubelens/kubelens/api/log/fakes"
	"github.com/stretchr/testify/assert"
)

func TestGetLogsDefault(t *testing.T) {
	ns := "fake"
	n := "test"

	c := setupClient(ns, n, false, false)

	lgs, err := c.Logs(LogOptions{
		Logger:    &logfakes.Logger{},
		Namespace: ns,
		PodName:   n,
		Tail:      100,
		Follow:    false,
		Context:   context.Background(),
	})

	assert.Nil(t, err)
	assert.NotNil(t, lgs)
}

func TestGetLogsMissingName(t *testing.T) {
	ns := "fake"
	n := "test"

	c := setupClient(ns, n, false, false)

	_, err := c.Logs(LogOptions{
		Logger:    &logfakes.Logger{},
		Namespace: ns,
		PodName:   "",
		Tail:      100,
		Follow:    false,
		Context:   context.Background(),
	})

	assert.Equal(t, "\nBad Request: podname must be provided when getting logs\n", err.Message)
	assert.Equal(t, 400, err.Code)
}

func TestGetLogsMissingNamespace(t *testing.T) {
	ns := "fake"
	n := "test"

	c := setupClient(ns, n, false, false)

	_, err := c.Logs(LogOptions{
		Logger:    &logfakes.Logger{},
		Namespace: "",
		PodName:   n,
		Tail:      100,
		Follow:    false,
		Context:   context.Background(),
	})

	assert.Equal(t, "\nBad Request: namespace must be provided when getting logs\n", err.Message)
	assert.Equal(t, 400, err.Code)
}

func TestGetLogsClientError(t *testing.T) {
	ns := "fake"
	n := "test"

	c := setupClient(ns, n, true, false)

	_, err := c.Logs(LogOptions{
		Logger:    &logfakes.Logger{},
		Namespace: ns,
		PodName:   "test-123",
		Tail:      100,
		Follow:    false,
		Context:   context.Background(),
	})

	assert.Equal(t, "\nInternal Server Error: GetClientSet Test Error\n", err.Message)
	assert.Equal(t, 500, err.Code)
}
