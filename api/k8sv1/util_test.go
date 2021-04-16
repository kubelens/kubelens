package k8sv1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLabelsContainSelectorTrue(t *testing.T) {
	labels := map[string]string{
		"app": "app-name",
	}

	match := labelsContainSelector("app-name", labels)

	assert.True(t, match)
}

func TestLabelsContainSelectorFalse(t *testing.T) {
	labels := map[string]string{
		"app": "app-name",
	}

	match := labelsContainSelector("app-name2", labels)

	assert.False(t, match)
}

func TestLabelsContainSelectorNilLabels(t *testing.T) {
	match := labelsContainSelector("app-name", nil)

	assert.False(t, match)
}

func TestLabelsContainSelectorNoSelector(t *testing.T) {
	labels := map[string]string{
		"app":          "app-name",
		"pod-template": "1234adsf1B",
	}

	match := labelsContainSelector("", labels)

	assert.False(t, match)
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
