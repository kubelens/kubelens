package k8sv1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateLabelSelector(t *testing.T) {
	selector := generateLabelSelector("value")

	assert.Equal(t, "app=value", selector)
}

func TestStringContainsSensitiveInfo(t *testing.T) {
	match := stringContainsSensitiveInfo("dbPassword")
	assert.True(t, match)

	match = stringContainsSensitiveInfo("apiKey")
	assert.True(t, match)

	match = stringContainsSensitiveInfo("someSecret")
	assert.True(t, match)

	match = stringContainsSensitiveInfo("random_thing")
	assert.False(t, match)
}
