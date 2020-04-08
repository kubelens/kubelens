/*
MIT License

Copyright (c) 2020 The KubeLens Authors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package log

import (
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Logger is an interface for logger
type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

// logger wraps logrus.Logger so that it can log messages sharing a common set of fields.
type logger struct {
	log    *logrus.Logger
	fields logrus.Fields
}

// New creates a logger instance with the specified logrus.Logger and the fields that should be added to every message.
func New(l *logrus.Logger, appname string) Logger {
	fields := logrus.Fields{
		"timestamp":  time.Now(),
		"request_id": uuid.New().String(),
		"appname":    appname,
	}
	return &logger{
		l,
		fields,
	}
}

func (l *logger) Debugf(format string, args ...interface{}) {
	go func(fields logrus.Fields, format string, args ...interface{}) {
		l.log.WithFields(l.fields).Debugf(format, args...)
	}(l.fields, format, args)
}

func (l *logger) Infof(format string, args ...interface{}) {
	go func(fields logrus.Fields, format string, args ...interface{}) {
		l.log.WithFields(l.fields).Infof(format, args...)
	}(l.fields, format, args)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	go func(fields logrus.Fields, format string, args ...interface{}) {
		l.log.WithFields(l.fields).Warnf(format, args...)
	}(l.fields, format, args)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	go func(fields logrus.Fields, format string, args ...interface{}) {
		l.log.WithFields(l.fields).Error(args...)
	}(l.fields, format, args)
}

func (l *logger) Debug(args ...interface{}) {
	go func(fields logrus.Fields, args ...interface{}) {
		l.log.WithFields(l.fields).Debug(args...)
	}(l.fields, args)
}

func (l *logger) Info(args ...interface{}) {
	go func(fields logrus.Fields, args ...interface{}) {
		l.log.WithFields(l.fields).Info(args...)
	}(l.fields, args)
}

func (l *logger) Warn(args ...interface{}) {
	go func(fields logrus.Fields, args ...interface{}) {
		l.log.WithFields(l.fields).Warn(args...)
	}(l.fields, args)
}

func (l *logger) Error(args ...interface{}) {
	go func(fields logrus.Fields, args ...interface{}) {
		l.log.WithFields(l.fields).Error(args...)
	}(l.fields, args)
}
