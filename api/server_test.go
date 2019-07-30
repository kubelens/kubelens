package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateServer(t *testing.T) {
	hs := createServer(nil)

	assert.Equal(t, ":39000", hs.Addr)
}
