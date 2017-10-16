// Copyright 2017 The go-rollbar Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rollbar

import (
	"fmt"
	"hash/adler32"
	"reflect"
	"strings"

	api "github.com/zchee/go-rollbar/api/v1"
)

// errorBody creates a rollbar error body with a given stack trace.
func errorBody(err error, stack Stack) *api.Body {
	message := "<nil>"
	if err != nil {
		message = err.Error()
	}

	return &api.Body{
		Trace: &api.Trace{
			Frames: stack,
			Exception: &api.Exception{
				Class:   errorClass(err),
				Message: message,
			},
		},
	}
}

// errorClass expands the function(class) name from err.
func errorClass(err error) string {
	if err == nil {
		return "<nil>"
	}

	fn := reflect.TypeOf(err).String()
	switch fn {
	case "":
		return "panic"
	case "*errors.errorString":
		checksum := adler32.Checksum([]byte(err.Error()))
		return fmt.Sprintf("{%x}", checksum)
	default:
		return strings.TrimPrefix(fn, "*")
	}
}

// Level level of stack trace.
type Level string

const (
	// DebugLevel logs are typically voluminous, and are usually disabled in production.
	DebugLevel Level = "debug"
	// InfoLevel is the default logging priority.
	InfoLevel Level = "info"
	// WarnLevel logs are more important than Info, but don't need individual human review.
	WarnLevel Level = "warning"
	// ErrorLevel logs are high-priority. If an application is running smoothly, it shouldn't generate any error-level logs.
	ErrorLevel Level = "error"
	// CriticalLevel logs are particularly important errors. In development the logger panics after writing the message.
	CriticalLevel Level = "critical"
)
