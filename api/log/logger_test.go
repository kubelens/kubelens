package log

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	l := New(&logrus.Logger{}, "testapp")

	p := func() {
		l.Debug("test log")
		l.Info("test log")
		l.Warn("test log")
		l.Error("test log")
		l.Debugf("formated %s", "test log")
		l.Infof("formated %s", "test log")
		l.Warnf("formated %s", "test log")
		l.Errorf("formated %s", "test log")
	}

	assert.NotPanics(t, p)
}
