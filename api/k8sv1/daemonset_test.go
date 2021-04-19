package k8sv1

import (
	"context"
	"testing"

	logfakes "github.com/kubelens/kubelens/api/log/fakes"
	"github.com/stretchr/testify/assert"
)

func TestDaemonSetsDefaultSuccess(t *testing.T) {
	c := setupClient("testns", "dstest1", false, false)

	d, err := c.DaemonSets(DaemonSetOptions{
		Logger:     &logfakes.Logger{},
		Namespace:  "testns",
		LinkedName: "dstest1",
		Context:    context.Background(),
	})

	assert.Nil(t, err)
	assert.Len(t, d, 1)
	assert.Equal(t, d[0].Namespace, "testns")
}

func TestGetDaemonSetsDefaultFail(t *testing.T) {
	c := setupClient("testns", "dstest2", true, true)

	_, err := c.DaemonSets(DaemonSetOptions{
		Logger:     &logfakes.Logger{},
		Namespace:  "testns",
		LinkedName: "dstest2",
		Context:    context.Background(),
	})

	assert.NotNil(t, err)
}

func TestDaemonSetDefaultSuccess(t *testing.T) {
	c := setupClient("testns", "dstest3", false, false)

	d, err := c.DaemonSet(DaemonSetOptions{
		Logger:     &logfakes.Logger{},
		Namespace:  "testns",
		Name:       "dstest3",
		LinkedName: "whatever",
		Context:    context.Background(),
	})

	assert.Nil(t, err)
	assert.NotNil(t, d)
	assert.Equal(t, d.Namespace, "testns")
}

func TestGetDaemonSetDefaultFail(t *testing.T) {
	c := setupClient("testns", "dstest4", true, true)

	_, err := c.DaemonSet(DaemonSetOptions{
		Logger:     &logfakes.Logger{},
		Namespace:  "testns",
		Name:       "dstest4",
		LinkedName: "dstest4",
		Context:    context.Background(),
	})

	assert.NotNil(t, err)
}
