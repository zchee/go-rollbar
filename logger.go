// Copyright 2017 The go-rollbar Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rollbar

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
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
	Fatal(context.Context, string, ...interface{})
}

type nilLogger struct{}

func (_ nilLogger) Debug(context.Context, string, ...interface{}) {}
func (_ nilLogger) Info(context.Context, string, ...interface{})  {}
func (_ nilLogger) Fatal(context.Context, string, ...interface{}) {}

type traceLogger struct {
	w io.Writer
}

// Debug outputs the debug log output.
func (l traceLogger) Debug(_ context.Context, f string, args ...interface{}) {
	fmt.Fprintf(l.w, spew.Sprintf(newLine(f), args...))
}

// Info outputs the info log output.
func (l traceLogger) Info(_ context.Context, f string, args ...interface{}) {
	fmt.Fprintf(l.w, newLine(f), args...)
}

// Fatal outputs the fatal log output, and exit 1.
func (l traceLogger) Fatal(_ context.Context, f string, args ...interface{}) {
	fmt.Fprintf(l.w, newLine(f), args...)
	os.Exit(1)
}

// newLine joins the new line if not present.
func newLine(s string) string {
	if !strings.HasSuffix(s, "\n") {
		s = s + "\n"
	}
	return s
}
