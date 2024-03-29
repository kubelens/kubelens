package k8sv1

import (
	"context"
	"net/http"
	"testing"

	logfakes "github.com/kubelens/kubelens/api/log/fakes"
	"github.com/stretchr/testify/assert"
)

func TestPodDefault(t *testing.T) {
	c := setupClient("default", "pod1", false, false)

	r, err := c.Pod(PodOptions{
		Logger:    &logfakes.Logger{},
		Name:      "pod1",
		Namespace: "default",
		Context:   context.Background(),
	})

	assert.Nil(t, err)
	assert.Equal(t, "pod1", r.Name)
}

func TestPodDetailForbidden(t *testing.T) {
	c := setupClient("default", "pod2", true, false)

	_, err := c.Pod(PodOptions{
		Logger:    &logfakes.Logger{},
		Name:      "pod2",
		Namespace: "default",
		Context:   context.Background(),
	})

	assert.Equal(t, http.StatusInternalServerError, err.Code)
	assert.Equal(t, "\nInternal Server Error: GetClientSet Test Error\n", err.Message)
}

func TestPodsDefault(t *testing.T) {
	c := setupClient("default", "pod3", false, false)

	r, err := c.Pods(PodOptions{
		Logger:     &logfakes.Logger{},
		Namespace:  "default",
		LinkedName: "pod3",
		Context:    context.Background(),
	})

	assert.Nil(t, err)
	assert.True(t, len(r) > 0)
}

func TestPodsDefaultWithFilters(t *testing.T) {
	c := setupClient("default", "pod4", false, false)

	r, err := c.Pods(PodOptions{
		Logger:     &logfakes.Logger{},
		Namespace:  "default",
		Name:       "test",
		LinkedName: "pod4",
		Context:    context.Background(),
	})

	assert.Nil(t, err)
	assert.True(t, len(r) > 0)
}
