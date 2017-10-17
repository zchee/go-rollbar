// Copyright 2017 The go-rollbar Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rollbar

const (
	Name    = "go-rollbar"
	Version = "0.0.0"

	language = "go"
)

type Rollbar interface {
	Debug(error) Call
	Info(error) Call
	Error(error) Call
	Warn(error) Call
	Critical(error) Call
}
