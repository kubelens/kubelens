package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrace(t *testing.T) {
	p := func() {
		Trace()
	}
	assert.NotPanics(t, p)
}
