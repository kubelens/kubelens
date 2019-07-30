package io

import (
	"net/http/httptest"
	"testing"

	"github.com/kubelens/kubelens/api/config"
	"github.com/stretchr/testify/assert"
)

func TestUpgraderCheckOriginNotAllowed(t *testing.T) {
	r := httptest.NewRequest("GET", "/io/test", nil)

	allowed := upgrader.CheckOrigin(r)

	assert.False(t, allowed)
}

func TestUpgraderCheckOriginAllowed(t *testing.T) {
	config.Set("../config/config.json")
	r := httptest.NewRequest("GET", "/io/test", nil)

	r.Host = "api.kubelens.local"

	allowed := upgrader.CheckOrigin(r)

	assert.True(t, allowed)
}
