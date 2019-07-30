package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	Set("config.json")

	assert.NotNil(t, C)
}

func TestSetNoFilePanics(t *testing.T) {
	p := func() {
		Set("fake")
	}

	assert.Panics(t, p)
}
