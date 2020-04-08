package k8sv1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanityCheckFail(t *testing.T) {
	c := setupClient("default", "test", false, false)

	success := c.SanityCheck()

	assert.True(t, success)
}

func TestSanityCheck(t *testing.T) {
	c := setupClient("default", "test", true, false)

	success := c.SanityCheck()

	assert.False(t, success)
}

func TestSanityCheckFailedGetServices(t *testing.T) {
	c := setupClient("default", "test", false, true)

	success := c.SanityCheck()

	assert.False(t, success)
}
