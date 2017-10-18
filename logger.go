// Copyright 2017 The go-rollbar Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rollbar

import (
	"fmt"
	"io"

	"golang.org/x/net/context"
)

// Logger is an interface for logging/tracing the client's
// execution.
//
// In particular, `Debug` will only be called if `WithDebug`
// is provided to the constructor.
type Logger interface {
	Debugf(context.Context, string, ...interface{})
	Infof(context.Context, string, ...interface{})
}

type nilLogger struct{}

func (nilLogger) Debugf(context.Context, string, ...interface{}) {}
func (nilLogger) Infof(context.Context, string, ...interface{})  {}

type traceLogger struct {
	w io.Writer
}

// Debug outputs the debug log output.
func (l traceLogger) Debugf(_ context.Context, format string, args ...interface{}) {
	fmt.Fprintf(l.w, format, args...)
}

// Info outputs the info log output.
func (l traceLogger) Infof(_ context.Context, format string, args ...interface{}) {
	fmt.Fprintf(l.w, format, args...)
}
