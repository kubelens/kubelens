package k8sv1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanityCheckFail(t *testing.T) {
	c := setupClient("default", "test", false, false)

	err := c.SanityCheck()

	assert.Nil(t, err)
}

func TestSanityCheck(t *testing.T) {
	c := setupClient("default", "test", true, false)

	err := c.SanityCheck()

	assert.NotNil(t, err)
}

func TestSanityCheckFailedGetServices(t *testing.T) {
	c := setupClient("default", "test", false, true)

	err := c.SanityCheck()

	assert.NotNil(t, err)
}
