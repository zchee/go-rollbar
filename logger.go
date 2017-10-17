// Copyright 2017 The go-rollbar Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rollbar

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/context"
)

// Logger is an interface for logging/tracing the client's
// execution.
//
// In particular, `Debug` will only be called if `WithDebug`
// is provided to the constructor.
type Logger interface {
	Debug(context.Context, string, ...interface{})
	Info(context.Context, string, ...interface{})
}

type nilLogger struct{}

func (_ nilLogger) Debug(context.Context, string, ...interface{}) {}
func (_ nilLogger) Info(context.Context, string, ...interface{})  {}

type traceLogger struct {
	w io.Writer
}

// Debug outputs the debug log output.
func (l traceLogger) Debug(_ context.Context, f string, args ...interface{}) {
	fmt.Fprintf(l.w, newLine(f), args...)
}

// Info outputs the info log output.
func (l traceLogger) Info(_ context.Context, f string, args ...interface{}) {
	fmt.Fprintf(l.w, newLine(f), args...)
}

// newLine joins the new line if not present.
func newLine(s string) string {
	if !strings.HasSuffix(s, "\n") {
		s = s + "\n"
	}
	return s
}
