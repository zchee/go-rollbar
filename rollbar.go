// Copyright 2017 The go-rollbar Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rollbar

import (
	"golang.org/x/net/context"
)

const (
	Name    = "go-rollbar"
	Version = "0.0.0"

	language = "go"
)

type Rollbar interface {
	Debug(context.Context, error, ...ErrorOption)
	Info(context.Context, error, ...ErrorOption)
	Error(context.Context, error, ...ErrorOption)
	Warn(context.Context, error, ...ErrorOption)
	Critical(context.Context, error, ...ErrorOption)
}
