package k8sv1

import (
	"context"
	"testing"

	logfakes "github.com/kubelens/kubelens/api/log/fakes"
	"github.com/stretchr/testify/assert"
)

func TestJobsDefaultSuccess(t *testing.T) {
	c := setupClient("testns", "jobs1", false, false)

	d, err := c.Jobs(JobOptions{
		Logger:     &logfakes.Logger{},
		Namespace:  "testns",
		LinkedName: "jobs1",
		Context:    context.Background(),
	})

	assert.Nil(t, err)
	assert.Len(t, d, 1)
	assert.Equal(t, "testns", d[0].Namespace)
}

func TestGetJobsDefaultFail(t *testing.T) {
	c := setupClient("testns", "jobs2", true, true)

	_, err := c.Jobs(JobOptions{
		Logger:     &logfakes.Logger{},
		Namespace:  "testns",
		LinkedName: "jobs2",
		Context:    context.Background(),
	})

	assert.NotNil(t, err)
}

func TestJobDefaultSuccess(t *testing.T) {
	c := setupClient("testns", "jobs3", false, false)

	d, err := c.Job(JobOptions{
		Logger:     &logfakes.Logger{},
		Namespace:  "testns",
		Name:       "jobs3",
		LinkedName: "whatever",
		Context:    context.Background(),
	})

	assert.Nil(t, err)
	assert.NotNil(t, d)
	assert.Equal(t, d.Namespace, "testns")
}

func TestGetJobDefaultFail(t *testing.T) {
	c := setupClient("testns", "test", true, true)

	_, err := c.Job(JobOptions{
		Logger:    &logfakes.Logger{},
		Namespace: "testns",
		Context:   context.Background(),
	})

	assert.NotNil(t, err)
}
